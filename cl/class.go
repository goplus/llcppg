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

	"github.com/goplus/llcppg/clang"
	lc "github.com/goplus/llcppg/lib/clang"
)

// -----------------------------------------------------------------------------

type classMethod struct {
	obj          *object      // declared inside of class
	outsideDecl  clang.Cursor // inline method declared outside of class
	manglingName string
	isPublic     bool
}

type classCtx struct {
	scopeCtx
	fields        []*types.Var
	publicMethods []*classMethod
	inPublic      bool
}

func compileClass(ctx *pkgCtx, scope *classCtx, cls clang.Cursor) {
	origName := clang.String(cls)
	if debugCompileDecl {
		log.Println("class", origName)
	}
	pkg := ctx.pkg
	pkgTypes := pkg.Types
	clsName, rewritten := ctx.getPubName(origName, -1)
	typDecl := pkg.NewTypeDefs().NewType(clsName, goNode(ctx, cls))
	typStruc := types.NewStruct(scope.fields, nil)
	typNamed := typDecl.InitType(pkg, typStruc)
	if rewritten {
		substObj(pkgTypes, pkgTypes.Scope(), origName, typNamed.Obj())
	}
	scope.reorder()
	for _, method := range scope.publicMethods {
		obj := method.obj
		if decl := method.outsideDecl; decl.Kind != 0 {
			compileFuncOrMethod(ctx, decl, obj, typNamed)
		} else {
			compileFuncOrMethod(ctx, obj.decl, obj, typNamed)
		}
	}
}

func loadClass(ctx *pkgCtx, cls clang.Cursor, defaultInPublic bool) {
	pkg := ctx.pkg
	pkgTypes := pkg.Types
	scope := &classCtx{
		overloads: make(map[string]*overloads),
		inPublic:  defaultInPublic,
	}
	clang.VisitChildren(cls, func(decl, parent clang.Cursor) clang.ChildVisitResult {
		loadClassMember(ctx, pkgTypes, scope, decl)
		return clang.Continue
	})
	ctx.compiles = append(ctx.compiles, func(ctx *pkgCtx) {
		compileClass(ctx, scope, cls)
	})
}

func loadClassMember(ctx *pkgCtx, pkg *types.Package, cls *classCtx, decl clang.Cursor) {
	switch decl.Kind {
	case lc.CursorCXXMethod, lc.CursorConstructor, lc.CursorDestructor:
		var name string
		switch decl.Kind {
		case lc.CursorConstructor:
			name = "XGo_Ctor"
		case lc.CursorDestructor:
			name = "XGo_Dtor"
		default:
			name = clang.String(decl)
		}
		obj := cls.addObject(name, decl)
		manglingName := clang.Mangling(decl)
		isPublic := cls.inPublic
		method := &classMethod{obj: obj, manglingName: manglingName, isPublic: isPublic}
		ctx.methods[manglingName] = method
		if isPublic {
			cls.publicMethods = append(cls.publicMethods, method)
		}

	case lc.CursorFieldDecl:
		origName := clang.String(decl)
		fldType := toType(ctx, pkg, decl.Type(), flagIsStructField)
		fldName := origName
		if cls.inPublic {
			fldName, _ = ctx.getPubName(origName, -1)
		}
		fld := types.NewField(goNodePos(ctx, decl), pkg, fldName, fldType, false)
		cls.fields = append(cls.fields, fld)

	case lc.CursorCXXAccessSpecifier:
		cls.inPublic = decl.CXXAccessSpecifier() == lc.CXXPublic

	default:
		log.Panicln("loadClassMember: unknown kind =", decl.Kind)
	}
}

func loadOutsideMethod(ctx *pkgCtx, outsideDecl clang.Cursor) {
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

// -----------------------------------------------------------------------------
