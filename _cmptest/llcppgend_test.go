package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"github.com/goplus/llcppg/config"
	"github.com/goplus/llpkgstore/upstream"
	"github.com/goplus/llpkgstore/upstream/installer/conan"
)

const llcppgGoVersion = "1.20.14"

var (
	// avoid conan install race condition
	conanInstallMutex sync.Mutex
	// avoid llgo run race condition
	llgoRunMutex sync.Mutex
)

type testCase struct {
	modpath  string
	dir      string
	cfgDir   string
	pkg      upstream.Package
	config   map[string]string // conan options
	demosDir string
}

var testCases = []testCase{
	{
		modpath: "github.com/goplus/llcppg/_cmptest/testdata/cjson/1.7.18/cjson",
		dir:     "./testdata/cjson/1.7.18",
		cfgDir:  "./cfgdata/cjson",
		pkg:     upstream.Package{Name: "cjson", Version: "1.7.18"},
		config: map[string]string{
			"options": "utils=True",
		},
		demosDir: "./testdata/cjson/demo",
	},
	{
		modpath: "github.com/goplus/llcppg/_cmptest/testdata/cjson/1.7.17/cjson",
		dir:     "./testdata/cjson/1.7.17",
		cfgDir:  "./cfgdata/cjson",
		pkg:     upstream.Package{Name: "cjson", Version: "1.7.17"},
		config: map[string]string{
			"options": "utils=True",
		},
		demosDir: "./testdata/cjson/demo",
	},
	{
		modpath: "github.com/goplus/llcppg/_cmptest/testdata/libxml2/2.13.6/libxml2",
		dir:     "./testdata/libxml2/2.13.6",
		cfgDir:  "./cfgdata/libxml2",
		pkg:     upstream.Package{Name: "libxml2", Version: "2.13.6"},
		config: map[string]string{
			"options": "iconv=False",
		},
		demosDir: "./testdata/libxml2/demo",
	},
	{
		modpath:  "github.com/goplus/llcppg/_cmptest/testdata/sqlite3/3.49.1/sqlite3",
		dir:      "./testdata/sqlite3/3.49.1",
		cfgDir:   "./cfgdata/sqlite3",
		pkg:      upstream.Package{Name: "sqlite3", Version: "3.49.1"},
		demosDir: "./testdata/sqlite3/demo",
	},
	{
		modpath:  "github.com/goplus/llcppg/_cmptest/testdata/zlib/1.3.1/zlib",
		dir:      "./testdata/zlib/1.3.1",
		cfgDir:   "./cfgdata/zlib",
		pkg:      upstream.Package{Name: "zlib", Version: "1.3.1"},
		demosDir: "./testdata/zlib/demo",
	},
	{
		modpath:  "github.com/goplus/llcppg/_cmptest/testdata/bzip3/1.5.1/bzip3",
		dir:      "./testdata/bzip3/1.5.1",
		cfgDir:   "./cfgdata/bzip3",
		pkg:      upstream.Package{Name: "bzip3", Version: "1.5.1"},
		demosDir: "./testdata/bzip3/demo",
	},
	{
		modpath:  "github.com/goplus/llcppg/_cmptest/testdata/cargs/1.2.0/cargs",
		dir:      "./testdata/cargs/1.2.0",
		cfgDir:   "./cfgdata/cargs",
		pkg:      upstream.Package{Name: "cargs", Version: "1.2.0"},
		demosDir: "./testdata/cargs/demo",
	},
	{
		modpath:  "github.com/goplus/llcppg/_cmptest/testdata/bzip2/1.0.8/bzip2",
		dir:      "./testdata/bzip2/1.0.8",
		cfgDir:   "./cfgdata/bzip2",
		pkg:      upstream.Package{Name: "bzip2", Version: "1.0.8"},
		demosDir: "./testdata/bzip2/demo",
	},
}

var mkdirTempLazily = sync.OnceValue(func() string {
	if env := os.Getenv("LLCPPG_TEST_LOG_DIR"); env != "" {
		return env
	}
	dir, err := os.MkdirTemp("", "test-log")
	if err != nil {
		panic(err)
	}
	return dir
})

func logFile(tc testCase) (*os.File, error) {
	caseName := fmt.Sprintf("%s-%s-llcppg-%s-%s", runtime.GOOS, runtime.GOARCH, tc.pkg.Name, tc.pkg.Version)
	dirPath := filepath.Join(mkdirTempLazily(), caseName)

	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return nil, err
	}

	return os.Create(filepath.Join(dirPath, fmt.Sprintf("%s.log", caseName)))
}

func TestEnd2End(t *testing.T) {
	for _, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("%s/%s", tc.pkg.Name, tc.pkg.Version), func(t *testing.T) {
			t.Parallel()
			testFrom(t, tc, false)
		})
	}
}

func TestEnd2EndLLCppcfg(t *testing.T) {
	null, err := os.OpenFile(os.DevNull, os.O_RDWR, 0644)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range testCases {
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
			if false {

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

func testFrom(t *testing.T, tc testCase, gen bool) {
	logFile, err := logFile(tc)
	if err != nil {
		t.Fatal(err)
	}
	defer logFile.Close()
	fmt.Printf("%s:%s log file: %s\n", tc.pkg.Name, tc.pkg.Version, logFile.Name())

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

	conanInstallMutex.Lock()
	_, err = conan.NewConanInstaller(tc.config).Install(tc.pkg, conanDir)
	conanInstallMutex.Unlock()
	if err != nil {
		t.Fatal(err)
	}

	cmd := command(logFile, resultDir, "llcppg", "-v", "-mod="+tc.modpath)
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
		diffCmd := command(logFile, wd, "git", "diff", "--no-index", dir, resultDir)
		err = diffCmd.Run()
		if err != nil {
			t.Fatal(err)
		}
	}
	runDemos(t, logFile, filepath.Join(wd, tc.demosDir), tc.pkg.Name, filepath.Join(dir, tc.pkg.Name), conanDir)
}

// pkgpath is the filepath use to replace the import path in demo's go.mod
func runDemos(t *testing.T, logFile *os.File, demosPath string, pkgname, pkgpath, pcPath string) {
	tempDemosPath, err := os.MkdirTemp("", "llcppg_end2end_test_demos_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDemosPath)
	err = os.CopyFS(tempDemosPath, os.DirFS(demosPath))
	if err != nil {
		t.Fatal(err)
	}

	goMod := command(logFile, tempDemosPath, "go", "mod", "init", "test")
	err = goMod.Run()
	if err != nil {
		t.Fatal(err)
	}

	replace := command(logFile, tempDemosPath, "go", "mod", "edit", "-replace", pkgname+"="+pkgpath)
	err = replace.Run()
	if err != nil {
		t.Fatal(err)
	}

	tidy := command(logFile, tempDemosPath, "go", "mod", "tidy")
	err = tidy.Run()
	if err != nil {
		t.Fatal(err)
	}

	demos, err := os.ReadDir(tempDemosPath)
	if err != nil {
		t.Fatal(err)
	}

	// only can lock out of loop,will got ld64.lld: error: undefined symbol: math.Float32bits
	llgoRunMutex.Lock()
	defer llgoRunMutex.Unlock()

	for _, demo := range demos {
		if !demo.IsDir() {
			continue
		}
		demoPath := filepath.Join(tempDemosPath, demo.Name())
		demoCmd := command(logFile, demoPath, "llgo", "run", ".")
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

func command(logFile *os.File, dir string, app string, args ...string) *exec.Cmd {
	cmd := exec.Command(app, args...)
	cmd.Dir = dir
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	return cmd
}
