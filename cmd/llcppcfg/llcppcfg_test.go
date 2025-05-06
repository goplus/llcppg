package main

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func recoverFn(fn func()) (ret any) {
	defer func() {
		ret = recover()
	}()
	fn()
	return
}

func readFile(filepath string) []byte {
	buf, _ := os.ReadFile(filepath)
	return buf
}

func TestLLCppcfg(t *testing.T) {

	llcppgFileName := filepath.Join("macos", "llcppg.cfg")
	if runtime.GOOS == "linux" {
		// cuurently, due to llcppcfg recognizing system path fail, all includes are empty for temporary tests.
		// TODO(ghl): fix it
		llcppgFileName = filepath.Join("linux", "llcppg.cfg")
	}

	cjsonCfgFilePath := filepath.Join("llcppgcfg", "cfg_test_data", "cjson", "conf", llcppgFileName)
	bdwgcCfgFilePath := filepath.Join("llcppgcfg", "cfg_test_data", "bdw-gc", "conf", llcppgFileName)
	libffiCfgFilePath := filepath.Join("llcppgcfg", "cfg_test_data", "libffi", "conf", llcppgFileName)
	libxsltCfgFilePath := filepath.Join("llcppgcfg", "cfg_test_data", "libxslt", "conf", llcppgFileName)

	type args struct {
		name           string
		tab            string
		exts           []string
		deps           []string
		excludeSubdirs []string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"libcjson",
			args{
				"libcjson",
				"true",
				[]string{".h"},
				[]string{},
				[]string{},
			},
			readFile(cjsonCfgFilePath),
			false,
		},
		{
			"bdw-gc",
			args{
				"bdw-gc",
				"true",
				[]string{".h"},
				[]string{},
				[]string{},
			},
			readFile(bdwgcCfgFilePath),
			false,
		},
		{
			"libxslt",
			args{
				"libxslt",
				"true",
				[]string{".h"},
				[]string{"c/os", "github.com/goplus/llpkg/libxml2@v1.0.0"},
				[]string{},
			},
			readFile(libxsltCfgFilePath),
			false,
		},
		{
			"libffi",
			args{
				"libffi",
				"true",
				[]string{".h"},
				[]string{},
				[]string{},
			},
			readFile(libffiCfgFilePath),
			false,
		},
		{
			"empty_name",
			args{
				"",
				"true",
				[]string{".h"},
				[]string{},
				[]string{},
			},
			nil,
			true,
		},
		{
			"normal_not_sort",
			args{
				"libcjson",
				"false",
				[]string{".h"},
				[]string{},
				[]string{},
			},
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = []string{
				"llcppcfg",
			}
			if len(tt.args.deps) > 0 {
				os.Args = append(os.Args, "-deps", strings.Join(tt.args.deps, " "))
			}
			if len(tt.args.excludeSubdirs) > 0 {
				os.Args = append(os.Args, "-excludes", strings.Join(tt.args.excludeSubdirs, " "))
			}
			if len(tt.args.exts) > 0 {
				os.Args = append(os.Args, "-exts", strings.Join(tt.args.exts, " "))
			}
			os.Args = append(os.Args, tt.args.name)

			// reset flags for the next test
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			ret := recoverFn(main)
			if ret != nil {
				if !tt.wantErr {
					t.Errorf("%v", ret)
				}
				return
			}
			defer os.Remove("llcppg.cfg")
			if tt.want == nil {
				return
			}
			b, err := os.ReadFile("llcppg.cfg")
			if err != nil {
				t.Error(err)
				return
			}
			if !bytes.Equal(b, tt.want) {
				t.Errorf("unexpected content: want %s got %s", string(tt.want), string(b))
			}
		})
	}
}
