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
	"github.com/goplus/llcppg/ast"
)

func compileMacro(ctx *blockCtx, macro *ast.Macro) {
	panic("todo")
}

// compileTypeDecl compiles a struct/union/class declaration.
func compileTypeDecl(ctx *blockCtx, decl *ast.TypeDecl, global, pub bool) {
	panic("todo")
}

// compileTypedefDecl compiles a typedef declaration.
func compileTypedefDecl(ctx *blockCtx, decl *ast.TypedefDecl, global, pub bool) {
	panic("todo")
}

// compileEnumTypeDecl compiles an enum type declaration.
func compileEnumTypeDecl(ctx *blockCtx, decl *ast.EnumTypeDecl, global bool) {
	panic("todo")
}

func compileFuncDecl(ctx *blockCtx, fn *ast.FuncDecl) {
	panic("todo")
	/*
		fnName, fnType := fn.Name, fn.Type
		if debugCompileDecl {
			log.Println("func", fnName, "-", fnType.QualType, fn.Loc.PresumedLine)
		}
		var hasName bool
		var params []*types.Var
		var results *types.Tuple
		for _, item := range fn.Inner {
			switch item.Kind {
			case ast.ParmVarDecl:
				if debugCompileDecl {
					log.Println("  => param", item.Name, "-", item.Type.QualType)
				}
				if item.Name != "" {
					hasName = true
				}
				params = append(params, newParam(ctx, item))
			case ast.CompoundStmt:
				// body = item
			case ast.BuiltinAttr, ast.FormatAttr, ast.AsmLabelAttr, ast.AvailabilityAttr, ast.ColdAttr, ast.DeprecatedAttr,
				ast.AlwaysInlineAttr, ast.WarnUnusedResultAttr, ast.NoThrowAttr, ast.NoInlineAttr, ast.AllocSizeAttr,
				ast.NonNullAttr, ast.ConstAttr, ast.PureAttr, ast.GNUInlineAttr, ast.ReturnsTwiceAttr, ast.NoSanitizeAttr,
				ast.RestrictAttr, ast.MSAllocatorAttr, ast.VisibilityAttr, ast.C11NoReturnAttr, ast.StrictFPAttr,
				ast.AllocAlignAttr, ast.DisableTailCallsAttr, ast.FormatArgAttr:
			default:
				log.Panicln("compileFunc: unknown kind =", item.Kind)
			}
		}
		variadic := fn.Variadic // TODO(xsw): check variadic
		if variadic {
			params = append(params, newVariadicParam(ctx, hasName))
		}
		pkg := ctx.pkg
		if tyRet := toType(ctx, fnType, parser.FlagGetRetType); ctypes.NotVoid(tyRet) {
			results = types.NewTuple(pkg.NewParam(token.NoPos, "", tyRet))
		}
		/* TODO(xsw): check
		sig := gogen.NewCSignature(types.NewTuple(params...), results, variadic)
		origName, rewritten := fnName, false
		if !ctx.inHeader && fn.StorageClass == ast.Static {
			fnName, rewritten = ctx.autoStaticName(origName), true
		} else {
			rewritten = ctx.getPubName(&fnName)
		}
	*/
	/* TODO(xsw): check
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
	}
	*/
}

/*
func newVariadicParam(ctx *blockCtx, hasName bool) *types.Var {
	panic("todo")
	name := ""
	if hasName {
		name = valistName
	}
	return types.NewParam(token.NoPos, ctx.pkg.Types, name, ctypes.Valist)
}

func newParam(ctx *blockCtx, decl *ast.Node) *types.Var {
	typ := toType(ctx, decl.Type, parser.FlagIsParam)
	// TODO(xsw): check
	// avoidKeyword(&decl.Name)
	return types.NewParam(ctx.goNodePos(decl), ctx.pkg.Types, decl.Name, typ)
}
*/
