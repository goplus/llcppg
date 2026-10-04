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
	"strings"

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

	if _, ok := ctx.typeAliasOf(cName); ok {
		// NOTE(xsw): typeAliasOf does the type alias when cName is first encountered.
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
		ctx.ignoref(featQuietIgnore, decl, "typedef %s: unsupported underlying type (%d: %v), ignored", cName, underlying.Kind, clang.String(underlying))
		ctx.ignoreType(cName, featQuietIgnore)
		return
	}
	if feats&featExplicitIgnore != 0 {
		ctx.ignoref(featExplicitIgnore, decl, "typedef %s: unsupported underlying type (%d: %v), ignored", cName, underlying.Kind, clang.String(underlying))
		ctx.ignoreType(cName, featExplicitIgnore)
		return
	}

	goName := ctx.typeName(cName, true)
	if tn, ok := tunder.(*types.Named); ok {
		if o := tn.Obj(); o.Pkg() == pkgTypes && o.Name() == goName {
			return // already defined
		}
	}

	defineTypedef(ctx, decl, cName, goName, scope, tunder, tparams, feats)
}

func isTypedefUnsupported(tunder types.Type) bool {
	switch tunder.(type) {
	case *types.TypeParam:
		return true
	}
	return false
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
		t := typDefs.NewType(goName, tparams, node).InitType(pkg, tunder)
		obj = t.Obj()
	} else {
		t := typDefs.NewType(goName, tparams, node).AliasType(pkg, tunder)
		obj = t.Obj()
	}

	ctx.types[cName] = typeObj{obj, 0}
	ctx.aliasTypeName(cName, goName)
}

// -----------------------------------------------------------------------------

func doAliasType(ctx *pkgCtx, cName string, tunder types.Type, scope *scopeCtx) *types.Alias {
	var tparams []*types.TypeParam
	if scope != nil {
		tparams = scope.tparams
	} else {
		tparams, tunder = typeParamsAndInstantiate(ctx, tunder)
	}
	goName := ctx.typeName(cName, true)
	t := ctx.pkg.NewTypeDefs().NewType(goName, tparams).AliasType(ctx.pkg, tunder)
	obj := t.Obj()
	ctx.types[cName] = typeObj{obj, 0}
	ctx.aliasTypeName(cName, goName)
	return t
}

func typeParamsAndInstantiate(ctx *pkgCtx, t types.Type) (tparams []*types.TypeParam, inst types.Type) {
	var tlist *types.TypeList
	var tplist *types.TypeParamList
	switch tt := t.(type) {
	case *types.Named:
		tplist, tlist = tt.TypeParams(), tt.TypeArgs()
	case *types.Alias:
		tplist, tlist = tt.TypeParams(), tt.TypeArgs()
	}
	if tlist.Len() > 0 || tplist.Len() == 0 {
		return nil, t
	}
	n := tplist.Len()
	tparams = make([]*types.TypeParam, n)
	targs := make([]types.Type, n)
	for i := range n {
		tp := cloneTypeParam(tplist.At(i))
		targs[i] = tp
		tparams[i] = tp
	}
	inst, _ = types.Instantiate(ctx.typeCtx(), t, targs, false)
	return
}

// alias = .StdBasicString[byte, c.Void, c.Void]
// alias = .Iterator
// alias = [_Derived, _Category, _Tp, _Distance, _Pointer, _Reference] = .Iterator[_Category, _Tp, _Distance, _Pointer, _Reference]
func (p *pkgCtx) aliasType(name, alias string) (ret types.Type, found bool) {
	var ok bool
	var scope *scopeCtx
	if alias[0] == '[' { // has typeParams
		pos := strings.IndexByte(alias, ']')
		if pos <= 0 {
			return
		}
		typParams := alias[1:pos]
		alias, ok = strings.CutPrefix(strings.TrimSpace(alias[pos+1:]), "=")
		if !ok {
			return
		}
		alias = strings.TrimLeft(alias, " \t")
		scope = &scopeCtx{tparams: p.goTypeParams(typParams)}
	}
	typArgs := ""
	if alias[len(alias)-1] == ']' { // has typeArgs
		pos := strings.IndexByte(alias, '[')
		if pos <= 0 {
			return
		}
		typArgs = alias[pos+1 : len(alias)-1]
		alias = alias[:pos]
	}
	obj, ok := p.goNamedTypeObj(alias, scope)
	if !ok {
		return
	}
	ret = obj.Type()
	if typArgs != "" {
		targs, ok := p.goTypeArgs(typArgs, scope)
		if !ok {
			return
		}
		inst, err := types.Instantiate(p.typeCtx(), ret, targs, true)
		if err != nil {
			return
		}
		ret = inst
	}
	return doAliasType(p, name, ret, scope), true
}

// -----------------------------------------------------------------------------
