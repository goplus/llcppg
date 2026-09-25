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
func loadEnum(ctx *pkgCtx, decl clang.Cursor, ns string) {
	name := clang.String(decl)
	origName := nameWithNS(name, ns)
	anonymous := decl.IsAnonymous() != 0
	if debugCompileDecl {
		log.Println("enum", origName, "anonymous:", anonymous)
	}

	pkg := ctx.pkg
	pkgTypes := pkg.Types

	// Preserve the enum type name for a named enum. An anonymous enum has no
	// name, so its constants stay untyped.
	var enumType types.Type
	if !anonymous {
		typDefs := pkg.NewTypeDefs()
		if doc := ctx.docCommentGroup(decl); doc != nil {
			typDefs.SetComments(doc)
		}
		// C enums decay to int; use the same C int type the rest of the
		// generator uses so enum-typed values interoperate with C APIs.
		typeName := ctx.typeName(origName, true)
		typDecl := typDefs.NewType(typeName, goNode(ctx, decl))
		typNamed := typDecl.InitType(pkg, ctx.basicTyp(cInt))
		ctx.addType(tagEnum, decl, typNamed)
		enumType = typNamed
	}

	defs := pkg.NewConstDefs(pkgTypes.Scope())
	// For an anonymous enum there is no type to carry the doc, so attach the
	// enum's doc comment to the generated const block instead.
	if anonymous {
		if doc := ctx.docCommentGroup(decl); doc != nil {
			defs.SetComments(doc)
		}
	}
	clang.VisitChildren(decl, func(item, parent clang.Cursor) clang.ChildVisitResult {
		if item.Kind != lc.CursorEnumConstantDecl {
			return clang.Continue
		}
		origName := nameWithNS(clang.String(item), ns)
		name := ctx.enumvalName(origName)
		val := int(item.EnumConstantDeclValue()) // TODO(xsw): use int64 for 64-bit enums?
		at := defs.NewPos()
		if doc := ctx.docCommentGroup(item); doc != nil {
			at.Doc = doc
		}
		defs.NewAt(at, func(cb *gogen.CodeBuilder) int {
			cb.Val(val)
			return 1
		}, 0, token.NoPos, enumType, name)
		return clang.Continue
	})
}

// -----------------------------------------------------------------------------
