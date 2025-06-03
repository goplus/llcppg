package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/goplus/llcppg/config"
	"github.com/goplus/llpkgstore/upstream"
	"github.com/goplus/llpkgstore/upstream/installer/conan"
)

const llcppgGoVersion = "1.20.14"

type testCase struct {
	modpath  string
	dir      string
	pkg      upstream.Package
	config   map[string]string // conan options
	demosDir string
}

var testCases = []testCase{
	{
		modpath: "github.com/goplus/llcppg/_cmptest/testdata/cjson/1.7.18/cjson",
		dir:     "./testdata/cjson/1.7.18",
		pkg:     upstream.Package{Name: "cjson", Version: "1.7.18"},
		config: map[string]string{
			"options": "utils=True",
		},
		demosDir: "./testdata/cjson/demo",
	},
	{
		modpath: "github.com/goplus/llcppg/_cmptest/testdata/cjson/1.7.17/cjson",
		dir:     "./testdata/cjson/1.7.17",
		pkg:     upstream.Package{Name: "cjson", Version: "1.7.17"},
		config: map[string]string{
			"options": "utils=True",
		},
		demosDir: "./testdata/cjson/demo",
	},
}

func TestEnd2End(t *testing.T) {
	for _, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("%s/%s", tc.pkg.Name, tc.pkg.Version), func(t *testing.T) {
			testFrom(t, tc, false)
		})
	}
}

func testFrom(t *testing.T, tc testCase, gen bool) {
	wd, _ := os.Getwd()
	dir := filepath.Join(wd, tc.dir)
	conanDir, err := os.MkdirTemp("", "llcppg_end2end_test_conan_dir_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(conanDir)

	resultDir, err := os.MkdirTemp("", "llcppg_end2end_test_result_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(resultDir)

	cfgPath := filepath.Join(wd, tc.dir, config.LLCPPG_CFG)
	cfg, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatal(err)
	}

	os.WriteFile(filepath.Join(resultDir, config.LLCPPG_CFG), cfg, os.ModePerm)
	_, err = conan.NewConanInstaller(tc.config).Install(tc.pkg, conanDir)
	if err != nil {
		t.Fatal(err)
	}

	cmd := command(resultDir, "llcppg", "-v", "-mod="+tc.modpath)
	cmd.Env = append(cmd.Env, goVerEnv())
	cmd.Env = append(cmd.Env, pcPathEnv(conanDir)...)

	err = cmd.Run()
	if err != nil {
		t.Fatal(err)
	}

	// llcppg.symb.json is a middle file
	os.Remove(filepath.Join(resultDir, config.LLCPPG_SYMB))

	if gen {
		os.RemoveAll(dir)
		os.Rename(resultDir, dir)
	} else {
		// check the result is the same as the expected result
		// when have diff,will got exit code 1
		diffCmd := command(wd, "git", "diff", "--no-index", dir, resultDir)
		err = diffCmd.Run()
		if err != nil {
			t.Fatal(err)
		}
	}
	runDemos(t, filepath.Join(wd, tc.demosDir), tc.pkg.Name, filepath.Join(dir, tc.pkg.Name), conanDir)
}

// pkgpath is the filepath use to replace the import path in demo's go.mod
func runDemos(t *testing.T, demosPath string, pkgname, pkgpath, pcPath string) {
	tempDemosPath, err := os.MkdirTemp("", "llcppg_end2end_test_demos_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDemosPath)
	err = os.CopyFS(tempDemosPath, os.DirFS(demosPath))
	if err != nil {
		t.Fatal(err)
	}

	goMod := command(tempDemosPath, "go", "mod", "init", "test")
	err = goMod.Run()
	if err != nil {
		t.Fatal(err)
	}

	replace := command(tempDemosPath, "go", "mod", "edit", "-replace", pkgname+"="+pkgpath)
	err = replace.Run()
	if err != nil {
		t.Fatal(err)
	}

	tidy := command(tempDemosPath, "go", "mod", "tidy")
	err = tidy.Run()
	if err != nil {
		t.Fatal(err)
	}

	demos, err := os.ReadDir(tempDemosPath)
	if err != nil {
		t.Fatal(err)
	}

	for _, demo := range demos {
		if !demo.IsDir() {
			continue
		}
		demoPath := filepath.Join(tempDemosPath, demo.Name())
		demoCmd := command(demoPath, "llgo", "run", ".")
		demoCmd.Env = append(demoCmd.Env, llgoEnv()...)
		demoCmd.Env = append(demoCmd.Env, pcPathEnv(pcPath)...)
		err = demoCmd.Run()
		if err != nil {
			t.Fatal(err)
		}
	}

}

func appendPCPath(path string) string {
	if env, ok := os.LookupEnv("PKG_CONFIG_PATH"); ok {
		return path + ":" + env
	}
	return path
}

// llgo env
func llgoEnv() []string {
	return []string{
		// for https://github.com/goplus/llgo/issues/1135
		"LLGO_RPATH_CHANGE=on",
	}
}

// env for pkg-config
func pcPathEnv(path string) []string {
	pcPath := fmt.Sprintf("PKG_CONFIG_PATH=%s", appendPCPath(path))
	return append(os.Environ(), pcPath)
}

// control the go version in output version
func goVerEnv() string {
	return fmt.Sprintf("GOTOOLCHAIN=go%s", llcppgGoVersion)
}

func command(dir string, app string, args ...string) *exec.Cmd {
	cmd := exec.Command(app, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}
