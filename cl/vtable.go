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
//   - a typed vtable struct "_xgo_vtable_X", annotated "//llgo:type C", with one
//     function-pointer field per virtual slot in vtable order. The first
//     parameter of every field is "this *X".
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
// A slot is either a named virtual method (decl set, the method cursor whose
// Go name and signature the field takes) or an unexported placeholder that only
// reserves the slot so later slots keep their index (decl null).
type vtableSlot struct {
	name string       // Go field name (method name, disambiguated for overloads)
	decl clang.Cursor // the virtual method cursor providing the signature
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
		if slot.decl.IsNull() == 0 {
			fldName = slot.name
			fldType = vtableSlotFunc(ctx, pkgTypes, recvPtr, slot.decl)
		} else {
			fldName = placeholderSlotName(i)
			fldType = types.Typ[types.UnsafePointer]
		}
		fields = append(fields, types.NewField(goNodePos(ctx, scope.decl), pkgTypes, fldName, fldType, false))
	}
	vtStruct := types.NewStruct(fields, nil)

	// //llgo:type C keeps the struct laid out exactly as the C++ vtable so the
	// reinterpret cast in the accessor is valid. The leading "\n" renders a
	// blank line before the directive, matching the spacing of the generated
	// method blocks.
	vtDecl := pkg.NewTypeDefs().SetComments(&ast.CommentGroup{
		List: []*ast.Comment{{Text: "\n//llgo:type C"}},
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
func vtableSlots(ctx *pkgCtx, scope *classCtx, cls clang.Cursor) []vtableSlot {
	var slots []vtableSlot
	if primary, ok := primaryBase(cls); ok {
		base := primary.Type().TypeDeclaration().Definition()
		slots = vtableSlots(ctx, nil, base)
	}
	own := ownVirtualMethods(cls)
	for _, m := range own {
		if idx := overriddenSlot(slots, m); idx >= 0 {
			slots[idx].name = vtableMethodName(ctx, scope, cls, m)
			slots[idx].decl = m
			continue
		}
		slots = append(slots, vtableSlot{
			name: vtableMethodName(ctx, scope, cls, m),
			decl: m,
		})
	}
	return slots
}

// ownVirtualMethods returns the virtual methods declared directly by cls (not
// inherited), in declaration order.
//
// TODO(#754): a virtual destructor (CursorDestructor with CXXMethodIsVirtual)
// occupies two Itanium slots (XGo_dtor, XGo_dtor_deleting) and non-public
// virtual methods need unexported placeholder slots to keep later indices
// aligned. Those are not emitted yet; a class whose only virtual member is a
// destructor therefore keeps its vptr field but gets no typed vtable, which is
// safe (no misaligned slots) though not yet friendly.
func ownVirtualMethods(cls clang.Cursor) []clang.Cursor {
	var methods []clang.Cursor
	clang.VisitChildren(cls, func(decl, parent clang.Cursor) clang.ChildVisitResult {
		switch decl.Kind {
		case lc.CursorCXXMethod:
			if decl.CXXMethodIsVirtual() != 0 && decl.CXXMethodIsStatic() == 0 {
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
		if s.decl.IsNull() != 0 {
			continue
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
		for _, base := range clang.OverriddenCursors(cur) {
			roots = append(roots, base)
			walk(base)
		}
	}
	walk(m)
	return roots
}

// vtableMethodName returns the Go field name of a virtual method slot. When the
// method belongs to the class currently being compiled (scope != nil and the
// method is one of scope's registered funcs) the already-computed overloaded Go
// name is reused; otherwise the name is derived the same way loadClassMember
// derives it.
func vtableMethodName(ctx *pkgCtx, scope *classCtx, cls clang.Cursor, m clang.Cursor) string {
	if scope != nil {
		if fn, ok := ctx.funcs[clang.Mangling(m)]; ok {
			name, _ := ctx.getPubName(fn.name, fn.order())
			return name
		}
	}
	name, _ := ctx.getPubName(clang.String(m), -1)
	return name
}

// placeholderSlotName returns the field name of an anonymous (reserved) slot at
// index i, e.g. "_xgo_slot3".
func placeholderSlotName(i int) string {
	return "_xgo_slot" + strconv.Itoa(i)
}

// -----------------------------------------------------------------------------
