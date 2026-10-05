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
	ns := cNS(decl)
	cName := cNameWithNS(clang.String(decl), ns)
	if feats := ctx.nsFeats(ns); feats&featAllIgnore != 0 {
		ctx.ignoref(feats, decl, "var %s: parent is ignored", cName)
		return
	}
	if varHasInitExpr(ctx, decl) {
		ctx.ignoref(featQuietIgnore, decl, "var %s: has initialized expression, ignored", cName)
		return
	}

	manglingName := clang.Mangling(decl)
	if _, ok := ctx.nameLookup(manglingName); !ok {
		ctx.ignoref(featExplicitIgnore, decl, "var %s: symbol not found in lib files, ignored", cName)
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

func varHasInitExpr(ctx *pkgCtx, v clang.Cursor) (hasInitExpr bool) {
	clang.VisitChildren(v, func(decl, parent clang.Cursor) clang.ChildVisitResult {
		switch decl.Kind {
		case lc.Cursor_InitListExpr, lc.Cursor_CallExpr, lc.Cursor_DeclRefExpr,
			lc.Cursor_UnaryOperator, lc.Cursor_BinaryOperator, lc.Cursor_UnexposedExpr,
			lc.Cursor_CXXStaticCastExpr:
			hasInitExpr = true
			return clang.Break
		case lc.Cursor_TypeRef, lc.Cursor_UnaryExpr, lc.Cursor_ParmDecl,
			lc.Cursor_IntegerLiteral:
		default:
			if lc.Cursor_FirstAttr > decl.Kind || decl.Kind > lc.Cursor_LastAttr {
				ctx.panicf(decl, "varHasInitExpr: unknown kind - %v", decl.Kind)
			}
		}
		return clang.Continue
	})
	return
}

// -----------------------------------------------------------------------------
