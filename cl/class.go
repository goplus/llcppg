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
	"log"

	"github.com/goplus/llcppg/clang"
	lc "github.com/goplus/llcppg/lib/clang"
)

// -----------------------------------------------------------------------------

// vptrName is the field name of the implicit vptr introduced by a polymorphic
// (has-virtual-methods) C++ class. It is a pointer-sized slot placed at the
// very beginning of the type layout, mirroring the C++ Itanium ABI.
//
// The field is unexported so it can coexist with the exported "XGo_vptr()"
// accessor method that returns the typed vtable (Go forbids a field and a
// method sharing a name). See vtable.go and issue goplus/llcppg#754.
const vptrName = "_xgo_vptr"

// vptrAccessorName is the exported method that returns the typed vtable for a
// polymorphic class, e.g. "func (p *X) XGo_vptr() *_xgo_vtable_X". It must
// differ from vptrName so the method and the field can coexist.
const vptrAccessorName = "XGo_vptr"

// dtorSlotName and dtorDeletingSlotName are the vtable field names of the two
// consecutive Itanium slots a virtual destructor occupies: the complete-object
// destructor followed by the deleting destructor. Both have the signature
// func(this *X). See vtable.go and issue goplus/llcppg#754.
//
// The lowercase "dtor" here is deliberately distinct from the "XGo_Dtor" method
// generated in loadClassMember (uppercase "D"): these are vtable slot fields,
// not the destructor method, so the differing case is intentional and not a typo.
const (
	dtorSlotName         = "XGo_dtor"
	dtorDeletingSlotName = "XGo_dtor_deleting"
)

type classCtx struct {
	scopeCtx
	decl          clang.Cursor
	typNamed      *types.Named
	fields        []*types.Var
	publicMethods []*funcObj
	inPublic      bool
	polymorphic   bool // declares or inherits virtual methods
	ownsVptr      bool // owns the vptr field (polymorphic with no primary base)
}

func compileClass(ctx *pkgCtx, scope *classCtx) {
	scope.reorder()
	if scope.polymorphic {
		// Emit the typed vtable and the XGo_vptr() accessor once method
		// overload ordering is finalized (reorder above), so vtable field
		// names match the generated method names.
		genVtable(ctx, scope, scope.ownsVptr)
	}
	for _, method := range scope.publicMethods {
		compileFuncOrMethod(ctx, method, scope)
	}
}

func loadClass(ctx *pkgCtx, cls clang.Cursor, ns string, defaultInPublic bool) {
	pkg := ctx.pkg
	pkgTypes := pkg.Types
	origName := ns + clang.String(cls)
	if debugCompileDecl {
		log.Println("class", origName)
	}
	clsName, rewritten := ctx.getPubName(origName, -1)
	typDecl := pkg.NewTypeDefs().NewType(clsName, goNode(ctx, cls))
	typNamed := typDecl.Type()
	if rewritten {
		substObj(pkgTypes, pkgTypes.Scope(), origName, typNamed.Obj())
	}
	ctx.types[clang.String(cls.Type())] = typNamed.Obj()
	scope := &classCtx{
		decl:      cls,
		typNamed:  typNamed,
		overloads: make(map[string]*overloads),
		inPublic:  defaultInPublic,
	}
	clang.VisitChildren(cls, func(decl, parent clang.Cursor) clang.ChildVisitResult {
		loadClassMember(ctx, pkgTypes, scope, origName, decl)
		return clang.Continue
	})
	// A virtual base subobject exists exactly once in the most-derived object and,
	// per the C++ Itanium ABI, is laid out after all non-virtual data. Append the
	// transitive set of distinct virtual bases as embedded fields at the tail of
	// the layout, deduplicated, so a class that virtually derives from a base keeps
	// a single shared base subobject at its tail. This runs after the non-virtual
	// members are collected above so the virtual bases follow them.
	//
	// This models the *complete-object* layout of a class that virtually inherits
	// (e.g. istream / ostream): vptr, own fields, then the shared virtual base. It
	// does NOT yet model the diamond *join* (e.g. iostream itself), where a class
	// combines two non-virtual bases that each carry the same virtual base: there
	// the ABI hoists the virtual base out of each base subobject into a single copy
	// at the most-derived tail, so the base subobject layout differs from the
	// base's complete-object layout and cannot be produced by embedding the base's
	// Go struct as-is. Reject that case loudly rather than emit a struct whose size
	// is wrong (it would double-store the shared base). See issue goplus/llcppg#759.
	if base, ok := nonVirtualBaseWithVirtualBase(cls); ok {
		panic("todo: diamond virtual base join is not supported yet (non-virtual base '" +
			clang.String(base) + "' itself has a virtual base); see goplus/llcppg#759")
	}
	appendVirtualBases(ctx, pkgTypes, scope, cls)
	// Establish the layout at offset 0, following the C++ Itanium ABI. A
	// polymorphic class shares its vptr with its primary base (the first
	// non-virtual *polymorphic* direct base in declaration order); that base is
	// laid out first so the shared vptr sits at offset 0. If the class is
	// polymorphic but has no such base to reuse, it introduces its own implicit
	// vptr at the start of the layout instead. A class with virtual bases is
	// itself polymorphic — it needs a vptr to locate the virtual-base subobjects —
	// so it owns a vptr when it has no non-virtual polymorphic primary base to
	// reuse one from.
	if primary, ok := primaryBase(cls); ok {
		hoistPrimaryBase(ctx, scope, primary)
		scope.polymorphic = true
	} else if isPolymorphic(cls) {
		ctx.forceImportUnsafe()
		vptr := types.NewField(goNodePos(ctx, cls), pkgTypes, vptrName, types.Typ[types.UnsafePointer], false)
		scope.fields = append([]*types.Var{vptr}, scope.fields...)
		scope.polymorphic = true
		scope.ownsVptr = true
	}
	typStruc := types.NewStruct(scope.fields, nil)
	typDecl.InitType(pkg, typStruc)
	ctx.compiles = append(ctx.compiles, func(ctx *pkgCtx) {
		compileClass(ctx, scope)
	})
}

func loadClassMember(ctx *pkgCtx, pkg *types.Package, cls *classCtx, origName string, decl clang.Cursor) {
	switch decl.Kind {
	case lc.CursorCXXMethod, lc.CursorConstructor, lc.CursorDestructor:
		var name string
		switch decl.Kind {
		case lc.CursorConstructor:
			name = "XGo_Ctor"
		case lc.CursorDestructor:
			name = "XGo_Dtor"
		default:
			// A static method is not a method: it has no implicit "this". Register
			// it as a global function whose name is prefixed by the enclosing class
			// name (the class name acts like a namespace); it is compiled with a
			// nil class in compileClass.
			if decl.CXXMethodIsStatic() != 0 {
				if cls.inPublic {
					loadGlobalFunc(ctx, &ctx.scopeCtx, decl, origName+"_")
				}
				return
			}
			name = clang.String(decl)
		}
		isPublic := cls.inPublic
		if isPublic {
			if fn, ok := cls.addFunc(ctx, name, decl); ok {
				cls.publicMethods = append(cls.publicMethods, fn)
			}
		}

	case lc.CursorFieldDecl:
		origName := clang.String(decl)
		fldType := toType(ctx, pkg, decl.Type(), flagIsStructField)
		fldName := origName
		if cls.inPublic {
			fldName, _ = ctx.getPubName(origName, -1)
		}
		fld := types.NewField(goNodePos(ctx, decl), pkg, fldName, fldType, false)
		cls.fields = append(cls.fields, fld)

	case lc.CursorVarDecl:
		// A static member variable is not a field: it has no per-instance
		// storage. Load it as a package-level variable whose name is prefixed
		// by the enclosing class name (the class name acts like a namespace),
		// mirroring how a static method is handled above.
		if cls.inPublic {
			loadVar(ctx, decl, origName+"_")
		}

	case lc.CursorCXXAccessSpecifier:
		cls.inPublic = decl.CXXAccessSpecifier() == lc.CXXPublic

	case lc.CursorEnumDecl:
		// An enum nested in a class only affects naming: its constants are
		// emitted as global consts prefixed by the enclosing class name (the
		// class name acts like a namespace), e.g. Color_Red.
		if cls.inPublic {
			loadEnum(ctx, decl, origName+"_")
		}

	case lc.CursorCXXBaseSpecifier:
		// A non-virtual base class is treated the same as a member variable (field)
		// - simply an embedded one at its declaration position. A virtual base is
		// instead laid out once at the tail of the most-derived object (see
		// appendVirtualBases), so it is skipped here.
		if decl.IsVirtualBase() != 0 {
			return
		}
		base := baseClass(ctx, decl)
		fld := types.NewField(goNodePos(ctx, decl), pkg, base.Name(), base.Type(), true)
		cls.fields = append(cls.fields, fld)

	default:
		log.Panicln("loadClassMember: unknown kind =", decl.Kind)
	}
}

// primaryBase returns the base-specifier cursor of the primary base class of
// cls, if any. Per the C++ Itanium ABI, the primary base is the first
// non-virtual *dynamic (polymorphic)* direct base in declaration order;
// non-polymorphic bases are skipped rather than disqualifying a later
// polymorphic one. A class shares its vptr (at offset 0) with its primary base,
// so no fresh vptr is introduced when one exists. This covers the four cases:
//   - No base, own virtual methods: no primary base -> introduces a vptr.
//   - Base without virtual methods + own virtual methods: no polymorphic base
//     to reuse -> introduces a vptr.
//   - Single polymorphic base: it is the primary base -> reuses its vptr.
//   - Multiple bases where an earlier one is non-polymorphic but a later one is
//     polymorphic: the polymorphic base is the primary base -> reuses its vptr
//     (and is laid out first so the shared vptr stays at offset 0).
func primaryBase(cls clang.Cursor) (spec clang.Cursor, ok bool) {
	clang.VisitChildren(cls, func(decl, parent clang.Cursor) clang.ChildVisitResult {
		if decl.Kind == lc.CursorCXXBaseSpecifier && decl.IsVirtualBase() == 0 {
			if b := decl.Type().TypeDeclaration().Definition(); b.IsNull() == 0 && isPolymorphic(b) {
				spec, ok = decl, true
				return clang.Break
			}
		}
		return clang.Continue
	})
	return
}

// hoistPrimaryBase moves the embedded field of the given primary base to the
// front of scope.fields, so its shared vptr sits at offset 0 (the ABI lays the
// primary base out first, regardless of its declaration order among bases).
func hoistPrimaryBase(ctx *pkgCtx, scope *classCtx, spec clang.Cursor) {
	name := baseClass(ctx, spec).Name()
	for i, f := range scope.fields {
		if f.Embedded() && f.Name() == name {
			if i != 0 {
				rest := make([]*types.Var, 0, len(scope.fields))
				rest = append(rest, scope.fields[:i]...)
				rest = append(rest, scope.fields[i+1:]...)
				scope.fields = append([]*types.Var{f}, rest...)
			}
			return
		}
	}
}

// appendVirtualBases appends the transitive set of distinct virtual base classes
// of cls as embedded fields at the tail of scope.fields, deduplicated so each
// virtual base subobject appears exactly once.
//
// Per the C++ Itanium ABI a virtual base is laid out once, after all non-virtual
// data of the most-derived object; this is what makes the diamond pattern (e.g.
// C++'s iostream, where istream and ostream both virtually derive from basic_ios)
// share a single base subobject instead of holding two independent copies. The
// bases are appended in the order collectVirtualBases discovers them (a
// depth-first walk of the base graph), matching the ABI's virtual-base ordering.
func appendVirtualBases(ctx *pkgCtx, pkg *types.Package, scope *classCtx, cls clang.Cursor) {
	seen := make(map[string]bool)
	for _, spec := range collectVirtualBases(cls, seen) {
		base := baseClass(ctx, spec)
		fld := types.NewField(goNodePos(ctx, spec), pkg, base.Name(), base.Type(), true)
		scope.fields = append(scope.fields, fld)
	}
}

// nonVirtualBaseWithVirtualBase reports the first non-virtual direct base of cls
// that itself has (directly or transitively) a virtual base. Embedding such a
// base's complete-object Go struct would carry a copy of the shared virtual base
// inside it, so appending the same virtual base again at the tail double-stores
// it — the diamond join the ABI resolves by hoisting a single copy to the most-
// derived object. This is not supported yet (see the caller in loadClass).
func nonVirtualBaseWithVirtualBase(cls clang.Cursor) (base clang.Cursor, ok bool) {
	clang.VisitChildren(cls, func(decl, parent clang.Cursor) clang.ChildVisitResult {
		if decl.Kind == lc.CursorCXXBaseSpecifier && decl.IsVirtualBase() == 0 {
			if b := decl.Type().TypeDeclaration().Definition(); b.IsNull() == 0 && hasVirtualBase(b) {
				base, ok = decl, true
				return clang.Break
			}
		}
		return clang.Continue
	})
	return
}

// hasVirtualBase reports whether cls has a virtual base directly or through any
// of its (virtual or non-virtual) bases.
func hasVirtualBase(cls clang.Cursor) bool {
	found := false
	clang.VisitChildren(cls, func(decl, parent clang.Cursor) clang.ChildVisitResult {
		if decl.Kind == lc.CursorCXXBaseSpecifier {
			if decl.IsVirtualBase() != 0 {
				found = true
				return clang.Break
			}
			if b := decl.Type().TypeDeclaration().Definition(); b.IsNull() == 0 && hasVirtualBase(b) {
				found = true
				return clang.Break
			}
		}
		return clang.Continue
	})
	return found
}

// collectVirtualBases returns the base-specifier cursors of the distinct virtual
// bases reachable from cls, in depth-first declaration order. seen tracks the
// fully-qualified name of each virtual base already collected so a base shared
// through several inheritance paths is emitted only once.
//
// Only direct virtual bases and, recursively, the virtual bases of a virtual
// base are collected. Virtual bases reached through a *non-virtual* base are not
// handled here: loadClass rejects that (the diamond join) up front, because
// embedding the non-virtual base's complete-object struct would already store a
// copy of the shared virtual base.
func collectVirtualBases(cls clang.Cursor, seen map[string]bool) []clang.Cursor {
	var specs []clang.Cursor
	clang.VisitChildren(cls, func(decl, parent clang.Cursor) clang.ChildVisitResult {
		if decl.Kind != lc.CursorCXXBaseSpecifier || decl.IsVirtualBase() == 0 {
			return clang.Continue
		}
		name := fullName(decl)
		if !seen[name] {
			seen[name] = true
			specs = append(specs, decl)
		}
		// A virtual base can itself declare virtual bases; recurse into it to hoist
		// those up as well (still deduplicated through seen).
		if base := decl.Type().TypeDeclaration().Definition(); base.IsNull() == 0 {
			specs = append(specs, collectVirtualBases(base, seen)...)
		}
		return clang.Continue
	})
	return specs
}

// isPolymorphic reports whether the class declares or inherits any virtual
// method, or has a virtual base — any of which means it carries (or shares) a
// vptr in its layout.
func isPolymorphic(cls clang.Cursor) bool {
	found := false
	clang.VisitChildren(cls, func(decl, parent clang.Cursor) clang.ChildVisitResult {
		switch decl.Kind {
		case lc.CursorCXXMethod, lc.CursorDestructor:
			if decl.CXXMethodIsVirtual() != 0 {
				found = true
				return clang.Break
			}
		case lc.CursorCXXBaseSpecifier:
			// A virtual base requires a vptr to locate the shared base subobject, so
			// the class is polymorphic regardless of whether that base has any virtual
			// methods. A non-virtual base only contributes polymorphism when it is
			// itself polymorphic.
			if decl.IsVirtualBase() != 0 {
				found = true
				return clang.Break
			}
			if b := decl.Type().TypeDeclaration().Definition(); b.IsNull() == 0 && isPolymorphic(b) {
				found = true
				return clang.Break
			}
		}
		return clang.Continue
	})
	return found
}

func baseClass(ctx *pkgCtx, decl clang.Cursor) *types.TypeName {
	t := decl.Type()
	// Depending on the libclang version, a base specifier's type may be reported
	// as an elaborated type (e.g. "struct Base") rather than the bare record;
	// unwrap it so the record lookup below works in both cases.
	if t.Kind == lc.TypeElaborated {
		t = t.NamedType()
	}
	switch t.Kind {
	case lc.TypeRecord:
		cName := fullName(decl)
		if t, ok := ctx.typeObj(cName); ok {
			return t
		}
	}
	panic("baseClass: unknown base class type - " + clang.String(t))
}

func loadOutsideMethod(ctx *pkgCtx, outsideDecl clang.Cursor) {
	manglingName := clang.Mangling(outsideDecl)
	if m, ok := ctx.funcs[manglingName]; ok {
		m.decl = outsideDecl
	} else {
		log.Panicln("method undeclared -", clang.DisplayName(outsideDecl))
	}
}

// -----------------------------------------------------------------------------
