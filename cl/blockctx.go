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

	"github.com/goplus/gogen"
	"github.com/goplus/lib/c"
	"github.com/goplus/llcppg/clang"
)

// -----------------------------------------------------------------------------

type node struct {
	pos token.Pos
	end token.Pos
	ctx *blockCtx
}

func (p *node) Pos() token.Pos {
	return p.pos
}

func (p *node) End() token.Pos {
	return p.end
}

/* TODO(xsw):
func goNode(ctx *blockCtx, v clang.Cursor) ast.Node {
	var pos, end c.Uint
	rg := v.Extent()
	rg.RangeStart().SpellingLocation(nil, nil, nil, &pos)
	rg.RangeEnd().SpellingLocation(nil, nil, nil, &end)
	base := ctx.file.Base()
	return &node{pos: token.Pos(int(pos) + base), end: token.Pos(int(end) + base), ctx: ctx}
}
*/

func goNodePos(ctx *blockCtx, v clang.Cursor) token.Pos {
	var pos c.Uint
	v.Extent().RangeStart().SpellingLocation(nil, nil, nil, &pos)
	return token.Pos(int(pos) + ctx.file.Base())
}

// -----------------------------------------------------------------------------

type nodeInterp struct {
	fset *token.FileSet
}

func (p *nodeInterp) Position(start token.Pos) token.Position {
	return p.fset.Position(start)
}

func (p *nodeInterp) LoadExpr(v ast.Node) string {
	panic("todo: nodeInterp.LoadExpr")
}

// -----------------------------------------------------------------------------

type blockCtx struct {
	pkg  *gogen.Package
	cb   *gogen.CodeBuilder
	fset *token.FileSet
	file *token.File
	c    gogen.PkgRef

	unsafeImported bool
}

func (p *blockCtx) forceImportUnsafe() {
	if !p.unsafeImported {
		p.unsafeImported = true
		p.pkg.ForceImport("unsafe")
	}
}

func (p *blockCtx) initFile(file Source) {
	src := file.TU.FileContents(file.Handle)
	p.file = p.fset.AddFile("", -1, len(src))
	p.file.SetLinesForContent(src)
}

func (p *blockCtx) getPubName(pfnName *string) (rewritten bool) {
	fnName := *pfnName
	pubName := cPubName(fnName)
	rewritten = fnName != pubName
	if rewritten {
		*pfnName = pubName
	}
	return
}

func cPubName(name string) string {
	if r := name[0]; 'a' <= r && r <= 'z' {
		r -= 'a' - 'A'
		return string(r) + name[1:]
	} else if r == '_' {
		return "X" + name
	}
	return name
}

// -----------------------------------------------------------------------------
