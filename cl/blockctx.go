/*
 * Copyright (c) 2025 The XGo Authors (xgo.dev). All rights reserved.
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

	"github.com/goplus/gogen"
	"github.com/goplus/llcppg/ast"
)

type blockCtx struct {
	pkg  *gogen.Package
	cb   *gogen.CodeBuilder
	fset *token.FileSet
}

func (ctx *blockCtx) inPkg(loc *ast.Location) bool {
	panic("todo")
}

/*
func (ctx *blockCtx) goNodePos(v *ast.Node) token.Pos {
	if rg := v.Range; rg != nil && ctx.file != nil {
		base := ctx.file.Base()
		return token.Pos(int(rg.Begin.Offset) + base)
	}
	return token.NoPos
}
*/
