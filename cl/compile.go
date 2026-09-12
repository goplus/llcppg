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

type PkgInfo struct {
}

// Package represents a generated Go package.
type Package struct {
	*gogen.Package
	pi *PkgInfo
}

// Reused specifies to reuse the Package instance between processing multiple C/C++ header files.
type Reused struct {
	pkg Package
}

// -----------------------------------------------------------------------------

// Config specifies the configuration for compiling C/C++ header files.
type Config struct {
	// Fset provides source position information for syntax trees and types.
	// If Fset is nil, Load will use a new fileset, but preserve Fset's value.
	Fset *token.FileSet

	// An Importer resolves import paths to Packages.
	Importer types.Importer

	// Include specifies include searching directories.
	Include []string

	// Reused specifies to reuse the Package instance between processing multiple C/C++ header files.
	*Reused

	// NameLookup looks up the archive path for a given mangling name. It returns the archive
	// path and a boolean indicating whether the lookup was successful.
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
	headerGoFile = "llcppg_header.i.go"
)

// NewPackage creates a new Package instance for the specified package path and name, using
// the provided Source and Config.
func NewPackage(pkgPath, pkgName string, file Source, conf *Config) (pkg Package, err error) {
	if reused := conf.Reused; reused != nil && reused.pkg.Package != nil {
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
		interp.fset = pkg.Fset
	}
	pkg.SetRedeclarable(true)
	pkg.pi, err = loadFile(pkg.Package, conf, file)
	return
}

// -----------------------------------------------------------------------------

func loadFile(p *gogen.Package, conf *Config, file Source) (pi *PkgInfo, err error) {
	c := p.Import("github.com/goplus/lib/c")
	ctx := &blockCtx{
		pkg: p, cb: p.CB(), fset: p.Fset, c: c,
	}
	ctx.initFile(file)
	_ = conf
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
	case lc.CursorVarDecl:
		// compileVarDecl(ctx, decl, global)
	case lc.CursorTypedefDecl:
		/* origName, pub := decl.Name, false
				if global {
					pub = ctx.getPubName(&decl.Name)
				}
				compileTypedef(ctx, decl, global, pub)
				if pub {
					substObj(ctx.pkg.Types, scope, origName, scope.Lookup(decl.Name))
				}
		case ast.RecordDecl:
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
		case ast.EmptyDecl:
		case ast.StaticAssertDecl:
			continue
		*/
	case lc.CursorEnumDecl:
		// compileEnum(ctx, decl, global)
	default:
		log.Panicln("compileDecl: unknown kind =", decl.Kind)
	}
}

// TODO(xsw): method support
func compileFunc(ctx *blockCtx, fn clang.Cursor) {
	fnName := clang.String(fn)
	if debugCompileDecl {
		log.Println("func", fnName, "-", clang.String(fn.Type()))
	}
	origName := fnName
	rewritten := ctx.getPubName(&fnName)
	n := fn.NumArguments()
	var params []*types.Var
	var results *types.Tuple
	for i := range n {
		item := fn.Argument(c.Uint(i))
		param := newParam(ctx, item, i)
		params = append(params, param)
	}
	variadic := fn.IsVariadic() != 0
	if variadic {
		params = append(params, newVariadicParam(ctx))
	}
	pkg := ctx.pkg
	retType := fn.ResultType()
	if retType.Kind != lc.TypeVoid {
		tyRet := toType(ctx, retType, flagRetType)
		results = types.NewTuple(pkg.NewParam(token.NoPos, "", tyRet, false))
	}
	sig := types.NewSignatureType(nil, nil, nil, types.NewTuple(params...), results, variadic)
	f, err := pkg.NewFuncWith(goNodePos(ctx, fn), fnName, sig, nil)
	if err != nil {
		log.Panicln("compileFunc:", fnName, err)
	}
	f.SetComments(pkg, &ast.CommentGroup{
		List: []*ast.Comment{
			{Text: "\n//go:linkname " + fnName + " C." + origName},
		},
	})
	// ctx.addExternFunc(fnName)
	if rewritten {
		scope := pkg.Types.Scope()
		substObj(pkg.Types, scope, origName, f)
	}
	/* origName, rewritten := fnName, false
	if !ctx.inHeader && fn.StorageClass == ast.Static {
		fnName, rewritten = ctx.autoStaticName(origName), true
	} else {
		rewritten = ctx.getPubName(&fnName)
	}
	if body != nil {
		if ctx.checkExists(fnName) {
			return
		}
		isMain := false
		if fnName == "main" && (results != nil || params != nil) {
			fnName, isMain = "_cgo_main", true
		}
		f, err := pkg.NewFuncWith(ctx.goNodePos(fn), fnName, sig, nil)
		if err != nil {
			log.Panicln("compileFunc:", err)
		}
		if rewritten { // for fnName is a recursive function
			scope := pkg.Types.Scope()
			substObj(pkg.Types, scope, origName, f.Obj())
			rewritten = false
		}
		cb := f.BodyStart(pkg)
		ctx.curfn = newFuncCtx(pkg, ctx.markComplicated(fnName, body), origName)
		compileSub(ctx, body)
		checkNeedReturn(ctx, body)
		ctx.curfn = nil
		cb.End()
		if isMain {
			var t *types.Var
			var entryParams *types.Tuple
			var entry = "main"
			var testMain = ctx.testMain
			if testMain {
				entry = "TestMain"
				testing := pkg.Import("testing")
				t = pkg.NewParam(token.NoPos, "t", types.NewPointer(testing.Ref("T").Type()))
				entryParams = types.NewTuple(t)
			}
			pkg.NewFunc(nil, entry, entryParams, nil, false).BodyStart(pkg)
			if results != nil {
				if testMain {
					// if _cgo_ret := _cgo_main(); _cgo_ret != 0 {
					//   t.Fatal("exit status", _cgo_ret)
					// }
					cb.If().DefineVarStart(token.NoPos, retName)
				} else {
					// os.Exit(int(_cgo_main()))
					cb.Val(pkg.Import("os").Ref("Exit")).Typ(types.Typ[types.Int])
				}
			}
			cb.Val(f.Obj())
			if params != nil {
				panic("TODO: main func with params")
			}
			cb.Call(len(params))
			if results != nil {
				if testMain {
					cb.EndInit(1)
					ret := cb.Scope().Lookup(retName)
					cb.Val(ret).Val(0).BinaryOp(token.NEQ).Then().
						Val(t).MemberVal("Fatal").Val("exit status").Val(ret).Call(2).EndStmt().
						End()
				} else {
					cb.Call(1).Call(1)
				}
			}
			cb.EndStmt().End()
		} else {
			delete(ctx.extfns, fnName)
		}
	} else if fn.IsUsed {
		f := types.NewFunc(ctx.goNodePos(fn), pkg.Types, fnName, sig)
		if pkg.Types.Scope().Insert(f) == nil {
			ctx.addExternFunc(fnName)
		}
	}
	if rewritten {
		scope := pkg.Types.Scope()
		substObj(pkg.Types, scope, origName, scope.Lookup(fnName))
	} */
}

var (
	tyValist types.Type = types.NewSlice(gogen.TyAny)
)

func newVariadicParam(ctx *blockCtx) *types.Var {
	return types.NewParam(token.NoPos, ctx.pkg.Types, "__llgo_va_list", tyValist)
}

func newParam(ctx *blockCtx, decl clang.Cursor, i c.Int) *types.Var {
	declName := clang.String(decl)
	declTyp := decl.Type()
	if debugCompileDecl {
		log.Println("  => param", declName, "-", clang.String(declTyp))
	}
	typ := toType(ctx, declTyp, flagIsParam)
	if declName != "" {
		avoidKeyword(&declName)
	} else {
		declName = "_llcppg_param" + strconv.Itoa(int(i)+1)
	}
	return types.NewParam(goNodePos(ctx, decl), ctx.pkg.Types, declName, typ)
}

// -----------------------------------------------------------------------------
