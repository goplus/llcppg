// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package astutil

import (
	"go/ast"
	"go/token"
	"reflect"
)

// Equal reports whether two nodes are structurally equal,
// ignoring source coordinates, [ast.Object], [ast.Scope], and comments.
//
// The operands x and y may be nil.
// A nil slice is not equal to an empty slice.
//
// The provided function determines whether two identifiers
// should be considered identical.
func Equal(x, y ast.Node, identical func(x, y *ast.Ident) bool) bool {
	if x == nil || y == nil {
		return x == y
	}
	return equal(reflect.ValueOf(x), reflect.ValueOf(y), identical)
}

// EqualSyntax reports whether x and y are equal.
// Identifiers are considered equal if they are spelled the same.
// Comments are ignored.
func EqualSyntax(x, y ast.Expr) bool {
	sameName := func(x, y *ast.Ident) bool { return x.Name == y.Name }
	return Equal(x, y, sameName)
}

func equal(x, y reflect.Value, identical func(x, y *ast.Ident) bool) bool {
	// Ensure types are the same
	if x.Type() != y.Type() {
		return false
	}
	switch x.Kind() {
	case reflect.Pointer:
		// Ignore metadata even when only one side has it, as with a field comment.
		switch x.Interface().(type) {
		case *ast.Object, *ast.Scope, *ast.CommentGroup:
			return true
		}
		if x.IsNil() || y.IsNil() {
			return x.IsNil() == y.IsNil()
		}
		switch t := x.Interface().(type) {
		case *ast.Ident:
			return identical(t, y.Interface().(*ast.Ident))
		// These positions also mark syntax, such as `type T = int` and `f(xs...)`.
		case *ast.TypeSpec:
			return t.Assign.IsValid() == y.Interface().(*ast.TypeSpec).Assign.IsValid() &&
				equal(x.Elem(), y.Elem(), identical)
		case *ast.CallExpr:
			return t.Ellipsis.IsValid() == y.Interface().(*ast.CallExpr).Ellipsis.IsValid() &&
				equal(x.Elem(), y.Elem(), identical)
		default:
			return equal(x.Elem(), y.Elem(), identical)
		}

	case reflect.Interface:
		if x.IsNil() || y.IsNil() {
			return x.IsNil() == y.IsNil()
		}
		return equal(x.Elem(), y.Elem(), identical)

	case reflect.Struct:
		for i := range x.NumField() {
			xf := x.Field(i)
			yf := y.Field(i)
			if _, ok := xf.Interface().(token.Pos); ok {
				continue
			}
			if !equal(xf, yf, identical) {
				return false
			}
		}
		return true

	case reflect.Slice:
		if x.IsNil() || y.IsNil() {
			return x.IsNil() == y.IsNil()
		}
		if x.Len() != y.Len() {
			return false
		}
		for i := range x.Len() {
			if !equal(x.Index(i), y.Index(i), identical) {
				return false
			}
		}
		return true

	case reflect.String:
		return x.String() == y.String()

	case reflect.Bool:
		return x.Bool() == y.Bool()

	case reflect.Int:
		return x.Int() == y.Int()

	default:
		panic(x)
	}
}
