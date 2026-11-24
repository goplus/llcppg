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

package tool

import (
	"go/ast"
	"go/token"
	"strings"
)

// -----------------------------------------------------------------------------
// Step 6: build copies of the chosen declarations in which every token.Pos is
// token.NoPos. Input positions are used only for ordering and never leak into
// the output. The copy preserves directive doc comments (dropping all other
// comments) and carries forward the Step-1 package-ref markers so that imports
// can be regenerated from the copied declarations.

// copyDecl returns a deep copy of a top-level declaration with all positions
// cleared and only directive comments retained.
func copyDecl(decl ast.Decl) ast.Decl {
	switch d := decl.(type) {
	case *ast.FuncDecl:
		return &ast.FuncDecl{
			Doc:  copyDirectiveDoc(d.Doc),
			Recv: copyFieldList(d.Recv),
			Name: copyIdent(d.Name),
			Type: copyFuncType(d.Type),
			Body: nil, // bodies are not emitted
		}
	case *ast.GenDecl:
		specs := make([]ast.Spec, len(d.Specs))
		for i, s := range d.Specs {
			specs[i] = copySpec(s)
		}
		return &ast.GenDecl{
			Doc:   copyDirectiveDoc(d.Doc),
			Tok:   d.Tok,
			Specs: specs,
		}
	}
	return decl
}

func copySpec(spec ast.Spec) ast.Spec {
	switch s := spec.(type) {
	case *ast.TypeSpec:
		return &ast.TypeSpec{
			Doc:        copyDirectiveDoc(s.Doc),
			Name:       copyIdent(s.Name),
			TypeParams: copyFieldList(s.TypeParams),
			Assign:     assignTok(s.Assign),
			Type:       copyExpr(s.Type),
		}
	case *ast.ValueSpec:
		names := make([]*ast.Ident, len(s.Names))
		for i, n := range s.Names {
			names[i] = copyIdent(n)
		}
		var values []ast.Expr
		if len(s.Values) > 0 {
			values = make([]ast.Expr, len(s.Values))
			for i, v := range s.Values {
				values[i] = copyExpr(v)
			}
		}
		return &ast.ValueSpec{
			Doc:    copyDirectiveDoc(s.Doc),
			Names:  names,
			Type:   copyExpr(s.Type),
			Values: values,
		}
	case *ast.ImportSpec:
		return &ast.ImportSpec{
			Name: copyIdent(s.Name),
			Path: copyBasicLit(s.Path),
		}
	}
	return spec
}

// assignTok returns a valid-but-NoPos marker for an alias assignment. token.Pos
// cannot be both "valid" and NoPos, so for an alias we use token.Pos(1) which is
// non-zero (IsValid() == true) yet carries no file information; printers that
// rely on relative positions still render "type A = B". A definition keeps
// token.NoPos.
func assignTok(assign token.Pos) token.Pos {
	if assign.IsValid() {
		return token.Pos(1)
	}
	return token.NoPos
}

func copyFuncType(t *ast.FuncType) *ast.FuncType {
	if t == nil {
		return nil
	}
	return &ast.FuncType{
		TypeParams: copyFieldList(t.TypeParams),
		Params:     copyFieldList(t.Params),
		Results:    copyFieldList(t.Results),
	}
}

func copyFieldList(fl *ast.FieldList) *ast.FieldList {
	if fl == nil {
		return nil
	}
	list := make([]*ast.Field, len(fl.List))
	for i, f := range fl.List {
		names := make([]*ast.Ident, len(f.Names))
		for j, n := range f.Names {
			names[j] = copyIdent(n)
		}
		list[i] = &ast.Field{
			Names: names,
			Type:  copyExpr(f.Type),
			Tag:   copyBasicLit(f.Tag),
		}
	}
	return &ast.FieldList{List: list}
}

func copyIdent(id *ast.Ident) *ast.Ident {
	if id == nil {
		return nil
	}
	out := ast.NewIdent(id.Name)
	// Carry forward the Step-1 package-ref marker so imports can be regenerated
	// from the copied declarations.
	if path, ok := packageRefOf(id); ok {
		out.Obj = &ast.Object{Kind: ast.Pkg, Name: id.Name, Data: path}
	}
	return out
}

func copyBasicLit(lit *ast.BasicLit) *ast.BasicLit {
	if lit == nil {
		return nil
	}
	return &ast.BasicLit{Kind: lit.Kind, Value: lit.Value}
}

// copyExpr deep-copies an expression, clearing positions.
func copyExpr(expr ast.Expr) ast.Expr {
	if expr == nil {
		return nil
	}
	switch e := expr.(type) {
	case *ast.Ident:
		return copyIdent(e)
	case *ast.SelectorExpr:
		return &ast.SelectorExpr{X: copyExpr(e.X), Sel: copyIdent(e.Sel)}
	case *ast.StarExpr:
		return &ast.StarExpr{X: copyExpr(e.X)}
	case *ast.ArrayType:
		return &ast.ArrayType{Len: copyExpr(e.Len), Elt: copyExpr(e.Elt)}
	case *ast.Ellipsis:
		return &ast.Ellipsis{Elt: copyExpr(e.Elt)}
	case *ast.StructType:
		return &ast.StructType{Fields: copyFieldList(e.Fields)}
	case *ast.FuncType:
		return copyFuncType(e)
	case *ast.InterfaceType:
		return &ast.InterfaceType{Methods: copyFieldList(e.Methods)}
	case *ast.MapType:
		return &ast.MapType{Key: copyExpr(e.Key), Value: copyExpr(e.Value)}
	case *ast.ChanType:
		return &ast.ChanType{Dir: e.Dir, Value: copyExpr(e.Value)}
	case *ast.ParenExpr:
		return &ast.ParenExpr{X: copyExpr(e.X)}
	case *ast.BasicLit:
		return copyBasicLit(e)
	case *ast.BinaryExpr:
		return &ast.BinaryExpr{X: copyExpr(e.X), Op: e.Op, Y: copyExpr(e.Y)}
	case *ast.UnaryExpr:
		return &ast.UnaryExpr{Op: e.Op, X: copyExpr(e.X)}
	case *ast.CallExpr:
		call := &ast.CallExpr{Fun: copyExpr(e.Fun), Args: copyExprs(e.Args)}
		if e.Ellipsis.IsValid() {
			call.Ellipsis = token.Pos(1)
		}
		return call
	case *ast.IndexExpr:
		return &ast.IndexExpr{X: copyExpr(e.X), Index: copyExpr(e.Index)}
	case *ast.IndexListExpr:
		return &ast.IndexListExpr{X: copyExpr(e.X), Indices: copyExprs(e.Indices)}
	case *ast.CompositeLit:
		return &ast.CompositeLit{Type: copyExpr(e.Type), Elts: copyExprs(e.Elts)}
	case *ast.KeyValueExpr:
		return &ast.KeyValueExpr{Key: copyExpr(e.Key), Value: copyExpr(e.Value)}
	case *ast.SliceExpr:
		sl := &ast.SliceExpr{
			X: copyExpr(e.X), Low: copyExpr(e.Low),
			High: copyExpr(e.High), Max: copyExpr(e.Max),
		}
		sl.Slice3 = e.Slice3
		return sl
	case *ast.TypeAssertExpr:
		return &ast.TypeAssertExpr{X: copyExpr(e.X), Type: copyExpr(e.Type)}
	}
	// Unknown expression kinds are returned as-is; they carry input positions,
	// which is acceptable because llcppg-generated type/value expressions are
	// limited to the kinds handled above.
	return expr
}

func copyExprs(exprs []ast.Expr) []ast.Expr {
	if len(exprs) == 0 {
		return nil
	}
	out := make([]ast.Expr, len(exprs))
	for i, e := range exprs {
		out[i] = copyExpr(e)
	}
	return out
}

// copyDirectiveDoc returns a copy of a doc comment containing only directive
// lines (//word:...), with all other comments dropped. //go:build lines are
// excluded because build constraints are managed by the merge. Returns nil when
// no directives remain.
func copyDirectiveDoc(cg *ast.CommentGroup) *ast.CommentGroup {
	if cg == nil {
		return nil
	}
	var list []*ast.Comment
	for _, c := range cg.List {
		if !strings.HasPrefix(c.Text, "//") {
			continue
		}
		body := c.Text[2:]
		if len(body) == 0 || body[0] == ' ' || body[0] == '\t' {
			continue
		}
		if !isDirective(body) || strings.HasPrefix(body, "go:build") {
			continue
		}
		list = append(list, &ast.Comment{Text: "//" + body})
	}
	if len(list) == 0 {
		return nil
	}
	return &ast.CommentGroup{List: list}
}
