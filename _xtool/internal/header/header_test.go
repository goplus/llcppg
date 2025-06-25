package header_test

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/goplus/lib/c/clang"
	clangutils "github.com/goplus/llcppg/_xtool/internal/clang"
	"github.com/goplus/llcppg/_xtool/internal/clangtool"
	"github.com/goplus/llcppg/_xtool/internal/header"
	llconfig "github.com/goplus/llcppg/config"
)

func TestPkgHfileInfo(t *testing.T) {
	cases := []struct {
		conf *llconfig.Config
		want *header.PkgHfilesInfo
	}{
		{
			conf: &llconfig.Config{
				CFlags:  "-I./testdata/hfile -I ./testdata/thirdhfile",
				Include: []string{"temp1.h", "temp2.h"},
			},
			want: &header.PkgHfilesInfo{
				Inters: []string{"testdata/hfile/temp1.h", "testdata/hfile/temp2.h"},
				Impls:  []string{"testdata/hfile/tempimpl.h"},
			},
		},
		{
			conf: &llconfig.Config{
				CFlags:  "-I./testdata/hfile -I ./testdata/thirdhfile",
				Include: []string{"temp1.h", "temp2.h"},
				Mix:     true,
			},
			want: &header.PkgHfilesInfo{
				Inters: []string{"testdata/hfile/temp1.h", "testdata/hfile/temp2.h"},
				Impls:  []string{},
			},
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			info := header.PkgHfileInfo(tc.conf.Include, strings.Fields(tc.conf.CFlags), tc.conf.Mix)
			if !reflect.DeepEqual(info.Inters, tc.want.Inters) {
				t.Fatalf("inter expected %v, but got %v", tc.want.Inters, info.Inters)
			}
			if !reflect.DeepEqual(info.Impls, tc.want.Impls) {
				t.Fatalf("impl expected %v, but got %v", tc.want.Impls, info.Impls)
			}

			thirdhfile, err := filepath.Abs("./testdata/thirdhfile/third.h")
			if err != nil {
				t.Fatalf("failed to get abs path: %w", err)
			}
			tfileFound := false
			stdioFound := false
			for _, tfile := range info.Thirds {
				absTfile, err := filepath.Abs(tfile)
				if err != nil {
					t.Fatalf("failed to get abs path: %w", err)
				}
				if absTfile == thirdhfile {
					tfileFound = true
				}
				if strings.HasSuffix(absTfile, "stdio.h") {
					stdioFound = true
				}
			}
			if !tfileFound || !stdioFound {
				t.Fatalf("third hfile or std hfile not found")
			}
		})
	}
}

func TestLongestPrefix(t *testing.T) {
	testCases := []struct {
		name string
		strs []string
		want string
	}{
		{
			name: "empty string 1",
			strs: []string{},
			want: "",
		},
		{
			name: "empty string 2",
			strs: []string{"", ""},
			want: ".",
		},
		{
			name: "one empty string(b)",
			strs: []string{"/a", ""},
			want: "",
		},
		{
			name: "one empty string(a)",
			strs: []string{"", "/a"},

			want: "",
		},
		// FIXME: substring bug
		// {
		// 	name: "b is substring of a",
		// 	strs: []string{"/usr/a/b", "/usr/a"},
		// 	want: "/usr/a",
		// },
		// {
		// 	name: "a is substring of b",
		// 	strs: []string{"/usr/c", "/usr/c/b"},
		// 	want: "/usr/c",
		// },
		{
			name: "normal case 1",
			strs: []string{"testdata/hfile/temp1.h", "testdata/thirdhfile/third.h"},
			want: "testdata",
		},
		{
			name: "normal case 2",
			strs: []string{"testdata/hfile/temp1.h", "testdata/hfile/third.h"},

			want: "testdata/hfile",
		},
		// FIXME: absolute path
		// {
		// 	name: "normal case 3",
		// 	strs: []string{"/opt/homebrew/Cellar/cjson/1.7.18/include/cJSON/cJSON.h", "/opt/homebrew/Cellar/cjson/1.7.18/include/cJSON.h", "/opt/homebrew/Cellar/cjson/1.7.18/include/zlib/zlib.h"},
		// 	want: "/opt/homebrew/Cellar/cjson/1.7.18/include",
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := header.CommonParentDir(tc.strs); got != tc.want {
				t.Fatalf("unexpected longest prefix: want %s got %s", tc.want, got)
			}
		})
	}
}

func benchmarkFn(fn func()) time.Duration {
	now := time.Now()

	fn()

	return time.Since(now)
}

func TestBenchmarkPkgHfileInfo(t *testing.T) {
	include := []string{"temp1.h", "temp2.h"}
	cflags := []string{"-I./testdata/hfile", "-I./testdata/thirdhfile"}
	t1 := benchmarkFn(func() {
		for i := 0; i < 100; i++ {
			pkgHfileInfo(include, cflags, false)
		}
	})

	t2 := benchmarkFn(func() {
		for i := 0; i < 100; i++ {
			header.PkgHfileInfo(include, cflags, false)
		}
	})

	fmt.Println("old PkgHfileInfo elapsed: ", t1, "new PkgHfileInfo elasped: ", t2)
}

func pkgHfileInfo(includes []string, args []string, mix bool) *header.PkgHfilesInfo {
	info := &header.PkgHfilesInfo{
		Inters: []string{},
		Impls:  []string{},
		Thirds: []string{},
	}
	outfile, err := os.CreateTemp("", "compose_*.h")
	if err != nil {
		panic(err)
	}
	defer os.Remove(outfile.Name())

	inters := make(map[string]struct{})
	others := []string{} // impl & third
	for _, f := range includes {
		content := "#include <" + f + ">"
		index, unit, err := clangutils.CreateTranslationUnit(&clangutils.Config{
			File: content,
			Temp: true,
			Args: args,
		})
		if err != nil {
			panic(err)
		}
		clangutils.GetInclusions(unit, func(inced clang.File, incins []clang.SourceLocation) {
			if len(incins) == 1 {
				filename := filepath.Clean(clang.GoString(inced.FileName()))
				info.Inters = append(info.Inters, filename)
				inters[filename] = struct{}{}
			}
		})
		unit.Dispose()
		index.Dispose()
	}

	clangtool.ComposeIncludes(includes, outfile.Name())
	index, unit, err := clangutils.CreateTranslationUnit(&clangutils.Config{
		File: outfile.Name(),
		Temp: false,
		Args: args,
	})
	defer unit.Dispose()
	defer index.Dispose()
	if err != nil {
		panic(err)
	}
	clangutils.GetInclusions(unit, func(inced clang.File, incins []clang.SourceLocation) {
		// not in the first level include maybe impl or third hfile
		filename := filepath.Clean(clang.GoString(inced.FileName()))
		_, inter := inters[filename]
		if len(incins) > 1 && !inter {
			others = append(others, filename)
		}
	})

	if mix {
		info.Thirds = others
		return info
	}

	root, err := filepath.Abs(header.CommonParentDir(info.Inters))
	if err != nil {
		panic(err)
	}
	for _, f := range others {
		file, err := filepath.Abs(f)
		if err != nil {
			panic(err)
		}
		if strings.HasPrefix(file, root) {
			info.Impls = append(info.Impls, f)
		} else {
			info.Thirds = append(info.Thirds, f)
		}
	}
	return info
}
