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
	"go/token"
	"go/types"
	"log"
	"maps"
	"path/filepath"
	"strings"

	"github.com/goplus/gogen"
	"github.com/goplus/llcppg/clang"
	lc "github.com/llarhub/clang-c"
)

const (
	DbgFlagCompileDecl = 1 << iota
	DbgFlagImport
	DbgFlagAll = DbgFlagCompileDecl | DbgFlagImport
)

var (
	debugCompileDecl bool
	debugImport      bool
)

func SetDebug(flags int) {
	debugCompileDecl = (flags & DbgFlagCompileDecl) != 0
	debugImport = (flags & DbgFlagImport) != 0
}

// -----------------------------------------------------------------------------

// Package represents a generated Go package from a C/C++ library.
type Package struct {
	*gogen.Package
	Wrap   *WrapFile
	Public []Entry // public entries
}

// -----------------------------------------------------------------------------

// Language specifies the programming language of the header file.
type Language int

const (
	LanguageC Language = iota
	LanguageCXX
)

// -----------------------------------------------------------------------------

// Config specifies the configuration for llcppg.
type Config struct {
	// Fset provides source position information for syntax trees and types (optional).
	// If Fset is nil, Load will use a new fileset, but preserve Fset's value.
	Fset *token.FileSet

	// An Importer resolves import paths to Packages (optional).
	Importer types.Importer

	// LLGoPackage specifies the value of the LLGoPackage constant in the generated
	// Go package (optional).
	LLGoPackage string

	// Language specifies the programming language of the header file (default LanguageC).
	Language Language

	// CFlags specifies the compiler flags to be used when compiling the wrapper file.
	// If not specified, llcppg will skip wrapping inline functions/methods.
	CFlags string

	// WrapFileHeader specifies the header content to be included at the top of the generated
	// wrapper file (optional).
	WrapFileHeader string

	// PackageOf returns the package path for a given header file. If ok is false, it means
	// we don't know the package path for the header file. If not specified, llcppg will
	// assume all header files belong to the same package.
	PackageOf func(headerFile string) (pkgPath string, ok bool)

	// PubFileLookup looks up the public file for a given package path. A public file is
	// a text file that contains a list of public C/C++ names and their corresponding Go
	// names (required).
	PubFileLookup func(pkgPath string) (pubFile string, ok bool)

	// NameLookup looks up the archive path for a given mangling name. It returns the
	// archive path and a boolean indicating whether the lookup was successful. If not
	// specified, llcppg uses a default lookup function that returns an empty archivePath
	// and true (it means any mangling name is considered found).
	NameLookup func(manglingName string) (archivePath string, ok bool)

	// Rename specifies a mapping of C/C++ names to Go names (optional). If a name is present
	// in the map, it will be renamed to the corresponding Go name.
	Rename map[string]string

	// TypeAbbr specifies a mapping of Go type name to its abbreviated name. The abbreviated
	// name will be used in function names (optional).
	TypeAbbr map[string]string

	// TypePrefix specifies the prefix to remove from C/C++ type names when generating Go
	// type names (optional).
	TypePrefix []string

	// EnumPrefix specifies the prefix to remove from C/C++ enum value names when generating
	// Go const names (optional).
	EnumPrefix []string

	// FuncPrefix specifies the prefix to remove from C/C++ global function names when
	// generating Go function names (optional).
	FuncPrefix []string

	// Class specifies a list of C/C++ typedef names to be treated as classes (optional).
	Class []string

	// NonClass specifies a list of C/C++ typedef names to be treated as non-classes (optional).
	NonClass []string

	// DefaultGoFile specifies default file name (optional).
	DefaultGoFile string

	// GenMultiGoFiles specifies whether to generate multiple Go files for each header file.
	GenMultiGoFiles bool

	// DontKeepDoc specifies whether to keep the documentation comments in the generated
	// Go package. If true, the documentation comments will be removed (optional).
	DontKeepDoc bool
}

// -----------------------------------------------------------------------------

// Source represents a source file to be processed by llcppg.
type Source = clang.TranslationUnit

// NewPackage loads a translation unit and generates a Go package with the given package
// path, name and configuration.
func NewPackage(pkgPath, pkgName string, files []Source, conf *Config) (ret Package, err error) {
	interp := &nodeInterp{}
	if conf == nil {
		conf = &Config{}
	}
	confGox := &gogen.Config{
		Fset:            conf.Fset,
		Importer:        conf.Importer,
		LoadNamed:       nil,
		HandleErr:       nil,
		NewBuiltin:      nil,
		NodeInterpreter: interp,
		CanImplicitCast: nil,
		DefaultGoFile:   conf.DefaultGoFile,
	}
	pkg := gogen.NewPackage(pkgPath, pkgName, confGox)
	pkg.SetRedeclarable(true)
	interp.fset = pkg.Fset

	llgo := pkg.NewConstDefs(pkg.Types.Scope())
	if llgoPkg := conf.LLGoPackage; llgoPkg != "" {
		llgo.New(func(cb *gogen.CodeBuilder) int {
			cb.Val(llgoPkg)
			return 1
		}, 0, token.NoPos, nil, "LLGoPackage")
	}

	c := pkg.Import("github.com/goplus/lib/c")
	nameLookup := conf.NameLookup
	if nameLookup == nil {
		nameLookup = defaultNameLookup
	}
	ctx := &pkgCtx{
		overloads: make(map[string]*overloads), pkg: pkg, cb: pkg.CB(), llgo: llgo,
		fset: pkg.Fset, c: c, lang: conf.Language, keepDoc: !conf.DontKeepDoc,
		cflags: conf.CFlags, wrapFileHeader: conf.WrapFileHeader, enumPrefix: conf.EnumPrefix,
		typeAbbr: conf.TypeAbbr, typePrefix: conf.TypePrefix, fnPrefix: conf.FuncPrefix,
		rename: conf.Rename, classes: conf.Class, nonClasses: conf.NonClass,
		pkgOf: conf.PackageOf, nameLookup: nameLookup, pubLookup: conf.PubFileLookup,
		fileBases: make(map[clang.File]int), funcs: make(map[string]*funcObj),
		macroVals: make(map[string]any), types: make(map[string]*types.TypeName),
		lastSeen: make(map[string]none), impPkgs: make(map[string]none),
	}
	loadFiles(ctx, files, pkgPath, conf.GenMultiGoFiles)
	ctx.compile()
	ret.Package = pkg
	ret.Wrap = ctx.wrap
	ret.Public = ctx.pubs
	return
}

func defaultNameLookup(manglingName string) (archivePath string, ok bool) {
	return "", true
}

// -----------------------------------------------------------------------------

func loadFiles(ctx *pkgCtx, files []Source, myPkgPath string, genMultiGoFiles bool) {
	pkg := ctx.pkg
	pkgOf := ctx.pkgOf
	scope := &ctx.scopeCtx
	lastSeen := ctx.lastSeen
	for _, f := range files {
		ctx.thisSeen = make(map[string]none) // reset for each file
		clang.VisitChildren(f.Cursor(), func(decl, parent clang.Cursor) clang.ChildVisitResult {
			if pkgOf != nil {
				at := clang.PresumedFile(decl.Location())
				if _, ok := lastSeen[at]; ok {
					return clang.Continue // already loaded
				}
				if pkgPath, ok := pkgOf(at); !ok || pkgPath != myPkgPath {
					return clang.Continue
				}
				if genMultiGoFiles {
					const goFileExt = ".go"
					fname := filepath.Base(at)
					if pos := strings.LastIndex(fname, "."); pos >= 0 {
						fname = fname[:pos] + goFileExt
					} else {
						fname += goFileExt
					}
					pkg.SetCurFile(fname, true)
				}
			}
			loadDecl(ctx, scope, decl, "")
			return clang.Continue
		})
		maps.Copy(lastSeen, ctx.thisSeen)
	}
	scope.reorder()
}

func loadDecl(ctx *pkgCtx, scope *scopeCtx, decl clang.Cursor, ns string) {
	switch decl.Kind {
	case lc.Cursor_FunctionDecl:
		loadGlobalFunc(ctx, scope, decl, ns)
	case lc.Cursor_ClassDecl, lc.Cursor_StructDecl:
		loadClass(ctx, decl, ns, decl.Kind)
	case lc.Cursor_CXXMethod, lc.Cursor_Constructor, lc.Cursor_Destructor:
		loadOutsideMethod(ctx, decl)
	case lc.Cursor_TypedefDecl:
		loadTypedef(ctx, decl, ns)
	case lc.Cursor_EnumDecl:
		loadEnum(ctx, decl, ns)
	case lc.Cursor_MacroDefinition:
		loadMacro(ctx, decl)
	case lc.Cursor_InclusionDirective:
		loadInclude(ctx, decl)
	case lc.Cursor_Namespace:
		loadNamespace(ctx, scope, decl, ns)
	case lc.Cursor_VarDecl:
		loadVar(ctx, decl, ns)
	case lc.Cursor_UnionDecl:
		loadUnion(ctx, decl, ns)
	case lc.Cursor_MacroExpansion:
		// noop
	default:
		log.Panicln("compileDecl: unknown kind =", decl.Kind)
	}
}

func loadNamespace(ctx *pkgCtx, scope *scopeCtx, namespace clang.Cursor, ns string) {
	ns = nsName(ns, clang.String(namespace))
	clang.VisitChildren(namespace, func(decl, parent clang.Cursor) clang.ChildVisitResult {
		loadDecl(ctx, scope, decl, ns)
		return clang.Continue
	})
}

// -----------------------------------------------------------------------------
