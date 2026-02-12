package convert

import (
	"bytes"
	goast "go/ast"
	"go/token"
	"go/types"
	"strings"
	"testing"

	"github.com/goplus/gogen"
	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/cl/internal/cltest"
	"github.com/goplus/llcppg/cl/nc"
	"github.com/goplus/llcppg/cl/nc/ncimpl"
	llcppg "github.com/goplus/llcppg/config"
	"github.com/goplus/llcppg/internal/name"
	ctoken "github.com/goplus/llcppg/token"
)

func emptyPkg(nc nc.NodeConverter) *Package {
	pnc := nc
	if pnc == nil {
		pnc = cltest.NC(&llcppg.Config{}, nil, cltest.NewConvSym())
	}
	pkg, err := NewPackage(pnc, &PackageConfig{
		PkgBase: PkgBase{
			PkgPath: ".",
			Pubs:    make(map[string]string),
		},
		Name:       "testpkg",
		GenConf:    &gogen.Config{},
		OutputDir:  "",
		LibCommand: "${pkg-config --libs xxx}",
	})
	if err != nil {
		panic(err)
	}
	return pkg
}

func TestTypeRefIncompleteFail(t *testing.T) {
	t.Run("ref tag incomplete fail", func(t *testing.T) {
		pkg := emptyPkg(nil)
		tempFile := &ncimpl.HeaderFile{
			File:     "temp.h",
			FileType: llcppg.Inter,
		}
		pkg.p.SetCurFile(tempFile.ToGoFileName("testpkg"), true)
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("Expected panic, got nil")
			}
		}()
		pkg.handleTyperefIncomplete(&ast.TagExpr{
			Tag: 0,
			Name: &ast.ScopingExpr{
				X: &ast.Ident{Name: "Bar"},
			},
		}, nil, "NewBar")
	})
}

func TestRedefPubName(t *testing.T) {
	pkg := emptyPkg(nil)
	tempFile := &ncimpl.HeaderFile{
		File:     "temp.h",
		FileType: llcppg.Inter,
	}
	pkg.p.SetCurFile(tempFile.ToGoFileName("testpkg"), true)
	// mock a function name which is not register in processsymbol
	pkg.p.NewFuncDecl(token.NoPos, "Foo", types.NewSignatureType(nil, nil, nil, types.NewTuple(), types.NewTuple(), false))
	pkg.p.NewFuncDecl(token.NoPos, "Bar", types.NewSignatureType(nil, nil, nil, types.NewTuple(), types.NewTuple(), false))
	t.Run("enum type redefine pubname", func(t *testing.T) {
		err := pkg.NewEnumTypeDecl("Foo", &ast.EnumTypeDecl{
			Object: ast.Object{
				Loc:  &ast.Location{File: "temp.h"},
				Name: nil,
			},
			Type: &ast.EnumType{
				Items: []*ast.EnumItem{
					{Name: &ast.Ident{Name: "Foo"}, Value: &ast.BasicLit{Kind: ast.IntLit, Value: "0"}},
				},
			},
		}, cltest.NC(&llcppg.Config{}, nil, cltest.NewConvSym()))
		if err == nil {
			t.Fatal("expect a error")
		}
	})
	t.Run("macro redefine pubname", func(t *testing.T) {
		err := pkg.NewMacro("Bar", &ast.Macro{
			Loc:    &ast.Location{File: "temp.h"},
			Name:   "Bar",
			Tokens: []*ast.Token{{Token: ctoken.IDENT, Lit: "Bar"}, {Token: ctoken.LITERAL, Lit: "1"}},
		})
		if err == nil {
			t.Fatal("expect a error")
		}
	})
}

func TestPubMethodName(t *testing.T) {
	name := types.NewTypeName(0, nil, "Foo", nil)
	named := types.NewNamed(name, nil, nil)
	ptrRecv := types.NewPointer(named)
	fnName := "Foo"
	pubName := pubMethodName(ptrRecv, &GoFuncSpec{GoSymbName: fnName, FnName: fnName, PtrRecv: true, IsMethod: true})
	if pubName != "(*Foo).Foo" {
		t.Fatal("Expected pubName to be '(*Foo).Foo', got", pubName)
	}
	valRecv := named
	pubName = pubMethodName(valRecv, &GoFuncSpec{GoSymbName: fnName, FnName: fnName, IsMethod: true})
	if pubName != "Foo.Foo" {
		t.Fatal("Expected pubName to be 'Foo.Foo', got", pubName)
	}

	unknownRecv := types.NewStruct(nil, []string{})
	pubName = pubMethodName(unknownRecv, &GoFuncSpec{GoSymbName: fnName, FnName: fnName, IsMethod: false})
	if pubName != "Foo" {
		t.Fatal("Expected pubName to be 'Foo', got", pubName)
	}
}

func TestGetNameType(t *testing.T) {
	named := types.NewNamed(types.NewTypeName(0, nil, "Foo", nil), nil, nil)
	ptrNamed := types.NewPointer(named)
	customSturct := types.NewStruct(nil, nil)

	namedRes := getNamedType(named)
	if namedRes != named {
		t.Fatal("Expected namedRes to be *types.Named, got", namedRes)
	}

	ptrNamedRes := getNamedType(ptrNamed)
	if ptrNamedRes != named {
		t.Fatal("Expected ptrNamedRes to be *types.Named, got", ptrNamedRes)
	}

	customRes := getNamedType(customSturct)
	if customRes != nil {
		t.Fatal("Expected nil, got", customRes)
	}
}

func TestMarkUseFail(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Expected panic, got nil")
		}
	}()
	pkg, err := NewPackage(cltest.NC(&llcppg.Config{}, nil, cltest.NewConvSym()), &PackageConfig{
		PkgBase: PkgBase{
			PkgPath: ".",
			Pubs:    make(map[string]string),
		},
		LibCommand: "${pkg-config --libs xxx}",
	})
	if err != nil {
		t.Fatal("NewPackage failed:", err)
	}
	pkg.markUseDeps(&PkgDepLoader{})
}

func TestProcessSymbol(t *testing.T) {
	toCamel := func(trimprefix []string) ncimpl.NameMethod {
		return func(cname string) string {
			return name.PubName(name.RemovePrefixedName(cname, trimprefix))
		}
	}
	toExport := func(trimprefix []string) ncimpl.NameMethod {
		return func(cname string) string {
			return name.ExportName(name.RemovePrefixedName(cname, trimprefix))
		}
	}
	sym := NewProcessSymbol()

	testCases := []struct {
		name         string
		trimPrefixes []string
		nameMethod   func(trimprefix []string) ncimpl.NameMethod
		expected     string
		expectChange bool
	}{
		{"lua_closethread", []string{"lua_", "luaL_"}, toCamel, "Closethread", true},
		{"luaL_checknumber", []string{"lua_", "luaL_"}, toCamel, "Checknumber", true},
		{"_gmp_err", []string{}, toCamel, "X_gmpErr", true},
		{"fn_123illegal", []string{"fn_"}, toCamel, "X123illegal", true},
		{"fts5_tokenizer", []string{}, toCamel, "Fts5Tokenizer", true},
		{"Fts5Tokenizer", []string{}, toCamel, "Fts5Tokenizer__1", true},
		{"normal_var", []string{}, toExport, "Normal_var", true},
		{"Cameled", []string{}, toExport, "Cameled", false},
	}
	for _, tc := range testCases {
		pubName := sym.Register(Node{name: tc.name, kind: TypeDecl}, tc.expected)
		if pubName != tc.expected {
			t.Errorf("Expected %s, but got %s", tc.expected, pubName)
		}
		if tc.expectChange && pubName == tc.name {
			t.Errorf("Expected Change, but got same name")
		}
	}
}

func TestNoEmptyConstGroupWhenAllEnumItemsSkipped(t *testing.T) {
	pnc := cltest.NC(&llcppg.Config{}, nil, cltest.NewConvSym())
	pkg := emptyPkg(pnc)
	tempFile := &ncimpl.HeaderFile{
		File:     "temp.h",
		FileType: llcppg.Inter,
	}
	pkg.p.SetCurFile(tempFile.ToGoFileName("testpkg"), true)

	items := []*ast.EnumItem{
		{Name: &ast.Ident{Name: "Red"}, Value: &ast.BasicLit{Kind: ast.IntLit, Value: "0"}},
		{Name: &ast.Ident{Name: "Green"}, Value: &ast.BasicLit{Kind: ast.IntLit, Value: "1"}},
	}

	// First enum: registers the items normally
	err := pkg.NewEnumTypeDecl("Color", &ast.EnumTypeDecl{
		Object: ast.Object{
			Loc:  &ast.Location{File: "temp.h"},
			Name: &ast.Ident{Name: "Color"},
		},
		Type: &ast.EnumType{Items: items},
	}, pnc)
	if err != nil {
		t.Fatal(err)
	}

	// Second enum: all items already registered, so all are skipped.
	// With lazy ConstGroup, no empty const() should be created.
	err = pkg.NewEnumTypeDecl("Color2", &ast.EnumTypeDecl{
		Object: ast.Object{
			Loc:  &ast.Location{File: "temp.h"},
			Name: &ast.Ident{Name: "Color2"},
		},
		Type: &ast.EnumType{Items: items},
	}, pnc)
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	err = pkg.p.WriteTo(&buf, "temp.go")
	if err != nil {
		t.Fatal(err)
	}
	output := buf.String()

	// There should be exactly one const block from the first enum.
	// No empty const() from the second enum.
	constCount := strings.Count(output, "const (")
	if constCount != 1 {
		t.Errorf("expected exactly 1 const block, got %d.\nOutput:\n%s", constCount, output)
	}

	// Also verify at the AST level: no const GenDecl with empty Specs
	// should exist in the package's corresponding ast.File.
	goFile := pkg.p.ASTFile("temp.go")
	if goFile == nil {
		t.Fatal("expected ast.File for temp.go, got nil")
	}
	constDeclCount := 0
	for _, decl := range goFile.Decls {
		genDecl, ok := decl.(*goast.GenDecl)
		if !ok || genDecl.Tok != token.CONST {
			continue
		}
		constDeclCount++
		if len(genDecl.Specs) == 0 {
			t.Error("found empty const declaration in AST (no Specs)")
		}
	}
	if constDeclCount != 1 {
		t.Errorf("expected exactly 1 const GenDecl in AST, got %d", constDeclCount)
	}
}
