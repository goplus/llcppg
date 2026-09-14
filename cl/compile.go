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
	"strconv"

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
}

// -----------------------------------------------------------------------------

// Config specifies the configuration for compiling C/C++ header files.
type Config struct {
	// Fset provides source position information for syntax trees and types.
	// If Fset is nil, Load will use a new fileset, but preserve Fset's value.
	Fset *token.FileSet

	// An Importer resolves import paths to Packages.
	Importer types.Importer

	LLGoPackage string
	Language    string
	CFlags      string

	// Reused specifies to reuse the Package instance between processing multiple C/C++
	// header files.
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
	case lc.CursorClassDecl:
		compileClass(ctx, decl)
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

func compileClass(ctx *blockCtx, cls clang.Cursor) {
	/* TODO(xsw):
	clang.VisitChildren(cls, func(decl, parent clang.Cursor) clang.ChildVisitResult {
		compileDecl(ctx, decl)
		return clang.Continue
	})
	*/
}

// TODO(xsw): method support
func compileFunc(ctx *blockCtx, fn clang.Cursor) {
	manglingName := clang.Mangling(fn)
	origName := clang.String(fn)
	if fn.IsFunctionInlined() != 0 {
		if ctx.cflags == "" {
			if debugCompileDecl {
				log.Println("inline func", origName, "- skipped")
			}
			return
		}
		manglingName = wrapInlineFunc(ctx, origName, fn)
	} else if _, ok := ctx.nameLookup(manglingName); !ok {
		if debugCompileDecl {
			log.Println("func", origName, "- skipped")
		}
		return
	}

	if debugCompileDecl {
		log.Println("func", origName, "-", clang.String(fn.Type()))
	}

	pkg := ctx.pkg
	pkgTypes := pkg.Types
	fnName, rewritten := ctx.getPubName(origName)
	params, variadic := newParams(ctx, pkgTypes, fn)
	results := toFuncResults(ctx, pkgTypes, fn.ResultType())
	sig := types.NewSignatureType(nil, nil, nil, params, results, variadic)
	f, err := pkg.NewFuncWith(goNodePos(ctx, fn), fnName, sig, nil)
	if err != nil {
		log.Panicln("compileFunc:", origName, err)
	}
	ctx.forceImportUnsafe()
	f.SetComments(pkg, &ast.CommentGroup{
		List: []*ast.Comment{
			{Text: "\n//go:linkname " + fnName + " C." + manglingName},
		},
	})
	if rewritten {
		scope := pkg.Types.Scope()
		substObj(pkg.Types, scope, origName, f)
	}
}

func wrapInlineFunc(ctx *blockCtx, origName string, fn clang.Cursor) string {
	if !ctx.hasWrapFile {
		ctx.hasWrapFile = true
		wrapExt := ".c"
		if ctx.lang != "c" {
			wrapExt = ".cpp"
		}
		wrapFile := "_wrap/" + ctx.pkg.Types.Name() + wrapExt
		llgoFiles := ctx.cflags + ": " + wrapFile
		ctx.reused.llgo.New(func(cb *gogen.CodeBuilder) int {
			cb.Val(llgoFiles)
			return 1
		}, 0, token.NoPos, nil, "LLGoFiles")
	}
	wrapName := "_llcppg_" + origName
	// TODO(xsw): wrap inline func
	_ = fn
	return wrapName
}

func newParams(ctx *blockCtx, pkg *types.Package, fn clang.Cursor) (ret *types.Tuple, variadic bool) {
	n := fn.NumArguments()
	var params []*types.Var
	for i := range n {
		item := fn.Argument(c.Uint(i))
		param := newParam(ctx, pkg, item, i)
		params = append(params, param)
	}
	variadic = fn.IsVariadic() != 0
	if variadic {
		params = append(params, newVariadicParam(pkg))
	}
	ret = types.NewTuple(params...)
	return
}

func newParam(ctx *blockCtx, pkg *types.Package, decl clang.Cursor, i c.Int) *types.Var {
	declName := clang.String(decl)
	declTyp := decl.Type()
	if debugCompileDecl {
		log.Println("  => param", declName, "-", clang.String(declTyp))
	}
	typ := toType(ctx, pkg, declTyp, flagIsParam)
	if declName != "" {
		avoidKeyword(&declName)
	} else {
		declName = "_llcppg_param" + strconv.Itoa(int(i)+1)
	}
	return types.NewParam(goNodePos(ctx, decl), pkg, declName, typ)
}

// -----------------------------------------------------------------------------
