package convert_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/goplus/gogen"
	"github.com/goplus/llcppg/_xtool/llcppsymg/tool/name"
	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/cl/internal/convert"
	"github.com/goplus/llcppg/cmd/gogensig/config"
	llcppg "github.com/goplus/llcppg/config"
	"github.com/goplus/llcppg/token"
	"github.com/goplus/mod/gopmod"
)

var dir string

func init() {
	convert.SetDebug(convert.DbgFlagAll)
	var err error
	dir, err = os.Getwd()
	if err != nil {
		panic(err)
	}
}

func TestUnionDecl(t *testing.T) {
	testCases := []genDeclTestCase{
		/*
			union  u
			{
			    int a;
			    long b;
			    long c;
			    bool f;
			};
		*/
		{
			name: "union u{int a; long b; long c; bool f;};",
			decl: &ast.TypeDecl{
				Name: &ast.Ident{Name: "u"},
				Type: &ast.RecordType{
					Tag: ast.Union,
					Fields: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{Name: "a"},
								},
								Type: &ast.BuiltinType{
									Kind: ast.Int},
							},
							{
								Names: []*ast.Ident{
									{Name: "b"},
								},
								Type: &ast.BuiltinType{
									Kind:  ast.Int,
									Flags: ast.Long,
								},
							},
							{
								Names: []*ast.Ident{
									{Name: "c"},
								},
								Type: &ast.BuiltinType{
									Kind:  ast.Int,
									Flags: ast.Long,
								},
							},
							{
								Names: []*ast.Ident{
									{Name: "f"},
								},
								Type: &ast.BuiltinType{
									Kind: ast.Bool,
								},
							},
						},
					},
				},
			},
			expected: `
package testpkg

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type U struct {
	B c.Long
}
`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testGenDecl(t, tc)
		})
	}
}

func TestToType(t *testing.T) {
	pkg, err := createTestPkg(&convert.PackageConfig{
		OutputDir:  "",
		LibCommand: "${pkg-config --libs libcjson}",
	})
	if err != nil {
		t.Fatal("NewPackage failed:", err)
	}

	testCases := []struct {
		name     string
		input    *ast.BuiltinType
		expected string
	}{
		{"Void", &ast.BuiltinType{Kind: ast.Void}, "github.com/goplus/lib/c.Void"},
		{"Bool", &ast.BuiltinType{Kind: ast.Bool}, "bool"},
		{"Char_S", &ast.BuiltinType{Kind: ast.Char, Flags: ast.Signed}, "github.com/goplus/lib/c.Char"},
		{"Char_U", &ast.BuiltinType{Kind: ast.Char, Flags: ast.Unsigned}, "github.com/goplus/lib/c.Char"},
		{"WChar", &ast.BuiltinType{Kind: ast.WChar}, "int16"},
		{"Char16", &ast.BuiltinType{Kind: ast.Char16}, "int16"},
		{"Char32", &ast.BuiltinType{Kind: ast.Char32}, "int32"},
		{"Short", &ast.BuiltinType{Kind: ast.Int, Flags: ast.Short}, "int16"},
		{"UShort", &ast.BuiltinType{Kind: ast.Int, Flags: ast.Short | ast.Unsigned}, "uint16"},
		{"Int", &ast.BuiltinType{Kind: ast.Int}, "github.com/goplus/lib/c.Int"},
		{"UInt", &ast.BuiltinType{Kind: ast.Int, Flags: ast.Unsigned}, "github.com/goplus/lib/c.Uint"},
		{"Long", &ast.BuiltinType{Kind: ast.Int, Flags: ast.Long}, "github.com/goplus/lib/c.Long"},
		{"ULong", &ast.BuiltinType{Kind: ast.Int, Flags: ast.Long | ast.Unsigned}, "github.com/goplus/lib/c.Ulong"},
		{"LongLong", &ast.BuiltinType{Kind: ast.Int, Flags: ast.LongLong}, "github.com/goplus/lib/c.LongLong"},
		{"ULongLong", &ast.BuiltinType{Kind: ast.Int, Flags: ast.LongLong | ast.Unsigned}, "github.com/goplus/lib/c.UlongLong"},
		{"Float", &ast.BuiltinType{Kind: ast.Float}, "github.com/goplus/lib/c.Float"},
		{"Double", &ast.BuiltinType{Kind: ast.Float, Flags: ast.Double}, "github.com/goplus/lib/c.Double"},
		{"ComplexFloat", &ast.BuiltinType{Kind: ast.Complex}, "complex64"},
		{"ComplexDouble", &ast.BuiltinType{Kind: ast.Complex, Flags: ast.Double}, "complex128"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, _ := pkg.ToType(tc.input)
			if result != nil && result.String() != tc.expected {
				t.Errorf("unexpected result:%s expected:%s", result.String(), tc.expected)
			}
		})
	}
}

func TestToTypeFail(t *testing.T) {
	pkg, err := createTestPkg(&convert.PackageConfig{
		OutputDir: "",
	})
	if err != nil {
		t.Fatal("NewPackage failed:", err)
	}
	_, err = pkg.ToType(&ast.Comment{Text: "test"})
	if err == nil {
		t.Fatal("Expect error but got nil")
	}
}

var tempFile = &convert.HeaderFile{
	File:     "/path/to/temp.go",
	FileType: llcppg.Inter,
}

func TestNewPackage(t *testing.T) {
	pkg, err := createTestPkg(&convert.PackageConfig{})
	if err != nil {
		t.Fatal("NewPackage failed:", err)
	}
	pkg.SetCurFile(tempFile)
	comparePackageOutput(t, pkg, `
package testpkg

import _ "unsafe"
	`)
}

func TestPackageWrite(t *testing.T) {
	verifyGeneratedFile := func(t *testing.T, expectedFilePath string) {
		t.Helper()
		if _, err := os.Stat(expectedFilePath); os.IsNotExist(err) {
			t.Fatalf("Expected output file does not exist: %s", expectedFilePath)
		}

		content, err := os.ReadFile(expectedFilePath)
		if err != nil {
			t.Fatalf("Unable to read generated file: %v", err)
		}

		expectedContent := "package testpkg"
		if !strings.Contains(string(content), expectedContent) {
			t.Errorf("Generated file content does not match expected.\nExpected:\n%s\nActual:\n%s", expectedContent, string(content))
		}
	}

	incPath := "mock_header.h"
	filePath := filepath.Join("/path", "to", incPath)
	genPath := name.HeaderFileToGo(filePath)

	headerFile := convert.NewHeaderFile(filePath, llcppg.Inter)

	t.Run("OutputToTempDir", func(t *testing.T) {
		tempDir, err := os.MkdirTemp(dir, "test_package_write")
		if err != nil {
			t.Fatalf("Failed to create temporary directory: %v", err)
		}
		defer os.RemoveAll(tempDir)

		pkg, err := createTestPkg(&convert.PackageConfig{
			OutputDir: tempDir,
		})
		if err != nil {
			t.Fatal("NewPackage failed:", err)
		}

		pkg.SetCurFile(headerFile)
		err = pkg.Write(filePath)
		if err != nil {
			t.Fatalf("Write method failed: %v", err)
		}

		expectedFilePath := filepath.Join(tempDir, genPath)
		verifyGeneratedFile(t, expectedFilePath)
	})

	t.Run("OutputToCurrentDir", func(t *testing.T) {
		testpkgDir := filepath.Join(dir, "testpkg")
		if err := os.MkdirAll(testpkgDir, 0755); err != nil {
			t.Fatalf("Failed to create testpkg directory: %v", err)
		}

		defer func() {
			// Clean up generated files and directory
			os.RemoveAll(testpkgDir)
		}()

		pkg, err := createTestPkg(&convert.PackageConfig{
			OutputDir: testpkgDir,
		})
		if err != nil {
			t.Fatal("NewPackage failed:", err)
		}
		pkg.SetCurFile(headerFile)
		err = pkg.Write(filePath)
		if err != nil {
			t.Fatalf("Write method failed: %v", err)
		}

		expectedFilePath := filepath.Join(testpkgDir, genPath)
		verifyGeneratedFile(t, expectedFilePath)
	})

	t.Run("InvalidOutputDir", func(t *testing.T) {
		testpkgDir := filepath.Join(dir, "testpkg")
		if err := os.MkdirAll(testpkgDir, 0755); err != nil {
			t.Fatalf("Failed to create testpkg directory: %v", err)
		}
		defer func() {
			os.RemoveAll(testpkgDir)
		}()
		pkg, err := createTestPkg(&convert.PackageConfig{
			OutputDir: testpkgDir,
		})
		if err != nil {
			t.Fatal("NewPackage failed:", err)
		}
		pkg.Config().OutputDir = "/nonexistent/directory"
		err = pkg.Write(incPath)
		if err == nil {
			t.Fatal("Expected an error for invalid output directory, but got nil")
		}
	})

	t.Run("WriteUnexistFile", func(t *testing.T) {
		pkg, err := createTestPkg(&convert.PackageConfig{})
		if err != nil {
			t.Fatal("NewPackage failed:", err)
		}
		err = pkg.Write("test1.h")
		if err == nil {
			t.Fatal("Expected an error for invalid output directory, but got nil")
		}
	})
}

func TestFuncDecl(t *testing.T) {
	testCases := []genDeclTestCase{
		{
			name: "empty func",
			decl: &ast.FuncDecl{
				Name:        &ast.Ident{Name: "foo"},
				MangledName: "foo",
				Type: &ast.FuncType{
					Params: nil,
					Ret:    &ast.BuiltinType{Kind: ast.Void},
				},
			},
			symbs: []config.SymbolEntry{
				{
					CppName:    "foo",
					MangleName: "foo",
					GoName:     "Foo",
				},
			},
			expected: `
package testpkg

import _ "unsafe"
//go:linkname Foo C.foo
func Foo()
`,
		},
		{
			name: "variadic func",
			decl: &ast.FuncDecl{
				Name:        &ast.Ident{Name: "foo"},
				MangledName: "foo",
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{Type: &ast.Variadic{}},
						},
					},
					Ret: &ast.BuiltinType{Kind: ast.Void},
				},
			},
			symbs: []config.SymbolEntry{
				{
					CppName:    "foo",
					MangleName: "foo",
					GoName:     "Foo",
				},
			},
			expected: `
package testpkg

import _ "unsafe"
//go:linkname Foo C.foo
func Foo(__llgo_va_list ...interface{})
`,
		},
		{
			name: "func not in symbol table",
			decl: &ast.FuncDecl{
				Name:        &ast.Ident{Name: "foo"},
				MangledName: "foo",
				Type: &ast.FuncType{
					Params: nil,
					Ret:    nil,
				},
			},
			expected: `
package testpkg

import _ "unsafe"
			`,
		},
		{
			name: "invalid function type",
			decl: &ast.FuncDecl{
				Name:        &ast.Ident{Name: "invalidFunc"},
				MangledName: "invalidFunc",
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{{Name: "a"}},
								Type:  &ast.BuiltinType{Kind: ast.Bool, Flags: ast.Long}, // invalid
							},
						},
					},
					Ret: nil,
				},
			},
			symbs: []config.SymbolEntry{
				{
					CppName:    "invalidFunc",
					MangleName: "invalidFunc",
					GoName:     "InvalidFunc",
				},
			},
			expectedErr: "NewFuncDecl: fail convert signature invalidFunc: not found in type map",
		},
		{
			name: "explict void return",
			decl: &ast.FuncDecl{
				Name:        &ast.Ident{Name: "foo"},
				MangledName: "foo",
				Type: &ast.FuncType{
					Params: nil,
					Ret:    &ast.BuiltinType{Kind: ast.Void},
				},
			},
			symbs: []config.SymbolEntry{
				{
					CppName:    "foo",
					MangleName: "foo",
					GoName:     "Foo",
				},
			},
			expected: `
package testpkg

import _ "unsafe"
//go:linkname Foo C.foo
func Foo()
`,
		},
		{
			name: "builtin type",
			decl: &ast.FuncDecl{
				Name:        &ast.Ident{Name: "foo"},
				MangledName: "foo",
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{Name: "a"},
								},
								Type: &ast.BuiltinType{
									Kind:  ast.Int,
									Flags: ast.Short | ast.Unsigned},
							},
							{
								Names: []*ast.Ident{
									{Name: "b"},
								},
								Type: &ast.BuiltinType{
									Kind: ast.Bool,
								},
							},
						},
					},
					Ret: &ast.BuiltinType{
						Kind:  ast.Float,
						Flags: ast.Double,
					},
				},
			},

			symbs: []config.SymbolEntry{
				{
					CppName:    "foo",
					MangleName: "foo",
					GoName:     "Foo",
				},
			},
			expected: `
package testpkg

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)
//go:linkname Foo C.foo
func Foo(a uint16, b bool) c.Double
`,
		},
		{
			name: "c builtin type",
			decl: &ast.FuncDecl{
				Name:        &ast.Ident{Name: "foo"},
				MangledName: "foo",
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{{Name: "a"}},
								Type:  &ast.BuiltinType{Kind: ast.Int, Flags: ast.Unsigned},
							},
							{
								Names: []*ast.Ident{{Name: "b"}},
								Type:  &ast.BuiltinType{Kind: ast.Int, Flags: ast.Long},
							},
						},
					},
					Ret: &ast.BuiltinType{Kind: ast.Int, Flags: ast.Long | ast.Unsigned},
				},
			},
			symbs: []config.SymbolEntry{
				{
					CppName:    "foo",
					MangleName: "foo",
					GoName:     "Foo",
				},
			},
			expected: `
package testpkg

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)
//go:linkname Foo C.foo
func Foo(a c.Uint, b c.Long) c.Ulong
`,
		},
		{
			name: "basic decl with c type",
			decl: &ast.FuncDecl{
				Name:        &ast.Ident{Name: "foo"},
				MangledName: "foo",
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{{Name: "a"}},
								Type:  &ast.BuiltinType{Kind: ast.Int, Flags: ast.Unsigned},
							},
							{
								Names: []*ast.Ident{{Name: "b"}},
								Type:  &ast.BuiltinType{Kind: ast.Int, Flags: ast.Long},
							},
						},
					},
					Ret: &ast.BuiltinType{Kind: ast.Int, Flags: ast.Long | ast.Unsigned},
				},
			},
			symbs: []config.SymbolEntry{
				{
					CppName:    "foo",
					MangleName: "foo",
					GoName:     "Foo",
				},
			},
			expected: `
package testpkg

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)
//go:linkname Foo C.foo
func Foo(a c.Uint, b c.Long) c.Ulong
`,
		},
		{
			name: "pointer type",
			decl: &ast.FuncDecl{
				Name:        &ast.Ident{Name: "foo"},
				MangledName: "foo",
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{{Name: "a"}},
								Type: &ast.PointerType{
									X: &ast.BuiltinType{Kind: ast.Int, Flags: ast.Unsigned},
								},
							},
							{
								Names: []*ast.Ident{{Name: "b"}},
								Type: &ast.PointerType{
									X: &ast.BuiltinType{Kind: ast.Int, Flags: ast.Long},
								},
							},
						},
					},
					Ret: &ast.PointerType{
						X: &ast.BuiltinType{
							Kind:  ast.Float,
							Flags: ast.Double,
						},
					},
				},
			},
			symbs: []config.SymbolEntry{
				{
					CppName:    "foo",
					MangleName: "foo",
					GoName:     "Foo",
				},
			},
			expected: `
package testpkg

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)
//go:linkname Foo C.foo
func Foo(a *c.Uint, b *c.Long) *c.Double
`,
		},
		{
			name: "void *",
			decl: &ast.FuncDecl{
				Name:        &ast.Ident{Name: "foo"},
				MangledName: "foo",
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{{Name: "a"}},
								Type: &ast.PointerType{
									X: &ast.BuiltinType{Kind: ast.Void},
								},
							},
						},
					},
					Ret: &ast.PointerType{
						X: &ast.BuiltinType{Kind: ast.Void},
					},
				},
			},
			symbs: []config.SymbolEntry{
				{
					CppName:    "foo",
					MangleName: "foo",
					GoName:     "Foo",
				},
			},
			expected: `
package testpkg

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)
//go:linkname Foo C.foo
func Foo(a c.Pointer) c.Pointer
`,
		},
		{
			name: "array",
			decl: &ast.FuncDecl{
				Name:        &ast.Ident{Name: "foo"},
				MangledName: "foo",
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{{Name: "a"}},
								// Uint[]
								Type: &ast.ArrayType{
									Elt: &ast.BuiltinType{Kind: ast.Int, Flags: ast.Unsigned},
								},
							},
							{
								Names: []*ast.Ident{{Name: "b"}},
								// Double[3]
								Type: &ast.ArrayType{
									Elt: &ast.BuiltinType{Kind: ast.Float, Flags: ast.Double},
									Len: &ast.BasicLit{Kind: ast.IntLit, Value: "3"},
								},
							},
						},
					},
					Ret: &ast.ArrayType{
						// char[3][4]
						Elt: &ast.ArrayType{
							Elt: &ast.BuiltinType{
								Kind:  ast.Char,
								Flags: ast.Signed,
							},
							Len: &ast.BasicLit{Kind: ast.IntLit, Value: "4"},
						},
						Len: &ast.BasicLit{Kind: ast.IntLit, Value: "3"},
					},
				},
			},
			symbs: []config.SymbolEntry{
				{
					CppName:    "foo",
					MangleName: "foo",
					GoName:     "Foo",
				},
			},
			cppgconf: &llcppg.Config{
				Name: "testpkg",
			},
			expected: `
package testpkg

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)
//go:linkname Foo C.foo
func Foo(a *c.Uint, b *c.Double) **c.Char
`,
		},
		{
			name: "error array param",
			decl: &ast.FuncDecl{
				Name:        &ast.Ident{Name: "foo"},
				MangledName: "foo",
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.ArrayType{
									Elt: &ast.BuiltinType{Kind: ast.Int, Flags: ast.Double},
								},
							},
						},
					},
					Ret: nil,
				},
			},
			symbs: []config.SymbolEntry{
				{
					CppName:    "foo",
					MangleName: "foo",
					GoName:     "Foo",
				},
			},
			expectedErr: "NewFuncDecl: fail convert signature foo: error convert elem type: not found in type map",
		},
		{
			name: "error return type",
			decl: &ast.FuncDecl{
				Name:        &ast.Ident{Name: "foo"},
				MangledName: "foo",
				Type: &ast.FuncType{
					Params: nil,
					Ret:    &ast.BuiltinType{Kind: ast.Bool, Flags: ast.Double},
				},
			},
			symbs: []config.SymbolEntry{
				{
					CppName:    "foo",
					MangleName: "foo",
					GoName:     "Foo",
				},
			},
			expectedErr: "NewFuncDecl: fail convert signature foo: error convert return type: not found in type map",
		},
		{
			name: "error nil param",
			decl: &ast.FuncDecl{
				Name:        &ast.Ident{Name: "foo"},
				MangledName: "foo",
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							nil,
						},
					},
					Ret: nil,
				},
			},
			symbs: []config.SymbolEntry{
				{
					CppName:    "foo",
					MangleName: "foo",
					GoName:     "Foo",
				},
			},
			expectedErr: "NewFuncDecl: fail convert signature foo: error convert type: unexpected nil field",
		},
		{
			name: "error receiver",
			decl: &ast.FuncDecl{
				DeclBase: ast.DeclBase{
					Loc: &ast.Location{File: tempFile.File},
				},
				Name:        &ast.Ident{Name: "foo"},
				MangledName: "foo",
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.BuiltinType{Kind: ast.Int, Flags: ast.Double},
							},
						},
					},
				},
			},
			symbs: []config.SymbolEntry{
				{
					CppName:    "foo",
					MangleName: "foo",
					GoName:     "(*Foo).foo",
				},
			},
			expectedErr: "NewFuncDecl: foo fail: newReceiver:failed to convert type: not found in type map",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testGenDecl(t, tc)
		})
	}
}

func TestStructDecl(t *testing.T) {
	testCases := []genDeclTestCase{
		// struct Foo {}
		{
			name: "empty struct",
			decl: &ast.TypeDecl{
				Name: &ast.Ident{Name: "Foo"},
				Type: &ast.RecordType{
					Tag:    ast.Struct,
					Fields: nil,
				},
			},
			expected: `
package testpkg

import _ "unsafe"

type Foo struct {
}`,
		},
		// invalid struct type
		{
			name: "invalid struct type",
			decl: &ast.TypeDecl{
				Name: &ast.Ident{Name: "InvalidStruct"},
				Type: &ast.RecordType{
					Tag: ast.Struct,
					Fields: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{{Name: "invalidField"}},
								Type:  &ast.BuiltinType{Kind: ast.Bool, Flags: ast.Long},
							},
						},
					},
				},
			},
			expectedErr: "NewTypeDecl: fail to complete type InvalidStruct: not found in type map",
		},
		// struct Foo { int a; double b; bool c; }
		{
			name: "struct field builtin type",
			decl: &ast.TypeDecl{
				Name: &ast.Ident{Name: "Foo"},
				Type: &ast.RecordType{
					Tag: ast.Struct,
					Fields: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{{Name: "a"}},
								Type: &ast.BuiltinType{
									Kind: ast.Int,
								},
							},
							{
								Names: []*ast.Ident{{Name: "b"}},
								Type: &ast.BuiltinType{
									Kind:  ast.Float,
									Flags: ast.Double,
								},
							},
							{
								Names: []*ast.Ident{{Name: "c"}},
								Type: &ast.BuiltinType{
									Kind: ast.Bool,
								},
							},
						},
					},
				},
			},
			expected: `
package testpkg

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type Foo struct {
	A c.Int
	B c.Double
	C bool
}`,
		},
		// struct Foo { int* a; double* b; bool* c;void* d; }
		{
			name: "struct field pointer",
			decl: &ast.TypeDecl{
				Name: &ast.Ident{Name: "Foo"},
				Type: &ast.RecordType{
					Tag: ast.Struct,
					Fields: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{{Name: "a"}},
								Type: &ast.PointerType{
									X: &ast.BuiltinType{
										Kind: ast.Int,
									},
								},
							},
							{
								Names: []*ast.Ident{{Name: "b"}},
								Type: &ast.PointerType{
									X: &ast.BuiltinType{
										Kind:  ast.Float,
										Flags: ast.Double,
									}},
							},
							{
								Names: []*ast.Ident{{Name: "c"}},
								Type: &ast.PointerType{
									X: &ast.BuiltinType{
										Kind: ast.Bool,
									},
								},
							},
							{
								Names: []*ast.Ident{{Name: "d"}},
								Type: &ast.PointerType{
									X: &ast.BuiltinType{
										Kind: ast.Void,
									},
								},
							},
						},
					},
				},
			},
			expected: `
package testpkg

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type Foo struct {
	A *c.Int
	B *c.Double
	C *bool
	D c.Pointer
}`},
		// struct Foo { char a[4]; int b[3][4]; }
		{
			name: "struct array field",
			decl: &ast.TypeDecl{
				Name: &ast.Ident{Name: "Foo"},
				Type: &ast.RecordType{
					Tag: ast.Struct,
					Fields: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{{Name: "a"}},
								Type: &ast.ArrayType{
									Elt: &ast.BuiltinType{
										Kind:  ast.Char,
										Flags: ast.Signed,
									},
									Len: &ast.BasicLit{
										Kind:  ast.IntLit,
										Value: "4",
									},
								},
							},
							{
								Names: []*ast.Ident{{Name: "b"}},
								Type: &ast.ArrayType{
									Elt: &ast.ArrayType{
										Elt: &ast.BuiltinType{
											Kind: ast.Int,
										},
										Len: &ast.BasicLit{Kind: ast.IntLit, Value: "4"},
									},
									Len: &ast.BasicLit{Kind: ast.IntLit, Value: "3"},
								},
							},
						},
					},
				},
			},
			expected: `
package testpkg

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type Foo struct {
	A [4]c.Char
	B [3][4]c.Int
}`},
		{
			name: "struct array field",
			decl: &ast.TypeDecl{
				Name: &ast.Ident{Name: "Foo"},
				Type: &ast.RecordType{
					Tag: ast.Struct,
					Fields: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{{Name: "a"}},
								Type: &ast.ArrayType{
									Elt: &ast.BuiltinType{
										Kind:  ast.Char,
										Flags: ast.Signed,
									},
									Len: &ast.BasicLit{
										Kind:  ast.IntLit,
										Value: "4",
									},
								},
							},
							{
								Names: []*ast.Ident{{Name: "b"}},
								Type: &ast.ArrayType{
									Elt: &ast.ArrayType{
										Elt: &ast.BuiltinType{
											Kind: ast.Int,
										},
										Len: &ast.BasicLit{Kind: ast.IntLit, Value: "4"},
									},
									Len: &ast.BasicLit{Kind: ast.IntLit, Value: "3"},
								},
							},
						},
					},
				},
			},
			expected: `
package testpkg

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type Foo struct {
	A [4]c.Char
	B [3][4]c.Int
}`},
		{
			name: "struct array field without len",
			decl: &ast.TypeDecl{
				Name: &ast.Ident{Name: "Foo"},
				Type: &ast.RecordType{
					Tag: ast.Struct,
					Fields: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{{Name: "a"}},
								Type: &ast.ArrayType{
									Elt: &ast.BuiltinType{
										Kind:  ast.Char,
										Flags: ast.Signed,
									},
								},
							},
						},
					},
				},
			},
			expectedErr: "NewTypeDecl: fail to complete type Foo: unsupport field with array without length",
		},
		{
			name: "struct array field without len",
			decl: &ast.TypeDecl{
				Name: &ast.Ident{Name: "Foo"},
				Type: &ast.RecordType{
					Tag: ast.Struct,
					Fields: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{{Name: "a"}},
								Type: &ast.ArrayType{
									Elt: &ast.BuiltinType{
										Kind:  ast.Char,
										Flags: ast.Signed,
									},
									Len: &ast.BuiltinType{Kind: ast.TypeKind(ast.Signed)}, //invalid
								},
							},
						},
					},
				},
			},
			expectedErr: "NewTypeDecl: fail to complete type Foo: can't determine the array length",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testGenDecl(t, tc)
		})
	}
}

func TestTypedefFunc(t *testing.T) {
	testCases := []genDeclTestCase{
		// typedef int (*Foo) (int a, int b);
		{
			name: "typedef func",
			decl: &ast.TypedefDecl{
				Name: &ast.Ident{Name: "Foo"},
				Type: &ast.PointerType{
					X: &ast.FuncType{
						Params: &ast.FieldList{
							List: []*ast.Field{
								{
									Type: &ast.BuiltinType{
										Kind: ast.Int,
									},
									Names: []*ast.Ident{{Name: "a"}},
								},
								{
									Type: &ast.BuiltinType{
										Kind: ast.Int,
									},
									Names: []*ast.Ident{{Name: "b"}},
								},
							},
						},
						Ret: &ast.BuiltinType{
							Kind: ast.Int,
						},
					},
				},
			},
			expected: `
package testpkg

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)
// llgo:type C
type Foo func(a c.Int, b c.Int) c.Int
`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testGenDecl(t, tc)
		})
	}
}

// Test Redefine error
func TestRedef(t *testing.T) {
	pkg, err := createTestPkg(&convert.PackageConfig{
		OutputDir: "",
		SymbolTable: config.CreateSymbolTable(
			[]config.SymbolEntry{
				{CppName: "Bar", MangleName: "Bar", GoName: "Bar"},
			},
		),
	})
	if err != nil {
		t.Fatal("NewPackage failed:", err)
	}
	pkg.SetCurFile(tempFile)

	flds := &ast.FieldList{
		List: []*ast.Field{
			{
				Type: &ast.BuiltinType{Kind: ast.Int},
			},
		},
	}
	pkg.NewTypeDecl(&ast.TypeDecl{
		Name: &ast.Ident{Name: "Foo"},
		Type: &ast.RecordType{
			Tag:    ast.Struct,
			Fields: flds,
		},
	})

	err = pkg.NewTypeDecl(&ast.TypeDecl{
		Name: &ast.Ident{Name: "Foo"},
		Type: &ast.RecordType{
			Tag:    ast.Struct,
			Fields: flds,
		},
	})
	if err != nil {
		t.Fatal("unexpect redefine err")
	}

	err = pkg.NewFuncDecl(&ast.FuncDecl{
		Name:        &ast.Ident{Name: "Bar"},
		MangledName: "Bar",
		Type: &ast.FuncType{
			Ret: &ast.BuiltinType{
				Kind: ast.Void,
			},
		},
	})
	if err != nil {
		t.Fatal("NewFuncDecl failed", err)
	}

	err = pkg.NewFuncDecl(&ast.FuncDecl{
		Name:        &ast.Ident{Name: "Bar"},
		MangledName: "Bar",
		Type:        &ast.FuncType{},
	})
	if err != nil {
		t.Fatal("unexpect redefine err")
	}

	err = pkg.NewFuncDecl(&ast.FuncDecl{
		Name:        &ast.Ident{Name: "Bar"},
		MangledName: "Bar",
		Type:        &ast.FuncType{},
	})
	if err != nil {
		t.Fatal("unexpect redefine err")
	}

	macro := &ast.Macro{
		Loc:    &ast.Location{File: tempFile.File},
		Name:   "MACRO_FOO",
		Tokens: []*ast.Token{{Token: token.IDENT, Lit: "MACRO_FOO"}, {Token: token.LITERAL, Lit: "1"}},
	}
	err = pkg.NewMacro(macro)
	if err != nil {
		t.Fatal("unexpect redefine err")
	}

	err = pkg.NewMacro(macro)
	if err != nil {
		t.Fatal("unexpect redefine err")
	}

	var buf bytes.Buffer
	err = pkg.GetGenPackage().WriteTo(&buf)
	if err != nil {
		t.Fatalf("WriteTo failed: %v", err)
	}

	expect := `
package testpkg

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type Foo struct {
	 c.Int
}
//go:linkname Bar C.Bar
func Bar()

const MACRO_FOO = 1
`
	comparePackageOutput(t, pkg, expect)
}

func TestRedefEnum(t *testing.T) {
	typDecl := &ast.TypeDecl{
		Name: &ast.Ident{Name: "Foo"},
		Type: &ast.RecordType{
			Tag: ast.Struct,
			Fields: &ast.FieldList{
				List: []*ast.Field{
					{Type: &ast.BuiltinType{Kind: ast.Int}},
				},
			},
		},
	}
	t.Run("redefine enum type", func(t *testing.T) {
		pkg, err := createTestPkg(&convert.PackageConfig{})
		if err != nil {
			t.Fatal("NewPackage failed:", err)
		}
		pkg.SetCurFile(tempFile)
		pkg.NewTypeDecl(typDecl)
		err = pkg.NewEnumTypeDecl(&ast.EnumTypeDecl{
			Name: &ast.Ident{Name: "Foo"},
			Type: &ast.EnumType{},
		})
		if err == nil {
			t.Fatalf("expect a redefine error")
		}
	})

	t.Run("redefine enum item", func(t *testing.T) {
		pkg, err := createTestPkg(&convert.PackageConfig{})
		if err != nil {
			t.Fatal("NewPackage failed:", err)
		}
		pkg.SetCurFile(tempFile)
		pkg.NewTypeDecl(typDecl)
		pkg.NewEnumTypeDecl(&ast.EnumTypeDecl{
			Name: nil,
			Type: &ast.EnumType{
				Items: []*ast.EnumItem{
					{Name: &ast.Ident{Name: "Foo"}, Value: &ast.BasicLit{Kind: ast.IntLit, Value: "0"}},
					// check if skip same name
					{Name: &ast.Ident{Name: "Foo"}, Value: &ast.BasicLit{Kind: ast.IntLit, Value: "0"}},
				},
			},
		})
		comparePackageOutput(t, pkg, `
package testpkg

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type Foo struct {
	 c.Int
}

const Foo__1 c.Int = 0
`)
	})
}

func TestRedefTypedef(t *testing.T) {
	pkg, err := createTestPkg(&convert.PackageConfig{
		SymbolTable: config.CreateSymbolTable(
			[]config.SymbolEntry{
				{CppName: "Foo", MangleName: "Foo", GoName: "Foo"},
			},
		),
	})
	if err != nil {
		t.Fatal("NewPackage failed:", err)
	}
	pkg.SetCurFile(tempFile)

	err = pkg.NewTypeDecl(&ast.TypeDecl{
		Name: &ast.Ident{Name: "Foo"},
		Type: &ast.RecordType{
			Tag:    ast.Struct,
			Fields: &ast.FieldList{},
		},
	})
	if err != nil {
		t.Fatal("NewTypeDecl failed", err)
	}
	err = pkg.NewTypedefDecl(&ast.TypedefDecl{
		Name: &ast.Ident{Name: "Foo"},
		Type: &ast.Ident{Name: "Foo"},
	})
	if err == nil {
		t.Fatal("expect a redefine error")
	}
}

func TestRedefineFunc(t *testing.T) {
	pkg, err := createTestPkg(&convert.PackageConfig{
		SymbolTable: config.CreateSymbolTable(
			[]config.SymbolEntry{
				{CppName: "Foo", MangleName: "Foo", GoName: "Foo"},
			},
		),
	})
	if err != nil {
		t.Fatal("NewPackage failed:", err)
	}
	pkg.SetCurFile(tempFile)

	err = pkg.NewTypeDecl(&ast.TypeDecl{
		Name: &ast.Ident{Name: "Foo"},
		Type: &ast.RecordType{
			Tag:    ast.Struct,
			Fields: &ast.FieldList{},
		},
	})
	if err != nil {
		t.Fatal("NewTypeDecl failed", err)
	}
	err = pkg.NewFuncDecl(&ast.FuncDecl{
		Name:        &ast.Ident{Name: "Foo"},
		MangledName: "Foo",
		Type:        &ast.FuncType{},
	})
	if err == nil {
		t.Fatal("expect a redefine error")
	}
}

func TestTypedef(t *testing.T) {
	testCases := []genDeclTestCase{
		// typedef double DOUBLE;
		{
			name: "typedef double",
			decl: &ast.TypedefDecl{
				Name: &ast.Ident{Name: "DOUBLE"},
				Type: &ast.BuiltinType{
					Kind:  ast.Float,
					Flags: ast.Double,
				},
			},
			expected: `
package testpkg

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type DOUBLE c.Double`,
		},
		// invalid typedef
		{
			name: "invalid typedef",
			decl: &ast.TypedefDecl{
				Name: &ast.Ident{Name: "INVALID"},
				Type: &ast.BuiltinType{
					Kind:  ast.Bool,
					Flags: ast.Double,
				},
			},
			expectedErr: "NewTypedefDecl:fail to convert type INVALID: not found in type map",
		},
		// typedef int INT;
		{
			name: "typedef int",
			decl: &ast.TypedefDecl{
				Name: &ast.Ident{Name: "INT"},
				Type: &ast.BuiltinType{
					Kind: ast.Int,
				},
			},
			expected: `
package testpkg

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type INT c.Int
			`,
		},
		{
			name: "typedef array",
			decl: &ast.TypedefDecl{
				Name: &ast.Ident{Name: "name"},
				Type: &ast.ArrayType{
					Elt: &ast.BuiltinType{
						Kind:  ast.Char,
						Flags: ast.Signed,
					},
					Len: &ast.BasicLit{Kind: ast.IntLit, Value: "5"},
				},
			},
			expected: `
package testpkg

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type Name [5]c.Char`,
		},
		// typedef void* ctx;
		{
			name: "typedef pointer",
			decl: &ast.TypedefDecl{
				Name: &ast.Ident{Name: "ctx"},
				Type: &ast.PointerType{
					X: &ast.BuiltinType{
						Kind: ast.Void,
					},
				},
			},
			expected: `
package testpkg

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type Ctx c.Pointer`,
		},

		// typedef char* name;
		{
			name: "typedef pointer",
			decl: &ast.TypedefDecl{
				Name: &ast.Ident{Name: "name"},
				Type: &ast.PointerType{
					X: &ast.BuiltinType{
						Kind:  ast.Char,
						Flags: ast.Signed,
					},
				},
			},
			expected: `
package testpkg

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type Name *c.Char
`,
		},
		{
			name: "typedef invalid pointer",
			decl: &ast.TypedefDecl{
				Name: &ast.Ident{Name: "name"},
				Type: &ast.PointerType{
					X: &ast.BuiltinType{
						Kind:  ast.Char,
						Flags: ast.Double,
					},
				},
			},
			expectedErr: "NewTypedefDecl:fail to convert type name: error convert baseType: not found in type map",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testGenDecl(t, tc)
		})
	}
}

func TestEnumDecl(t *testing.T) {
	testCases := []genDeclTestCase{
		{
			name: "enum",
			decl: &ast.EnumTypeDecl{
				Name: &ast.Ident{Name: "Color"},
				Type: &ast.EnumType{
					Items: []*ast.EnumItem{
						{Name: &ast.Ident{Name: "Red"}, Value: &ast.BasicLit{Kind: ast.IntLit, Value: "0"}},
						{Name: &ast.Ident{Name: "Green"}, Value: &ast.BasicLit{Kind: ast.IntLit, Value: "1"}},
						{Name: &ast.Ident{Name: "Blue"}, Value: &ast.BasicLit{Kind: ast.IntLit, Value: "2"}},
					},
				},
			},
			expected: `
package testpkg

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type Color c.Int

const (
	Red   Color = 0
	Green Color = 1
	Blue  Color = 2
)
`,
		},
		{
			name: "anonymous enum",
			decl: &ast.EnumTypeDecl{
				Name: nil,
				Type: &ast.EnumType{
					Items: []*ast.EnumItem{
						{Name: &ast.Ident{Name: "red"}, Value: &ast.BasicLit{Kind: ast.IntLit, Value: "0"}},
						{Name: &ast.Ident{Name: "green"}, Value: &ast.BasicLit{Kind: ast.IntLit, Value: "1"}},
						{Name: &ast.Ident{Name: "blue"}, Value: &ast.BasicLit{Kind: ast.IntLit, Value: "2"}},
					},
				},
			},
			expected: `
package testpkg

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

const (
	Red   c.Int = 0
	Green c.Int = 1
	Blue  c.Int = 2
)`,
		},
		{
			name: "invalid enum item",
			decl: &ast.EnumTypeDecl{
				Name: nil,
				Type: &ast.EnumType{
					Items: []*ast.EnumItem{
						{Name: &ast.Ident{Name: "red"}, Value: &ast.ArrayType{Elt: &ast.BuiltinType{Kind: ast.Bool}}},
					},
				},
			},
			expectedErr: "NewEnumTypeDecl: <nil> fail: createEnumItems:fail to convert *ast.ArrayType to int",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testGenDecl(t, tc)
		})
	}
}

func TestIdentRefer(t *testing.T) {
	pkg, err := createTestPkg(&convert.PackageConfig{})
	if err != nil {
		t.Fatal("NewPackage failed:", err)
	}
	pkg.SetCurFile(&convert.HeaderFile{
		File:     "/path/to/stdio.h",
		FileType: llcppg.Third,
	})
	pkg.NewTypedefDecl(&ast.TypedefDecl{
		DeclBase: ast.DeclBase{
			Loc: &ast.Location{File: "/path/to/stdio.h"},
		},
		Name: &ast.Ident{Name: "undefType"},
		Type: &ast.BuiltinType{
			Kind:  ast.Char,
			Flags: ast.Signed,
		},
	})
	pkg.SetCurFile(&convert.HeaderFile{
		File:     "/path/to/notsys.h",
		FileType: llcppg.Inter,
	})
	t.Run("undef sys ident ref", func(t *testing.T) {
		err := pkg.NewTypeDecl(&ast.TypeDecl{
			DeclBase: ast.DeclBase{
				Loc: &ast.Location{File: "/path/to/notsys.h"},
			},
			Name: &ast.Ident{Name: "Foo"},
			Type: &ast.RecordType{
				Tag: ast.Struct,
				Fields: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{{Name: "notfound"}},
							Type: &ast.Ident{
								Name: "undefType",
							},
						},
					},
				},
			},
		})
		if err == nil {
			t.Fatal("expect a error")
		}
	})
	t.Run("undef tag ident ref", func(t *testing.T) {
		err := pkg.NewTypeDecl(&ast.TypeDecl{
			Name: &ast.Ident{Name: "Bar"},
			Type: &ast.RecordType{
				Tag: ast.Struct,
				Fields: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{{Name: "notfound"}},
							Type: &ast.TagExpr{
								Tag: ast.Class,
								Name: &ast.Ident{
									Name: "undefType",
								},
							},
						},
					},
				},
			},
		})
		if err == nil {
			t.Fatal("expect a error")
		}
	})
	t.Run("type alias", func(t *testing.T) {
		pkg, err := createTestPkg(&convert.PackageConfig{})
		if err != nil {
			t.Fatal("NewPackage failed:", err)
		}
		pkg.SetCurFile(tempFile)
		err = pkg.NewTypedefDecl(&ast.TypedefDecl{
			Name: &ast.Ident{Name: "typ_int8_t"},
			Type: &ast.BuiltinType{
				Kind:  ast.Char,
				Flags: ast.Signed,
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		err = pkg.NewTypeDecl(&ast.TypeDecl{
			Name: &ast.Ident{Name: "Foo"},
			Type: &ast.RecordType{
				Tag: ast.Struct,
				Fields: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{{Name: "a"}},
							Type: &ast.Ident{
								Name: "typ_int8_t",
							},
						},
					},
				},
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		expect := `
package testpkg

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type TypInt8T c.Char

type Foo struct {
	A TypInt8T
}
`
		comparePackageOutput(t, pkg, expect)
	})
}

func TestForwardDecl(t *testing.T) {
	pkg, err := createTestPkg(&convert.PackageConfig{
		OutputDir: "",
		SymbolTable: config.CreateSymbolTable(
			[]config.SymbolEntry{
				{CppName: "Bar", MangleName: "Bar", GoName: "Bar"},
			},
		),
	})
	if err != nil {
		t.Fatal("NewPackage failed:", err)
	}
	pkg.SetCurFile(tempFile)

	forwardDecl := &ast.TypeDecl{
		Name: &ast.Ident{Name: "Foo"},
		Type: &ast.RecordType{
			Tag:    ast.Struct,
			Fields: &ast.FieldList{},
		},
	}
	// forward decl
	err = pkg.NewTypeDecl(forwardDecl)
	if err != nil {
		t.Fatalf("Forward decl failed: %v", err)
	}

	// complete decl
	err = pkg.NewTypeDecl(&ast.TypeDecl{
		Name: &ast.Ident{Name: "Foo"},
		Type: &ast.RecordType{
			Tag: ast.Struct,
			Fields: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{{Name: "a"}},
						Type:  &ast.BuiltinType{Kind: ast.Int},
					},
				},
			},
		},
	})

	if err != nil {
		t.Fatalf("NewTypeDecl failed: %v", err)
	}

	err = pkg.NewTypeDecl(forwardDecl)

	if err != nil {
		t.Fatalf("NewTypeDecl failed: %v", err)
	}

	expect := `
package testpkg

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type Foo struct {
	A c.Int
}
`
	comparePackageOutput(t, pkg, expect)
}

type genDeclTestCase struct {
	name        string
	decl        ast.Decl
	symbs       []config.SymbolEntry
	cppgconf    *llcppg.Config
	expected    string
	expectedErr string
}

func testGenDecl(t *testing.T, tc genDeclTestCase) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Fatal("unexpect panic", r)
		}
	}()
	var libCommand string
	var deps []string
	if tc.cppgconf != nil {
		libCommand = tc.cppgconf.Libs
		deps = tc.cppgconf.Deps
	}
	pkg, err := createTestPkg(&convert.PackageConfig{
		SymbolTable: config.CreateSymbolTable(tc.symbs),
		LibCommand:  libCommand,
		PkgBase: convert.PkgBase{
			Deps: deps,
		},
	})
	if err != nil {
		t.Fatal("NewPackage failed:", err)
	}
	if pkg == nil {
		t.Fatal("NewPackage failed")
	}
	pkg.SetCurFile(tempFile)
	switch d := tc.decl.(type) {
	case *ast.TypeDecl:
		err = pkg.NewTypeDecl(d)
	case *ast.TypedefDecl:
		err = pkg.NewTypedefDecl(d)
	case *ast.FuncDecl:
		err = pkg.NewFuncDecl(d)
	case *ast.EnumTypeDecl:
		err = pkg.NewEnumTypeDecl(d)
	default:
		t.Errorf("Unsupported declaration type: %T", tc.decl)
		return
	}
	if tc.expectedErr != "" {
		compareError(t, err, tc.expectedErr)
	} else {
		if err != nil {
			t.Errorf("Declaration generation failed: %v", err)
		} else {
			comparePackageOutput(t, pkg, tc.expected)
		}
	}
}

// compare error
func compareError(t *testing.T, err error, expectErr string) {
	t.Helper()
	if err == nil {
		t.Errorf("Expected error containing %q, but got nil", expectErr)
	} else if !strings.Contains(err.Error(), expectErr) {
		t.Errorf("Expected error contain %q, but got %q", expectErr, err.Error())
	}
}

func createTestPkg(cfg *convert.PackageConfig) (*convert.Package, error) {
	if cfg.SymbolTable == nil {
		cfg.SymbolTable = config.CreateSymbolTable([]config.SymbolEntry{})
	}
	if cfg.LibCommand == "" {
		cfg.LibCommand = "${pkg-config --libs xxx}"
	}
	return convert.NewPackage(&convert.PackageConfig{
		PkgBase: convert.PkgBase{
			PkgPath: ".",
			Deps:    cfg.Deps,
			Pubs:    make(map[string]string),
		},
		Name:           "testpkg",
		GenConf:        &gogen.Config{},
		OutputDir:      cfg.OutputDir,
		SymbolTable:    cfg.SymbolTable,
		LibCommand:     cfg.LibCommand,
		TrimPrefixes:   cfg.TrimPrefixes,
		KeepUnderScore: cfg.KeepUnderScore,
	})
}

// compares the output of a gogen.Package with the expected
func comparePackageOutput(t *testing.T, pkg *convert.Package, expect string) {
	t.Helper()
	// For Test,The Test package's header filename same as package name
	buf, err := pkg.WriteToBuffer("temp.go")
	if err != nil {
		t.Fatalf("WriteTo failed: %v", err)
	}
	expectedStr := strings.TrimSpace(expect)
	actualStr := strings.TrimSpace(buf.String())
	if expectedStr != actualStr {
		t.Errorf("does not match expected.\nExpected:\n%s\nGot:\n%s", expectedStr, actualStr)
	}
}

/** multiple package test **/

func TestTypeClean(t *testing.T) {
	pkg, err := createTestPkg(&convert.PackageConfig{
		OutputDir: "",
		SymbolTable: config.CreateSymbolTable(
			[]config.SymbolEntry{
				{CppName: "Func1", MangleName: "Func1", GoName: "Func1"},
				{CppName: "Func2", MangleName: "Func2", GoName: "Func2"},
			},
		),
	})
	if err != nil {
		t.Fatal("NewPackage failed:", err)
	}
	testCases := []struct {
		addType    func()
		headerFile string
		incPath    string
		newType    string
	}{
		{
			addType: func() {
				pkg.NewTypeDecl(&ast.TypeDecl{
					Name: &ast.Ident{Name: "Foo1"},
					Type: &ast.RecordType{Tag: ast.Struct},
				})
			},
			headerFile: "/path/to/file1.h",
			incPath:    "file1.h",
			newType:    "Foo1",
		},
		{
			addType: func() {
				pkg.NewTypedefDecl(&ast.TypedefDecl{
					Name: &ast.Ident{Name: "Bar2"},
					Type: &ast.BuiltinType{Kind: ast.Int},
				})
			},
			headerFile: "/path/to/file2.h",
			incPath:    "file2.h",
			newType:    "Bar2",
		},
		{
			addType: func() {
				pkg.NewFuncDecl(&ast.FuncDecl{
					Name: &ast.Ident{Name: "Func1"}, MangledName: "Func1",
					Type: &ast.FuncType{Params: nil, Ret: &ast.BuiltinType{Kind: ast.Void}},
				})
			},
			headerFile: "/path/to/file3.h",
			incPath:    "file3.h",
			newType:    "Func1",
		},
	}

	for i, tc := range testCases {
		pkg.SetCurFile(&convert.HeaderFile{
			File:     tc.headerFile,
			FileType: llcppg.Inter,
		})
		tc.addType()

		goFileName := name.HeaderFileToGo(tc.headerFile)
		buf, err := pkg.WriteToBuffer(goFileName)
		if err != nil {
			t.Fatal(err)
		}
		result := buf.String()

		if !strings.Contains(result, tc.newType) {
			t.Errorf("Case %d: Generated type does not contain %s", i, tc.newType)
		}

		for j := 0; j < i; j++ {
			oldType := testCases[j].newType
			if strings.Contains(result, oldType) {
				t.Errorf("Case %d: Previously added type %s (from case %d) still exists", i, oldType, j)
			}
		}
	}
}

func TestHeaderFileToGo(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal",
			input:    "/path/to/sys/dirent.h",
			expected: "dirent.go",
		},
		{
			name:     "sys",
			input:    "/path/to/sys/_pthread/_pthread_types.h",
			expected: "X_pthread_types.go",
		},
		{
			name:     "sys",
			input:    "/path/to/_types.h",
			expected: "X_types.go",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := name.HeaderFileToGo(tc.input)
			if result != tc.expected {
				t.Errorf("Expected %s, but got %s", tc.expected, result)
			}
		})
	}
}

func TestImport(t *testing.T) {
	t.Run("invalid include path", func(t *testing.T) {
		p := &convert.Package{}
		genPkg := gogen.NewPackage(".", "include", nil)
		mod, err := gopmod.Load(".")
		if err != nil {
			t.Fatal(err)
		}
		deps := []string{
			"github.com/goplus/llcppg/cl/internal/convert/testdata/invalidpath",
			"github.com/goplus/llcppg/cl/internal/convert/testdata/partfinddep",
		}
		p.PkgInfo = convert.NewPkgInfo(".", ".", deps, nil)
		loader := convert.NewPkgDepLoader(mod, genPkg)
		depPkgs, err := loader.LoadDeps(p.PkgInfo)
		p.PkgInfo.Deps = depPkgs
		if err != nil {
			t.Fatal(err)
		}
		_, err = loader.Import("github.com/goplus/invalidpkg")
		if err == nil {
			t.Fatal("expected error")
		}
	})
	t.Run("invalid pub file", func(t *testing.T) {
		_, err := createTestPkg(&convert.PackageConfig{
			OutputDir: ".",
			PkgBase: convert.PkgBase{
				Deps: []string{
					"github.com/goplus/llcppg/cl/internal/convert/testdata/invalidpub",
				},
			},
		})
		if err == nil {
			t.Fatal("NewPackage failed:", err)
		}
	})
	t.Run("invalid dep", func(t *testing.T) {
		_, err := createTestPkg(&convert.PackageConfig{
			OutputDir: ".",
			PkgBase: convert.PkgBase{
				Deps: []string{
					"github.com/goplus/llcppg/cl/internal/convert/testdata/invaliddep",
				},
			},
		})
		if err == nil {
			t.Fatal("NewPackage failed:", err)
		}
	})
	t.Run("same type register", func(t *testing.T) {
		_, err := createTestPkg(&convert.PackageConfig{
			OutputDir: ".",
			PkgBase: convert.PkgBase{
				Deps: []string{
					"github.com/goplus/llcppg/cl/internal/convert/testdata/cjson",
					"github.com/goplus/llcppg/cl/internal/convert/testdata/cjsonbool",
				},
			},
		})
		if err != nil {
			t.Fatal("NewPackage failed:", err)
		}
	})
}

func TestUnkownHfile(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("Expect Error")
		}
	}()
	convert.NewHeaderFile("/path/to/foo.h", 0).ToGoFileName("Pkg")
}
