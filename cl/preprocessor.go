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
)

// -----------------------------------------------------------------------------

func loadMacro(ctx *pkgCtx, decl clang.Cursor) {
	if decl.IsMacroFunctionLike() != 0 {
		return
	}
	origName := clang.String(decl)
	tu := clang.TU(decl)
	tokens, dispose := tu.Tokenize(decl.Extent())
	defer dispose()
	if debugCompileDecl {
		log.Println("macro", origName, "-", len(tokens), "tokens")
	}
	if len(tokens) > 1 {
		if v, ok := evalConstExpr(ctx, tu, tokens[1:]); ok {
			pkg := ctx.pkg
			pkgTypes := pkg.Types
			ctx.macroVals[origName] = v
			name, _ := ctx.getPubName(origName, -1)
			pkg.NewConstDefs(pkgTypes.Scope()).New(func(cb *gogen.CodeBuilder) int {
				cb.Val(v)
				return 1
			}, 0, token.NoPos, nil, name)
		}
	}
}

// -----------------------------------------------------------------------------
