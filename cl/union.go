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
	"log"

	"github.com/goplus/llcppg/clang"
	lc "github.com/goplus/llcppg/lib/clang"
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
func loadUnion(ctx *pkgCtx, decl clang.Cursor, ns string) {
	if decl.IsCursorDefinition() == 0 {
		return
	}

	origName := ns + clang.String(decl)
	if debugCompileDecl {
		log.Println("union", origName)
	}

	pkg := ctx.pkg
	uName, _ := ctx.getPubName(origName, -1)
	typDecl := pkg.NewTypeDefs().NewType(uName, goNode(ctx, decl))
	typNamed := typDecl.Type()

	// Register the union type under its C spelling so members that reference it
	// (e.g. "union U *next") and typedef aliases resolve to the same X.
	ctx.types[clang.String(decl.Type())] = typNamed.Obj()

	typ := decl.Type()
	storage, ok := unionStorageType(typ)
	if !ok {
		// A union with no body (GNU empty union) has no storage: emit an empty
		// struct with no accessors.
		typDecl.InitType(pkg, types.NewStruct(nil, nil))
		return
	}
	ctx.forceImportUnsafe()
	typDecl.InitType(pkg, unionStruct(ctx, decl, storage))

	// Collect the members that get an accessor, then generate the accessors in
	// the compile phase (after every type is registered) so member types
	// referencing other records resolve regardless of declaration order.
	// Bit-fields are skipped; an under-aligned union (!aligned) generates no
	// accessors at all.
	var members []clang.Cursor
	clang.VisitChildren(decl, func(m, parent clang.Cursor) clang.ChildVisitResult {
		switch m.Kind {
		case lc.CursorFieldDecl:
			members = append(members, m)
		}
		return clang.Continue
	})

	recvPtr := types.NewPointer(typNamed)
	ctx.compiles = append(ctx.compiles, func(ctx *pkgCtx) {
		for _, m := range members {
			genUnionAccessor(ctx, recvPtr, m)
		}
	})
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
	case lc.TypeFloat, lc.TypeDouble, lc.TypeLongDouble, lc.TypeFloat128, lc.TypeFloat16:
		if int64(typ.SizeOf()) != width {
			return false
		}
		*found = true
		return true
	case lc.TypeConstantArray:
		return walkFloatLeaves(typ.ArrayElementType(), width, found)
	case lc.TypeElaborated:
		return walkFloatLeaves(typ.NamedType(), width, found)
	case lc.TypeRecord:
		decl := typ.TypeDeclaration()
		ok := true
		clang.VisitChildren(decl, func(m, parent clang.Cursor) clang.ChildVisitResult {
			if m.Kind != lc.CursorFieldDecl {
				return clang.Continue
			}
			if !walkFloatLeaves(m.Type(), width, found) {
				ok = false
				return clang.Break
			}
			return clang.Continue
		})
		return ok
	default:
		return false
	}
}

// genUnionAccessor emits, for a member foo of type T,
//
//	func (p *X) XGof_ref_foo() *T { return (*T)(unsafe.Pointer(p)) }
func genUnionAccessor(ctx *pkgCtx, recvPtr types.Type, m clang.Cursor) {
	pkg := ctx.pkg
	pkgTypes := pkg.Types
	member := clang.String(m)

	fldType := toType(ctx, pkgTypes, m.Type(), flagIsStructField)

	name := unionRefPrefix + member
	retType := types.NewPointer(fldType)
	recv := types.NewParam(token.NoPos, pkgTypes, "p", recvPtr)
	results := types.NewTuple(types.NewParam(token.NoPos, pkgTypes, "", retType))
	sig := types.NewSignatureType(recv, nil, nil, nil, results, false)

	f, err := pkg.NewFuncWith(goNodePos(ctx, m), name, sig, nil)
	if err != nil {
		log.Panicln("genUnionAccessor:", member, err)
	}
	cb := f.BodyStart(pkg)
	// return (*T)(unsafe.Pointer(p))
	cb.Typ(retType).
		Typ(types.Typ[types.UnsafePointer]).VarVal("p").
		Call(1).
		Call(1).
		Return(1).End()
}

// -----------------------------------------------------------------------------
