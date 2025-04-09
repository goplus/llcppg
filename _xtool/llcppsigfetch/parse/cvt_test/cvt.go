package cvttest

import (
	"fmt"
	"os"
	"sort"
	"unsafe"

	"github.com/goplus/lib/c"
	"github.com/goplus/lib/c/clang"
	"github.com/goplus/llcppg/_xtool/llcppsigfetch/parse"
	"github.com/goplus/llcppg/_xtool/llcppsymg/clangutils"
	"github.com/goplus/llcppg/llcppg"
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
		RunTestWithConfig(&parse.ParseConfig{
			Conf: &llcppg.Config{
				Cplusplus: true,
				Include:   include,
				CFlags:    cflags,
			},
		})
		os.Remove(tempIncFile)
	}
}

func RunTestWithConfig(config *parse.ParseConfig) {
	cvt, err := parse.Do(config)
	if err != nil {
		panic(err)
	}
	result := MarshalPkg(cvt.Pkg)
	str := result.Print()
	c.Printf(c.Str("%s\n\n"), str)
	cjson.FreeCStr(unsafe.Pointer(str))
	result.Delete()
}

// for test order map
func MarshalPkg(pkg *llcppg.Pkg) *cjson.JSON {
	root := cjson.Object()
	root.SetItem(c.Str("File"), parse.MarshalASTFile(pkg.File))
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

type GetTypeOptions struct {
	TypeCode string // e.g. "char*", "char**"

	// ExpectTypeKind specifies the expected type kind (optional)
	// Use clang.Type_Invalid to accept any type (default behavior)
	// *For complex types (when <complex.h> is included), specifying this is crucial
	// to filter out the correct type, as there will be multiple VarDecl fields present
	ExpectTypeKind clang.TypeKind

	// Args contains additional compilation arguments passed to Clang (optional)
	// These are appended after the default language-specific arguments
	// Example: []string{"-std=c++11"}
	Args []string

	// IsCpp indicates whether the code should be treated as C++ (true) or C (false)
	// This affects the default language arguments passed to Clang:
	// - For C++: []string{"-x", "c++"}
	// - For C:   []string{"-x", "c"}
	// *For complex C types, C Must be specified
	IsCpp bool
}

// GetType returns the clang.Type of the given type code
// Need to dispose the index and unit after using
// e.g. index.Dispose(), unit.Dispose()
func GetType(option *GetTypeOptions) (clang.Type, *clang.Index, *clang.TranslationUnit) {
	code := fmt.Sprintf("%s placeholder;", option.TypeCode)
	index, unit, err := clangutils.CreateTranslationUnit(&clangutils.Config{
		File:  code,
		Temp:  true,
		Args:  option.Args,
		IsCpp: option.IsCpp,
	})
	if err != nil {
		panic(err)
	}
	cursor := unit.Cursor()
	var typ clang.Type
	clangutils.VisitChildren(cursor, func(child, parent clang.Cursor) clang.ChildVisitResult {
		if child.Kind == clang.CursorVarDecl && (option.ExpectTypeKind == clang.TypeInvalid || option.ExpectTypeKind == child.Type().Kind) {
			typ = child.Type()
			return clang.ChildVisit_Break
		}
		return clang.ChildVisit_Continue
	})
	return typ, index, unit
}
