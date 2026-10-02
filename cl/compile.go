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

	"github.com/goplus/gogen"
	"github.com/goplus/llcppg/clang"
	lc "github.com/llarhub/clang-c"
)

const (
	DbgFlagCompileDecl = 1 << iota
	DbgFlagMajorProc
	DbgFlagQuietIgnore
	DbgFlagAll = DbgFlagCompileDecl | DbgFlagMajorProc | DbgFlagQuietIgnore
)

var (
	debugCompileDecl bool
	debugMajorProc   bool
	debugQuietIgnore bool
)

func SetDebug(flags int) {
	debugCompileDecl = (flags & DbgFlagCompileDecl) != 0
	debugMajorProc = (flags & DbgFlagMajorProc) != 0
	debugQuietIgnore = (flags & DbgFlagQuietIgnore) != 0
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

	// NSIgnore specifies a list of C/C++ namespaces to be ignored (optional).
	NSIgnore []string

	// MacroIgnore specifies a list of C/C++ macro names to be ignored (optional).
	MacroIgnore []string

	// TypeIgnore specifies a list of C/C++ type names to be ignored (optional).
	TypeIgnore []string

	// TypeAlias specifies a mapping of C/C++ type names to another C/C++ name (optional).
	TypeAlias map[string]string

	// TypeAbbr specifies a mapping of Go type name to its abbreviated name(s). The abbreviated
	// name(s) can be a name or name list. They will be used in function names (optional).
	TypeAbbr map[string]any

	// TypeAbbrSuffix specifies the suffix to remove from Go type name when generating function
	// names. It is only valid for types that are not present in TypeAbbr (optional).
	TypeAbbrSuffix []string

	// TypePrefix/TypeSuffix specifies the prefix/suffix to remove from C/C++ type names
	// when generating Go type names (optional).
	TypePrefix, TypeSuffix []string

	// EnumPrefix specifies the prefix to remove from C/C++ enum value names when generating
	// Go const names (optional).
	EnumPrefix []string

	// MacroPrefix specifies the prefix to remove from C/C++ macro names when generating
	// Go const names (optional).
	MacroPrefix []string

	// FuncPrefix specifies the prefix to remove from C/C++ global function names when
	// generating Go function names (optional).
	FuncPrefix []string

	// VarPrefix specifies the prefix to remove from C/C++ global variable names when
	// generating Go variable names (optional).
	VarPrefix []string

	// NSPrefix specifies the prefix to remove from C/C++ namespace names when generating
	// Go package names (optional).
	NSPrefix []string

	// Class specifies a list of C/C++ typedef names to be treated as classes (optional).
	Class []string

	// NonClass specifies a list of C/C++ typedef names to be treated as non-classes (optional).
	NonClass []string

	// DefaultGoFile specifies default file name (optional).
	DefaultGoFile string

	// GoFileOf returns the Go file name for a given header file (optional). If ok is false,
	// it means we don't know the Go file name for the header file and this header file will
	// be skipped. If not specified, llcppg will use the default Go file name (DefaultGoFile)
	// for all header files.
	GoFileOf func(headerFile string) (string, bool)

	// FailFast specifies whether to fail fast on errors (optional). If FailFast is non-zero,
	// llcppg will stop after printing errors `FailFast` times.
	FailFast int

	// LoadLibcPubFile specifies whether to load the pubFile for libc package (optional).
	LoadLibcPubFile bool

	// UseStdRecvName specifies whether to use a standard receiver name (self) for C functions
	// converted to Go methods (optional). If false, the receiver name is taken from the first
	// parameter name. This option does not affect C++ instance methods, whose receiver is
	// always named this.
	UseStdRecvName bool

	// DontKeepDoc specifies whether to keep the documentation comments in the generated
	// Go package. If true, the documentation comments will be removed (optional).
	DontKeepDoc bool
}

const (
	libcPkgPath = "github.com/goplus/lib/c"
)

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
		}, 0, 0, nil, "LLGoPackage")
	}

	nameLookup := conf.NameLookup
	if nameLookup == nil {
		nameLookup = defaultNameLookup
	}
	rename := conf.Rename
	if rename == nil {
		rename = make(map[string]string)
	}
	typdecls := make(map[string]typDecl)
	ctx := &pkgCtx{
		overloads: make(map[string]*overloads), pkg: pkg, cb: pkg.CB(),
		llgo: llgo, fset: pkg.Fset, lang: conf.Language, failFast: conf.FailFast,
		keepDoc: !conf.DontKeepDoc, stdRecvName: conf.UseStdRecvName,
		cflags: conf.CFlags, wrapFileHeader: conf.WrapFileHeader,
		typeAbbr: conf.TypeAbbr, typeAbbrSuffix: conf.TypeAbbrSuffix,
		typePrefix: conf.TypePrefix, typeSuffix: conf.TypeSuffix, typeAlias: conf.TypeAlias,
		fnPrefix: conf.FuncPrefix, enumPrefix: conf.EnumPrefix, rename: rename,
		nsPrefix: conf.NSPrefix, macroPrefix: conf.MacroPrefix, varPrefix: conf.VarPrefix,
		nsIgnore: conf.NSIgnore, macroIgnore: conf.MacroIgnore, typeIgnores: conf.TypeIgnore,
		classes: conf.Class, nonClasses: conf.NonClass, typdecls: typdecls,
		pkgOf: conf.PackageOf, nameLookup: nameLookup, pubLookup: conf.PubFileLookup,
		fileBases: make(map[clang.File]int), fns: make(map[string]*funcObj),
		macroVals: make(map[string]any), types: make(map[string]typeObj),
		lastSeen: make(map[string]none), impPkgs: make(map[string]none),
	}

	if conf.LoadLibcPubFile {
		ctx.c.Types = ctx.importPkg(libcPkgPath)
	} else {
		ctx.c = pkg.Import(libcPkgPath)
	}

	loadFiles(ctx, files, pkgPath, conf.GoFileOf)

	// NOTE(xsw): should complete uninitialized typDecls before compiling
	if debugMajorProc {
		log.Println("==> complete uninitialized type declarations")
	}
	for _, typDecl := range typdecls {
		if typDecl.State() == gogen.TyStateUninited {
			typDecl.InitType(pkg, types.NewStruct(nil, nil))
		}
	}

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

func loadFiles(ctx *pkgCtx, files []Source, myPkgPath string, goFileOf func(headerFile string) (string, bool)) {
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
				pkgPath, ok := pkgOf(at)
				if !ok || pkgPath != myPkgPath {
					return clang.Continue
				}
				if goFileOf != nil {
					fname, ok := goFileOf(at)
					if !ok {
						return clang.Continue
					}
					pkg.SetCurFile(fname, true)
				}
			}
			loadDecl(ctx, scope, decl)
			return clang.Continue
		})
		maps.Copy(lastSeen, ctx.thisSeen)
	}
	scope.reorder()
}

func loadDecl(ctx *pkgCtx, scope *scopeCtx, decl clang.Cursor) {
	switch decl.Kind {
	case lc.Cursor_FunctionDecl:
		loadGlobalFunc(ctx, scope, decl)
	case lc.Cursor_ClassDecl, lc.Cursor_StructDecl:
		loadClass(ctx, decl, decl.Kind, nil)
	case lc.Cursor_CXXMethod, lc.Cursor_Constructor, lc.Cursor_Destructor:
		loadOutsideMethod(ctx, decl)
	case lc.Cursor_TypedefDecl, lc.Cursor_TypeAliasDecl, lc.Cursor_TypeAliasTemplateDecl:
		loadTypedef(ctx, decl, nil)
	case lc.Cursor_EnumDecl:
		loadEnum(ctx, decl)
	case lc.Cursor_MacroDefinition:
		loadMacro(ctx, decl)
	case lc.Cursor_InclusionDirective:
		loadInclude(ctx, decl)
	case lc.Cursor_Namespace:
		loadNamespace(ctx, scope, decl)
	case lc.Cursor_VarDecl:
		loadVar(ctx, decl)
	case lc.Cursor_UnionDecl:
		loadUnion(ctx, decl)
	case lc.Cursor_LinkageSpec: // extern "C" { ... }
		loadLinkageSpec(ctx, scope, decl)
	case lc.Cursor_MacroExpansion, lc.Cursor_StaticAssert, lc.Cursor_UsingDeclaration:
		// noop
	case lc.Cursor_FunctionTemplate:
		// TODO(xsw): ignore for now
	case lc.Cursor_ClassTemplate, lc.Cursor_ClassTemplatePartialSpecialization:
		loadTemplateClass(ctx, decl, nil)
	case lc.Cursor_UnexposedDecl, lc.Cursor_UnexposedAttr:
		// noop
	default:
		ctx.panicf(decl, "loadDecl: unknown kind - %v", decl.Kind)
	}
}

func loadNamespace(ctx *pkgCtx, scope *scopeCtx, namespace clang.Cursor) {
	cName := cNameOf(namespace)
	if ctx.isNSIgnored(cName) {
		if debugCompileDecl {
			ctx.logf(namespace, "namespace %s - ignored", cName)
		}
		return
	}
	clang.VisitChildren(namespace, func(decl, parent clang.Cursor) clang.ChildVisitResult {
		loadDecl(ctx, scope, decl)
		return clang.Continue
	})
}

func loadLinkageSpec(ctx *pkgCtx, scope *scopeCtx, linkage clang.Cursor) {
	clang.VisitChildren(linkage, func(decl, parent clang.Cursor) clang.ChildVisitResult {
		loadDecl(ctx, scope, decl)
		return clang.Continue
	})
}

// -----------------------------------------------------------------------------
