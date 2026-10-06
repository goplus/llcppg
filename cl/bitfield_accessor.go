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

	"github.com/goplus/gogen"
)

// Accessor and helper emission for C/C++ bit-fields. See bitfield.go for the
// layout half and issue goplus/llcppg#770 for the design.
//
// For every named bit-field "foo" of width w at bit offset o (from the start of
// X), this file emits the virtual-field getter/setter pair
//
//	func (p *X) XGof_get_foo() T  { return T(_xgo_bitget(unsafe.Pointer(p), o, w)) }
//	func (p *X) XGof_set_foo(v T) { _xgo_bitset(unsafe.Pointer(p), o, w, uint64(v)) }
//
// using _xgo_bitget_signed for a signed T, and a bool test / 0-1 store for a
// bool T. The three unexported helpers _xgo_bitget, _xgo_bitget_signed and
// _xgo_bitset are emitted once per package (ensureBitHelpers); they operate on
// bits numbered from the least significant bit of each byte (little-endian).

const (
	bitGetPrefix = "XGof_get_"
	bitSetPrefix = "XGof_set_"
)

// -----------------------------------------------------------------------------

// genBitAccessor emits the getter/setter pair for one named bit-field.
//
// A bit-field whose name collides with a real field or method of X (or is
// otherwise not a usable Go identifier) gets no accessor and is reported, per
// section 7 of the proposal: its bits still belong to the run, so the layout is
// unaffected. The getter is attempted first; if it fails the setter is skipped
// too, so the virtual field is never left half-defined.
func genBitAccessor(ctx *pkgCtx, recvPtr types.Type, b bitAccessor) {
	if err := genBitGetter(ctx, recvPtr, b); err != nil {
		ctx.ignoref(featQuietIgnore, b.decl, "bit-field %s: no accessor (%v)", b.name, err)
		return
	}
	if err := genBitSetter(ctx, recvPtr, b); err != nil {
		ctx.ignoref(featQuietIgnore, b.decl, "bit-field %s: setter skipped (%v)", b.name, err)
	}
}

func genBitGetter(ctx *pkgCtx, recvPtr types.Type, b bitAccessor) error {
	pkg := ctx.pkg
	pkgTypes := pkg.Types

	recv := types.NewParam(token.NoPos, pkgTypes, "p", recvPtr)
	results := types.NewTuple(types.NewParam(token.NoPos, pkgTypes, "", b.goType))
	sig := types.NewSignatureType(recv, nil, nil, nil, results, false)
	f, err := pkg.NewFuncWith(token.NoPos, bitGetPrefix+b.name, sig, nil)
	if err != nil {
		return err
	}
	f.SetComments(pkg, bitDocComment(bitFieldDoc(b.decl)))

	cb := f.BodyStart(pkg)
	if b.isBool {
		bitfieldGetCall(ctx, cb, "Unsigned", b)
		cb.Val(0).BinaryOp(token.NEQ)
	} else {
		helper := "Unsigned"
		if b.signed {
			helper = "Signed"
		}
		cb.Typ(b.goType)
		bitfieldGetCall(ctx, cb, helper, b)
		cb.Call(1)
	}
	cb.Return(1).End()
	return nil
}

func genBitSetter(ctx *pkgCtx, recvPtr types.Type, b bitAccessor) error {
	pkg := ctx.pkg
	pkgTypes := pkg.Types

	recv := types.NewParam(token.NoPos, pkgTypes, "p", recvPtr)
	params := types.NewTuple(types.NewParam(token.NoPos, pkgTypes, "v", b.goType))
	sig := types.NewSignatureType(recv, nil, nil, params, nil, false)
	f, err := pkg.NewFuncWith(token.NoPos, bitSetPrefix+b.name, sig, nil)
	if err != nil {
		return err
	}
	f.SetComments(pkg, bitDocComment(bitFieldDoc(b.decl)))

	cb := f.BodyStart(pkg)
	// bitfield.Set(unsafe.Pointer(p), o, w, uint64(v))  [bool: uint64(b2u(v))]
	cb.Val(ctx.bitfieldRef("Set"))
	pushUnsafePtrOfP(ctx, cb)
	cb.Val(int(b.off))
	cb.Val(int(b.width))
	cb.Typ(uint64T())
	if b.isBool {
		pushBoolToUint64(ctx, cb)
	} else {
		cb.VarVal("v")
	}
	cb.Call(1) // uint64(...)
	cb.Call(4) // bitfield.Set(...)
	cb.EndStmt().End()
	return nil
}

func bitfieldGetCall(ctx *pkgCtx, cb *gogen.CodeBuilder, helper string, b bitAccessor) {
	cb.Val(ctx.bitfieldRef(helper))
	pushUnsafePtrOfP(ctx, cb)
	cb.Val(int(b.off))
	cb.Val(int(b.width))
	cb.Call(3)
}

// pushUnsafePtrOfP pushes "unsafe.Pointer(p)".
func pushUnsafePtrOfP(ctx *pkgCtx, cb *gogen.CodeBuilder) {
	cb.Typ(ctx.unsafePointer()).VarVal("p").Call(1)
}

// pushBoolToUint64 pushes an expression that evaluates the bool parameter "v"
// to a uint64 (1 for true, 0 for false), via an immediately-invoked closure
// since Go has no direct bool-to-integer conversion.
func pushBoolToUint64(ctx *pkgCtx, cb *gogen.CodeBuilder) {
	pkg := ctx.pkg
	results := types.NewTuple(types.NewParam(token.NoPos, pkg.Types, "", uint64T()))
	cb.NewClosure(nil, results, false).BodyStart(pkg).
		If().VarVal("v").Then().
		/**/ Val(1).Return(1).EndStmt().
		End().
		Val(0).Return(1).
		End()
	cb.Call(0)
}

// -----------------------------------------------------------------------------

func (p *pkgCtx) bitfieldRef(name string) types.Object {
	return p.pkg.Import("github.com/qiniu/x/bitfield").Ref(name)
}

func uintptrT() types.Type { return types.Typ[types.Uintptr] }
func uint64T() types.Type  { return types.Typ[types.Uint64] }
func uint8T() types.Type   { return types.Typ[types.Uint8] }

func bitHelperSig(ctx *pkgCtx, extraParam *types.Var, results *types.Tuple) *types.Signature {
	pkgTypes := ctx.pkg.Types
	ps := []*types.Var{
		types.NewParam(token.NoPos, pkgTypes, "base", ctx.unsafePointer()),
		types.NewParam(token.NoPos, pkgTypes, "off", uintptrT()),
		types.NewParam(token.NoPos, pkgTypes, "width", uintptrT()),
	}
	if extraParam != nil {
		ps = append(ps, extraParam)
	}
	return types.NewSignatureType(nil, nil, nil, types.NewTuple(ps...), results, false)
}

// pushByteSlice pushes "unsafe.Slice((*uint8)(base), (off+width+7)/8)".
func pushByteSlice(ctx *pkgCtx, cb *gogen.CodeBuilder) {
	cb.Val(ctx.pkg.Import("unsafe").Ref("Slice"))
	cb.Typ(types.NewPointer(uint8T())).VarVal("base").Call(1) // (*uint8)(base)
	// (off + width + 7) / 8
	cb.VarVal("off").VarVal("width").BinaryOp(token.ADD).
		Val(7).BinaryOp(token.ADD).
		Val(8).BinaryOp(token.QUO)
	cb.Call(2)
}

// pushBitTest pushes "s[(off+i)/8] & (uint8(1) << ((off+i)%8))".
func pushBitTest(ctx *pkgCtx, cb *gogen.CodeBuilder) {
	// s[(off+i)/8]
	cb.VarVal("s")
	cb.VarVal("off").VarVal("i").BinaryOp(token.ADD).Val(8).BinaryOp(token.QUO)
	cb.Index(1, 1)
	// uint8(1) << ((off+i)%8)
	cb.Typ(uint8T()).Val(1).Call(1)
	cb.VarVal("off").VarVal("i").BinaryOp(token.ADD).Val(8).BinaryOp(token.REM)
	cb.BinaryOp(token.SHL)
	// s[...] & (...)
	cb.BinaryOp(token.AND)
}

// bitDocComment builds a doc-only comment group holding a single "// text" line
// above the declaration (the leading "\n" leaves a blank line before it).
func bitDocComment(text string) *ast.CommentGroup {
	return &ast.CommentGroup{List: []*ast.Comment{{Text: "\n// " + text}}}
}
