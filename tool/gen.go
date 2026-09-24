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
	"fmt"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/goplus/gogen/packages"
	"github.com/goplus/gogen/packages/cache"
	"github.com/goplus/llcppg/cl"
	"github.com/goplus/llcppg/clang"
	"github.com/goplus/llcppg/tool/listth"
	"github.com/goplus/mod"
	"github.com/goplus/mod/xgomod"

	lc "github.com/goplus/llcppg/lib/clang"
)

const (
	DbgFlagLoadSource = 1 << iota
	DbgFlagLoadDeps
	DbgFlagAll = DbgFlagLoadSource | DbgFlagLoadDeps
)

var (
	debugLoadSource bool
	debugLoadDeps   bool
)

func SetDebug(flags int) {
	debugLoadSource = (flags & DbgFlagLoadSource) != 0
	debugLoadDeps = (flags & DbgFlagLoadDeps) != 0
}

// -----------------------------------------------------------------------------

type Config struct {
	Name           string            `json:"Name"`
	LLGoPackage    string            `json:"LLGoPackage"`
	WrapFileHeader string            `json:"WrapFileHeader"`
	CFlags         string            `json:"CFlags"`
	Language       string            `json:"Language"`   // c, c++, etc.
	Dir            string            `json:"Dir"`        // dir or dir/... (recursive)
	Deps           []string          `json:"Deps"`       // dependencies (package paths)
	Class          []string          `json:"Class"`      // typedef names to be treated as classes
	TypePrefix     []string          `json:"TypePrefix"` // type prefix to remove
	FuncPrefix     []string          `json:"FuncPrefix"` // global function prefix to remove
	Rename         map[string]string `json:"Rename"`     // renaming of C/C++ names to Go names
}

// LoadConf loads the llcppg configuration.
func LoadConf(filename string) (conf Config, err error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &conf)
	return
}

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

const includeSuffix = string(os.PathSeparator) + "include"

// topHeaders lists the top-level header files according to the configuration. If the
// Dir field ends with "/...", it will recursively list all header files in the directory
// and its subdirectories. includeDirs[0] will be set to the include directory for the
// package.
func topHeaders(dir, workDir string, includeDirs []string) (headerFiles []string, err error) {
	recursive := strings.HasSuffix(dir, "/...")
	if recursive {
		dir = dir[:len(dir)-4]
	}
	if !filepath.IsAbs(dir) {
		dir = filepath.Join(workDir, dir)
	}
	incDir := dir
	if pos := strings.LastIndex(incDir, includeSuffix); pos >= 0 {
		incDir = incDir[:pos+len(includeSuffix)]
	}
	includeDirs[0] = incDir
	if debugLoadSource {
		log.Println("==> includeDirs:", includeDirs)
	}
	return listth.TopHeaders(dir, recursive, false, includeDirs)
}

// NewPackage loads the source files and converts them into a Go package according to the
// configuration.
func (cfg *Config) NewPackage(pkgPath, workDir, stdlibDir string, index clang.Index) (ret cl.Package, lang cl.Language, err error) {
	lang, ok := cfg.Lang()
	if !ok {
		err = fmt.Errorf("invalid language: %q", cfg.Language)
		return
	}

	if debugLoadSource {
		log.Println("==> workDir:", workDir)
	}

	mod, err := LoadModuleFrom(workDir)
	if err != nil {
		return
	}

	fset := token.NewFileSet()
	imp := packages.NewImporter(fset, workDir)

	deps := cfg.Deps
	if len(deps) > 0 {
		c := cache.New(pkgHash)
		c.Prepare(workDir, deps...)
		imp.SetCache(c)
	}

	includeDirs, pkgPaths := mod.includeDirs(imp, deps, 2)
	includeDirs[1] = stdlibDir
	pkgPaths[0] = pkgPath
	pkgPaths[1] = "github.com/goplus/lib/c"

	// includeDirs[0] is set by topHeaders
	topHeaders, err := topHeaders(cfg.Dir, workDir, includeDirs)
	if err != nil {
		return
	}
	/* if debugLoadSource {
		log.Println("==> topHeaders:", topHeaders)
	} */

	files := ParseSources(index, topHeaders, includeDirs, cfg.Language)
	defer DisposeSources(files)

	if debugLoadSource {
		dumpSources(topHeaders, files)
	}

	ret, err = cl.NewPackage(pkgPath, cfg.Name, files, &cl.Config{
		Fset:           fset,
		Importer:       imp,
		LLGoPackage:    cfg.LLGoPackage,
		Language:       lang,
		CFlags:         cfg.CFlags,
		WrapFileHeader: cfg.WrapFileHeader,
		NameLookup:     nil,
		PubFileLookup:  mod.PubFileLookup,
		PackageOf: func(headerFile string) (pkgPath string, ok bool) {
			for i, includeDir := range includeDirs {
				if strings.HasPrefix(headerFile, includeDir) {
					return pkgPaths[i], true
				}
			}
			return
		},
	})
	return
}

func pkgHash(pkgPath string, self bool) string {
	// don't need calc hash since we only load all packages once.
	return cache.HashSkip
}

// -----------------------------------------------------------------------------

// Module represents a llcppg module.
type Module struct {
	mod *xgomod.Module
}

// LoadModuleFrom loads a llcppg module from the specified directory.
func LoadModuleFrom(dir string) (ret Module, err error) {
	_, gomod, err := mod.FindGoMod(dir)
	if err != nil {
		return
	}
	m, err := xgomod.LoadFrom(gomod, "")
	return Module{mod: m}, err
}

// PubFileLookup looks up the public file for the given package.
func (p Module) PubFileLookup(pkgPath string) (pubFile string, ok bool) {
	pkg, err := p.mod.Lookup(pkgPath)
	if err == nil {
		pubFile, ok = filepath.Join(pkg.Dir, "llcppg.pub"), true
	}
	return
}

// includeDirs returns the include directories and package paths for the given dependencies.
// The reserved parameter specifies the number of reserved slots in the returned slices.
func (p Module) includeDirs(imp *packages.Importer, deps []string, reserved int) (incDirs, pkgPaths []string) {
	incDirs = make([]string, reserved, len(deps)+reserved)
	pkgPaths = make([]string, reserved, len(deps)+reserved)
	for _, dep := range deps {
		pkgTypes, err := imp.Import(dep)
		if err == nil {
			if o := pkgTypes.Scope().Lookup("LLGoFiles"); o != nil {
				pkg, err := p.mod.Lookup(dep)
				if err == nil {
					incDirs = append(incDirs, filepath.Join(pkg.Dir, "_wrap/include"))
					pkgPaths = append(pkgPaths, dep)
				}
			}
		}
	}
	return
}

// -----------------------------------------------------------------------------

// ParseSources parses the given source files and returns the translation units corresponding
// to those files.
func ParseSources(index clang.Index, headerFiles, includeDirs []string, lang string) []cl.Source {
	options := lc.DefaultDiagnosticDisplayOptions()
	n := len(includeDirs)
	flags := make([]string, n+2)
	for i, dir := range includeDirs {
		flags[i] = "-I" + dir
	}
	flags[n] = "-x"
	flags[n+1] = lang
	files := make([]cl.Source, 0, len(headerFiles))
	for tu := range index.ParseTranslationUnits(clang.DetailedPreprocessingRecord, headerFiles, flags...) {
		files = append(files, tu)
		tu.VisitDiagnostics(func(diag clang.Diagnostic) {
			fmt.Fprintln(os.Stderr, diag.Format(options))
		})
	}
	return files
}

// DisposeSources disposes the given translation units.
func DisposeSources(sources []cl.Source) {
	for _, tu := range sources {
		tu.Dispose()
	}
}

// -----------------------------------------------------------------------------
