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
	"strconv"
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

	cflags string
	lang   Language

	wrapFileHeader string

	fnPrefix    []string
	varPrefix   []string
	enumPrefix  []string
	macroPrefix []string
	typePrefix  []string
	classes     []string          // typedef names to be treated as classes
	nonClasses  []string          // typedef names to be treated as non-classes
	typeAbbr    map[string]string // Go type name => abbreviated name, used in function names
	rename      map[string]string // C/C++ name => Go name
	typeIgnores []string          // C/C++ type names to be ignored
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

const (
	anonPrefix = "_llcppg_anon_"
)

func (p *pkgCtx) nextAnonName() string {
	name := anonPrefix + strconv.Itoa(p.anonSeq)
	p.anonSeq++
	return name
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

func (p *pkgCtx) isTypeIgnored(cName string) bool {
	return contains(cName, p.typeIgnores)
}

func (p *pkgCtx) isNSIgnored(ns string) bool {
	return contains(ns, p.nsIgnores)
}

func (p *pkgCtx) globalName(name string, trimPrefixs []string) string {
	if v, ok := p.rename[name]; ok {
		return v // special case
	}
	return p.cstyleToGo(rmPrefix(name, trimPrefixs), true)
}

func (p *pkgCtx) fieldName(name string, public bool) string {
	return p.cstyleToGo(name, public)
}

func (p *pkgCtx) varName(name string) string {
	return p.globalName(name, p.varPrefix)
}

func (p *pkgCtx) macroName(name string) string {
	return p.globalName(name, p.macroPrefix)
}

func (p *pkgCtx) enumvalName(name, ns string) string {
	if strings.HasSuffix(ns, "_") {
		return p.globalName(ns, p.enumPrefix) + p.cstyleToGo(name, true)
	}
	return p.globalName(nameWithNS(name, ns), p.enumPrefix)
}

func (p *pkgCtx) typeName(name string, _ bool) string {
	return p.globalName(name, p.typePrefix)
}

func (p *pkgCtx) funcName(name string, order int, typName, typCName string, global, _ bool) string {
	if v, ok := p.rename[name]; ok {
		return v // special case
	}
	if global {
		name = rmPrefix(name, p.fnPrefix)
		if typCName != "" {
			// remove typCName prefix & suffix
			if before, ok := strings.CutSuffix(name, typCName); ok {
				name = strings.TrimSuffix(before, "_")
			} else if after, ok := strings.CutPrefix(name, typCName+"_"); ok {
				name = after
			} else {
				name = strings.TrimPrefix(name, typName+"_")
			}
		}
	} else {
		// don't remove type name suffix for a method
		typName = ""
	}
	if !strings.HasPrefix(name, "XGo_") { // avoid rewriting XGo_xxx names
		name = p.cstyleToGo(name, true)
		if typName != "" {
			if v, ok := p.typeAbbr[typName]; ok { // Go type name => abbreviated name
				typName = v
			}
			name = cutMethodPrefix(strings.TrimSuffix(name, typName), typName)
		}
	}
	if order >= 0 {
		name = name + "__" + strconv.FormatInt(int64(order), 36)
	}
	return name
}

func (p *pkgCtx) cstyleToGo(cName string, public bool) string {
	rename := p.rename
	parts := strings.Split(cName, "_")
	if isAllUpperStart(parts) {
		if parts[0] == "" && public {
			return "X" + cName
		}
		return cName
	}
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

func (p *pkgCtx) nsName(ns, inner string) string {
	if v, ok := p.rename[inner]; ok {
		inner = v // special case
	}
	return nsName(ns, inner)
}

func nsName(ns, inner string) string {
	if ns != "" {
		c := ns[len(ns)-1]
		if 'A' <= c && c <= 'Z' {
			ns += "_"
		}
	}
	if c := inner[0]; 'a' <= c && c <= 'z' {
		return ns + string(c-'a'+'A') + inner[1:]
	}
	return ns + inner
}

func nameWithNS(name, ns string) string {
	if ns == "" {
		return name
	}
	return nsName(ns, name)
}

func isAllUpperStart(parts []string) bool {
	for _, part := range parts {
		if part != "" {
			if r := part[0]; 'a' <= r && r <= 'z' {
				return false
			}
		}
	}
	return true
}

func cutMethodPrefix(name, objName string) string {
	name = cutPrefix(name, objName)
	name = cutPrefix(name, "Get")
	return cutPrefix(name, objName)
}

func cutPrefix(name, prefix string) string {
	after, ok := strings.CutPrefix(name, prefix)
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

func contains(v string, names []string) bool {
	for _, name := range names {
		if name == v {
			return true
		}
	}
	return false
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

type scopeCtx struct {
	overloads map[string]*overloads // name => overload items
	tparams   []*types.TypeParam    // type parameters, only for class/struct scope
	parent    *scopeCtx
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

func funcUSR(decl clang.Cursor) string {
	return clang.USR(decl)
}

func funcDisplayName(decl clang.Cursor) string {
	return clang.DisplayName(decl)
}

// -----------------------------------------------------------------------------
