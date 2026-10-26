/*
 * Copyright (c) 2026 The XGo Authors (xgo.dev). All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cl

import (
	"go/types"

	"github.com/goplus/llcppg/clang"
	lc "github.com/llarhub/clang-c"
)

// -----------------------------------------------------------------------------
//
// Token-level parsing of a forwarding macro body for automatic MacroFunc
// detection (see macrofunc_auto.go). Everything here works on the raw libclang
// tokens of the macro definition; no C preprocessor expansion is performed, so
// the rules stay syntactic and predictable.

// toks is a slice of libclang tokens paired with the translation unit needed to
// spell them. Every parse helper takes and returns toks so a sub-slice stays
// spellable.
type toks struct {
	tu   clang.TranslationUnit
	list []lc.Token
}

func (t toks) len() int                { return len(t.list) }
func (t toks) spell(i int) string      { return t.tu.TokenSpelling(t.list[i]) }
func (t toks) kind(i int) lc.TokenKind { return t.list[i].Kind() }
func (t toks) slice(lo, hi int) toks   { return toks{tu: t.tu, list: t.list[lo:hi]} }

// macroFuncBody returns the tokens that make up the body of a function-like
// macro definition, i.e. everything after the parameter list. The caller must
// invoke dispose when done to free the tokens. Comments are left in place and
// skipped by the parse helpers.
func macroFuncBody(decl clang.Cursor) (body toks, dispose func(), ok bool) {
	tu := clang.TU(decl)
	tokens, disp := tu.Tokenize(decl.Extent())
	dispose = disp

	// tokens[0] is the macro name; tokens[1] must be "(" for a function-like
	// macro. Walk to the matching ")" that closes the parameter list.
	if len(tokens) < 3 || tu.TokenSpelling(tokens[1]) != "(" {
		return
	}
	depth := 0
	i := 1
	for ; i < len(tokens); i++ {
		switch tu.TokenSpelling(tokens[i]) {
		case "(":
			depth++
		case ")":
			depth--
			if depth == 0 {
				i++ // first body token
				goto done
			}
		}
	}
	return
done:
	body = dropComments(toks{tu: tu, list: tokens[i:]})
	ok = body.len() > 0
	return
}

// dropComments returns t without any comment tokens.
func dropComments(t toks) toks {
	hasComment := false
	for i := range t.list {
		if t.kind(i) == lc.Token_Comment {
			hasComment = true
			break
		}
	}
	if !hasComment {
		return t
	}
	out := make([]lc.Token, 0, len(t.list))
	for i := range t.list {
		if t.kind(i) != lc.Token_Comment {
			out = append(out, t.list[i])
		}
	}
	return toks{tu: t.tu, list: out}
}

// -----------------------------------------------------------------------------

// parseForwardingCall checks that body is a single function call whose callee
// is a plain identifier (optionally wrapped in redundant parentheses), and
// returns the callee's C name together with the argument token groups (split at
// top-level commas).
//
// Accepted shapes (after stripping redundant outer parentheses around the whole
// body and around the callee identifier):
//
//	callee(arg0, arg1, ...)
//	(callee)(arg0, ...)
//
// A callee that is itself a call, a member access, a cast, a conditional, a
// comma expression, etc. is rejected.
func parseForwardingCall(body toks) (callee string, argGroups []toks, ok bool) {
	t := stripOuterParens(body)
	n := t.len()
	if n == 0 {
		return
	}

	i := 0
	// Optional "(" ident ")" around the callee identifier.
	if t.spell(0) == "(" {
		if n >= 3 && t.kind(1) == lc.Token_Identifier && t.spell(2) == ")" {
			callee = t.spell(1)
			i = 3
		} else {
			return
		}
	} else {
		if t.kind(0) != lc.Token_Identifier {
			return
		}
		callee = t.spell(0)
		i = 1
	}

	// Next token must be the call's opening "(".
	if i >= n || t.spell(i) != "(" {
		return
	}
	open := i
	// Find the matching ")" for the call; it must be the last token.
	depth := 0
	closeIdx := -1
	for j := open; j < n; j++ {
		switch t.spell(j) {
		case "(":
			depth++
		case ")":
			depth--
			if depth == 0 {
				closeIdx = j
				j = n
			}
		}
	}
	if closeIdx != n-1 {
		return // trailing tokens after the call, or unbalanced
	}

	inner := t.slice(open+1, closeIdx)
	if inner.len() > 0 {
		argGroups = splitTopLevel(inner)
	}
	ok = true
	return
}

// splitTopLevel splits a token slice on commas that are not nested inside
// parentheses/brackets/braces, returning each argument's tokens.
func splitTopLevel(t toks) []toks {
	var groups []toks
	depth := 0
	start := 0
	for j := 0; j < t.len(); j++ {
		switch t.spell(j) {
		case "(", "[", "{":
			depth++
		case ")", "]", "}":
			depth--
		case ",":
			if depth == 0 {
				groups = append(groups, t.slice(start, j))
				start = j + 1
			}
		}
	}
	groups = append(groups, t.slice(start, t.len()))
	return groups
}

// -----------------------------------------------------------------------------

// stripOuterParens removes redundant parentheses that wrap the entire token
// slice, e.g. "((x))" => "x". It only strips a leading "(" whose matching ")"
// is the last token.
func stripOuterParens(t toks) toks {
	for t.len() >= 2 && t.spell(0) == "(" {
		depth := 0
		matched := -1
		for j := 0; j < t.len(); j++ {
			switch t.spell(j) {
			case "(":
				depth++
			case ")":
				depth--
				if depth == 0 {
					matched = j
					j = t.len()
				}
			}
		}
		if matched != t.len()-1 {
			return t // the opening paren does not wrap the whole slice
		}
		t = t.slice(1, t.len()-1)
	}
	return t
}

// singleIdent reports whether t is exactly one identifier (optionally wrapped
// in redundant parentheses) and returns its spelling.
func singleIdent(t toks) (string, bool) {
	t = stripOuterParens(t)
	if t.len() == 1 && t.kind(0) == lc.Token_Identifier {
		return t.spell(0), true
	}
	return "", false
}

// parseSizeof recognizes a sizeof(T) expression, optionally wrapped in a scalar
// cast "(ScalarType)sizeof(T)", and returns the generated Go type for T plus the
// Go cast type (nil when there is no cast). It returns ok == false when the
// expression is not a plain sizeof over a type llcppg generates as a complete
// Go type.
func parseSizeof(ctx *pkgCtx, t toks) (szType, castTo types.Type, ok bool) {
	t = stripOuterParens(t)
	n := t.len()
	if n == 0 {
		return
	}

	// Optional leading scalar cast: "(" TypeName ")".
	if t.spell(0) == "(" {
		// find the matching ")"
		depth := 0
		close := -1
		for j := 0; j < n; j++ {
			switch t.spell(j) {
			case "(":
				depth++
			case ")":
				depth--
				if depth == 0 {
					close = j
					j = n
				}
			}
		}
		// A cast is "(Type) sizeof(...)" where the type tokens are between the
		// parens and "sizeof" immediately follows the closing paren.
		if close > 1 && close+1 < n && t.spell(close+1) == "sizeof" {
			castName := joinTypeTokens(t.slice(1, close))
			ct, found := scalarCastType(ctx, castName)
			if !found {
				return
			}
			castTo = ct
			t = t.slice(close+1, n)
			n = t.len()
		}
	}

	// Now t must be "sizeof" "(" TypeName ")".
	if n < 4 || t.spell(0) != "sizeof" || t.spell(1) != "(" || t.spell(n-1) != ")" {
		return
	}
	typeName := joinTypeTokens(t.slice(2, n-1))
	st, found := generatedCompleteType(ctx, typeName)
	if !found {
		return
	}
	szType = st
	ok = true
	return
}

// mentionsParam reports whether any identifier token in t is one of the macro
// parameters named in paramUse.
func mentionsParam(ctx *pkgCtx, t toks, paramUse map[string]int) bool {
	_ = ctx
	for i := 0; i < t.len(); i++ {
		if t.kind(i) == lc.Token_Identifier {
			if _, ok := paramUse[t.spell(i)]; ok {
				return true
			}
		}
	}
	return false
}

// evalConstTokens evaluates a constant expression given as a token slice.
func evalConstTokens(ctx *pkgCtx, t toks) (any, bool) {
	return evalConstExpr(ctx, t.tu, t.list)
}

// -----------------------------------------------------------------------------

// joinTypeTokens joins the identifier/keyword tokens of a type name into a
// single C type name (e.g. "unsigned int", "z_stream"). Pointer "*" and
// qualifiers are preserved with surrounding spaces collapsed.
func joinTypeTokens(t toks) string {
	s := ""
	for i := 0; i < t.len(); i++ {
		sp := t.spell(i)
		if s != "" && sp != "*" {
			s += " "
		}
		s += sp
	}
	return s
}

// scalarBasicKind maps a C scalar type name (as spelled in a cast) to a basic
// kind. Only the scalar spellings that can legitimately cast a sizeof result
// are listed.
var scalarBasicKind = map[string]basicKind{
	"char":               cChar,
	"int":                cInt,
	"unsigned":           cUint,
	"unsigned int":       cUint,
	"long":               cLong,
	"unsigned long":      cUlong,
	"long long":          cLongLong,
	"unsigned long long": cUlongLong,
	"float":              cFloat,
	"double":             cDouble,
	"long double":        cLongDouble,
}

// scalarCastType resolves a scalar cast type name (e.g. "int") to its Go type.
func scalarCastType(ctx *pkgCtx, name string) (types.Type, bool) {
	if k, ok := scalarBasicKind[name]; ok {
		return ctx.basicTyp(k), true
	}
	// A typedef'd scalar type is also acceptable as a cast target when it is a
	// type llcppg generates.
	if o, ok := ctx.getTypeObj(name, new(int)); ok {
		return o.Type(), true
	}
	return nil, false
}

// generatedCompleteType resolves a C type name used in sizeof(T) to the Go type
// llcppg generates for it, requiring that the type is a complete generated type
// (so unsafe.Sizeof over its Go form is meaningful). A type llcppg does not
// generate, and a record type llcppg emitted as an empty struct because its C
// definition was incomplete/opaque (unsafe.Sizeof would yield 0, not the real
// size), return ok == false.
func generatedCompleteType(ctx *pkgCtx, name string) (types.Type, bool) {
	// A scalar type name (e.g. sizeof(int)) resolves to a basic Go type.
	if k, ok := scalarBasicKind[name]; ok {
		return ctx.basicTyp(k), true
	}
	var feats int
	o, ok := ctx.getTypeObj(name, &feats)
	if !ok || o == nil {
		return nil, false
	}
	t := o.Type()
	// A record llcppg emits for an incomplete/opaque C type is a zero-field Go
	// struct whose size is 0, which does not match the C layout. Reject it so
	// the macro is skipped rather than generating a wrong unsafe.Sizeof.
	if s, ok := types.Unalias(t).Underlying().(*types.Struct); ok && s.NumFields() == 0 {
		return nil, false
	}
	return t, true
}

// -----------------------------------------------------------------------------
