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

	// PresumedFiles specifies the list of files that are presumed to be included in the
	// compilation. This is used to determine which files are considered part of the
	// package being compiled (optional).
	PresumedFiles []string
}

// -----------------------------------------------------------------------------

const (
	headerGoFile = "llcppg.i.go"
)

// NewPackage loads a translation unit and generates a Go package with the given package
// path, name and configuration.
func NewPackage(pkgPath, pkgName string, conf *Config, tu clang.TranslationUnit, files ...string) (ret Package, err error) {
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
		DefaultGoFile:   headerGoFile,
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
		pkg: pkg, cb: pkg.CB(), llgo: llgo, fset: pkg.Fset, tu: tu, c: c,
		lang: conf.Language, cflags: conf.CFlags, wrapFileHeader: conf.WrapFileHeader,
		nameLookup: nameLookup, methods: make(map[string]*classMethod),
	}
	ctx.initFiles(files)
	loadFiles(ctx)
	ctx.compile()
	ret.Package = pkg
	ret.Wrap = ctx.wrap
	return
}

func defaultNameLookup(manglingName string) (archivePath string, ok bool) {
	return "", true
}

// -----------------------------------------------------------------------------

func loadFiles(ctx *pkgCtx) {
	scope := &scopeCtx{
		overloads: make(map[string]*overloads),
	}
	clang.VisitChildren(ctx.tu.Cursor(), func(decl, parent clang.Cursor) clang.ChildVisitResult {
		loadDecl(ctx, scope, decl)
		return clang.Continue
	})
	scope.reorder()
}

func loadDecl(ctx *pkgCtx, scope *scopeCtx, decl clang.Cursor) {
	/* if global {
		ctx.logFile(decl)
		if decl.IsImplicit || ctx.inDepPkg {
			continue
		}
	} */
	switch decl.Kind {
	case lc.CursorFunctionDecl:
		loadGlobalFunc(ctx, scope, decl)
	case lc.CursorClassDecl, lc.CursorStructDecl:
		defaultInPublic := decl.Kind == lc.CursorStructDecl
		loadClass(ctx, decl, defaultInPublic)
	case lc.CursorCXXMethod, lc.CursorConstructor, lc.CursorDestructor:
		loadOutsideMethod(ctx, decl)
	case lc.CursorVarDecl:
		// compileVarDecl(ctx, decl, global)
	case lc.CursorTypedefDecl:
		loadTypedef(ctx, decl)
	case lc.CursorEnumDecl:
		// compileEnum(ctx, decl, global)
	default:
		log.Panicln("compileDecl: unknown kind =", decl.Kind)
	}
}

func loadTypedef(ctx *pkgCtx, decl clang.Cursor) {
	ctx.compiles = append(ctx.compiles, func(ctx *pkgCtx) {
		origName := clang.String(decl)
		pkg := ctx.pkg
		pkgTypes := pkg.Types
		underlying := decl.TypedefDeclUnderlyingType()
		if debugCompileDecl {
			log.Println("typedef", origName, "-", clang.String(underlying))
		}
		tunder := toType(ctx, pkgTypes, underlying, flagIsTypedef)
		name, rewritten := ctx.getPubName(origName, -1)
		t := pkg.NewTypeDefs().AliasType(name, tunder)
		if rewritten {
			pkgTypes.Scope().Insert(types.NewTypeName(token.NoPos, pkgTypes, origName, t))
		}
	})
}

// -----------------------------------------------------------------------------
