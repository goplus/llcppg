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
	"sort"
	"strconv"

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
	wrap *WrapFile
	fset *token.FileSet
	tu   clang.TranslationUnit
	c    gogen.PkgRef

	cflags string
	lang   Language

	wrapFileHeader string

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

func (p *pkgCtx) getPubName(fnName string, order int) (pubName string, rewritten bool) {
	pubName = cPubName(fnName)
	if order >= 0 {
		return pubName + "__" + strconv.FormatInt(int64(order), 36), true
	}
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

func (p *overloads) reorder() {
	items := p.items
	if len(items) > 1 {
		sort.SliceStable(items, func(i, j int) bool {
			a, b := items[i].decl, items[j].decl
			na, nb := a.NumArguments(), b.NumArguments()
			if na != nb {
				return na < nb
			}
			for k := range c.Uint(na) {
				ta, tb := a.Argument(k).Type(), b.Argument(k).Type()
				if ret := cmpType(ta, tb); ret != 0 {
					return ret < 0
				}
			}
			return false
		})
	}
}

type object struct {
	name      string
	decl      clang.Cursor
	overloads *overloads
}

// order returns the order of the object in the overloads list.
// -1 means no order (only one overload, or not found).
func (p *object) order() int {
	items := p.overloads.items
	if len(items) > 1 {
		for i, obj := range items {
			if obj == p {
				return i
			}
		}
	}
	return -1
}

type scopeCtx struct {
	overloads map[string]*overloads // name => overload items
}

func (p *scopeCtx) addObject(name string, decl clang.Cursor) *object {
	obj := &object{
		name: name,
		decl: decl,
	}
	ovs, ok := p.overloads[name]
	if ok {
		ovs.items = append(ovs.items, obj)
	} else {
		ovs = &overloads{items: []*object{obj}}
		p.overloads[name] = ovs
	}
	obj.overloads = ovs
	return obj
}

func (p *scopeCtx) reorder() {
	for _, o := range p.overloads {
		o.reorder()
	}
}

// -----------------------------------------------------------------------------
