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
	// Establish the layout at offset 0, following the C++ Itanium ABI. A
	// polymorphic class shares its vptr with its primary base (the first
	// non-virtual *polymorphic* direct base in declaration order); that base is
	// laid out first so the shared vptr sits at offset 0. If the class is
	// polymorphic but has no such base to reuse, it introduces its own implicit
	// vptr at the start of the layout instead.
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
		var fldType types.Type
		var anonymous bool
		var ft = decl.Type()
		if ft.Kind == lc.TypeRecord {
			if ftd := ft.TypeDeclaration(); ftd.IsAnonymous() != 0 {
				fldType, anonymous = emitUnion(ctx, ftd, ctx.nextAnonUnionName()), true
			}
		}
		if !anonymous {
			fldType = toType(ctx, pkg, ft, flagIsStructField)
		}
		origName := clang.String(decl)
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

	case lc.CursorTypedefDecl:
		// A typedef nested in a class acts like one nested in a namespace: it
		// only affects naming, so it is emitted as a package-level type alias
		// prefixed by the enclosing class name (the class name acts like a
		// namespace), e.g. Bar_iterator.
		if cls.inPublic {
			loadTypedef(ctx, decl, origName+"_")
		}

	case lc.CursorCXXBaseSpecifier:
		// A base class is treated the same as a member variable (field) - simply
		// an embedded one. Virtual base classes are not supported for now.
		if decl.IsVirtualBase() != 0 {
			panic("todo: virtual base class is not supported")
		}
		base := baseClass(ctx, decl)
		fld := types.NewField(goNodePos(ctx, decl), pkg, base.Name(), base.Type(), true)
		cls.fields = append(cls.fields, fld)

	case lc.CursorClassDecl, lc.CursorStructDecl:
		switch {
		case decl.IsAnonymousRecordDecl() != 0:
			// TODO(xsw):
		case decl.IsAnonymous() != 0:
			// noop
		default:
			if cls.inPublic {
				defaultInPublic := decl.Kind == lc.CursorStructDecl
				loadClass(ctx, decl, origName+"_", defaultInPublic)
			}
		}
	case lc.CursorUnionDecl:
		switch {
		case decl.IsAnonymousRecordDecl() != 0:
			hoisted := emitUnion(ctx, decl, ctx.nextAnonUnionName())
			fld := types.NewField(goNodePos(ctx, decl), pkg, hoisted.Obj().Name(), hoisted, true)
			cls.fields = append(cls.fields, fld)
		case decl.IsAnonymous() != 0:
			// noop
		default:
			if cls.inPublic {
				loadUnion(ctx, decl, origName+"_")
			}
		}

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

// isPolymorphic reports whether the class declares or inherits any virtual
// method, i.e. whether it has (or shares) a vptr in its layout.
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
		cName := clang.String(decl.Type())
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
