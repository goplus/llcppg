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
)

// -----------------------------------------------------------------------------

func loadTypedef(ctx *pkgCtx, decl clang.Cursor, scope *scopeCtx) {
	ns := cNS(decl)
	cName := cNameWithNS(clang.String(decl), ns)
	if ctx.isConfTypeIgnored(cName) {
		ctx.ignoref(featQuietIgnore, decl, "typedef %s: ignored by config", cName)
		ctx.ignoreType(cName, featQuietIgnore)
		return
	}

	goName := ctx.typeName(cName, true)
	if t, ok := ctx.typeAliasOf(cName); ok {
		defineTypedef(ctx, decl, cName, goName, scope, t, nil, 0)
		return
	}

	if feats := ctx.nsFeats(ns); feats&featAllIgnore != 0 {
		ctx.ignoref(feats, decl, "typedef %s: its parent is ignored, ignored too", cName)
		ctx.ignoreType(cName, feats)
		return
	}

	pkg := ctx.pkg
	pkgTypes := pkg.Types
	tparams, quietIgnore := newTypeParams(ctx, pkgTypes, decl)
	if quietIgnore {
		ctx.ignoref(featQuietIgnore, decl, "typedef %s: unsupported template params, ignored", cName)
		ctx.ignoreType(cName, featQuietIgnore)
		return
	}

	underlying := decl.TypedefDeclUnderlyingType()
	if debugCompileDecl {
		ctx.logf(decl, "typedef %ss: %s", cName, clang.String(underlying))
	}

	feats := 0
	tunder := toTypeEx(ctx, pkgTypes, underlying, flagIsTypeDef, &feats, scope)
	if feats&featQuietIgnore != 0 || isTypedefUnsupported(tunder) {
		ctx.ignoref(featQuietIgnore, decl, "typedef %s: unsupported underlying type (%v), ignored", cName, clang.String(underlying))
		ctx.ignoreType(cName, featQuietIgnore)
		return
	}
	if feats&featExplicitIgnore != 0 {
		ctx.ignoref(featExplicitIgnore, decl, "typedef %s: unsupported underlying type (%v), ignored", cName, clang.String(underlying))
		ctx.ignoreType(cName, featExplicitIgnore)
		return
	}

	if tn, ok := tunder.(*types.Named); ok {
		if o := tn.Obj(); o.Pkg() == pkgTypes && o.Name() == goName {
			return // already defined
		}
	}

	defineTypedef(ctx, decl, cName, goName, scope, tunder, tparams, feats)
}

func defineTypedef(ctx *pkgCtx, decl clang.Cursor, cName, goName string, scope *scopeCtx, tunder types.Type, tparams []*types.TypeParam, feats int) {
	pkg := ctx.pkg
	typDefs := pkg.NewTypeDefs()
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
	tparams = scope.typeParams(tparams)

	var obj *types.TypeName
	if isClass {
		t := typDefs.NewType(goName, node).InitType(pkg, tunder, tparams...)
		obj = t.Obj()
	} else {
		t := typDefs.AliasTypeEx(goName, tunder, tparams, node)
		obj = t.Obj()
	}

	ctx.types[cName] = typeObj{obj, 0}
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
