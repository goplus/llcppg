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
	"go/ast"
	"go/types"
	"log"

	"github.com/goplus/gogen"
	"github.com/goplus/lib/c"
	"github.com/goplus/llcppg/clang"
)

// -----------------------------------------------------------------------------
//
// Automatic detection of forwarding macros (see goplus/llcppg#1109).
//
// A function-like macro that is not listed in MacroFunc is examined on its own.
// When its body is a single plain forwarding call to a function llcppg already
// generates, and every argument is either one of the macro's own parameters or
// a compile-time constant, the macro is turned into an ordinary Go function
// whose body is that one call. No C wrapper and no //go:linkname are emitted:
// the generated function calls the callee's own generated Go declaration
// directly, filling in the constants.
//
// Detection rules (all must hold; a macro that fails any rule is silently
// skipped, logged only in verbose mode, and can still be handled explicitly via
// MacroFunc):
//
//  1. The body is a single function call whose callee is a plain identifier
//     (optionally parenthesized).
//  2. The callee is a non-variadic function llcppg has parsed and generated.
//  3. Every argument is exactly one macro parameter, or a constant (a value
//     llcppg can evaluate, or a sizeof(T) over a generated complete type).
//  4. Every macro parameter is passed exactly once.

// compiledFn records a generated Go function so a forwarding macro can call it.
// recv is non-nil when the function was generated as a method (its first C
// parameter became the receiver).
type compiledFn struct {
	fn       *types.Func
	recv     *types.Var // non-nil => method; the receiver is the first C param
	variadic bool
}

// -----------------------------------------------------------------------------

// autoMacroFunc attempts automatic detection for a function-like macro that is
// not listed in MacroFunc. The real work is deferred to the dedicated
// autoMacros pass, which runs after every function overload has been registered
// and compiled (so the callee's Go declaration is recorded in ctx.compiledFns).
func autoMacroFunc(ctx *pkgCtx, decl clang.Cursor) {
	ctx.addAutoMacroUnit(func(ctx *pkgCtx) {
		compileAutoMacroFunc(ctx, decl)
	})
}

func compileAutoMacroFunc(ctx *pkgCtx, decl clang.Cursor) {
	macroName := clang.String(decl)

	paramNames, variadic, ok := macroParamNames(decl)
	if !ok {
		return
	}
	if variadic {
		ctx.autoSkip(decl, macroName, "variadic macros are not supported")
		return
	}

	body, dispose, ok := macroFuncBody(decl)
	if !ok {
		return // not a function-like macro definition with a body
	}
	defer dispose()

	calleeName, argGroups, ok := parseForwardingCall(body)
	if !ok {
		ctx.autoSkip(decl, macroName, "body is not a single forwarding call")
		return
	}

	// Rule 2: the callee must be a single non-variadic generated function.
	ref, ok := ctx.lookupFuncDecl(calleeName)
	if !ok {
		ctx.autoSkip(decl, macroName, "callee %s is not a generated function", calleeName)
		return
	}
	cfn, ok := ctx.compiledFns[objUSR(ref)]
	if !ok {
		ctx.autoSkip(decl, macroName, "callee %s was not generated in this package", calleeName)
		return
	}
	if cfn.variadic || ref.IsVariadic() != 0 {
		ctx.autoSkip(decl, macroName, "callee %s is variadic", calleeName)
		return
	}

	nRef := int(ref.NumArguments())
	if len(argGroups) != nRef {
		ctx.autoSkip(decl, macroName, "callee %s expects %d arguments, macro passes %d",
			calleeName, nRef, len(argGroups))
		return
	}

	// Classify every argument as either a macro parameter or a constant, and
	// enforce rule 4 (each macro parameter passed exactly once).
	paramUse := make(map[string]int, len(paramNames))
	for _, p := range paramNames {
		paramUse[p] = 0
	}
	args := make([]autoArg, nRef)
	for i := range nRef {
		refArg := ref.Argument(c.Uint(i))
		arg, ok := classifyArg(ctx, decl, macroName, argGroups[i], refArg, paramUse)
		if !ok {
			return // classifyArg already logged
		}
		args[i] = arg
	}
	for _, p := range paramNames {
		if n := paramUse[p]; n != 1 {
			ctx.autoSkip(decl, macroName, "macro parameter %s is passed %d times (must be exactly once)", p, n)
			return
		}
	}

	genAutoMacroFunc(ctx, decl, macroName, paramNames, args, ref, cfn)
}

// -----------------------------------------------------------------------------

// autoArg is one classified macro-body argument, in callee-parameter order.
type autoArg struct {
	param    string        // non-empty => a macro parameter with this name
	typ      types.Type    // the callee parameter's Go type (used for params and conversions)
	kind     autoConstKind // for constants
	val      any           // evaluated constant value (numeric/string)
	sizeType types.Type    // for sizeof: the generated Go type T (composite literal T{})
	castTo   types.Type    // explicit cast type around the argument, if any
}

type autoConstKind int

const (
	autoParam       autoConstKind = iota // a macro parameter, not a constant
	autoConstValue                       // an evaluated numeric or string value
	autoConstSizeof                      // a sizeof(T) expression
)

// classifyArg decides whether one macro-body argument is a macro parameter or a
// constant, validating it against the callee parameter refArg whose type it
// uses. It updates paramUse when the argument is a macro parameter.
func classifyArg(ctx *pkgCtx, decl clang.Cursor, macroName string, arg toks, refArg clang.Cursor, paramUse map[string]int) (autoArg, bool) {
	refTyp := toType(ctx, ctx.pkg.Types, refArg.Type(), flagIsParam, &ctx.scopeCtx)

	// A bare identifier that is one of the macro's parameters.
	if name, ok := singleIdent(arg); ok {
		if _, isParam := paramUse[name]; isParam {
			paramUse[name]++
			return autoArg{param: name, typ: refTyp, kind: autoParam}, true
		}
	}

	// sizeof(T) [optionally cast to a scalar type].
	if szType, castTo, ok := parseSizeof(ctx, arg); ok {
		return autoArg{kind: autoConstSizeof, typ: refTyp, sizeType: szType, castTo: castTo}, true
	}

	// A value llcppg can evaluate at generation time (literal, constant macro,
	// cast or arithmetic over them). The tokens must not mention a macro
	// parameter (otherwise it is neither a plain parameter nor a constant).
	if mentionsParam(ctx, arg, paramUse) {
		ctx.autoSkip(decl, macroName, "argument mentions a macro parameter but is not a plain parameter")
		return autoArg{}, false
	}
	if v, ok := evalConstTokens(ctx, stripOuterParens(arg)); ok {
		return autoArg{kind: autoConstValue, typ: refTyp, val: v}, true
	}

	ctx.autoSkip(decl, macroName, "argument is neither a macro parameter nor a constant llcppg can evaluate")
	return autoArg{}, false
}

// -----------------------------------------------------------------------------

// genAutoMacroFunc generates the Go function for a detected forwarding macro.
func genAutoMacroFunc(ctx *pkgCtx, decl clang.Cursor, macroName string, paramNames []string, args []autoArg, ref clang.Cursor, cfn *compiledFn) {
	pkg := ctx.pkg
	pkgTypes := pkg.Types

	fnName := ctx.funcName(macroName, -1, "", "", true, true)

	// Name collision: a Go name derived from a macro may clash with another
	// generated symbol. Skip (with a warning) rather than emit a redeclaration.
	if o := pkgTypes.Scope().Lookup(fnName); o != nil {
		log.Printf("[WARN] macroFunc %s: Go name %s collides with %s, skipped\n",
			macroName, fnName, o.Name())
		return
	}

	// The generated function's parameters are the macro's own parameters, each
	// typed by the callee parameter it forwards to. Map param name => Go type.
	paramType := make(map[string]types.Type, len(paramNames))
	for _, a := range args {
		if a.kind == autoParam {
			paramType[a.param] = a.typ
		}
	}
	params := make([]*types.Var, len(paramNames))
	for i, name := range paramNames {
		params[i] = types.NewParam(goNodePos(ctx, decl), pkgTypes, avoidKeyword(name), paramType[name])
	}

	results := toFuncResults(ctx, pkgTypes, ref.ResultType(), new(int), &ctx.scopeCtx)

	sig := types.NewSignatureType(nil, nil, nil, types.NewTuple(params...), results, false)
	f, err := pkg.NewFuncWith(goNodePos(ctx, decl), fnName, sig, nil)
	if err != nil {
		ctx.errorf(decl, "macroFunc %s: %v", macroName, err)
		return
	}
	f.SetComments(pkg, autoMacroDoc(ctx, decl, ref))

	cb := f.BodyStart(pkg)
	hasResult := results.Len() > 0

	if cfn.recv == nil {
		// Package-level function callee: callee(a0, a1, ...).
		cb.Val(cfn.fn)
		for _, a := range args {
			pushAutoArg(ctx, cb, a)
		}
		cb.Call(len(args))
	} else {
		// Method callee: recvArg.Method(a1, a2, ...).
		pushAutoArg(ctx, cb, args[0])
		cb.MemberVal(cfn.fn.Name(), 1)
		for _, a := range args[1:] {
			pushAutoArg(ctx, cb, a)
		}
		cb.Call(len(args) - 1)
	}

	if hasResult {
		cb.Return(1)
	} else {
		cb.EndStmt()
	}
	cb.End()
}

// pushAutoArg pushes one call argument onto the code builder.
func pushAutoArg(ctx *pkgCtx, cb *gogen.CodeBuilder, a autoArg) {
	switch a.kind {
	case autoParam:
		cb.VarVal(avoidKeyword(a.param))
	case autoConstValue:
		if s, ok := a.val.(string); ok {
			// string constant => c.Str("...")
			cb.Val(ctx.c.Ref("Str")).Val(s).Call(1)
		} else {
			// numeric constant: evaluated value, converted to the callee's
			// parameter type so it fits (e.g. c.Int(-1)).
			cb.Typ(a.typ).Val(a.val).Call(1)
		}
	case autoConstSizeof:
		// unsafe.Sizeof(T{}) converted to the cast type (if any) or the
		// callee's parameter type (unsafe.Sizeof yields uintptr).
		//
		// The expression is built as a raw AST element rather than through
		// gogen's unsafe.Sizeof instruction: that instruction eagerly computes
		// the constant size via go/types SizesFor(runtime.Compiler, ...), which
		// is unavailable under the llgo toolchain that runs llcppg and would
		// crash at generation time. We only need to emit the call; its value is
		// computed by the Go compiler for the target platform (the whole point
		// of routing sizeof through unsafe.Sizeof).
		target := a.castTo
		if target == nil {
			target = a.typ
		}
		ctx.forceImportUnsafe()
		cb.Typ(target)
		cb.Val(sizeofElem(ctx, cb, a.sizeType))
		cb.Call(1) // target(unsafe.Sizeof(T{}))
	}
}

// sizeofElem builds a raw "unsafe.Sizeof(<zero-literal of typ>)" expression as a
// gogen element of type uintptr. The zero literal is the composite literal T{}
// for struct/record types, or a generic zero value for other types.
//
// The "unsafe.Sizeof" selector is obtained from gogen (by briefly pushing and
// popping the reference) rather than hand-built, so its package identifier
// carries the ast.Object gogen needs to emit a proper `import "unsafe"` (not a
// blank import) when the generated AST is finalized. The call expression itself
// is assembled by hand to bypass gogen's unsafe.Sizeof instruction, which would
// otherwise eagerly compute the size through go/types SizesFor — unavailable
// under the llgo toolchain that runs llcppg — and crash at generation time.
func sizeofElem(ctx *pkgCtx, cb *gogen.CodeBuilder, typ types.Type) *gogen.Element {
	zero := ctx.pkg.Zero(typ) // e.g. T{}, nil, 0 — carries the AST in .Val

	// Obtain gogen's "unsafe.Sizeof" selector expression (with a package ident
	// whose ast.Object marks the import as used).
	cb.Val(ctx.pkg.Unsafe().Ref("Sizeof"))
	fn := cb.InternalStack().Pop().Val

	call := &ast.CallExpr{Fun: fn, Args: []ast.Expr{zero.Val}}
	return &gogen.Element{Val: call, Type: types.Typ[types.Uintptr]}
}

// autoMacroDoc resolves the doc comment for a detected macro function: the
// macro's own doc comment if present, otherwise the callee's. The two are never
// merged.
func autoMacroDoc(ctx *pkgCtx, decl clang.Cursor, ref clang.Cursor) *ast.CommentGroup {
	if doc := ctx.docComments(decl); len(doc) != 0 {
		return docCommentGroup(doc)
	}
	return docCommentGroup(ctx.docComments(ref))
}

// -----------------------------------------------------------------------------

func (p *pkgCtx) autoSkip(decl clang.Cursor, macroName, format string, args ...any) {
	if debugCompileDecl {
		p.logf(decl, "macroFunc %s (auto): skipped - "+format, append([]any{macroName}, args...)...)
	}
}
