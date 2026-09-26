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
	"strconv"
	"strings"

	"github.com/goplus/lib/c"
	"github.com/goplus/llcppg/clang"
	lc "github.com/llarhub/clang-c"
)

// -----------------------------------------------------------------------------

func loadGlobalFunc(ctx *pkgCtx, scope *scopeCtx, decl clang.Cursor, ns string) {
	name := nameWithNS(clang.String(decl), ns)
	if obj, ok := scope.addFunc(ctx, name, decl); ok {
		ctx.addCompileUnit(func(ctx *pkgCtx) {
			compileFuncOrMethod(ctx, obj, nil)
		})
	}
}

// compileFuncOrMethod compiles a C/C++ function or method into a Go function.
//
// When cls is nil, it is compiled as a receiver-less global function with a
// //go:linkname directive. A C++ static method is compiled this way too: it
// has no implicit "this", so it is loaded as a global function (its Go name is
// prefixed by the enclosing class name, which acts like a namespace). When cls
// is non-nil, it is an instance method compiled with a "this" receiver.
func compileFuncOrMethod(ctx *pkgCtx, obj *funcObj, cls *classCtx) {
	fn := obj.decl
	pkg := ctx.pkg
	pkgTypes := pkg.Types
	origName := obj.name
	feats := 0
	params, variadic := newParams(ctx, pkgTypes, fn, &feats)
	results := toFuncResults(ctx, pkgTypes, fn.ResultType(), &feats)
	if feats&featIgnored != 0 {
		if debugCompileDecl {
			log.Println("func", origName, "- ignored")
		}
		return
	}

	manglingName := clang.Mangling(fn)
	if fn.IsFunctionInlined() != 0 {
		if ctx.cflags == "" {
			if debugCompileDecl {
				log.Println("inline func", origName, "- skipped")
			}
			return
		}
		manglingName = wrapInlineFunc(ctx, manglingName, fn, cls)
	} else if _, ok := ctx.nameLookup(manglingName); !ok {
		if debugCompileDecl {
			log.Println("func", origName, "- skipped")
		}
		return
	}

	if debugCompileDecl {
		log.Println("func", origName, "-", clang.String(fn.Type()))
	}

	var recv *types.Var
	var typRecv *types.Named // if tryToMethod succeeded, this is the recv
	var typName, typCName string
	var nameInPkg string
	if cls == nil {
		// TODO(xsw): llgo bugfix - to support method with callback
		if ctx.lang == LanguageC && feats&featHasCallback == 0 {
			// try to method for C global functions
			params, recv, typRecv, typName = tryToMethod(ctx, pkgTypes, params)
			if typRecv != nil {
				recvCType := fn.Argument(0).Type()
				if recvCType.Kind == lc.Type_Pointer {
					recvCType = recvCType.Pointee()
				}
				typCName = clang.String(recvCType.Unqualified())
				if pos := strings.LastIndex(typCName, " "); pos >= 0 {
					// remove type tag prefix, e.g. "struct Foo" => "Foo"
					typCName = typCName[pos+1:]
				}
			}
		}
	} else {
		typNamed := cls.typNamed
		recv = types.NewParam(token.NoPos, pkgTypes, "this", types.NewPointer(typNamed))
		typName = typNamed.Obj().Name()
	}
	fnName := ctx.funcName(origName, obj.order(), typName, typCName, cls == nil, true)
	if recv == nil {
		nameInPkg = fnName
	} else {
		if typRecv != nil {
			if existMember(typRecv, fnName) {
				newName := ctx.funcName(origName, obj.order(), "", typCName, true, true)
				log.Printf("==> member %s.%s already exists, rename to %s\n", typName, fnName, newName)
				fnName = newName
			}
		}
		if _, ok := recv.Type().(*types.Pointer); ok {
			nameInPkg = "(*" + typName + ")." + fnName
		} else {
			nameInPkg = typName + "." + fnName
		}
	}
	sig := types.NewSignatureType(recv, nil, nil, types.NewTuple(params...), results, variadic)
	f, err := pkg.NewFuncWith(goNodePos(ctx, fn), fnName, sig, nil)
	if err != nil {
		log.Panicln("compileFunc:", origName, err)
	}

	if recv == nil {
		ctx.forceImportUnsafe()
		f.SetComments(pkg, ctx.directiveComments(fn,
			"\n//go:linkname "+nameInPkg+" C."+manglingName))
	} else {
		f.SetComments(pkg, ctx.directiveComments(fn,
			"\n// llgo:link "+nameInPkg+" C."+manglingName))
		cb := f.BodyStart(pkg)
		if n := results.Len(); n > 0 {
			retType := results.At(0).Type()
			if types.Identical(recv.Type(), retType) {
				cb.Val(recv)
			} else {
				cb.ZeroLit(retType)
			}
			cb.Return(1)
		}
		cb.End()
	}
}

func existMember(typ *types.Named, name string) bool {
	for i := range typ.NumMethods() {
		if typ.Method(i).Name() == name {
			return true
		}
	}
	if s, ok := typ.Underlying().(*types.Struct); ok {
		for i := range s.NumFields() {
			if s.Field(i).Name() == name {
				return true
			}
		}
	}
	return false
}

const (
	llgoSupportAliasAsRecv = false
)

func tryToMethod(ctx *pkgCtx, pkgTypes *types.Package, params []*types.Var) ([]*types.Var, *types.Var, *types.Named, string) {
	if len(params) > 0 {
		first := params[0]
		t := first.Type()
		if len(params) == 2 && params[1].Type() == t {
			// don't convert to method if the first two params have the same type
			return params, nil, nil, ""
		}
		if tp, ok := t.(*types.Pointer); ok {
			t = tp.Elem()
		}
		switch t := t.(type) {
		case *types.Named:
			if t.Obj().Pkg() == pkgTypes {
				return params[1:], first, t, t.Obj().Name()
			}
		case *types.Alias:
			if t.Obj().Pkg() != pkgTypes {
				break
			}
			if llgoSupportAliasAsRecv {
				ta := types.Unalias(t)
				if tp, ok := ta.(*types.Pointer); ok {
					ta = tp.Elem()
				}
				if tn, ok := ta.(*types.Named); ok {
					return params[1:], first, tn, t.Obj().Name()
				}
			} else {
				ta := types.Unalias(t)
				first = types.NewParam(first.Pos(), pkgTypes, first.Name(), ta)
				if tp, ok := ta.(*types.Pointer); ok {
					ta = tp.Elem()
				}
				if tn, ok := ta.(*types.Named); ok && tn.Obj().Pkg() == pkgTypes {
					// add type abbreviation for alias type
					ctx.typeAbbr[tn.Obj().Name()] = t.Obj().Name()
					return params[1:], first, tn, tn.Obj().Name()
				}
			}
		}
	}
	return params, nil, nil, ""
}

func newParams(ctx *pkgCtx, pkg *types.Package, fn clang.Cursor, feats *int) (params []*types.Var, variadic bool) {
	n := fn.NumArguments()
	for i := range n {
		item := fn.Argument(c.Uint(i))
		param := newParam(ctx, pkg, item, i, feats)
		params = append(params, param)
	}
	variadic = fn.IsVariadic() != 0
	if variadic {
		params = append(params, newVariadicParam(pkg))
	}
	return
}

func newParam(ctx *pkgCtx, pkg *types.Package, decl clang.Cursor, i c.Int, feats *int) *types.Var {
	declName := clang.String(decl)
	declTyp := decl.Type()
	if debugCompileDecl {
		log.Println("  => param", declName, "-", clang.String(declTyp))
	}
	typ := toTypeEx(ctx, pkg, declTyp, flagIsParam, feats)
	if declName != "" {
		avoidKeyword(&declName)
	} else {
		declName = "_llcppg_param" + strconv.Itoa(int(i)+1)
	}
	return types.NewParam(goNodePos(ctx, decl), pkg, declName, typ)
}

// -----------------------------------------------------------------------------
