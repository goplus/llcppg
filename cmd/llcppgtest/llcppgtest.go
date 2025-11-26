package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/goplus/llcppg/cmd/llcppgtest/demo"
)

type FlagName string

const (
	RFlagName        FlagName = "r"
	RandFlagName     FlagName = "rand"
	AFlagName        FlagName = "a"
	AllFlagName      FlagName = "all"
	VFlagName        FlagName = "v"
	VfetchFlagName   FlagName = "vfetch"
	VsymgFlagName    FlagName = "vsymg"
	VgogenFlagName   FlagName = "vgogen"
	CppFlagName      FlagName = "cpp"
	ExtsFlagName     FlagName = "exts"
	ExcludesFlagName FlagName = "excludes"
	HFlagName        FlagName = "h"
	HelpFlagName     FlagName = "help"
	DemosFlagName    FlagName = "demos"
	DemoFlagName     FlagName = "demo"
	DemoConfName     FlagName = "conf"
)

func RunCommandWithOut(out *io.PipeWriter, dir, cmdName string, args ...string) {
	defer out.Close()
	cmd := exec.Command(cmdName, args...)
	cmd.Stdout = out
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = dir
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func RunCommandInDir(dir string, done func(error), name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Dir = dir
	err := cmd.Run()
	if done != nil {
		done(err)
	}
}

func RunCommand(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func PkgList(r io.Reader) []string {
	var pkgs []string
	scan := bufio.NewScanner(r)
	for scan.Scan() {
		lineBuf := bytes.NewBufferString(scan.Text())
		lineScan := bufio.NewScanner(lineBuf)
		lineScan.Split(bufio.ScanWords)
		firstWord := ""
		for lineScan.Scan() {
			text := lineScan.Text()
			if len(firstWord) == 0 {
				firstWord = text
				pkgs = append(pkgs, firstWord)
			}
		}
	}
	return pkgs
}

func getPkgs() []string {
	wd, _ := os.Getwd()
	r, w := io.Pipe()
	go RunCommandWithOut(w, wd, "pkg-config", "--list-all")
	pkgs := PkgList(r)
	return pkgs
}

type runPkgMode int

const (
	withCpp runPkgMode = 1 << iota
	withSigfetchVerbose
	withSymgVerbose
	withExts
	withExcludes
)

type runAppMode int

const (
	runSelected runAppMode = iota
	runRand
	runAll
	runDemos
	runDemo
)

type RunConfig struct {
	runMode  runPkgMode
	exts     string
	excludes string
}

func NewRunConfig(runMode runPkgMode, exts, excludes string) *RunConfig {
	return &RunConfig{runMode: runMode, exts: exts, excludes: excludes}
}

func runPkgs(pkgs []string, cfg *RunConfig) {
	wd, _ := os.Getwd()
	wg := sync.WaitGroup{}
	wg.Add(len(pkgs))
	llcppcfgArg := []string{}
	if cfg.runMode&withCpp != 0 {
		llcppcfgArg = append(llcppcfgArg, fmt.Sprintf("-%s", CppFlagName))
	}
	if cfg.runMode&withExts != 0 {
		llcppcfgArg = append(llcppcfgArg, fmt.Sprintf("-%s=%s", ExtsFlagName, cfg.exts))
	}
	if cfg.runMode&withExcludes != 0 {
		llcppcfgArg = append(llcppcfgArg, fmt.Sprintf("-%s=%s", ExcludesFlagName, cfg.excludes))
	}
	llcppgArg := []string{}
	if cfg.runMode&withSigfetchVerbose != 0 {
		llcppgArg = append(llcppgArg, fmt.Sprintf("-%s", VfetchFlagName))
	}
	if cfg.runMode&withSymgVerbose != 0 {
		llcppgArg = append(llcppgArg, fmt.Sprintf("-%s", VsymgFlagName))
	}
	var runs []string
	for _, pkg := range pkgs {
		dir := "./out/" + pkg
		RunCommand("mkdir", "-p", dir)
		RunCommand("cd", dir)
		curDir := wd + "/out/" + pkg
		RunCommandInDir(curDir, func(error) {
			runs = append(runs, pkg)
			go RunCommandInDir(curDir, func(error) {
				wg.Done()
			}, "llcppg", llcppgArg...)
		}, "llcppcfg", append(llcppcfgArg, pkg)...)
	}
	wg.Wait()
	fmt.Printf("llcppgtest run %v finished!\n", runs)
}

func randIndex(maxInt int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(maxInt)
}

func runPkg(cfg *RunConfig) {
	pkgs := getPkgs()
	idx := randIndex(len(pkgs))
	pkg := pkgs[idx]
	fmt.Printf("***start test %s\n", pkg)
	runPkgs([]string{pkg}, cfg)
}

func printHelp() {
	helpString := fmt.Sprintf(`
	llcppgtest is used to test llcppg
usage: llcppgtest [-%s|-%s|-%s|-%s] [-%s|-%s|-%s] [-%s|-%s|-%s] [-%s|-%s] pkgname
	   llcppgtest -%s <path>    # test all first-level demo directories
       llcppgtest -%s <path>     # test specific demo directory
`, RFlagName, RandFlagName, AFlagName, AllFlagName, VFlagName,
		VfetchFlagName, VsymgFlagName, CppFlagName, ExtsFlagName,
		ExcludesFlagName, HFlagName, HelpFlagName, DemosFlagName, DemoFlagName)
	fmt.Println(helpString)
	flag.PrintDefaults()
}

func main() {
	rand := false
	flag.BoolVar(&rand, string(RFlagName), false, "same as -rand")
	flag.BoolVar(&rand, string(RandFlagName), false, "select one pkg of pkg-config --list-all to test")
	all := false
	flag.BoolVar(&all, string(AFlagName), false, "same as -all")
	flag.BoolVar(&all, string(AllFlagName), false, "test all pkgs of pkg-config --list-all")
	v := false
	flag.BoolVar(&v, string(VFlagName), false, "enable verbose of llcppsigfetch and llcppsymg")
	vSig := false
	flag.BoolVar(&vSig, string(VfetchFlagName), false, "enable verbose of llcppsigfetch")
	vSym := false
	flag.BoolVar(&vSym, string(VsymgFlagName), false, "enable verbose of llcppsymg")
	vGogen := false
	flag.BoolVar(&vGogen, string(VgogenFlagName), false, "enable verbose of gogensig")
	cpp := false
	flag.BoolVar(&cpp, string(CppFlagName), false, "if it is a cpp library")
	exts := ""
	flag.StringVar(&exts, string(ExtsFlagName), ".h", "for all headers with ext of exts to generate .go, for example -exts=\".h .hh .cpp .hpp\"")
	excludes := ""
	flag.StringVar(&excludes, string(ExcludesFlagName), "", "for all internal implementation directors that you want to excludes from -I include director to handle. For example -excludes=\"internal impl\"")
	help := false
	flag.BoolVar(&help, string(HFlagName), false, "print help message")
	flag.BoolVar(&help, string(HelpFlagName), false, "print help message")
	demosPath := flag.String(string(DemosFlagName), "", "test all first-level demo directories in the specified path")
	demoPath := flag.String(string(DemoFlagName), "", "test the specified demo directory")
	confPath := flag.String(string(DemoConfName), "", "path to the config directory relative to demo directory, defaults to \".\" if empty")
	flag.Parse()

	if help || len(os.Args) == 1 {
		printHelp()
		return
	}

	runMode := 0
	if cpp {
		runMode |= int(withCpp)
	}
	if vSig {
		runMode |= int(withSigfetchVerbose)
	}
	if vSym {
		runMode |= int(withSymgVerbose)
	}
	if v {
		runMode |= int(withSigfetchVerbose)
		runMode |= int(withSymgVerbose)
	}

	appMode := runSelected
	if rand {
		appMode = runRand
	}
	if all {
		appMode = runAll
	}

	if *demosPath != "" {
		appMode = runDemos
	}
	if *demoPath != "" {
		appMode = runDemo
	}

	cfg := NewRunConfig(runPkgMode(runMode), exts, excludes)

	switch {
	case appMode == runRand:
		runPkg(cfg)
	case appMode == runAll:
		pkgs := getPkgs()
		runPkgs(pkgs, cfg)
	case appMode == runDemos:
		if err := demo.RunAllGenPkgDemos(*demosPath, *confPath); err != nil {
			panic(err)
		}
	case appMode == runDemo:
		if err := demo.RunGenPkgDemo(*demoPath, *confPath); err != nil {
			panic(err)
		}
	default:
		if len(flag.Args()) > 0 {
			arg := flag.Arg(0)
			fmt.Printf("***start test %s\n", arg)
			runPkgs([]string{arg}, cfg)
		} else {
			printHelp()
		}
	}
}
