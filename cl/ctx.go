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
	"log"

	"github.com/goplus/gogen"
	"github.com/goplus/lib/c"
	"github.com/goplus/llcppg/clang"
)

// -----------------------------------------------------------------------------

type node struct {
	pos token.Pos
	end token.Pos
	ctx *pkgCtx
}

func (p *node) Pos() token.Pos {
	return p.pos
}

func (p *node) End() token.Pos {
	return p.end
}

func goNode(ctx *pkgCtx, v clang.Cursor) ast.Node {
	var file clang.File
	var pos, end c.Uint
	rg := v.Extent()
	rg.RangeStart().SpellingLocation(&file, nil, nil, &pos)
	base, ok := ctx.fileBases[file]
	if !ok {
		return nil
	}
	rg.RangeEnd().SpellingLocation(nil, nil, nil, &end)
	return &node{pos: token.Pos(int(pos) + base), end: token.Pos(int(end) + base), ctx: ctx}
}

func goNodePos(ctx *pkgCtx, v clang.Cursor) token.Pos {
	var file clang.File
	var pos c.Uint
	v.Extent().RangeStart().SpellingLocation(&file, nil, nil, &pos)
	base, ok := ctx.fileBases[file]
	if !ok {
		return token.NoPos
	}
	return token.Pos(int(pos) + base)
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

type compileFunc = func(ctx *pkgCtx)

type pkgCtx struct {
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

	methods  map[string]*classMethod // manglingName => class
	compiles []compileFunc

	unsafeImported bool
}

func (p *pkgCtx) forceImportUnsafe() {
	if !p.unsafeImported {
		p.unsafeImported = true
		p.pkg.ForceImport("unsafe")
	}
}

func (p *pkgCtx) initFiles(files []string) {
	fset := p.fset
	tu := p.tu
	fileBases := make(map[clang.File]int)
	for _, filename := range files {
		f := tu.File(filename)
		if f == clang.InvalidFile {
			log.Println("[WARN]", filename, "is not included in the translation unit")
			continue
		}
		src := p.tu.FileContents(f)
		tf := fset.AddFile(filename, -1, len(src))
		tf.SetLinesForContent(src)
		fileBases[f] = tf.Base()
	}
	p.fileBases = fileBases
}

func (p *pkgCtx) compile() {
	for _, compile := range p.compiles {
		compile(p)
	}
}

func (p *pkgCtx) getPubName(fnName string) (pubName string, rewritten bool) {
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

type overloads struct {
	items []*object
}

type object struct {
	name      string
	decl      clang.Cursor
	overloads *overloads
	idx       int // index in overloads.items
}

type scopeCtx struct {
	ns        string
	overloads map[string]*overloads // name => overload items
}

func (p *scopeCtx) addObject(decl clang.Cursor) *object {
	name := clang.String(decl)
	obj := &object{
		name: name,
		decl: decl,
	}
	ovs, ok := p.overloads[name]
	if ok {
		obj.idx = len(ovs.items)
		ovs.items = append(ovs.items, obj)
	} else {
		ovs = &overloads{items: []*object{obj}}
		p.overloads[name] = ovs
	}
	obj.overloads = ovs
	return obj
}

// -----------------------------------------------------------------------------
