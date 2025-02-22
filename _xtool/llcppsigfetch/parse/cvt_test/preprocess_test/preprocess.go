package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/goplus/llcppg/_xtool/llcppsigfetch/parse"
	test "github.com/goplus/llcppg/_xtool/llcppsigfetch/parse/cvt_test"
	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/llcppg"
	"github.com/goplus/llgo/c"
)

func main() {
	TestDefine()
	TestInclude()
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

func TestInclude() {
	testCases := []string{
		`#include "foo.h"`,
		// `#include <limits.h>`, //  Standard libraries are mostly platform-dependent
	}
	test.RunTest("TestInclude", testCases)
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
		if f.IncPath == "sys/types.h" {
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
	context, err := parse.Do(&llcppg.Config{
		Include: []string{"inc.h"},
		CFlags:  "-I./testdata/sysinc",
	})
	if err != nil {
		panic(err)
	}

	if len(context.FileSet) < 2 {
		panic("expect 2 files")
	}
	if context.FileSet[0].IsSys {
		panic("entry file is not system header")
	}

	includePath := context.FileSet[0].Doc.Includes[0].Path
	if strings.HasSuffix(includePath, "stdio.h") && filepath.IsAbs(includePath) {
		fmt.Println("stdio.h is absolute path")
	}

	for i := 1; i < len(context.FileSet); i++ {
		if !context.FileSet[i].IsSys {
			panic(fmt.Errorf("include file is not system header: %s", context.FileSet[i].Path))
		}
		for _, decl := range context.FileSet[i].Doc.Decls {
			switch decl := decl.(type) {
			case *ast.TypeDecl:
			case *ast.EnumTypeDecl:
			case *ast.FuncDecl:
			case *ast.TypedefDecl:
				if decl.DeclBase.Loc.File != context.FileSet[i].Path {
					fmt.Println("Decl is not in the file", decl.DeclBase.Loc.File, "expect", context.FileSet[i].Path)
				}
			}
		}
	}
	fmt.Println("include files are all system headers")
}

func TestMacroExpansionOtherFile() {
	c.Printf(c.Str("TestMacroExpansionOtherFile:\n"))
	test.RunTestWithConfig(&llcppg.Config{
		Include: []string{"ref.h"},
		CFlags:  "-I./testdata/macroexpan",
	})
}
