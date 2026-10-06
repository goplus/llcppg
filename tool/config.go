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

type Config struct {
	Name             string            `json:"Name"`             // required, sub package name should start with '/'
	Language         string            `json:"Language"`         // c, c++, etc. required
	Dir              string            `json:"Dir"`              // dir or dir/... (recursive), required
	Files            []string          `json:"Files"`            // selected header files relative to Dir; overrides auto-discovery, optional.
	Stdlib           string            `json:"Stdlib"`           // C stdlib include dir, optional
	LLGoPackage      string            `json:"LLGoPackage"`      // optional
	CFlags           string            `json:"CFlags"`           // optional
	Deps             []string          `json:"Deps"`             // dependencies (module paths), optional
	Class            []string          `json:"Class"`            // typedef names to be treated as classes
	NonClass         []string          `json:"NonClass"`         // typedef names to be treated as non-classes
	NSPrefix         []string          `json:"NSPrefix"`         // C/C++ namespace prefix to remove
	MethodCheck      []string          `json:"MethodCheck"`      // C/C++ method check list
	FuncPrefix       []string          `json:"FuncPrefix"`       // C/C++ function name prefix to remove
	VarPrefix        []string          `json:"VarPrefix"`        // C/C++ variable name prefix to remove
	EnumPrefix       []string          `json:"EnumPrefix"`       // C/C++ enum value prefix to remove
	MacroPrefix      []string          `json:"MacroPrefix"`      // C/C++ macro name prefix to remove
	TypePrefix       []string          `json:"TypePrefix"`       // C/C++ type name prefix to remove
	TypeSuffix       []string          `json:"TypeSuffix"`       // C/C++ type name suffix to remove
	TypeAbbrSuffix   []string          `json:"TypeAbbrSuffix"`   // Go type abbr suffix to remove, only valid for types that are not present in TypeAbbr
	TypeAbbr         map[string]any    `json:"TypeAbbr"`         // Go type name to its abbr(s), used in function names
	Rename           map[string]string `json:"Rename"`           // renaming of C/C++ names to Go names
	TypeAlias        map[string]string `json:"TypeAlias"`        // C/C++ type name to a Go type name in pkgPath.Name format (pkgPath can be empty if Name is in current package), optional
	TypeIgnore       []string          `json:"TypeIgnore"`       // C/C++ type names to ignore
	FuncIgnore       []string          `json:"FuncIgnore"`       // C/C++ function names to ignore
	MacroIgnore      []string          `json:"MacroIgnore"`      // C/C++ macro names to ignore
	NSIgnore         []string          `json:"NSIgnore"`         // C/C++ namespaces to ignore
	Pkgs             []string          `json:"Pkgs"`             // sub-packages to generate, optional
	FailFast         int               `json:"FailFast"`         // exit on first N errors, optional
	ForceCamelCase   bool              `json:"ForceCamelCase"`   // convert capitalized, underscore-joined names to camel case, except names whose final segment is all-uppercase (acronyms/macros), which are kept as-is.
	IgnoreInline     bool              `json:"IgnoreInline"`     // quietly ignore inline functions
	NoManglingIgnore bool              `json:"NoManglingIgnore"` // quietly ignore functions with no mangled symbol
	GroupSubdir      bool              `json:"GroupSubdir"`      // treats sub-directory files as a single file. Deprecated: use GroupSubdirBy instead.
	GroupSubdirBy    string            `json:"GroupSubdirBy"`    // criterion to group sub-directory files by, e.g., "dir" or "fname". `GroupSubdir = true` is equivalent to `GroupSubdirBy = "dir"`.
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
