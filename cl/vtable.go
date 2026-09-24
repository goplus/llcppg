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
	"strconv"

	"github.com/goplus/llcppg/clang"
	lc "github.com/goplus/llcppg/lib/clang"
)

// -----------------------------------------------------------------------------
//
// vtable-aware bindings for C++ virtual methods.
//
// For every polymorphic class X (one that declares or inherits virtual
// methods) llcppg emits, in addition to the struct that mirrors the C++ data
// layout:
//
//   - a typed vtable struct `_xgo_vtable_X`, annotated `// llgo:type C`, with one
//     field per virtual slot in vtable order. A public method's slot is a
//     function-pointer field whose first parameter is "this *X"; a reserved or
//     non-public slot is instead an unexported "_xgo_slotN unsafe.Pointer"
//     placeholder that only holds the slot's position (see vtableSlot).
//   - an accessor method "func (p *X) XGo_vptr() *_xgo_vtable_X" that reinterprets
//     the pointer stored at offset 0 as the typed vtable.
//
// The vptr slot itself keeps the layout it already had; the only change is that
// the exported field "XGo_vptr" becomes the unexported field vptrName
// ("_xgo_vptr") so it can coexist with the "XGo_vptr()" accessor (Go forbids a
// field and a method sharing a name). Only the class that owns the vptr (i.e.
// has no primary base) declares that field.
//
// See issue goplus/llcppg#754 for the full proposal.

// vtableName returns the Go name of the typed vtable struct for a class named
// clsName, e.g. "Shape" -> "_xgo_vtable_Shape".
func vtableName(clsName string) string {
	return "_xgo_vtable_" + clsName
}

// vtableSlot describes one entry of a class vtable.
//
// A slot is rendered as a named function-pointer field when named is true (its
// Go name is name and its signature comes from decl); otherwise it is rendered
// as an unexported "_xgo_slot<N> unsafe.Pointer" placeholder that only reserves
// the slot so later slots keep their index. A placeholder is used for a slot
// whose method is non-public: a private/protected virtual method still occupies
// a vtable slot, but exposing a callable field for it would leak an inaccessible
// member.
//
// decl is always the backing virtual method cursor (never null); a placeholder
// is distinguished by named == false, not by an absent cursor, so its signature
// is still available if a later pass needs it.
type vtableSlot struct {
	name  string       // Go field name (method name, disambiguated for overloads)
	decl  clang.Cursor // the virtual method cursor providing the signature (never null)
	named bool         // render as a named field (public) vs a placeholder
}

func vtableSlotOf(ctx *pkgCtx, m clang.Cursor, named bool) vtableSlot {
	if !named {
		return vtableSlot{decl: m}
	}
	return vtableSlot{
		name:  vtableMethodName(ctx, m),
		decl:  m,
		named: named,
	}
}

// genVtable emits the "_xgo_vtable_X" struct and the "XGo_vptr()" accessor for a
// polymorphic class. ownsVptr reports whether the class declares its own vptr
// field (no primary base); when false the class shares its primary base's vptr,
// which sits at offset 0.
func genVtable(ctx *pkgCtx, scope *classCtx, ownsVptr bool) {
	slots := vtableSlots(ctx, scope, scope.decl)
	if len(slots) == 0 {
		return
	}
	pkg := ctx.pkg
	pkgTypes := pkg.Types
	clsNamed := scope.typNamed
	clsName := clsNamed.Obj().Name()
	recvPtr := types.NewPointer(clsNamed)

	ctx.forceImportUnsafe()

	// The vtable struct: one field per slot. A named slot is a function pointer
	// whose first parameter is "this *X"; an anonymous slot is an unexported
	// unsafe.Pointer placeholder.
	fields := make([]*types.Var, 0, len(slots))
	for i, slot := range slots {
		var fldType types.Type
		var fldName string
		if slot.named {
			fldName = slot.name
			fldType = vtableSlotFunc(ctx, pkgTypes, recvPtr, slot.decl)
		} else {
			fldName = placeholderSlotName(i)
			fldType = types.Typ[types.UnsafePointer]
		}
		fields = append(fields, types.NewField(goNodePos(ctx, scope.decl), pkgTypes, fldName, fldType, false))
	}
	vtStruct := types.NewStruct(fields, nil)

	// `llgo:type C` keeps the struct laid out exactly as the C++ vtable so the
	// reinterpret cast in the accessor is valid. The leading "\n" renders a
	// blank line before the directive, matching the spacing of the generated
	// method blocks.
	vtDecl := pkg.NewTypeDefs().SetComments(&ast.CommentGroup{
		List: []*ast.Comment{{Text: "\n// llgo:type C"}},
	}).NewType(vtableName(clsName), goNode(ctx, scope.decl))
	vtNamed := vtDecl.InitType(pkg, vtStruct)
	vtPtr := types.NewPointer(vtNamed)

	genVptrAccessor(ctx, recvPtr, vtPtr, ownsVptr)
}

// genVptrAccessor emits "func (p *X) XGo_vptr() *_xgo_vtable_X".
//
// When the class owns its vptr the body is:
//
//	return (*_xgo_vtable_X)(p._xgo_vptr)
//
// otherwise the class shares its primary base's vptr at offset 0, so the body
// reads that pointer directly:
//
//	return (*_xgo_vtable_X)(*(*unsafe.Pointer)(unsafe.Pointer(p)))
func genVptrAccessor(ctx *pkgCtx, recvPtr, vtPtr types.Type, ownsVptr bool) {
	pkg := ctx.pkg
	pkgTypes := pkg.Types
	tyUP := types.Typ[types.UnsafePointer]

	recv := types.NewParam(token.NoPos, pkgTypes, "p", recvPtr)
	results := types.NewTuple(types.NewParam(token.NoPos, pkgTypes, "", vtPtr))
	sig := types.NewSignatureType(recv, nil, nil, nil, results, false)

	f, err := pkg.NewFuncWith(token.NoPos, vptrAccessorName, sig, nil)
	if err != nil {
		panic("genVptrAccessor: " + err.Error())
	}
	cb := f.BodyStart(pkg)
	cb.Typ(vtPtr) // (*_xgo_vtable_X)( ...
	if ownsVptr {
		// p._xgo_vptr
		cb.VarVal("p").MemberVal(vptrName, 1)
	} else {
		// *(*unsafe.Pointer)(unsafe.Pointer(p))
		cb.Typ(types.NewPointer(tyUP)). // (*unsafe.Pointer)(
						Typ(tyUP).VarVal("p").Call(1). //   unsafe.Pointer(p)
						Call(1).                       // )
						Star()                         // *
	}
	cb.Call(1).Return(1).End() // )
}

// vtableSlotFunc builds the function-pointer type of a named vtable slot:
// func(this *X, <params>) <result>, mirroring the method signature but with the
// receiver turned into an explicit leading "this" parameter.
func vtableSlotFunc(ctx *pkgCtx, pkg *types.Package, recvPtr types.Type, fn clang.Cursor) types.Type {
	this := types.NewParam(token.NoPos, pkg, "this", recvPtr)
	rest, variadic := newParams(ctx, pkg, fn)
	params := make([]*types.Var, 0, 1+rest.Len())
	params = append(params, this)
	for i := 0; i < rest.Len(); i++ {
		params = append(params, rest.At(i))
	}
	results := toFuncResults(ctx, pkg, fn.ResultType())
	sig := types.NewSignatureType(nil, nil, nil, types.NewTuple(params...), results, variadic)
	return sig
}

// vtableSlots computes the vtable of cls in ABI (declaration) order.
//
// If cls has a primary base its slots come first (with "this" re-typed to the
// derived class); a virtual method of cls that overrides an inherited method
// takes over that slot instead of adding a new one. New virtual methods follow,
// in declaration order.
//
// Two slot rules beyond a plain public method (see issue goplus/llcppg#754):
//
//   - A virtual destructor occupies two consecutive Itanium slots: "XGo_dtor"
//     (complete-object destructor) followed by "XGo_dtor_deleting". Both have the
//     signature func(this *X).
//   - A non-public (private/protected) virtual method still occupies a slot but
//     is rendered as an unexported placeholder, so later slots keep their index
//     without exposing an inaccessible member.
func vtableSlots(ctx *pkgCtx, scope *classCtx, cls clang.Cursor) []vtableSlot {
	var slots []vtableSlot
	if primary, ok := primaryBase(cls); ok {
		base := primary.Type().TypeDeclaration().Definition()
		slots = vtableSlots(ctx, nil, base)
	}
	for _, m := range ownVirtualMethods(cls) {
		public := m.CXXAccessSpecifier() == lc.CXXPublic
		if m.Kind == lc.CursorDestructor {
			// The destructor's two slots reuse the inherited destructor slots when
			// present (a base virtual destructor always seeds a matching pair), so a
			// derived destructor overrides rather than appends.
			if idx := overriddenSlot(slots, m); idx >= 0 {
				// The complete-object slot and the deleting slot are an inseparable
				// Itanium pair, always seeded contiguously by dtorSlot below, so a
				// matched override must have both. If the deleting slot is missing we
				// would silently leave XGo_dtor_deleting pointing at the stale base
				// entry (wrong "this" type and public flag); fail loudly instead so a
				// broken invariant surfaces rather than emitting an incorrect vtable.
				if idx+1 >= len(slots) {
					panic("vtableSlots: virtual destructor override missing its deleting slot; the Itanium destructor pair is not contiguous")
				}
				setDtorSlot(&slots[idx], dtorSlotName, m, public)
				setDtorSlot(&slots[idx+1], dtorDeletingSlotName, m, public)
				continue
			}
			slots = append(slots,
				dtorSlot(dtorSlotName, m, public),
				dtorSlot(dtorDeletingSlotName, m, public),
			)
			continue
		}
		if idx := overriddenSlot(slots, m); idx >= 0 {
			slots[idx] = vtableSlotOf(ctx, m, public)
			continue
		}
		slots = append(slots, vtableSlotOf(ctx, m, public))
	}
	return slots
}

// dtorSlot builds one of the two slots of a virtual destructor. Both slots take
// the destructor cursor for their func(this *X) signature; only a public
// destructor is exposed as a named field.
func dtorSlot(name string, m clang.Cursor, public bool) vtableSlot {
	return vtableSlot{name: name, decl: m, named: public}
}

// setDtorSlot rewrites an inherited destructor slot in place when a derived
// class re-declares its virtual destructor.
func setDtorSlot(s *vtableSlot, name string, m clang.Cursor, public bool) {
	s.name, s.decl, s.named = name, m, public
}

// ownVirtualMethods returns the virtual methods and virtual destructor declared
// directly by cls (not inherited), in declaration order. Static methods are not
// virtual, so they never appear.
func ownVirtualMethods(cls clang.Cursor) []clang.Cursor {
	var methods []clang.Cursor
	clang.VisitChildren(cls, func(decl, parent clang.Cursor) clang.ChildVisitResult {
		switch decl.Kind {
		case lc.CursorCXXMethod:
			if decl.CXXMethodIsVirtual() != 0 && decl.CXXMethodIsStatic() == 0 {
				methods = append(methods, decl)
			}
		case lc.CursorDestructor:
			if decl.CXXMethodIsVirtual() != 0 {
				methods = append(methods, decl)
			}
		}
		return clang.Continue
	})
	return methods
}

// overriddenSlot returns the index of the inherited slot that m overrides, or
// -1 if m introduces a new slot.
//
// Rather than re-deriving the C++ override rules by matching names and
// parameter types, this asks libclang directly: clang_getOverriddenCursors
// reports the base-class virtual methods a method immediately overrides, and it
// already accounts for signature (including cv-qualifiers, ref-qualifiers and
// covariant returns) the way the Itanium ABI does. libclang only returns the
// *immediate* overrides, so we walk the graph transitively to reach the base
// method that actually seeded the inherited slot.
func overriddenSlot(slots []vtableSlot, m clang.Cursor) int {
	roots := overriddenRoots(m)
	if len(roots) == 0 {
		return -1
	}
	for i, s := range slots {
		// Every slot is seeded with a real backing cursor (see vtableSlot), so a
		// null decl means a slot was built without one — a broken invariant. Assert
		// rather than skip so it surfaces instead of silently missing an override.
		if s.decl.IsNull() != 0 {
			panic("overriddenSlot: vtable slot has a null decl; every slot must carry its backing virtual method cursor")
		}
		for _, r := range roots {
			if s.decl.Equal(r) != 0 {
				return i
			}
		}
	}
	return -1
}

// overriddenRoots returns the transitive set of virtual methods that m
// overrides, walking clang_getOverriddenCursors until it reaches methods with
// no further overrides. The inherited slot list is built from each base's own
// method cursors, so an inherited slot's decl is one of these roots exactly when
// m overrides it.
func overriddenRoots(m clang.Cursor) []clang.Cursor {
	var roots []clang.Cursor
	var walk func(cur clang.Cursor)
	walk = func(cur clang.Cursor) {
		overridden, dispose := clang.OverriddenCursors(cur)
		defer dispose()
		for _, base := range overridden {
			roots = append(roots, base)
			walk(base)
		}
	}
	walk(m)
	return roots
}

func vtableMethodName(ctx *pkgCtx, m clang.Cursor) string {
	manglingName := clang.Mangling(m)
	if fn, ok := ctx.funcs[manglingName]; ok {
		return ctx.funcName(fn.name, fn.order(), false, true)
	}
	panic("vtableMethodName: method not found - " + manglingName)
}

// placeholderSlotName returns the field name of an anonymous (reserved) slot at
// index i, e.g. "_xgo_slot3".
func placeholderSlotName(i int) string {
	return "_xgo_slot" + strconv.Itoa(i)
}

// -----------------------------------------------------------------------------
