package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/goplus/llcppg/cmd/llcppgtest/demo"
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

func RemoveFile(fileName string) {
	cmd := exec.Command("rm", fileName)
	outBuf := bytes.NewBufferString("")
	cmd.Stderr = outBuf
	cmd.Stdout = outBuf
	err := cmd.Run()
	if err != nil {
		log.Println("")
	}
}

func PkgList(r io.Reader) []string {
	pkgs := make([]string, 0)
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
	withSort
	withSigfetchVerbose
	withSymgVerbose
)

type runAppMode int

const (
	runSelected runAppMode = iota
	runRand
	runAll
	runDemos
)

func runPkgs(pkgs []string, runMode runPkgMode) {
	wd, _ := os.Getwd()
	wg := sync.WaitGroup{}
	wg.Add(len(pkgs))
	llcppcfgArg := []string{}
	if runMode&withCpp != 0 {
		llcppcfgArg = append(llcppcfgArg, "-cpp")
	}
	if runMode&withSort != 0 {
		llcppcfgArg = append(llcppcfgArg, "-sort")
	}
	llcppgArg := []string{}
	if runMode&withSigfetchVerbose != 0 {
		llcppgArg = append(llcppgArg, "-vfetch")
	}
	if runMode&withSymgVerbose != 0 {
		llcppgArg = append(llcppgArg, "-vsym")
	}
	runs := make([]string, 0)
	for _, pkg := range pkgs {
		dir := "./out/" + pkg
		llcppsymgFile := filepath.Join(dir, "llcppg.symb.json")
		RemoveFile(llcppsymgFile)
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

func runPkg(runMode runPkgMode) {
	pkgs := getPkgs()
	idx := randIndex(len(pkgs))
	pkg := pkgs[idx]
	fmt.Printf("***start test %s\n", pkg)
	runPkgs([]string{pkg}, runMode)
}

func printHelp() {
	helpString := `llcppgtest is used to test llcppg
usage: llcppgtest [-r|-rand|-a|-all] [-v|-vfetch|-vsym] [-cpp|-sort] [-h|-help] pkgname`
	fmt.Println(helpString)
	flag.PrintDefaults()
}

func main() {
	rand := false
	flag.BoolVar(&rand, "r", false, "same as -rand")
	flag.BoolVar(&rand, "rand", false, "select one pkg of pkg-config --list-all to test")
	all := false
	flag.BoolVar(&all, "a", false, "same as -all")
	flag.BoolVar(&all, "all", false, "test all pkgs of pkg-config --list-all")
	v := false
	flag.BoolVar(&v, "v", false, "enable verbose of llcppsigfetch and llcppsymg")
	vSig := false
	flag.BoolVar(&vSig, "vfetch", false, "enable verbose of llcppsigfetch")
	vSym := false
	flag.BoolVar(&vSym, "vsym", false, "enable verbose of llcppsymg")
	cpp := false
	flag.BoolVar(&cpp, "cpp", false, "if it is a cpp library")
	sortIncludes := false
	flag.BoolVar(&sortIncludes, "sort", false, "sort includes in llcppg.cfg")
	help := false
	flag.BoolVar(&help, "h", false, "print help message")
	flag.BoolVar(&help, "help", false, "print help message")
	demoPath := flag.String("demos", "", "test all first-level demo directories in the specified path")
	flag.Parse()

	if help || len(os.Args) == 1 {
		printHelp()
		return
	}

	runMode := 0
	if cpp {
		runMode |= int(withCpp)
	}
	if sortIncludes {
		runMode |= int(withSort)
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

	if *demoPath != "" {
		appMode = runDemos
	}
	switch {
	case appMode == runRand:
		runPkg(runPkgMode(runMode))
	case appMode == runAll:
		pkgs := getPkgs()
		runPkgs(pkgs, runPkgMode(runMode))
	case appMode == runDemos:
		demo.TestDemos(*demoPath)
	default:
		if len(flag.Args()) > 0 {
			arg := flag.Arg(0)
			fmt.Printf("***start test %s\n", arg)
			runPkgs([]string{arg}, runPkgMode(runMode))
		} else {
			printHelp()
		}
	}
}
