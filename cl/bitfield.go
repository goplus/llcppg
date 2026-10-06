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
	"go/types"

	"github.com/goplus/llcppg/clang"
	lc "github.com/llarhub/clang-c"
)

// bitFieldMember represents a single bit-field with its metadata
type bitFieldMember struct {
	cursor    clang.Cursor  // The FieldDecl cursor
	name      string        // C name
	offset    int64         // Bit offset from start of struct/union
	width     int64         // Bit width  
	typ       types.Type    // Go type of the declared C type
	signed    bool          // Whether type is signed
}

// bitFieldStorage represents the storage for a run of bit-fields
type bitFieldStorage struct {
	index      int             // 0-based index for naming
	byteOffset int64           // Where this storage starts (in bytes)
	byteSize   int64           // Size in bytes
	members    []*bitFieldMember // Named members that fit in this storage
}

// hasBitFields checks if a struct/union has any bit-fields
func hasBitFields(cursor clang.Cursor) bool {
	result := false
	clang.VisitChildren(cursor, func(child, parent clang.Cursor) clang.ChildVisitResult {
		if child.Kind == lc.Cursor_FieldDecl && child.FieldDeclBitWidth() > 0 {
			result = true
			return clang.Break
		}
		return clang.Continue
	})
	return result
}

// collectBitFieldMembers collects all bit-field members from a struct/union
func collectBitFieldMembers(ctx *pkgCtx, pkgTypes *types.Package, cursor clang.Cursor) []*bitFieldMember {
	var members []*bitFieldMember
	
	clang.VisitChildren(cursor, func(child, parent clang.Cursor) clang.ChildVisitResult {
		if child.Kind != lc.Cursor_FieldDecl {
			return clang.Continue
		}
		
		width := child.FieldDeclBitWidth()
		if width < 0 {
			return clang.Continue // Not a bit-field (regular field)
		}
		
		offset := child.OffsetOfField() // In bits
		name := clang.String(child)
		
		// Try to convert the declared type
		var typ types.Type
		var signed bool
		var feats int
		
		typ = toTypeEx(ctx, pkgTypes, child.Type(), flagIsVarDef, &feats, nil)
		if feats&featAllIgnore != 0 {
			typ = nil  // Can't convert this type
		}
		
		// Determine signedness
		if typ != nil {
			if basic, ok := typ.(*types.Basic); ok {
				signed = (basic.Info() & types.IsUnsigned) == 0
			} else if named, ok := typ.(*types.Named); ok {
				if basic, ok := named.Underlying().(*types.Basic); ok {
					signed = (basic.Info() & types.IsUnsigned) == 0
				}
			}
		}
		
		members = append(members, &bitFieldMember{
			cursor: child,
			name:   name,
			offset: offset,
			width:  int64(width),
			typ:    typ,
			signed: signed,
		})
		
		return clang.Continue
	})
	
	return members
}

// groupBitFieldStorage groups bit-field members into storage runs.
// Returns storage runs and which members are named (have names and convertible types).
func groupBitFieldStorage(members []*bitFieldMember) ([]*bitFieldStorage, []*bitFieldMember) {
	if len(members) == 0 {
		return nil, nil
	}
	
	// Sort members by offset (should already be in order from visiting, but be safe)
	// Group by byte offset ranges
	
	var storages []*bitFieldStorage
	var named []*bitFieldMember
	
	// Process members from first to last
	var currentStorage *bitFieldStorage
	lastEndByte := int64(-1)
	storageIndex := 0
	
	for _, member := range members {
		startByte := member.offset / 8
		
		// Check if we need to start a new storage
		// New storage if:
		// 1. This is the first member
		// 2. There's a gap (startByte > lastEndByte)
		if currentStorage == nil || startByte > lastEndByte {
			// Finalize previous storage if any
			if currentStorage != nil {
				storages = append(storages, currentStorage)
				storageIndex++
			}
			
			// Start new storage
			currentStorage = &bitFieldStorage{
				index:      storageIndex,
				byteOffset: startByte,
				members:    []*bitFieldMember{},
			}
		}
		
		// Add member to current storage
		currentStorage.members = append(currentStorage.members, member)
		
		// Update end position
		endByte := (member.offset + member.width + 7) / 8
		if endByte > lastEndByte {
			lastEndByte = endByte
		}
		
		// Track named members for accessor generation
		if member.name != "" && member.typ != nil {
			named = append(named, member)
		}
	}
	
	// Finalize last storage
	if currentStorage != nil {
		currentStorage.byteSize = lastEndByte - currentStorage.byteOffset
		storages = append(storages, currentStorage)
	}
	
	return storages, named
}

// genBitFieldStorageField creates a storage field for a bit-field storage run
func genBitFieldStorageField(ctx *pkgCtx, storage *bitFieldStorage) *types.Var {
	fieldName := fmt.Sprintf("_xgo_bits_%d", storage.index)
	arrayType := types.NewArray(types.Typ[types.Uint8], storage.byteSize)
	return types.NewField(0, ctx.pkg.Types, fieldName, arrayType, false)
}

// genBitFieldAccessors generates getter/setter methods for named bit-fields
func genBitFieldAccessors(ctx *pkgCtx, recvPtr types.Type, storages []*bitFieldStorage) {
	for _, storage := range storages {
		for _, member := range storage.members {
			if member.name == "" || member.typ == nil {
				continue  // Skip unnamed or unconvertible
			}
			
			genBitFieldGetter(ctx, recvPtr, storage, member)
			genBitFieldSetter(ctx, recvPtr, storage, member)
		}
	}
}

// genBitFieldGetter generates a getter for a bit-field member
func genBitFieldGetter(ctx *pkgCtx, recvPtr types.Type, storage *bitFieldStorage, member *bitFieldMember) {
	pkg := ctx.pkg
	pkgTypes := pkg.Types
	
	name := fmt.Sprintf("XGof_get_%s", member.name)
	recv := types.NewParam(0, pkgTypes, "p", recvPtr)
	results := types.NewTuple(types.NewParam(0, pkgTypes, "", member.typ))
	sig := types.NewSignatureType(recv, nil, nil, nil, results, false)
	
	f, err := pkg.NewFuncWith(0, name, sig, nil)
	if err != nil {
		ctx.panicf(member.cursor, "failed to create bitfield getter %s: %v", member.name, err)
	}
	
	// For now, just return a zero value
	// TODO: Implement actual bit extraction
	cb := f.BodyStart(pkg)
	cb.ZeroLit(member.typ).Return(1).End()
}

// genBitFieldSetter generates a setter for a bit-field member
func genBitFieldSetter(ctx *pkgCtx, recvPtr types.Type, storage *bitFieldStorage, member *bitFieldMember) {
	pkg := ctx.pkg
	pkgTypes := pkg.Types
	
	name := fmt.Sprintf("XGof_set_%s", member.name)
	recv := types.NewParam(0, pkgTypes, "p", recvPtr)
	params := types.NewTuple(types.NewParam(0, pkgTypes, "v", member.typ))
	sig := types.NewSignatureType(recv, nil, nil, params, nil, false)
	
	f, err := pkg.NewFuncWith(0, name, sig, nil)
	if err != nil {
		ctx.panicf(member.cursor, "failed to create bitfield setter %s: %v", member.name, err)
	}
	
	// For now, just empty body
	// TODO: Implement actual bit setting
	cb := f.BodyStart(pkg)
	cb.End()
}

// ensureBitFieldHelpers generates the bit manipulation helper functions
func ensureBitFieldHelpers(ctx *pkgCtx) {
	if ctx.hasBitFieldHelpers {
		return
	}
	ctx.hasBitFieldHelpers = true
	
	pkg := ctx.pkg
	pkgTypes := pkg.Types
	
	// _xgo_bitget helper
	{
		params := types.NewTuple(
			types.NewParam(0, pkgTypes, "base", ctx.unsafePointer()),
			types.NewParam(0, pkgTypes, "off", types.Typ[types.Uintptr]),
			types.NewParam(0, pkgTypes, "width", types.Typ[types.Uintptr]),
		)
		results := types.NewTuple(types.NewParam(0, pkgTypes, "", types.Typ[types.Uint64]))
		sig := types.NewSignatureType(nil, nil, nil, params, results, false)
		
		f, err := pkg.NewFuncWith(0, "_xgo_bitget", sig, nil)
		if err != nil {
			panic(fmt.Sprintf("failed to create _xgo_bitget: %v", err))
		}
		
		cb := f.BodyStart(pkg)
		cb.ZeroLit(types.Typ[types.Uint64]).Return(1).End()
	}
	
	// _xgo_bitget_signed helper
	{
		params := types.NewTuple(
			types.NewParam(0, pkgTypes, "base", ctx.unsafePointer()),
			types.NewParam(0, pkgTypes, "off", types.Typ[types.Uintptr]),
			types.NewParam(0, pkgTypes, "width", types.Typ[types.Uintptr]),
		)
		results := types.NewTuple(types.NewParam(0, pkgTypes, "", types.Typ[types.Int64]))
		sig := types.NewSignatureType(nil, nil, nil, params, results, false)
		
		f, err := pkg.NewFuncWith(0, "_xgo_bitget_signed", sig, nil)
		if err != nil {
			panic(fmt.Sprintf("failed to create _xgo_bitget_signed: %v", err))
		}
		
		cb := f.BodyStart(pkg)
		cb.ZeroLit(types.Typ[types.Int64]).Return(1).End()
	}
	
	// _xgo_bitset helper
	{
		params := types.NewTuple(
			types.NewParam(0, pkgTypes, "base", ctx.unsafePointer()),
			types.NewParam(0, pkgTypes, "off", types.Typ[types.Uintptr]),
			types.NewParam(0, pkgTypes, "width", types.Typ[types.Uintptr]),
			types.NewParam(0, pkgTypes, "v", types.Typ[types.Uint64]),
		)
		sig := types.NewSignatureType(nil, nil, nil, params, nil, false)
		
		f, err := pkg.NewFuncWith(0, "_xgo_bitset", sig, nil)
		if err != nil {
			panic(fmt.Sprintf("failed to create _xgo_bitset: %v", err))
		}
		
		cb := f.BodyStart(pkg)
		cb.End()
	}
}
