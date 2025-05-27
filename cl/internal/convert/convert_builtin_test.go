package convert

import (
	"errors"
	"go/types"
	"os"
	"strings"
	"testing"

	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/cl/internal/cltest"
	"github.com/goplus/llcppg/cl/nc"
	llcppg "github.com/goplus/llcppg/config"
)

func basicConverter(pkg *ast.File, nc nc.NodeConverter) *Converter {
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
	nodeConverter := nc
	if nodeConverter == nil {
		nodeConverter = cltest.NC(cfg, nil, cltest.NewConvSym())
	}
	p := pkg
	if p == nil {
		p = &ast.File{}
	}
	converter, err := NewConverter(&Config{
		PkgPath:   ".",
		PkgName:   "test",
		OutputDir: tempDir,
		Pkg:       p,
		NC:        nodeConverter,
		Deps:      cfg.Deps,
		Libs:      cfg.Libs,
	})
	if err != nil {
		panic(err)
	}
	return converter
}

func TestPkgFail(t *testing.T) {
	converter := basicConverter(nil, nil)
	defer os.RemoveAll(converter.GenPkg.conf.OutputDir)

	t.Run("ProcessFail", func(t *testing.T) {
		defer func() {
			checkPanic(t, recover(), "File \"noexist.h\" not found in FileMap")
		}()
		pkg := &ast.File{
			Decls: []ast.Decl{
				&ast.TypeDecl{
					Object: ast.Object{
						Loc: &ast.Location{
							File: "noexist.h",
						},
					},
				},
			},
		}
		cvt := basicConverter(pkg, cltest.NC(&llcppg.Config{},
			map[string]*llcppg.FileInfo{
				"exist.h": {
					FileType: llcppg.Inter,
				},
			},
			cltest.NewConvSym(cltest.SymbolEntry{
				CppName:    "Foo",
				MangleName: "Foo",
				GoName:     "Foo",
			}),
		))
		cvt.Convert()
	})

	t.Run("Complete fail", func(t *testing.T) {
		ctx := converter.GenPkg
		ctx.p.SetCurFile("temp.go", true)
		ctx.incompleteTypes.Add(&Incomplete{cname: "Bar", file: ctx.p.CurFile(), getType: func() (types.Type, error) {
			return nil, errors.New("Mock Err")
		}})
		err := converter.Complete()
		checkError(t, err, "Complete Fail: Mock Err")
	})
}

func TestProcessWithError(t *testing.T) {
	declLoc := &ast.Location{
		File: "exist.h",
	}
	pkg := &ast.File{
		Decls: []ast.Decl{
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
		},
	}

	converter := basicConverter(pkg, cltest.NC(&llcppg.Config{},
		map[string]*llcppg.FileInfo{
			"exist.h": {
				FileType: llcppg.Inter,
			},
		},
		cltest.NewConvSym(cltest.SymbolEntry{
			CppName:    "Foo",
			MangleName: "Foo",
			GoName:     "Foo",
		}),
	))

	err := converter.Process()
	checkError(t, err, "NewTypedefDecl: Foo fail")
}

func checkError(t *testing.T, err error, expectedPrefix string) {
	if err == nil {
		t.Fatalf("Expected error, but got nil")
	}
	if !strings.HasPrefix(err.Error(), expectedPrefix) {
		t.Fatalf("Expected error %s, but got: %s", expectedPrefix, err.Error())
	}
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
