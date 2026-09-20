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
	decl          clang.Cursor
	typNamed      *types.Named
	fields        []*types.Var
	publicMethods []*classMethod
	staticFuncs   []*object // static methods, compiled as global functions
	inPublic      bool
}

func compileClass(ctx *pkgCtx, scope *classCtx) {
	scope.reorder()
	for _, method := range scope.publicMethods {
		obj := method.obj
		if decl := method.outsideDecl; decl.Kind != 0 {
			compileFuncOrMethod(ctx, decl, obj, scope)
		} else {
			compileFuncOrMethod(ctx, obj.decl, obj, scope)
		}
	}
	// A static method has no implicit "this"; it is compiled as a receiver-less
	// global function (cls == nil), so its enclosing class name acts like a
	// namespace prefix (matching free functions declared inside a namespace).
	for _, obj := range scope.staticFuncs {
		compileFuncOrMethod(ctx, obj.decl, obj, nil)
	}
}

func loadClass(ctx *pkgCtx, cls clang.Cursor, ns string, defaultInPublic bool) {
	pkg := ctx.pkg
	pkgTypes := pkg.Types
	origName := ns + clang.String(cls)
	if debugCompileDecl {
		log.Println("class", origName)
	}
	clsName, rewritten := ctx.getPubName(origName, -1)
	typDecl := pkg.NewTypeDefs().NewType(clsName, goNode(ctx, cls))
	typNamed := typDecl.Type()
	if rewritten {
		substObj(pkgTypes, pkgTypes.Scope(), origName, typNamed.Obj())
	}
	ctx.objects[clang.String(cls.Type())] = typNamed.Obj()
	scope := &classCtx{
		decl:      cls,
		typNamed:  typNamed,
		overloads: make(map[string]*overloads),
		inPublic:  defaultInPublic,
	}
	clang.VisitChildren(cls, func(decl, parent clang.Cursor) clang.ChildVisitResult {
		loadClassMember(ctx, pkgTypes, scope, origName, decl)
		return clang.Continue
	})
	typStruc := types.NewStruct(scope.fields, nil)
	typDecl.InitType(pkg, typStruc)
	ctx.compiles = append(ctx.compiles, func(ctx *pkgCtx) {
		compileClass(ctx, scope)
	})
}

func loadClassMember(ctx *pkgCtx, pkg *types.Package, cls *classCtx, origName string, decl clang.Cursor) {
	switch decl.Kind {
	case lc.CursorCXXMethod, lc.CursorConstructor, lc.CursorDestructor:
		// A static method is not a method: it has no implicit "this". Register
		// it as a global function whose name is prefixed by the enclosing class
		// name (the class name acts like a namespace); it is compiled with a
		// nil class in compileClass.
		if decl.Kind == lc.CursorCXXMethod && decl.CXXMethodIsStatic() != 0 {
			if cls.inPublic {
				obj := cls.addObject(origName+"_"+clang.String(decl), decl)
				cls.staticFuncs = append(cls.staticFuncs, obj)
			}
			return
		}
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

	case lc.CursorCXXBaseSpecifier:
		// A base class is treated the same as a member variable (field) - simply
		// an embedded one. Virtual base classes are not supported for now.
		if decl.IsVirtualBase() != 0 {
			panic("todo: virtual base class is not supported")
		}
		base := baseClass(ctx, decl)
		fld := types.NewField(goNodePos(ctx, decl), pkg, base.Name(), base.Type(), true)
		cls.fields = append(cls.fields, fld)

	default:
		log.Panicln("loadClassMember: unknown kind =", decl.Kind)
	}
}

func baseClass(ctx *pkgCtx, decl clang.Cursor) *types.TypeName {
	t := decl.Type()
	switch t.Kind {
	case lc.TypeRecord:
		cName := fullName(decl)
		if t, ok := ctx.typeObj(cName); ok {
			return t
		}
	}
	panic("baseClass: unknown base class type - " + clang.String(t))
}

func loadOutsideMethod(ctx *pkgCtx, outsideDecl clang.Cursor) {
	// A static method is loaded as a global function, not a classMethod, so it
	// is absent from ctx.methods; its in-class declaration is authoritative.
	if outsideDecl.CXXMethodIsStatic() != 0 {
		return
	}
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
