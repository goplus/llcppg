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
		ctx.logf(decl, "typedef %s: %s", cName, clang.String(underlying))
	}

	feats := 0
	tunder := toTypeEx(ctx, pkgTypes, underlying, flagIsTypeDef, &feats, scope)
	if feats&featQuietIgnore != 0 {
		ctx.ignoref(featQuietIgnore, decl, "typedef %s: unsupported underlying type (%d: %v), ignored", cName, underlying.Kind, clang.String(underlying))
		ctx.ignoreType(cName, featQuietIgnore)
		return
	}
	if feats&featExplicitIgnore != 0 {
		ctx.ignoref(featExplicitIgnore, decl, "typedef %s: unsupported underlying type (%d: %v), ignored", cName, underlying.Kind, clang.String(underlying))
		ctx.ignoreType(cName, featExplicitIgnore)
		return
	}

	if tp, ok := tunder.(*types.TypeParam); ok {
		// A typedef whose underlying type is a bare template parameter, e.g.
		//   template <class _Tp> class shared_ptr { typedef _Tp element_type; };
		// cannot be emitted as a Go top-level type (Go has no
		// "type X[_Tp any] = _Tp"). But it is NOT unsupported: references to
		// it, such as the field "element_type* __ptr_", must resolve to the
		// template parameter itself. Register it transparently so lookups
		// return _Tp, without emitting a declaration and without ignoring it
		// (ignoring would poison every field that uses it and, through them,
		// the whole enclosing class). See issue goplus/llcppg#985.
		defineTypedefToTypeParam(ctx, decl, cName, tp)
		return
	}

	goName := ctx.typeName(cName, true)
	if tn, ok := tunder.(*types.Named); ok {
		if o := tn.Obj(); o.Pkg() == pkgTypes && o.Name() == goName {
			return // already defined
		}
	}

	// A C name can be typedef'd more than once across headers, e.g. a forward
	// typedef "typedef struct _Foo Foo;" in one header and the defining typedef
	// "typedef struct _Foo { ... } Foo;" in another. Both reach here and both
	// resolve to goName. The first one is emitted by defineTypedef below, which
	// inserts goName into the package scope; a later one finds it there. If the
	// two underlying types match, the later typedef is a harmless redeclaration
	// and is reused; if they differ, report it as a redeclaration error instead
	// of emitting a second Go type with the same name. This reuses the existing
	// scope info rather than keeping a separate map. See issue goplus/llcppg#1001.
	if o := pkgTypes.Scope().Lookup(goName); o != nil {
		eunder, _ := ctx.effectiveUnder(cName, tunder)
		if prev, ok := o.(*types.TypeName); ok && types.Identical(prev.Type(), eunder) {
			return // same underlying type, reuse the existing definition
		}
		ctx.errorf(decl, "%s redeclared in this block\n\t%v: other declaration of %s",
			goName, ctx.position(o.Pos()), goName)
		return
	}

	defineTypedef(ctx, decl, cName, goName, scope, tunder, tparams, feats)
}

// defineTypedefToTypeParam registers a typedef whose underlying type is a bare
// template parameter (e.g. "typedef _Tp element_type;" inside a class template)
// as a transparent alias to that parameter. No Go top-level declaration is
// emitted - Go cannot express "type X[_Tp any] = _Tp" - but the C/C++ name is
// made resolvable so that fields and other members referring to it resolve to
// the template parameter instead of being treated as unsupported. See issue
// goplus/llcppg#985.
func defineTypedefToTypeParam(ctx *pkgCtx, decl clang.Cursor, cName string, tp *types.TypeParam) {
	obj := types.NewTypeName(goNodePos(ctx, decl), ctx.pkg.Types, tp.Obj().Name(), tp)
	ctx.types[cName] = typeObj{obj, 0}
	if debugCompileDecl {
		ctx.logf(decl, "==> addType %s: %v (transparent type param)", cName, tp)
	}
}

func defineTypedef(ctx *pkgCtx, decl clang.Cursor, cName, goName string, scope *scopeCtx, tunder types.Type, tparams []*types.TypeParam, feats int) {
	pkg := ctx.pkg
	typDefs := pkg.NewTypeDefs()
	if doc := ctx.directiveTypeC(decl, feats&featHasCallback != 0); doc != nil {
		typDefs.SetComments(doc)
	}

	var isClass bool
	tunder, isClass = ctx.effectiveUnder(cName, tunder)

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

// effectiveUnder maps a typedef's resolved underlying type to the type actually
// emitted for it: unsafe.Pointer becomes uintptr (as a named type, unless cName
// is configured as a non-class), void becomes c.Void, and a cName configured as
// a class is emitted as a named type. The returned isClass reports whether the
// typedef is emitted as a named type (InitType) rather than a type alias
// (AliasType). Keeping this in one place lets the redeclaration check in
// loadTypedef compare against the same effective type defineTypedef emits. See
// issue goplus/llcppg#1001.
func (ctx *pkgCtx) effectiveUnder(cName string, tunder types.Type) (types.Type, bool) {
	switch tunder {
	case ctx.unsafePointer():
		if !contains(cName, ctx.nonClasses) {
			return types.Typ[types.Uintptr], true // unsafe.Pointer => uintptr
		}
		return tunder, false
	case tyVoid:
		return ctx.basicTyp(cVoid), false
	}
	return tunder, contains(cName, ctx.classes)
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

// alias = void (ignore this type, as an empty class)
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
		if alias == "void" {
			p.types[name] = typeObj{nil, featTyIgnore}
			return tyIgnore, true
		}
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
