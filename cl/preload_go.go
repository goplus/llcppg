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
)

// -----------------------------------------------------------------------------

func loadGoFiles(ctx *pkgCtx, goFiles map[string]*ast.File) {
	if len(goFiles) == 0 {
		return
	}
	pkg := ctx.pkg
	pkgTypes := pkg.Types
	for goFile, f := range goFiles {
		pkg.SetCurFile(goFile, true)
		fileScope := types.NewScope(pkgTypes.Scope(), f.Pos(), f.End(), goFile)
		loadGoFile(ctx, pkgTypes, fileScope, f)
	}
}

func loadGoFile(ctx *pkgCtx, pkg *types.Package, scope *types.Scope, f *ast.File) {
	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			switch d.Tok {
			case token.TYPE:
				for _, spec := range d.Specs {
					t := spec.(*ast.TypeSpec)
					loadGoType(ctx, pkg, scope, t)
				}
			case token.IMPORT:
				for _, spec := range d.Specs {
					imp := spec.(*ast.ImportSpec)
					loadGoImportSpec(ctx, pkg, scope, imp)
				}
			}
		}
	}
}

func loadGoImportSpec(ctx *pkgCtx, pkg *types.Package, scope *types.Scope, imp *ast.ImportSpec) {
	path := toString(ctx, imp.Path, "invalid import path")
	if path == "" {
		return
	}

	imported := ctx.importPkg(path)
	if imported == nil {
		return
	}

	var pos token.Pos
	var name string
	if v := imp.Name; v != nil {
		pos, name = v.Pos(), v.Name
		if name == "_" || name == "." {
			return
		}
	} else {
		pos, name = imp.Path.Pos(), imported.Name()
	}
	pkgName := types.NewPkgName(pos, pkg, name, imported)
	insertObj(ctx, scope, pkgName)
}

func loadGoType(ctx *pkgCtx, pkg *types.Package, scope *types.Scope, t *ast.TypeSpec) {
	tName := t.Name
	name := tName.Name
	if t.Assign != token.NoPos { // alias type
		if t.TypeParams != nil {
			return // NOTE(xsw): doesn't support type parameters
		}
		ctx.addLoadUnit(func(ctx *pkgCtx) {
			aliasGoType(ctx, pkg, t.Pos(), name, goType(ctx, scope, t.Type))
		})
		return
	}
	named := newGoType(ctx, pkg, t.Pos(), name)
	ctx.addLoadUnit(func(ctx *pkgCtx) {
		typ := goType(ctx, scope, t.Type)
		named.SetUnderlying(typ.Underlying())
	})
}

func insertObj(ctx *pkgCtx, scope *types.Scope, obj types.Object) {
	if old := scope.Insert(obj); old != nil {
		pos, name := obj.Pos(), obj.Name()
		prev := ctx.fset.Position(old.Pos())
		ctx.goErrorf(
			pos, pos+token.Pos(len(name)),
			"%s redeclared in this block\n\tprevious declaration at %v\n", name, prev)
	}
}

func aliasGoType(ctx *pkgCtx, pkg *types.Package, pos token.Pos, name string, under types.Type) *types.Alias {
	typName := types.NewTypeName(pos, pkg, name, nil)
	insertObj(ctx, pkg.Scope(), typName)
	return types.NewAlias(typName, under)
}

func newGoType(ctx *pkgCtx, pkg *types.Package, pos token.Pos, name string) *types.Named {
	typName := types.NewTypeName(pos, pkg, name, nil)
	insertObj(ctx, pkg.Scope(), typName)
	return types.NewNamed(typName, nil, nil)
}

// -----------------------------------------------------------------------------

func goType(ctx *pkgCtx, scope *types.Scope, typ ast.Expr) (t types.Type) {
	ret, err := tryType(ctx, scope, typ)
	if err != nil {
		log.Println(err)
		return types.Typ[types.Invalid]
	}
	return ret
}

func tryType(ctx *pkgCtx, scope *types.Scope, typ ast.Expr) (t types.Type, _ error) {
	switch v := typ.(type) {
	case *ast.Ident:
		return tryIdentType(ctx, scope, v)
	case *ast.StarExpr:
		elem, err := tryType(ctx, scope, v.X)
		if err != nil {
			return nil, err
		}
		return types.NewPointer(elem), nil
	case *ast.StructType:
		return tryStructType(ctx, scope, v)
	case *ast.SelectorExpr:
		return tryExternalType(ctx, scope, v)
	case *ast.ArrayType:
		return goArrayType(ctx, scope, v), nil
	case *ast.ParenExpr:
		return tryType(ctx, scope, v.X)
	default:
		return nil, ctx.newCodeErrorf(v.Pos(), v.End(), "goType unexpected: %T", v)
	}
}

func tryIdentType(ctx *pkgCtx, scope *types.Scope, ident *ast.Ident) (typ types.Type, err error) {
	_, o := scope.LookupParent(ident.Name, token.NoPos)
	if o != nil {
		if t, ok := o.(*types.TypeName); ok {
			return t.Type(), nil
		}
	}
	err = ctx.newCodeErrorf(ident.Pos(), ident.End(), "%s is not a type", ident.Name)
	return
}

func goArrayType(ctx *pkgCtx, scope *types.Scope, v *ast.ArrayType) types.Type {
	elem := goType(ctx, scope, v.Elt)
	if v.Len == nil {
		return types.NewSlice(elem)
	}
	return types.NewArray(elem, toInt64(ctx, v.Len, "non-constant array bound"))
}

func tryStructType(ctx *pkgCtx, scope *types.Scope, v *ast.StructType) (typ *types.Struct, err error) {
	pkg := ctx.pkg.Types
	fieldList := v.Fields.List
	fields := make([]*types.Var, 0, len(fieldList))
	for _, field := range fieldList {
		if len(field.Names) == 0 {
			return nil, ctx.newCodeErrorf(field.Pos(), field.End(), "embedded fields are not supported")
		}
		typ := goType(ctx, scope, field.Type)
		for _, name := range field.Names {
			fld := types.NewField(name.NamePos, pkg, name.Name, typ, false)
			fields = append(fields, fld)
		}
	}
	return types.NewStruct(fields, nil), nil
}

func tryExternalType(ctx *pkgCtx, scope *types.Scope, v *ast.SelectorExpr) (typ types.Type, err error) {
	id, ok := v.X.(*ast.Ident)
	if !ok {
		return nil, ctx.newCodeErrorf(v.Pos(), v.End(), "not a type")
	}
	if pkgName, ok := scope.Lookup(id.Name).(*types.PkgName); ok {
		o := pkgName.Imported().Scope().Lookup(v.Sel.Name)
		if t, ok := o.(*types.TypeName); ok {
			return t.Type(), nil
		}
	}
	return nil, ctx.newCodeErrorf(v.Pos(), v.End(), "%s.%s is not a type", id.Name, v.Sel.Name)
}

// -----------------------------------------------------------------------------

func toInt64(ctx *pkgCtx, e ast.Expr, prompt string) int64 {
	switch e := e.(type) {
	case *ast.BasicLit:
		if e.Kind == token.INT {
			i, err := strconv.ParseInt(e.Value, 10, 64)
			if err == nil {
				return i
			}
		}
		ctx.goErrorf(e.Pos(), e.End(), "%s: %s", prompt, e.Value)
	default:
		ctx.goErrorf(e.Pos(), e.End(), "%s", prompt)
	}
	return -1
}

func toString(ctx *pkgCtx, l *ast.BasicLit, prompt string) string {
	if l.Kind == token.STRING {
		s, err := strconv.Unquote(l.Value)
		if err == nil {
			return s
		}
	}
	ctx.goErrorf(l.Pos(), l.End(), "%s: %s", prompt, l.Value)
	return ""
}

// -----------------------------------------------------------------------------
