package demo

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	llcppg "github.com/goplus/llcppg/config"
)

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

func demoLogDir(demoDir string) string {
	dirName := fmt.Sprintf("%s-%s-llcppg-%s", runtime.GOOS, runtime.GOARCH, filepath.Base(demoDir))
	return filepath.Join(mkdirTempLazily(), dirName)
}

func logFile(demoDir string) (*os.File, error) {
	dirName := demoLogDir(demoDir)

	err := os.MkdirAll(dirName, 0755)
	if err != nil {
		return nil, err
	}

	return os.Create(filepath.Join(dirName, "all.log"))
}

// runSingleDemo tests a single LLCPPG conversion case in the given demo directory.
// The testing process includes:
//  1. Reading and validating the llcppg.cfg configuration file
//     if the current OS is not macOS, it will look up the `conf/{OS}/llcppg.cfg` and skip if llcppg.cfg dones't exist.
//  2. Running LLCPPG to generate the converted package
//  3. Verifying the generated package can be built using llgo
//  4. Running example programs in the demo subdirectory that use the generated package
//
// Directory structure (take _llcppgtest/lua as an example):
// _llcppgtest/lua/           - Demo root directory
//
//		├── llcppg.cfg          - LLCPPG configuration file
//	 	├── conf         		- Configuration directory for speficied platforms
//		|	└── linux			- Linux Platform
//		| 	  	 └── llcppg.cfg - LLCPPG configuration file (linux platform)
//		├── out/                - Generated package output directory
//		└── demo/               - Example programs directory
//		    ├── example1/       - First example program
//		    └── example2/       - Second example program
//
// The function will panic if any step in the testing process fails.
//
// Parameters:
//   - demoRoot: Path to the root directory of a single demo case
//   - confDir: Path to the configuration directory relative to demoRoot, defaults to "." if empty
func RunGenPkgDemo(demoRoot string, confDir string) (paniced bool) {
	fmt.Printf("Testing demo: %s\n", demoRoot)

	absPath, err := filepath.Abs(demoRoot)
	if err != nil {
		log.Panicf("failed to get absolute path for %s: %v", demoRoot, err)
	}
	demoPkgName := filepath.Base(absPath)

	tempLog, err := logFile(demoRoot)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s: log file: %s\n", demoPkgName, tempLog.Name())

	defer func() {
		paniced = recover() != nil
	}()

	if runtime.GOOS == "linux" && confDir == "" {
		confDir = filepath.Join("conf", "linux")
	}

	if confDir == "" {
		confDir = "."
	}

	configPath := filepath.Join(absPath, confDir)
	configFile := filepath.Join(configPath, llcppg.LLCPPG_CFG)
	fmt.Printf("Looking for config file at: %s\n", configFile)

	if _, err = os.Stat(configFile); os.IsNotExist(err) {
		log.Panicf("config file not found: %s", configFile)
	}

	llcppgArgs := []string{"-v", "-mod", demoPkgName}

	outDir := filepath.Join(absPath, "out")
	if err = os.MkdirAll(outDir, 0755); err != nil {
		log.Panicf("failed to create output directory: %v", err)
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
			log.Panicf("%s: failed to read config file: %v", demoPkgName, err)
		}
		if err = os.WriteFile(dst, content, 0600); err != nil {
			log.Panicf("%s: failed to write config file: %v", demoPkgName, err)
		}
	}

	// run llcppg to gen pkg
	if err = runCommand(tempLog, outDir, "llcppg", llcppgArgs...); err != nil {
		log.Panicf("%s: llcppg execution failed: %v", demoPkgName, err)
	}
	fmt.Printf("%s: llcppg execution success\n", demoPkgName)

	// check if the gen pkg is ok
	genPkgDir := filepath.Join(outDir, demoPkgName)
	if err = runCommand(tempLog, genPkgDir, "go", "fmt"); err != nil {
		log.Panicf("%s: go fmt failed in %s: %v", demoPkgName, genPkgDir, err)
	}

	if err = runCommand(tempLog, genPkgDir, "llgo", "build", "."); err != nil {
		log.Panicf("%s: llgo build failed in %s: %v", demoPkgName, genPkgDir, err)
	}
	fmt.Printf("%s: llgo build success\n", demoPkgName)

	demosPath := filepath.Join(demoRoot, "demo")
	// init mods to test package,because the demo is dependent on the gen pkg
	if err = runCommand(tempLog, demoRoot, "go", "mod", "init", "demo"); err != nil {
		log.Panicf("go mod init failed in %s: %v", demoRoot, err)
	}
	if err = runCommand(tempLog, demoRoot, "go", "mod", "edit", "-replace", demoPkgName+"="+"./out/"+demoPkgName); err != nil {
		log.Panicf("go mod edit failed in %s: %v", demoRoot, err)
	}
	if err = runCommand(tempLog, demoRoot, "go", "mod", "tidy"); err != nil {
		log.Panicf("go mod tidy failed in %s: %v", demoRoot, err)
	}
	defer os.Remove(filepath.Join(absPath, "go.mod"))
	defer os.Remove(filepath.Join(absPath, "go.sum"))

	fmt.Printf("testing demos in %s\n", demosPath)
	// run the demo
	var demos []os.DirEntry
	demos, err = os.ReadDir(demosPath)
	if err != nil {
		log.Panicf("%s: failed to read demo directory: %v", demoPkgName, err)
	}
	for _, demo := range demos {
		if demo.IsDir() {
			fmt.Printf("%s: Running demo: %s\n", demoPkgName, demo.Name())
			if demoErr := runCommand(tempLog, filepath.Join(demosPath, demo.Name()), "llgo", "run", "."); demoErr != nil {
				log.Panicf("%s: failed to run demo: %s: %v", demoPkgName, demo.Name(), demoErr)
			}
		}
	}

	return
}

// Get all first-level directories containing llcppg.cfg
func getFirstLevelDemos(baseDir string, confDir string) []string {
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		log.Panicf("failed to read directory: %v", err)
	}

	if runtime.GOOS == "linux" && confDir == "" {
		confDir = filepath.Join("conf", "linux")
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
		log.Panicf("specified path is not a directory or does not exist: %s", baseDir)
	}

	demos := getFirstLevelDemos(baseDir, confDir)
	if len(demos) == 0 {
		log.Panicf("no directories containing llcppg.cfg found in %s", baseDir)
	}

	failedDemosCh := make(chan string, len(demos))
	// Test each demo
	for _, demo := range demos {
		demo := demo

		go func() {
			if paniced := RunGenPkgDemo(demo, confDir); !paniced {
				failedDemosCh <- ""
			} else {
				failedDemosCh <- demo
			}
		}()
	}

	var failedDemos []string
	for range demos {
		demoDir := <-failedDemosCh

		if demoDir != "" {
			failedDemos = append(failedDemos, demoDir)
		}
	}

	if len(failedDemos) > 0 {
		fmt.Println("Failed generated package demos:", strings.Join(failedDemos, ","))
	}
}

func runCommand(logFile *os.File, dir, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	return cmd.Run()
}
