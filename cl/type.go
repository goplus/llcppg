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
	"go/token"
	"go/types"
	"log"
	"strconv"

	"github.com/goplus/gogen"
	"github.com/goplus/lib/c"
	"github.com/goplus/llcppg/clang"
	lc "github.com/goplus/llcppg/lib/clang"
)

// -----------------------------------------------------------------------------

const (
	flagIsParam = 1 << iota
	flagIsVarDef
	flagIsTypeDef
	flagRetType
)

var (
	tyVoid = types.Typ[types.UntypedNil]
)

func newPointer(typ types.Type) types.Type {
	switch t := typ.(type) {
	case *types.Basic:
		if t == tyVoid {
			return types.Typ[types.UnsafePointer]
		}
	case *types.Named:
		/* TODO(xsw):
		if typ == ValistTag {
			return Valist
		} */
	}
	return types.NewPointer(typ)
}

func toType(ctx *pkgCtx, pkg *types.Package, typ lc.Type, flags int) types.Type {
	switch typ.Kind {
	case lc.TypeVoid:
		return tyVoid
	case lc.TypeBool:
		return types.Typ[types.Bool]
	case lc.TypeCharS:
		return ctx.c.Ref("Char").Type()
	case lc.TypeSChar:
		return types.Typ[types.Int8]
	case lc.TypeCharU, lc.TypeUChar:
		return types.Typ[types.Uint8]
	case lc.TypeShort:
		return types.Typ[types.Int16]
	case lc.TypeUShort:
		return types.Typ[types.Uint16]
	case lc.TypeInt:
		return ctx.c.Ref("Int").Type()
	case lc.TypeUInt:
		return ctx.c.Ref("Uint").Type()
	case lc.TypeLong:
		return ctx.c.Ref("Long").Type()
	case lc.TypeULong:
		return ctx.c.Ref("Ulong").Type()
	case lc.TypeLongLong:
		return ctx.c.Ref("LongLong").Type()
	case lc.TypeULongLong:
		return ctx.c.Ref("UlongLong").Type()
	case lc.TypeFloat:
		return ctx.c.Ref("Float").Type()
	case lc.TypeDouble:
		return ctx.c.Ref("Double").Type()
	case lc.TypePointer:
		elem := typ.PointeeType()
		if elem.Kind == lc.TypeFunctionProto {
			return toFuncType(ctx, pkg, elem)
		}
		// flagIsParam only governs the outermost type of a parameter, so clear
		// it before recursing so inner arrays are not wrongly decayed.
		pointee := toType(ctx, pkg, elem, flagIsTypeDef)
		return newPointer(pointee)
	case lc.TypeFunctionProto:
		return toFuncType(ctx, pkg, typ)
	case lc.TypeEnum:
		cName := clang.String(typ)
		if t, ok := ctx.typeOf(cName); ok {
			return t
		}
	case lc.TypeRecord, lc.TypeTypedef:
		// A record type registers under its declaration's type spelling
		// (clang.String(decl.Type()); see loadClass/emitUnion). Resolve through
		// the same key via the type's declaration cursor, which also covers a
		// tagless inline union hoisted to "_llcppg_union_<n>" whose tag-less
		// fullName would otherwise miss. See issue goplus/llcppg#775.
		cName := clang.String(typ.TypeDeclaration().Type())
		if t, ok := ctx.typeOf(cName); ok {
			return t
		}
	case lc.TypeElaborated:
		cName := clang.String(typ.NamedType())
		if t, ok := ctx.typeOf(cName); ok {
			return t
		}
	case lc.TypeConstantArray:
		// A fixed-size C array T[N] is a true array only when it has real
		// storage, e.g. as a struct field. As a function parameter it is a
		// pseudo-array that decays to a pointer T*, so honor that here since
		// libclang reports the parameter type as an array, not a pointer.
		//
		// Decay applies only to the outermost array, so clear flagIsParam
		// before recursing; otherwise a nested array like int matrix[3][4]
		// would decay its inner [4] too, yielding **c.Int instead of *[4]c.Int.
		elem := toType(ctx, pkg, typ.ArrayElementType(), flagIsTypeDef)
		if flags&flagIsParam != 0 {
			return newPointer(elem)
		}
		return types.NewArray(elem, int64(typ.ArraySize()))
	case lc.TypeIncompleteArray, lc.TypeVariableArray:
		// T[] (and VLAs) have no known extent, so they behave like T*. Clear
		// flagIsParam before recursing since decay applies only to this level.
		elem := toType(ctx, pkg, typ.ArrayElementType(), flagIsTypeDef)
		return newPointer(elem)
	default:
		log.Println("==> toType: unknown Kind -", typ.Kind)
	}
	panic("todo: toType " + clang.String(typ))
}

func toFuncType(ctx *pkgCtx, pkg *types.Package, fn lc.Type) *types.Signature {
	params, variadic := toFuncParams(ctx, pkg, fn)
	results := toFuncResults(ctx, pkg, fn.ResultType())
	return types.NewSignatureType(nil, nil, nil, params, results, variadic)
}

func toFuncParams(ctx *pkgCtx, pkg *types.Package, fn lc.Type) (ret *types.Tuple, variadic bool) {
	n := fn.NumArgTypes()
	var params []*types.Var
	for i := range n {
		item := fn.ArgType(c.Uint(i))
		tyParam := toType(ctx, pkg, item, flagIsParam)
		nameParam := "_llcppg_param" + strconv.Itoa(int(i)+1)
		params = append(params, types.NewParam(token.NoPos, pkg, nameParam, tyParam))
	}
	variadic = fn.IsFunctionTypeVariadic() != 0
	if variadic {
		params = append(params, newVariadicParam(pkg))
	}
	ret = types.NewTuple(params...)
	return
}

var (
	tyValist types.Type = types.NewSlice(gogen.TyAny)
)

func newVariadicParam(pkg *types.Package) *types.Var {
	return types.NewParam(token.NoPos, pkg, "__llgo_va_list", tyValist)
}

func toFuncResults(ctx *pkgCtx, pkg *types.Package, retType lc.Type) (results *types.Tuple) {
	if retType.Kind != lc.TypeVoid {
		tyRet := toType(ctx, pkg, retType, flagRetType)
		results = types.NewTuple(types.NewParam(token.NoPos, pkg, "", tyRet))
	}
	return
}

// -----------------------------------------------------------------------------

func cmpType(ta, tb lc.Type) int {
	// TODO(xsw): c++ overload support
	return int(ta.Kind - tb.Kind)
}

// -----------------------------------------------------------------------------
