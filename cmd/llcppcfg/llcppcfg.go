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
usage: llcppcfg [-cpp|-help|-expand|-sort] libname`)
	flag.PrintDefaults()
}

func main() {
	var cpp, help, expand, sortByDep bool
	flag.BoolVar(&cpp, "cpp", false, "if it is c++ lib")
	flag.BoolVar(&help, "help", false, "print help message")
	flag.BoolVar(&expand, "expand", false, "expand pkg-config command to result")
	flag.BoolVar(&sortByDep, "sort", true, "expand every cflag and list it's include files and sort include files by dependency")
	extsString := ""
	flag.StringVar(&extsString, "exts", ".h", "extra include file extensions for example \"-exts=.h .hpp .hh\"")
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

	fnToCfgExpandMode := func() llcppgcfg.CfgMode {
		if expand {
			return llcppgcfg.ExpandMode
		} else if sortByDep {
			return llcppgcfg.SortMode
		}
		return llcppgcfg.NormalMode
	}
	exts := strings.Fields(extsString)
	buf, err := llcppgcfg.GenCfg(name, cpp, fnToCfgExpandMode(), exts)
	if err != nil {
		log.Fatal(err)
	}
	outFile := "./llcppg.cfg"
	err = os.WriteFile(outFile, buf.Bytes(), 0600)
	if err != nil {
		log.Fatal(err)
	}
}
