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
	"maps"
	"os"
	"slices"
	"strings"

	"github.com/goplus/llcppg/tool/pputil"
)

type none struct{}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(`usage: ppresolve <header-dir>[/...]`)
		return
	}
	dir := os.Args[1]
	recursive := strings.HasSuffix(dir, "/...")
	if recursive {
		dir = dir[:len(dir)-4]
	}

	includeDirs := []string{"."}
	unresolvedIncludes := make(map[string]none)
	for file, err := range pputil.ListFiles(dir, recursive) {
		if err != nil {
			fmt.Fprintln(os.Stderr, "[ERROR]", err)
			return
		}
		if pputil.IsHeaderFile(file.Path) {
			includeFiles, err := pputil.ListIncludes(file.Path, includeDirs)
			if err != nil {
				fmt.Fprintln(os.Stderr, "[ERROR]", err)
				return
			}
			for includeFile, ok := range includeFiles {
				if !ok {
					unresolvedIncludes[includeFile] = none{}
				}
			}
		}
	}
	for _, includeFile := range slices.Sorted(maps.Keys(unresolvedIncludes)) {
		fmt.Println(includeFile)
	}
}
