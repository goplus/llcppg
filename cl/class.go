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
	"go/types"
	"log"

	"github.com/goplus/gogen"
	"github.com/goplus/llcppg/clang"
	lc "github.com/goplus/llcppg/lib/clang"
)

// -----------------------------------------------------------------------------

type classMethod struct {
	decl         clang.Cursor // declared inside of class
	outsideDecl  clang.Cursor // inline method declared outside of class
	manglingName string
	isPublic     bool
}

type classCtx struct {
	typDecl       *gogen.TypeDecl
	fields        []*types.Var
	publicMethods []*classMethod
	inPublic      bool
}

func compileClass(ctx *blockCtx, cls clang.Cursor, defaultInPublic bool) {
	origName := clang.String(cls)
	if debugCompileDecl {
		log.Println("class", origName)
	}

	pkg := ctx.pkg
	pkgTypes := pkg.Types
	clsName, rewritten := ctx.getPubName(origName)
	typDecl := pkg.NewTypeDefs().NewType(clsName, goNode(ctx, cls))

	ctxCls := &classCtx{
		typDecl:  typDecl,
		inPublic: defaultInPublic,
	}
	clang.VisitChildren(cls, func(decl, parent clang.Cursor) clang.ChildVisitResult {
		compileClassMember(ctx, pkgTypes, ctxCls, decl)
		return clang.Continue
	})

	typStruc := types.NewStruct(ctxCls.fields, nil)
	typNamed := typDecl.InitType(pkg, typStruc)
	if rewritten {
		scope := pkgTypes.Scope()
		substObj(pkgTypes, scope, origName, typNamed.Obj())
	}

	ctx.clTasks = append(ctx.clTasks, func() {
		compilePublicMethods(ctx, ctxCls, typNamed)
	})
}

func compileOutsideMethod(ctx *blockCtx, outsideDecl clang.Cursor) {
	manglingName := clang.Mangling(outsideDecl)
	if m, ok := ctx.methods[manglingName]; ok {
		if m.outsideDecl.Kind == 0 {
			m.outsideDecl = outsideDecl
		} else {
			log.Panicln("method redeclared -", clang.DisplayName(outsideDecl))
		}
	} else {
		log.Panicln("method undeclared -", clang.DisplayName(outsideDecl))
	}
}

func compilePublicMethods(ctx *blockCtx, cls *classCtx, typNamed *types.Named) {
	for _, method := range cls.publicMethods {
		if method.outsideDecl.Kind != 0 {
			compileFunc(ctx, method.outsideDecl, typNamed)
		} else {
			compileFunc(ctx, method.decl, typNamed)
		}
	}
}

func compileClassMember(ctx *blockCtx, pkg *types.Package, cls *classCtx, decl clang.Cursor) {
	switch decl.Kind {
	case lc.CursorCXXMethod, lc.CursorConstructor, lc.CursorDestructor:
		manglingName := clang.Mangling(decl)
		isPublic := cls.inPublic
		method := &classMethod{decl: decl, manglingName: manglingName, isPublic: isPublic}
		ctx.methods[manglingName] = method
		if isPublic {
			cls.publicMethods = append(cls.publicMethods, method)
		}

	case lc.CursorFieldDecl:
		origName := clang.String(decl)
		fldType := toType(ctx, pkg, decl.Type(), flagIsStructField)
		fldName := origName
		if cls.inPublic {
			fldName, _ = ctx.getPubName(origName)
		}
		fld := types.NewField(goNodePos(ctx, decl), pkg, fldName, fldType, false)
		cls.fields = append(cls.fields, fld)

	case lc.CursorCXXAccessSpecifier:
		cls.inPublic = decl.CXXAccessSpecifier() == lc.CXXPublic

	default:
		log.Panicln("compileClassMember: unknown kind =", decl.Kind)
	}
}

// -----------------------------------------------------------------------------
