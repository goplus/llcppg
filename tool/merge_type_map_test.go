// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in tool/internal/astutil/LICENSE.

package tool

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

var (
	tPStr1    = &ast.StarExpr{X: ast.NewIdent("string")}
	tPStr2    = &ast.StarExpr{X: ast.NewIdent("string")}
	tChanInt1 = &ast.ChanType{Dir: ast.RECV, Value: ast.NewIdent("int")}
	tChanInt2 = &ast.ChanType{Dir: ast.RECV, Value: ast.NewIdent("int")}
)

func checkEqualButNotIdentical(t *testing.T, x, y ast.Node, comment string) {
	if !equalASTKey(x, y) {
		t.Errorf("%s: not equal: %v, %v", comment, x, y)
	}
	if x == y {
		t.Errorf("%s: identical: %v, %v", comment, x, y)
	}
}

func TestASTMapAxioms(t *testing.T) {
	checkEqualButNotIdentical(t, tPStr1, tPStr2, "tPstr{1,2}")
	checkEqualButNotIdentical(t, tChanInt1, tChanInt2, "tChanInt{1,2}")
}

func TestASTMap(t *testing.T) {
	var tmap *astMap

	// All methods but Set are safe on (*T)(nil).
	tmap.len()
	tmap.at(tPStr1)
	tmap.delete(tPStr1)
	tmap.keysString()
	_ = tmap.string()

	tmap = new(astMap)

	// Length of empty map.
	if l := tmap.len(); l != 0 {
		t.Errorf("Len() on empty Map: got %d, want 0", l)
	}
	// At of missing key.
	if v := tmap.at(tPStr1); v != nil {
		t.Errorf("At() on empty Map: got %v, want nil", v)
	}
	// Deletion of missing key.
	if tmap.delete(tPStr1) {
		t.Errorf("Delete() on empty Map: got true, want false")
	}
	// Set of new key.
	if prev := tmap.set(tPStr1, "*string"); prev != nil {
		t.Errorf("Set() on empty Map returned non-nil previous value %s", prev)
	}

	// Now: {*string: "*string"}

	// Length of non-empty map.
	if l := tmap.len(); l != 1 {
		t.Errorf("Len(): got %d, want 1", l)
	}
	// At via insertion key.
	if v := tmap.at(tPStr1); v != "*string" {
		t.Errorf("At(): got %q, want \"*string\"", v)
	}
	// At via equal key.
	if v := tmap.at(tPStr2); v != "*string" {
		t.Errorf("At(): got %q, want \"*string\"", v)
	}
	// Iteration over sole entry.
	tmap.iterate(func(key ast.Node, value any) {
		if key != tPStr1 {
			t.Errorf("Iterate: key: got %v, want %v", key, tPStr1)
		}
		if want := "*string"; value != want {
			t.Errorf("Iterate: value: got %s, want %s", value, want)
		}
	})

	// Set with key equal to present one.
	if prev := tmap.set(tPStr2, "*string again"); prev != "*string" {
		t.Errorf("Set() previous value: got %s, want \"*string\"", prev)
	}

	// Set of another association.
	if prev := tmap.set(tChanInt1, "<-chan int"); prev != nil {
		t.Errorf("Set() previous value: got %s, want nil", prev)
	}

	// Now: {*string: "*string again", <-chan int: "<-chan int"}

	want1 := "{*string: \"*string again\", <-chan int: \"<-chan int\"}"
	want2 := "{<-chan int: \"<-chan int\", *string: \"*string again\"}"
	if s := tmap.string(); s != want1 && s != want2 {
		t.Errorf("String(): got %s, want %s", s, want1)
	}

	want1 = "{*string, <-chan int}"
	want2 = "{<-chan int, *string}"
	if s := tmap.keysString(); s != want1 && s != want2 {
		t.Errorf("KeysString(): got %s, want %s", s, want1)
	}

	// Keys().
	I := equalASTKey
	switch k := tmap.keys(); {
	case I(k[0], tChanInt1) && I(k[1], tPStr1): // ok
	case I(k[1], tChanInt1) && I(k[0], tPStr1): // ok
	default:
		t.Errorf("Keys(): got %v, want %s", k, want2)
	}

	if l := tmap.len(); l != 2 {
		t.Errorf("Len(): got %d, want 2", l)
	}
	// At via original key.
	if v := tmap.at(tPStr1); v != "*string again" {
		t.Errorf("At(): got %q, want \"*string again\"", v)
	}
	hamming := 1
	tmap.iterate(func(key ast.Node, value any) {
		switch {
		case I(key, tChanInt1):
			hamming *= 2 // ok
		case I(key, tPStr1):
			hamming *= 3 // ok
		}
	})
	if hamming != 6 {
		t.Errorf("Iterate: hamming: got %d, want %d", hamming, 6)
	}

	if v := tmap.at(tChanInt2); v != "<-chan int" {
		t.Errorf("At(): got %q, want \"<-chan int\"", v)
	}
	// Deletion with key equal to present one.
	if !tmap.delete(tChanInt2) {
		t.Errorf("Delete() of existing key: got false, want true")
	}

	// Now: {*string: "*string again"}

	if l := tmap.len(); l != 1 {
		t.Errorf("Len(): got %d, want 1", l)
	}
	// Deletion again.
	if !tmap.delete(tPStr2) {
		t.Errorf("Delete() of existing key: got false, want true")
	}

	// Now: {}

	if l := tmap.len(); l != 0 {
		t.Errorf("Len(): got %d, want %d", l, 0)
	}
	if s := tmap.string(); s != "{}" {
		t.Errorf("String(): got %q, want %q", s, "{}")
	}
}

func TestASTMapGenerics(t *testing.T) {
	const src = `
package p

// Basic defined types.
type T1 int
type T2 int

// Identical methods.
func (T1) M(int) {}
func (T2) M(int) {}

// A constraint interface.
type C interface {
	~int | string
}

type I interface {
}

// A generic type.
type G[P C] int

// Generic functions with different type parameter names.
func Fa1[P C](p P) {}
func Fa2[Q C](q Q) {}

// Map the parameter and local variable in Fb1 before the function type.
func Fb1[P any](x *P) {
	var y *P // Map this first.
	_ = y
}
func Fb2[Q any](x *Q) {
}

// G1 and G2 are mutually recursive and have identical method types.
type G1[P any] struct{
	Field *G2[P]
}
func (G1[P]) M(G1[P], G2[P]) {}
type G2[Q any] struct{
	Field *G1[Q]
}
func (G2[P]) M(G1[P], G2[P]) {}

// Method type expressions on different generic types are different.
var ME1 = G1[int].M
var ME2 = G2[int].M

// ME1Type spells out the type of the method expression ME1.
var ME1Type func(G1[int], G1[int], G2[int])

// Examples from issue #51314
type Constraint[T any] any
func Foo[T Constraint[T]]() {}
func Fn[T1 ~*T2, T2 ~*T1](t1 T1, t2 T2) {}

// Bar and Baz have different syntax for the constraints.
func Bar[P Constraint[P]]() {}
func Baz[Q any]() {} // The underlying type of Constraint[P] is any.
func Quux[Q interface{ quux() }]() {}


type Issue56048_I interface{ m() interface { Issue56048_I } }
var Issue56048 = Issue56048_I.m

type Issue56048_Ib interface{ m() chan []*interface { Issue56048_Ib } }
var Issue56048b = Issue56048_Ib.m

// Non-generic alias
type NonAlias int
type Alias1 = NonAlias
type Alias2 = NonAlias

// Generic alias (requires go1.23)
// type SetOfInt = map[int]bool
// type Set[T comparable] = map[K]bool
// type SetOfInt2 = Set[int]
`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "p.go", src, 0)
	if err != nil {
		t.Fatal(err)
	}

	typeSpecs := make(map[string]*ast.TypeSpec)
	funcs := make(map[string]*ast.FuncDecl)
	vars := make(map[string]*ast.ValueSpec)
	for _, decl := range file.Decls {
		switch decl := decl.(type) {
		case *ast.GenDecl:
			for _, spec := range decl.Specs {
				switch spec := spec.(type) {
				case *ast.TypeSpec:
					typeSpecs[spec.Name.Name] = spec
				case *ast.ValueSpec:
					vars[spec.Names[0].Name] = spec
				}
			}
		case *ast.FuncDecl:
			name := decl.Name.Name
			if decl.Recv != nil {
				recv := decl.Recv.List[0].Type
				if instance, ok := recv.(*ast.IndexExpr); ok {
					recv = instance.X
				}
				name = recv.(*ast.Ident).Name + "." + name
			}
			funcs[name] = decl
		}
	}

	expr := func(src string) ast.Expr {
		t.Helper()
		node, err := parser.ParseExpr(src)
		if err != nil {
			t.Fatal(err)
		}
		return node
	}

	ci := typeSpecs["C"].Type.(*ast.InterfaceType)
	fa1, fa2 := funcs["Fa1"].Type, funcs["Fa2"].Type
	fb1, fb2 := funcs["Fb1"], funcs["Fb2"]
	fb1y := fb1.Body.List[0].(*ast.DeclStmt).Decl.(*ast.GenDecl).Specs[0].(*ast.ValueSpec).Type

	tmap := new(astMap)
	steps := []struct {
		node     ast.Node
		name     string
		newEntry bool
	}{
		{typeSpecs["T1"], "T1", true},
		{typeSpecs["T2"], "T2", true},
		{typeSpecs["G"], "G", true},
		{typeSpecs["C"], "C", true},
		{ci, "CI", true},
		{ci.Methods.List[0].Type, "U", true},
		{typeSpecs["I"], "I", true},
		{typeSpecs["I"].Type, "II", true},

		// Method type nodes leave out the receiver, e.g. T1.M and T2.M.
		{funcs["T1.M"].Type, "T1M", true},
		{funcs["T2.M"].Type, "T2M", false},

		// Identical instances share an entry, e.g. separately parsed G[int].
		{expr("G[int]"), "GInt1", true},
		{expr("G[int]"), "GInt2", false},
		{expr("G[string]"), "GStr", true},

		// Parameter names are structural, so P and Q distinguish Fa1 and Fa2.
		{fa1, "Fa1", true},
		{fa2, "Fa2", true},
		{fa1.TypeParams.List[0].Names[0], "Fa1P", true},
		{fa2.TypeParams.List[0].Names[0], "Fa2Q", true},

		{fb1y, "Fb1y", true},
		{fb1.Type.Params.List[0].Type, "Fb1x", false},
		{fb2.Type.Params.List[0].Type, "Fb2x", true},
		{fb1.Type, "Fb1", true},
		{fb2.Type, "Fb2", true},

		{typeSpecs["G1"], "G1", true},
		{funcs["G1.M"].Type, "G1M", true},
		{typeSpecs["G2"], "G2", true},
		{expr("func(G1[int], G2[int])"), "G1IntM1", true},
		{expr("func(G1[int], G2[int])"), "G1IntM2", false},
		{expr("func(G1[string], G2[string])"), "G1StrM", true},
		{expr("func(G1[int], G2[int])"), "G2IntM", false},

		// A selector G1[int].M and its declared signature have different ASTs.
		{vars["ME1"].Values[0], "ME1", true},
		{vars["ME1Type"].Type, "ME1Type", true},
		{vars["ME2"].Values[0], "ME2", true},

		// Recursive type constraints such as T Constraint[T] are finite syntax.
		{typeSpecs["Constraint"], "Constraint", true},
		{funcs["Foo"].Type, "Foo", true},
		{funcs["Fn"].Type, "Fn", true},
		{funcs["Bar"].Type, "Bar", true},
		{funcs["Baz"].Type, "Baz", true},
		{funcs["Quux"].Type, "Quux", true},

		{typeSpecs["Issue56048_I"].Type, "Issue56048", true},
		{typeSpecs["Issue56048_Ib"].Type, "Issue56048b", true},

		// A type definition and aliases such as type Alias1 = NonAlias have distinct syntax.
		{typeSpecs["NonAlias"], "NonAlias", true},
		{typeSpecs["Alias1"], "Alias1", true},
		{typeSpecs["Alias2"], "Alias2", true},
	}

	for _, step := range steps {
		existing := tmap.at(step.node)
		if (existing == nil) != step.newEntry {
			t.Errorf("at(%s) = %v, want new entry: %t", step.name, existing, step.newEntry)
		}
		tmap.set(step.node, step.name)
	}
}
