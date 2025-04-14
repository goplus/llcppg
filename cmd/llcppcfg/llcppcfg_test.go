package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func readFile(filepath string) *bytes.Buffer {
	buf, err := os.ReadFile(filepath)
	if err != nil {
		return bytes.NewBufferString("")
	}
	return bytes.NewBuffer(buf)
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
			args := []string{"run", "."}
			if len(tt.args.deps) > 0 {
				args = append(args, "-deps", strings.Join(tt.args.deps, " "))
			}
			if len(tt.args.excludeSubdirs) > 0 {
				args = append(args, "-excludes", strings.Join(tt.args.excludeSubdirs, " "))
			}
			if len(tt.args.exts) > 0 {
				args = append(args, "-exts", strings.Join(tt.args.exts, " "))
			}
			args = append(args, tt.args.name)

			cmd := exec.Command("go", args...)
			ret, err := cmd.CombinedOutput()
			if err != nil {
				if !tt.wantErr {
					t.Error(string(ret))
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
			if !bytes.Equal(b, tt.want.Bytes()) {
				t.Errorf("unexpected content: want %s got %s", tt.want.String(), string(b))
			}
		})
	}
}
