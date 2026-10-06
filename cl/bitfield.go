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

	"github.com/goplus/llcppg/clang"
	lc "github.com/llarhub/clang-c"
)

// Support for C/C++ bit-fields (layout). See issue goplus/llcppg#770.
//
// Go has no bit-fields, so a C struct S that declares bit-fields is converted
// to a Go struct X that keeps the size, alignment and member offsets of S:
// every run of consecutive bit-fields is stored as a byte array
// "_xgo_bits_<N>" at the byte position where C puts those bits, regular members
// keep their natural conversion, and explicit padding and an "_xgo_align"
// marker are inserted where Go's own layout would otherwise diverge from C.
//
// This file implements part 1 of the proposal (the Go representation). The
// per-field getter/setter accessors (part 2) build on the offsets and widths
// collected here; see planBitFieldLayout's accessor list.
//
// Bits are numbered from the least significant bit of each byte, matching the
// little-endian x86-64 / AArch64 ABIs (assumption 2 of the proposal). The byte
// offsets and widths come from libclang (OffsetOfField, FieldDeclBitWidth), so
// the result is correct for whatever ABI and packing the headers use.

// -----------------------------------------------------------------------------

// Reserved names used by the generated bit-field storage. The "_xgo_" prefix
// keeps them unexported so they never collide with a C-derived exported member.
const (
	bitsFieldPrefix = "_xgo_bits_"  // "_xgo_bits_<N>" run storage field
	bitAlignName    = "_xgo_align"  // "_xgo_align [0]uintW" alignment marker
	bitOpaqueName   = "_xgo_opaque" // "_xgo_opaque [S]uint8" opaque storage
)

// bitAccessor describes one named bit-field that is eligible for a getter/setter
// pair. off and width are the bit offset (from the start of X) and the width in
// bits, as reported by libclang; signed and isBool come from the declared type.
type bitAccessor struct {
	decl   clang.Cursor
	name   string     // C name, kept verbatim (the virtual field name)
	goType types.Type // T, the Go type of the declared C type
	off    int64      // bit offset from the start of X
	width  int64      // width in bits
	signed bool       // declared type is a signed integer / enum
	isBool bool       // declared type is Go bool (_Bool / C++ bool)
}

// bitLayout is the result of planning the storage of a record that contains at
// least one bit-field: the Go fields that replace the record's members (regular
// fields, run storage arrays and explicit padding) plus the named bit-fields
// that are eligible for accessors.
//
// accessors carries the exact offset/width/signedness of each named bit-field.
// The getter/setter accessors (emitted in bitfield_accessor.go) build directly
// on this list.
type bitLayout struct {
	fields    []*types.Var
	accessors []bitAccessor
	// opaque is set when the layout cannot be represented field-by-field (a
	// packed struct, or a run overlapping another member). The caller then emits
	// an opaque byte array of the right size/alignment and no accessors.
	opaque bool
	reason string // human-readable reason, set when opaque
	size   int64  // size of S, in bytes
	align  int64  // alignment of S, in bytes
}

// hasBitField reports whether the record declaration decl has at least one
// bit-field member. Only such records go through the bit-field layout path;
// every other struct keeps its ordinary field-by-field conversion.
func hasBitField(decl clang.Cursor) bool {
	found := false
	clang.VisitChildren(decl, func(m, parent clang.Cursor) clang.ChildVisitResult {
		if m.Kind == lc.Cursor_FieldDecl && m.IsBitField() != 0 {
			found = true
			return clang.Break
		}
		return clang.Continue
	})
	return found
}

// -----------------------------------------------------------------------------

// member is one child of the record in declaration order, classified as a
// regular field or a bit-field. For a bit-field, off is its bit offset and
// width its declared width; for a regular field, off is its byte offset.
type member struct {
	decl   clang.Cursor
	name   string
	goType types.Type
	isBit  bool
	off    int64
	width  int64
	signed bool
	isBool bool
}

// planBitFieldLayout converts the members of a bit-field-bearing record into Go
// fields that preserve the C layout, following section 2 of the proposal.
//
// toGoType converts a member's declared C type to its Go type T, returning
// ok=false when the type is unsupported/ignored (the caller supplies it so this
// file stays independent of the type-conversion flags). A member with an
// unsupported type cannot be laid out field-by-field, so the record falls back
// to opaque storage.
func planBitFieldLayout(ctx *pkgCtx, pos token.Pos, decl clang.Cursor, toGoType func(lc.Type) (types.Type, bool)) bitLayout {
	pkgTypes := ctx.pkg.Types
	typ := decl.Type()
	out := bitLayout{
		align: int64(typ.AlignOf()),
		size:  int64(typ.SizeOf()),
	}
	if out.size <= 0 || out.align <= 0 {
		return out // no known layout (forward declaration): nothing to place
	}

	// Bits are allocated from the least significant bit of a byte only on a
	// little-endian target. On big-endian the layout is still preserved but the
	// bit numbering differs, so no accessor is generated (open question 3).
	accessible := !isBigEndianTarget(typ)

	var members []member
	packed := false
	unsupported := false
	clang.VisitChildren(decl, func(m, parent clang.Cursor) clang.ChildVisitResult {
		if m.Kind != lc.Cursor_FieldDecl {
			return clang.Continue
		}
		mt := m.Type()
		goType, ok := toGoType(mt)
		if !ok {
			unsupported = true
			return clang.Break
		}
		if m.IsBitField() != 0 {
			off := int64(m.OffsetOfField())
			if off < 0 {
				// A negative offset is a libclang CXTypeLayoutError code
				// (incomplete/dependent member); the layout is not known.
				packed = true
			}
			members = append(members, member{
				decl:   m,
				name:   clang.String(m),
				goType: goType,
				isBit:  true,
				off:    off,
				width:  int64(m.FieldDeclBitWidth()),
				signed: isSignedType(mt),
				isBool: mt.Kind == lc.Type_Bool,
			})
		} else {
			off := int64(m.OffsetOfField())
			if off < 0 || off%8 != 0 {
				packed = true
			}
			members = append(members, member{
				decl:   m,
				name:   ctx.fieldName(clang.String(m), isPublic(m)),
				goType: goType,
				off:    off / 8,
			})
		}
		return clang.Continue
	})
	if unsupported {
		out.opaque, out.reason = true, "a member has an unsupported type"
		return out
	}
	if packed {
		out.opaque, out.reason = true, "a member is not byte-aligned or has no known offset (packed struct or incomplete member)"
		return out
	}

	sizes := goSizes()
	var fields []*types.Var
	var accessors []bitAccessor
	curOff := int64(0) // next free byte offset in the Go struct being built
	runIndex := 0

	// emitPad ensures the next item lands at byte offset target. Go will itself
	// round curOff up to the item's natural alignment, so explicit "_ [k]uint8"
	// padding is inserted only for the gap C leaves beyond that aligned offset.
	// It returns false when even the aligned offset is already past target, which
	// means the layout is not representable field-by-field (a packed struct).
	emitPad := func(target, align int64) bool {
		aligned := roundUp(curOff, align)
		if aligned > target {
			return false
		}
		if target > aligned {
			fields = append(fields, padField(pkgTypes, pos, target-aligned))
		}
		curOff = target
		return true
	}

	for i := 0; i < len(members); {
		m := members[i]
		if !m.isBit {
			if !emitPad(m.off, sizes.Alignof(m.goType)) {
				out.opaque, out.reason = true, fmt.Sprintf("member %q does not fit at its C offset", m.name)
				return out
			}
			fields = append(fields, types.NewField(pos, pkgTypes, m.name, m.goType, false))
			curOff += sizes.Sizeof(m.goType)
			i++
			continue
		}

		// Collect a maximal run of consecutive bit-fields and its bit span.
		j := i
		var lo, hi int64 = -1, -1
		for j < len(members) && members[j].isBit {
			b := members[j]
			if b.width > 0 {
				if lo < 0 || b.off < lo {
					lo = b.off
				}
				if b.off+b.width > hi {
					hi = b.off + b.width
				}
			}
			j++
		}
		if lo < 0 {
			i = j // run of only unnamed/zero-width bit-fields: no storage
			continue
		}
		byteLo := lo / 8
		n := (hi+7)/8 - byteLo
		if n <= 0 {
			// Defensive: a well-formed run always spans >= 1 byte. A non-positive
			// span means libclang reported an inconsistent offset/width, so fall
			// back to opaque storage rather than building an invalid array type.
			out.opaque, out.reason = true, "a bit-field run has an inconsistent span"
			return out
		}
		if !emitPad(byteLo, 1) { // a run storage array has alignment 1
			out.opaque, out.reason = true, "a bit-field run overlaps another member"
			return out
		}
		name := fmt.Sprintf("%s%d", bitsFieldPrefix, runIndex)
		fields = append(fields, types.NewField(pos, pkgTypes, name, types.NewArray(types.Typ[types.Uint8], n), false))
		curOff += n
		runIndex++

		if accessible {
			for k := i; k < j; k++ {
				b := members[k]
				if b.width == 0 || b.name == "" {
					continue // zero-width or unnamed: padding, no accessor
				}
				if b.width > 64 {
					// The generated accessors move the value through uint64/int64,
					// so a bit-field wider than 64 bits (e.g. __int128 : 100) would
					// silently truncate. Keep its bits in the run storage but emit
					// no accessor.
					ctx.logf(b.decl, "[WARN] bit-field %q width %d > 64: no accessor", b.name, b.width)
					continue
				}
				accessors = append(accessors, bitAccessor{
					decl:   b.decl,
					name:   b.name,
					goType: b.goType,
					off:    b.off,
					width:  b.width,
					signed: b.signed,
					isBool: b.isBool,
				})
			}
		}
		i = j
	}

	// Raise the alignment of X to A when the byte arrays alone cannot: a
	// zero-length array "_xgo_align [0]uintW" placed first forces alignment W.
	if goStructAlign(fields, sizes) < out.align {
		elem, warn := alignElem(out.align)
		if warn {
			ctx.logf(decl, "[WARN] bit-field struct alignment %d > 8 is not fully representable", out.align)
		}
		marker := types.NewField(pos, pkgTypes, bitAlignName, types.NewArray(elem, 0), false)
		fields = append([]*types.Var{marker}, fields...)
	}

	// Append trailing padding so the Go size matches sizeof(S). If the Go struct
	// is already larger than S (its own alignment/tail-padding rounding diverged
	// from C, e.g. for an over-aligned record), the field-by-field form cannot
	// preserve the size, so fall back to opaque storage.
	if used := goStructSize(fields, sizes); used < out.size {
		fields = append(fields, padField(pkgTypes, pos, out.size-used))
	} else if used > out.size {
		out.opaque, out.reason = true, fmt.Sprintf("Go layout size %d exceeds sizeof %d", used, out.size)
		return out
	}

	out.fields = fields
	out.accessors = accessors
	return out
}

// padField builds an explicit padding field "_ [n]uint8".
func padField(pkg *types.Package, pos token.Pos, n int64) *types.Var {
	return types.NewField(pos, pkg, "_", types.NewArray(types.Typ[types.Uint8], n), false)
}

// alignElem returns the element type of the "_xgo_align [0]uintW" marker for a
// struct of alignment align bytes. For align > 8 no Go scalar is wide enough, so
// it falls back to uint64 and reports warn=true.
func alignElem(align int64) (elem types.Type, warn bool) {
	switch align {
	case 1:
		return types.Typ[types.Uint8], false
	case 2:
		return types.Typ[types.Uint16], false
	case 4:
		return types.Typ[types.Uint32], false
	case 8:
		return types.Typ[types.Uint64], false
	default:
		return types.Typ[types.Uint64], true
	}
}

// -----------------------------------------------------------------------------

// roundUp rounds n up to the next multiple of align (align is a power of two
// >= 1).
func roundUp(n, align int64) int64 {
	if align <= 1 {
		return n
	}
	return (n + align - 1) / align * align
}

// goSizes returns the Go type sizes used to lay out generated structs. The
// bit-field design targets 64-bit little-endian ABIs (proposal assumption 2),
// so an 8-byte word / max-align model is used regardless of the host GOARCH.
func goSizes() types.Sizes {
	return &types.StdSizes{WordSize: 8, MaxAlign: 8}
}

// goStructAlign is the alignment Go computes for a struct made of the fields.
func goStructAlign(fields []*types.Var, sizes types.Sizes) int64 {
	a := int64(1)
	for _, f := range fields {
		if fa := sizes.Alignof(f.Type()); fa > a {
			a = fa
		}
	}
	return a
}

// goStructSize is the size Go computes for a struct made of the fields.
func goStructSize(fields []*types.Var, sizes types.Sizes) int64 {
	return sizes.Sizeof(types.NewStruct(fields, nil))
}

// isSignedType reports whether the declared C type of a bit-field is a signed
// integer. For an enum the signedness is that of its underlying type.
func isSignedType(t lc.Type) bool {
	switch t.Kind {
	case lc.Type_Char_S, lc.Type_SChar, lc.Type_Short, lc.Type_Int,
		lc.Type_Long, lc.Type_LongLong, lc.Type_Int128:
		return true
	case lc.Type_Enum:
		return isSignedType(t.Declaration().EnumDeclIntegerType())
	case lc.Type_Elaborated:
		return isSignedType(t.Named())
	case lc.Type_Typedef:
		return isSignedType(t.Canonical())
	default:
		return false
	}
}

// isBigEndianTarget reports whether the record's target is big-endian. libclang
// does not expose endianness directly; the layout is preserved either way, and
// accessors are only skipped on big-endian (open question 3). The common
// targets (x86-64, AArch64) are little-endian, so this returns false; a future
// change can thread the real target triple through here.
func isBigEndianTarget(typ lc.Type) bool {
	_ = typ
	return false
}

// bitFieldDoc returns the original C declaration of a bit-field for its doc
// comment, e.g. "unsigned int mode : 3".
func bitFieldDoc(decl clang.Cursor) string {
	return fmt.Sprintf("%s %s : %d", clang.String(decl.Type()), clang.String(decl), int64(decl.FieldDeclBitWidth()))
}

// -----------------------------------------------------------------------------

// hasBaseOrNestedField reports whether the record has a C++ base specifier or an
// anonymous record/field that the bit-field layout path does not handle. Such
// records fall back to the ordinary field-by-field conversion (which does not
// yet support bit-fields in them). A plain C struct/union has only FieldDecl
// members and nested type/enum/method declarations, none of which trip this.
func hasBaseOrNestedField(cls clang.Cursor) bool {
	found := false
	clang.VisitChildren(cls, func(decl, parent clang.Cursor) clang.ChildVisitResult {
		switch decl.Kind {
		case lc.Cursor_CXXBaseSpecifier:
			found = true
			return clang.Break
		case lc.Cursor_FieldDecl:
			// An anonymous-record field (an inline struct/union used as a field's
			// type) needs the hoisting the ordinary path does; the bit-field path
			// handles only named, convertible field types.
			if decl.Type().Declaration().IsAnonymous() != 0 {
				found = true
				return clang.Break
			}
		case lc.Cursor_StructDecl, lc.Cursor_ClassDecl, lc.Cursor_UnionDecl:
			if decl.IsAnonymousRecordDecl() != 0 {
				found = true
				return clang.Break
			}
		}
		return clang.Continue
	})
	return found
}

// initBitFieldType builds the Go struct for a bit-field-bearing record from its
// planned layout, installs it as the type's underlying struct, and schedules the
// getter/setter accessors for the compile phase. It also loads nested type,
// enum and method declarations through the ordinary member path so they are not
// dropped. It returns false (leaving the type uninitialized) when the record has
// no known layout, so the caller can fall back to the ordinary conversion.
func initBitFieldType(ctx *pkgCtx, typDecl typDecl, this *classCtx, cls clang.Cursor) bool {
	pkg := ctx.pkg
	pos := goNodePos(ctx, cls)

	toGoType := func(t lc.Type) (types.Type, bool) {
		var feats int
		ret := toTypeEx(ctx, pkg.Types, t, flagIsVarDef, &feats, this.scope())
		return ret, feats&featAllIgnore == 0
	}
	layout := planBitFieldLayout(ctx, pos, cls, toGoType)
	if layout.size <= 0 {
		return false // forward declaration: let the caller decide
	}

	if layout.opaque {
		ctx.logf(cls, "[WARN] bit-field record laid out as opaque storage: %s", layout.reason)
		typDecl.InitType(pkg, opaqueBitStruct(ctx, pos, cls, layout))
		return true
	}

	typDecl.InitType(pkg, types.NewStruct(layout.fields, nil))

	// Load nested type/enum/method declarations (but not FieldDecls, already
	// handled by the layout) so they are emitted like for an ordinary record.
	var feats int
	clang.VisitChildren(cls, func(decl, parent clang.Cursor) clang.ChildVisitResult {
		if decl.Kind == lc.Cursor_FieldDecl {
			return clang.Continue
		}
		loadClassMember(ctx, pkg.Types, this, typDecl.Type().Obj().Name(), decl, &feats)
		return clang.Continue
	})

	// Emit the getter/setter pair for every named bit-field in the compile phase
	// (after every type is registered) so member types that reference other
	// records resolve regardless of declaration order.
	if len(layout.accessors) > 0 {
		recvPtr := types.NewPointer(typDecl.Type())
		accessors := layout.accessors
		ctx.addCompileUnit(func(ctx *pkgCtx) {
			ctx.ensureBitHelpers()
			for _, b := range accessors {
				genBitAccessor(ctx, recvPtr, b)
			}
		})
	}

	// Run any methods collected into the class ctx (compileClassImpl is not
	// scheduled on this path, so method bodies are emitted here).
	if len(this.publicMethods) > 0 {
		ctx.addCompileUnit(func(ctx *pkgCtx) {
			compileClassImpl(ctx, this)
		})
	}
	return true
}

// opaqueBitStruct builds the fallback storage for a record whose layout cannot
// be represented field-by-field: an "_xgo_align [0]uintW" marker followed by an
// "_xgo_opaque [S]uint8" byte array of the record's exact size.
func opaqueBitStruct(ctx *pkgCtx, pos token.Pos, cls clang.Cursor, layout bitLayout) *types.Struct {
	pkg := ctx.pkg.Types
	var fields []*types.Var
	elem, warn := alignElem(layout.align)
	if warn {
		ctx.logf(cls, "[WARN] bit-field struct alignment %d > 8 is not fully representable", layout.align)
	}
	fields = append(fields, types.NewField(pos, pkg, bitAlignName, types.NewArray(elem, 0), false))
	fields = append(fields, types.NewField(pos, pkg, bitOpaqueName, types.NewArray(types.Typ[types.Uint8], layout.size), false))
	return types.NewStruct(fields, nil)
}
