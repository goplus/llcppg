package main

import (
	"fmt"
	"strings"

	"github.com/goplus/lib/c"
	"github.com/goplus/llcppg/_xtool/llcppsigfetch/parse"
	test "github.com/goplus/llcppg/_xtool/llcppsigfetch/parse/cvt_test"
	llcppg "github.com/goplus/llcppg/config"
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
	checkFileMap := func(conf *parse.Config, pkg *llcppg.Pkg) {
		found := false
		for path, info := range pkg.FileMap {
			if strings.HasSuffix(path, "sys/types.h") && info.FileType == llcppg.Third {
				found = true
			}
		}
		if !found {
			panic("sys/types.h not found")
		} else {
			fmt.Println("sys/types.h include path found")
		}
	}
	err := parse.Do(&parse.Config{
		Conf: &llcppg.Config{
			Include: []string{"sys.h"},
			CFlags:  "-I./testdata/sysinc",
		},
		Exec: checkFileMap,
	})
	if err != nil {
		panic(err)
	}

}

func TestSystemHeader() {
	fmt.Println("=== TestSystemHeader ===")
	checkFileMap := func(conf *parse.Config, pkg *llcppg.Pkg) {
		for path, info := range pkg.FileMap {
			if path != "testdata/sysinc/inc.h" && info.FileType != llcppg.Third {
				panic(fmt.Errorf("include file is not third header: %s", path))
			}
		}
		fmt.Println("include files are all system headers")
	}
	err := parse.Do(&parse.Config{
		Conf: &llcppg.Config{
			Include: []string{"inc.h"},
			CFlags:  "-I./testdata/sysinc",
		},
		Exec: checkFileMap,
	})
	if err != nil {
		panic(err)
	}
}

func TestMacroExpansionOtherFile() {
	c.Printf(c.Str("=== TestMacroExpansionOtherFile ===\n"))
	test.RunTestWithConfig(&parse.Config{
		Conf: &llcppg.Config{
			Include: []string{"ref.h"},
			CFlags:  "-I./testdata/macroexpan",
		},
	})
}
