package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/goplus/llcppg/_xtool/llcppsigfetch/parse"
	test "github.com/goplus/llcppg/_xtool/llcppsigfetch/parse/cvt_test"
	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/llcppg"
	"github.com/goplus/llgo/c"
)

func main() {
	TestDefine()
	TestSystemHeader()
	TestInclusionMap()
	TestMacroExpansionOtherFile()
}

func TestDefine() {
	testCases := []string{
		`#define DEBUG`,
		`#define OK 1`,
		`#define SQUARE(x) ((x) * (x))`,
	}
	test.RunTest("TestDefine", testCases)
}

func TestInclusionMap() {
	fmt.Println("=== TestInclusionMap ===")
	context, err := parse.Do(&llcppg.Config{
		Include: []string{"sys.h"},
		CFlags:  "-I./testdata/sysinc",
	})
	if err != nil {
		panic(err)
	}
	found := false
	for _, f := range context.FileSet {
		if strings.HasSuffix(f.Path, "sys/types.h") {
			found = true
		}
	}
	if !found {
		panic("sys/types.h not found")
	} else {
		fmt.Println("sys/types.h include path found")
	}
}

func TestSystemHeader() {
	fmt.Println("=== TestSystemHeader ===")
	temp := path.Join(os.TempDir(), "temp.h")
	os.WriteFile(temp, []byte("#include <stdio.h>"), 0644)
	converter, err := parse.NewConverterX(&parse.Config{
		Cfg: &clangutils.Config{
			File: temp,
		},
		CombinedFile: temp,
	})
	if err != nil {
		panic(err)
	}
	pkg, err := converter.ConvertX()
	if err != nil {
		panic(err)
	}

	for path, info := range pkg.FileMap {
		if !info.IsSys {
			panic(fmt.Errorf("include file is not system header: %s", path))
		}
	}

	for _, decl := range pkg.File.Decls {
		switch decl := decl.(type) {
		case *ast.TypeDecl:
		case *ast.EnumTypeDecl:
		case *ast.FuncDecl:
		case *ast.TypedefDecl:
			if _, ok := pkg.FileMap[decl.DeclBase.Loc.File]; !ok {
				fmt.Println("Decl is not Found in the fileMap", decl.DeclBase.Loc.File)
			}
		}
	}
	fmt.Println("include files are all system headers")
}

func TestMacroExpansionOtherFile() {
	c.Printf(c.Str("TestMacroExpansionOtherFile:\n"))
	test.RunTestWithConfig(&parse.ParseConfig{
		Conf: &llcppg.Config{
			Include: []string{"ref.h"},
			CFlags:  "-I./testdata/macroexpan",
		},
	})
}
