package convert

import (
	"errors"
	"go/types"
	"os"
	"strings"
	"testing"

	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/llcppg"
)

func TestPkgFail(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal("Getwd failed:", err)
	}
	tempDir, err := os.MkdirTemp(dir, "test_package_write_unwritable")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	converter, err := NewConverter(&Config{
		PkgName:   "test",
		SymbFile:  "",
		CfgFile:   "",
		OutputDir: tempDir,
		Pkg: &llcppg.Pkg{
			File: &ast.File{
				Decls: []ast.Decl{},
			},
			FileMap: map[string]*llcppg.FileInfo{},
		},
	})
	t.Run("FmtFail", func(t *testing.T) {
		defer func() {
			checkPanic(t, recover(), "go fmt:")
		}()
		converter.Fmt()
	})
	t.Run("ProcessFail", func(t *testing.T) {
		defer func() {
			checkPanic(t, recover(), "File \"noexist.h\" not found in FileMap")
		}()
		converter.Pkg.File.Decls = append(converter.Pkg.File.Decls, &ast.TypeDecl{
			DeclBase: ast.DeclBase{
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
	t.Run("WriteLinkFileFail", func(t *testing.T) {
		defer func() {
			checkPanic(t, recover(), "WriteLinkFile:")
		}()
		converter.Write()
	})
	t.Run("WritePubFileFail", func(t *testing.T) {
		defer func() {
			checkPanic(t, recover(), "WritePubFile:")
		}()
		converter.GenPkg.conf.OutputDir = "/nonexistent_directory/test.txt"
		converter.GenPkg.Pubs = map[string]string{"test": "Test"}
		converter.Write()
	})
	t.Run("WritePkgFilesFail", func(t *testing.T) {
		defer func() {
			checkPanic(t, recover(), "WritePkgFiles:")
		}()
		converter.GenPkg.incompleteTypes.Add(&Incomplete{cname: "Bar", file: &HeaderFile{
			File:     "/path/to/temp.go",
			FileType: llcppg.Inter,
		}, getType: func() (types.Type, error) {
			return nil, errors.New("Mock Err")
		}})
		if err != nil {
			t.Fatal("NewAstConvert Fail")
		}
		converter.Write()
	})
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
