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
	"fmt"
	"go/token"
	"go/types"
	"unsafe"

	"github.com/goplus/llcppg/clang"
	lc "github.com/llarhub/clang-c"
)

// -----------------------------------------------------------------------------

// unionStorageName is the single unexported field of a generated union struct
// X. Its array type gives X the same size and alignment as the C union, while
// leaving X free of any exported field that could clash with a "XGof_ref_*"
// accessor. See issue goplus/llcppg#775.
const unionStorageName = "_xgo_union"

// unionRefPrefix is the reserved prefix of the accessor methods generated for
// each accessible union member: "XGof_ref_<member>" returns a typed pointer
// into the union's storage. The "XGof_"/"XGo_" prefixes are reserved, so a
// method derived from a C function must never start with them.
const unionRefPrefix = "XGof_ref_"

type unionCtx struct {
	scopeCtx
	decl clang.Cursor

	members     []unionMember
	parentTypes map[clang.Cursor]types.Type
}

type unionMember struct {
	decl clang.Cursor
	typ  types.Type
}

func (p *unionCtx) scope() *scopeCtx {
	return (*scopeCtx)(unsafe.Pointer(p))
}

// loadUnion translates a C/C++ union declaration into Go declarations.
//
// For a union U it emits a Go struct X with a single unexported storage field
// (see unionStorageName) sized and aligned like U, plus one accessor method
//
//	func (p *X) XGof_ref_foo() *T { return (*T)(unsafe.Pointer(p)) }
//
// for each accessible member foo of a convertible type T. Member names are kept
// verbatim (no PascalCase, no prefix trimming).
//
// Whether the union is global, in a namespace, or nested inside a class only
// affects naming: ns carries the enclosing prefix, so the union type name goes
// through getPubName like a struct's.
func loadUnion(ctx *pkgCtx, cName string, decl clang.Cursor) {
	if debugCompileDecl {
		ctx.logf(decl, "union %s", cName)
	}

	var typDecl, ok = ctx.typdecls[cName]
	if !ok {
		goName := ctx.typeName(cName, true)
		typDecl = newType(ctx, decl, cName, goName)
		ctx.typdecls[cName] = typDecl
	}

	if decl.IsCursorDefinition() == 0 {
		return // declaration only, no definition
	}

	feats := 0
	initUnionType(ctx, decl, typDecl, &feats)
	if feats&featAllIgnore != 0 {
		ctx.ignoref(feats, decl, "union %s: unsupported features, ignored", cName)
		ctx.ignoreType(cName, feats)
	}
}

func emitUnion(ctx *pkgCtx, decl clang.Cursor, goName string, feats *int) *types.Named {
	typDecl := newType(ctx, decl, "", goName)
	initUnionType(ctx, decl, typDecl, feats)
	return typDecl.Type()
}

func initUnionType(ctx *pkgCtx, decl clang.Cursor, typDecl typDecl, feats *int) {
	pkg := ctx.pkg
	typ := decl.Type()
	storage, ok := unionStorageType(typ)
	if !ok {
		// A union with no body (GNU empty union) has no storage: emit an empty
		// struct with no accessors.
		typDecl.InitType(pkg, types.NewStruct(nil, nil))
		return
	}
	this := &unionCtx{
		decl:        decl,
		parentTypes: make(map[clang.Cursor]types.Type),
	}
	clang.VisitChildren(decl, func(m, parent clang.Cursor) clang.ChildVisitResult {
		return loadUnionMember(ctx, pkg.Types, this, m, parent, feats)
	})
	if *feats&featAllIgnore != 0 {
		typDecl.Delete()
		return
	}

	// Members overlap: union { int i; double d; } needs one storage field.
	typDecl.InitType(pkg, unionStruct(ctx, decl, storage))
	recvPtr := types.NewPointer(typDecl.Type())
	ctx.addCompileUnit(func(ctx *pkgCtx) {
		compileUnionImpl(ctx, this, recvPtr)
	})
}

func loadUnionMember(ctx *pkgCtx, pkg *types.Package, this *unionCtx, decl, parent clang.Cursor, feats *int) clang.ChildVisitResult {
	switch decl.Kind {
	case lc.Cursor_FieldDecl:
		// union U { union { int b; }; };
		// b is promoted without a parentTyp.
		// union U { struct { int a; enum { E } e; }; };
		// a, e: parentTyp is the anonymous struct containing them.
		// union U { struct { int tag; struct { int a; }; union { int b; }; }; };
		// tag, a, b: parentTyp is the outer anonymous struct.
		if parentTyp, ok := this.parentTypes[parent]; ok {
			this.members = append(this.members, unionMember{decl: decl})
			this.parentTypes[decl] = parentTyp
			return clang.Continue
		}
		var fldType types.Type

		ft := decl.Type()
		if ftd := ft.Declaration(); ftd.IsAnonymous() != 0 && ftd.Kind != lc.Cursor_UnionDecl {
			switch {
			case ft.Kind == lc.Type_Record:
				var ok bool
				fldType, ok = emitClass(ctx, ftd, ctx.nextAnonName(), this.scope())
				if !ok {
					*feats |= featExplicitIgnore
					return clang.Continue
				}
			case ft.Kind == lc.Type_Enum:
				fldType = emitEnum(ctx, ftd, ctx.nextAnonName())
			default:
				ctx.panicf(ftd, "unknown anonymous field type (%d: %s)", ft.Kind, clang.String(ft))
			}
		} else {
			if fldType = toTypeEx(ctx, pkg, ft, flagIsVarDef, feats, this.scope()); *feats&featAllIgnore != 0 {
				return clang.Continue
			}
		}
		this.members = append(this.members, unionMember{decl: decl, typ: fldType})
	case lc.Cursor_VarDecl:
		if isPublic(decl) {
			loadVar(ctx, decl)
		}
	case lc.Cursor_UnionDecl:
		// a union-nested fieldless union, like,
		// for field union, handle in Cursor_FieldDecl branch
		// union{ union { } }
		if decl.IsAnonymousRecordDecl() != 0 {
			// union U { union { int x; }; };
			// Inner union: no parentTyp; recurse to collect x as U's member.
			// union U { struct { int tag; union { int x; }; }; };
			// Inner union: parentTyp is the struct containing tag; recurse using that view.
			if parentTyp, ok := this.parentTypes[parent]; ok {
				this.parentTypes[decl] = parentTyp
			}
			return clang.Recurse
		} else if decl.IsAnonymous() == 0 {
			// The parent struct view already loads V:
			// union U { struct { union V { int x; } v; }; };
			if this.parentTypes[parent] == nil {
				loadUnion(ctx, cNameOf(decl), decl)
			}
		}
	case lc.Cursor_ClassDecl, lc.Cursor_StructDecl:
		// a union-nested fieldless struct, like,
		// union{ struct{ int x; } }
		if decl.IsAnonymousRecordDecl() != 0 {
			// union U { struct { int x; }; };
			// This struct: no parentTyp; emit its own Go view, then recurse to collect x.
			// union U { struct { int tag; struct { int x; }; }; };
			// Inner struct: parentTyp is the struct containing tag; reuse it and recurse.
			parentTyp, ok := this.parentTypes[parent]
			if !ok {
				parentTyp, ok = emitClass(ctx, decl, ctx.nextAnonName(), this.scope())
				if !ok {
					*feats |= featExplicitIgnore
					return clang.Continue
				}
			}
			this.parentTypes[decl] = parentTyp
			return clang.Recurse
		} else if decl.IsAnonymous() == 0 {
			// The parent struct view already loads S:
			// union U { struct { struct S { int x; } s; }; };
			if this.parentTypes[parent] == nil {
				nested := newClassCtx(ctx, decl, this.scope())
				loadClass(ctx, cNameOf(decl), nested, decl)
			}
		}
	case lc.Cursor_EnumDecl:
		// a union-nested fieldless enum, like
		// union{ enum{ } }
		if isPublic(decl) {
			loadEnum(ctx, decl)
		}

	case lc.Cursor_TypedefDecl, lc.Cursor_TypeAliasDecl, lc.Cursor_TypeAliasTemplateDecl:
		// a union-nested typedef
		if isPublic(decl) {
			loadTypedef(ctx, decl, this.scope())
		}
	}
	return clang.Continue
}

// unionStruct builds the "type X struct { _xgo_union <storage> }" definition.
func unionStruct(ctx *pkgCtx, decl clang.Cursor, storage types.Type) *types.Struct {
	fld := types.NewField(goNodePos(ctx, decl), ctx.pkg.Types, unionStorageName, storage, false)
	return types.NewStruct([]*types.Var{fld}, nil)
}

// unionStorageType selects the array element type and length for a union's
// storage field from its size S and alignment A (both in bytes):
//
//   - A in {1,2,4,8}: element width is A, length is S/A. The element is
//     float32/float64 only when every scalar leaf of the union is a float of
//     width A (so by-value ABI classification stays correct); otherwise uintN.
//   - A > 8 (long double, __int128, aligned(16)): [S/8]uint64 plus a warning.
//     Exact alignment for these is out of scope (issue goplus/llcppg#775).
//
// It returns t == nil when the union has no storage (size 0 or a forward
// declaration with no known layout); the caller then emits an empty struct.
//
// aligned is false when A is not a power of two the generator maps to an element
// width (or A > 8); the caller then emits a byte array and skips accessors,
// since a reinterpret cast into an under-aligned storage would be unsound.
func unionStorageType(typ lc.Type) (t types.Type, aligned bool) {
	size := int64(typ.SizeOf())
	align := int64(typ.AlignOf())
	if size <= 0 || align <= 0 {
		return nil, false
	}
	switch align {
	case 1, 2, 4, 8:
		var elem types.Type
		if align == 4 && allFloatLeaves(typ, 4) {
			elem = types.Typ[types.Float32]
		} else if align == 8 && allFloatLeaves(typ, 8) {
			elem = types.Typ[types.Float64]
		} else {
			elem = uintOfWidth(align)
		}
		return types.NewArray(elem, size/align), true
	default:
		panic(fmt.Sprintf("[WARN] union alignment %d > 8 is not fully supported", align))
	}
}

func uintOfWidth(width int64) types.Type {
	switch width {
	case 1:
		return types.Typ[types.Uint8]
	case 2:
		return types.Typ[types.Uint16]
	case 4:
		return types.Typ[types.Uint32]
	default:
		return types.Typ[types.Uint64]
	}
}

// allFloatLeaves reports whether every scalar leaf of typ is a floating-point
// type of exactly width bytes. Arrays and nested structs/unions are flattened
// recursively; any non-float leaf, or a float of a different width, makes the
// whole union non-float so it falls back to an integer element. An empty
// aggregate has no float leaves and is therefore not all-float.
func allFloatLeaves(typ lc.Type, width int64) bool {
	found := false
	if !walkFloatLeaves(typ, width, &found) {
		return false
	}
	return found
}

func walkFloatLeaves(typ lc.Type, width int64, found *bool) bool {
	switch typ.Kind {
	case lc.Type_Float, lc.Type_Double, lc.Type_LongDouble, lc.Type_Float128, lc.Type_Float16:
		if int64(typ.SizeOf()) != width {
			return false
		}
		*found = true
		return true
	case lc.Type_ConstantArray:
		return walkFloatLeaves(typ.ArrayElement(), width, found)
	case lc.Type_Elaborated:
		return walkFloatLeaves(typ.Named(), width, found)
	case lc.Type_Record:
		decl := typ.Declaration()
		ok := true
		clang.VisitChildren(decl, func(m, parent clang.Cursor) clang.ChildVisitResult {
			switch m.Kind {
			case lc.Cursor_FieldDecl:
				if !walkFloatLeaves(m.Type(), width, found) {
					ok = false
					return clang.Break
				}
			case lc.Cursor_StructDecl, lc.Cursor_UnionDecl:
				if m.IsAnonymousRecordDecl() != 0 {
					return clang.Recurse
				}
			}
			return clang.Continue
		})
		return ok
	default:
		return false
	}
}

func compileUnionImpl(ctx *pkgCtx, this *unionCtx, recvPtr types.Type) {
	for _, member := range this.members {
		// union U { int a; union { int b; }; struct { int c; union { int d; }; }; };
		// a, b: no parentTyp; c, d: parentTyp is the anonymous struct containing c.
		if parentTyp, ok := this.parentTypes[member.decl]; ok {
			genUnionNestedStructAccessor(ctx, recvPtr, parentTyp, member.decl)
		} else {
			genUnionAccessor(ctx, recvPtr, member.typ, member.decl)
		}
	}
}

// An accessor's callback result needs its own calling convention.
// union U { int (*cb)(int); }; needs XGof_ref_cb() *_llcppg_anon_0.
func unionCallbackType(ctx *pkgCtx, typ types.Type, decl clang.Cursor) types.Type {
	sig, ok := typ.(*types.Signature)
	if !ok {
		return typ
	}
	featCallback := featHasCdeclCallback
	if decl.Type().Canonical().Pointee().FunctionTypeCallingConv() == lc.CallingConv_X86StdCall {
		featCallback = featHasStdcallCallback
	}
	typDecl := newType(ctx, decl, "", ctx.nextAnonName())
	typDecl.defs.SetComments(doDirectiveCallback(nil, featCallback))
	typDecl.InitType(ctx.pkg, sig)
	return typDecl.Type()
}

// genUnionAccessor emits, for a member foo of type T,
//
//	func (p *X) XGof_ref_foo() *T { return (*T)(unsafe.Pointer(p)) }
func genUnionAccessor(ctx *pkgCtx, recvPtr types.Type, fieldType types.Type, decl clang.Cursor) {
	pkg := ctx.pkg
	pkgTypes := pkg.Types
	member := clang.String(decl)

	name := unionRefPrefix + member
	fieldType = unionCallbackType(ctx, fieldType, decl)
	retType := types.NewPointer(fieldType)
	recv := types.NewParam(0, pkgTypes, "p", recvPtr)
	results := types.NewTuple(types.NewParam(0, pkgTypes, "", retType))
	sig := types.NewSignatureType(recv, nil, nil, nil, results, false)

	f, err := pkg.NewFuncWith(goNodePos(ctx, decl), name, sig, nil)
	if err != nil {
		ctx.panicf(decl, "union field %s: genUnionAccessor failed - %v", member, err)
	}
	cb := f.BodyStart(pkg)
	// return (*T)(unsafe.Pointer(p))
	cb.Typ(retType).
		Typ(ctx.unsafePointer()).VarVal("p").
		Call(1).
		Call(1).
		Return(1).End()
}

// genUnionNestedStructAccessor emits, for a member foo of type T in an
// anonymous struct view S, as in union { struct { int foo; }; },
//
//	func (p *X) XGof_ref_foo() *T { return &(*S)(unsafe.Pointer(p)).Foo }
//
// For union { struct { int tag; union { int foo; }; }; }, foo is promoted
// through the embedded union's accessor, so it emits
//
//	func (p *X) XGof_ref_foo() *T { return (*S)(unsafe.Pointer(p)).XGof_ref_foo() }
func genUnionNestedStructAccessor(ctx *pkgCtx, recvPtr types.Type, parentTyp types.Type, m clang.Cursor) {
	pkg := ctx.pkg
	pkgTypes := pkg.Types
	member := clang.String(m)

	name := unionRefPrefix + member
	selector := ctx.fieldName(member, isPublic(m))
	obj, _, _ := types.LookupFieldOrMethod(types.NewPointer(parentTyp), true, pkgTypes, selector)
	field := obj != nil
	var retType types.Type
	if field {
		retType = types.NewPointer(obj.Type())
	} else {
		selector = name
		obj, _, _ = types.LookupFieldOrMethod(types.NewPointer(parentTyp), true, pkgTypes, selector)
		retType = obj.Type().(*types.Signature).Results().At(0).Type()
	}
	elemType := retType.(*types.Pointer).Elem()
	callbackType := unionCallbackType(ctx, elemType, m)
	cast := callbackType != elemType
	if cast {
		retType = types.NewPointer(callbackType)
	}
	recv := types.NewParam(0, pkgTypes, "p", recvPtr)
	results := types.NewTuple(types.NewParam(0, pkgTypes, "", retType))
	sig := types.NewSignatureType(recv, nil, nil, nil, results, false)

	f, err := pkg.NewFuncWith(goNodePos(ctx, m), name, sig, nil)
	if err != nil {
		ctx.panicf(m, "union field %s: genUnionNestedStructAccessor failed - %v", member, err)
	}
	cb := f.BodyStart(pkg)
	if cast {
		cb.Typ(retType)
	}
	cb.Typ(types.NewPointer(parentTyp)).Typ(ctx.unsafePointer()).VarVal("p").Call(1).Call(1)
	if field {
		cb.MemberRef(selector).UnaryOp(token.AND)
	} else {
		cb.MemberVal(selector, 1).Call(0)
	}
	if cast {
		cb.Call(1)
	}
	cb.Return(1).End()
}

// -----------------------------------------------------------------------------
