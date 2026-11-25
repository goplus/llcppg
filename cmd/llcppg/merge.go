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
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/goplus/llcppg/tool"
	"github.com/goplus/llcppg/tool/merge"
	"github.com/goplus/llcppg/tool/pputil"
)

func runMerge() {
	args := flag.Args()
	if len(args) < 1 || len(args) > 2 {
		help()
		os.Exit(1)
	}
	destDir := args[0]
	srcDir := "."
	if len(args) == 2 {
		srcDir = args[1]
	}
	if err := mergePkg(destDir, srcDir); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func mergePkg(destDir, srcDir string) error {
	platforms, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}

	fset := token.NewFileSet()
	pkgs := make([]*merge.PlatformPackage, 0, len(platforms))
	for _, platform := range platforms {
		if !platform.IsDir() {
			continue
		}
		tags := strings.Split(platform.Name(), "-")
		if len(tags) != 2 || tags[0] == "" || tags[1] == "" {
			continue
		}

		platformDir := filepath.Join(srcDir, platform.Name())
		var names []string
		for file, err := range pputil.ListFiles(platformDir, false) {
			if err != nil {
				return err
			}
			if strings.HasSuffix(file.Name(), ".go") {
				names = append(names, file.Name())
			}
		}
		if len(names) == 0 {
			return fmt.Errorf("no Go inputs in %s", platformDir)
		}
		files, err := tool.LoadGoFiles(fset, platformDir, names, parser.ParseComments)
		if err != nil {
			return err
		}
		pkgs = append(pkgs, &merge.PlatformPackage{Tags: tags, Files: files})
	}
	if len(pkgs) == 0 {
		return fmt.Errorf("no platform packages in %s", srcDir)
	}

	merged, err := merge.Merge(pkgs)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}
	for _, name := range slices.Sorted(maps.Keys(merged)) {
		if err := writeMergedFile(filepath.Join(destDir, name), fset, merged[name]); err != nil {
			return err
		}
	}
	return nil
}

func writeMergedFile(path string, fset *token.FileSet, input *ast.File) error {
	var source bytes.Buffer
	if input.Doc != nil {
		for _, comment := range input.Doc.List {
			fmt.Fprintln(&source, comment.Text)
		}
		source.WriteByte('\n')
	}
	fmt.Fprintln(&source, "package", input.Name.Name)

	// Print each declaration with its source positions so //go:linkname Dlopen C.dlopen stays with Dlopen.
	for _, decl := range input.Decls {
		source.WriteByte('\n')
		if err := format.Node(&source, fset, decl); err != nil {
			return err
		}
		source.WriteByte('\n')
	}
	formatted, err := format.Source(source.Bytes())
	if err != nil {
		return err
	}
	return os.WriteFile(path, formatted, 0644)
}
