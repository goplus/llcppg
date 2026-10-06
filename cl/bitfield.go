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

	"github.com/goplus/gogen"
	"github.com/goplus/llcppg/clang"
	lc "github.com/llarhub/clang-c"
)

// bitFieldInfo holds metadata about a single C/C++ bit-field.
type bitFieldInfo struct {
	cursor    clang.Cursor // The FieldDecl cursor
	name      string       // C name of the field
	offset    int64        // Bit offset from start of struct/union
	width     int64        // Bit width
	fieldType types.Type   // Go type of the declared C type
	signed    bool         // Whether the declared type is signed
	isBool    bool         // Whether the type is _Bool or C++ bool
}

// bitFieldRun represents a maximal sequence of consecutive bit-fields.
type bitFieldRun struct {
	fields       []*bitFieldInfo // All fields in this run (named and unnamed)
	storageIndex int             // Index for naming (_xgo_bits_<N>)
	byteOffset   int64           // Byte offset where the run starts
	byteLen      int64           // Number of bytes spanned by the run
}

// hasBitFields reports whether the record declaration has any bit-fields.
func hasBitFields(cursor clang.Cursor) bool {
	found := false
	clang.VisitChildren(cursor, func(decl, parent clang.Cursor) clang.ChildVisitResult {
		if decl.Kind == lc.Cursor_FieldDecl && decl.FieldDeclBitWidth() != 0 {
			found = true
			return clang.Break
		}
		return clang.Continue
	})
	return found
}

// collectBitFields collects all bit-fields from a struct/union cursor
// and returns them grouped into runs of consecutive bit-fields.
func collectBitFields(ctx *pkgCtx, pkgTypes *types.Package, cursor clang.Cursor) []*bitFieldRun {
	var allFields []*bitFieldInfo
	
	// First pass: collect all bit-fields
	clang.VisitChildren(cursor, func(decl, parent clang.Cursor) clang.ChildVisitResult {
		if decl.Kind != lc.Cursor_FieldDecl {
			return clang.Continue
		}
		
		bitWidth := decl.FieldDeclBitWidth()
		if bitWidth == 0 {
			return clang.Continue // Not a bit-field
		}
		
		offset := decl.OffsetOfField() // Offset in bits
		origName := clang.String(decl)
		
		// Try to convert the field type
		var fldType types.Type
		var signed, isBool bool
		var feats int
		
		fldCType := decl.Type()
		fldType = toTypeEx(ctx, pkgTypes, fldCType, flagIsVarDef, &feats, nil)
		if feats&featAllIgnore != 0 {
			// Can't convert this type
			fldType = nil
		}
		
		// Determine signedness for later use
		if fldType != nil {
			switch ut := fldType.(type) {
			case *types.Basic:
				info := ut.Info()
				signed = (info & types.IsUnsigned) == 0
				isBool = (info & types.IsBoolean) != 0
			case *types.Named:
				if basic, ok := ut.Underlying().(*types.Basic); ok {
					info := basic.Info()
					signed = (info & types.IsUnsigned) == 0
					isBool = (info & types.IsBoolean) != 0
				}
			}
		}
		
		allFields = append(allFields, &bitFieldInfo{
			cursor:    decl,
			name:      origName,
			offset:    offset,
			width:     int64(bitWidth),
			fieldType: fldType,
			signed:    signed,
			isBool:    isBool,
		})
		
		return clang.Continue
	})
	
	if len(allFields) == 0 {
		return nil
	}
	
	// Second pass: group into runs
	var runs []*bitFieldRun
	var currentRun *bitFieldRun
	
	for _, field := range allFields {
		if currentRun == nil {
			// Start first run
			currentRun = &bitFieldRun{
				fields:       []*bitFieldInfo{field},
				storageIndex: 0,
				byteOffset:   field.offset / 8,
			}
		} else {
			// Check if field is consecutive with current run
			// Fields are consecutive if they're in the same or adjacent bytes
			lastField := currentRun.fields[len(currentRun.fields)-1]
			lastEnd := (lastField.offset + lastField.width + 7) / 8
			thisStart := field.offset / 8
			
			if thisStart < lastEnd {
				// Same run
				currentRun.fields = append(currentRun.fields, field)
			} else {
				// End current run and start new one
				runs = append(runs, currentRun)
				currentRun = &bitFieldRun{
					fields:       []*bitFieldInfo{field},
					storageIndex: len(runs),
					byteOffset:   thisStart,
				}
			}
		}
	}
	
	// Add the final run
	if currentRun != nil {
		runs = append(runs, currentRun)
	}
	
	// Calculate byte length for each run
	for _, run := range runs {
		minByte := run.byteOffset
		maxByte := int64(0)
		for _, f := range run.fields {
			if f.width > 0 {
				end := (f.offset + f.width + 7) / 8
				if end > maxByte {
					maxByte = end
				}
			}
		}
		run.byteLen = maxByte - minByte
	}
	
	return runs
}

// genBitFieldStorage generates a storage field for a bit-field run.
func genBitFieldStorage(ctx *pkgCtx, run *bitFieldRun) *types.Var {
	pkgTypes := ctx.pkg.Types
	fieldName := fmt.Sprintf("_xgo_bits_%d", run.storageIndex)
	arrayType := types.NewArray(types.Typ[types.Uint8], run.byteLen)
	return types.NewField(0, pkgTypes, fieldName, arrayType, false)
}

// genBitFieldAccessors generates getter/setter methods for bit-fields in a run.
func genBitFieldAccessors(ctx *pkgCtx, recvPtr types.Type, run *bitFieldRun) {
	pkg := ctx.pkg
	pkgTypes := pkg.Types
	
	for _, field := range run.fields {
		if field.name == "" || field.fieldType == nil {
			// Skip unnamed and unconvertible bit-fields
			continue
		}
		
		// Check for name conflicts with existing fields/methods
		if recvPtrNamed, ok := recvPtr.(*types.Pointer); ok {
			if recvNamed, ok := recvPtrNamed.Elem().(*types.Named); ok {
				if _, exists := recvNamed.Underlying().(*types.Struct); exists {
					// Check if field name collides
					if isField, exists := existMember(recvNamed, field.name); exists {
						if isField {
							// Skip if collides with field name
							continue
						}
					}
				}
			}
		}
		
		// Generate getter
		getterName := fmt.Sprintf("XGof_get_%s", field.name)
		recv := types.NewParam(0, pkgTypes, "p", recvPtr)
		results := types.NewTuple(types.NewParam(0, pkgTypes, "", field.fieldType))
		getterSig := types.NewSignatureType(recv, nil, nil, nil, results, false)
		
		getterFunc, err := pkg.NewFuncWith(0, getterName, getterSig, nil)
		if err != nil {
			ctx.panicf(field.cursor, "bitfield %s: failed to create getter - %v", field.name, err)
		}
		
		genBitFieldGetterBody(ctx, getterFunc, field)
		
		// Generate setter
		setterName := fmt.Sprintf("XGof_set_%s", field.name)
		params := types.NewTuple(types.NewParam(0, pkgTypes, "v", field.fieldType))
		setterSig := types.NewSignatureType(recv, nil, nil, params, nil, false)
		
		setterFunc, err := pkg.NewFuncWith(0, setterName, setterSig, nil)
		if err != nil {
			ctx.panicf(field.cursor, "bitfield %s: failed to create setter - %v", field.name, err)
		}
		
		genBitFieldSetterBody(ctx, setterFunc, field)
	}
}

// genBitFieldGetterBody generates the body for a bit-field getter.
// Generated code: return <type>(_xgo_bitget[_signed](unsafe.Pointer(p), offset, width))
func genBitFieldGetterBody(ctx *pkgCtx, f *gogen.Func, field *bitFieldInfo) {
	// For now, stub implementation that returns zero
	// TODO: Implement actual bit extraction logic
	cb := f.BodyStart(ctx.pkg)
	cb.ZeroLit(field.fieldType).Return(1).End()
}

// genBitFieldSetterBody generates the body for a bit-field setter.
// Generated code: _xgo_bitset(unsafe.Pointer(p), offset, width, uint64(v))
func genBitFieldSetterBody(ctx *pkgCtx, f *gogen.Func, field *bitFieldInfo) {
	// For now, stub implementation that does nothing
	// TODO: Implement actual bit setting logic
	cb := f.BodyStart(ctx.pkg)
	cb.End()
}

// ensureBitFieldHelpers generates the bit manipulation helper functions
// if there are any bit-fields in the package.
func ensureBitFieldHelpers(ctx *pkgCtx) {
	if ctx.hasBitFieldHelpers {
		return
	}
	ctx.hasBitFieldHelpers = true
	
	pkg := ctx.pkg
	pkgTypes := pkg.Types
	
	// Helper 1: _xgo_bitget(base unsafe.Pointer, off, width uintptr) uint64
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
		
		// Simplified stub for now
		cb := f.BodyStart(pkg)
		cb.Val(uint64(0)).Return(1).End()
	}
	
	// Helper 2: _xgo_bitget_signed(base unsafe.Pointer, off, width uintptr) int64
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
		
		// Simplified stub for now
		cb := f.BodyStart(pkg)
		cb.Val(int64(0)).Return(1).End()
	}
	
	// Helper 3: _xgo_bitset(base unsafe.Pointer, off, width uintptr, v uint64)
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
		
		// Simplified stub for now
		cb := f.BodyStart(pkg)
		cb.End()
	}
}
