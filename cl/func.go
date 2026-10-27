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
	"github.com/qiniu/x/ctype"

	lc "github.com/llarhub/clang-c"
)

// -----------------------------------------------------------------------------

const operatorPrefix = "operator"

func isOperator(baseName string) bool {
	if len(baseName) > len(operatorPrefix) && strings.HasPrefix(baseName, operatorPrefix) {
		c := baseName[len(operatorPrefix)]
		return !ctype.Is(ctype.CSYMBOL_FIRST_CHAR|ctype.SPACE_BAR, rune(c))
	}
	switch baseName {
	case "new", "delete", "new[]", "delete[]":
		return true
	}
	return false
}

// -----------------------------------------------------------------------------

func loadGlobalFunc(ctx *pkgCtx, obj *overloadObj) {
	ctx.addCompileUnit(func(ctx *pkgCtx) {
		compileFuncOrMethod(ctx, obj, nil)
	})
}

// -----------------------------------------------------------------------------

func isTemplateFunc(fn clang.Cursor) bool {
	if fn.Kind == lc.Cursor_FunctionTemplate {
		return true
	}
	sct := fn.SpecializedCursorTemplate()
	return sct.IsNull() == 0
}

// compileFuncOrMethod compiles a C/C++ function or method into a Go function.
//
// When cls is nil, it is compiled as a receiver-less global function with a
// //go:linkname directive. A C++ static method is compiled this way too: it
// has no implicit "this", so it is loaded as a global function (its Go name is
// prefixed by the enclosing class name, which acts like a namespace). When cls
// is non-nil, it is an instance method compiled with a "this" receiver.
func compileFuncOrMethod(ctx *pkgCtx, obj *overloadObj, this *classCtx) {
	fn := obj.decl
	ns := cNS(fn)
	cName := cNameWithNS(clang.String(fn), ns)
	if isTemplateFunc(fn) {
		ctx.ignoref(featQuietIgnore, fn, "func %s: is template, ignored", cName)
		return
	}

	if ctx.isConfFuncIgnored(cName) {
		ctx.ignoref(featQuietIgnore, fn, "func %s: ignored by config", cName)
		return
	}

	isOp := isOperator(obj.name)
	if debugCompileDecl {
		fnType := clang.String(fn.Type())
		ctx.logf(fn, "func %s: %s - isOp: %v", cName, fnType, isOp)
	}
	if isOp {
		return // TODO(xsw): support operator
	}

	if feats := ctx.nsFeats(ns); feats&featAllIgnore != 0 {
		ctx.ignoref(feats, fn, "func %s: its parent is ignored, ignored too", cName)
		ctx.ignoreType(cName, feats)
		return
	}

	manglingName := clang.Mangling(fn)
	if manglingName == "" {
		ctx.ignoref(NoManglingIgnore, fn, "func %s: no mangled symbol, ignored", cName)
		return
	}

	doCompileFunc(ctx, cName, manglingName, clangFunc{Cursor: fn}, obj, this)
}

func doCompileFunc[funcT fnCursor](ctx *pkgCtx, cName, manglingName string, fn funcT, obj *overloadObj, this *classCtx) {
	name := obj.name

	pkg := ctx.pkg
	pkgTypes := pkg.Types

	feats := 0
	scope := this.scope()
	params, variadic := newParamsEx(ctx, pkgTypes, fn, &feats, scope)
	results := toFuncResults(ctx, pkgTypes, fn.ResultType(), &feats, scope)
	if feats&featAllIgnore != 0 {
		ctx.exIgnoref(feats, fn, "func %s: function with unsupported features, ignored", cName)
		return
	}

	if fn.IsFunctionInlined() != 0 {
		if ctx.cflags == "" {
			ctx.exIgnoref(InlineFuncIgnore, fn, "func %s: inline function but no CFlags in config, ignored", cName)
			return
		}
		manglingName = wrapInlineFunc(ctx, manglingName, fn, this)
	} else if _, ok := ctx.nameLookup(manglingName); !ok {
		ctx.exIgnoref(featExplicitIgnore, fn, "func %s: symbol not found in lib files, ignored", cName)
		return
	}

	var fnName string
	var recv *types.Var
	var typRecv *types.Named // if tryToMethod succeeded, this is the recv
	var typName, typCName string
	var nameInPkg string
	if this == nil {
		if ctx.lang == LanguageC {
			if creator, clsRecv, obj, mayMethod, mayCreator := ctx.creatorCheck(pkgTypes, cName); mayCreator {
				if creator != "" {
					if newResults, ok := tryNewResults(pkgTypes, obj, results); ok {
						results = newResults
						creator = ctx.cstyleToGo(creator, false, true)
						if mayMethod {
							params, recv, typRecv, typName = tryToMethod(ctx, pkgTypes, fn, params, mayMethod)
							if typRecv != nil {
								fnName = creator
								if clsRecv != "" {
									if objRecv := pkgTypes.Scope().Lookup(clsRecv); objRecv != nil {
										tryNewRecv(fn, &recv, &typRecv, &typName, &typCName, clsRecv, objRecv)
									}
								}
								typRecv = nil // skip redeclaration check
							}
						} else {
							fnName = creator
						}
					}
				}
			} else if cls, obj, mayClass := ctx.methodCheck(pkgTypes, cName); mayClass {
				// try to method for C global functions
				params, recv, typRecv, typName = tryToMethod(ctx, pkgTypes, fn, params, cls != "")
				if typRecv != nil {
					tryNewRecv(fn, &recv, &typRecv, &typName, &typCName, cls, obj)
				}
			}
		}
	} else {
		typNamed := this.typNamed
		recv = types.NewParam(fn.NodePos(ctx), pkgTypes, "this", types.NewPointer(typNamed))
		typName = typNamed.Obj().Name()
	}

	if fnName == "" {
		fnName = ctx.funcName(name, obj.order(), typName, typCName, this == nil, true)
	}

	if recv == nil {
		nameInPkg = fnName
	} else {
		if typRecv != nil {
			if pos, isField, exists := findMember(typRecv, fnName); exists {
				newName := ctx.funcName(name, obj.order(), "", typCName, true, true)
				if newName == fnName {
					if isField {
						newName = "Get" + fnName
					} else {
						ctx.exErrorf(fn, "%s redeclared in this block\n\t%v: other declaration of %s", fnName, ctx.position(pos), fnName)
						return
					}
				}
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
	f, err := pkg.NewFuncWith(fn.NodePos(ctx), fnName, sig, nil)
	if err != nil {
		ctx.exErrorf(fn, "%v", err)
		return
	}
	if obj.decl.Kind != 0 {
		ctx.compiledFns[objUSR(obj.decl)] = f
	}

	if recv == nil {
		ctx.forceImportUnsafe()
		f.SetComments(pkg, ctx.directiveCommentsEx(fn,
			"\n//go:linkname "+nameInPkg+" C."+manglingName))
	} else {
		f.SetComments(pkg, ctx.directiveCommentsEx(fn,
			"\n//llgo:link "+nameInPkg+" C."+manglingName))
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

func tryNewRecv[funcT fnCursor](fn funcT, precv **types.Var, typRecv **types.Named, typName, typCName *string, cls string, obj types.Object) {
	recvCType := fn.ArgumentType(0)
	isPtr := recvCType.Kind == lc.Type_Pointer
	if isPtr {
		recvCType = recvCType.Pointee()
	}
	*typCName = trimTypeTag(clang.String(recvCType.Unqualified()))
	if isPtr && cls != "" && cls != *typName {
		if t := obj.Type(); isTypeFromBase(t, *typRecv) {
			recv := *precv
			*precv = types.NewParam(recv.Pos(), recv.Pkg(), recv.Name(), types.NewPointer(t))
			*typRecv = types.Unalias(t).(*types.Named)
			*typName = cls
		}
	}
}

func isTypeFromBase(t types.Type, base *types.Named) bool {
	if t, ok := types.Unalias(t).(*types.Named); ok {
		if t == base {
			return true
		}
		if s, ok := t.Underlying().(*types.Struct); ok {
			if s.NumFields() > 0 {
				return isTypeFromBase(s.Field(0).Type(), base)
			}
		}
	}
	return false
}

func findMember(typ *types.Named, name string) (pos token.Pos, isField bool, exists bool) {
	for i := range typ.NumMethods() {
		m := typ.Method(i)
		if m.Name() == name {
			return m.Pos(), false, true
		}
	}
	if s, ok := typ.Underlying().(*types.Struct); ok {
		for i := range s.NumFields() {
			f := s.Field(i)
			if f.Name() == name {
				return f.Pos(), true, true
			}
		}
	}
	return token.NoPos, false, false
}

func tryNewResults(pkgTypes *types.Package, obj types.Object, results *types.Tuple) (ret *types.Tuple, mayResult bool) {
	if obj == nil || results.Len() != 1 {
		return
	}
	realRet, ok := types.Unalias(results.At(0).Type()).(*types.Pointer)
	if !ok {
		return
	}
	if base, ok := types.Unalias(realRet.Elem()).(*types.Named); ok {
		if t := obj.Type(); isTypeFromBase(t, base) {
			result := types.NewParam(token.NoPos, pkgTypes, "", types.NewPointer(t))
			return types.NewTuple(result), true
		}
	}
	return
}

const (
	c2goMethodRecvName = "self"
)

func tryToMethod[funcT fnCursor](ctx *pkgCtx, pkgTypes *types.Package, fn funcT, params []*types.Var, shouldMethod bool) ([]*types.Var, *types.Var, *types.Named, string) {
	if len(params) > 0 {
		first := params[0]
		t := first.Type()
		if !shouldMethod && len(params) == 2 && types.Identical(params[1].Type(), t) && sameCV(fn) {
			// don't convert to method if the two params have the same type
			return params, nil, nil, ""
		}
		tPtr, isPtr := t.(*types.Pointer)
		if isPtr {
			t = tPtr.Elem()
		}
		switch t := t.(type) {
		case *types.Named:
			if nonClass(ctx, pkgTypes, t.Obj()) {
				break
			}
			if ctx.stdRecvName {
				first = types.NewParam(first.Pos(), pkgTypes, c2goMethodRecvName, first.Type())
			}
			return params[1:], first, t, t.Obj().Name()
		case *types.Alias:
			if nonClass(ctx, pkgTypes, t.Obj()) {
				break
			}
			ta := types.Unalias(t)
			if !isPtr {
				if tp, ok := ta.(*types.Pointer); ok {
					ta = types.Unalias(tp.Elem())
				}
			}
			if tn, ok := ta.(*types.Named); ok && tn.Obj().Pkg() == pkgTypes {
				if ctx.stdRecvName {
					first = types.NewParam(first.Pos(), pkgTypes, c2goMethodRecvName, first.Type())
				}
				return params[1:], first, tn, t.Obj().Name()
			}
		}
	}
	return params, nil, nil, ""
}

func sameCV[funcT fnCursor](fn funcT) bool {
	arg0 := fn.ArgumentType(0)
	arg1 := fn.ArgumentType(1)
	if arg0.Kind == lc.Type_Pointer {
		arg0 = arg0.Pointee()
		arg1 = arg1.Pointee()
	}
	return arg0.IsConstQualified() == arg1.IsConstQualified()
}

func nonClass(ctx *pkgCtx, pkgTypes *types.Package, tn *types.TypeName) bool {
	if tn.Pkg() != pkgTypes {
		return true
	}
	return contains(tn.Name(), ctx.nonClasses)
}

type fnCursor interface {
	nodeCursor
	NumArguments() c.Int
	ArgumentType(i c.Int) lc.Type
	Argument(i c.Int) (typ lc.Type, name string)
	ResultType() lc.Type
	IsVariadic() c.Uint
	IsFunctionInlined() c.Uint
}

type clangFunc struct {
	clangNode
}

func (fn clangFunc) ArgumentType(i c.Int) lc.Type {
	return fn.Cursor.Argument(c.Uint(i)).Type()
}

func (fn clangFunc) Argument(i c.Int) (typ lc.Type, name string) {
	item := fn.Cursor.Argument(c.Uint(i))
	return item.Type(), clang.String(item)
}

func newParams(ctx *pkgCtx, pkg *types.Package, fn clang.Cursor, feats *int, scope *scopeCtx) (params []*types.Var, variadic bool) {
	return newParamsEx(ctx, pkg, clangFunc{Cursor: fn}, feats, scope)
}

func newParamsEx[funcT fnCursor](ctx *pkgCtx, pkg *types.Package, fn funcT, feats *int, scope *scopeCtx) (params []*types.Var, variadic bool) {
	n := fn.NumArguments()
	for i := range n {
		param := newParam(ctx, pkg, fn, i, feats, scope)
		params = append(params, param)
	}
	variadic = fn.IsVariadic() != 0
	if variadic {
		params = append(params, newVariadicParam(pkg))
	}
	return
}

func newParam[funcT fnCursor](ctx *pkgCtx, pkg *types.Package, fn funcT, i c.Int, feats *int, scope *scopeCtx) *types.Var {
	declTyp, declName := fn.Argument(i)
	if debugCompileDecl {
		log.Println("  => param", declName, "-", clang.String(declTyp))
	}
	typ := toTypeEx(ctx, pkg, declTyp, flagIsParam, feats, scope)
	if declName != "" {
		declName = avoidKeyword(declName)
	} else {
		declName = "_llcppg_param" + strconv.Itoa(int(i)+1)
	}
	return types.NewParam(fn.NodePos(ctx), pkg, declName, typ)
}

func avoidKeyword(name string) string {
	switch name {
	case "map", "type", "range", "chan", "var", "func", "go", "select",
		"defer", "package", "import", "interface", "fallthrough":
		return name + "_"
	}
	return name
}

// -----------------------------------------------------------------------------
