package ncimpl

import (
	"errors"
	"testing"

	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/cl/nc"
	llconfig "github.com/goplus/llcppg/config"
)

func TestHeaderFile(t *testing.T) {
	testCases := []struct {
		name     string
		file     string
		fileType llconfig.FileType
		pkgName  string
		expected string
		inCurPkg bool
	}{
		{
			name:     "Inter file",
			file:     "/path/to/header.h",
			fileType: llconfig.Inter,
			pkgName:  "testpkg",
			expected: "header.go",
			inCurPkg: true,
		},
		{
			name:     "Impl file",
			file:     "/path/to/impl.h",
			fileType: llconfig.Impl,
			pkgName:  "testpkg",
			expected: "testpkg_autogen.go",
			inCurPkg: true,
		},
		{
			name:     "Third file",
			file:     "/path/to/third.h",
			fileType: llconfig.Third,
			pkgName:  "testpkg",
			expected: "testpkg_autogen.go",
			inCurPkg: false,
		},
		{
			name:     "Underscore file",
			file:     "/path/to/_types.h",
			fileType: llconfig.Inter,
			pkgName:  "testpkg",
			expected: "X_types.go",
			inCurPkg: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hf := NewHeaderFile(tc.file, tc.fileType)
			result := hf.ToGoFileName(tc.pkgName)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
			if hf.InCurPkg() != tc.inCurPkg {
				t.Errorf("Expected InCurPkg to be %v, got %v", tc.inCurPkg, hf.InCurPkg())
			}
		})
	}
}

func TestThirdTypeLoc(t *testing.T) {
	locMap := &ThirdTypeLoc{}
	ident1 := &ast.Ident{Name: "TypeA"}
	loc1 := &ast.Location{File: "/path/to/typeA.h"}
	locMap.Add(ident1, loc1)

	file, ok := locMap.Lookup("TypeA")
	if !ok {
		t.Error("Expected to find TypeA, but didn't")
	}
	if file != "/path/to/typeA.h" {
		t.Errorf("Expected /path/to/typeA.h, got %s", file)
	}

	_, ok = locMap.Lookup("TypeB")
	if ok {
		t.Error("Unexpectedly found TypeB, which shouldn't exist")
	}

	locMap.Add(ident1, &ast.Location{File: "/different/path.h"})
	cvt := &Converter{
		locMap: *locMap,
	}
	file, _ = cvt.Lookup("TypeA")
	if file != "/path/to/typeA.h" {
		t.Error("Location for TypeA should not have changed")
	}
}

func TestConverterNameTransformation(t *testing.T) {
	converter := &Converter{
		PkgName:        "testpkg",
		TrimPrefixes:   []string{"prefix_"},
		KeepUnderScore: false,
		Pubs:           map[string]string{"predefined": "CustomName", "KEEP": ""},
	}

	testCases := []struct {
		name     string
		input    string
		expected string
		method   string
	}{
		{"Simple declName", "simple_name", "SimpleName", "declName"},
		{"With prefix declName", "prefix_name", "Name", "declName"},
		{"With underscore declName", "_hidden", "X_hidden", "declName"},
		{"Predefined declName", "predefined", "CustomName", "declName"},
		{"Keep Origin Name", "KEEP", "KEEP", "declName"},

		{"Simple constName", "SIMPLE_NAME", "SIMPLE_NAME", "constName"},
		{"With prefix constName", "prefix_NAME", "NAME", "constName"},
		{"With underscore constName", "_hidden", "X_hidden", "constName"},
		{"Keep Origin Name", "KEEP", "KEEP", "constName"},
		{"Predefined constName", "predefined", "CustomName", "constName"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var result string
			if tc.method == "declName" {
				result = converter.declName(tc.input)
			} else {
				result = converter.constName(tc.input)
			}

			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

const interFile = "/path/to/inter.h"
const implFile = "/path/to/impl.h"
const thirdFile = "/path/to/third.h"

var fileMap = map[string]*llconfig.FileInfo{
	interFile: {FileType: llconfig.Inter},
	implFile:  {FileType: llconfig.Impl},
	thirdFile: {FileType: llconfig.Third},
}

func TestConverterConvFile(t *testing.T) {
	converter := &Converter{
		PkgName: "testpkg",
		FileMap: fileMap,
	}

	testCases := []struct {
		name     string
		file     string
		expected string
		ok       bool
	}{
		{"Inter file", interFile, "inter.go", true},
		{"Impl file", implFile, "testpkg_autogen.go", true},
		{"Third file", thirdFile, "testpkg_autogen.go", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			goFile, ok := converter.convFile(tc.file, nil)

			if goFile.FileName != tc.expected {
				t.Errorf("Expected file %s, got %s", tc.expected, goFile.FileName)
			}

			if ok != tc.ok {
				t.Errorf("Expected ok to be %v, got %v", tc.ok, ok)
			}
		})
	}

	t.Run("Missing file", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic")
			}
		}()

		converter.convFile("/path/to/missing.h", nil)
	})
}

func mockSymConv(name *ast.Object, mangleName string) (string, error) {
	if mangleName == "validFunc" {
		return "ValidFunc", nil
	}
	return "", errors.New("symbol not found")
}

func TestConverterIsPublic(t *testing.T) {
	converter := &Converter{
		KeepUnderScore: false,
	}

	if converter.IsPublic("_hidden") {
		t.Error("Unexpectedly found public when keepUnderScore is false")
	}

	if !converter.IsPublic("public") {
		t.Error("Unexpectedly found private when keepUnderScore is false")
	}

	converter.KeepUnderScore = true

	if !converter.IsPublic("_hidden") {
		t.Error("Unexpectedly found private when keepUnderScore is true")
	}
}

func TestConv(t *testing.T) {
	testCases := []struct {
		name         string
		convNode     ast.Node // ast.Macro,ast.Decl...
		pubs         map[string]string
		convSym      func(name *ast.Object, mangleName string) (string, error)
		trimPrefixes []string
		expectFile   string
		expectName   string
		exec         func(conv *Converter, convNode ast.Node) (string, *nc.GoFile, error)
		expectErr    error
	}{
		{
			name: "Macro",
			convNode: &ast.Macro{
				Name: "MACRO_NAME",
				Loc:  &ast.Location{File: interFile},
			},
			exec: func(conv *Converter, convNode ast.Node) (string, *nc.GoFile, error) {
				return conv.ConvMacro(interFile, convNode.(*ast.Macro))
			},
			expectFile: "inter.go",
			expectName: "MACRO_NAME",
			expectErr:  nil,
		},
		{
			name: "Macro Without Prefix",
			convNode: &ast.Macro{
				Name: "MACRO_NAME",
				Loc:  &ast.Location{File: interFile},
			},
			exec: func(conv *Converter, convNode ast.Node) (string, *nc.GoFile, error) {
				return conv.ConvMacro(interFile, convNode.(*ast.Macro))
			},
			trimPrefixes: []string{"MACRO_"},
			expectFile:   "inter.go",
			expectName:   "NAME",
			expectErr:    nil,
		},
		{
			name: "Predefined macro",
			convNode: &ast.Macro{
				Name: "MACRO_CONST",
				Loc:  &ast.Location{File: interFile},
			},
			pubs: map[string]string{"MACRO_CONST": "CustomMacro"},
			exec: func(conv *Converter, convNode ast.Node) (string, *nc.GoFile, error) {
				return conv.ConvMacro(interFile, convNode.(*ast.Macro))
			},
			expectFile: "inter.go",
			expectName: "CustomMacro",
			expectErr:  nil,
		},
		{
			name: "Not Cur Pkg Macro",
			convNode: &ast.Macro{
				Name: "MACRO_NAME",
				Loc:  &ast.Location{File: thirdFile},
			},
			exec: func(conv *Converter, convNode ast.Node) (string, *nc.GoFile, error) {
				return conv.ConvMacro(thirdFile, convNode.(*ast.Macro))
			},
			expectFile: "testpkg_autogen.go",
			expectErr:  nc.ErrSkip,
		},
		{
			name: "Function",
			convNode: &ast.FuncDecl{
				Object: ast.Object{
					Name: &ast.Ident{Name: "func_name"},
					Loc:  &ast.Location{File: interFile},
				},
				MangledName: "validFunc",
			},
			convSym: mockSymConv,
			exec: func(conv *Converter, convNode ast.Node) (string, *nc.GoFile, error) {
				return conv.ConvDecl(interFile, convNode.(*ast.FuncDecl))
			},
			expectFile: "inter.go",
			expectName: "ValidFunc",
			expectErr:  nil,
		},
		{
			name: "Function in third file",
			convNode: &ast.FuncDecl{
				Object: ast.Object{
					Name: &ast.Ident{Name: "func_name"},
					Loc:  &ast.Location{File: thirdFile},
				},
			},
			exec: func(conv *Converter, convNode ast.Node) (string, *nc.GoFile, error) {
				return conv.ConvDecl(thirdFile, convNode.(*ast.FuncDecl))
			},
			expectFile: "testpkg_autogen.go",
			expectErr:  nc.ErrSkip,
		},
		{
			name: "Function No Symbol",
			convNode: &ast.FuncDecl{
				Object: ast.Object{
					Name: &ast.Ident{Name: "func_name"},
					Loc:  &ast.Location{File: "/path/to/inter.h"},
				},
				MangledName: "noSymbolFunc",
			},
			exec: func(conv *Converter, convNode ast.Node) (string, *nc.GoFile, error) {
				return conv.ConvDecl(interFile, convNode.(*ast.FuncDecl))
			},
			convSym:    mockSymConv,
			expectFile: "inter.go",
			expectErr:  nc.ErrSkip,
		},
		{
			name: "Enum",
			convNode: &ast.EnumTypeDecl{
				Object: ast.Object{
					Name: &ast.Ident{Name: "enum_name"},
					Loc:  &ast.Location{File: interFile},
				},
			},
			exec: func(conv *Converter, convNode ast.Node) (string, *nc.GoFile, error) {
				return conv.ConvDecl(interFile, convNode.(*ast.EnumTypeDecl))
			},
			expectFile: "inter.go",
			expectName: "EnumName",
			expectErr:  nil,
		},
		{
			name: "Anony Enum",
			convNode: &ast.EnumTypeDecl{
				Object: ast.Object{
					Name: nil,
					Loc:  &ast.Location{File: interFile},
				},
			},
			exec: func(conv *Converter, convNode ast.Node) (string, *nc.GoFile, error) {
				return conv.ConvDecl(interFile, convNode.(*ast.EnumTypeDecl))
			},
			expectFile: "inter.go",
			expectErr:  nil,
		},
		{
			name: "Type",
			convNode: &ast.TypeDecl{
				Object: ast.Object{
					Name: &ast.Ident{Name: "type_name"},
					Loc:  &ast.Location{File: interFile},
				},
			},
			exec: func(conv *Converter, convNode ast.Node) (string, *nc.GoFile, error) {
				return conv.ConvDecl(interFile, convNode.(*ast.TypeDecl))
			},
			expectFile: "inter.go",
			expectName: "TypeName",
			expectErr:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			converter := &Converter{
				PkgName:      "testpkg",
				FileMap:      fileMap,
				Pubs:         tc.pubs,
				ConvSym:      tc.convSym,
				TrimPrefixes: tc.trimPrefixes,
			}
			goName, goFile, err := tc.exec(converter, tc.convNode)
			if goName != tc.expectName {
				t.Errorf("Expected %s, got %s", tc.expectName, goName)
			}
			if goFile.FileName != tc.expectFile {
				t.Errorf("Expected %s, got %s", tc.expectFile, goFile.FileName)
			}
			if err != tc.expectErr {
				t.Errorf("Expected %v, got %v", tc.expectErr, err)
			}
		})
	}
}

func TestConvEnumItem(t *testing.T) {
	cvt := &Converter{
		PkgName: "testpkg",
		FileMap: fileMap,
		ConvSym: mockSymConv,
	}
	goName, err := cvt.ConvEnumItem(&ast.EnumTypeDecl{
		Object: ast.Object{
			Name: &ast.Ident{Name: "enum_name"},
			Loc:  &ast.Location{File: interFile},
		},
	}, &ast.EnumItem{
		Name: &ast.Ident{Name: "ENUM_ITEM_NAME"},
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if goName != "ENUM_ITEM_NAME" {
		t.Errorf("Expected EnumItem, got %s", goName)
	}
}

func TestConvTagExpr(t *testing.T) {
	cvt := &Converter{
		PkgName: "testpkg",
		FileMap: fileMap,
		ConvSym: mockSymConv,
	}
	goName := cvt.ConvTagExpr("type_name")
	if goName != "TypeName" {
		t.Errorf("Expected TypeName, got %s", goName)
	}
}
