package main

import (
	"bytes"
	"flag"
	"fmt"
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

func readFile(filepath string) *bytes.Buffer {
	buf, err := os.ReadFile(filepath)
	if err != nil {
		return bytes.NewBufferString("")
	}
	return bytes.NewBuffer(buf)
}

func TestLLCppcfg(t *testing.T) {
	llcppgFileName := "llcppg.cfg"
	if runtime.GOOS == "linux" {
		llcppgFileName = "llcppg_linux.cfg"
	}
	cjsonCfgFilePath := filepath.Join("llcppgcfg", "cfg_test_data", "cjson", llcppgFileName)
	bdwgcCfgFilePath := filepath.Join("llcppgcfg", "cfg_test_data", "bdw-gc", llcppgFileName)
	libffiCfgFilePath := filepath.Join("llcppgcfg", "cfg_test_data", "libffi", llcppgFileName)
	libxsltCfgFilePath := filepath.Join("llcppgcfg", "cfg_test_data", "libxslt", llcppgFileName)

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
		want    *bytes.Buffer
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
				"-exts", fmt.Sprintf(`"%s"`, strings.Join(tt.args.exts, " ")),
				"-deps", fmt.Sprintf(`"%s"`, strings.Join(tt.args.deps, " ")),
				"-excludes", fmt.Sprintf(`"%s"`, strings.Join(tt.args.excludeSubdirs, " ")),
				"-tab", tt.args.tab,
				"-cpp", "false",
				tt.args.name,
			}
			// reset flag
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			ret := recoverFn(main)

			if ret != nil {
				t.Errorf("%v", ret)
				return
			}
			defer os.Remove("llcppg.cfg")
			b, err := os.ReadFile("llcppg.cfg")
			if err != nil {
				t.Error(err)
				return
			}
			if !bytes.Equal(b, tt.want.Bytes()) {
				t.Errorf("unexpected content: want %s got %s", tt.want.String(), string(b))
			}
		})
	}
}
