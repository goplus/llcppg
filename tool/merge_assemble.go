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
	"go/ast"
	"go/token"
	"sort"
	"strconv"
	"strings"
)

// -----------------------------------------------------------------------------
// Step 4 + 6: decide shared vs per-platform and assemble the output.

// entry is one declaration placed into an output file, carrying the Offset used
// to sort it within that file.
type entry struct {
	decl ast.Decl
	off  offsetKey
}

// assemble produces the merged output files from the collected symbols.
func assemble(pkgs []*PlatformPackage, pkgName string, syms map[string]*symbol, localNames map[importPath]string) (map[string]*ast.File, error) {
	// Deterministic symbol order for stable iteration (keys are unique).
	keys := make([]string, 0, len(syms))
	for k := range syms {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// outFiles maps output file name => entries to place in it.
	outFiles := map[string][]entry{}
	// platformTags maps a generated z_ file name => its [GOOS, GOARCH] tags, so
	// the build constraint can be written for it (Step 4).
	platformTags := map[string][2]string{}

	addEntry := func(name string, e entry) {
		outFiles[name] = append(outFiles[name], e)
	}

	for _, key := range keys {
		s := syms[key]
		if isShared(pkgs, s) {
			// Step 4: emit one shared copy into $(Fname).go.
			off, star := s.sharedOffset(pkgs)
			home := s.homeFile()
			src := s.byPlat[star].decl
			addEntry(home, entry{decl: copyDecl(src), off: off})
		} else {
			// Platform-specific: emit each platform's own copy into its z_ file.
			for pi, p := range pkgs {
				pl, ok := s.byPlat[pi]
				if !ok {
					continue
				}
				fname := platformFileName(p)
				platformTags[fname] = [2]string{p.Tags[0], p.Tags[1]}
				off := s.platOffset(pkgs, pi)
				addEntry(fname, entry{decl: copyDecl(pl.decl), off: off})
			}
		}
	}

	// Build the ast.File for each output file.
	result := map[string]*ast.File{}
	for _, fname := range sortedKeys(outFiles) {
		entries := outFiles[fname]
		if len(entries) == 0 {
			continue // a file with no declarations is not emitted
		}
		// Sort by Offset, ascending. Offsets are unique, so the order is total.
		sort.SliceStable(entries, func(a, b int) bool {
			return entries[a].off.less(entries[b].off)
		})
		decls := make([]ast.Decl, 0, len(entries)+1)
		for _, e := range entries {
			decls = append(decls, e.decl)
		}

		file := &ast.File{
			Name:  ast.NewIdent(pkgName),
			Decls: decls,
		}
		// Imports: regenerate from the import-path refs actually used.
		if importDecl := buildImports(decls, localNames); importDecl != nil {
			file.Decls = append([]ast.Decl{importDecl}, file.Decls...)
		}
		// Build constraint on generated per-platform files (Step 4). The z_ file
		// name already implies the constraint, but it is written explicitly so
		// the constraint is visible and independent of file-name parsing.
		if tags, ok := platformTags[fname]; ok {
			cg := &ast.CommentGroup{List: []*ast.Comment{
				{Text: "//go:build " + tags[0] + " && " + tags[1]},
			}}
			file.Doc = cg
			file.Comments = []*ast.CommentGroup{cg}
		}
		result[fname] = file
	}
	return result, nil
}

// isShared reports whether a symbol is identical on all input platforms (it
// exists on every platform and all copies are equal under Step 5).
func isShared(pkgs []*PlatformPackage, s *symbol) bool {
	if len(s.byPlat) != len(pkgs) {
		return false // missing on at least one platform
	}
	var first ast.Decl
	for pi := range pkgs {
		pl := s.byPlat[pi]
		if first == nil {
			first = pl.decl
			continue
		}
		if !declEqual(first, pl.decl) {
			return false
		}
	}
	return true
}

// platformFileName returns the generated per-platform file name
// z_$(GOOS)_$(GOARCH).go for a platform.
func platformFileName(p *PlatformPackage) string {
	return platformFilePrefix + p.Tags[0] + "_" + p.Tags[1] + ".go"
}

func sortedKeys(m map[string][]entry) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// -----------------------------------------------------------------------------
// Step 6: imports and build constraints.

// buildImports regenerates the import declaration for an output file from the
// extern package refs used in its declarations. A file that uses any
// //go:linkname directive additionally imports _ "unsafe".
func buildImports(decls []ast.Decl, localNames map[importPath]string) *ast.GenDecl {
	used := map[importPath]bool{}
	for _, d := range decls {
		collectUsedPaths(d, used)
	}
	needUnsafe := usesLinkname(decls)

	// Collect import specs: real package refs plus a blank unsafe import when a
	// linkname directive is present and unsafe is not already imported.
	type imp struct {
		path  importPath
		local string
		blank bool
	}
	var imps []imp
	hasUnsafe := false
	for path := range used {
		local := localNames[path]
		if local == "" {
			local = pkgNameOf(string(path))
		}
		if path == "unsafe" {
			hasUnsafe = true
		}
		imps = append(imps, imp{path: path, local: local})
	}
	if needUnsafe && !hasUnsafe {
		imps = append(imps, imp{path: "unsafe", local: "_", blank: true})
	}
	if len(imps) == 0 {
		return nil
	}
	// Deterministic order: by import path.
	sort.Slice(imps, func(a, b int) bool { return imps[a].path < imps[b].path })

	specs := make([]ast.Spec, 0, len(imps))
	for _, im := range imps {
		spec := &ast.ImportSpec{
			Path: &ast.BasicLit{Kind: token.STRING, Value: strconv.Quote(string(im.path))},
		}
		// Emit an explicit name when it differs from the default, or for blank.
		if im.blank {
			spec.Name = ast.NewIdent("_")
		} else if im.local != "" && im.local != pkgNameOf(string(im.path)) {
			spec.Name = ast.NewIdent(im.local)
		}
		specs = append(specs, spec)
	}
	return &ast.GenDecl{Tok: token.IMPORT, Specs: specs}
}

// collectUsedPaths records every extern package ref path used within a
// declaration (as annotated in Step 1).
func collectUsedPaths(node ast.Node, used map[importPath]bool) {
	ast.Inspect(node, func(n ast.Node) bool {
		if sel, ok := n.(*ast.SelectorExpr); ok {
			if id, ok := sel.X.(*ast.Ident); ok {
				if path, isRef := packageRefOf(id); isRef {
					used[path] = true
				}
			}
		}
		return true
	})
}

// usesLinkname reports whether any declaration carries a //go:linkname directive.
func usesLinkname(decls []ast.Decl) bool {
	for _, d := range decls {
		var doc *ast.CommentGroup
		switch x := d.(type) {
		case *ast.FuncDecl:
			doc = x.Doc
		case *ast.GenDecl:
			doc = x.Doc
		}
		for _, dir := range directivesOf(doc) {
			if strings.HasPrefix(dir, "go:linkname") {
				return true
			}
		}
	}
	return false
}
