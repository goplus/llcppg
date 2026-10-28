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

	"github.com/goplus/lib/c"
	"github.com/goplus/llcppg/clang"
	lc "github.com/llarhub/clang-c"
)

// -----------------------------------------------------------------------------

func loadMacroFunc(ctx *pkgCtx, decl clang.Cursor) {
	ctx.addPostUnit(func(ctx *pkgCtx) {
		compileMacroFunc(ctx, decl)
	})
}

// -----------------------------------------------------------------------------

func compileMacroFunc(ctx *pkgCtx, decl clang.Cursor) {
	macroName := clang.String(decl)

	tu := clang.TU(decl)
	tokens, dispose := tu.Tokenize(decl.Extent())
	defer dispose()

	paramNames, iNext, variadic, ok := macroParamNames(tu, tokens)
	if !ok || variadic {
		return
	}

	if refName, ok := ctx.macroFunc[macroName]; ok {
		doCompileMacroFunc(ctx, decl, macroName, refName, paramNames)
	} else {
		emitMacroFunc(ctx, decl, macroName, paramNames, tu, tokens[iNext:])
	}
}

func convMacroFuncCallArgs(ctx *pkgCtx, args []funcCallArg) bool {
	for i, arg := range args {
		switch arg.kind {
		case fca_IdentResolved:
			// Already resolved, nothing to do
		case fca_Ident:
			// Unresolved identifier, cannot convert
			return false
		case fca_Sizeof:
			feats := 0
			if o, ok := ctx.getTypeObj(arg.v.(string), &feats); ok {
				args[i].v = o.Type()
			} else {
				return false
			}
		}
	}
	return true
}

func doCompileAutoMacroFunc(ctx *pkgCtx, decl clang.Cursor, macroName string, paramNames []string, v funcCallExpr) {
	ref, ok := ctx.lookupFuncDecl(v.fn)
	if !ok {
		return
	}

	if fnRef, ok := ctx.compiledFns[objUSR(ref)]; ok {
		typName := ""
		sigCall := fnRef.Signature()
		recv := sigCall.Recv()
		if recv != nil {
			if len(paramNames) == 0 || !v.isMethodCall(paramNames[0]) {
				return
			}
			typName = typeNameOf(recv.Type())
			paramNames = paramNames[1:]
			v.args = v.args[1:]
		}
		callParams := sigCall.Params()
		if callParams.Len() != len(v.args) {
			return
		}
		pkg := ctx.pkg
		pkgTypes := pkg.Types
		node := goNode(ctx, decl)
		fnPos := node.Pos()
		params := make([]*types.Var, len(paramNames))
		for i, pname := range paramNames {
			iParam := v.checkParam(pname)
			if iParam < 0 {
				return
			}
			typ := callParams.At(iParam).Type()
			param := types.NewParam(fnPos, pkgTypes, avoidKeyword(pname), typ)
			v.args[iParam] = funcCallArg{fca_IdentResolved, param}
			params[i] = param
		}
		if !convMacroFuncCallArgs(ctx, v.args) {
			return
		}
		results := sigCall.Results()
		sig := types.NewSignatureType(recv, nil, nil, types.NewTuple(params...), results, false)
		fnName := ctx.funcName(macroName, -1, typName, typName, true, true)
		f, err := pkg.NewFuncWith(fnPos, fnName, sig, nil)
		if err != nil {
			ctx.errorf(decl, "%v", err)
			return
		}
		if doc := docCommentGroup(decl); doc != nil {
			f.SetComments(pkg, doc)
		}
		cb := f.BodyStart(pkg, node)
		if recv != nil {
			cb.Val(recv, node)
			cb.MemberVal(fnRef.Name(), 0, node)
		} else {
			cb.Val(fnRef.Func, node)
		}
		for i, arg := range v.args {
			if arg.kind == fca_Sizeof { // IntType( unsafe.Sizeof( T ) )
				intT := callParams.At(i).Type()
				t := arg.v.(types.Type)
				cb.Typ(intT, node).
					Val(ctx.unsafeSizeof(), node).ZeroLit(t).CallWith(1, 0, 0, node).
					CallWith(1, 0, 0, node)
			} else if s, ok := arg.v.(string); ok { // c.Str("...")
				cb.Val(ctx.cStr()).Val(s, node).CallWith(1, 0, 0, node)
			} else {
				cb.Val(arg.v, node)
			}
		}
		cb.CallWith(len(v.args), 0, 0, node)
		if results.Len() > 0 {
			cb.Return(1, node)
		}
		cb.EndStmt().End(node)
	}
}

func doCompileMacroFunc(ctx *pkgCtx, decl clang.Cursor, macroName, refName string, paramNames []string) {
	ref, ok := ctx.lookupFuncDecl(refName)
	if !ok {
		ctx.errorf(decl, "macroFunc %s: reference function %s not found", macroName, refName)
		return
	}

	params, ok := matchMacroParams(ctx, decl, macroName, refName, paramNames, ref)
	if !ok {
		return
	}

	fn := &clangMacroFunc{
		Cursor: decl,
		params: params,
		ref:    ref,
	}
	doCompileFunc(ctx, macroName, macroName, fn, &overloadObj{name: macroName}, nil)
}

// macroParamNames parses the parameter list of a function-like macro from its
// tokens. For "#define z_gzgetc(g) ..." the tokens are "z_gzgetc", "(", "g",
// ")", ... so the parameters are the identifiers between the first "(" (which
// immediately follows the macro name for a function-like macro) and the
// matching ")".
func macroParamNames(tu clang.TranslationUnit, tokens []lc.Token) (names []string, iNext int, variadic, ok bool) {
	// tokens[0] is the macro name; tokens[1] must be "(" for a function-like
	// macro. Walk until the matching ")".
	if len(tokens) < 3 || tu.TokenSpelling(tokens[1]) != "(" {
		return
	}
	for i := 2; i < len(tokens); i++ {
		tok := tokens[i]
		switch tok.Kind() {
		case lc.Token_Punctuation:
			switch tu.TokenSpelling(tok) {
			case ",":
				// noop
			case ")":
				iNext, ok = i+1, true
				return
			case "...":
				variadic = true
			default:
				return
			}
		case lc.Token_Identifier:
			names = append(names, tu.TokenSpelling(tok))
		case lc.Token_Comment:
			// noop
		default:
			return
		}
	}
	return
}

// macroFuncParam describes one parameter of a function-like macro listed in
// MacroFunc, together with the reference-function parameter whose type it uses.
type macroFuncParam struct {
	name string  // the macro's own parameter name
	typ  lc.Type // the type supplied by the matched reference-function parameter
}

// matchMacroParams resolves, for every macro parameter, which reference-function
// parameter supplies its type. Name matches are resolved first for all
// parameters; the remaining parameters fall back to position, skipping any
// reference parameter already claimed by a name match.
func matchMacroParams(ctx *pkgCtx, decl clang.Cursor, macroName, refName string, paramNames []string, ref clang.Cursor) ([]macroFuncParam, bool) {
	nRef := int(ref.NumArguments())
	refArgs := make([]clang.Cursor, nRef)
	refByName := make(map[string]int, nRef)
	for i := range nRef {
		arg := ref.Argument(c.Uint(i))
		refArgs[i] = arg
		if name := clang.String(arg); name != "" {
			refByName[name] = i
		}
	}

	params := make([]macroFuncParam, len(paramNames))
	claimed := make([]bool, nRef)
	resolved := make([]bool, len(paramNames))

	// 1) match by name.
	for i, name := range paramNames {
		params[i].name = name
		if j, ok := refByName[name]; ok && !claimed[j] {
			params[i].typ = refArgs[j].Type()
			claimed[j] = true
			resolved[i] = true
		}
	}

	// 2) match remaining by position, skipping claimed reference parameters.
	for i := range paramNames {
		if resolved[i] {
			continue
		}
		if i < nRef && !claimed[i] {
			params[i].typ = refArgs[i].Type()
			claimed[i] = true
			resolved[i] = true
			continue
		}
		ctx.errorf(decl, "macroFunc %s: cannot resolve parameter %s from reference function %s",
			macroName, paramNames[i], refName)
		return nil, false
	}
	return params, true
}

// -----------------------------------------------------------------------------

type clangMacroFunc struct {
	clangNode
	params []macroFuncParam
	ref    clang.Cursor
}

func (fn *clangMacroFunc) RawComment() string {
	raw := clang.RawCommentsOf(fn.Cursor)
	return raw.Text()
}

func (fn *clangMacroFunc) NumArguments() c.Int {
	return c.Int(len(fn.params))
}

func (fn *clangMacroFunc) ArgumentType(i c.Int) lc.Type {
	return fn.params[i].typ
}

func (fn *clangMacroFunc) Argument(i c.Int) (typ lc.Type, name string) {
	param := fn.params[i]
	return param.typ, param.name
}

func (fn *clangMacroFunc) ResultType() lc.Type {
	return fn.ref.ResultType()
}

func (fn *clangMacroFunc) IsVariadic() c.Uint {
	return 0 // not variadic
}

func (fn *clangMacroFunc) IsFunctionInlined() c.Uint {
	return 1 // always inlined
}

// -----------------------------------------------------------------------------
