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

package cstdlib

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/goplus/llcppg/xtool/env/clang"
)

// -----------------------------------------------------------------------------

var stdlibDirs []string

// Dirs returns the list of directories where the C standard library headers
// are located.
func Dirs() []string {
	if stdlibDirs == nil {
		stdlibDirs = findStdlibDirs()
	}
	return stdlibDirs
}

func findStdlibDirs() []string {
	var stderr bytes.Buffer
	cmd := exec.Command(clang.Path(), "-x", "c", "-E", "-v", "-")
	cmd.Stdin = bytes.NewReader(nil)
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	return parseStdlibDirs(stderr.String())
}

const (
	startMarker = "#include <...> search starts here:\n "
	endMarker   = "\nEnd of search list."
)

func parseStdlibDirs(output string) []string {
	start := strings.Index(output, startMarker)
	if start < 0 {
		panic("failed to find start marker in clang output")
	}
	output = output[start+len(startMarker):]
	end := strings.Index(output, endMarker)
	if end < 0 {
		panic("failed to find end marker in clang output")
	}
	dirs := strings.SplitN(output[:end], "\n ", 3)
	if len(dirs) > 2 {
		dirs = dirs[:2] // only keep the first two directories
	}
	return dirs
}

// -----------------------------------------------------------------------------
