package convert_test

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/goplus/gogen"
	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/cl/internal/cltest"
	"github.com/goplus/llcppg/cl/internal/convert"
	"github.com/goplus/llcppg/cl/nc"
	"github.com/goplus/llcppg/cl/nc/ncimpl"
	llcppg "github.com/goplus/llcppg/config"
	"github.com/goplus/llcppg/internal/name"
	"github.com/goplus/llcppg/token"
	"github.com/goplus/mod/xgomod"
)

var dir string

var tempFile = &ncimpl.HeaderFile{
	File:     "/path/to/temp.h",
	FileType: llcppg.Inter,
}

var pkgname = "testpkg"

func SetGoFile(ctx *convert.Package, goFile string) {
	pkg := ctx.Pkg()
	pkg.SetCurFile(goFile, true)
	pkg.Unsafe().MarkForceUsed(pkg)
}

func SetTempFile(ctx *convert.Package) {
	SetGoFile(ctx, tempFile.ToGoFileName(pkgname))
}

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
				Object: ast.Object{
					Name: &ast.Ident{Name: "u"},
				},
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
	pkg, err := createTestPkg(nil, &convert.PackageConfig{
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
	pkg, err := createTestPkg(nil, &convert.PackageConfig{
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

func TestNewPackage(t *testing.T) {
	pkg, err := createTestPkg(nil, &convert.PackageConfig{})
	if err != nil {
		t.Fatal("NewPackage failed:", err)
	}
	SetTempFile(pkg)
	comparePackageOutput(t, pkg, `
package testpkg

import _ "unsafe"
	`)
}

func TestFuncDecl(t *testing.T) {
	testCases := []genDeclTestCase{
		{
			name: "empty func",
			decl: &ast.FuncDecl{
				Object: ast.Object{
					Name: &ast.Ident{Name: "foo"},
				},
				MangledName: "foo",
				Type: &ast.FuncType{
					Params: &ast.FieldList{List: []*ast.Field{}},
					Ret:    &ast.BuiltinType{Kind: ast.Void},
				},
			},
			symbs: []llcppg.SymbolInfo{
				{
					Mangle: "foo",
					CPP:    "foo",
					Go:     "Foo",
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
				Object: ast.Object{
					Name: &ast.Ident{Name: "foo"},
				},
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
			symbs: []llcppg.SymbolInfo{
				{
					Mangle: "foo",
					CPP:    "foo",
					Go:     "Foo",
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
			name: "invalid function type",
			decl: &ast.FuncDecl{
				Object: ast.Object{
					Name: &ast.Ident{Name: "invalidFunc"},
				},
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
			symbs: []llcppg.SymbolInfo{
				{
					Mangle: "invalidFunc",
					CPP:    "invalidFunc",
					Go:     "InvalidFunc",
				},
			},
			expectedErr: "NewFuncDecl: fail convert signature invalidFunc: not found in type map",
		},
		{
			name: "explict void return",
			decl: &ast.FuncDecl{
				Object: ast.Object{
					Name: &ast.Ident{Name: "foo"},
				},
				MangledName: "foo",
				Type: &ast.FuncType{
					Params: &ast.FieldList{List: []*ast.Field{}},
					Ret:    &ast.BuiltinType{Kind: ast.Void},
				},
			},
			symbs: []llcppg.SymbolInfo{
				{
					Mangle: "foo",
					CPP:    "foo",
					Go:     "Foo",
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
				Object: ast.Object{
					Name: &ast.Ident{Name: "foo"},
				},
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

			symbs: []llcppg.SymbolInfo{
				{
					Mangle: "foo",
					CPP:    "foo",
					Go:     "Foo",
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
				Object: ast.Object{
					Name: &ast.Ident{Name: "foo"},
				},
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
			symbs: []llcppg.SymbolInfo{
				{
					Mangle: "foo",
					CPP:    "foo",
					Go:     "Foo",
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
				Object: ast.Object{
					Name: &ast.Ident{Name: "foo"},
				},
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
			symbs: []llcppg.SymbolInfo{
				{
					Mangle: "foo",
					CPP:    "foo",
					Go:     "Foo",
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
				Object: ast.Object{
					Name: &ast.Ident{Name: "foo"},
				},
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
			symbs: []llcppg.SymbolInfo{
				{
					Mangle: "foo",
					CPP:    "foo",
					Go:     "Foo",
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
				Object: ast.Object{
					Name: &ast.Ident{Name: "foo"},
				},
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
			symbs: []llcppg.SymbolInfo{
				{
					Mangle: "foo",
					CPP:    "foo",
					Go:     "Foo",
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
				Object: ast.Object{
					Name: &ast.Ident{Name: "foo"},
				},
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
			symbs: []llcppg.SymbolInfo{
				{
					Mangle: "foo",
					CPP:    "foo",
					Go:     "Foo",
				},
			},
			cppgconf: &llcppg.Config{
				Name: pkgname,
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
				Object: ast.Object{
					Name: &ast.Ident{Name: "foo"},
				},
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
			symbs: []llcppg.SymbolInfo{
				{
					Mangle: "foo",
					CPP:    "foo",
					Go:     "Foo",
				},
			},
			expectedErr: "NewFuncDecl: fail convert signature foo: error convert elem type: not found in type map",
		},
		{
			name: "error return type",
			decl: &ast.FuncDecl{
				Object: ast.Object{
					Name: &ast.Ident{Name: "foo"},
				},
				MangledName: "foo",
				Type: &ast.FuncType{
					Params: &ast.FieldList{List: []*ast.Field{}},
					Ret:    &ast.BuiltinType{Kind: ast.Bool, Flags: ast.Double},
				},
			},
			symbs: []llcppg.SymbolInfo{
				{
					Mangle: "foo",
					CPP:    "foo",
					Go:     "Foo",
				},
			},
			expectedErr: "NewFuncDecl: fail convert signature foo: error convert return type: not found in type map",
		},
		{
			name: "error nil param",
			decl: &ast.FuncDecl{
				Object: ast.Object{
					Name: &ast.Ident{Name: "foo"},
				},
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
			symbs: []llcppg.SymbolInfo{
				{
					Mangle: "foo",
					CPP:    "foo",
					Go:     "Foo",
				},
			},
			expectedErr: "NewFuncDecl: fail convert signature foo: error convert type: unexpected nil field",
		},
		{
			name: "error receiver",
			decl: &ast.FuncDecl{
				Object: ast.Object{
					Loc:  &ast.Location{File: tempFile.File},
					Name: &ast.Ident{Name: "foo"},
				},
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
			symbs: []llcppg.SymbolInfo{
				{
					Mangle: "foo",
					CPP:    "foo",
					Go:     "(*Foo).foo",
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
				Object: ast.Object{
					Name: &ast.Ident{Name: "Foo"},
				},
				Type: &ast.RecordType{
					Tag:    ast.Struct,
					Fields: &ast.FieldList{},
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
				Object: ast.Object{
					Name: &ast.Ident{Name: "InvalidStruct"},
				},
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
				Object: ast.Object{
					Name: &ast.Ident{Name: "Foo"},
				},
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
				Object: ast.Object{
					Name: &ast.Ident{Name: "Foo"},
				},
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
				Object: ast.Object{
					Name: &ast.Ident{Name: "Foo"},
				},
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
				Object: ast.Object{
					Name: &ast.Ident{Name: "Foo"},
				},
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
				Object: ast.Object{
					Name: &ast.Ident{Name: "Foo"},
				},
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
				Object: ast.Object{
					Name: &ast.Ident{Name: "Foo"},
				},
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
				Object: ast.Object{
					Name: &ast.Ident{Name: "Foo"},
				},
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
	nc := cltest.NC(&llcppg.Config{}, nil, cltest.NewConvSym(
		llcppg.SymbolInfo{
			Mangle: "Bar",
			CPP:    "Bar",
			Go:     "Bar",
		},
	))
	pkg, err := createTestPkg(nc, &convert.PackageConfig{
		OutputDir: "",
	})
	if err != nil {
		t.Fatal("NewPackage failed:", err)
	}
	SetTempFile(pkg)

	flds := &ast.FieldList{
		List: []*ast.Field{
			{
				Type: &ast.BuiltinType{Kind: ast.Int},
			},
		},
	}
	pkg.NewTypeDecl("Foo", &ast.TypeDecl{
		Object: ast.Object{
			Name: &ast.Ident{Name: "Foo"},
		},
		Type: &ast.RecordType{
			Tag:    ast.Struct,
			Fields: flds,
		},
	}, nc)

	err = pkg.NewTypeDecl("Foo", &ast.TypeDecl{
		Object: ast.Object{
			Name: &ast.Ident{Name: "Foo"},
		},
		Type: &ast.RecordType{
			Tag:    ast.Struct,
			Fields: flds,
		},
	}, nc)
	if err != nil {
		t.Fatal("unexpect redefine err")
	}

	err = pkg.NewFuncDecl("Bar", &ast.FuncDecl{
		Object: ast.Object{
			Name: &ast.Ident{Name: "Bar"},
		},
		MangledName: "Bar",
		Type: &ast.FuncType{
			Ret: &ast.BuiltinType{
				Kind: ast.Void,
			},
			Params: &ast.FieldList{List: []*ast.Field{}},
		},
	})
	if err != nil {
		t.Fatal("NewFuncDecl failed", err)
	}

	err = pkg.NewFuncDecl("Bar", &ast.FuncDecl{
		Object: ast.Object{
			Name: &ast.Ident{Name: "Bar"},
		},
		MangledName: "Bar",
		Type: &ast.FuncType{
			Params: &ast.FieldList{List: []*ast.Field{}},
		},
	})
	if err != nil {
		t.Fatal("unexpect redefine err")
	}

	err = pkg.NewFuncDecl("Bar", &ast.FuncDecl{
		Object: ast.Object{
			Name: &ast.Ident{Name: "Bar"},
		},
		MangledName: "Bar",
		Type: &ast.FuncType{
			Params: &ast.FieldList{List: []*ast.Field{}},
		},
	})
	if err != nil {
		t.Fatal("unexpect redefine err")
	}

	macro := &ast.Macro{
		Loc:    &ast.Location{File: tempFile.File},
		Name:   "MACRO_FOO",
		Tokens: []*ast.Token{{Token: token.IDENT, Lit: "MACRO_FOO"}, {Token: token.LITERAL, Lit: "1"}},
	}
	err = pkg.NewMacro("MACRO_FOO", macro)
	if err != nil {
		t.Fatal("unexpect redefine err")
	}

	err = pkg.NewMacro("MACRO_FOO", macro)
	if err != nil {
		t.Fatal("unexpect redefine err")
	}

	var buf bytes.Buffer
	err = pkg.Pkg().WriteTo(&buf)
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
	typeName := "Foo"
	typDecl := &ast.TypeDecl{
		Object: ast.Object{
			Name: &ast.Ident{Name: typeName},
		},
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
		nc := cltest.NC(&llcppg.Config{}, nil, cltest.NewConvSym())
		pkg, err := createTestPkg(nc, &convert.PackageConfig{})
		if err != nil {
			t.Fatal("NewPackage failed:", err)
		}
		SetTempFile(pkg)
		pkg.NewTypeDecl(typeName, typDecl, nc)
		err = pkg.NewEnumTypeDecl(typeName, &ast.EnumTypeDecl{
			Object: ast.Object{
				Name: &ast.Ident{Name: typeName},
			},
			Type: &ast.EnumType{},
		}, nc)
		if err == nil {
			t.Fatalf("expect a redefine error")
		}
	})

	t.Run("redefine enum item", func(t *testing.T) {
		nc := cltest.NC(&llcppg.Config{}, nil, cltest.NewConvSym())
		pkg, err := createTestPkg(nc, &convert.PackageConfig{})
		if err != nil {
			t.Fatal("NewPackage failed:", err)
		}
		SetTempFile(pkg)
		pkg.NewTypeDecl(typeName, typDecl, nc)
		pkg.NewEnumTypeDecl(typeName, &ast.EnumTypeDecl{
			Object: ast.Object{
				Name: nil,
			},
			Type: &ast.EnumType{
				Items: []*ast.EnumItem{
					{Name: &ast.Ident{Name: "Foo"}, Value: &ast.BasicLit{Kind: ast.IntLit, Value: "0"}},
					// check if skip same name
					{Name: &ast.Ident{Name: "Foo"}, Value: &ast.BasicLit{Kind: ast.IntLit, Value: "0"}},
				},
			},
		}, nc)
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
	nc := cltest.NC(&llcppg.Config{}, nil, cltest.NewConvSym(llcppg.SymbolInfo{
		Mangle: "Foo",
		CPP:    "Foo",
		Go:     "Foo",
	}))
	pkg, err := createTestPkg(nc, &convert.PackageConfig{})
	if err != nil {
		t.Fatal("NewPackage failed:", err)
	}
	SetTempFile(pkg)

	err = pkg.NewTypeDecl("Foo", &ast.TypeDecl{
		Object: ast.Object{
			Name: &ast.Ident{Name: "Foo"},
		},
		Type: &ast.RecordType{
			Tag:    ast.Struct,
			Fields: &ast.FieldList{},
		},
	}, nc)
	if err != nil {
		t.Fatal("NewTypeDecl failed", err)
	}
	err = pkg.NewTypedefDecl("Foo", &ast.TypedefDecl{
		Object: ast.Object{
			Name: &ast.Ident{Name: "Foo"},
		},
		Type: &ast.Ident{Name: "Foo"},
	}, nc)
	if err == nil {
		t.Fatal("expect a redefine error")
	}
}

func TestRedefineFunc(t *testing.T) {
	nc := cltest.NC(&llcppg.Config{}, nil, cltest.NewConvSym(llcppg.SymbolInfo{
		Mangle: "Foo",
		CPP:    "Foo",
		Go:     "Foo",
	}))
	pkg, err := createTestPkg(nc, &convert.PackageConfig{})
	if err != nil {
		t.Fatal("NewPackage failed:", err)
	}
	SetTempFile(pkg)

	err = pkg.NewTypeDecl("Foo", &ast.TypeDecl{
		Object: ast.Object{
			Name: &ast.Ident{Name: "Foo"},
		},
		Type: &ast.RecordType{
			Tag:    ast.Struct,
			Fields: &ast.FieldList{},
		},
	}, nc)
	if err != nil {
		t.Fatal("NewTypeDecl failed", err)
	}
	err = pkg.NewFuncDecl("Foo", &ast.FuncDecl{
		Object: ast.Object{
			Name: &ast.Ident{Name: "Foo"},
		},
		MangledName: "Foo",
		Type: &ast.FuncType{
			Params: &ast.FieldList{List: []*ast.Field{}},
		},
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
				Object: ast.Object{
					Name: &ast.Ident{Name: "DOUBLE"},
				},
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
				Object: ast.Object{
					Name: &ast.Ident{Name: "INVALID"},
				},
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
				Object: ast.Object{
					Name: &ast.Ident{Name: "INT"},
				},
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
				Object: ast.Object{
					Name: &ast.Ident{Name: "name"},
				},
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
				Object: ast.Object{
					Name: &ast.Ident{Name: "ctx"},
				},
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
				Object: ast.Object{
					Name: &ast.Ident{Name: "name"},
				},
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
				Object: ast.Object{
					Name: &ast.Ident{Name: "name"},
				},
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
				Object: ast.Object{
					Name: &ast.Ident{Name: "Color"},
				},
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
				Object: ast.Object{
					Name: nil,
				},
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

func TestTypeAlias(t *testing.T) {
	nc := cltest.NC(&llcppg.Config{}, nil, cltest.NewConvSym())
	pkg, err := createTestPkg(nc, &convert.PackageConfig{})
	if err != nil {
		t.Fatal("NewPackage failed:", err)
	}
	SetTempFile(pkg)
	err = pkg.NewTypedefDecl("TypInt8T", &ast.TypedefDecl{
		Object: ast.Object{
			Name: &ast.Ident{Name: "typ_int8_t"},
		},
		Type: &ast.BuiltinType{
			Kind:  ast.Char,
			Flags: ast.Signed,
		},
	}, nc)
	if err != nil {
		t.Fatal(err)
	}
	err = pkg.NewTypeDecl("Foo", &ast.TypeDecl{
		Object: ast.Object{
			Name: &ast.Ident{Name: "Foo"},
		},
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
	}, nc)
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
}

func TestForwardDecl(t *testing.T) {
	nc := cltest.NC(&llcppg.Config{}, nil, cltest.NewConvSym(
		llcppg.SymbolInfo{
			Mangle: "Bar",
			CPP:    "Bar",
			Go:     "Bar",
		},
	))
	pkg, err := createTestPkg(nc, &convert.PackageConfig{
		OutputDir: "",
	})
	if err != nil {
		t.Fatal("NewPackage failed:", err)
	}
	SetTempFile(pkg)

	forwardDecl := &ast.TypeDecl{
		Object: ast.Object{
			Name: &ast.Ident{Name: "Foo"},
		},
		Type: &ast.RecordType{
			Tag:    ast.Struct,
			Fields: nil,
		},
	}
	// forward decl
	err = pkg.NewTypeDecl("Foo", forwardDecl, nc)
	if err != nil {
		t.Fatalf("Forward decl failed: %v", err)
	}

	// complete decl
	err = pkg.NewTypeDecl("Foo", &ast.TypeDecl{
		Object: ast.Object{
			Name: &ast.Ident{Name: "Foo"},
		},
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
	}, nc)

	if err != nil {
		t.Fatalf("NewTypeDecl failed: %v", err)
	}

	err = pkg.NewTypeDecl("Foo", forwardDecl, nc)

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
	symbs       []llcppg.SymbolInfo
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
	if tc.cppgconf == nil {
		tc.cppgconf = &llcppg.Config{Name: pkgname}
	}
	fileMap := make(map[string]*llcppg.FileInfo)
	fileMap["/path/to/temp.h"] = &llcppg.FileInfo{
		FileType: llcppg.Inter,
	}

	nc := cltest.NC(tc.cppgconf, fileMap, cltest.NewConvSym(tc.symbs...))
	pkg, err := createTestPkg(nc, &convert.PackageConfig{
		LibCommand: libCommand,
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
	goName, _, err := nc.ConvDecl("/path/to/temp.h", tc.decl)
	if err != nil {
		t.Fatalf("ConvDecl failed: %v", err)
	}
	SetTempFile(pkg)
	switch d := tc.decl.(type) {
	case *ast.TypeDecl:
		err = pkg.NewTypeDecl(goName, d, nc)
	case *ast.TypedefDecl:
		err = pkg.NewTypedefDecl(goName, d, nc)
	case *ast.FuncDecl:
		err = pkg.NewFuncDecl(goName, d)
	case *ast.EnumTypeDecl:
		err = pkg.NewEnumTypeDecl(goName, d, nc)
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

// // compare error
func compareError(t *testing.T, err error, expectErr string) {
	t.Helper()
	if err == nil {
		t.Errorf("Expected error containing %q, but got nil", expectErr)
	} else if !strings.Contains(err.Error(), expectErr) {
		t.Errorf("Expected error contain %q, but got %q", expectErr, err.Error())
	}
}

func createTestPkg(nc nc.NodeConverter, cfg *convert.PackageConfig) (*convert.Package, error) {
	pnc := nc
	if pnc == nil {
		pnc = cltest.NC(&llcppg.Config{}, nil, cltest.NewConvSym())
	}
	if cfg.LibCommand == "" {
		cfg.LibCommand = "${pkg-config --libs xxx}"
	}
	return convert.NewPackage(pnc, &convert.PackageConfig{
		PkgBase: convert.PkgBase{
			PkgPath: ".",
			Deps:    cfg.Deps,
			Pubs:    make(map[string]string),
		},
		Name:       pkgname,
		GenConf:    &gogen.Config{},
		OutputDir:  cfg.OutputDir,
		LibCommand: cfg.LibCommand,
	})
}

// // compares the output of a gogen.Package with the expected
func comparePackageOutput(t *testing.T, pkg *convert.Package, expect string) {
	t.Helper()
	// For Test,The Test package's header filename same as package name
	var buf bytes.Buffer
	err := pkg.Pkg().WriteTo(&buf, "temp.go")
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
	nc := cltest.NC(&llcppg.Config{}, nil, cltest.NewConvSym(
		llcppg.SymbolInfo{
			Mangle: "Func1",
			CPP:    "Func1",
			Go:     "Func1",
		},
		llcppg.SymbolInfo{
			Mangle: "Func2",
			CPP:    "Func2",
			Go:     "Func2",
		},
	))
	pkg, err := createTestPkg(nc, &convert.PackageConfig{
		OutputDir: "",
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
				pkg.NewTypeDecl("Foo1", &ast.TypeDecl{
					Object: ast.Object{
						Name: &ast.Ident{Name: "Foo1"},
					},
					Type: &ast.RecordType{
						Tag:    ast.Struct,
						Fields: &ast.FieldList{}},
				}, nc)
			},
			headerFile: "/path/to/file1.h",
			incPath:    "file1.h",
			newType:    "Foo1",
		},
		{
			addType: func() {
				pkg.NewTypedefDecl("Bar2", &ast.TypedefDecl{
					Object: ast.Object{
						Name: &ast.Ident{Name: "Bar2"},
					},
					Type: &ast.BuiltinType{Kind: ast.Int},
				}, nc)
			},
			headerFile: "/path/to/file2.h",
			incPath:    "file2.h",
			newType:    "Bar2",
		},
		{
			addType: func() {
				pkg.NewFuncDecl("Func1", &ast.FuncDecl{
					Object: ast.Object{
						Name: &ast.Ident{Name: "Func1"},
					},
					MangledName: "Func1",
					Type:        &ast.FuncType{Params: &ast.FieldList{List: []*ast.Field{}}, Ret: &ast.BuiltinType{Kind: ast.Void}},
				})
			},
			headerFile: "/path/to/file3.h",
			incPath:    "file3.h",
			newType:    "Func1",
		},
	}

	for i, tc := range testCases {
		hfile := &ncimpl.HeaderFile{
			File:     tc.headerFile,
			FileType: llcppg.Inter,
		}
		SetGoFile(pkg, hfile.ToGoFileName(pkgname))
		tc.addType()

		var buf bytes.Buffer
		goFileName := name.HeaderFileToGo(tc.headerFile)
		pkg.Pkg().WriteTo(&buf, goFileName)
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
		mod, err := xgomod.Load(".")
		if err != nil {
			t.Fatal(err)
		}
		deps := []string{
			"github.com/goplus/llcppg/cl/internal/convert/testdata/invalidpath",
			"github.com/goplus/llcppg/cl/internal/convert/testdata/partfinddep",
		}
		p.PkgInfo = convert.NewPkgInfo(".", deps, nil)
		loader := convert.NewPkgDepLoader(mod, genPkg)
		depPkgs, err := loader.LoadDeps(p.PkgInfo)
		p.PkgInfo.Deps = depPkgs
		if err != nil && !errors.Is(err, llcppg.ErrConfig) {
			t.Fatal(err)
		}
		_, err = loader.Import("github.com/goplus/invalidpkg")
		if err == nil {
			t.Fatal("expected error")
		}
	})
	t.Run("invalid pub file", func(t *testing.T) {
		_, err := createTestPkg(nil, &convert.PackageConfig{
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
		_, err := createTestPkg(nil, &convert.PackageConfig{
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
		_, err := createTestPkg(nil, &convert.PackageConfig{
			OutputDir: ".",
			PkgBase: convert.PkgBase{
				Deps: []string{
					"github.com/goplus/llcppg/cl/internal/convert/testdata/cjson",
					"github.com/goplus/llcppg/cl/internal/convert/testdata/cjsonbool",
				},
			},
		})
		if err != nil && !errors.Is(err, llcppg.ErrConfig) {
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
	ncimpl.NewHeaderFile("/path/to/foo.h", 0).ToGoFileName("Pkg")
}

func TestNewPackageLinkWithoutLibCommand(t *testing.T) {
	_, err := convert.NewPackage(nil, &convert.PackageConfig{
		PkgBase: convert.PkgBase{
			PkgPath: ".",
		},
		Name:    pkgname,
		GenConf: &gogen.Config{},
	})
	if err != nil {
		t.Fatal("Unexpect Error")
	}
}
