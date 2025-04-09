package main

import (
	"fmt"
	"strings"

	"github.com/goplus/lib/c"
	"github.com/goplus/llcppg/_xtool/llcppsigfetch/parse"
	test "github.com/goplus/llcppg/_xtool/llcppsigfetch/parse/cvt_test"
	"github.com/goplus/llcppg/llcppg"
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
	cvt, err := parse.Do(&parse.ParseConfig{
		Conf: &llcppg.Config{
			Include: []string{"sys.h"},
			CFlags:  "-I./testdata/sysinc",
		},
	})
	if err != nil {
		panic(err)
	}
	found := false
	for path, info := range cvt.Pkg.FileMap {
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

func TestSystemHeader() {
	fmt.Println("=== TestSystemHeader ===")
	cvt, err := parse.Do(&parse.ParseConfig{
		Conf: &llcppg.Config{
			Include: []string{"inc.h"},
			CFlags:  "-I./testdata/sysinc",
		},
	})
	if err != nil {
		panic(err)
	}

	for path, info := range cvt.Pkg.FileMap {
		if path != "./testdata/sysinc/inc.h" && info.FileType != llcppg.Third {
			panic(fmt.Errorf("include file is not third header: %s", path))
		}
	}

	fmt.Println("include files are all system headers")
}

func TestMacroExpansionOtherFile() {
	c.Printf(c.Str("=== TestMacroExpansionOtherFile ===\n"))
	test.RunTestWithConfig(&parse.ParseConfig{
		Conf: &llcppg.Config{
			Include: []string{"ref.h"},
			CFlags:  "-I./testdata/macroexpan",
		},
	})
}
