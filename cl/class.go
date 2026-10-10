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
	"strconv"
	"unsafe"

	"github.com/goplus/llcppg/clang"
	lc "github.com/llarhub/clang-c"
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
// polymorphic class, e.g. "func (p *X) XGo_vptr() *X_vtable_XXX". It must
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
	publicMethods []*overloadObj
	polymorphic   bool // declares or inherits virtual methods
	ownsVptr      bool // owns the vptr field (polymorphic with no primary base)
}

func (p *classCtx) scope() *scopeCtx {
	return (*scopeCtx)(unsafe.Pointer(p))
}

func compileClassImpl(ctx *pkgCtx, this *classCtx) {
	this.reorder()
	if this.polymorphic {
		// Emit the typed vtable and the XGo_vptr() accessor once method
		// overload ordering is finalized (reorder above), so vtable field
		// names match the generated method names.
		genVtable(ctx, this, this.ownsVptr)
	}
	for _, method := range this.publicMethods {
		compileFuncOrMethod(ctx, method, this)
	}
}

// -----------------------------------------------------------------------------

func newTypeParams(ctx *pkgCtx, pkg *types.Package, cls clang.Cursor) (ret []*types.TypeParam, quietIgnore bool) {
	idx := 0
	ret = make([]*types.TypeParam, 0, 2)
	clang.VisitChildren(cls, func(decl, parent clang.Cursor) clang.ChildVisitResult {
		switch decl.Kind {
		case lc.Cursor_TemplateTypeParameter:
			if isParameterPack(decl) {
				quietIgnore = true
				return clang.Break
			}
			idx++
			name := clang.String(decl)
			if name == "" {
				name = "_llcppg_tparam" + strconv.Itoa(idx)
			}
			objName := types.NewTypeName(0, pkg, name, nil)
			ret = append(ret, types.NewTypeParam(objName, ctx.any()))
		case lc.Cursor_NonTypeTemplateParameter, lc.Cursor_TemplateTemplateParameter:
			quietIgnore = true
			fallthrough
		case lc.Cursor_CXXMethod, lc.Cursor_FunctionTemplate, lc.Cursor_Constructor,
			lc.Cursor_Destructor, lc.Cursor_ConversionFunction, lc.Cursor_FieldDecl,
			lc.Cursor_CXXBaseSpecifier, lc.Cursor_TypedefDecl, lc.Cursor_TypeAliasDecl,
			lc.Cursor_TypeAliasTemplateDecl, lc.Cursor_VarDecl, lc.Cursor_ClassDecl,
			lc.Cursor_CXXAccessSpecifier, lc.Cursor_FriendDecl, lc.Cursor_UsingDeclaration,
			lc.Cursor_EnumDecl, lc.Cursor_StructDecl:
			return clang.Break
		}
		return clang.Continue
	})
	return
}

func isParameterPack(decl clang.Cursor) bool { // <class... T>
	tu := clang.TU(decl)
	extent := decl.Extent()
	tokens, dispose := tu.Tokenize(extent)
	defer dispose()
	for _, token := range tokens {
		if token.Kind() == lc.Token_Punctuation && tu.TokenSpelling(token) == "..." {
			return true
		}
	}
	return false
}

func loadTemplateClass(ctx *pkgCtx, cName string, this *classCtx, obj *overloadObj) {
	cls := obj.decl
	if ctx.isConfTypeIgnored(cName) {
		ctx.ignoref(featQuietIgnore, cls, "template class %s: ignored by config", cName)
		ctx.ignoreType(cName, featQuietIgnore)
		return
	}

	order := obj.order()
	if debugCompileDecl {
		ctx.logf(cls, "template class %s: order - %d", cName, order)
	}

	pkg := ctx.pkg
	pkgTypes := pkg.Types
	tparams, quietIgnore := newTypeParams(ctx, pkgTypes, cls)
	if quietIgnore || order >= 0 {
		ctx.ignoref(featQuietIgnore, cls, "class %s: unsupported template params, ignored", cName)
		ctx.ignoreType(cName, featQuietIgnore)
		return
	}

	var typDecl, ok = ctx.typdecls[cName]
	if !ok {
		if cls.IsCursorDefinition() == 0 || cls.NumTemplateArguments() > 0 {
			return
		}
		goName := ctx.typeName(cName, true)
		typDecl = newType(ctx, cls, cName, goName, tparams...)
		ctx.typdecls[cName] = typDecl
	}

	if this == nil {
		ctx.ignoref(featQuietIgnore, cls, "class %s: no definition, ignored", cName)
		return
	}

	goName := typDecl.Type().Obj().Name()
	initClassType(ctx, typDecl, this, cName, goName, tparams)
}

// -----------------------------------------------------------------------------

func loadClass(ctx *pkgCtx, cName string, this *classCtx, cls clang.Cursor) {
	if debugCompileDecl {
		ctx.logf(cls, "%s", tagStrvals[cls.Kind]+cName)
	}

	var typDecl, ok = ctx.typdecls[cName]
	if !ok {
		goName := ctx.typeName(cName, true)
		typDecl = newType(ctx, cls, cName, goName)
		ctx.typdecls[cName] = typDecl
	}

	if this == nil {
		return // declaration only, no definition
	}

	goName := typDecl.Type().Obj().Name()
	initClassType(ctx, typDecl, this, cName, goName, nil)
}

// -----------------------------------------------------------------------------

func newType(ctx *pkgCtx, cls clang.Cursor, cName, goName string, tparams ...*types.TypeParam) (ret typDecl) {
	ret.defs = ctx.pkg.NewTypeDefs()
	ret.TypeDecl = ret.defs.NewType(goName, tparams, goNode(ctx, cls))
	if cName != "" {
		ctx.addType(cName, cls, ret.Type())
	}
	return
}

func initClassType(ctx *pkgCtx, typDecl typDecl, this *classCtx, cName, goName string, tparams []*types.TypeParam) bool {
	feats := 0
	initClassTypeEx(ctx, typDecl, this, goName, tparams, &feats)
	if feats&featAllIgnore != 0 {
		ctx.ignoref(feats, this.decl, "class %s: unsupported features, ignored", cName)
		ctx.ignoreType(cName, feats)
		typDecl.Delete()
		return false
	}
	return true
}

func initClassTypeEx(ctx *pkgCtx, typDecl typDecl, this *classCtx, goName string, tparams []*types.TypeParam, feats *int) {
	cls := this.decl

	// Attach the doc at the TypeDefs (GenDecl) level rather than on the
	// TypeSpec; see the note in loadTypedef for why a spec-level doc renders as
	// "type// doc" here.
	if doc := ctx.docCommentGroup(cls); doc != nil {
		typDecl.defs.SetComments(doc)
	}

	this.typNamed = typDecl.Type()
	this.tparams = tparams

	pkg := ctx.pkg
	pkgTypes := pkg.Types

	// A record that declares bit-fields cannot be laid out field-by-field in the
	// usual way: several bit-fields share bytes, so each run is stored as a byte
	// array that preserves the C size, alignment and member offsets, and every
	// named bit-field is exposed through an XGof_get_/XGof_set_ accessor pair
	// (see bitfield.go and issue goplus/llcppg#770). This path handles plain
	// C-style records (no bases, not polymorphic, no type params); anything more
	// complex falls through to the ordinary conversion below.
	if len(tparams) == 0 && !isPolymorphic(cls) && !hasBaseOrNestedField(cls) && hasBitField(cls) {
		if initBitFieldType(ctx, typDecl, this, cls) {
			return
		}
	}

	clang.VisitChildren(cls, func(decl, parent clang.Cursor) clang.ChildVisitResult {
		loadClassMember(ctx, pkgTypes, this, goName, decl, feats)
		return clang.Continue
	})
	if *feats&featAllIgnore != 0 {
		return
	}
	// A record that holds a callback (function pointer) field must carry the
	// "// llgo:type C" directive so those callbacks use the C calling
	// convention, mirroring how function-pointer typedefs are handled (see
	// defineTypedef). featHasCallback is set by toTypeEx while visiting the
	// members above, so it is only known now; upgrade the doc-only comment set
	// at the top of this function to the directive form.
	if *feats&featHasCallback != 0 {
		typDecl.defs.SetComments(ctx.directiveTypeC(cls, true))
	}
	// Establish the layout at offset 0, following the C++ Itanium ABI. A
	// polymorphic class shares its vptr with its primary base (the first
	// non-virtual *polymorphic* direct base in declaration order); that base is
	// laid out first so the shared vptr sits at offset 0. If the class is
	// polymorphic but has no such base to reuse, it introduces its own implicit
	// vptr at the start of the layout instead.
	if primary, ok := primaryBase(cls); ok {
		if !hoistPrimaryBase(ctx, this, primary, feats) {
			return
		}
		this.polymorphic = true
	} else if isPolymorphic(cls) {
		voidptr := ctx.unsafePointer()
		vptr := types.NewField(goNodePos(ctx, cls), pkgTypes, vptrName, voidptr, false)
		this.fields = append([]*types.Var{vptr}, this.fields...)
		this.polymorphic = true
		this.ownsVptr = true
	}
	typStruc := types.NewStruct(this.fields, nil)
	typDecl.InitType(pkg, typStruc)
	ctx.addCompileUnit(func(ctx *pkgCtx) {
		compileClassImpl(ctx, this)
	})
}

func emitClass(ctx *pkgCtx, cls clang.Cursor, goName string, parent *scopeCtx) *types.Named {
	typDecl := newType(ctx, cls, "", goName)
	this := &classCtx{
		decl:      cls,
		overloads: make(map[string]*overloads),
		parent:    parent,
	}
	if !initClassType(ctx, typDecl, this, "", goName, nil) {
		ctx.panicf(cls, "class %s: unsupported feature, failed to initialize class", goName)
	}
	return typDecl.Type()
}

func loadClassMember(ctx *pkgCtx, pkg *types.Package, this *classCtx, goName string, decl clang.Cursor, feats *int) {
	switch decl.Kind {
	case lc.Cursor_CXXMethod, lc.Cursor_FunctionTemplate,
		lc.Cursor_Constructor, lc.Cursor_Destructor, lc.Cursor_ConversionFunction:
		// noop: have been preloaded in newClassCtx

	case lc.Cursor_FieldDecl:
		var fldType types.Type
		var ft = decl.Type()
		ftd := ft.Declaration()
		switch {
		case ftd.IsAnonymous() == 0:
			if fldType = toTypeEx(ctx, pkg, ft, flagIsVarDef, feats, this.scope()); *feats&featAllIgnore != 0 {
				return
			}
		case ftd.Kind == lc.Cursor_UnionDecl:
			fldType = toType(ctx, pkg, ft, flagIsVarDef, this.scope())
		case ft.Kind == lc.Type_Record:
			fldType = emitClass(ctx, ftd, ctx.nextAnonName(), this.scope())
		case ft.Kind == lc.Type_Enum:
			fldType = emitEnum(ctx, ftd, ctx.nextAnonName())
		default:
			ctx.panicf(ftd, "unknown anonymous field type (%d: %s)", ft.Kind, clang.String(ft))
		}
		origName := clang.String(decl)
		fldName := ctx.fieldName(origName, isPublic(decl))
		fld := types.NewField(goNodePos(ctx, decl), pkg, fldName, fldType, false)
		this.fields = append(this.fields, fld)

	case lc.Cursor_VarDecl:
		if isPublic(decl) {
			loadVar(ctx, decl)
		}

	case lc.Cursor_CXXAccessSpecifier, lc.Cursor_FriendDecl,
		lc.Cursor_StaticAssert, lc.Cursor_UsingDeclaration:
		// noop

	case lc.Cursor_EnumDecl:
		// An enum nested in a class only affects naming: its constants are
		// emitted as global consts prefixed by the enclosing class name (the
		// class name acts like a namespace), e.g. Color_Red.
		if isPublic(decl) {
			loadEnum(ctx, decl)
		}

	case lc.Cursor_TypedefDecl, lc.Cursor_TypeAliasDecl, lc.Cursor_TypeAliasTemplateDecl:
		// A typedef nested in a class acts like one nested in a namespace: it
		// only affects naming, so it is emitted as a package-level type alias
		// prefixed by the enclosing class name (the class name acts like a
		// namespace), e.g. Bar_iterator.
		if isPublic(decl) {
			loadTypedef(ctx, decl, this.scope())
		}

	case lc.Cursor_CXXBaseSpecifier:
		// A base class is treated the same as a member variable (field) - simply
		// an embedded one. Virtual base classes are not supported for now.
		if decl.IsVirtualBase() != 0 {
			panic("todo: virtual base class is not supported")
		}
		if typ, name, ok := baseClass(ctx, this, decl, feats); ok && typ != tyIgnore {
			_, isTypeParam := typ.(*types.TypeParam) // typeParam can't be embedded
			fld := types.NewField(goNodePos(ctx, decl), pkg, name, typ, !isTypeParam)
			this.fields = append(this.fields, fld)
		}

	case lc.Cursor_ClassDecl, lc.Cursor_StructDecl:
		switch {
		case decl.IsAnonymousRecordDecl() != 0:
			hoisted := emitClass(ctx, decl, ctx.nextAnonName(), this.scope())
			fld := types.NewField(goNodePos(ctx, decl), pkg, hoisted.Obj().Name(), hoisted, true)
			this.fields = append(this.fields, fld)
		case decl.IsAnonymous() != 0:
			// noop
		default:
			nested := newClassCtx(ctx, decl, this.scope())
			loadClass(ctx, cNameOf(decl), nested, decl)
		}

	case lc.Cursor_ClassTemplate:
		nested := newClassCtx(ctx, decl, this.scope())
		loadTemplateClass(ctx, cNameOf(decl), nested, &overloadObj{
			decl: decl,
		})

	case lc.Cursor_ClassTemplatePartialSpecialization:
		// noop

	case lc.Cursor_UnionDecl:
		switch {
		case decl.IsAnonymousRecordDecl() != 0:
			hoisted := emitUnion(ctx, decl, ctx.nextAnonName())
			fld := types.NewField(goNodePos(ctx, decl), pkg, hoisted.Obj().Name(), hoisted, true)
			this.fields = append(this.fields, fld)
		case decl.IsAnonymous() != 0:
			// noop
		default:
			// A named nested union is emitted at package level for the same
			// reason as a named nested class/struct above: a field may use it
			// as its type even when declared in a private section.
			loadUnion(ctx, decl)
		}

	case lc.Cursor_TemplateTypeParameter, lc.Cursor_NonTypeTemplateParameter,
		lc.Cursor_TemplateTemplateParameter, lc.Cursor_TypeRef:
		// noop

	default:
		if decl.Kind >= lc.Cursor_FirstAttr && decl.Kind <= lc.Cursor_LastAttr {
			return // noop
		}
		ctx.panicf(decl, "class %s: unknown child node kind - %v", goName, decl.Kind)
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
		if decl.Kind == lc.Cursor_CXXBaseSpecifier && decl.IsVirtualBase() == 0 {
			if b := decl.Type().Declaration().Definition(); b.IsNull() == 0 && isPolymorphic(b) {
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
func hoistPrimaryBase(ctx *pkgCtx, this *classCtx, spec clang.Cursor, feats *int) bool {
	if _, name, ok := baseClass(ctx, this, spec, feats); ok {
		for i, f := range this.fields {
			if f.Embedded() && f.Name() == name {
				if i != 0 {
					rest := make([]*types.Var, 0, len(this.fields))
					rest = append(rest, this.fields[:i]...)
					rest = append(rest, this.fields[i+1:]...)
					this.fields = append([]*types.Var{f}, rest...)
				}
				break
			}
		}
		return true
	}
	return false
}

// isPolymorphic reports whether the class declares or inherits any virtual
// method, i.e. whether it has (or shares) a vptr in its layout.
func isPolymorphic(cls clang.Cursor) bool {
	found := false
	clang.VisitChildren(cls, func(decl, parent clang.Cursor) clang.ChildVisitResult {
		switch decl.Kind {
		case lc.Cursor_CXXMethod, lc.Cursor_Destructor:
			if decl.CXXMethodIsVirtual() != 0 {
				found = true
				return clang.Break
			}
		case lc.Cursor_CXXBaseSpecifier:
			if b := decl.Type().Declaration().Definition(); b.IsNull() == 0 && isPolymorphic(b) {
				found = true
				return clang.Break
			}
		}
		return clang.Continue
	})
	return found
}

func baseClass(ctx *pkgCtx, this *classCtx, decl clang.Cursor, feats *int) (typ types.Type, name string, found bool) {
	t := decl.Type()
	// Depending on the libclang version, a base specifier's type may be reported
	// as an elaborated type (e.g. "struct Base") rather than the bare record;
	// unwrap it so the record lookup below works in both cases.
	if t.Kind == lc.Type_Elaborated {
		t = t.Named()
	}
	switch t.Kind {
	case lc.Type_Record, lc.Type_Typedef:
		typ, found = namedType(ctx, t, feats)
		if found {
			name = goNamedTypeName(typ)
			return
		}
	case lc.Type_Unexposed:
		typ, found = unexposedType(ctx, t, feats, this.scope())
		if found {
			name = goNamedTypeName(typ)
		}
		return
	}
	ctx.panicf(decl, "baseClass %s: unknown base class - %s (%d)", cNameOf(decl), clang.String(t), t.Kind)
	return
}

func goNamedTypeName(typ types.Type) string {
	switch typ := typ.(type) {
	case *types.Named:
		return typ.Obj().Name()
	case *types.Alias:
		return typ.Obj().Name()
	case *types.Basic:
		return typ.Name()
	case *types.TypeParam:
		return typ.Obj().Name()
	}
	panic("goNamedTypeName: unreachable")
}

// -----------------------------------------------------------------------------
