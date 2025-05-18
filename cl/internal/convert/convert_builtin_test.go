package convert

import (
	"errors"
	"go/types"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/cl/internal/cltest"
	"github.com/goplus/llcppg/cmd/gogensig/config"
	llcppg "github.com/goplus/llcppg/config"
)

func basicConverter() *Converter {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	tempDir, err := os.MkdirTemp(dir, "test_package_write_unwritable")
	if err != nil {
		panic(err)
	}

	cfg := &llcppg.Config{
		Libs: "${pkg-config --libs xxx}",
	}
	// todo: remove this,convert not read llcppg.cfg directly
	cfgPath := filepath.Join(tempDir, llcppg.LLCPPG_CFG)
	err = config.CreateJSONFile(cfgPath, cfg)
	if err != nil {
		panic(err)
	}

	converter, err := NewConverter(&Config{
		PkgName:   "test",
		ConvSym:   cltest.NewConvSym(),
		CfgFile:   cfgPath,
		OutputDir: tempDir,
		Pkg: &llcppg.Pkg{
			File: &ast.File{
				Decls: []ast.Decl{},
			},
			FileMap: map[string]*llcppg.FileInfo{},
		},
	})
	if err != nil {
		panic(err)
	}
	return converter
}

func TestPkgFail(t *testing.T) {
	converter := basicConverter()
	defer os.RemoveAll(converter.GenPkg.conf.OutputDir)
	t.Run("ProcessFail", func(t *testing.T) {
		defer func() {
			checkPanic(t, recover(), "File \"noexist.h\" not found in FileMap")
		}()
		converter.Pkg.File.Decls = append(converter.Pkg.File.Decls, &ast.TypeDecl{
			Object: ast.Object{
				Loc: &ast.Location{
					File: "noexist.h",
				},
			},
		})
		converter.Pkg.FileMap["exist.h"] = &llcppg.FileInfo{
			FileType: llcppg.Inter,
		}
		converter.Process()
	})

	t.Run("Complete fail", func(t *testing.T) {
		defer func() {
			checkPanic(t, recover(), "Complete Fail: Mock Err")
		}()
		converter.GenPkg.incompleteTypes.Add(&Incomplete{cname: "Bar", file: &HeaderFile{
			File:     "temp.h",
			FileType: llcppg.Inter,
		}, getType: func() (types.Type, error) {
			return nil, errors.New("Mock Err")
		}})
		converter.Complete()
	})
}

func TestProcessWithError(t *testing.T) {
	defer func() {
		checkPanic(t, recover(), "NewTypedefDecl: Foo fail")
	}()
	converter := basicConverter()
	converter.GenPkg.conf.ConvSym = cltest.NewConvSym(cltest.SymbolEntry{
		CppName:    "Foo",
		MangleName: "Foo",
		GoName:     "Foo",
	})
	declLoc := &ast.Location{
		File: "exist.h",
	}
	converter.Pkg.File.Decls = []ast.Decl{
		&ast.FuncDecl{
			Object: ast.Object{
				Loc: declLoc,
				Name: &ast.Ident{
					Name: "Foo",
				},
			},
			MangledName: "Foo",
			Type: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
						{Type: &ast.Ident{Name: "int"}},
					},
				},
				Ret: &ast.Ident{Name: "int"},
			},
		},
		&ast.TypedefDecl{
			Object: ast.Object{
				Loc: declLoc,
				Name: &ast.Ident{
					Name: "Foo",
				},
			},
			Type: &ast.Ident{
				Name: "Foo",
			},
		},
	}
	converter.Pkg.FileMap["exist.h"] = &llcppg.FileInfo{
		FileType: llcppg.Inter,
	}
	converter.Process()
}

func checkPanic(t *testing.T, r interface{}, expectedPrefix string) {
	if r == nil {
		t.Errorf("Expected panic, but got: %v", r)
	} else {
		if !strings.HasPrefix(r.(string), expectedPrefix) {
			t.Errorf("Expected panic %s, but got: %v", expectedPrefix, r)
		}
	}
}
