package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/goplus/llcppg/cmd/llcppcfg/llcppgcfg"
)

func printUsageLog() {
	fmt.Println(`llcppcfg is to generate llcpp.cfg file.
usage: llcppcfg [-cpp|-help|-expand] libname`)
}

func printHelp() {
	printUsageLog()
	flag.PrintDefaults()
}

func main() {
	var cpp, help, expand bool
	flag.BoolVar(&cpp, "cpp", false, "if it is c++ lib")
	flag.BoolVar(&help, "help", false, "print help message")
	flag.BoolVar(&expand, "expand", false, "expand pkg-config command to result")
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
	buf, err := llcppgcfg.GenCfg(name, cpp, expand)
	if err != nil {
		log.Fatal(err)
	}
	outFile := "./llcppg.cfg"
	err = os.WriteFile(outFile, buf.Bytes(), 0600)
	if err != nil {
		log.Fatal(err)
	}
}
