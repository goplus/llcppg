package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/goplus/llcppg/cmd/llcppcfg/gen"
)

func printHelp() {
	log.Println(`llcppcfg is to generate llcppg.cfg file.
usage: llcppcfg [-cpp|-tab|-excludes|-exts|-help] libname`)
	flag.PrintDefaults()
}

func main() {
	var dependencies string
	var cpp, help, tab bool
	flag.BoolVar(&cpp, "cpp", false, "if it is C++ lib")
	flag.BoolVar(&help, "help", false, "print help message")
	flag.BoolVar(&tab, "tab", true, "generate .cfg config file with tab indent")
	extsString := ""
	flag.StringVar(&extsString, "exts", ".h", "extra include file extensions for example -exts=\".h .hpp .hh\"")
	excludes := ""
	flag.StringVar(&excludes, "excludes", "", "exclude all header files in subdir of include expamle -excludes=\"internal impl\"")
	flag.StringVar(&dependencies, "deps", "", "deps for autofilling dependencies")
	flag.Usage = printHelp
	flag.Parse()
	if help || len(os.Args) <= 1 {
		flag.Usage()
		return
	}
	name := ""
	if len(flag.Args()) > 0 {
		name = flag.Arg(0)
	}

	exts := strings.Fields(extsString)
	deps := strings.Fields(dependencies)

	excludeSubdirs := []string{}
	if len(excludes) > 0 {
		excludeSubdirs = strings.Fields(excludes)
	}
	var flag gen.FlagMode
	if cpp {
		flag |= gen.WithCpp
	}
	if tab {
		flag |= gen.WithTab
	}
	buf, err := gen.Do(gen.NewConfig(name, flag, exts, deps, excludeSubdirs))
	if err != nil {
		panic(err)
	}
	outFile := "./llcppg.cfg"
	err = os.WriteFile(outFile, buf, 0600)
	if err != nil {
		panic(err)
	}
}
