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
	"strconv"
	"strings"

	"github.com/goplus/gogen"
	"github.com/goplus/lib/c"
	"github.com/goplus/llcppg/clang"
	lc "github.com/llarhub/clang-c"
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

func newPointer(ctx *pkgCtx, typ types.Type) types.Type {
	switch t := typ.(type) {
	case *types.Basic:
		if t == tyVoid {
			return ctx.unsafePointer()
		}
	case *types.Named:
		/* TODO(xsw):
		if typ == ValistTag {
			return Valist
		} */
	}
	return types.NewPointer(typ)
}

const (
	featHasCallback = 1 << iota
	featExplicitIgnore
	featQuietIgnore
	featAllIgnore = featExplicitIgnore | featQuietIgnore
)

const (
	QuietIgnore = featQuietIgnore
)

var (
	InlineFuncIgnore = featExplicitIgnore
	NoManglingIgnore = featExplicitIgnore
)

func toType(ctx *pkgCtx, pkg *types.Package, typ lc.Type, flags int, scope *scopeCtx) types.Type {
	var feats int
	ret := toTypeEx(ctx, pkg, typ, flags, &feats, scope)
	if feats&featAllIgnore != 0 {
		panic("unsupported type - " + clang.String(typ))
	}
	return ret
}

func toTypeEx(ctx *pkgCtx, pkg *types.Package, typ lc.Type, flags int, feats *int, scope *scopeCtx) types.Type {
	switch typ.Kind {
	case lc.Type_Void:
		return tyVoid
	case lc.Type_Bool:
		return types.Typ[types.Bool]
	case lc.Type_Char_S:
		return ctx.basicTyp(cChar)
	case lc.Type_SChar:
		return types.Typ[types.Int8]
	case lc.Type_Char_U, lc.Type_UChar:
		return types.Typ[types.Uint8]
	case lc.Type_Short:
		return types.Typ[types.Int16]
	case lc.Type_UShort:
		return types.Typ[types.Uint16]
	case lc.Type_Int:
		return ctx.basicTyp(cInt)
	case lc.Type_UInt:
		return ctx.basicTyp(cUint)
	case lc.Type_Long:
		return ctx.basicTyp(cLong)
	case lc.Type_ULong:
		return ctx.basicTyp(cUlong)
	case lc.Type_LongLong:
		return ctx.basicTyp(cLongLong)
	case lc.Type_ULongLong:
		return ctx.basicTyp(cUlongLong)
	case lc.Type_Float:
		return ctx.basicTyp(cFloat)
	case lc.Type_Double:
		return ctx.basicTyp(cDouble)
	case lc.Type_Pointer:
		elem := typ.Pointee()
		if elem.Kind == lc.Type_FunctionProto {
			*feats |= featHasCallback
			return toFuncType(ctx, pkg, elem, feats)
		}
		// flagIsParam only governs the outermost type of a parameter, so clear
		// it before recursing so inner arrays are not wrongly decayed.
		pointee := toTypeEx(ctx, pkg, elem, flagIsTypeDef, feats, scope)
		return newPointer(ctx, pointee)
	case lc.Type_LValueReference, lc.Type_RValueReference: // TODO(xsw): check RVRef is a pointer
		elem := typ.NonReference()
		pointee := toTypeEx(ctx, pkg, elem, flagIsTypeDef, feats, scope)
		return newPointer(ctx, pointee)
	case lc.Type_FunctionProto:
		*feats |= featHasCallback
		return toFuncType(ctx, pkg, typ, feats)
	case lc.Type_Enum, lc.Type_Record, lc.Type_Typedef:
		if t, ok := namedType(ctx, typ, feats); ok {
			return t
		}
	case lc.Type_Elaborated:
		if t, ok := namedType(ctx, typ.Named(), feats); ok {
			return t
		}
	case lc.Type_IncompleteArray, lc.Type_VariableArray:
		// T[] (and VLAs) have no known extent, so they behave like T*. Clear
		// flagIsParam before recursing since decay applies only to this level.
		elem := toTypeEx(ctx, pkg, typ.ArrayElement(), flagIsTypeDef, feats, scope)
		return newPointer(ctx, elem)
	case lc.Type_ConstantArray:
		// A fixed-size C array T[N] is a true array only when it has real
		// storage, e.g. as a struct field. As a function parameter it is a
		// pseudo-array that decays to a pointer T*, so honor that here since
		// libclang reports the parameter type as an array, not a pointer.
		//
		// Decay applies only to the outermost array, so clear flagIsParam
		// before recursing; otherwise a nested array like int matrix[3][4]
		// would decay its inner [4] too, yielding **c.Int instead of *[4]c.Int.
		elem := toTypeEx(ctx, pkg, typ.ArrayElement(), flagIsTypeDef, feats, scope)
		if flags&flagIsParam != 0 {
			return newPointer(ctx, elem)
		}
		return types.NewArray(elem, int64(typ.ArraySize()))
	case lc.Type_Unexposed:
		if t, ok := unexposedType(ctx, typ, feats, scope); ok {
			return t
		}
	case lc.Type_LongDouble:
		return ctx.basicTyp(cLongDouble)
	case lc.Type_WChar:
		return ctx.basicTyp(cWcharT)
	case lc.Type_BlockPointer, lc.Type_Invalid:
		*feats |= featQuietIgnore // will always be ignored
		return types.Typ[types.Invalid]
	}
	if *feats&featQuietIgnore == 0 {
		ctx.logtf(typ, "toType: unsupported type - %s (%d: %s)", cTypeName(typ), typ.Kind, clang.String(typ))
		*feats |= featExplicitIgnore
	}
	return types.Typ[types.Invalid]
}

func namedType(ctx *pkgCtx, typ lc.Type, feats *int) (ret types.Type, found bool) {
	cName := cTypeName(typ)
	if t, ok := ctx.typeAliasOf(cName); ok {
		return t, true
	}
	if o, ok := ctx.getTypeObj(cName, feats); ok {
		return o.Type(), true
	}
	return
}

func unexposedType(ctx *pkgCtx, typ lc.Type, feats *int, scope *scopeCtx) (ret types.Type, found bool) {
	name := clang.String(typ)
	if t, ok := ctx.typeAliasOf(name); ok { // typeAlias supported in config
		return t, true
	}
	if strings.HasPrefix(name, "typename ") || hasPredefinedDirective(name) {
		// typename XXX | clang predefined directive
		*feats |= featQuietIgnore
		return
	}
	name = removeCV(name)
	if o, ok := scope.lookupTypeObj(name); ok { // typeParams
		return o.Type(), true
	}
	cName := cTypeName(typ)
	if o, ok := ctx.getTypeObj(cName, feats); ok {
		return o.Type(), true
	}
	return
}

func toFuncType(ctx *pkgCtx, pkg *types.Package, fn lc.Type, feats *int) *types.Signature {
	params, variadic := toFuncParams(ctx, pkg, fn, feats, nil)
	results := toFuncResults(ctx, pkg, fn.Result(), feats, nil)
	return types.NewSignatureType(nil, nil, nil, params, results, variadic)
}

func toFuncParams(ctx *pkgCtx, pkg *types.Package, fn lc.Type, feats *int, scope *scopeCtx) (ret *types.Tuple, variadic bool) {
	n := fn.NumArgTypes()
	var params []*types.Var
	for i := range n {
		item := fn.Arg(c.Uint(i))
		tyParam := toTypeEx(ctx, pkg, item, flagIsParam, feats, scope)
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

func toFuncResults(ctx *pkgCtx, pkg *types.Package, retType lc.Type, feats *int, scope *scopeCtx) (results *types.Tuple) {
	if retType.Kind != lc.Type_Void {
		tyRet := toTypeEx(ctx, pkg, retType, flagRetType, feats, scope)
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

// basicKind describes the kind of basic type.
type basicKind int

const (
	cVoid basicKind = iota
	cChar
	cInt
	cUint
	cLong
	cUlong
	cLongLong
	cUlongLong
	cFloat
	cDouble
	cLongDouble
	cWcharT
	cBasicMax
)

var ctypBasic = [cBasicMax]string{
	cVoid:       "Void",
	cChar:       "Char",
	cInt:        "Int",
	cUint:       "Uint",
	cLong:       "Long",
	cUlong:      "Ulong",
	cLongLong:   "LongLong",
	cUlongLong:  "UlongLong",
	cFloat:      "Float",
	cDouble:     "Double",
	cLongDouble: "LongDouble",
	cWcharT:     "WcharT",
}

// -----------------------------------------------------------------------------

type typeTag = lc.CursorKind

const (
	tagStruct typeTag = lc.Cursor_StructDecl
	tagUnion  typeTag = lc.Cursor_UnionDecl
	tagClass  typeTag = lc.Cursor_ClassDecl
	tagEnum   typeTag = lc.Cursor_EnumDecl
)

var tagStrvals = [...]string{
	tagStruct: "struct ",
	tagUnion:  "union ",
	tagClass:  "class ",
	tagEnum:   "enum ",
}

// remove type tag prefix, e.g. "struct Foo" => "Foo"
func trimTypeTag(typCName string) string {
	if pos := strings.LastIndex(typCName, " "); pos >= 0 {
		typCName = typCName[pos+1:]
	}
	return typCName
}

// const T, volatile T, const volatile T => T
func removeCV(name string) string {
	for {
		pos := strings.IndexByte(name, ' ')
		if pos > 0 {
			switch name[:pos] {
			case "const", "volatile":
				name = name[pos+1:]
				continue
			}
		}
		return name
	}
}

// __remove_cv(T), __is_integral(T), etc
func hasPredefinedDirective(name string) bool {
	epos := strings.IndexByte(name, '(')
	if epos > 0 {
		if pos := strings.LastIndex(name[:epos], "__"); pos >= 0 {
			if _, ok := clangPredefinedDirectives[name[pos+2:epos]]; ok {
				return true
			}
		}
	}
	return false
}

var clangPredefinedDirectives = map[string]none{
	"decay":                  {},
	"remove_cv":              {},
	"remove_const":           {},
	"remove_volatile":        {},
	"remove_reference_t":     {},
	"remove_extent":          {},
	"remove_all_extents":     {},
	"remove_pointer":         {},
	"add_pointer":            {},
	"add_lvalue_reference":   {},
	"add_rvalue_reference":   {},
	"has_virtual_destructor": {},
	"is_pointer":             {},
	"is_reference":           {},
	"is_const":               {},
	"is_array":               {},
	"is_enum":                {},
	"is_class":               {},
	"is_integral":            {},
	"is_same":                {},
	"is_base_of":             {},
	"is_constructible":       {},
	"is_abstract":            {},
}

// -----------------------------------------------------------------------------
