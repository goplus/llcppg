/*
 * Copyright (c) 2025 The XGo Authors (xgo.dev). All rights reserved.
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
	"reflect"

	"github.com/goplus/gogen"
	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/cl/internal/convert"
)

/* TODO(xsw): remove
var (
	ErrInvalidArg = errors.New("invalid argument")
)

var (
	debugCompileDecl bool = true
)
*/

type Package struct {
	*gogen.Package
	*convert.PkgInfo // TODO(xsw): check
}

type Config struct {
	// Fset provides source position information for syntax trees and types.
	// If Fset is nil, Load will use a new fileset, but preserve Fset's value.
	Fset *token.FileSet

	// An Importer resolves import paths to Packages.
	Importer types.Importer
}

const (
	headerGoFile = "llcppg_header.go"
)

// NewPackage create a Go package from C/C++ header files AST.
// All the C/C++ header files have been parsed and merged into a single AST.
func NewPackage(pkgPath, pkgName string, file *ast.File, conf *Config) (pkg Package, err error) {
	confGox := &gogen.Config{
		Fset:            conf.Fset,
		Importer:        conf.Importer,
		LoadNamed:       nil,
		HandleErr:       nil,
		NewBuiltin:      nil,
		NodeInterpreter: nil, // TODO(xsw): check
		CanImplicitCast: nil, // TODO(xsw): check
		DefaultGoFile:   headerGoFile,
	}
	pkg.Package = gogen.NewPackage(pkgPath, pkgName, confGox)
	pkg.SetRedeclarable(true)
	err = loadFile(pkg.Package, conf, file)
	return
}

func loadFile(p *gogen.Package, _ *Config, file *ast.File) (err error) {
	/* TODO(xsw): check
	srcFile := conf.SrcFile
	if srcFile != "" {
		srcFile, _ = filepath.Abs(srcFile)
	}
	*/
	ctx := &blockCtx{
		pkg: p, cb: p.CB(), fset: p.Fset,
		/* TODO(xsw): check
		unnameds: make(map[ast.ID]unnamedType),
		gblvars:  make(map[string]*gogen.VarDefs),
		ignored:  conf.Ignored,
		public:   conf.Public,
		srcdir:   filepath.Dir(srcFile),
		srcfile:  srcFile,
		src:      conf.Src,
		bfm:      conf.BuiltinFuncMode,
		testMain: conf.TestMain,
		*/
	}
	/* TODO(xsw): check
	baseDir, _ := filepath.Abs(conf.Dir)
	ctx.initMultiFileCtl(p, baseDir, conf)
	ctx.initCTypes()
	ctx.initFile()
	ctx.initPublicFrom(baseDir, conf, file)
	for _, ign := range ctx.ignored {
		if ctx.getPubName(&ign) {
			ctx.ignored = append(ctx.ignored, ign)
		}
	}
	*/
	for _, macro := range file.Macros {
		if ctx.inPkg(macro.Loc) {
			compileMacro(ctx, macro)
		}
	}
	compileDeclStmt(ctx, file, true)
	/* TODO(xsw): check
	if conf.NeedPkgInfo {
		pkgInfo := ctx.PkgInfo // make a copy: don't keep a ref to blockCtx
		pi = &pkgInfo
	}
	*/
	return
}

func compileDeclStmt(ctx *blockCtx, node *ast.File, global bool) {
	// scope := ctx.cb.Scope()
	for _, decl := range node.Decls {
		/* TODO(xsw): check
		if global {
			ctx.logFile(decl)
			if decl.IsImplicit || ctx.inDepPkg {
				continue
			}
		}
		*/
		switch decl := decl.(type) {
		//case ast.VarDecl:
		// compileVarDecl(ctx, decl, global)
		case *ast.TypedefDecl:
			compileTypedefDecl(ctx, decl, global, false)
			/* TODO(xsw): check
			origName, pub := decl.Name, false
			if global {
				pub = ctx.getPubName(&decl.Name)
			}
			compileTypedef(ctx, decl, global, pub)
			if pub {
				substObj(ctx.pkg.Types, scope, origName, scope.Lookup(decl.Name))
			}
			*/
		case *ast.TypeDecl:
			compileTypeDecl(ctx, decl, global, false)
			/* TODO(xsw): check
			pub := false
			name, suKind := ctx.getSuName(decl, decl.TagUsed)
			origName := name
			if global {
				if suKind == suAnonymous {
					// pub = true if this is a public typedef
					pub = i+1 < n && isPubTypedef(ctx, node.Inner[i+1])
				} else {
					pub = ctx.getPubName(&name)
					if decl.CompleteDefinition && ctx.checkExists(name) {
						continue
					}
				}
			}
			typ, del := compileStructOrUnion(ctx, name, decl, pub)
			if suKind != suAnonymous {
				if pub {
					substObj(ctx.pkg.Types, scope, origName, scope.Lookup(name))
				}
				break
			}
			ctx.unnameds[decl.ID] = unnamedType{typ: typ, del: del}
			for i+1 < n {
				next := node.Inner[i+1]
				if next.Kind == ast.VarDecl {
					if ret, ok := checkAnonymous(ctx, scope, typ, next); ok {
						compileVarWith(ctx, ret, next)
						i++
						continue
					}
				}
				break
			}
			*/
		case *ast.EnumTypeDecl:
			compileEnumTypeDecl(ctx, decl, global)
		// case ast.EmptyDecl:
		case *ast.FuncDecl:
			if global {
				compileFuncDecl(ctx, decl)
				continue
			}
			//	fallthrough
			//case ast.StaticAssertDecl:
			continue
		default:
			log.Panicln("compileDeclStmt: unknown decl -", reflect.TypeOf(decl))
		}
	}
}
