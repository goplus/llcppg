package convert

import (
	"errors"
	"go/types"
	"testing"

	"github.com/goplus/gogen"
	"github.com/goplus/llcppg/ast"
	cfg "github.com/goplus/llcppg/cmd/gogensig/config"
	"github.com/goplus/llcppg/llcppg"
)

func TestTypeRefIncompleteFail(t *testing.T) {
	pkg := NewPackage(&PackageConfig{
		PkgBase: PkgBase{
			PkgPath:  ".",
			CppgConf: &llcppg.Config{},
			Pubs:     make(map[string]string),
		},
		Name:        "testpkg",
		GenConf:     &gogen.Config{},
		OutputDir:   "",
		SymbolTable: cfg.CreateSymbolTable([]cfg.SymbolEntry{}),
	})
	tempFile := &HeaderFile{
		File:     "temp.h",
		FileType: llcppg.Inter,
	}
	pkg.SetCurFile(tempFile)

	t.Run("write pkg fail", func(t *testing.T) {
		pkg.incompleteTypes.Add(&Incomplete{cname: "Bar", file: tempFile, getType: func() (types.Type, error) {
			return nil, errors.New("Mock Err")
		}})
		err := pkg.WritePkgFiles()
		if err == nil {
			t.Fatal("Expect Error")
		}
		pkg.incompleteTypes.Complete("Bar")
	})

	t.Run("defer write third type not found", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("Expected panic, got nil")
			}
		}()
		pkg.locMap.Add(&ast.Ident{Name: "Bar"}, &ast.Location{File: "Bar"})
		pkg.incompleteTypes.Add(&Incomplete{cname: "Bar"})
		err := pkg.NewTypedefDecl(&ast.TypedefDecl{
			Name: &ast.Ident{Name: "Foo"},
			Type: &ast.TagExpr{
				Name: &ast.Ident{Name: "Bar"},
			},
		})
		if err != nil {
			t.Fatal("NewTypedefDecl failed:", err)
		}
		pkg.incompleteTypes.Complete("Bar")
		pkg.WritePkgFiles()
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
	pkg := NewPackage(&PackageConfig{
		PkgBase: PkgBase{
			PkgPath: ".",
			CppgConf: &llcppg.Config{
				TrimPrefixes: []string{"prefix1", "prefix2"},
			},
			Pubs: make(map[string]string),
		},
		Name:        "testpkg",
		GenConf:     &gogen.Config{},
		OutputDir:   "",
		SymbolTable: &cfg.SymbolTable{},
	})

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
