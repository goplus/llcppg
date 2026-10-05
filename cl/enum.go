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
	"go/types"

	"github.com/goplus/gogen"
	"github.com/goplus/llcppg/clang"
	lc "github.com/llarhub/clang-c"
)

// -----------------------------------------------------------------------------

// loadEnum translates a C/C++ enum declaration into Go declarations.
//
// A named enum is emitted as a Go named type (whose underlying type is the C
// int the enum decays to) followed by a const block whose constants are typed
// with that named type, so the enum type name is preserved. An anonymous enum
// has no type name to preserve, so only untyped consts are emitted.
//
// Whether the enum is global, in a namespace, or nested inside a class only
// affects naming: ns carries the enclosing namespace/class prefix (e.g. "bar_"
// or "Shape_"), so a constant Red becomes bar_Red / Shape_Red and then goes
// through getPubName for the final Go name.
func loadEnum(ctx *pkgCtx, decl clang.Cursor) {
	ns := ctx.cNS(decl)
	hasName := decl.IsAnonymous() == 0
	scoped := decl.EnumDeclIsScoped() != 0

	var cName string
	if hasName {
		cName = cNameWithNS(cBaseName(decl), ns)
	}

	if debugCompileDecl {
		ctx.logf(decl, "enum %s - hasName: %v, scoped: %v", cName, hasName, scoped)
	}

	var ok bool
	var typDecl typDecl
	if hasName {
		typDecl, ok = ctx.typdecls[cName]
		if !ok {
			goName := ctx.typeName(cName, true)
			typDecl = newEnumType(ctx, decl, cName, goName)
			ctx.typdecls[cName] = typDecl
		}
	}

	if decl.IsCursorDefinition() == 0 {
		return // declaration only, no definition
	}

	if scoped {
		ns = cName
	}
	initEnumvals(ctx, decl, typDecl, ns)
}

func emitEnum(ctx *pkgCtx, decl clang.Cursor, goName string) *types.Named {
	typDecl := newEnumType(ctx, decl, "", goName)
	initEnumvals(ctx, decl, typDecl, ctx.cNS(decl))
	return typDecl.Type()
}

func newEnumType(ctx *pkgCtx, decl clang.Cursor, cName, goName string) typDecl {
	typDecl := newType(ctx, decl, cName, goName)

	pkg := ctx.pkg
	pkgTypes := pkg.Types

	feats := 0
	typ := decl.EnumDeclIntegerType()
	underType := toTypeEx(ctx, pkgTypes, typ, flagIsTypeDef, &feats, nil)
	if feats&featAllIgnore != 0 {
		ctx.panicf(decl, "enum %s: unsupported underlying type %s (%d: %s)", cNameOf(decl), cTypeName(typ), typ.Kind, clang.String(typ))
	}
	typDecl.InitType(pkg, underType)
	return typDecl
}

func initEnumvals(ctx *pkgCtx, decl clang.Cursor, typDecl typDecl, ns string) {
	pkg := ctx.pkg
	pkgTypes := pkg.Types
	defs := pkg.NewConstDefs(pkgTypes.Scope())

	var enumType types.Type
	if typDecl.defs != nil {
		if doc := ctx.docCommentGroup(decl); doc != nil {
			typDecl.defs.SetComments(doc)
		}
		enumType = typDecl.Type()
	} else {
		if doc := ctx.docCommentGroup(decl); doc != nil {
			defs.SetComments(doc)
		}
	}

	clang.VisitChildren(decl, func(item, parent clang.Cursor) clang.ChildVisitResult {
		if item.Kind != lc.Cursor_EnumConstantDecl {
			return clang.Continue
		}
		name := ctx.enumvalName(clang.String(item), ns)
		val := int(item.EnumConstantDeclValue()) // TODO(xsw): use int64 for 64-bit enums?
		at := defs.NewPos()
		if doc := ctx.docCommentGroup(item); doc != nil {
			at.Doc = doc
		}
		pos := goNodePos(ctx, item)
		defs.NewAt(at, func(cb *gogen.CodeBuilder) int {
			cb.Val(val)
			return 1
		}, 0, pos, enumType, name)
		return clang.Continue
	})
}

// -----------------------------------------------------------------------------
