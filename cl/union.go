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
	"log"
	"strconv"

	"github.com/goplus/llcppg/clang"
	lc "github.com/goplus/llcppg/lib/clang"
)

// -----------------------------------------------------------------------------
//
// C/C++ union bindings.
//
// Go has no union type, so a C union U is bound as a Go struct X that owns a
// single storage field named unionField ("_xgo_union"), whose array type is
// chosen to reproduce U's size and alignment. For every accessible member foo
// of U, llcppg emits an accessor
//
//	func (p *X) XGo_union_foo() *T { return (*T)(unsafe.Pointer(p)) }
//
// that reinterprets the storage as the member's Go type T. One method serves
// both reads and writes; reading a member other than the one last written is
// type punning, exactly as in C. See issue goplus/llcppg#764 for the full
// proposal.

// unionField is the single unexported storage field of a generated union
// struct. It coexists with the exported "XGo_union_*" accessor methods (Go
// forbids a field and a method sharing a name; the differing shape avoids any
// clash).
const unionField = "_xgo_union"

// unionAccessorPrefix namespaces the member accessors. The C member name is
// kept verbatim after this prefix (no PascalCase, no keyword escaping, no
// prefix trimming): a member "type" becomes "XGo_union_type". Methods llcppg
// derives from C functions never start with "XGo_", so accessors cannot collide
// with them, and C forbids duplicate member names within a union.
const unionAccessorPrefix = "XGo_union_"

// unionMember is one accessible direct member of a union for which an accessor
// is generated.
type unionMember struct {
	name string     // verbatim C member name
	typ  types.Type // Go type T of the member (accessor returns *T)
	off  int64      // byte offset of the member within the union storage
	decl clang.Cursor
}

// loadUnion translates a C/C++ union declaration into a Go struct type plus its
// member accessors.
//
// name, when non-empty, is the fully-qualified C name the union should take
// (used when a "typedef union { ... } Name;" gives a tagless union its name);
// otherwise the union is named by its tag, following the same rules as structs.
//
// A hoisted (tagless, inline) union passed via hoistUnionField supplies name ==
// "" and a hoisted Go name is used instead; that path goes through
// loadUnionAs. Callers that name the union by tag or typedef use loadUnion,
// which derives the Go name from the C name with getPubName.
func loadUnion(ctx *pkgCtx, decl clang.Cursor, ns, name string) *types.Named {
	origName := name
	if origName == "" {
		origName = ns + clang.String(decl)
	}
	goName, rewritten := ctx.getPubName(origName, -1)
	return loadUnionAs(ctx, decl, goName, origName, rewritten, name != "")
}

// loadUnionAs is the core of loadUnion, taking the already-resolved Go type name
// goName. origName is the C name used for registration/substObj; rewritten marks
// that goName differs from origName (so the C name is also inserted into scope);
// registerName registers the union under origName in ctx.types so a field of
// that C/typedef name resolves through toType.
func loadUnionAs(ctx *pkgCtx, decl clang.Cursor, goName, origName string, rewritten, registerName bool) *types.Named {
	pkg := ctx.pkg
	pkgTypes := pkg.Types

	if debugCompileDecl {
		log.Println("union", origName, "=>", goName)
	}

	typDecl := pkg.NewTypeDefs().NewType(goName, goNode(ctx, decl))
	typNamed := typDecl.Type()
	if rewritten {
		substObj(pkgTypes, pkgTypes.Scope(), origName, typNamed.Obj())
	}
	// Register under the record type spelling (as loadClass does) so a union
	// used as a struct field resolves through toType, and under the typedef name
	// spelling when one is given so "Value" fields resolve too.
	ctx.types[clang.String(decl.Type())] = typNamed.Obj()
	if registerName {
		ctx.types[origName] = typNamed.Obj()
	}

	// An opaque union (forward-declared, empty, or zero-sized) has no layout to
	// mirror and no members to expose: bind it as an empty struct, like an
	// opaque struct, so "union U *" still maps to "*U" in signatures.
	utyp := decl.Type()
	size := utyp.SizeOf()
	if decl.Definition().IsNull() != 0 || size <= 0 {
		typDecl.InitType(pkg, types.NewStruct(nil, nil))
		return typNamed
	}

	ctx.forceImportUnsafe()

	storage := unionStorage(ctx, utyp, goName)
	field := types.NewField(goNodePos(ctx, decl), pkgTypes, unionField, storage, false)
	typDecl.InitType(pkg, types.NewStruct([]*types.Var{field}, nil))

	members := unionMembers(ctx, pkgTypes, decl, 0)
	recvPtr := types.NewPointer(typNamed)
	ctx.compiles = append(ctx.compiles, func(ctx *pkgCtx) {
		for _, m := range members {
			genUnionAccessor(ctx, recvPtr, m)
		}
	})
	return typNamed
}

// unionStorage returns the Go type of the "_xgo_union" storage field, chosen so
// that the struct's size and alignment equal the union's (proposal §1).
func unionStorage(ctx *pkgCtx, utyp lc.Type, uName string) types.Type {
	align := utyp.AlignOf()
	size := utyp.SizeOf()

	// Step 3: an alignment larger than 8 has no matching Go scalar. Preserve the
	// size with [S/8]uint64 (alignment lowered to 8) and warn; such a union must
	// not be relied on where its full alignment matters.
	if align > 8 {
		log.Printf("warning: union %s has alignment %d > 8; emitting [%d]uint64 which lowers alignment to 8\n", uName, align, size/8)
		return types.NewArray(types.Typ[types.Uint64], size/8)
	}

	// Steps 1 & 2: element width is the alignment; the element is a float of
	// that width when every scalar leaf is a float of that width, else an
	// unsigned integer of that width. The length is size/align.
	elem := unionElem(align, unionAllFloat(utyp.TypeDeclaration(), align))
	return types.NewArray(elem, size/align)
}

// unionElem returns the storage element type for the given element width in
// bytes. When allFloat is true a floating-point element is used (only valid for
// widths 4 and 8), otherwise an unsigned integer of that width.
func unionElem(width int64, allFloat bool) types.Type {
	if allFloat {
		switch width {
		case 4:
			return types.Typ[types.Float32]
		case 8:
			return types.Typ[types.Float64]
		}
	}
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

// unionAllFloat reports whether every scalar leaf of the union is a
// floating-point type of exactly the given width (float with width 4, double
// with width 8). Arrays contribute their element type; nested structs and
// unions contribute their own leaves. Only then is a floating-point storage
// element used, so that by-value ABI classification of an all-float union is
// preserved (proposal §1 step 2).
func unionAllFloat(decl clang.Cursor, width int64) bool {
	allFloat := true
	any := false
	var walk func(t lc.Type)
	walk = func(t lc.Type) {
		t = t.CanonicalType()
		switch t.Kind {
		case lc.TypeConstantArray, lc.TypeIncompleteArray, lc.TypeVariableArray:
			walk(t.ArrayElementType())
			return
		case lc.TypeFloat:
			any = true
			if width != 4 {
				allFloat = false
			}
			return
		case lc.TypeDouble:
			any = true
			if width != 8 {
				allFloat = false
			}
			return
		case lc.TypeRecord:
			rec := t.TypeDeclaration()
			clang.VisitChildren(rec, func(cur, parent clang.Cursor) clang.ChildVisitResult {
				if cur.Kind == lc.CursorFieldDecl {
					walk(cur.Type())
				}
				return clang.Continue
			})
			return
		}
		allFloat = false
	}
	clang.VisitChildren(decl, func(cur, parent clang.Cursor) clang.ChildVisitResult {
		if cur.Kind == lc.CursorFieldDecl {
			walk(cur.Type())
		}
		return clang.Continue
	})
	return any && allFloat
}

// unionMembers collects the accessible direct members of a union at the given
// byte base offset. A C11 anonymous struct/union member contributes its own
// members promoted to this union (at their field offsets, added to base); a
// bit-field or an unconvertible member is reported and skipped (proposal §6).
func unionMembers(ctx *pkgCtx, pkg *types.Package, decl clang.Cursor, base int64) (members []unionMember) {
	clang.VisitChildren(decl, func(cur, parent clang.Cursor) clang.ChildVisitResult {
		if cur.Kind != lc.CursorFieldDecl {
			return clang.Continue
		}
		name := clang.String(cur)
		if cur.FieldDeclBitWidth() >= 0 {
			log.Printf("union member %q: bit-field, no accessor generated\n", name)
			return clang.Continue
		}
		// An anonymous struct/union member exposes its own members directly
		// (proposal §5.3, direct case): recurse at this member's offset.
		if name == "" {
			off := base + cur.OffsetOfField()/8
			members = append(members, unionMembers(ctx, pkg, cur.Type().TypeDeclaration(), off)...)
			return clang.Continue
		}
		typ := toUnionMemberType(ctx, pkg, cur.Type())
		if typ == nil {
			log.Printf("union member %q: cannot convert C type %q, no accessor generated\n", name, clang.String(cur.Type()))
			return clang.Continue
		}
		off := base + cur.OffsetOfField()/8
		members = append(members, unionMember{name: name, typ: typ, off: off, decl: cur})
		return clang.Continue
	})
	return
}

// toUnionMemberType converts a union member's C type to its Go type T (the same
// mapping used for a struct field of that type), returning nil when the type is
// not convertible so the caller can skip it with a comment.
func toUnionMemberType(ctx *pkgCtx, pkg *types.Package, t lc.Type) (typ types.Type) {
	defer func() {
		if recover() != nil {
			typ = nil
		}
	}()
	return toType(ctx, pkg, t, flagIsStructField)
}

// genUnionAccessor emits the member accessor:
//
//	func (p *X) XGo_union_foo() *T { return (*T)(unsafe.Pointer(p)) }
//
// For a member at offset 0 (every direct member of a union) the body is the
// plain reinterpret cast; a promoted member at a non-zero offset uses
// unsafe.Add(unsafe.Pointer(p), off). The accessor's doc comment carries the
// original C member declaration.
func genUnionAccessor(ctx *pkgCtx, recvPtr types.Type, m unionMember) {
	pkg := ctx.pkg
	pkgTypes := pkg.Types
	tyUP := types.Typ[types.UnsafePointer]
	ptrT := types.NewPointer(m.typ)

	recv := types.NewParam(token.NoPos, pkgTypes, "p", recvPtr)
	results := types.NewTuple(types.NewParam(token.NoPos, pkgTypes, "", ptrT))
	sig := types.NewSignatureType(recv, nil, nil, nil, results, false)

	f, err := pkg.NewFuncWith(goNodePos(ctx, m.decl), unionAccessorPrefix+m.name, sig, nil)
	if err != nil {
		panic("genUnionAccessor: " + err.Error())
	}
	f.SetComments(pkg, &ast.CommentGroup{
		List: []*ast.Comment{{Text: "\n// " + clang.String(m.decl.Type()) + " " + m.name}},
	})
	cb := f.BodyStart(pkg)
	cb.Typ(ptrT) // (*T)(
	if m.off == 0 {
		cb.Typ(tyUP).VarVal("p").Call(1) // unsafe.Pointer(p)
	} else {
		// unsafe.Add(unsafe.Pointer(p), off)
		cb.Val(ctx.unsafeAdd()).
			Typ(tyUP).VarVal("p").Call(1).
			Val(int(m.off)).
			Call(2)
	}
	cb.Call(1).Return(1).End() // )
}

// -----------------------------------------------------------------------------

// hoistUnionName returns the next "_llcppg_union_<n>" name for a hoisted
// (tagless, inline) union and advances the package-wide counter (proposal §5).
func hoistUnionName(ctx *pkgCtx) string {
	name := "_llcppg_union_" + strconv.Itoa(ctx.unionSeq)
	ctx.unionSeq++
	return name
}

// hoistUnionField handles a struct field whose type is a tagless union declared
// inline (proposal §5). Such a union has no name to attach accessors to, so it
// is hoisted into a named Go type "_llcppg_union_<n>". Two shapes are handled:
//
//   - §5.1 A field with a name keeps its converted name and takes the hoisted
//     type.
//   - §5.2 A C11 anonymous union (no field name) is embedded, so Go promotes the
//     hoisted type's accessors onto the parent struct.
//
// It returns ok == false when decl is not a tagless-inline-union field (a
// tagged or typedef'd union field resolves through toType as usual).
func hoistUnionField(ctx *pkgCtx, pkg *types.Package, cls *classCtx, decl clang.Cursor, origName string) (fld *types.Var, ok bool) {
	rec := decl.Type().TypeDeclaration()
	if rec.Kind != lc.CursorUnionDecl {
		return nil, false
	}
	// Only a tagless inline union is hoisted; a tagged union (union Tag { ... })
	// is a named type resolved through toType.
	if clang.String(rec) != "" && rec.IsAnonymous() == 0 {
		return nil, false
	}

	named := loadUnionAs(ctx, rec, hoistUnionName(ctx), "", false, false)

	// §5.2: a C11 anonymous union member has no field name; embed the hoisted
	// type so its accessors are promoted onto the parent.
	if origName == "" || decl.IsAnonymousRecordDecl() != 0 {
		return types.NewField(goNodePos(ctx, decl), pkg, named.Obj().Name(), named, true), true
	}

	// §5.1: a named field keeps its converted name and takes the hoisted type.
	fldName := origName
	if cls.inPublic {
		fldName, _ = ctx.getPubName(origName, -1)
	}
	return types.NewField(goNodePos(ctx, decl), pkg, fldName, named, false), true
}

// -----------------------------------------------------------------------------
