package demo

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	llcppg "github.com/goplus/llcppg/config"
)

// runSingleDemo tests a single LLCPPG conversion case in the given demo directory.
// The testing process includes:
// 1. Reading and validating the llcppg.cfg configuration file
// 2. Running LLCPPG to generate the converted package
// 3. Verifying the generated package can be built using llgo
// 4. Running example programs in the demo subdirectory that use the generated package
//
// Directory structure (take _llcppgtest/lua as an example):
// _llcppgtest/lua/           - Demo root directory
//
//	├── llcppg.cfg          - LLCPPG configuration file
//	├── out/                - Generated package output directory
//	└── demo/               - Example programs directory
//	    ├── example1/       - First example program
//	    └── example2/       - Second example program
//
// The function will panic if any step in the testing process fails.
//
// Parameters:
//   - demoRoot: Path to the root directory of a single demo case
//   - confDir: Path to the configuration directory relative to demoRoot, defaults to "." if empty
func RunGenPkgDemo(demoRoot string, confDir string) {
	fmt.Printf("Testing demo: %s\n", demoRoot)

	absPath, err := filepath.Abs(demoRoot)
	if err != nil {
		panic(fmt.Sprintf("failed to get absolute path for %s: %v", demoRoot, err))
	}
	demoPkgName := filepath.Base(absPath)

	if confDir == "" {
		confDir = "."
	}

	configPath := filepath.Join(absPath, confDir)
	configFile := filepath.Join(configPath, llcppg.LLCPPG_CFG)
	fmt.Printf("Looking for config file at: %s\n", configFile)

	if _, err = os.Stat(configFile); os.IsNotExist(err) {
		panic(fmt.Sprintf("config file not found: %s", configFile))
	}

	llcppgArgs := []string{"-v", "-mod", demoPkgName}

	outDir := filepath.Join(absPath, "out")
	if err = os.MkdirAll(outDir, 0755); err != nil {
		panic(fmt.Sprintf("failed to create output directory: %v", err))
	}
	defer os.RemoveAll(outDir)

	// copy configs to out dir
	cfgFiles := []string{llcppg.LLCPPG_CFG, llcppg.LLCPPG_PUB, llcppg.LLCPPG_SYMB}
	for _, cfg := range cfgFiles {
		src := filepath.Join(configPath, cfg)
		dst := filepath.Join(outDir, cfg)
		var content []byte
		content, err = os.ReadFile(src)
		if err != nil {
			if os.IsNotExist(err) && cfg != llcppg.LLCPPG_CFG {
				continue
			}
			panic(fmt.Sprintf("failed to read config file: %v", err))
		}
		if err = os.WriteFile(dst, content, 0600); err != nil {
			panic(fmt.Sprintf("failed to write config file: %v", err))
		}
	}

	// run llcppg to gen pkg
	if err = runCommand(outDir, "llcppg", llcppgArgs...); err != nil {
		panic(fmt.Sprintf("llcppg execution failed: %v", err))
	}
	fmt.Printf("llcppg execution success\n")

	// check if the gen pkg is ok
	genPkgDir := filepath.Join(outDir, demoPkgName)
	if err = runCommand(genPkgDir, "go", "fmt"); err != nil {
		panic(fmt.Sprintf("go fmt failed in %s: %v", genPkgDir, err))
	}

	if err = runCommand(genPkgDir, "llgo", "build", "."); err != nil {
		panic(fmt.Sprintf("llgo build failed in %s: %v", genPkgDir, err))
	}
	fmt.Printf("llgo build success\n")

	demosPath := filepath.Join(demoRoot, "demo")
	// init mods to test package,because the demo is dependent on the gen pkg
	if err = runCommand(demoRoot, "go", "mod", "init", "demo"); err != nil {
		panic(fmt.Sprintf("go mod init failed in %s: %v", demoRoot, err))
	}
	if err = runCommand(demoRoot, "go", "mod", "edit", "-replace", demoPkgName+"="+"./out/"+demoPkgName); err != nil {
		panic(fmt.Sprintf("go mod edit failed in %s: %v", demoRoot, err))
	}
	if err = runCommand(demoRoot, "go", "mod", "tidy"); err != nil {
		panic(fmt.Sprintf("go mod tidy failed in %s: %v", demoRoot, err))
	}
	defer os.Remove(filepath.Join(absPath, "go.mod"))
	defer os.Remove(filepath.Join(absPath, "go.sum"))

	fmt.Printf("testing demos in %s\n", demosPath)
	// run the demo
	var demos []os.DirEntry
	demos, err = os.ReadDir(demosPath)
	if err != nil {
		panic(fmt.Sprintf("failed to read demo directory: %v", err))
	}
	for _, demo := range demos {
		if demo.IsDir() {
			fmt.Printf("Running demo: %s\n", demo.Name())
			if demoErr := runCommand(filepath.Join(demosPath, demo.Name()), "llgo", "run", "."); demoErr != nil {
				panic(fmt.Sprintf("failed to run demo: %s: %v", demo.Name(), demoErr))
			}
		}
	}
}

// Get all first-level directories containing llcppg.cfg
func getFirstLevelDemos(baseDir string, confDir string) []string {
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		panic(fmt.Sprintf("failed to read directory: %v", err))
	}

	var demos []string
	for _, entry := range entries {
		if entry.IsDir() {
			demoRoot := filepath.Join(baseDir, entry.Name())
			configPath := filepath.Join(demoRoot, confDir, llcppg.LLCPPG_CFG)
			if _, err := os.Stat(configPath); err == nil {
				demos = append(demos, demoRoot)
			}
		}
	}
	return demos
}

func RunAllGenPkgDemos(baseDir string, confDir string) {
	fmt.Printf("Starting generated package tests in directory: %s\n", baseDir)

	stat, err := os.Stat(baseDir)
	if err != nil || !stat.IsDir() {
		panic(fmt.Sprintf("specified path is not a directory or does not exist: %s", baseDir))
	}

	demos := getFirstLevelDemos(baseDir, confDir)
	if len(demos) == 0 {
		panic(fmt.Sprintf("no directories containing llcppg.cfg found in %s", baseDir))
	}

	// Test each demo
	for _, demo := range demos {
		RunGenPkgDemo(demo, confDir)
	}
	fmt.Println("All generated package demos passed:", strings.Join(demos, ","))
}

func runCommand(dir, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
