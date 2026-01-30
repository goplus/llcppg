/*
 * Copyright (c) 2025 The XGo Authors (xgo.dev). All rights reserved.
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
	"go/token"

	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llgo/xtool/clang/preprocessor"
)

// Config represents the configuration for parsing C/C++ source files.
type Config struct {
	Compiler    string // default: clang++
	PPFlag      string // default: -E
	BaseDir     string // base of include searching directory, should be absolute path
	IncludeDirs []string
	Defines     []string
	Flags       []string

	Mode Mode // parsing mode
}

// ParseFile parses a C/C++ source file and returns the corresponding AST.
// Allow fset to be nil, in which case a new FileSet will be created.
func ParseFile(fset *token.FileSet, srcFile, interFile string, conf *Config) (f *ast.File, err error) {
	mode := ParseAllComments
	var ppconf *preprocessor.Config
	if conf == nil {
		conf = new(Config)
	}
	cc := conf.Compiler
	if cc == "" {
		cc = "clang++"
	}
	ppconf = &preprocessor.Config{
		Compiler:    cc,
		PPFlag:      conf.PPFlag,
		BaseDir:     conf.BaseDir,
		IncludeDirs: conf.IncludeDirs,
		Defines:     conf.Defines,
		Flags:       conf.Flags,
	}
	mode = conf.Mode
	if interFile == "" {
		interFile = srcFile + ".i"
	}
	err = preprocessor.Do(srcFile, interFile, ppconf)
	if err != nil {
		return
	}
	return ParseIntermediateFile(fset, interFile, mode)
}
