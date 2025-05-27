package gen

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	llcppg "github.com/goplus/llcppg/config"
)

func Test_emptyStringError_Error(t *testing.T) {
	tests := []struct {
		name string
		p    *emptyStringError
		want string
	}{
		{
			"newEmptyStringError",
			newEmptyStringError("param"),
			"param can't be empty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.p.Error(); got != tt.want {
				t.Errorf("emptyStringError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newEmptyStringError(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want *emptyStringError
	}{
		{
			"newEmptyStringError",
			args{
				"param",
			},
			&emptyStringError{"param"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.args.name, func(t *testing.T) {
			if got := newEmptyStringError(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newEmptyStringError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isExcludeDir(t *testing.T) {
	type args struct {
		relPath        string
		excludeSubdirs []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"isExcludeDir",
			args{
				filepath.Join("a", "b", "c", "d"),
				[]string{"a"},
			},
			true,
		},
		{
			"isExcludeDir",
			args{
				filepath.Join("a", "b", "c", "d"),
				[]string{"b"},
			},
			false,
		},
		{
			"isExcludeDir",
			args{
				"ab",
				[]string{"a"},
			},
			false,
		},
		{
			"isExcludeDir",
			args{
				"a",
				[]string{"a"},
			},
			true,
		},
		{
			"isExcludeDir",
			args{
				filepath.Join("a", ""),
				[]string{"a"},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.args.relPath, func(t *testing.T) {
			if got := isExcludeDir(tt.args.relPath, tt.args.excludeSubdirs); got != tt.want {
				t.Errorf("isExcludeDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExpandName(t *testing.T) {
	type args struct {
		name   string
		dir    string
		cfgKey llcppCfgKey
	}
	tests := []struct {
		name             string
		args             args
		wantExpandPrefix string
	}{
		{
			"cflags",
			args{
				"libcjson",
				"",
				cfgCflagsKey,
			},
			"-I/",
		},
		{
			"libs",
			args{
				"libcjson",
				"",
				cfgLibsKey,
			},
			"-",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExpand := ExpandName(tt.args.name, tt.args.dir, tt.args.cfgKey)
			if !strings.HasPrefix(gotExpand, tt.wantExpandPrefix) {
				t.Errorf("ExpandName() gotExpand = %v, want %v", gotExpand, tt.wantExpandPrefix)
			}
		})
	}
}

func Test_findDepSlice(t *testing.T) {
	type args struct {
		lines []string
	}
	tests := []struct {
		name  string
		args  args
		want  []string
		want1 string
	}{
		{
			"cjson",
			args{
				[]string{"cJSON_Utils.o:", "cjson/cJSON_Utils.h", "cjson/cJSON.h"},
			},
			[]string{"cjson/cJSON.h"},
			"cJSON_Utils.o:cjson/cJSON_Utils.h",
		},
		{
			"cjson",
			args{
				[]string{"cJSON_Utils.o:cjson/cJSON_Utils.h", "cjson/cJSON.h"},
			},
			[]string{"cjson/cJSON.h"},
			"cJSON_Utils.o:cjson/cJSON_Utils.h",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := findDepSlice(tt.args.lines)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findDepSlice() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("findDepSlice() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

type parseEntry struct {
	isDir bool
	name  string
	typ   fs.FileMode
}

func (p *parseEntry) Name() string {
	return p.name
}

func (p *parseEntry) IsDir() bool {
	return p.isDir
}

func (p *parseEntry) Type() fs.FileMode {
	return p.typ
}

func (p *parseEntry) Info() (os.FileInfo, error) {
	return nil, nil
}

func Test_parseFileEntry(t *testing.T) {
	cflags, inc := newCflags("cfg_test_data/cjson/include/")
	path := inc + "cJSON.h"
	internalPath := filepath.Join(inc, "internal")
	objFile := NewObjFileString("cJSON.o:cJSON.h")
	objFile.Deps = []string{}
	type args struct {
		cflags         string
		trimStr        string
		path           string
		d              fs.DirEntry
		exts           []string
		excludeSubdirs []string
	}
	tests := []struct {
		name string
		args args
		want *ObjFile
	}{
		{
			"cjson_ok",
			args{
				cflags,
				inc,
				path,
				&parseEntry{name: "cJSON.h"},
				[]string{".h"},
				[]string{},
			},
			objFile,
		},
		{
			"cjson_exts_not_found",
			args{
				cflags,
				inc,
				path,
				&parseEntry{name: "cJSON.h"},
				[]string{".hh"},
				[]string{},
			},
			nil,
		},
		{
			"cjson_not_rel_path",
			args{
				cflags,
				inc,
				"cJSON.h",
				&parseEntry{name: "cJSON.h"},
				[]string{".h"},
				[]string{},
			},
			objFile,
		},
		{
			"cjson_exclude_dir",
			args{
				cflags,
				inc,
				filepath.Join(internalPath, "a.h"),
				&parseEntry{name: "a.h"},
				[]string{".h"},
				[]string{"internal"},
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := parseFileEntry(tt.args.cflags, tt.args.trimStr, tt.args.path, tt.args.d, tt.args.exts, tt.args.excludeSubdirs)
			if tt.want != nil && got != nil && !got.IsEqual(tt.want) {
				t.Errorf("parseFileEntry() = %v, want %v", got, tt.want)
			}
			if tt.want == nil {
				if got != nil {
					t.Errorf("parseFileEntry() = %v, want nil", got)
				}
			}
		})
	}
}

func Test_parseCFlagsEntry(t *testing.T) {
	cflags, inc := newCflags("cfg_test_data/cjson/include/")
	type args struct {
		cflags         string
		cflag          string
		exts           []string
		excludeSubdirs []string
	}
	tests := []struct {
		name string
		args args
		want *CflagEntry
	}{
		{
			"cjson_prefix_with_I",
			args{
				cflags,
				cflags,
				[]string{".h"},
				[]string{"internal"},
			},
			&CflagEntry{
				Include: inc,
				ObjFiles: []*ObjFile{
					{OFile: "cJSON_Utils.o", HFile: "cJSON_Utils.h", Deps: []string{"cJSON.h"}},
					{OFile: "cJSON.o", HFile: "cJSON.h", Deps: []string{}},
				},
				InvalidObjFiles: []*ObjFile{},
			},
		},
		{
			"cjson_not_prefix_with_I",
			args{
				cflags,
				inc,
				[]string{".h"},
				[]string{},
			},
			nil,
		},
		{
			"cjson_walk_dir_err",
			args{
				cflags,
				"-I/a/b/c",
				[]string{".h"},
				[]string{},
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseCFlagsEntry(tt.args.cflags, tt.args.cflag, tt.args.exts, tt.args.excludeSubdirs)
			if got != nil && tt.want != nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseCFlagsEntry() = %v, want %v", got, tt.want)
			}
			if tt.want == nil {
				if got != nil {
					t.Errorf("parseCFlagsEntry() = %v, want nil", got)
				}
			}
		})
	}
}

func Test_sortIncludes(t *testing.T) {
	cflags, _ := newCflags("cfg_test_data/cjson/include/")
	depsCflags, _ := newCflags("cfg_test_data/deps/")
	cfg := &llcppg.Config{
		Name:   "libcjson",
		CFlags: cflags,
		Libs:   "$(pkg-config --libs libcjson)",
	}
	depCfg := &llcppg.Config{
		Name:   "deps",
		CFlags: depsCflags,
		Libs:   "",
	}
	type args struct {
		expandCflags   string
		cfg            *llcppg.Config
		exts           []string
		excludeSubdirs []string
	}
	tests := []struct {
		name        string
		args        args
		wantInclude []string
	}{
		{
			"cjson",
			args{
				cflags,
				cfg,
				[]string{".h"},
				[]string{"internal"},
			},
			[]string{
				"cJSON_Utils.h",
				"cJSON.h",
			},
		},
		{
			"deps/case1",
			args{
				filepath.Join(depsCflags, "case1"),
				depCfg,
				[]string{".h"},
				[]string{},
			},
			[]string{
				"a.h",
				"b.h",
			},
		},
		{
			"deps/case2",
			args{
				filepath.Join(depsCflags, "case2"),
				depCfg,
				[]string{".h"},
				[]string{},
			},
			[]string{
				"a.h",
				"b.h",
				"c.h",
			},
		},
		{
			"deps/case3_recircle",
			args{
				filepath.Join(depsCflags, "case3"),
				depCfg,
				[]string{".h"},
				[]string{},
			},
			[]string{},
		},
		{
			"deps/case4_recircle",
			args{
				filepath.Join(depsCflags, "case4"),
				depCfg,
				[]string{".h"},
				[]string{},
			},
			[]string{
				"a.h",
				"b.h",
				"f.h",
				"c.h",
				"e.h",
				"g.h",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortIncludes(tt.args.expandCflags, tt.args.cfg, tt.args.exts, tt.args.excludeSubdirs)
			if !reflect.DeepEqual(tt.args.cfg.Include, tt.wantInclude) {
				t.Errorf("sortIncludes() = %v, want %v", tt.args.cfg.Include, tt.wantInclude)
			}
		})
	}
}

func TestNewLLCppConfig(t *testing.T) {
	type args struct {
		name string
		flag FlagMode
	}
	tests := []struct {
		name string
		args args
		want *llcppg.Config
	}{
		{
			"libcjson",
			args{
				"libcjson",
				WithTab,
			},
			&llcppg.Config{
				Name:           "libcjson",
				CFlags:         "$(pkg-config --cflags libcjson)",
				Libs:           "$(pkg-config --libs libcjson)",
				Cplusplus:      false,
				KeepUnderScore: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newLLCppgConfig(tt.args.name, tt.args.flag); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLLCppgConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNormalizePackageName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"lua5.4",
			args{
				"lua5.4",
			},
			"lua5_4",
		},
		{
			"tree-sitter",
			args{
				"tree-sitter",
			},
			"tree_sitter",
		},
		{
			"python-3.12-embed",
			args{
				"python-3.12-embed",
			},
			"python_3_12_embed",
		},
		{
			"libmpdec++",
			args{
				"libmpdec++",
			},
			"libmpdec",
		},
		{
			"startWithDigit",
			args{
				"3json",
			},
			"_3json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizePackageName(tt.args.name); got != tt.want {
				t.Errorf("NormalizePackageName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func newCflags(relDir string) (cflags string, path string) {
	dir, file := filepath.Split(relDir)
	wd, _ := os.Getwd()
	inc := joinPath(wd, dir)
	if len(file) > 0 {
		inc += file
	}
	trimStr := "-I" + inc
	return trimStr, inc
}

func joinPath(dir, rel string) string {
	path := filepath.Join(dir, rel)
	if !strings.HasSuffix(path, string(filepath.Separator)) {
		path += string(filepath.Separator)
	}
	return path
}

func Test_getDir(t *testing.T) {
	type args struct {
		relPath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"getDir/",
			args{
				relPath: "a/",
			},
			"a",
		},
		{
			"getDir",
			args{
				relPath: "a",
			},
			"a",
		},
		{
			"getDir/a/b/c/d",
			args{
				relPath: "a/b/c/d",
			},
			"a",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDir(tt.args.relPath); got != tt.want {
				t.Errorf("getDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getClangArgs(t *testing.T) {
	cflags, inc := newCflags("cfg_test_data/same_rel")
	multiCflags := fmt.Sprintf("%s %s", filepath.Join(cflags, "libcjson/include"), filepath.Join(cflags, "stdcjson/include"))
	type args struct {
		cflags  string
		relpath string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"cjson",
			args{
				cflags:  multiCflags,
				relpath: "cjson.h",
			},
			[]string{"-I/libcjson/include", "-I/stdcjson/include", "-MM", "cjson.h"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getClangArgs(tt.args.cflags, tt.args.relpath)
			if tt.name == "cjson" {
				got[0] = strings.ReplaceAll(got[0], inc, "")
				got[1] = strings.ReplaceAll(got[1], inc, "")
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getClangArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}
