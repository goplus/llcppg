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
	"os"
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

	typePrefix []string
	fnPrefix   []string
	rename     map[string]string
	classes    map[string]none // typedef names to be treated as classes

	nameLookup func(manglingName string) (archivePath string, ok bool)
	pubLookup  func(pkgPath string) (pubFile string, ok bool)

	pkgOf func(headerFile string) (pkgPath string, ok bool)

	fileBases map[clang.File]int // clang.File => base

	macroVals map[string]any             // macroName => value
	funcs     map[string]*funcObj        // manglingName => func object
	types     map[string]*types.TypeName // c/c++ fullName => type name object
	lastSeen  map[string]none            // last seen include file set (loaded include files)

	compiles []compileFunc
	pubs     []Entry

	// anonUnionSeq is the per-package counter that names tagless inline unions
	// hoisted to a Go type "_llcppg_union_<n>", starting at 0 in source order.
	// See issue goplus/llcppg#775 (the union proposal's D-series).
	anonUnionSeq  int
	anonStructSeq int

	unsafeImported bool
	keepDoc        bool
}

const (
	anonUnionPrefix  = "_llcppg_union_"
	anonStructPrefix = "_llcppg_struct_"
)

// nextAnonUnionName returns the next hoisted-union type name
// "_llcppg_union_<n>", incrementing the per-package counter. A tagless inline
// union (with or without a field name) is hoisted to such a type; a tagged or
// named union consumes no number. See issue goplus/llcppg#775.
func (p *pkgCtx) nextAnonUnionName() string {
	name := anonUnionPrefix + strconv.Itoa(p.anonUnionSeq)
	p.anonUnionSeq++
	return name
}

func (p *pkgCtx) nextAnonStructName() string {
	name := anonStructPrefix + strconv.Itoa(p.anonStructSeq)
	p.anonStructSeq++
	return name
}

func (p *pkgCtx) forceImportUnsafe() {
	if !p.unsafeImported {
		p.unsafeImported = true
		p.pkg.ForceImport("unsafe")
	}
}

func (p *pkgCtx) importPkg(pkgPath string) {
	pubFile, ok := p.pubLookup(pkgPath)
	if !ok {
		log.Panicln("[ERROR] pubFile not found for", pkgPath)
	}
	if debugCompileDecl {
		log.Println("==> importPkg", pkgPath)
	}
	pkg := p.pkg.Import(pkgPath)
	entries, err := loadPubFile(pubFile)
	if err != nil {
		if os.IsNotExist(err) {
			return // ignore missing pub file
		}
		log.Panicln("[ERROR] loadPubFile failed:", err)
	}
	scope := pkg.Types.Scope()
	for e := range entries {
		switch e.Kind {
		case 'T': // type
			if e.GoName == "" {
				e.GoName = p.typeName(strings.ReplaceAll(e.Name, "::", "_"), true)
			}
			if o := scope.Lookup(e.GoName); o != nil {
				if t, ok := o.(*types.TypeName); ok {
					p.types[e.Name] = t
				}
			}
		default:
			panic("todo: importPubFile " + e.Name)
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
	if o, ok := p.types[cName]; ok {
		return o, true
	}
	return nil, false
}

func (p *pkgCtx) typeOf(cName string) (types.Type, bool) {
	if o, ok := p.types[cName]; ok {
		return o.Type(), true
	}
	return nil, false
}

func (p *pkgCtx) typeName(name string, _ bool) string {
	return p.cstyleToGo(rmPrefix(name, p.typePrefix), true)
}

func (p *pkgCtx) funcName(name string, order int, typNamed *types.Named, global, _ bool) string {
	if global {
		name = rmPrefix(name, p.fnPrefix)
	} else {
		// don't remove type name suffix for a method
		typNamed = nil
	}
	if !strings.HasPrefix(name, "XGo_") { // avoid rewriting XGo_xxx names
		name = p.cstyleToGo(name, true)
		if typNamed != nil {
			objName := typNamed.Obj().Name()
			name = cutMethodPrefix(strings.TrimSuffix(name, objName), objName)
		}
	}
	if order >= 0 {
		name = name + "__" + strconv.FormatInt(int64(order), 36)
	}
	return name
}

func (p *pkgCtx) fieldName(name string, public bool) string {
	return p.cstyleToGo(name, public)
}

func (p *pkgCtx) varName(name string) string {
	return p.cstyleToGo(name, true)
}

func (p *pkgCtx) macroName(name string) string {
	return p.cstyleToGo(name, true)
}

func (p *pkgCtx) enumvalName(name string) string {
	return p.cstyleToGo(name, true)
}

func (p *pkgCtx) cstyleToGo(cName string, public bool) string {
	rename := p.rename
	parts := strings.Split(cName, "_")
	for i := 0; i < len(parts); i++ {
		part := parts[i]
		if part == "" {
			if i == 0 && public {
				parts[i] = "X_"
				i++ // skip next part
			} else {
				parts[i] = "_"
			}
		} else if v, ok := rename[part]; ok {
			parts[i] = v
		} else if i > 0 || public {
			parts[i] = cPubName(part)
		}
	}
	return strings.Join(parts, "")
}

func cutMethodPrefix(name, objName string) string {
	after, ok := strings.CutPrefix(name, "Get")
	if ok {
		after, ok = strings.CutPrefix(after, objName)
	} else {
		after, ok = strings.CutPrefix(name, objName)
	}
	if ok && after != "" {
		if c := after[0]; 'A' <= c && c <= 'Z' {
			return after
		}
	}
	return name
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

func rmPrefix(name string, prefix []string) string {
	for _, pfx := range prefix {
		if strings.HasPrefix(name, pfx) {
			return name[len(pfx):]
		}
	}
	return name
}

// -----------------------------------------------------------------------------

type overloads struct {
	items []*funcObj
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

type funcObj struct {
	name      string
	decl      clang.Cursor
	overloads *overloads

	manglingName string
}

// order returns the order of the object in the overloads list.
// -1 means no order (only one overload, or not found).
func (p *funcObj) order() int {
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

func (p *scopeCtx) addFunc(ctx *pkgCtx, name string, decl clang.Cursor) (*funcObj, bool) {
	manglingName := clang.Mangling(decl)
	if fn, ok := ctx.funcs[manglingName]; ok { // re-declared
		if decl.IsFunctionInlined() != 0 {
			fn.decl = decl // use the latest inline decl
		}
		return fn, false
	}

	obj := &funcObj{
		name:         name,
		decl:         decl,
		manglingName: manglingName,
	}
	ovs, ok := p.overloads[name]
	if ok {
		ovs.items = append(ovs.items, obj)
	} else {
		ovs = &overloads{items: []*funcObj{obj}}
		p.overloads[name] = ovs
	}
	obj.overloads = ovs
	ctx.funcs[manglingName] = obj
	return obj, true
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
