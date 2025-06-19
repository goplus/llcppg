package parse_test

import (
	"encoding/json"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/goplus/llcppg/_xtool/internal/clangtool"
	"github.com/goplus/llcppg/_xtool/internal/parser"
	"github.com/goplus/llcppg/_xtool/llcppsigfetch/internal/parse"
	llcppg "github.com/goplus/llcppg/config"
)

func TestInclusionMap(t *testing.T) {
	checkFileMap := func(conf *parse.Config, pkg *llcppg.Pkg) {
		found := false
		for path, info := range pkg.FileMap {
			if strings.HasSuffix(path, "sys/types.h") && info.FileType == llcppg.Third {
				found = true
			}
		}
		if !found {
			t.Fatalf("sys/types.h not found")
		}
	}

	combinedFile := createCombineFile([]string{"sys.h"})
	defer os.Remove(combinedFile)
	err := parse.Do(&parse.Config{
		CombinedFile: combinedFile,
		Conf: &llcppg.Config{
			Include: []string{"sys.h"},
			CFlags:  "-I./testdata/sysinc/hfile",
		},
		Exec: checkFileMap,
	})
	if err != nil {
		panic(err)
	}

}

func TestSystemHeader(t *testing.T) {
	checkFileMap := func(conf *parse.Config, pkg *llcppg.Pkg) {
		for path, info := range pkg.FileMap {
			if path != "testdata/sysinc/hfile/inc.h" && info.FileType != llcppg.Third {
				t.Fatalf("include file is not third header: %s", path)
			}
		}
	}
	combinedFile := createCombineFile([]string{"inc.h"})
	defer os.Remove(combinedFile)
	err := parse.Do(&parse.Config{
		CombinedFile: combinedFile,
		Conf: &llcppg.Config{
			Include: []string{"inc.h"},
			CFlags:  "-I./testdata/sysinc/hfile",
		},
		Exec: checkFileMap,
	})
	if err != nil {
		panic(err)
	}
}

func TestMacroExpansionOtherFile(t *testing.T) {
	combinedFile := createCombineFile([]string{"ref.h"})
	defer os.Remove(combinedFile)
	conf := &parse.Config{
		CombinedFile: combinedFile,
		Conf: &llcppg.Config{
			Include: []string{"ref.h"},
			CFlags:  "-I./testdata/macroexpan/hfile",
		},
	}
	testFrom(t, conf, "testdata/macroexpan", false)
}

func testFrom(t *testing.T, conf *parse.Config, dir string, gen bool) {
	pkg := parseWithConfig(conf)

	result := marshalPkg(pkg)
	output, _ := json.MarshalIndent(&result, "", "  ")

	expectFile := filepath.Join(dir, "expect.json")
	if gen {
		err := os.WriteFile(expectFile, output, os.ModePerm)
		if err != nil {
			t.Fatal("WriteFile failed:", err)
		}
	} else {
		json, err := os.ReadFile(expectFile)
		if err != nil {
			t.Fatal("ReadExpectFile failed:", err)
		}
		expect := string(json)
		if expect != string(output) {
			t.Fatalf("expect %s, got %s", expect, string(output))
		}
	}
}

func parseWithConfig(config *parse.Config) (res *llcppg.Pkg) {
	config.Exec = func(conf *parse.Config, pkg *llcppg.Pkg) {
		res = pkg
	}
	err := parse.Do(config)
	if err != nil {
		panic(err)
	}
	return
}

// for test order map
func marshalPkg(pkg *llcppg.Pkg) map[string]any {
	return map[string]any{
		"File":    parser.XMarshalASTFile(pkg.File),
		"FileMap": marshalFileMap(pkg.FileMap),
	}
}

// for test order map
func marshalFileMap(fmap map[string]*llcppg.FileInfo) map[string]any {
	root := make(map[string]any)
	keys := slices.Collect(maps.Keys(fmap))

	sort.Strings(keys)

	for _, path := range keys {
		root[path] = parse.MarshalFileInfo(fmap[path])
	}
	return root
}

func createCombineFile(includes []string) string {
	combinedFile, err := os.CreateTemp("./", "compose_*.h")
	if err != nil {
		panic(err)
	}
	err = clangtool.ComposeIncludes(includes, combinedFile.Name())
	if err != nil {
		panic(err)
	}
	return combinedFile.Name()
}
