package parser_test

import (
	"testing"

	"github.com/xgo-dev/llcppg/ast"
	"github.com/xgo-dev/llcppg/parser"
)

func TestParseFile(t *testing.T) {
	astFile, err := parser.ParseFile(nil, "./testdata/func/hfile/forwarddecl.h", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	if astFile == nil {
		t.Fatal("astFile is nil")
	}
	if len(astFile.Decls) == 0 {
		t.Fatal("expect decls, but got 0")
	}
	expectNames := []string{"foo0", "foo1", "foo2"}
	for i, decl := range astFile.Decls {
		if _, ok := decl.(*ast.FuncDecl); !ok {
			t.Fatalf("expect FuncDecl, but got %T", decl)
		}
		if decl.(*ast.FuncDecl).Name.Name != expectNames[i] {
			t.Fatalf("expect %s, but got %s", expectNames[i], decl.(*ast.FuncDecl).Name.Name)
		}
	}
}
