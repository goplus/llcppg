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
	"sort"
	"strings"

	"github.com/goplus/llcppg/tool/pputil"
)

func check(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

// -----------------------------------------------------------------------------

type headerInfo struct {
	deps      []string
	hasExtern bool
}

const includeSuffix = string(os.PathSeparator) + "include"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: llfilter <header-dir> <hasExtern>")
		return
	}

	headerDir, err := filepath.Abs(os.Args[1])
	check(err)

	fHasExtern := false
	if len(os.Args) > 2 {
		arg := os.Args[2]
		fHasExtern = arg != "false" && arg != "0"
	}

	incDir := headerDir
	if pos := strings.LastIndex(incDir, includeSuffix); pos >= 0 {
		incDir = incDir[:pos+len(includeSuffix)]
	}
	includeDirs := []string{incDir}

	ret := make(map[string]headerInfo)
	headerDirPrefix := headerDir + string(os.PathSeparator)
	for file, err := range pputil.ListFilterFiles(headerDir, true, pputil.FilterPublicHeaderFile) {
		check(err)

		includeFiles, err := pputil.ListIncludes(file.Path, includeDirs)
		check(err)

		var deps []string
		var hasExtern bool
		for includeFile, found := range includeFiles {
			if found {
				if strings.HasPrefix(includeFile, headerDirPrefix) {
					deps = append(deps, includeFile)
				} else {
					hasExtern = true
				}
			}
		}
		ret[file.Path] = headerInfo{
			deps:      deps,
			hasExtern: hasExtern,
		}
	}

	var selFiles []string
	for file := range ret {
		if fHasExtern == hasExtern(ret, file) {
			selFiles = append(selFiles, file)
		}
	}
	sort.Strings(selFiles)
	sep := ","
	for i, file := range selFiles {
		if i == len(selFiles)-1 {
			sep = ""
		}
		fmt.Printf("\t\t%q%s\n", file[len(headerDirPrefix):], sep)
	}
}

func hasExtern(ret map[string]headerInfo, file string) bool {
	info, ok := ret[file]
	if !ok {
		panic("file not found: " + file)
	}
	if info.hasExtern {
		return true
	}
	for _, dep := range info.deps {
		if hasExtern(ret, dep) {
			return true
		}
	}
	return false
}

// -----------------------------------------------------------------------------
