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
	"go/ast"
	"go/token"
	"go/types"
	"log"
)

func loadGoFiles(ctx *pkgCtx, goFiles map[string]*ast.File) {
}

func newGoType(pkg *types.Package, pos token.Pos, name string) *types.Named {
	typName := types.NewTypeName(pos, pkg, name, nil)
	if old := pkg.Scope().Insert(typName); old != nil {
		log.Panicf("%s redeclared in this block\n\tprevious declaration at %v\n", name, "<TODO>")
	}
	return types.NewNamed(typName, nil, nil)
}
