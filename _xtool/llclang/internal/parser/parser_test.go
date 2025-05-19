package parser_test

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"unsafe"

	"github.com/goplus/lib/c"
	"github.com/goplus/lib/c/clang"
	"github.com/goplus/llcppg/_xtool/llclang/internal/parser"
	"github.com/goplus/llcppg/ast"

	"github.com/goplus/llpkg/cjson"
)

func TestParser(t *testing.T) {
	cases := []string{"class", "comment", "enum", "func", "scope", "struct", "typedef", "union", "macro", "forwarddecl1", "forwarddecl2"}
	// https://github.com/goplus/llgo/issues/1114
	// todo(zzy):use os.ReadDir
	for _, folder := range cases {
		testFrom(t, filepath.Join("testdata", folder), "temp.h", false)
	}
}

func testFrom(t *testing.T, dir string, filename string, gen bool) {
	var expect string
	var err error
	if !gen {
		json, err := os.ReadFile(filepath.Join(dir, "expect.json"))
		if err != nil {
			t.Fatal("ReadExpectFile failed:", err)
		}
		expect = string(json)
	}
	ast, err := parser.Do(&parser.ConverterConfig{
		File:  filepath.Join(dir, filename),
		IsCpp: true,
		Args:  []string{"-fparse-all-comments"},
	})
	if err != nil {
		t.Fatal("Do failed:", err)
	}
	// https://github.com/goplus/llgo/issues/1116
	// astJson, err := json.MarshalIndent(ast, "", "  ")
	// todo(zzy):use json.Marshal
	if err != nil {
		t.Fatal("MarshalIndent failed:", err)
	}
	json := parser.MarshalASTFile(ast)
	output := json.Print()
	actual := c.GoString(output)
	defer cjson.FreeCStr(unsafe.Pointer(output))
	defer json.Delete()

	if gen {
		err = os.WriteFile(filepath.Join(dir, "expect.json"), []byte(actual), os.ModePerm)
		if err != nil {
			t.Fatal("WriteFile failed:", err)
		}
	} else if expect != actual {
		t.Fatalf("expect %s, got %s", expect, actual)
	}
}

func TestNonBuiltinTypes(t *testing.T) {
	tests := []struct {
		TypeCode      string
		ExpectTypeStr string
		expr          ast.Expr
	}{
		{
			TypeCode:      "int*",
			ExpectTypeStr: "int *",
			expr: &ast.PointerType{
				X: &ast.BuiltinType{
					Kind: ast.Int,
				},
			},
		},
		{
			TypeCode:      "int***",
			ExpectTypeStr: "int ***",
			expr: &ast.PointerType{
				X: &ast.PointerType{
					X: &ast.PointerType{
						X: &ast.BuiltinType{Kind: ast.Int},
					},
				},
			},
		},
		{
			TypeCode:      "int[]",
			ExpectTypeStr: "int[]",
			expr: &ast.ArrayType{
				Elt: &ast.BuiltinType{Kind: ast.Int},
			},
		},
		{
			TypeCode:      "int[10]",
			ExpectTypeStr: "int[10]",
			expr: &ast.ArrayType{
				Elt: &ast.BuiltinType{Kind: ast.Int},
				Len: &ast.BasicLit{
					Kind:  ast.IntLit,
					Value: "10",
				},
			},
		},
		{
			TypeCode:      "int[3][4]",
			ExpectTypeStr: "int[3][4]",
			expr: &ast.ArrayType{
				Elt: &ast.ArrayType{
					Elt: &ast.BuiltinType{Kind: ast.Int},
					Len: &ast.BasicLit{
						Kind:  ast.IntLit,
						Value: "4",
					},
				},
				Len: &ast.BasicLit{
					Kind:  ast.IntLit,
					Value: "3",
				},
			},
		},
		{
			TypeCode:      "int&",
			ExpectTypeStr: "int &",
			expr: &ast.LvalueRefType{
				X: &ast.BuiltinType{Kind: ast.Int},
			},
		},
		{
			TypeCode:      "int&&",
			ExpectTypeStr: "int &&",
			expr: &ast.RvalueRefType{
				X: &ast.BuiltinType{Kind: ast.Int},
			},
		},
		{
			TypeCode: `struct Foo {};
		               Foo`,
			ExpectTypeStr: "Foo",
			expr: &ast.Ident{
				Name: "Foo",
			},
		},
		{
			TypeCode: `struct Foo {};
		               struct Foo`,
			ExpectTypeStr: "struct Foo",
			expr: &ast.TagExpr{
				Tag: ast.Struct,
				Name: &ast.Ident{
					Name: "Foo",
				},
			},
		},
		{
			TypeCode: `struct {
						 int x;
					   }`,
			ExpectTypeStr: "struct (unnamed struct at temp.h:1:1)",
			expr: &ast.RecordType{
				Tag: ast.Struct,
				Fields: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								{Name: "x"},
							},
							Type:   &ast.BuiltinType{Kind: ast.Int},
							Access: ast.Public,
						},
					},
				},
				Methods: []*ast.FuncDecl{},
			},
		},
		{
			TypeCode: `union Foo {};
		               Foo`,
			ExpectTypeStr: "Foo",
			expr: &ast.Ident{
				Name: "Foo",
			},
		},
		{
			TypeCode: `union Foo {};
		               union Foo`,
			ExpectTypeStr: "union Foo",
			expr: &ast.TagExpr{
				Tag: ast.Union,
				Name: &ast.Ident{
					Name: "Foo",
				},
			},
		},
		{
			TypeCode: `union {
						int x;
					   }`,
			ExpectTypeStr: "union (unnamed union at temp.h:1:1)",
			expr: &ast.RecordType{
				Tag: ast.Union,
				Fields: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								{Name: "x"},
							},
							Access: ast.Public,
							Type:   &ast.BuiltinType{Kind: ast.Int},
						},
					},
				},
				Methods: []*ast.FuncDecl{},
			},
		},
		{
			TypeCode: `enum Foo {};
		               Foo`,
			ExpectTypeStr: "Foo",
			expr: &ast.Ident{
				Name: "Foo",
			},
		},
		{
			TypeCode: `enum Foo {};
		               enum Foo`,
			ExpectTypeStr: "enum Foo",
			expr: &ast.TagExpr{
				Tag: ast.Enum,
				Name: &ast.Ident{
					Name: "Foo",
				},
			},
		},
		{
			TypeCode:      `enum { x = 42 }`,
			ExpectTypeStr: "enum (unnamed enum at temp.h:1:1)",
			expr: &ast.EnumType{
				Items: []*ast.EnumItem{
					{
						Name: &ast.Ident{
							Name: "x",
						},
						Value: &ast.BasicLit{
							Kind:  ast.IntLit,
							Value: "42",
						},
					},
				},
			},
		},
		{
			TypeCode: `class Foo {};
		               Foo`,
			ExpectTypeStr: "Foo",
			expr: &ast.Ident{
				Name: "Foo",
			},
		},
		{
			TypeCode: `class Foo {};
		               class Foo`,
			ExpectTypeStr: "class Foo",
			expr: &ast.TagExpr{
				Tag: ast.Class,
				Name: &ast.Ident{
					Name: "Foo",
				},
			},
		},
		{
			TypeCode: `class {
						int x;
					   }`,
			ExpectTypeStr: "class (unnamed class at temp.h:1:1)",
			expr: &ast.RecordType{
				Tag: ast.Class,
				Fields: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								{Name: "x"},
							},
							Access: ast.Private,
							Type:   &ast.BuiltinType{Kind: ast.Int},
						},
					},
				},
				Methods: []*ast.FuncDecl{},
			},
		},
		{
			TypeCode: `namespace a {
						 namespace b {
						   class c {
						   };
						 }
					   }
					   a::b::c`,
			ExpectTypeStr: "a::b::c",
			expr: &ast.ScopingExpr{
				Parent: &ast.ScopingExpr{
					Parent: &ast.Ident{
						Name: "a",
					},
					X: &ast.Ident{
						Name: "b",
					},
				},
				X: &ast.Ident{
					Name: "c",
				},
			},
		},
		{
			TypeCode: `namespace a {
					     namespace b {
						   class c {
						   };
						 }
					   }
					   class a::b::c`,
			ExpectTypeStr: "class a::b::c",
			expr: &ast.TagExpr{
				Tag: ast.Class,
				Name: &ast.ScopingExpr{
					Parent: &ast.ScopingExpr{
						Parent: &ast.Ident{
							Name: "a",
						},
						X: &ast.Ident{
							Name: "b",
						},
					},
					X: &ast.Ident{
						Name: "c",
					},
				},
			},
		},
		{
			TypeCode:      `int (*p)(int, int);`,
			ExpectTypeStr: "int (*)(int, int)",
			expr: &ast.PointerType{
				X: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.BuiltinType{Kind: ast.Int},
							},
							{
								Type: &ast.BuiltinType{Kind: ast.Int},
							},
						},
					},
					Ret: &ast.BuiltinType{Kind: ast.Int},
				},
			},
		},
	}

	for _, tc := range tests {
		typ, index, unit := GetType(&GetTypeOptions{
			TypeCode: tc.TypeCode,
			IsCpp:    true,
		})
		converter := &parser.Converter{}
		expr := converter.ProcessType(typ)
		json := parser.MarshalASTExpr(expr)
		str := json.Print()
		typstr := typ.String()
		if typGoStr := c.GoString(typstr.CStr()); typGoStr != tc.ExpectTypeStr {
			t.Fatalf("expect %s , got %s", tc.ExpectTypeStr, typGoStr)
		}
		if !reflect.DeepEqual(expr, tc.expr) {
			t.Fatalf("%s expect %#v, got %#v", tc.ExpectTypeStr, tc.expr, expr)
		}

		typstr.Dispose()
		cjson.FreeCStr(unsafe.Pointer(str))
		json.Delete()
		index.Dispose()
		unit.Dispose()
	}
}

type GetTypeOptions struct {
	TypeCode string // e.g. "char*", "char**"

	// ExpectTypeKind specifies the expected type kind (optional)
	// Use clang.Type_Invalid to accept any type (default behavior)
	// *For complex types (when <complex.h> is included), specifying this is crucial
	// to filter out the correct type, as there will be multiple VarDecl fields present
	ExpectTypeKind clang.TypeKind

	// Args contains additional compilation arguments passed to Clang (optional)
	// These are appended after the default language-specific arguments
	// Example: []string{"-std=c++11"}
	Args []string

	// IsCpp indicates whether the code should be treated as C++ (true) or C (false)
	// This affects the default language arguments passed to Clang:
	// - For C++: []string{"-x", "c++"}
	// - For C:   []string{"-x", "c"}
	// *For complex C types, C Must be specified
	IsCpp bool
}

// GetType returns the clang.Type of the given type code
// Need to dispose the index and unit after using
// e.g. index.Dispose(), unit.Dispose()
func GetType(option *GetTypeOptions) (clang.Type, *clang.Index, *clang.TranslationUnit) {
	code := fmt.Sprintf("%s placeholder;", option.TypeCode)
	index, unit, err := parser.CreateTranslationUnit(&parser.LibClangConfig{
		File:  code,
		Temp:  true,
		Args:  option.Args,
		IsCpp: option.IsCpp,
	})
	if err != nil {
		panic(err)
	}
	cursor := unit.Cursor()
	var typ clang.Type
	parser.VisitChildren(cursor, func(child, parent clang.Cursor) clang.ChildVisitResult {
		if child.Kind == clang.CursorVarDecl && (option.ExpectTypeKind == clang.TypeInvalid || option.ExpectTypeKind == child.Type().Kind) {
			typ = child.Type()
			return clang.ChildVisit_Break
		}
		return clang.ChildVisit_Continue
	})
	return typ, index, unit
}
