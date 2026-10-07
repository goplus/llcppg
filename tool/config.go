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

package tool

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/goplus/llcppg/cl"
	"github.com/goplus/llcppg/xtool/env/cstdlib"
)

// -----------------------------------------------------------------------------

// Config represents the configuration for the llcppg tool.
type Config struct {
	// Package Name (required). Sub-package name should start with '/'.
	Name string `json:"Name"`

	// Language (required). c, c++, etc.
	Language string `json:"Language"`

	// dir or dir/... (recursive), required
	Dir string `json:"Dir"`

	// selected header files relative to Dir; overrides auto-discovery, optional.
	Files []string `json:"Files"`

	// C stdlib include dir, optional
	Stdlib string `json:"Stdlib"`

	// LLGoPackage specifies the value of the LLGoPackage constant in the generated
	// Go package (optional).
	LLGoPackage string `json:"LLGoPackage"`

	// CFlags specifies the compiler flags to be used when compiling the wrapper file.
	// If not specified, llcppg will skip wrapping inline functions/methods.
	CFlags string `json:"CFlags"`

	// dependencies (module paths), optional
	Deps []string `json:"Deps"`

	// Class specifies a list of C/C++ typedef names to be treated as classes (optional).
	Class []string `json:"Class"`

	// 1) for typedef: NonClass specifies a list of C/C++ typedef names to be treated
	//    as non-classes.
	// 2) for global function to method: NonClass specifies a list of Go type names to
	//    be treated as non-classes.
	NonClass []string `json:"NonClass"`

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

	// VarPrefix specifies the prefix to remove from C/C++ global variable names when
	// generating Go variable names (optional).
	VarPrefix []string `json:"VarPrefix"`

	// EnumPrefix specifies the prefix to remove from C/C++ enum value names when generating
	// Go const names (optional).
	EnumPrefix []string `json:"EnumPrefix"`

	// MacroPrefix specifies the prefix to remove from C/C++ macro names when generating
	// Go const names (optional).
	MacroPrefix []string `json:"MacroPrefix"`

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

	// TypeAlias specifies a mapping of C/C++ type names to a Go type name in pkgPath.Name
	// format (pkgPath can be empty if Name is in current package), optional.
	TypeAlias map[string]string `json:"TypeAlias"`

	// TypeIgnore specifies a list of C/C++ type names to be ignored (optional).
	TypeIgnore []string `json:"TypeIgnore"`

	// FuncIgnore specifies a list of C/C++ function names to be ignored (optional).
	FuncIgnore []string `json:"FuncIgnore"`

	// MacroIgnore specifies a list of C/C++ macro names to be ignored (optional).
	MacroIgnore []string `json:"MacroIgnore"`

	// NSIgnore specifies a list of C/C++ namespaces to be ignored (optional).
	NSIgnore []string `json:"NSIgnore"`

	// Sub-packages to generate, optional
	Pkgs []string `json:"Pkgs"`

	// exit on first N errors, optional
	FailFast int `json:"FailFast"`

	// convert capitalized, underscore-joined names to camel case, except names whose final
	// segment is all-uppercase (acronyms/macros), which are kept as-is.
	ForceCamelCase bool `json:"ForceCamelCase"`

	// quietly ignore inline functions
	IgnoreInline bool `json:"IgnoreInline"`

	// quietly ignore functions with no mangled symbol
	NoManglingIgnore bool `json:"NoManglingIgnore"`

	// treats sub-directory files as a single file.
	// Deprecated: use GroupSubdirBy instead.
	GroupSubdir bool `json:"GroupSubdir"`

	// criterion to group sub-directory files by, e.g., "dir" or "fname". `GroupSubdir = true`
	// is equivalent to `GroupSubdirBy = "dir"`.
	GroupSubdirBy string `json:"GroupSubdirBy"`
}

// -----------------------------------------------------------------------------

// LoadConf loads the llcppg configuration.
func LoadConf(filename string) (cfg Config, err error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &cfg)
	return
}

// -----------------------------------------------------------------------------

// Apply applies the parent configuration to the current configuration.
func (cfg *Config) Apply(parent *Config) {
	cfg.Name = parent.Name + cfg.Name
	if cfg.Language == "" {
		cfg.Language = parent.Language
	}
	if cfg.Stdlib == "" {
		cfg.Stdlib = parent.Stdlib
	}
	if cfg.LLGoPackage == "" {
		cfg.LLGoPackage = parent.LLGoPackage
	}
	if cfg.CFlags == "" {
		cfg.CFlags = parent.CFlags
	}
	if cfg.Deps == nil {
		cfg.Deps = parent.Deps
	}
}

// -----------------------------------------------------------------------------

// Lang returns the language of the configuration.
func (cfg *Config) Lang() (lang cl.Language, ok bool) {
	switch cfg.Language {
	case "c":
		return cl.LanguageC, true
	case "c++":
		return cl.LanguageCXX, true
	default:
		return 0, false
	}
}

// -----------------------------------------------------------------------------

// StdlibDirs returns the standard library include directories.
func (cfg *Config) StdlibDirs(workDir string) []string {
	var stdlibDirs []string
	if stdlibDir := cfg.Stdlib; stdlibDir != "" {
		if !filepath.IsAbs(stdlibDir) {
			stdlibDir = filepath.Join(workDir, stdlibDir)
		}
		stdlibDirs = []string{stdlibDir}
	} else {
		stdlibDirs = cstdlib.Dirs()
	}
	return stdlibDirs
}

// -----------------------------------------------------------------------------
