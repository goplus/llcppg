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

func loadTypedef(ctx *pkgCtx, decl clang.Cursor, scope *scopeCtx) {
	cName := clang.String(decl.Type())
	cFullName := cNameOf(decl)
	if ctx.isTypeIgnored(cFullName) {
		ctx.ignoref(featQuietIgnore, decl, "typedef %s: ignored by config", cFullName)
		ctx.ignoreType(cName, featQuietIgnore)
		return
	}

	if decl.Kind == lc.Cursor_TypeAliasTemplateDecl {
		ctx.logf(decl, "typedef %s: template type alias, skipped", cFullName)
		return
	}

	pkg := ctx.pkg
	pkgTypes := pkg.Types
	underlying := decl.TypedefDeclUnderlyingType()
	if debugCompileDecl {
		ctx.logf(decl, "typedef %s (%s) - %s", cFullName, cName, clang.String(underlying))
	}

	feats := 0
	tunder := toTypeEx(ctx, pkgTypes, underlying, flagIsTypeDef, &feats, scope)
	if feats&featQuietIgnore != 0 || isTypedefUnsupported(tunder) {
		ctx.ignoref(featQuietIgnore, decl, "typedef %s: unsupported underlying type, ignored", cFullName)
		ctx.ignoreType(cName, featQuietIgnore)
		return
	}
	if feats&featExplicitIgnore != 0 {
		ctx.logf(decl, "typedef %s: unsupported underlying type, skipped", cFullName)
		ctx.ignoreType(cName, featExplicitIgnore)
		return
	}

	goName := ctx.typeName(cFullName, true)
	if tn, ok := tunder.(*types.Named); ok {
		if o := tn.Obj(); o.Pkg() == pkgTypes && o.Name() == goName {
			return // already defined
		}
	}

	var obj *types.TypeName
	var typDefs = pkg.NewTypeDefs()
	if doc := ctx.directiveTypeC(decl, feats&featHasCallback != 0); doc != nil {
		typDefs.SetComments(doc)
	}
	var isClass bool
	switch tunder {
	case ctx.unsafePointer():
		if !contains(cName, ctx.nonClasses) {
			tunder, isClass = types.Typ[types.Uintptr], true // unsafe.Pointer => uintptr
		}
	case tyVoid:
		tunder = ctx.basicTyp(cVoid)
	default:
		isClass = contains(cName, ctx.classes)
	}

	node := goNode(ctx, decl)
	tparams := scope.typeParams(nil)
	if isClass {
		t := typDefs.NewType(goName, node).InitType(pkg, tunder, tparams...)
		obj = t.Obj()
	} else {
		t := typDefs.AliasTypeEx(goName, tunder, tparams, node)
		obj = t.Obj()
	}

	ctx.types[cName] = typeObj{obj, feats}
	ctx.aliasTypeName(cName, goName)
}

func isTypedefUnsupported(tunder types.Type) bool {
	switch tunder.(type) {
	case *types.TypeParam:
		return true
	}
	return false
}

// -----------------------------------------------------------------------------
