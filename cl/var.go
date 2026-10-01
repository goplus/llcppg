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
	"log"

	"github.com/goplus/llcppg/clang"
	lc "github.com/llarhub/clang-c"
)

// -----------------------------------------------------------------------------

// loadVar loads a C/C++ variable declaration (a global variable, a variable
// within a namespace, or a static variable within a class) and schedules it to
// be compiled into a Go variable.
//
// The Go name is prefixed by ns, which encodes the enclosing namespaces and/or
// class (the class name acts like a namespace, mirroring how static methods are
// handled). Non-static class member variables are fields, not VarDecls, and are
// handled separately in loadClassMember.
func loadVar(ctx *pkgCtx, decl clang.Cursor, ns string) {
	ctx.addCompileUnit(func(ctx *pkgCtx) {
		compileVar(ctx, decl, ns)
	})
}

// compileVar compiles a C/C++ variable into a Go variable declared with a
// //go:linkname directive that binds it to the external C/C++ symbol, similar
// to how a global function is compiled in compileFuncOrMethod.
//
// A const-qualified variable (e.g. const int x;) is handled the same way: in
// terms of C/C++ semantics a const is just an immutable var, and there is no
// difference in the underlying binding logic. It maps to the same Go var linked
// to the same C/C++ symbol; Go has no way to express a linkname-bound immutable
// value, so the const-ness is not reflected in the generated declaration. The
// const qualifier is a property of the type, not a distinct type kind, so
// toType (which switches on the type kind) already yields the unqualified Go
// type without any special casing here.
func compileVar(ctx *pkgCtx, decl clang.Cursor, ns string) {
	localName := clang.String(decl)
	if varHasInitExpr(ctx, decl) {
		ctx.logf(decl, "var %s: has initialized expression, skipped", localName)
		return
	}

	manglingName := clang.Mangling(decl)
	if _, ok := ctx.nameLookup(manglingName); !ok {
		ctx.logf(decl, "var %s: symbol not found in lib files, skipped", localName)
		return
	}

	pkg := ctx.pkg
	pkgTypes := pkg.Types

	origName := nameWithNS(localName, ns)
	if debugCompileDecl {
		kind := "var"
		if decl.Type().IsConstQualified() != 0 {
			kind = "const var"
		}
		log.Println(kind, origName, "-", clang.String(decl.Type()))
	}

	feats := 0
	typ := toTypeEx(ctx, pkgTypes, decl.Type(), flagIsVarDef, &feats, nil)
	if feats&featAllIgnore != 0 {
		ctx.ignoref(feats, decl, "var %s: unsupported type, ignored", localName)
		return
	}

	goName := ctx.varName(origName)
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
