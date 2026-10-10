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

package docserver

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/doc"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

// pkg is the loaded documentation for the package in a directory, together with
// the pieces the renderer needs: the file set used for printing declarations,
// any non-fatal parse errors, and the resolved import path.
type pkg struct {
	Doc        *doc.Package
	FileSet    *token.FileSet
	ImportPath string
	Dir        string // absolute directory the package was read from
	// ParseErrors holds per-file parse errors. They are non-fatal: the page
	// renders whatever parsed and lists these at the top, because a
	// half-broken generation result is exactly when the author wants to look.
	ParseErrors []string
}

// load reads the Go package in dir and builds its documentation. allDecls
// selects doc.AllDecls so unexported declarations are included.
//
// It returns a usage error (and no pkg) when dir contains no Go files for the
// current build context or more than one package. Parse errors in otherwise
// present files are not fatal: they are collected on the returned pkg and the
// documentation is built from the files that parsed.
func load(dir string, allDecls bool) (*pkg, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	// Use go/build to select the files that belong to the package for the
	// current GOOS/GOARCH, honouring file-name suffixes and //go:build
	// constraints that generated bindings commonly use. Test files are
	// excluded (we never serve _test packages).
	ctxt := build.Default
	bpkg, err := ctxt.ImportDir(absDir, 0)
	if err != nil {
		if _, ok := err.(*build.NoGoError); ok {
			return nil, fmt.Errorf("llcppg: no Go files in %s for %s/%s", absDir, ctxt.GOOS, ctxt.GOARCH)
		}
		if m, ok := err.(*build.MultiplePackageError); ok {
			return nil, fmt.Errorf("llcppg: directory %s contains more than one package (%s); run llcppg -doc in the directory of a single package", absDir, strings.Join(uniqueStrings(m.Packages), ", "))
		}
		return nil, fmt.Errorf("llcppg: cannot read package in %s: %w", absDir, err)
	}

	// goFiles are the non-test .go files build selected, in directory order.
	goFiles := make([]string, 0, len(bpkg.GoFiles)+len(bpkg.CgoFiles))
	goFiles = append(goFiles, bpkg.GoFiles...)
	goFiles = append(goFiles, bpkg.CgoFiles...)
	sort.Strings(goFiles)
	if len(goFiles) == 0 {
		return nil, fmt.Errorf("llcppg: no Go files in %s for %s/%s", absDir, ctxt.GOOS, ctxt.GOARCH)
	}

	fset := token.NewFileSet()
	files := make([]*ast.File, 0, len(goFiles))
	var parseErrors []string
	for _, name := range goFiles {
		path := filepath.Join(absDir, name)
		f, perr := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if perr != nil {
			parseErrors = append(parseErrors, perr.Error())
		}
		// ParseFile returns a partial AST even on error; render what we got.
		if f != nil {
			files = append(files, f)
		}
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("llcppg: could not parse any Go file in %s:\n%s", absDir, strings.Join(parseErrors, "\n"))
	}

	importPath := resolveImportPath(absDir)

	mode := doc.Mode(0)
	if allDecls {
		mode |= doc.AllDecls
	}
	dpkg, err := doc.NewFromFiles(fset, files, importPath, mode)
	if err != nil {
		return nil, fmt.Errorf("llcppg: cannot build documentation for %s: %w", absDir, err)
	}

	return &pkg{
		Doc:         dpkg,
		FileSet:     fset,
		ImportPath:  importPath,
		Dir:         absDir,
		ParseErrors: parseErrors,
	}, nil
}

// resolveImportPath determines the import path to show and to pass to
// doc.NewFromFiles. If a go.mod is found in absDir or a parent, the import path
// is the module path joined with the relative directory. Otherwise it falls
// back to the directory's base name. Only the module line of go.mod is needed,
// so no third-party parser is required.
func resolveImportPath(absDir string) string {
	dir := absDir
	for {
		mod := filepath.Join(dir, "go.mod")
		if data, err := os.ReadFile(mod); err == nil {
			if mp := modulePath(data); mp != "" {
				rel, rerr := filepath.Rel(dir, absDir)
				if rerr != nil || rel == "." || rel == "" {
					return mp
				}
				return path.Join(mp, filepath.ToSlash(rel))
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return filepath.Base(absDir)
}

// modulePath extracts the module path from go.mod content, reading only the
// `module` line. It returns "" when no module directive is present.
func modulePath(data []byte) string {
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "//") {
			continue
		}
		rest, ok := strings.CutPrefix(line, "module")
		if !ok {
			continue
		}
		rest = strings.TrimSpace(rest)
		if rest == "" {
			continue
		}
		// A module path may be quoted in go.mod.
		rest = strings.Trim(rest, "`\"")
		// Drop a trailing inline comment.
		if i := strings.Index(rest, "//"); i >= 0 {
			rest = strings.TrimSpace(rest[:i])
		}
		if rest != "" {
			return rest
		}
	}
	return ""
}

// uniqueStrings returns the distinct elements of s, preserving first-seen order.
func uniqueStrings(s []string) []string {
	seen := make(map[string]bool, len(s))
	out := s[:0:0]
	for _, v := range s {
		if !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}
