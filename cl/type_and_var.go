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

	"github.com/goplus/gogen"
	"github.com/goplus/llcppg/clang"
	lc "github.com/goplus/llcppg/lib/clang"
)

// -----------------------------------------------------------------------------

const (
	flagIsParam = 1 << iota
	flagIsStructField
	flagIsExtern
	flagIsTypedef
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
	case *types.Signature:
		panic("todo: newPointer for signature")
		/* if gogen.IsCSignature(t) {
			return types.NewSignature(nil, t.Params(), t.Results(), t.Variadic())
		} */
	case *types.Named:
		panic("todo: newPointer for named type")
		/* if typ == ValistTag {
			return Valist
		} */
	}
	return types.NewPointer(typ)
}

func toType(ctx *blockCtx, typ lc.Type, flags int) types.Type {
	switch typ.Kind {
	case lc.TypeCharS:
		return ctx.c.Ref("Char").Type()
	case lc.TypeInt:
		return ctx.c.Ref("Int").Type()
	case lc.TypeUInt:
		return ctx.c.Ref("Uint").Type()
	case lc.TypePointer:
		pointee := toType(ctx, typ.PointeeType(), flags)
		return newPointer(pointee)
	default:
		log.Println("==> toType: unknown Kind -", typ.Kind)
	}
	panic("todo: toType " + clang.String(typ))
}

// -----------------------------------------------------------------------------

func substObj(pkg *types.Package, scope *types.Scope, origName string, real types.Object) {
	old := scope.Insert(gogen.NewSubst(token.NoPos, pkg, origName, real))
	if old != nil {
		if t, ok := old.Type().(*gogen.TySubst); ok {
			t.Real = real
		} else {
			log.Panicln(origName, "redefined")
		}
	}
}

// -----------------------------------------------------------------------------

func avoidKeyword(name *string) {
	switch *name {
	case "map", "type", "range", "chan", "var", "func", "go", "select",
		"defer", "package", "import", "interface", "fallthrough":
		*name += "_"
	}
}

// -----------------------------------------------------------------------------
