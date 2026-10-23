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

package cltest

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/goplus/llcppg/cl"
	"github.com/qiniu/x/test/testutil"
)

// -----------------------------------------------------------------------------

// Config represents the lltest configuration.
type Config struct {
	LLGoPackage    string   `json:"LLGoPackage"`
	WrapFileHeader string   `json:"WrapFileHeader"`
	CFlags         string   `json:"CFlags"`
	Files          []string `json:"Files"`

	// NSPrefix specifies the prefix to remove from C/C++ namespace names when generating
	// Go package names (optional).
	NSPrefix []string `json:"NSPrefix"` // C/C++ namespace prefix to remove

	// Type Creator Detection (optional). See https://github.com/xgo-dev/llcppg/issues/960.
	NewCheck []string `json:"NewCheck"`

	// Type Method Detection (optional). See https://github.com/xgo-dev/llcppg/issues/955.
	MethodCheck []string `json:"MethodCheck"`

	// FuncPrefix specifies the prefix to remove from C/C++ global function names when
	// generating Go function names (optional).
	FuncPrefix []string `json:"FuncPrefix"`

	// EnumPrefix specifies the prefix to remove from C/C++ enum value names when generating
	// Go const names (optional).
	EnumPrefix []string `json:"EnumPrefix"`

	// FieldAlias overrides the Go name and type of C/C++ fields (optional).
	// See https://github.com/xgo-dev/llcppg/issues/1019.
	FieldAlias map[string]string `json:"FieldAlias"`

	// TypeDef is an object whose keys are file names and whose values are arrays of type
	// definitions (optional). It defines Go types that should exist in the generated
	// package but is missing from the C/C++ headers.
	// See https://github.com/xgo-dev/llcppg/issues/1079.
	TypeDef map[string][]cl.TypeDef `json:"TypeDef"`

	// FuncDef is an object whose keys are file names and whose values are arrays of
	// function definitions (optional). It declares Go functions and methods that should
	// exist in the generated package, each linked to a C symbol.
	// See https://github.com/goplus/llcppg/issues/1083.
	FuncDef map[string][]cl.FuncDef `json:"FuncDef"`

	TypeAlias  map[string]string `json:"TypeAlias"`
	TypeIgnore []string          `json:"TypeIgnore"`

	// TypePrefix/TypeSuffix specifies the prefix/suffix to remove from C/C++ type names
	// when generating Go type names (optional).
	TypePrefix []string `json:"TypePrefix"`
	TypeSuffix []string `json:"TypeSuffix"`

	// TypeAbbr specifies abbreviated name for Go type names and will be used in function
	// names (optional). See https://github.com/xgo-dev/llcppg/issues/958.
	TypeAbbr []string `json:"TypeAbbr"`

	// Rename specifies a mapping of C/C++ names to Go names (optional). If a name is present
	// in the map, it will be renamed to the corresponding Go name.
	Rename map[string]string `json:"Rename"`

	MacroIgnore []string `json:"MacroIgnore"`

	// MacroFunc maps a function-like C macro name to the name of a reference
	// function declared in the processed headers (optional).
	MacroFunc map[string]string `json:"MacroFunc"`

	// MacroGroup declares groups of related constant macros emitted together as a
	// single Go `const ( ... )` block (optional). It maps the first macro of a
	// group to its last macro (original C macro names, before MacroPrefix
	// stripping or renaming); the first and last macro must live in the same
	// header file. See https://github.com/goplus/llcppg/issues/1064.
	MacroGroup map[string]string `json:"MacroGroup"`

	// MacroPrefix specifies the prefix to remove from C/C++ macro names when
	// generating Go const names (optional).
	MacroPrefix []string `json:"MacroPrefix"`

	// Define specifies user constants to inject into the generated package (optional).
	Define map[string]any `json:"Define"`

	// convert capitalized, underscore-joined names to camel case, except names whose final
	// segment is all-uppercase (acronyms/macros), which are kept as-is.
	ForceCamelCase bool `json:"ForceCamelCase"`

	LoadLibcPubFile bool `json:"LoadLibcPubFile"`
}

// LoadConf loads the lltest configuration.
func LoadConf(filename string) (conf Config, err error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &conf)
	return
}

// -----------------------------------------------------------------------------

// TestFromDir runs the provided test function for each subdirectory of the specified
// relative directory. If sel is not empty, only the subdirectory matching sel will be
// tested. If relDir is a single-level directory, all its immediate subdirectories will
// be considered. If relDir is not a single-level directory, only the specified directory
// itself will be tested.
func TestFromDir(t *testing.T, sel, relDir string, testFunc func(t *testing.T, pkgDir string)) {
	testutil.FromDir(t, sel, relDir, testFunc)
}

// -----------------------------------------------------------------------------

var langs = [...]string{
	cl.LanguageC:   "c",
	cl.LanguageCXX: "c++",
}

// LanguageOf returns the string representation of the specified programming language.
func LanguageOf(lang cl.Language) string {
	return langs[lang]
}

// -----------------------------------------------------------------------------

// MockNameLookup is a mock implementation of the NameLookup function. It returns a
// fixed archive path and true for any input.
func MockNameLookup(manglingName string) (archivePath string, ok bool) {
	return "libfoo.a", true
}

// -----------------------------------------------------------------------------
