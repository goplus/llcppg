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

	"github.com/goplus/gogen"
	"github.com/goplus/llcppg/clang"
	lc "github.com/goplus/llcppg/lib/clang"
)

const (
	DbgFlagCompileDecl = 1 << iota
	DbgFlagLoadDeps
	DbgFlagAll = DbgFlagCompileDecl | DbgFlagLoadDeps
)

var (
	debugCompileDecl bool
	debugLoadDeps    bool
)

func SetDebug(flags int) {
	debugCompileDecl = (flags & DbgFlagCompileDecl) != 0
	debugLoadDeps = (flags & DbgFlagLoadDeps) != 0
}

// -----------------------------------------------------------------------------

// Package represents a generated Go package.
type Package struct {
	*gogen.Package
	Wrap *WrapFile
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

	// NameLookup looks up the archive path for a given mangling name. It returns the
	// archive path and a boolean indicating whether the lookup was successful. If not
	// specified, llcppg uses a default lookup function that returns an empty archivePath
	// and true (it means any mangling name is considered found).
	NameLookup func(manglingName string) (archivePath string, ok bool)

	// PackageOf returns the package path for a given header file. If ok is false, it means
	// we don't know the package path for the header file. If not specified, llcppg will assume
	// all header files belong to the same package.
	PackageOf func(headerFile string) (pkgPath string, ok bool)
}

// -----------------------------------------------------------------------------

// Source represents a source file to be processed by llcppg.
type Source struct {
	TU clang.TranslationUnit
}

// -----------------------------------------------------------------------------

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
		DefaultGoFile:   "",
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
		pkg: pkg, cb: pkg.CB(), llgo: llgo, fset: pkg.Fset, c: c, lang: conf.Language,
		cflags: conf.CFlags, wrapFileHeader: conf.WrapFileHeader, nameLookup: nameLookup,
		fileBases: make(map[clang.File]int), methods: make(map[string]*classMethod),
		macroVals: make(map[string]any), types: make(map[string]types.Type),
	}
	loadFiles(ctx, files, pkgPath, conf.PackageOf)
	ctx.compile()
	ret.Package = pkg
	ret.Wrap = ctx.wrap
	return
}

func defaultNameLookup(manglingName string) (archivePath string, ok bool) {
	return "", true
}

// -----------------------------------------------------------------------------

func loadFiles(ctx *pkgCtx, files []Source, myPkgPath string, pkgOf func(headerFile string) (pkgPath string, ok bool)) {
	scope := &scopeCtx{
		overloads: make(map[string]*overloads),
	}
	for _, file := range files {
		tu := file.TU
		clang.VisitChildren(tu.Cursor(), func(decl, parent clang.Cursor) clang.ChildVisitResult {
			if pkgOf != nil {
				at := clang.PresumedFile(decl.Location())
				if pkgPath, ok := pkgOf(at); !ok || pkgPath != myPkgPath {
					return clang.Continue
				}
			}
			loadDecl(ctx, scope, decl, "")
			return clang.Continue
		})
	}
	scope.reorder()
}

func loadDecl(ctx *pkgCtx, scope *scopeCtx, decl clang.Cursor, ns string) {
	/* if global {
		ctx.logFile(decl)
		if decl.IsImplicit || ctx.inDepPkg {
			continue
		}
	} */
	switch decl.Kind {
	case lc.CursorFunctionDecl:
		loadGlobalFunc(ctx, scope, decl, ns)
	case lc.CursorClassDecl, lc.CursorStructDecl:
		defaultInPublic := decl.Kind == lc.CursorStructDecl
		loadClass(ctx, decl, ns, defaultInPublic)
	case lc.CursorCXXMethod, lc.CursorConstructor, lc.CursorDestructor:
		loadOutsideMethod(ctx, decl)
	case lc.CursorTypedefDecl:
		loadTypedef(ctx, decl, ns)
	case lc.CursorEnumDecl:
		// compileEnum(ctx, decl, global)
	case lc.CursorMacroDefinition:
		loadMacro(ctx, decl)
	case lc.CursorInclusionDirective:
		// noop
	case lc.CursorNamespace:
		loadNamespace(ctx, scope, decl, ns)
	case lc.CursorVarDecl:
		// compileVarDecl(ctx, decl, global)
	default:
		log.Panicln("compileDecl: unknown kind =", decl.Kind)
	}
}

func loadNamespace(ctx *pkgCtx, scope *scopeCtx, namespace clang.Cursor, ns string) {
	ns = ns + clang.String(namespace) + "_"
	clang.VisitChildren(namespace, func(decl, parent clang.Cursor) clang.ChildVisitResult {
		loadDecl(ctx, scope, decl, ns)
		return clang.Continue
	})
}

// -----------------------------------------------------------------------------
