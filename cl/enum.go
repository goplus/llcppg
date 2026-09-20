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
	"log"

	"github.com/goplus/gogen"
	"github.com/goplus/llcppg/clang"
	lc "github.com/goplus/llcppg/lib/clang"
)

// -----------------------------------------------------------------------------

// compileEnum translates a C/C++ enum declaration into Go const declarations.
//
// Only the enum constants are emitted (as Go consts); the enum type itself is
// not generated. Whether the enum is global, in a namespace, or nested inside a
// class only affects the naming of each constant: ns carries the enclosing
// namespace/class prefix (e.g. "bar_" or "Color_"), so a constant Red becomes
// bar_Red / Color_Red and then goes through getPubName for the final Go name.
func compileEnum(ctx *pkgCtx, decl clang.Cursor, ns string) {
	if debugCompileDecl {
		log.Println("enum", ns+clang.String(decl))
	}
	pkg := ctx.pkg
	scope := pkg.Types.Scope()
	defs := pkg.NewConstDefs(scope)
	clang.VisitChildren(decl, func(item, parent clang.Cursor) clang.ChildVisitResult {
		if item.Kind != lc.CursorEnumConstantDecl {
			return clang.Continue
		}
		origName := ns + clang.String(item)
		name, _ := ctx.getPubName(origName, -1)
		val := int(item.EnumConstantDeclValue())
		defs.New(func(cb *gogen.CodeBuilder) int {
			cb.Val(val)
			return 1
		}, 0, token.NoPos, nil, name)
		return clang.Continue
	})
}

// -----------------------------------------------------------------------------
