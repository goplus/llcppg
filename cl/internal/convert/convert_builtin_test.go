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

func basicConverter(nc nc.NodeConverter) *Converter {
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
	converter, err := NewConverter(&Config{
		PkgPath:   ".",
		PkgName:   "test",
		OutputDir: tempDir,
		Pkg: &ast.File{
			Decls: []ast.Decl{},
		},
		NC:   nodeConverter,
		Deps: cfg.Deps,
		Libs: cfg.Libs,
	})
	if err != nil {
		panic(err)
	}
	return converter
}

func TestPkgFail(t *testing.T) {
	converter := basicConverter(nil)
	defer os.RemoveAll(converter.GenPkg.conf.OutputDir)

	/* todo(zzy): move to NC to test name fetch
	t.Run("ProcessFail", func(t *testing.T) {
		defer func() {
			checkPanic(t, recover(), "File \"noexist.h\" not found in FileMap")
		}()
		converter.Pkg.Decls = append(converter.Pkg.Decls, &ast.TypeDecl{
			Object: ast.Object{
				Loc: &ast.Location{
					File: "noexist.h",
				},
			},
		})
		converter.FileMap["exist.h"] = &llcppg.FileInfo{
			FileType: llcppg.Inter,
		}
		converter.Process()
	})
	*/

	t.Run("Complete fail", func(t *testing.T) {
		defer func() {
			checkPanic(t, recover(), "Complete Fail: Mock Err")
		}()
		ctx := converter.GenPkg
		ctx.p.SetCurFile("temp.go", true)
		ctx.incompleteTypes.Add(&Incomplete{cname: "Bar", file: ctx.p.CurFile(), getType: func() (types.Type, error) {
			return nil, errors.New("Mock Err")
		}})
		converter.Complete()
	})
}

func TestProcessWithError(t *testing.T) {
	defer func() {
		checkPanic(t, recover(), "NewTypedefDecl: Foo fail")
	}()
	converter := basicConverter(cltest.NC(&llcppg.Config{},
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
	declLoc := &ast.Location{
		File: "exist.h",
	}
	converter.Pkg.Decls = []ast.Decl{
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
