package _cmptest

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/goplus/llpkgstore/upstream"
	"github.com/goplus/llpkgstore/upstream/installer/conan"
)

type cfgTestCase struct {
	cfgDir string
	pkg    upstream.Package
	config map[string]string // conan options
}

var cfgTestCases = []cfgTestCase{
	{
		cfgDir: "./cfgdata/cjson",
		pkg:    upstream.Package{Name: "cjson", Version: "1.7.18"},
		config: map[string]string{
			"options": "utils=True",
		},
	},
	{
		cfgDir: "./cfgdata/cjson",
		pkg:    upstream.Package{Name: "cjson", Version: "1.7.17"},
		config: map[string]string{
			"options": "utils=True",
		},
	},
	{
		cfgDir: "./cfgdata/libxml2",
		pkg:    upstream.Package{Name: "libxml2", Version: "2.13.6"},
		config: map[string]string{
			"options": "iconv=False",
		},
	},
	{
		cfgDir: "./cfgdata/sqlite3",
		pkg:    upstream.Package{Name: "sqlite3", Version: "3.49.1"},
	},
	{
		cfgDir: "./cfgdata/zlib",
		pkg:    upstream.Package{Name: "zlib", Version: "1.3.1"},
	},
	{
		cfgDir: "./cfgdata/bzip3",
		pkg:    upstream.Package{Name: "bzip3", Version: "1.5.1"},
	},
	{
		cfgDir: "./cfgdata/cargs",
		pkg:    upstream.Package{Name: "cargs", Version: "1.2.0"},
	},
	{
		cfgDir: "./cfgdata/bzip2",
		pkg:    upstream.Package{Name: "bzip2", Version: "1.0.8"},
	},
	{
		cfgDir: "./cfgdata/libtool",
		pkg:    upstream.Package{Name: "libtool", Version: "2.4.7"},
	},
}

/*
​​The cfgdata directory is used by llcppcfg to generate end-to-end tests. Its directory structure is as follows.

	cfgdata
	├── libxml2
	│   └── {{OS}}
	│       └── 2.13.6
	└── cjson
	    └── {{OS}}
	        └── 1.7.18

Due to inconsistencies in header file paths across different systems, the expected test files are platform-based.​​
*/
func TestEnd2EndLLCppcfg(t *testing.T) {
	null, err := os.OpenFile(os.DevNull, os.O_RDWR, 0644)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range cfgTestCases {
		t.Run(fmt.Sprintf("%s/%s", tc.pkg.Name, tc.pkg.Version), func(t *testing.T) {
			conanDir, err := os.MkdirTemp("", "llcppg_end2end_test_conan_dir_*")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(conanDir)

			resultDir, err := os.MkdirTemp("", "llcppg_end2end_llcppcfg_gen_result_*")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(resultDir)

			stderr := os.Stderr

			os.Stderr = null
			actualPcFiles, err := conan.NewConanInstaller(tc.config).Install(tc.pkg, conanDir)
			os.Stderr = stderr

			if err != nil {
				t.Fatal(err)
			}

			pcFileName := actualPcFiles[0]

			cmd := exec.Command("llcppcfg", pcFileName)
			cmd.Env = os.Environ()
			cmd.Env = append(cmd.Env, pcPathEnv(conanDir)...)
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout

			platformCfgDir := filepath.Join(tc.cfgDir, runtime.GOOS, tc.pkg.Version)

			gen := false
			// generate config only
			if gen {
				os.MkdirAll(platformCfgDir, 0700)
				cmd.Dir = platformCfgDir

				if err = cmd.Run(); err != nil {
					t.Fatal(err)
				}
				return
			}

			cmd.Dir = resultDir
			if err = cmd.Run(); err != nil {
				t.Fatal(err)
			}

			cmd = exec.Command("git", "diff", "--no-index", resultDir, platformCfgDir)
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout
			err = cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
