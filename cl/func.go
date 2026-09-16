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
	"go/token"
	"go/types"
	"log"
	"strconv"

	"github.com/goplus/lib/c"
	"github.com/goplus/llcppg/clang"
)

// -----------------------------------------------------------------------------

func loadGlobalFunc(ctx *pkgCtx, scope *scopeCtx, decl clang.Cursor, ns string) {
	name := ns + clang.String(decl)
	obj := scope.addObject(name, decl)
	ctx.compiles = append(ctx.compiles, func(ctx *pkgCtx) {
		compileFuncOrMethod(ctx, decl, obj, nil)
	})
}

func compileFuncOrMethod(ctx *pkgCtx, fn clang.Cursor, obj *object, cls *classCtx) {
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

	var recv *types.Var
	var nameInPkg string
	var fnName, rewritten = ctx.getPubName(origName, obj.order())
	if cls == nil {
		nameInPkg = fnName
	} else {
		typNamed := cls.typDecl.Type()
		nameInPkg = "(*" + typNamed.Obj().Name() + ")." + fnName
		recv = types.NewParam(token.NoPos, pkgTypes, "this", types.NewPointer(typNamed))
	}

	params, variadic := newParams(ctx, pkgTypes, fn)
	results := toFuncResults(ctx, pkgTypes, fn.ResultType())
	sig := types.NewSignatureType(recv, nil, nil, params, results, variadic)
	f, err := pkg.NewFuncWith(goNodePos(ctx, fn), fnName, sig, nil)
	if err != nil {
		log.Panicln("compileFunc:", origName, err)
	}

	if cls == nil {
		ctx.forceImportUnsafe()
		f.SetComments(pkg, &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "\n//go:linkname " + nameInPkg + " C." + manglingName},
			},
		})
		if rewritten {
			substObj(pkgTypes, pkgTypes.Scope(), origName, f)
		}
	} else {
		f.SetComments(pkg, &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "\n// llgo:link " + nameInPkg + " C." + manglingName},
			},
		})
		cb := f.BodyStart(pkg)
		if n := results.Len(); n > 0 {
			cb.ZeroLit(results.At(0).Type()).Return(1)
		}
		cb.End()
	}
}

func newParams(ctx *pkgCtx, pkg *types.Package, fn clang.Cursor) (ret *types.Tuple, variadic bool) {
	n := fn.NumArguments()
	var params []*types.Var
	for i := range n {
		item := fn.Argument(c.Uint(i))
		param := newParam(ctx, pkg, item, i)
		params = append(params, param)
	}
	variadic = fn.IsVariadic() != 0
	if variadic {
		params = append(params, newVariadicParam(pkg))
	}
	ret = types.NewTuple(params...)
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
