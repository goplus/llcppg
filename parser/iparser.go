/*
 * Copyright (c) 2025 The GoPlus Authors (goplus.org). All rights reserved.
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

package parser

import (
	"bytes"
	"go/token"
	"os"
	"os/exec"
	"strconv"

	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/internal/unmarshal"
)

// Mode represents the parsing mode.
type Mode int

const (
	ParseC           Mode = 1 << iota // default parser with C++
	ParseAllComments                  // clang -fparse-all-comments
)

// ParseIntermediateFile parses an intermediate file (*.i) and returns the corresponding AST.
// Allow fset to be nil, in which case a new FileSet will be created.
func ParseIntermediateFile(fset *token.FileSet, filename string, mode Mode) (f *ast.File, err error) {
	var stdout bytes.Buffer
	cmd := exec.Command("llclang", "parseIntermediateFile", filename, strconv.Itoa(int(mode)))
	cmd.Stderr = os.Stderr
	cmd.Stdout = &stdout
	err = cmd.Run()
	if err != nil {
		return
	}
	file, err := unmarshal.File(stdout.Bytes())
	if err != nil {
		return nil, err
	}
	return file, nil
}
