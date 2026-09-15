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

func goNode(ctx *blockCtx, v clang.Cursor) ast.Node {
	var file clang.File
	var pos, end c.Uint
	rg := v.Extent()
	rg.RangeStart().SpellingLocation(&file, nil, nil, &pos)
	rg.RangeEnd().SpellingLocation(nil, nil, nil, &end)
	base := ctx.fileBases[file]
	return &node{pos: token.Pos(int(pos) + base), end: token.Pos(int(end) + base), ctx: ctx}
}

func goNodePos(ctx *blockCtx, v clang.Cursor) token.Pos {
	var file clang.File
	var pos c.Uint
	v.Extent().RangeStart().SpellingLocation(&file, nil, nil, &pos)
	return token.Pos(int(pos) + ctx.fileBases[file])
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
	llgo *gogen.ConstDefs
	wrap *wrapFile
	fset *token.FileSet
	tu   clang.TranslationUnit
	c    gogen.PkgRef

	cflags string
	lang   Language

	nameLookup func(manglingName string) (archivePath string, ok bool)

	fileBases map[clang.File]int // clang.File => base

	methods map[string]*classMethod // manglingName => class
	clTasks []func()

	unsafeImported bool
}

func (p *blockCtx) forceImportUnsafe() {
	if !p.unsafeImported {
		p.unsafeImported = true
		p.pkg.ForceImport("unsafe")
	}
}

func (p *blockCtx) initFiles(files []string) {
	fset := p.fset
	tu := p.tu
	fileBases := make(map[clang.File]int)
	for _, filename := range files {
		f := tu.File(filename)
		src := p.tu.FileContents(f)
		tf := fset.AddFile(filename, -1, len(src))
		tf.SetLinesForContent(src)
		fileBases[f] = tf.Base()
	}
	p.fileBases = fileBases
}

func (p *blockCtx) getPubName(fnName string) (pubName string, rewritten bool) {
	pubName = cPubName(fnName)
	rewritten = fnName != pubName
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
