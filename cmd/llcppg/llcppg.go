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
	"fmt"
	"os"

	"github.com/goplus/llcppg/clang"
	"github.com/goplus/llcppg/tool"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: llcppg <dest-gopkg-dir> [<src-header-files-and-cfg-dir>]")
		return
	}
	destDir := os.Args[1]
	srcDir := "."
	if len(os.Args) >= 3 {
		srcDir = os.Args[2]
	}

	idx := clang.CreateIndex(0, 0)
	defer idx.Dispose()
	err := tool.Gen(destDir, srcDir, idx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
