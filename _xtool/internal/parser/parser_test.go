package parser_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/goplus/lib/c"
	clangutils "github.com/goplus/llcppg/_xtool/internal/clang"
	"github.com/goplus/llcppg/_xtool/internal/clangtool"
	clang "github.com/goplus/llcppg/_xtool/internal/libclang"
	"github.com/goplus/llcppg/_xtool/internal/parser"
	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llgo/xtool/clang/preprocessor"
)

func TestParserCppMode(t *testing.T) {
	cases := []string{"class", "comment", "enum", "func", "scope", "struct", "typedef", "union", "macro", "forwarddecl1", "forwarddecl2", "include", "typeof", "forward_vs_empty", "nestedenum_cpp"}
	// https://github.com/goplus/llgo/issues/1114
	// todo(zzy):use os.ReadDir
	for _, folder := range cases {
		t.Run(folder, func(t *testing.T) {
			testFrom(t, filepath.Join("testdata", folder), "temp.h", true, false)
		})
	}
}

func TestParserCMode(t *testing.T) {
	cases := []string{"enum", "struct", "union", "macro", "include", "typeof", "named_nested_struct", "forward_vs_empty", "nestedenum"}
	for _, folder := range cases {
		t.Run(folder, func(t *testing.T) {
			testFrom(t, filepath.Join("testdata", folder), "temp.h", false, false)
		})
	}
}

func testFrom(t *testing.T, dir string, filename string, isCpp, gen bool) {
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
		IsCpp: isCpp,
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
	js := parser.XMarshalASTFile(ast)
	output, _ := json.MarshalIndent(&js, "", "  ")

	if gen {
		err = os.WriteFile(filepath.Join(dir, "expect.json"), output, os.ModePerm)
		if err != nil {
			t.Fatal("WriteFile failed:", err)
		}
	} else if expect != string(output) {
		t.Fatalf("expect %s, got %s", expect, string(output))
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
			TypeCode: `struct Foo { enum Bar {} k; };
						enum Bar`,
			ExpectTypeStr: "enum Bar",
			expr: &ast.TagExpr{
				Tag:  ast.Enum,
				Name: &ast.Ident{Name: "Bar"},
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
		t.Run(tc.ExpectTypeStr, func(t *testing.T) {
			typ, index, unit := GetType(&GetTypeOptions{
				TypeCode: tc.TypeCode,
				IsCpp:    true,
			})
			converter := &parser.Converter{}
			expr := converter.ProcessType(typ)
			typstr := typ.String()
			if typGoStr := c.GoString(typstr.CStr()); typGoStr != tc.ExpectTypeStr {
				t.Fatalf("expect %s , got %s", tc.ExpectTypeStr, typGoStr)
			}
			if !reflect.DeepEqual(expr, tc.expr) {
				t.Fatalf("%s expect %#v, got %#v", tc.ExpectTypeStr, tc.expr, expr)
			}

			typstr.Dispose()

			index.Dispose()
			unit.Dispose()
		})
	}
}

func TestBuiltinType(t *testing.T) {
	tests := []struct {
		name     string
		typ      clang.Type
		expected ast.BuiltinType
	}{
		{"Void", btType(clang.TypeVoid), ast.BuiltinType{Kind: ast.Void}},
		{"Bool", btType(clang.TypeBool), ast.BuiltinType{Kind: ast.Bool}},
		{"Char_S", btType(clang.TypeCharS), ast.BuiltinType{Kind: ast.Char, Flags: ast.Signed}},
		{"Char_U", btType(clang.TypeCharU), ast.BuiltinType{Kind: ast.Char, Flags: ast.Unsigned}},
		{"Char16", btType(clang.TypeChar16), ast.BuiltinType{Kind: ast.Char16}},
		{"Char32", btType(clang.TypeChar32), ast.BuiltinType{Kind: ast.Char32}},
		{"WChar", btType(clang.TypeWChar), ast.BuiltinType{Kind: ast.WChar}},
		{"Short", btType(clang.TypeShort), ast.BuiltinType{Kind: ast.Int, Flags: ast.Short}},
		{"UShort", btType(clang.TypeUShort), ast.BuiltinType{Kind: ast.Int, Flags: ast.Short | ast.Unsigned}},
		{"Int", btType(clang.TypeInt), ast.BuiltinType{Kind: ast.Int}},
		{"UInt", btType(clang.TypeUInt), ast.BuiltinType{Kind: ast.Int, Flags: ast.Unsigned}},
		{"Long", btType(clang.TypeLong), ast.BuiltinType{Kind: ast.Int, Flags: ast.Long}},
		{"ULong", btType(clang.TypeULong), ast.BuiltinType{Kind: ast.Int, Flags: ast.Long | ast.Unsigned}},
		{"LongLong", btType(clang.TypeLongLong), ast.BuiltinType{Kind: ast.Int, Flags: ast.LongLong}},
		{"ULongLong", btType(clang.TypeULongLong), ast.BuiltinType{Kind: ast.Int, Flags: ast.LongLong | ast.Unsigned}},
		{"Int128", btType(clang.TypeInt128), ast.BuiltinType{Kind: ast.Int128}},
		{"UInt128", btType(clang.TypeUInt128), ast.BuiltinType{Kind: ast.Int128, Flags: ast.Unsigned}},
		{"Float", btType(clang.TypeFloat), ast.BuiltinType{Kind: ast.Float}},
		{"Half", btType(clang.TypeHalf), ast.BuiltinType{Kind: ast.Float16}},
		{"Float16", btType(clang.TypeFloat16), ast.BuiltinType{Kind: ast.Float16}},
		{"Double", btType(clang.TypeDouble), ast.BuiltinType{Kind: ast.Float, Flags: ast.Double}},
		{"LongDouble", btType(clang.TypeLongDouble), ast.BuiltinType{Kind: ast.Float, Flags: ast.Long | ast.Double}},
		{"Float128", btType(clang.TypeFloat128), ast.BuiltinType{Kind: ast.Float128}},
		{"Complex", getComplexType(0), ast.BuiltinType{Kind: ast.Complex}},
		{"Complex", getComplexType(ast.Double), ast.BuiltinType{Flags: ast.Double, Kind: ast.Complex}},
		{"Complex", getComplexType(ast.Long | ast.Double), ast.BuiltinType{Flags: ast.Long | ast.Double, Kind: ast.Complex}},
		{"Unknown", btType(clang.TypeIbm128), ast.BuiltinType{Kind: ast.Void}},
	}

	converter := &parser.Converter{}
	converter.Convert()
	for _, bt := range tests {
		res := converter.ProcessBuiltinType(bt.typ)
		if res.Kind != bt.expected.Kind {
			t.Fatalf("%s Kind mismatch:got %d want %d, \n", bt.name, res.Kind, bt.expected.Kind)
		}
		if res.Flags != bt.expected.Flags {
			t.Fatalf("%s Flags mismatch:got %d,want %d\n", bt.name, res.Flags, bt.expected.Flags)
		}
	}
}

// Char's Default Type in macos is signed char & in linux is unsigned char
// So we only confirm the char's kind is char & flags is unsigned or signed
func TestChar(t *testing.T) {
	typ, index, transunit := GetType(&GetTypeOptions{
		TypeCode: "char",
		IsCpp:    false,
	})
	converter := &parser.Converter{}
	expr := converter.ProcessType(typ)
	if btType, ok := expr.(*ast.BuiltinType); ok {
		if btType.Kind == ast.Char {
			if btType.Flags != ast.Signed && btType.Flags != ast.Unsigned {
				t.Fatal("Char's flags is not signed or unsigned")
			}
		}
	} else {
		t.Fatal("Char's expr is not a builtin type")
	}
	index.Dispose()
	transunit.Dispose()
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
	index, unit, err := clangutils.CreateTranslationUnit(&clangutils.Config{
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
	clangutils.VisitChildren(cursor, func(child, parent clang.Cursor) clang.ChildVisitResult {
		if child.Kind == clang.CursorVarDecl && (option.ExpectTypeKind == clang.TypeInvalid || option.ExpectTypeKind == child.Type().Kind) {
			typ = child.Type()
			return clang.ChildVisit_Break
		}
		return clang.ChildVisit_Continue
	})
	return typ, index, unit
}

func btType(kind clang.TypeKind) clang.Type {
	return clang.Type{Kind: kind}
}

// get complex type from source code parsed
func getComplexType(flag ast.TypeFlag) clang.Type {
	var typeStr string
	if flag&(ast.Long|ast.Double) == (ast.Long | ast.Double) {
		typeStr = "long double"
	} else if flag&ast.Double != 0 {
		typeStr = "double"
	} else {
		typeStr = "float"
	}

	code := fmt.Sprintf("#include <complex.h>\n%s complex", typeStr)

	// todo(zzy):free index and unit after test
	typ, _, _ := GetType(&GetTypeOptions{
		TypeCode:       code,
		ExpectTypeKind: clang.TypeComplex,
		IsCpp:          false,
	})

	return typ
}

func TestPreprocess(t *testing.T) {
	combinedFile, err := os.CreateTemp("./", "compose_*.h")
	if err != nil {
		panic(err)
	}
	defer os.Remove(combinedFile.Name())

	clangtool.ComposeIncludes([]string{"main.h", "compat.h"}, combinedFile.Name())

	efile, err := os.CreateTemp("", "temp_*.i")
	if err != nil {
		panic(err)
	}
	defer os.Remove(efile.Name())

	ppconf := &preprocessor.Config{
		Compiler: "clang",
		Flags:    []string{"-I./_testdata/hfile"},
	}
	err = preprocessor.Do(combinedFile.Name(), efile.Name(), ppconf)
	if err != nil {
		t.Fatal(err)
	}

	config := &clangutils.Config{
		File:  efile.Name(),
		Temp:  false,
		IsCpp: false,
	}

	var str strings.Builder

	visit(config, func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		switch cursor.Kind {
		case clang.CursorEnumDecl, clang.CursorStructDecl, clang.CursorUnionDecl, clang.CursorTypedefDecl:
			var filename clang.String
			var line, column c.Uint
			cursor.Location().PresumedLocation(&filename, &line, &column)
			str.WriteString("TypeKind: ")
			str.WriteString(clang.GoString(cursor.Kind.String()))
			str.WriteString(" Name: ")
			str.WriteString(clang.GoString(cursor.String()))
			str.WriteString("\n")
			str.WriteString("Location: ")
			str.WriteString(fmt.Sprintf("%s:%d:%d\n", path.Base(c.GoString(filename.CStr())), line, column))
		}
		return clang.ChildVisit_Continue
	})

	expect := `
TypeKind: StructDecl Name: A
Location: main.h:3:16
TypeKind: TypedefDecl Name: A
Location: main.h:6:3
TypeKind: TypedefDecl Name: B
Location: compat.h:3:11
TypeKind: TypedefDecl Name: C
Location: main.h:8:11
`

	compareOutput(t, expect, str.String())
}

func visit(config *clangutils.Config, visitFunc func(cursor, parent clang.Cursor) clang.ChildVisitResult) {
	index, unit, err := clangutils.CreateTranslationUnit(config)
	if err != nil {
		panic(err)
	}
	cursor := unit.Cursor()
	clangutils.VisitChildren(cursor, visitFunc)
	index.Dispose()
	unit.Dispose()
}

func compareOutput(t *testing.T, expected, actual string) {
	expected = strings.TrimSpace(expected)
	actual = strings.TrimSpace(actual)
	if expected != actual {
		t.Fatalf("Test failed: expected \n%s \ngot \n%s", expected, actual)
	}
}

func TestPostOrderVisitChildren(t *testing.T) {
	config := &clangutils.Config{
		File:  "./testdata/named_nested_struct/temp.h",
		Temp:  false,
		IsCpp: false,
	}

	name := make(map[string]bool)
	visit(config, func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		if cursor.Kind == clang.CursorStructDecl {
			if !name[clang.GoString(cursor.String())] {
				name[clang.GoString(cursor.String())] = true
				file, line, column := clangutils.GetPresumedLocation(cursor.Location())
				fmt.Println("StructDecl Name:", clang.GoString(cursor.String()), file, line, column)
			}
		}
		return clang.ChildVisit_Recurse
	})

	index, unit, err := clangutils.CreateTranslationUnit(config)
	if err != nil {
		panic(err)
	}
	defer index.Dispose()
	defer unit.Dispose()

	childStr := make([]string, 6)
	childs := parser.PostOrderVisitChildren(unit.Cursor(), func(child, parent clang.Cursor) bool {
		return child.Kind == clang.CursorStructDecl
	})
	for i, child := range childs {
		childStr[i] = clang.GoString(child.String())
	}
	expect := []string{"c", "d", "b", "f", "e", "a"}
	if !reflect.DeepEqual(expect, childStr) {
		fmt.Println("Unexpected child order:", childStr)
	}
}

func TestEmptyDeclVsForwardDecl(t *testing.T) {
	config := &clangutils.Config{
		File:  "./testdata/forward_vs_empty/temp.h",
		Temp:  false,
		IsCpp: false,
	}

	type isDefinition = bool
	var decl map[string]isDefinition = map[string]isDefinition{
		"ForwardOnly": false,
		"EmptyStruct": true,
	}

	visit(config, func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		if cursor.Kind == clang.CursorStructDecl {
			sdecl := clang.GoString(cursor.String())
			if _, ok := decl[sdecl]; ok {
				isDefine := cursor.IsCursorDefinition() != 0
				if isDefine != decl[sdecl] {
					t.Fatalf("StructDecl %s isDefinition expect %v, got %v", sdecl, decl[sdecl], isDefine)
				}
			}
			fmt.Println("StructDecl Name:", clang.GoString(cursor.String()), "isDefinition:", cursor.IsCursorDefinition() != 0)
		}
		return clang.ChildVisit_Recurse
	})
}
