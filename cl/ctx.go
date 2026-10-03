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
	"os"
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

	nsPrefix    []string
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
	typeAlias   map[string]string // C/C++ type name => another C/C++ type name
	typeIgnores []string          // C/C++ type names to be ignored
	macroIgnore []string          // C/C++ macro names to be ignored
	nsIgnore    []string          // C/C++ namespace names to be ignored

	nameLookup func(manglingName string) (archivePath string, ok bool)
	pubLookup  func(pkgPath string) (pubFile string, ok bool)

	pkgOf func(headerFile string) (pkgPath string, ok bool)

	fileBases map[clang.File]int // clang.File => base

	macroVals map[string]any          // macroName => value
	ovobjs    map[string]*overloadObj // objUSR => overload object
	types     map[string]typeObj      // c/c++ fullName => type name object (include external types)
	typdecls  map[string]typDecl      // c/c++ fullName => type declaration object (only local types)
	impPkgs   map[string]none         // imported package path set
	lastSeen  map[string]none         // last seen include file set (loaded include files)
	thisSeen  map[string]none         // include file set seen in this translation unit

	loads    []compileUnit
	compiles []compileUnit
	pubs     []Entry

	failFast int
	anonSeq  int

	stdRecvName bool
	keepDoc     bool
}

func (p *pkgCtx) ignoref(feats int, decl clang.Cursor, format string, args ...any) {
	if feats&featQuietIgnore == 0 || debugQuietIgnore {
		p.logf(decl, format, args...)
		p.failFast--
		if p.failFast == 0 {
			os.Exit(1)
		}
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

func (p *pkgCtx) addLoadUnit(f compileFunc) {
	p.loads = append(p.loads, compileUnit{
		fn: f,
		at: p.pkg.CurFile(),
	})
}

func (p *pkgCtx) addCompileUnit(f compileFunc) {
	p.compiles = append(p.compiles, compileUnit{
		fn: f,
		at: p.pkg.CurFile(),
	})
}

func (p *pkgCtx) compile() {
	doCompile(p, &p.loads)
	if debugMajorProc {
		log.Println("==> complete uninitialized type declarations")
	}
	pkg := p.pkg
	for _, typDecl := range p.typdecls {
		if typDecl.State() == gogen.TyStateUninited {
			typDecl.InitType(pkg, types.NewStruct(nil, nil))
		}
	}
	doCompile(p, &p.compiles)
}

func doCompile(p *pkgCtx, pcompiles *[]compileUnit) {
	pkg := p.pkg
	for {
		compiles := *pcompiles
		if len(compiles) == 0 {
			break
		}
		*pcompiles = nil
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
	if _, ok := p.rename[cName]; !ok {
		// insert alias name if not already present
		p.rename[cName] = goName
	}
}

func (p *pkgCtx) addType(cName string, decl clang.Cursor, typNamed *types.Named) {
	typObj := typeObj{typNamed.Obj(), 0}
	p.types[cName] = typObj
	p.aliasTypeName(cName, typObj.Name())

	if debugCompileDecl {
		p.logf(decl, "==> addType %s: %v", cName, typNamed)
	}
}

func (p *pkgCtx) ignoreType(cName string, feats int) {
	if o, ok := p.types[cName]; ok {
		if o.TypeName != nil && o.Pkg() == p.pkg.Types { // don't ignore external type
			o.feats |= feats
			p.types[cName] = o
		}
		return
	}
	p.types[cName] = typeObj{nil, feats}
}

func (p *pkgCtx) isConfTypeIgnored(cName string) bool {
	return contains(cName, p.typeIgnores)
}

func (p *pkgCtx) nsFeats(ns string) int {
	if ns != "" {
		if o, ok := p.types[ns]; ok {
			return o.feats
		}
	}
	return 0
}

func (p *pkgCtx) getTypeObj(cName string, feats *int) (ret *types.TypeName, ok bool) {
	o, ok := p.types[cName]
	if ok {
		ret, ok = o.TypeName, o.feats&featAllIgnore == 0
		*feats |= o.feats
	}
	return
}

func (p *pkgCtx) goNamedTypeObj(typName string) (ret *types.TypeName, ok bool) {
	var scope *types.Scope
	pos := strings.IndexByte(typName, '.')
	if pos < 0 {
		scope = types.Universe
	} else {
		var pkg *types.Package
		var pkgPath = typName[:pos]
		switch pkgPath {
		case "":
			pkg = p.pkg.Types
		case "c":
			pkg = p.c.Types
		default:
			pkg = p.pkg.TryImport(pkgPath).Types
			if pkg == nil {
				return
			}
		}
		scope = pkg.Scope()
		typName = typName[pos+1:]
	}
	o := scope.Lookup(typName)
	if o != nil {
		ret, ok = o.(*types.TypeName)
	}
	return
}

func (p *pkgCtx) goTypeArgs(typArgs string) (ret []types.Type, found bool) {
	parts := strings.Split(typArgs, ",")
	ret = make([]types.Type, len(parts))
	for i, part := range parts {
		o, ok := p.goNamedTypeObj(strings.TrimSpace(part))
		if !ok {
			return
		}
		ret[i] = o.Type()
	}
	found = true
	return
}

func (p *pkgCtx) goNamedType(goName string) (ret types.Type, found bool) {
	typArgs := ""
	if goName[len(goName)-1] == ']' { // has typeArgs
		pos := strings.IndexByte(goName, '[')
		if pos <= 0 {
			return
		}
		typArgs = goName[pos+1 : len(goName)-1]
		goName = goName[:pos]
	}
	obj, ok := p.goNamedTypeObj(goName)
	if !ok {
		return
	}
	ret = obj.Type()
	if typArgs != "" {
		targs, ok := p.goTypeArgs(typArgs)
		if !ok {
			return
		}
		inst, err := types.Instantiate(p.typeCtx(), ret, targs, true)
		if err != nil {
			return
		}
		ret = inst
	}
	found = true
	return
}

func (p *pkgCtx) typeAliasOf(name string) (ret types.Type, found bool) {
	newName, ok := p.typeAlias[name]
	if !ok || newName == "" {
		return
	}
	ret, found = p.goNamedType(newName)
	if !found {
		log.Printf("[WARN] alias %s => %s in config: target type not found\n", name, newName)
	}
	return
}

func (p *pkgCtx) isMacroIgnored(cName string) bool {
	return contains(cName, p.macroIgnore)
}

func (p *pkgCtx) isNSIgnored(cName string) bool {
	return contains(cName, p.nsIgnore)
}

func (p *pkgCtx) typeCtx() *types.Context {
	return nil // TODO(xsw): check this
}

// -----------------------------------------------------------------------------
