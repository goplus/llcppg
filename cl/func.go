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
	"go/token"
	"go/types"
	"log"
	"strconv"

	"github.com/goplus/lib/c"
	"github.com/goplus/llcppg/clang"
)

// -----------------------------------------------------------------------------

func loadGlobalFunc(ctx *pkgCtx, scope *scopeCtx, decl clang.Cursor, ns string) {
	name := nameWithNS(clang.String(decl), ns)
	if obj, ok := scope.addFunc(ctx, name, decl); ok {
		ctx.compiles = append(ctx.compiles, func(ctx *pkgCtx) {
			compileFuncOrMethod(ctx, obj, nil)
		})
	}
}

// compileFuncOrMethod compiles a C/C++ function or method into a Go function.
//
// When cls is nil, it is compiled as a receiver-less global function with a
// //go:linkname directive. A C++ static method is compiled this way too: it
// has no implicit "this", so it is loaded as a global function (its Go name is
// prefixed by the enclosing class name, which acts like a namespace). When cls
// is non-nil, it is an instance method compiled with a "this" receiver.
func compileFuncOrMethod(ctx *pkgCtx, obj *funcObj, cls *classCtx) {
	fn := obj.decl
	manglingName := clang.Mangling(fn)
	origName := obj.name
	if fn.IsFunctionInlined() != 0 {
		if ctx.cflags == "" {
			if debugCompileDecl {
				log.Println("inline func", origName, "- skipped")
			}
			return
		}
		manglingName = wrapInlineFunc(ctx, manglingName, fn, cls)
	} else if _, ok := ctx.nameLookup(manglingName); !ok {
		if debugCompileDecl {
			log.Println("func", origName, "- skipped")
		}
		return
	}

	if debugCompileDecl {
		log.Println("func", origName, "-", clang.String(fn.Type()))
	}

	pkg := ctx.pkg
	pkgTypes := pkg.Types
	params, variadic := newParams(ctx, pkgTypes, fn)
	results := toFuncResults(ctx, pkgTypes, fn.ResultType())

	var recv *types.Var
	var typNamed *types.Named
	var nameInPkg string
	if cls == nil {
		if ctx.lang == LanguageC {
			// try to method for C functions
			params, recv, typNamed = tryToMethod(pkgTypes, params)
		}
	} else {
		typNamed = cls.typNamed
		recv = types.NewParam(token.NoPos, pkgTypes, "this", types.NewPointer(typNamed))
	}
	fnName := ctx.funcName(origName, obj.order(), typNamed, cls == nil, true)
	if recv == nil {
		nameInPkg = fnName
	} else {
		objName := typNamed.Obj().Name()
		if _, ok := recv.Type().(*types.Pointer); ok {
			nameInPkg = "(*" + objName + ")." + fnName
		} else {
			nameInPkg = objName + "." + fnName
		}
	}
	sig := types.NewSignatureType(recv, nil, nil, types.NewTuple(params...), results, variadic)
	f, err := pkg.NewFuncWith(goNodePos(ctx, fn), fnName, sig, nil)
	if err != nil {
		log.Panicln("compileFunc:", origName, err)
	}

	if recv == nil {
		ctx.forceImportUnsafe()
		f.SetComments(pkg, ctx.directiveComments(fn,
			"\n//go:linkname "+nameInPkg+" C."+manglingName))
	} else {
		f.SetComments(pkg, ctx.directiveComments(fn,
			"\n// llgo:link "+nameInPkg+" C."+manglingName))
		cb := f.BodyStart(pkg)
		if n := results.Len(); n > 0 {
			cb.ZeroLit(results.At(0).Type()).Return(1)
		}
		cb.End()
	}
}

func tryToMethod(pkgTypes *types.Package, params []*types.Var) ([]*types.Var, *types.Var, *types.Named) {
	if len(params) > 0 {
		first := params[0]
		t := first.Type()
		if len(params) == 2 && params[1].Type() == t {
			// don't convert to method if the first two params have the same type
			return params, nil, nil
		}
		if tp, ok := t.(*types.Pointer); ok {
			t = tp.Elem()
		}
		if tn, ok := t.(*types.Named); ok && tn.Obj().Pkg() == pkgTypes {
			return params[1:], first, tn // can be a method
		}
	}
	return params, nil, nil
}

func newParams(ctx *pkgCtx, pkg *types.Package, fn clang.Cursor) (params []*types.Var, variadic bool) {
	n := fn.NumArguments()
	for i := range n {
		item := fn.Argument(c.Uint(i))
		param := newParam(ctx, pkg, item, i)
		params = append(params, param)
	}
	variadic = fn.IsVariadic() != 0
	if variadic {
		params = append(params, newVariadicParam(pkg))
	}
	return
}

func newParam(ctx *pkgCtx, pkg *types.Package, decl clang.Cursor, i c.Int) *types.Var {
	declName := clang.String(decl)
	declTyp := decl.Type()
	if debugCompileDecl {
		log.Println("  => param", declName, "-", clang.String(declTyp))
	}
	typ := toType(ctx, pkg, declTyp, flagIsParam)
	if declName != "" {
		avoidKeyword(&declName)
	} else {
		declName = "_llcppg_param" + strconv.Itoa(int(i)+1)
	}
	return types.NewParam(goNodePos(ctx, decl), pkg, declName, typ)
}

// -----------------------------------------------------------------------------
