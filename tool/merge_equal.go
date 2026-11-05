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
// Step 5: When are two declarations "the same"?
//
// Comparison works on the input ASTs as annotated in Step 1. Positions and
// comments (except directives) are ignored; function bodies, parameter/result/
// receiver names are ignored for functions; selector expressions are compared by
// import path when both sides are extern package refs.

// declEqual reports whether two declarations are identical under the merge rules.
func declEqual(a, b ast.Decl) bool {
	switch x := a.(type) {
	case *ast.FuncDecl:
		y, ok := b.(*ast.FuncDecl)
		return ok && funcEqual(x, y)
	case *ast.GenDecl:
		y, ok := b.(*ast.GenDecl)
		return ok && genEqual(x, y)
	}
	return false
}

func funcEqual(a, b *ast.FuncDecl) bool {
	if nameOf(a.Name) != nameOf(b.Name) {
		return false
	}
	if !directivesEqual(a.Doc, b.Doc) {
		return false
	}
	if a.Type == nil || b.Type == nil {
		return (a.Type == nil) == (b.Type == nil)
	}
	// Receiver type (for methods); names are ignored.
	if !recvEqual(a.Recv, b.Recv) {
		return false
	}
	// Type parameters.
	if !fieldListEqual(a.Type.TypeParams, b.Type.TypeParams) {
		return false
	}
	// Parameter and result types, including variadic-ness; names are ignored.
	if !fieldListEqual(a.Type.Params, b.Type.Params) {
		return false
	}
	if !fieldListEqual(a.Type.Results, b.Type.Results) {
		return false
	}
	return true
}

func recvEqual(a, b *ast.FieldList) bool {
	aHas := a != nil && len(a.List) > 0
	bHas := b != nil && len(b.List) > 0
	if aHas != bHas {
		return false
	}
	if !aHas {
		return true
	}
	// Compare only the receiver type, not its name.
	return exprEqual(a.List[0].Type, b.List[0].Type)
}

func genEqual(a, b *ast.GenDecl) bool {
	if a.Tok != b.Tok {
		return false
	}
	switch a.Tok {
	case token.CONST:
		// The whole const block is compared as one unit.
		if !directivesEqual(a.Doc, b.Doc) {
			return false
		}
		if len(a.Specs) != len(b.Specs) {
			return false
		}
		for i := range a.Specs {
			if !constSpecEqual(a.Specs[i], b.Specs[i]) {
				return false
			}
		}
		return true
	case token.TYPE:
		if len(a.Specs) != 1 || len(b.Specs) != 1 {
			return len(a.Specs) == len(b.Specs) && allTypeSpecsEqual(a, b)
		}
		if !directivesEqual(a.Doc, b.Doc) {
			return false
		}
		return typeSpecEqual(a.Specs[0], b.Specs[0])
	case token.VAR:
		if len(a.Specs) != 1 || len(b.Specs) != 1 {
			return len(a.Specs) == len(b.Specs) && allValueSpecsEqual(a, b)
		}
		if !directivesEqual(a.Doc, b.Doc) {
			return false
		}
		return varSpecEqual(a.Specs[0], b.Specs[0])
	}
	return false
}

func allTypeSpecsEqual(a, b *ast.GenDecl) bool {
	for i := range a.Specs {
		if !typeSpecEqual(a.Specs[i], b.Specs[i]) {
			return false
		}
	}
	return directivesEqual(a.Doc, b.Doc)
}

func allValueSpecsEqual(a, b *ast.GenDecl) bool {
	for i := range a.Specs {
		if !varSpecEqual(a.Specs[i], b.Specs[i]) {
			return false
		}
	}
	return directivesEqual(a.Doc, b.Doc)
}

func typeSpecEqual(a, b ast.Spec) bool {
	x, ok := a.(*ast.TypeSpec)
	if !ok {
		return false
	}
	y, ok := b.(*ast.TypeSpec)
	if !ok {
		return false
	}
	if nameOf(x.Name) != nameOf(y.Name) {
		return false
	}
	// Alias (type A = B) vs definition.
	if x.Assign.IsValid() != y.Assign.IsValid() {
		return false
	}
	if !fieldListEqual(x.TypeParams, y.TypeParams) {
		return false
	}
	return exprEqual(x.Type, y.Type)
}

func varSpecEqual(a, b ast.Spec) bool {
	x, ok := a.(*ast.ValueSpec)
	if !ok {
		return false
	}
	y, ok := b.(*ast.ValueSpec)
	if !ok {
		return false
	}
	if len(x.Names) != len(y.Names) {
		return false
	}
	for i := range x.Names {
		if nameOf(x.Names[i]) != nameOf(y.Names[i]) {
			return false
		}
	}
	if !exprEqual(x.Type, y.Type) {
		return false
	}
	return exprListEqual(x.Values, y.Values)
}

func constSpecEqual(a, b ast.Spec) bool {
	// Same structure as a var spec (names, declared type, value expressions),
	// but applied member by member so iota/implicit repetition are preserved.
	return varSpecEqual(a, b)
}

// -----------------------------------------------------------------------------
// Expression comparison.

func exprEqual(a, b ast.Expr) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	switch x := a.(type) {
	case *ast.Ident:
		y, ok := b.(*ast.Ident)
		return ok && x.Name == y.Name
	case *ast.SelectorExpr:
		y, ok := b.(*ast.SelectorExpr)
		if !ok {
			return false
		}
		return selectorEqual(x, y)
	case *ast.StarExpr:
		y, ok := b.(*ast.StarExpr)
		return ok && exprEqual(x.X, y.X)
	case *ast.ArrayType:
		y, ok := b.(*ast.ArrayType)
		return ok && exprEqual(x.Len, y.Len) && exprEqual(x.Elt, y.Elt)
	case *ast.Ellipsis:
		y, ok := b.(*ast.Ellipsis)
		return ok && exprEqual(x.Elt, y.Elt)
	case *ast.StructType:
		y, ok := b.(*ast.StructType)
		return ok && structFieldsEqual(x.Fields, y.Fields)
	case *ast.FuncType:
		y, ok := b.(*ast.FuncType)
		if !ok {
			return false
		}
		return fieldListEqual(x.TypeParams, y.TypeParams) &&
			fieldListEqual(x.Params, y.Params) &&
			fieldListEqual(x.Results, y.Results)
	case *ast.InterfaceType:
		y, ok := b.(*ast.InterfaceType)
		return ok && structFieldsEqual(x.Methods, y.Methods)
	case *ast.MapType:
		y, ok := b.(*ast.MapType)
		return ok && exprEqual(x.Key, y.Key) && exprEqual(x.Value, y.Value)
	case *ast.ChanType:
		y, ok := b.(*ast.ChanType)
		return ok && x.Dir == y.Dir && exprEqual(x.Value, y.Value)
	case *ast.ParenExpr:
		y, ok := b.(*ast.ParenExpr)
		return ok && exprEqual(x.X, y.X)
	case *ast.BasicLit:
		y, ok := b.(*ast.BasicLit)
		return ok && x.Kind == y.Kind && x.Value == y.Value
	case *ast.BinaryExpr:
		y, ok := b.(*ast.BinaryExpr)
		return ok && x.Op == y.Op && exprEqual(x.X, y.X) && exprEqual(x.Y, y.Y)
	case *ast.UnaryExpr:
		y, ok := b.(*ast.UnaryExpr)
		return ok && x.Op == y.Op && exprEqual(x.X, y.X)
	case *ast.CallExpr:
		y, ok := b.(*ast.CallExpr)
		if !ok {
			return false
		}
		return exprEqual(x.Fun, y.Fun) && exprListEqual(x.Args, y.Args) &&
			(x.Ellipsis.IsValid() == y.Ellipsis.IsValid())
	case *ast.IndexExpr:
		y, ok := b.(*ast.IndexExpr)
		return ok && exprEqual(x.X, y.X) && exprEqual(x.Index, y.Index)
	case *ast.IndexListExpr:
		y, ok := b.(*ast.IndexListExpr)
		return ok && exprEqual(x.X, y.X) && exprListEqual(x.Indices, y.Indices)
	case *ast.CompositeLit:
		y, ok := b.(*ast.CompositeLit)
		return ok && exprEqual(x.Type, y.Type) && exprListEqual(x.Elts, y.Elts)
	case *ast.KeyValueExpr:
		y, ok := b.(*ast.KeyValueExpr)
		return ok && exprEqual(x.Key, y.Key) && exprEqual(x.Value, y.Value)
	case *ast.SliceExpr:
		y, ok := b.(*ast.SliceExpr)
		return ok && exprEqual(x.X, y.X) && exprEqual(x.Low, y.Low) &&
			exprEqual(x.High, y.High) && exprEqual(x.Max, y.Max)
	case *ast.TypeAssertExpr:
		y, ok := b.(*ast.TypeAssertExpr)
		return ok && exprEqual(x.X, y.X) && exprEqual(x.Type, y.Type)
	}
	// Unknown / unsupported expression kinds compare unequal unless identical
	// pointers, which keeps the comparison conservative.
	return a == b
}

// selectorEqual compares two selector expressions (Step 5).
func selectorEqual(a, b *ast.SelectorExpr) bool {
	ax, aIsID := a.X.(*ast.Ident)
	bx, bIsID := b.X.(*ast.Ident)
	var aPath, bPath importPath
	var aRef, bRef bool
	if aIsID {
		aPath, aRef = packageRefOf(ax)
	}
	if bIsID {
		bPath, bRef = packageRefOf(bx)
	}
	switch {
	case aRef && bRef:
		// Both extern package refs: compare import path and selector name only.
		return aPath == bPath && nameOf(a.Sel) == nameOf(b.Sel)
	case aRef != bRef:
		// Only one side is a package ref: different.
		return false
	default:
		// Ordinary field/method selection: structural comparison.
		return exprEqual(a.X, b.X) && nameOf(a.Sel) == nameOf(b.Sel)
	}
}

func exprListEqual(a, b []ast.Expr) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !exprEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}

// fieldListEqual compares two field lists by name, type, tag and embedding;
// field order matters. For function signatures the field "names" are parameter
// names and are therefore ignored; callers that need name comparison (struct
// fields) use fieldListEqualNamed.
func fieldListEqual(a, b *ast.FieldList) bool {
	return fieldListCompare(a, b, false)
}

func fieldListCompare(a, b *ast.FieldList, compareNames bool) bool {
	if a == nil || b == nil {
		return (a == nil) == (b == nil)
	}
	if len(a.List) != len(b.List) {
		return false
	}
	for i := range a.List {
		fa, fb := a.List[i], b.List[i]
		if compareNames {
			if !identNamesEqual(fa.Names, fb.Names) {
				return false
			}
			if tagValue(fa.Tag) != tagValue(fb.Tag) {
				return false
			}
		} else {
			// For signatures, a grouped field "a, b int" declares as many
			// parameters as it has names (or one if anonymous); the arity must
			// match even though the names are ignored.
			if fieldArity(fa) != fieldArity(fb) {
				return false
			}
		}
		if !exprEqual(fa.Type, fb.Type) {
			return false
		}
	}
	return true
}

// fieldArity returns how many entries a signature field represents.
func fieldArity(f *ast.Field) int {
	if len(f.Names) == 0 {
		return 1
	}
	return len(f.Names)
}

func identNamesEqual(a, b []*ast.Ident) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if nameOf(a[i]) != nameOf(b[i]) {
			return false
		}
	}
	return true
}

func tagValue(tag *ast.BasicLit) string {
	if tag == nil {
		return ""
	}
	return tag.Value
}

func nameOf(id *ast.Ident) string {
	if id == nil {
		return ""
	}
	return id.Name
}

// -----------------------------------------------------------------------------
// Struct field comparison (names, types, tags, embedding matter).

// structFieldsEqual compares struct fields by name, type, tag and embedding.
func structFieldsEqual(a, b *ast.FieldList) bool {
	return fieldListCompare(a, b, true)
}

// -----------------------------------------------------------------------------
// Directives (Step 5): collected from a doc comment, compared as an ordered list
// of normalized lines. //go:build lines are not symbol directives.

func directivesEqual(a, b *ast.CommentGroup) bool {
	da := directivesOf(a)
	db := directivesOf(b)
	if len(da) != len(db) {
		return false
	}
	for i := range da {
		if da[i] != db[i] {
			return false
		}
	}
	return true
}

// directivesOf extracts directive lines (//word:...) from a comment group,
// normalizing insignificant whitespace, and excluding //go:build lines.
func directivesOf(cg *ast.CommentGroup) []string {
	if cg == nil {
		return nil
	}
	var out []string
	for _, c := range cg.List {
		text := c.Text
		if !strings.HasPrefix(text, "//") {
			continue // block comments are never directives
		}
		body := text[2:]
		if len(body) == 0 || body[0] == ' ' || body[0] == '\t' {
			continue // space after // => ordinary comment, not a directive
		}
		if !isDirective(body) {
			continue
		}
		if strings.HasPrefix(body, "go:build") {
			continue // managed by the merge itself
		}
		out = append(out, normalizeDirective(body))
	}
	return out
}

// isDirective reports whether a comment body has the form "word:..." with the
// colon appearing before any whitespace (e.g. "go:linkname", "llgo:link").
func isDirective(body string) bool {
	colon := strings.IndexByte(body, ':')
	if colon <= 0 {
		return false
	}
	word := body[:colon]
	for i := 0; i < len(word); i++ {
		c := word[i]
		isLetter := (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
		isDigit := c >= '0' && c <= '9'
		if !isLetter && !isDigit && c != '_' {
			return false
		}
	}
	return true
}

// normalizeDirective collapses runs of whitespace to single spaces and trims the
// ends, so insignificant whitespace in a directive does not matter.
func normalizeDirective(body string) string {
	return strings.Join(strings.Fields(body), " ")
}
