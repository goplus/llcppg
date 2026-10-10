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
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/goplus/gogen"
	"github.com/goplus/llcppg/cl"
	"github.com/goplus/llcppg/clang"
	"github.com/goplus/llcppg/internal/docserver"
	"github.com/goplus/llcppg/tool"
	"github.com/goplus/llcppg/tool/initpkg"
)

var (
	verbose  = flag.Bool("v", false, "enable verbose output")
	debug    = flag.Bool("debug", false, "enable debug output")
	initRepo = flag.Bool("init", false, "bootstrap a new binding repository from the template; the module name is taken from the optional argument, or inferred from the current directory name")

	docMode = flag.Bool("doc", false, "serve the documentation of the Go package in the current directory and open it in a browser")
	docHTTP = flag.String("http", docserver.DefaultAddr, "listen address for -doc (loopback only unless an explicit non-loopback host is given)")
	docOpen = flag.Bool("open", true, "launch the default browser for -doc; use -open=false to only print the URL")
	docAll  = flag.Bool("all", false, "include unexported declarations in -doc")
)

func main() {
	flag.Parse()

	// -doc is a mutually exclusive mode: it neither bootstraps a repository
	// nor converts headers, so combining it with -init (or passing conversion
	// arguments) is a usage error.
	if *docMode {
		if *initRepo {
			fmt.Fprintln(os.Stderr, "llcppg: -doc cannot be combined with -init")
			os.Exit(1)
		}
		if args := flag.Args(); len(args) > 0 {
			fmt.Fprintln(os.Stderr, "llcppg: -doc takes no arguments; run it in the directory of the package to document")
			os.Exit(1)
		}
		runDoc()
		return
	}

	if *initRepo {
		// The module name is optional: `llcppg -init <module-name>` uses the
		// given name, `llcppg -init` infers it from the current directory.
		module := ""
		if args := flag.Args(); len(args) > 0 {
			module = args[0]
		}
		dir, err := cacheDir()
		check(err)
		if err := initpkg.Init(".", module, &initpkg.Config{CacheDir: dir}); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("usage: llcppg [-v -debug] <dest-gopkg-dir> [<src-header-files-and-cfg-dir>]")
		fmt.Println("       llcppg -init [<module-name>]")
		fmt.Println("       llcppg -doc [-http addr] [-open=false] [-all]")
		return
	}
	destDir := args[0]
	srcDir := "."
	if len(args) >= 2 {
		srcDir = args[1]
	}

	log.SetFlags(0)
	if *debug {
		cl.SetDebug(cl.DbgFlagAll)
		tool.SetDebug(tool.DbgFlagAll)
		gogen.SetDebug(gogen.DbgFlagAll)
	} else if *verbose {
		cl.SetDebug(cl.DbgFlagMajorProc)
		tool.SetDebug(tool.DbgFlagSettings)
	}

	idx := clang.CreateIndex(1, 1)
	defer idx.Dispose()
	err := Gen(destDir, srcDir, idx, *verbose || *debug)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// runDoc serves the documentation of the Go package in the current directory
// (`llcppg -doc`). It blocks until interrupted (SIGINT), then shuts down
// gracefully. Usage errors (no Go files, more than one package, a bad listen
// address) are reported and exit non-zero.
func runDoc() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	err := docserver.Run(ctx, ".", docserver.Options{
		Addr:        *docHTTP,
		OpenBrowser: *docOpen,
		AllDecls:    *docAll,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// CacheDir returns the directory llcppg uses to cache downloadable resources
// such as the project template. The LLCPPG_CACHE environment variable overrides
// the default, which is <os.UserCacheDir>/llcppg.
func cacheDir() (string, error) {
	base, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("locate user cache dir: %w", err)
	}
	return filepath.Join(base, "llcppg"), nil
}

func check(err error) {
	if err != nil {
		log.Panicln(err)
	}
}
