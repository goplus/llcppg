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
	"github.com/goplus/llcppg/clang"
	lc "github.com/llarhub/clang-c"
)

// -----------------------------------------------------------------------------

func loadVar(ctx *pkgCtx, decl clang.Cursor) {
	ctx.addCompileUnit(func(ctx *pkgCtx) {
		compileVar(ctx, decl)
	})
}

func compileVar(ctx *pkgCtx, decl clang.Cursor) {
	cName := cNameOf(decl)
	if varHasInitExpr(ctx, decl) {
		ctx.logf(decl, "var %s: has initialized expression, skipped", cName)
		return
	}

	manglingName := clang.Mangling(decl)
	if _, ok := ctx.nameLookup(manglingName); !ok {
		ctx.logf(decl, "var %s: symbol not found in lib files, skipped", cName)
		return
	}

	if debugCompileDecl {
		ctx.logf(decl, "var %s: %s", cName, clang.String(decl.Type()))
	}

	pkg := ctx.pkg
	pkgTypes := pkg.Types

	feats := 0
	typ := toTypeEx(ctx, pkgTypes, decl.Type(), flagIsVarDef, &feats, nil)
	if feats&featAllIgnore != 0 {
		ctx.ignoref(feats, decl, "var %s: unsupported type, ignored", cName)
		return
	}

	goName := ctx.varName(cName)
	ctx.forceImportUnsafe()
	defs := pkg.NewVarDefs(pkgTypes.Scope()).SetComments(
		ctx.directiveComments(decl, "\n//go:linkname "+goName+" C."+manglingName))
	defs.New(goNodePos(ctx, decl), typ, goName)
}

func varHasInitExpr(_ *pkgCtx, v clang.Cursor) (hasInitExpr bool) {
	clang.VisitChildren(v, func(decl, parent clang.Cursor) clang.ChildVisitResult {
		switch decl.Kind {
		case lc.Cursor_InitListExpr, lc.Cursor_CallExpr:
			hasInitExpr = true
			return clang.Break
		default:
			// ctx.panicf(decl, "varHasInitExpr: unknown kind - %v", decl.Kind)
		}
		return clang.Continue
	})
	return
}

// -----------------------------------------------------------------------------
