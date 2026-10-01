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
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"log"
	"sort"
	"strings"

	"github.com/goplus/gogen"
	"github.com/goplus/lib/c"
	"github.com/goplus/llcppg/clang"
	lc "github.com/llarhub/clang-c"
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
	rg.Start().Spelling(&file, nil, nil, &pos)
	if file == clang.InvalidFile {
		return nil
	}
	rg.End().Spelling(nil, nil, nil, &end)
	base := ctx.getFileBase(v, file)
	return &node{pos: token.Pos(int(pos) + base), end: token.Pos(int(end) + base), ctx: ctx}
}

func goNodePos(ctx *pkgCtx, v clang.Cursor) token.Pos {
	var file clang.File
	var pos c.Uint
	v.Extent().Start().Spelling(&file, nil, nil, &pos)
	if file == clang.InvalidFile {
		return token.NoPos
	}
	base := ctx.getFileBase(v, file)
	return token.Pos(int(pos) + base)
}

// -----------------------------------------------------------------------------

type nodeInterp struct {
	fset *token.FileSet
}

func (p *nodeInterp) LoadExpr(v ast.Node) string {
	panic("todo: nodeInterp.LoadExpr")
}

// -----------------------------------------------------------------------------

type none struct{}

type compileFunc = func(ctx *pkgCtx)
type compileUnit struct {
	fn compileFunc
	at *gogen.File
}

type typeObj struct {
	*types.TypeName
	feats int
}

type typDecl struct {
	*gogen.TypeDecl
	defs *gogen.TypeDefs
}

type pkgCtx struct {
	scopeCtx
	pkg  *gogen.Package
	cb   *gogen.CodeBuilder
	llgo *gogen.ConstDefs
	wrap *WrapFile
	fset *token.FileSet
	c    gogen.PkgRef
	ctyp [cBasicMax]types.Type
	tany types.Type
	tptr types.Type

	cflags string
	lang   Language

	wrapFileHeader string

	typeAbbrSuffix []string

	fnPrefix    []string
	varPrefix   []string
	enumPrefix  []string
	macroPrefix []string
	typePrefix  []string
	typeSuffix  []string
	classes     []string          // typedef names to be treated as classes
	nonClasses  []string          // typedef names to be treated as non-classes
	typeAbbr    map[string]any    // Go type name => abbreviated name(s), used in function names
	rename      map[string]string // C/C++ name => Go name
	typeIgnores []string          // C/C++ type names to be ignored
	macroIgnore []string          // C/C++ macro names to be ignored
	nsIgnores   []string          // Go style namespace names to be ignored

	nameLookup func(manglingName string) (archivePath string, ok bool)
	pubLookup  func(pkgPath string) (pubFile string, ok bool)

	pkgOf func(headerFile string) (pkgPath string, ok bool)

	fileBases map[clang.File]int // clang.File => base

	macroVals map[string]any      // macroName => value
	fns       map[string]*funcObj // fnUSR => func object
	types     map[string]typeObj  // c/c++ fullName => type name object (include external types)
	typdecls  map[string]typDecl  // c/c++ fullName => type declaration object (only local types)
	impPkgs   map[string]none     // imported package path set
	lastSeen  map[string]none     // last seen include file set (loaded include files)
	thisSeen  map[string]none     // include file set seen in this translation unit

	compiles []compileUnit
	pubs     []Entry

	anonSeq int

	stdRecvName bool
	keepDoc     bool
}

func (p *pkgCtx) ignoref(feats int, decl clang.Cursor, format string, args ...any) {
	if feats&featQuietIgnore == 0 || debugQuietIgnore {
		p.logf(decl, format, args...)
	}
}

func (p *pkgCtx) logf(decl clang.Cursor, format string, args ...any) {
	pos := p.fset.Position(goNodePos(p, decl))
	log.Printf("%s: %s", pos, fmt.Sprintf(format, args...))
}

func (p *pkgCtx) panicf(decl clang.Cursor, format string, args ...any) {
	pos := p.fset.Position(goNodePos(p, decl))
	log.Panicf("%s: %s", pos, fmt.Sprintf(format, args...))
}

func (p *pkgCtx) logtf(typ lc.Type, format string, args ...any) {
	p.logf(typ.Declaration(), format, args...)
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

func (p *pkgCtx) addCompileUnit(f compileFunc) {
	p.compiles = append(p.compiles, compileUnit{
		fn: f,
		at: p.pkg.CurFile(),
	})
}

func (p *pkgCtx) compile() {
	pkg := p.pkg
	for {
		compiles := p.compiles
		if len(compiles) == 0 {
			break
		}
		p.compiles = nil
		for _, c := range compiles {
			pkg.RestoreCurFile(c.at)
			c.fn(p)
		}
	}
}

func (p *pkgCtx) unsafePointer() types.Type {
	if p.tptr == nil {
		p.tptr = p.pkg.Import("unsafe").Ref("Pointer").Type()
	}
	return p.tptr
}

func (p *pkgCtx) any() types.Type {
	if p.tany == nil {
		p.tany = types.Universe.Lookup("any").Type()
	}
	return p.tany
}

func (p *pkgCtx) basicTyp(kind basicKind) types.Type {
	typ := p.ctyp[kind]
	if typ == nil {
		typ = p.c.Ref(ctypBasic[kind]).Type()
	}
	return typ
}

func (p *pkgCtx) aliasTypeName(cName, goName string) {
	cName = trimTypeTag(cName)
	if _, ok := p.rename[cName]; !ok {
		// insert alias name if not already present
		p.rename[cName] = goName
	}
}

func (p *pkgCtx) addType(kind typeTag, decl clang.Cursor, typNamed *types.Named) {
	var cName string
	if decl.Kind == lc.Cursor_ClassTemplate { // TODO(xsw): check if this is correct
		cName = clang.String(decl)
	} else {
		cName = clang.String(decl.Type())
	}
	typObj := typeObj{typNamed.Obj(), 0}
	p.types[cName] = typObj
	p.aliasTypeName(cName, typObj.Name())

	if debugCompileDecl {
		log.Println("==> addType", cName, typObj.Name())
	}

	// name of typedef <tag> may be "m" instead of "<tag> m"
	tag := tagStrvals[kind]
	if strings.HasPrefix(cName, tag) {
		p.types[cName[len(tag):]] = typObj
	} else {
		p.types[tag+cName] = typObj
	}
}

func (p *pkgCtx) ignoreType(cName string, feats int) {
	if _, ok := p.types[cName]; !ok {
		p.types[cName] = typeObj{nil, feats}
	}
}

func (p *pkgCtx) isTypeIgnored(cName string) bool {
	return contains(trimTypeTag(cName), p.typeIgnores)
}

func (p *pkgCtx) typeObj(cName string) (*types.TypeName, bool) {
	if o, ok := p.types[cName]; ok {
		return o.TypeName, true
	}
	return nil, false
}

func (p *pkgCtx) typeOf(cName string) (types.Type, bool) {
	if o, ok := p.types[cName]; ok {
		return o.Type(), true
	}
	return nil, false
}

func (p *pkgCtx) isMacroIgnored(cName string) bool {
	return contains(cName, p.macroIgnore)
}

func (p *pkgCtx) isNSIgnored(cName string) bool {
	return contains(cName, p.nsIgnores)
}

// -----------------------------------------------------------------------------

type templateClass struct {
	name      string       // go name
	decl      clang.Cursor // AST object
	overloads *overloads
}

// order returns the order of the template class in the overloads list.
// -1 means no order (only one overload, or not found).
func (p *templateClass) order() int {
	items := p.overloads.classes
	if len(items) > 1 {
		for i, obj := range items {
			if obj == p {
				return i
			}
		}
	}
	return -1
}

type funcObj struct {
	name      string       // go name
	decl      clang.Cursor // AST object
	overloads *overloads

	isOperator bool
}

// order returns the order of the object in the overloads list.
// -1 means no order (only one overload, or not found).
func (p *funcObj) order() int {
	items := p.overloads.fns
	if len(items) > 1 {
		for i, obj := range items {
			if obj == p {
				return i
			}
		}
	}
	return -1
}

type overloads struct {
	fns     []*funcObj
	classes []*templateClass
}

func (p *overloads) reorder() {
	items := p.fns
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

func cloneTypeName(obj *types.TypeName) *types.TypeName {
	return types.NewTypeName(obj.Pos(), obj.Pkg(), obj.Name(), obj.Type())
}

func cloneTypeParams(tparams []*types.TypeParam) []*types.TypeParam {
	if len(tparams) == 0 {
		return nil
	}
	ret := make([]*types.TypeParam, len(tparams))
	for i, tp := range tparams {
		obj := cloneTypeName(tp.Obj())
		ret[i] = types.NewTypeParam(obj, tp.Constraint())
	}
	return ret
}

func concatTypeParams(a, b []*types.TypeParam) []*types.TypeParam {
	if len(a) == 0 {
		return b
	}
	a = cloneTypeParams(a)
	if len(b) == 0 {
		return a
	}
	ret := make([]*types.TypeParam, 0, len(a)+len(b))
	ret = append(ret, a...)
	return append(ret, b...)
}

type scopeCtx struct {
	overloads map[string]*overloads // name => overload items
	tparams   []*types.TypeParam    // type parameters, only for class/struct scope
	parent    *scopeCtx
}

func (p *scopeCtx) typeParams(in []*types.TypeParam) []*types.TypeParam {
	for p != nil {
		in = concatTypeParams(p.tparams, in)
		p = p.parent
	}
	return in
}

func (p *scopeCtx) lookupType(name string) (types.Type, bool) {
	for p != nil {
		for _, t := range p.tparams {
			if t.Obj().Name() == name {
				return t, true
			}
		}
		p = p.parent
	}
	return nil, false
}

func (p *scopeCtx) addTemplateClass(_ *pkgCtx, name string, decl clang.Cursor) (*templateClass, bool) {
	if debugCompileDecl {
		log.Println("==> addTemplateClass", name, "-", clang.DisplayName(decl))
	}
	// TODO(xsw): check if the class is already added
	obj := &templateClass{
		name: name,
		decl: decl,
	}
	ovs, ok := p.overloads[name]
	if ok {
		ovs.classes = append(ovs.classes, obj)
	} else {
		ovs = &overloads{classes: []*templateClass{obj}}
		p.overloads[name] = ovs
	}
	obj.overloads = ovs
	return obj, true
}

func (p *scopeCtx) addFunc(ctx *pkgCtx, name string, decl clang.Cursor, isOp bool) (*funcObj, bool) {
	fnUSR := funcUSR(decl)
	if debugCompileDecl {
		log.Println("==> addFunc", funcDisplayName(decl), "- USR:", fnUSR)
	}

	if fn, ok := ctx.fns[fnUSR]; ok { // re-declared
		if decl.IsFunctionInlined() != 0 {
			fn.decl = decl // use the latest inline decl
		}
		return fn, false
	}

	obj := &funcObj{
		name:       name,
		decl:       decl,
		isOperator: isOp,
	}
	ovs, ok := p.overloads[name]
	if ok {
		ovs.fns = append(ovs.fns, obj)
	} else {
		ovs = &overloads{fns: []*funcObj{obj}}
		p.overloads[name] = ovs
	}
	obj.overloads = ovs
	ctx.fns[fnUSR] = obj

	return obj, true
}

func (p *scopeCtx) reorder() {
	for _, o := range p.overloads {
		o.reorder()
	}
}

// -----------------------------------------------------------------------------
