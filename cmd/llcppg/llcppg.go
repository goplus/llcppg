/*
 * Copyright (c) 2026 The XGo Authors (xgo.dev). All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/goplus/llcppg/internal/docserver"
)

var (
	verbose  = flag.Bool("v", false, "enable verbose output")
	debug    = flag.Bool("debug", false, "enable debug output")
	initRepo = flag.Bool("init", false, "bootstrap a new binding repository from the template; the module name is taken from the optional argument, or inferred from the current directory name")

	defines multiFlag // -d "NAME1=VALUE1; NAME2=VALUE2", repeatable

	docMode     = flag.Bool("doc", false, "serve the documentation of the Go package in the current directory and open it in a browser")
	docHTTP     = flag.String("http", docserver.DefaultAddr, "listen address for -doc (loopback only unless an explicit non-loopback host is given)")
	docOpen     = flag.Bool("open", true, "launch the default browser for -doc; use -open=false to only print the URL")
	docAll      = flag.Bool("all", false, "include unexported declarations in -doc")
	docAutoExit = flag.Bool("autoexit", true, "for -doc, exit when the auto-opened page is closed; use -autoexit=false to keep serving until Ctrl-C")
)

func init() {
	flag.Var(&defines, "d", "inject constants into the generated package: -d \"NAME1=VALUE1; NAME2=VALUE2\"; may be given more than once")
}

// multiFlag collects a repeatable string flag's values in the order given.
type multiFlag []string

func (m *multiFlag) String() string { return strings.Join(*m, "; ") }

func (m *multiFlag) Set(v string) error {
	*m = append(*m, v)
	return nil
}

const (
	cmdDoc = 1 << iota
	cmdInit
	cmdGen = 0
)

func getCmd() (cmd int) {
	if *docMode {
		cmd |= cmdDoc
	}
	if *initRepo {
		cmd |= cmdInit
	}
	return
}

func main() {
	flag.Parse()

	switch getCmd() {
	case cmdGen:
		runGen()
	case cmdDoc:
		runDoc()
	case cmdInit:
		runInit()
	default:
		help()
		os.Exit(1)
	}
}

func help() {
	fmt.Println("usage: llcppg [-v -debug] [-d \"NAME=VALUE; ...\"] <dest-gopkg-dir> [<src-header-files-and-cfg-dir>]")
	fmt.Println("       llcppg -init [<module-name>]")
	fmt.Println("       llcppg -doc [-http addr] [-open=false] [-all] [-autoexit=false]")
}

func check(err error) {
	if err != nil {
		log.Panicln(err)
	}
}
