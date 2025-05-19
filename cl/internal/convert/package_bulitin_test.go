package convert

import (
	"go/token"
	"go/types"
	"testing"

	"github.com/goplus/gogen"
	"github.com/goplus/llcppg/_xtool/llcppsymg/tool/name"
	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/cl/internal/cltest"
	llcppg "github.com/goplus/llcppg/config"
	ctoken "github.com/goplus/llcppg/token"
	"github.com/goplus/mod/gopmod"
)

func emptyPkg() *Package {
	mod, err := gopmod.Load(".")
	if err != nil {
		panic(err)
	}
	pkg, err := NewPackage(&PackageConfig{
		PkgBase: PkgBase{
			PkgPath: ".",
			Pubs:    make(map[string]string),
		},
		Name:       "testpkg",
		Mod:        mod,
		GenConf:    &gogen.Config{},
		ConvSym:    cltest.NewConvSym(),
		LibCommand: "${pkg-config --libs xxx}",
	})
	if err != nil {
		panic(err)
	}
	return pkg
}

func TestTypeRefIncompleteFail(t *testing.T) {
	pkg := emptyPkg()
	tempFile := &HeaderFile{
		File:     "temp.h",
		FileType: llcppg.Inter,
	}
	pkg.SetCurFile(tempFile)

	t.Run("defer write third type not found", func(t *testing.T) {
		pkg.locMap.Add(&ast.Ident{Name: "Bar"}, &ast.Location{File: "Bar"})
		pkg.incompleteTypes.Add(&Incomplete{cname: "Bar"})
		err := pkg.NewTypedefDecl(&ast.TypedefDecl{
			Object: ast.Object{
				Name: &ast.Ident{Name: "Foo"},
			},
			Type: &ast.TagExpr{
				Name: &ast.Ident{Name: "Bar"},
			},
		})
		if err != nil {
			t.Fatal("NewTypedefDecl failed:", err)
		}
		pkg.incompleteTypes.Complete("Bar")
		err = pkg.Complete()
		if err == nil {
			t.Fatal("expect a error")
		}
	})
	t.Run("ref tag incomplete fail", func(t *testing.T) {
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
	pkg := emptyPkg()
	tempFile := &HeaderFile{
		File:     "temp.h",
		FileType: llcppg.Inter,
	}
	pkg.SetCurFile(tempFile)
	// mock a function name which is not register in processsymbol
	pkg.p.NewFuncDecl(token.NoPos, "Foo", types.NewSignatureType(nil, nil, nil, types.NewTuple(), types.NewTuple(), false))
	pkg.p.NewFuncDecl(token.NoPos, "Bar", types.NewSignatureType(nil, nil, nil, types.NewTuple(), types.NewTuple(), false))
	t.Run("enum type redefine pubname", func(t *testing.T) {
		err := pkg.NewEnumTypeDecl(&ast.EnumTypeDecl{
			Object: ast.Object{
				Loc:  &ast.Location{File: "temp.h"},
				Name: nil,
			},
			Type: &ast.EnumType{
				Items: []*ast.EnumItem{
					{Name: &ast.Ident{Name: "Foo"}, Value: &ast.BasicLit{Kind: ast.IntLit, Value: "0"}},
				},
			},
		})
		if err == nil {
			t.Fatal("expect a error")
		}
	})
	t.Run("macro redefine pubname", func(t *testing.T) {
		err := pkg.NewMacro(&ast.Macro{
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

func TestTrimPrefixes(t *testing.T) {
	mod, err := gopmod.Load(".")
	if err != nil {
		panic(err)
	}
	pkg, err := NewPackage(&PackageConfig{
		PkgBase: PkgBase{
			PkgPath: ".",
			Pubs:    make(map[string]string),
		},
		Name:         "testpkg",
		GenConf:      &gogen.Config{},
		ConvSym:      cltest.NewConvSym(),
		TrimPrefixes: []string{"prefix1", "prefix2"},
		LibCommand:   "${pkg-config --libs xxx}",
		Mod:          mod,
	})
	if err != nil {
		t.Fatal("NewPackage failed:", err)
	}

	pkg.curFile = &HeaderFile{
		FileType: llcppg.Inter,
	}

	result := pkg.trimPrefixes()
	expected := []string{"prefix1", "prefix2"}
	if len(result) != len(expected) || (len(result) > 0 && result[0] != expected[0]) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	pkg.curFile.FileType = llcppg.Third
	result = pkg.trimPrefixes()
	if len(result) != 0 {
		t.Errorf("Expected Empty TrimPrefix")
	}
}

func TestMarkUseFail(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Expected panic, got nil")
		}
	}()
	pkg, err := NewPackage(&PackageConfig{
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
	toCamel := func(trimprefix []string) NameMethod {
		return func(cname string) string {
			return name.PubName(name.RemovePrefixedName(cname, trimprefix))
		}
	}
	toExport := func(trimprefix []string) NameMethod {
		return func(cname string) string {
			return name.ExportName(name.RemovePrefixedName(cname, trimprefix))
		}
	}
	sym := NewProcessSymbol()

	testCases := []struct {
		name         string
		trimPrefixes []string
		nameMethod   func(trimprefix []string) NameMethod
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
