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
	"os"

	"github.com/goplus/llcppg/cl"
	"github.com/goplus/llcppg/clang"
	"github.com/goplus/llcppg/tool"
)

var (
	verbose = flag.Bool("v", false, "enable verbose output")
	debug   = flag.Bool("debug", false, "enable debug output")
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("usage: llcppg [-v -debug] <dest-gopkg-dir> [<src-header-files-and-cfg-dir>]")
		return
	}
	destDir := args[0]
	srcDir := "."
	if len(args) >= 2 {
		srcDir = args[1]
	}

	if *debug {
		cl.SetDebug(cl.DbgFlagAll)
		tool.SetDebug(tool.DbgFlagAll)
	} else if *verbose {
		tool.SetDebug(tool.DbgFlagSettings)
	}

	idx := clang.CreateIndex(1, 1)
	defer idx.Dispose()
	err := tool.Gen(destDir, srcDir, idx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
