package cvttest

import (
	"os"
	"sort"
	"unsafe"

	"github.com/goplus/lib/c"
	clangutils "github.com/goplus/llcppg/_xtool/internal/clang"
	"github.com/goplus/llcppg/_xtool/internal/parser"
	"github.com/goplus/llcppg/_xtool/llcppsigfetch/parse"
	llcppg "github.com/goplus/llcppg/config"
	"github.com/goplus/llpkg/cjson"
)

func RunTest(testName string, testCases []string) {
	for i, content := range testCases {
		tempIncFile := clangutils.TEMP_FILE
		if err := os.WriteFile(tempIncFile, []byte(content), 0644); err != nil {
			panic(err)
		}
		include := []string{tempIncFile}
		cflags := "-I./"
		c.Printf(c.Str("%s Case %d:\n"), c.AllocaCStr(testName), c.Int(i+1))
		RunTestWithConfig(&parse.Config{
			Conf: &llcppg.Config{
				Cplusplus: true,
				Include:   include,
				CFlags:    cflags,
			},
		})
		os.Remove(tempIncFile)
	}
}

func RunTestWithConfig(config *parse.Config) {
	config.Exec = func(conf *parse.Config, pkg *llcppg.Pkg) {
		result := MarshalPkg(pkg)
		str := result.Print()
		c.Printf(c.Str("%s\n\n"), str)
		cjson.FreeCStr(unsafe.Pointer(str))
		result.Delete()
	}
	err := parse.Do(config)
	if err != nil {
		panic(err)
	}
}

// for test order map
func MarshalPkg(pkg *llcppg.Pkg) *cjson.JSON {
	root := cjson.Object()
	root.SetItem(c.Str("File"), parser.MarshalASTFile(pkg.File))
	root.SetItem(c.Str("FileMap"), MarshalFileMap(pkg.FileMap))
	return root
}

// for test order map
func MarshalFileMap(fmap map[string]*llcppg.FileInfo) *cjson.JSON {
	root := cjson.Object()
	keys := make([]string, 0, len(fmap))
	for path := range fmap {
		keys = append(keys, path)
	}
	sort.Strings(keys)
	for _, path := range keys {
		root.SetItem(c.AllocaCStr(path), parse.MarshalFileInfo(fmap[path]))
	}
	return root
}
