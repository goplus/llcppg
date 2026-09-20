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
	"sort"
	"strconv"
	"strings"

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
	rg.RangeEnd().SpellingLocation(nil, nil, nil, &end)
	base := ctx.getFileBase(v, file)
	return &node{pos: token.Pos(int(pos) + base), end: token.Pos(int(end) + base), ctx: ctx}
}

func goNodePos(ctx *pkgCtx, v clang.Cursor) token.Pos {
	var file clang.File
	var pos c.Uint
	v.Extent().RangeStart().SpellingLocation(&file, nil, nil, &pos)
	base := ctx.getFileBase(v, file)
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

type none struct{}
type compileFunc = func(ctx *pkgCtx)

type pkgCtx struct {
	scopeCtx
	pkg  *gogen.Package
	cb   *gogen.CodeBuilder
	llgo *gogen.ConstDefs
	wrap *WrapFile
	fset *token.FileSet
	c    gogen.PkgRef

	cflags string
	lang   Language

	wrapFileHeader string

	nameLookup func(manglingName string) (archivePath string, ok bool)
	pubLookup  func(pkgPath string) (pubFile string, ok bool)

	pkgOf func(headerFile string) (pkgPath string, ok bool)

	fileBases map[clang.File]int // clang.File => base

	macroVals map[string]any          // macroName => value
	methods   map[string]*classMethod // manglingName => class
	objects   map[string]types.Object // c/c++ fullName => object
	includes  map[string]none         // includeFile Set

	compiles []compileFunc
	pubs     []PublicEntry

	unsafeImported bool
}

func (p *pkgCtx) forceImportUnsafe() {
	if !p.unsafeImported {
		p.unsafeImported = true
		p.pkg.ForceImport("unsafe")
	}
}

func (p *pkgCtx) importPkg(pkgPath string) {
	if debugCompileDecl {
		log.Println("==> importPkg", pkgPath)
	}
	if pubFile, ok := p.pubLookup(pkgPath); ok {
		if debugCompileDecl {
			log.Println("==> pubFile", pubFile)
		}
		pkg := p.pkg.Import(pkgPath)
		scope := pkg.Types.Scope()
		if it, _, e := loadPubFile(pubFile); e == nil {
			for cName, goName := range it {
				if goName == "" {
					goName, _ = p.getPubName(strings.ReplaceAll(cName, "::", "_"), -1)
				}
				if o := scope.Lookup(goName); o != nil {
					p.objects[cName] = o
				}
			}
		}
	}
}

func (p *pkgCtx) getFileBase(c clang.Cursor, file clang.File) int {
	base, ok := p.fileBases[file]
	if !ok {
		src := clang.TU(c).FileContents(file)
		tf := p.fset.AddFile(clang.FileName(file), -1, len(src))
		tf.SetLinesForContent(src)
		base = tf.Base()
		p.fileBases[file] = base
	}
	return base
}

func (p *pkgCtx) compile() {
	for _, compile := range p.compiles {
		compile(p)
	}
}

func (p *pkgCtx) typeObj(cName string) (*types.TypeName, bool) {
	if o, ok := p.objects[cName]; ok {
		if t, ok := o.(*types.TypeName); ok {
			return t, true
		}
	}
	return nil, false
}

func (p *pkgCtx) typeOf(cName string) (types.Type, bool) {
	if o, ok := p.objects[cName]; ok {
		if t, ok := o.(*types.TypeName); ok {
			return t.Type(), true
		}
	}
	return nil, false
}

func (p *pkgCtx) getPubName(cName string, order int) (pubName string, rewritten bool) {
	pubName = cPubName(cName)
	if order >= 0 {
		return pubName + "__" + strconv.FormatInt(int64(order), 36), true
	}
	rewritten = cName != pubName
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

func substObj(pkg *types.Package, scope *types.Scope, origName string, real types.Object) {
	old := scope.Insert(gogen.NewSubst(token.NoPos, pkg, origName, real))
	if old != nil {
		if t, ok := old.Type().(*gogen.TySubst); ok {
			t.Real = real
		} else {
			log.Panicln(origName, "redefined")
		}
	}
}

// -----------------------------------------------------------------------------

func avoidKeyword(name *string) {
	switch *name {
	case "map", "type", "range", "chan", "var", "func", "go", "select",
		"defer", "package", "import", "interface", "fallthrough":
		*name += "_"
	}
}

// -----------------------------------------------------------------------------
