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
	"sort"

	"github.com/goplus/lib/c"
	"github.com/goplus/llcppg/clang"
	lc "github.com/llarhub/clang-c"
)

func cloneTypeName(obj *types.TypeName) *types.TypeName {
	return types.NewTypeName(obj.Pos(), obj.Pkg(), obj.Name(), obj.Type())
}

func cloneTypeParams(tparams []*types.TypeParam) []*types.TypeParam {
	if len(tparams) == 0 {
		return nil
	}
	ret := make([]*types.TypeParam, len(tparams))
	for i, tp := range tparams {
		obj := cloneTypeName(tp.Obj())
		ret[i] = types.NewTypeParam(obj, tp.Constraint())
	}
	return ret
}

func concatTypeParams(a, b []*types.TypeParam) []*types.TypeParam {
	if len(a) == 0 {
		return b
	}
	a = cloneTypeParams(a)
	if len(b) == 0 {
		return a
	}
	ret := make([]*types.TypeParam, 0, len(a)+len(b))
	ret = append(ret, a...)
	return append(ret, b...)
}

// -----------------------------------------------------------------------------

type scopeCtx struct {
	overloads map[string]*overloads // cName => overload items
	tparams   []*types.TypeParam    // type parameters, only for class/struct scope
	parent    *scopeCtx
}

func (p *scopeCtx) typeParams(in []*types.TypeParam) []*types.TypeParam {
	for p != nil {
		in = concatTypeParams(p.tparams, in)
		p = p.parent
	}
	return in
}

func (p *scopeCtx) lookupType(name string) (types.Type, bool) {
	for p != nil {
		for _, t := range p.tparams {
			if t.Obj().Name() == name {
				return t, true
			}
		}
		p = p.parent
	}
	return nil, false
}

func (p *scopeCtx) addOverloadObj(ctx *pkgCtx, name string, decl clang.Cursor) (obj *overloadObj, isNew bool) {
	oUSR := objUSR(decl)
	if debugCompileDecl {
		ctx.logf(decl, "==> addOveerloadObj %s - USR: %s", name, oUSR)
	}

	obj, isOld := ctx.ovobjs[oUSR]
	if !isOld {
		isNew = true
		obj = &overloadObj{
			name: name,
			decl: decl,
		}
		ovs, ok := p.overloads[name]
		if ok {
			ovs.items = append(ovs.items, obj)
		} else {
			ovs = &overloads{items: []*overloadObj{obj}}
			p.overloads[name] = ovs
		}
		obj.overloads = ovs
		ctx.ovobjs[oUSR] = obj
	}
	if decl.IsCursorDefinition() != 0 { // TODO(xsw): redefinition
		obj.decl = decl // use the last definition
	}
	return
}

func (p *scopeCtx) reorder() {
	for _, o := range p.overloads {
		o.reorder()
	}
}

// -----------------------------------------------------------------------------

type overloadObj struct {
	name      string       // local name in scope
	decl      clang.Cursor // AST object
	overloads *overloads
}

// order returns the order of the overloadObj in the overloads list.
// -1 means no order (only one overload, or not found).
func (p *overloadObj) order() int {
	ovs := p.overloads
	if ovs == nil {
		return -1
	}
	items := ovs.items
	if len(items) > 1 {
		for i, obj := range items {
			if obj == p {
				return i
			}
		}
	}
	return -1
}

type overloads struct {
	items []*overloadObj
}

func (p *overloads) reorder() {
	items := p.items
	if len(items) > 1 {
		sort.SliceStable(items, func(i, j int) bool {
			a, b := items[i].decl, items[j].decl
			// NumArguments reports -1 for cursors that are not functions or
			// methods (e.g. a class template and its partial specializations,
			// which share one overload group under the same C/C++ name). Clamp
			// to 0 so such cursors are treated as having no arguments: the
			// argument-based ordering below is only meaningful for function
			// overloads, and feeding -1 into the unsigned loop counter c.Uint(na)
			// would otherwise wrap to ~4.3 billion iterations and hang the
			// generator. See issue goplus/llcppg#894.
			na, nb := a.NumArguments(), b.NumArguments()
			if na < 0 {
				na = 0
			}
			if nb < 0 {
				nb = 0
			}
			if na != nb {
				return na < nb
			}
			for k := range c.Uint(na) {
				ta, tb := a.Argument(k).Type(), b.Argument(k).Type()
				if ret := cmpType(ta, tb); ret != 0 {
					return ret < 0
				}
			}
			return false
		})
	}
}

// -----------------------------------------------------------------------------

func loadDecl(ctx *pkgCtx, scope *scopeCtx, decl clang.Cursor) {
	switch decl.Kind {
	case lc.Cursor_ClassDecl, lc.Cursor_StructDecl:
		preloadClass(ctx, scope, decl)
	case lc.Cursor_FunctionDecl:
		preloadGlobalFunc(ctx, scope, decl)
	case lc.Cursor_CXXMethod, lc.Cursor_Constructor, lc.Cursor_Destructor:
		ctx.addLoadUnit(func(ctx *pkgCtx) {
			loadOutsideMethod(ctx, decl)
		})
	case lc.Cursor_TypedefDecl, lc.Cursor_TypeAliasDecl, lc.Cursor_TypeAliasTemplateDecl:
		ctx.addLoadUnit(func(ctx *pkgCtx) {
			loadTypedef(ctx, decl, nil)
		})
	case lc.Cursor_EnumDecl:
		loadEnum(ctx, decl)
	case lc.Cursor_MacroDefinition:
		loadMacro(ctx, decl)
	case lc.Cursor_InclusionDirective:
		loadInclude(ctx, decl)
	case lc.Cursor_Namespace:
		loadNamespace(ctx, scope, decl)
	case lc.Cursor_VarDecl:
		ctx.addLoadUnit(func(ctx *pkgCtx) {
			loadVar(ctx, decl)
		})
	case lc.Cursor_UnionDecl:
		ctx.addLoadUnit(func(ctx *pkgCtx) {
			loadUnion(ctx, decl)
		})
	case lc.Cursor_LinkageSpec: // extern "C" { ... }
		loadLinkageSpec(ctx, scope, decl)
	case lc.Cursor_MacroExpansion, lc.Cursor_StaticAssert, lc.Cursor_UsingDeclaration:
		// noop
	case lc.Cursor_FunctionTemplate:
		// TODO(xsw):
	case lc.Cursor_ClassTemplate, lc.Cursor_ClassTemplatePartialSpecialization:
		preloadTemplateClass(ctx, scope, decl)
	case lc.Cursor_UnexposedDecl, lc.Cursor_UnexposedAttr:
		// noop
	default:
		ctx.panicf(decl, "loadDecl: unknown kind - %v", decl.Kind)
	}
}

func loadNamespace(ctx *pkgCtx, scope *scopeCtx, namespace clang.Cursor) {
	cName := cNameOf(namespace)
	if ctx.isNSIgnored(cName) {
		if debugCompileDecl {
			ctx.logf(namespace, "namespace %s - ignored", cName)
		}
		return
	}
	clang.VisitChildren(namespace, func(decl, parent clang.Cursor) clang.ChildVisitResult {
		loadDecl(ctx, scope, decl)
		return clang.Continue
	})
}

func loadLinkageSpec(ctx *pkgCtx, scope *scopeCtx, linkage clang.Cursor) {
	clang.VisitChildren(linkage, func(decl, parent clang.Cursor) clang.ChildVisitResult {
		loadDecl(ctx, scope, decl)
		return clang.Continue
	})
}

func loadOutsideMethod(ctx *pkgCtx, outsideDecl clang.Cursor) {
	fnUSR := objUSR(outsideDecl)
	if m, ok := ctx.ovobjs[fnUSR]; ok {
		// TODO(xsw): check another definition of the same method
		m.decl = outsideDecl
	} else {
		ctx.ignoref(featExplicitIgnore, outsideDecl, "[WARN] method undeclared - %s", cNameOf(outsideDecl))
	}
}

func newClassCtx(ctx *pkgCtx, cls clang.Cursor, parent *scopeCtx) *classCtx {
	if cls.IsCursorDefinition() == 0 {
		return nil // not a definition
	}
	this := &classCtx{
		decl:      cls,
		overloads: make(map[string]*overloads),
		parent:    parent,
	}
	clang.VisitChildren(cls, func(decl, parent clang.Cursor) clang.ChildVisitResult {
		switch decl.Kind {
		case lc.Cursor_CXXMethod, lc.Cursor_FunctionTemplate,
			lc.Cursor_Constructor, lc.Cursor_Destructor, lc.Cursor_ConversionFunction:
			if !isPublic(decl) {
				break
			}
			var name string
			switch decl.Kind {
			case lc.Cursor_Constructor:
				name = "XGo_Ctor"
			case lc.Cursor_Destructor:
				name = "XGo_Dtor"
			default:
				if decl.CXXMethodIsStatic() != 0 {
					name = cNameOf(decl) // static method name
				} else {
					name = clang.String(decl) // method name
				}
			}
			if decl.CXXMethodIsStatic() != 0 {
				preloadGlobalFunc(ctx, &ctx.scopeCtx, decl)
				break
			}
			fn, isNew := this.addOverloadObj(ctx, name, decl)
			if isNew {
				this.publicMethods = append(this.publicMethods, fn)
			}
		}
		return clang.Continue
	})
	return this
}

func isPublic(decl clang.Cursor) bool {
	return decl.CXXAccessSpecifier() == lc.CXXPublic
}

func preloadClass(ctx *pkgCtx, scope *scopeCtx, cls clang.Cursor) {
	cName := cNameOf(cls)
	if cls.NumTemplateArguments() >= 0 { // note: -1 means not a template class
		scope.addOverloadObj(ctx, cName, cls)
		return
	}
	this := newClassCtx(ctx, cls, scope)
	ctx.addLoadUnit(func(ctx *pkgCtx) {
		loadClass(ctx, cName, this, cls)
	})
}

func preloadTemplateClass(ctx *pkgCtx, scope *scopeCtx, cls clang.Cursor) {
	cName := cNameOf(cls)
	obj, isNew := scope.addOverloadObj(ctx, cName, cls)
	if isNew {
		ctx.addLoadUnit(func(ctx *pkgCtx) {
			this := newClassCtx(ctx, obj.decl, scope)
			loadTemplateClass(ctx, cName, this, obj)
		})
	}
}

func preloadGlobalFunc(ctx *pkgCtx, scope *scopeCtx, fn clang.Cursor) {
	cName := cNameOf(fn)
	obj, isNew := scope.addOverloadObj(ctx, cName, fn)
	if isNew {
		ctx.addLoadUnit(func(ctx *pkgCtx) {
			loadGlobalFunc(ctx, obj)
		})
	}
}

// -----------------------------------------------------------------------------
