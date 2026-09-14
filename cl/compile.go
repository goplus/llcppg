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
	"github.com/goplus/lib/c"
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
}

// Reused specifies to reuse the Package instance between processing multiple C/C++
// header files.
type Reused struct {
	pkg  Package
	llgo *gogen.ConstDefs
	wrap *wrapFile
}

// -----------------------------------------------------------------------------

// Language specifies the programming language of the header file.
type Language int

const (
	LanguageC Language = iota
	LanguageCXX
)

// -----------------------------------------------------------------------------

// Config specifies the configuration for compiling header files.
type Config struct {
	// Fset provides source position information for syntax trees and types.
	// If Fset is nil, Load will use a new fileset, but preserve Fset's value.
	Fset *token.FileSet

	// An Importer resolves import paths to Packages.
	Importer types.Importer

	// LLGoPackage specifies the value of the LLGoPackage constant in the generated
	// Go package.
	LLGoPackage string

	// Language specifies the programming language of the header file.
	Language Language

	// CFlags specifies the compiler flags to be used when compiling the wrapper file.
	CFlags string

	// Reused specifies to reuse the Package instance between processing multiple header
	// files.
	*Reused

	// NameLookup looks up the archive path for a given mangling name. It returns the
	// archive path and a boolean indicating whether the lookup was successful.
	NameLookup func(manglingName string) (archivePath string, ok bool)
}

// -----------------------------------------------------------------------------

// Source represents a C/C++ header to compile.
type Source struct {
	TU           clang.TranslationUnit
	Handle       clang.File
	PresumedFile *c.Char
}

// -----------------------------------------------------------------------------

const (
	headerGoFile = "llcppg.i.go"
)

// NewPackage creates a new Package instance for the specified package path and name, using
// the provided Source and Config.
func NewPackage(pkgPath, pkgName string, file Source, conf *Config) (pkg Package, err error) {
	reused := conf.Reused
	if reused == nil {
		reused = new(Reused)
	}
	if reused.pkg.Package != nil {
		pkg = reused.pkg
	} else {
		interp := &nodeInterp{}
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
		pkg.Package = gogen.NewPackage(pkgPath, pkgName, confGox)
		reused.llgo = pkg.Package.NewConstDefs(pkg.Types.Scope())
		interp.fset = pkg.Fset
		if llgoPkg := conf.LLGoPackage; llgoPkg != "" {
			reused.llgo.New(func(cb *gogen.CodeBuilder) int {
				cb.Val(llgoPkg)
				return 1
			}, 0, token.NoPos, nil, "LLGoPackage")
		}
	}
	pkg.SetRedeclarable(true)
	err = loadFile(pkg.Package, conf, file, reused)
	reused.pkg = pkg
	return
}

// -----------------------------------------------------------------------------

func loadFile(p *gogen.Package, conf *Config, file Source, reused *Reused) (err error) {
	c := p.Import("github.com/goplus/lib/c")
	ctx := &blockCtx{
		pkg: p, cb: p.CB(), fset: p.Fset, c: c,
		lang: conf.Language, cflags: conf.CFlags,
		reused: reused, nameLookup: conf.NameLookup,
	}
	ctx.initFile(file)
	clang.VisitChildren(file.TU.Cursor(), func(decl, parent clang.Cursor) clang.ChildVisitResult {
		compileDecl(ctx, decl)
		return clang.Continue
	})
	return
}

func compileDecl(ctx *blockCtx, decl clang.Cursor) {
	/* if global {
		ctx.logFile(decl)
		if decl.IsImplicit || ctx.inDepPkg {
			continue
		}
	} */
	switch decl.Kind {
	case lc.CursorFunctionDecl:
		compileFunc(ctx, decl)
	case lc.CursorClassDecl, lc.CursorStructDecl:
		defaultInPublic := decl.Kind == lc.CursorStructDecl
		compileClass(ctx, decl, defaultInPublic)
	case lc.CursorVarDecl:
		// compileVarDecl(ctx, decl, global)
	case lc.CursorTypedefDecl:
		// TODO(xsw)
	case lc.CursorEnumDecl:
		// compileEnum(ctx, decl, global)
	default:
		log.Panicln("compileDecl: unknown kind =", decl.Kind)
	}
}

// -----------------------------------------------------------------------------
