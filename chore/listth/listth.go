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
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/goplus/llcppg/tool/listth"
)

func check(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

// -----------------------------------------------------------------------------

func parseArgs(args []string) (string, []string) {
	if len(args) > 1 {
		var includeDirs []string
		for i, arg := range args[1:] {
			if strings.HasPrefix(arg, "-I") {
				incDir, err := filepath.Abs(arg[2:])
				check(err)
				includeDirs = append(includeDirs, incDir)
			} else {
				if len(args)-2 != i {
					break
				}
				headerDir, err := filepath.Abs(arg)
				check(err)
				return headerDir, includeDirs
			}
		}
	}
	return "", nil
}

// -----------------------------------------------------------------------------

func main() {
	headerDir, includeDirs := parseArgs(os.Args)
	if headerDir == "" {
		fmt.Println("usage: listth [-I<include-dir> ...] <header-dir>")
		return
	}

	files, err := listth.TopHeaders(headerDir, false, true, includeDirs)
	check(err)

	for _, file := range files {
		fmt.Println(file)
	}
}

// -----------------------------------------------------------------------------
