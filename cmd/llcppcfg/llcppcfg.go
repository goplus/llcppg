package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/goplus/llcppg/cmd/llcppcfg/llcppgcfg"
)

func printHelp() {
	log.Println(`llcppcfg is to generate llcppg.cfg file.
usage: llcppcfg [-cpp|-tab|-excludes|-exts|-help] libname`)
	flag.PrintDefaults()
}

func main() {
	var cpp, help, tab bool
	flag.BoolVar(&cpp, "cpp", false, "if it is c++ lib")
	flag.BoolVar(&help, "help", false, "print help message")
	flag.BoolVar(&tab, "tab", true, "generate .cfg config file with tab indent")
	extsString := ""
	flag.StringVar(&extsString, "exts", ".h", "extra include file extensions for example -exts=\".h .hpp .hh\"")
	excludes := ""
	flag.StringVar(&excludes, "excludes", "", "exclude all header files in subdir of include expamle -excludes=\"internal impl\"")
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
	excludeSubdirs := []string{}
	if len(excludes) > 0 {
		excludeSubdirs = strings.Fields(excludes)
	}
	var flag llcppgcfg.FlagMode
	if cpp {
		flag |= llcppgcfg.WithCpp
	}
	if tab {
		flag |= llcppgcfg.WithTab
	}
	buf, err := llcppgcfg.GenCfg(llcppgcfg.NewGenConfig(name, flag, exts, excludeSubdirs))
	if err != nil {
		log.Fatal(err)
	}
	outFile := "./llcppg.cfg"
	err = os.WriteFile(outFile, buf.Bytes(), 0600)
	if err != nil {
		log.Fatal(err)
	}
}
