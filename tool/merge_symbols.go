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
	"fmt"
	"go/ast"
	"go/token"
	"sort"
)

// -----------------------------------------------------------------------------
// Step 2 + 3: Collect the global symbols and their placement.

// placement records, for one symbol on one platform, the declaration node, the
// file it came from, and its offset within that file.
type placement struct {
	platform int      // index into the (sorted) platform slice
	decl     ast.Decl // the declaration node (func, split spec wrapped in GenDecl, or const block)
	file     string   // file name the declaration came from
	offset   int      // byte offset of the declaration within that file (Step 3)
}

// symbol is one global symbol across all platforms that declare it.
type symbol struct {
	key    string
	byPlat map[int]placement // platform index => placement
}

// offsetKey is the lexicographically-comparable Offset of a placement:
// (file name, platform string, file offset).
type offsetKey struct {
	fname    string
	platform string
	foffset  int
}

func (a offsetKey) less(b offsetKey) bool {
	if a.fname != b.fname {
		return a.fname < b.fname
	}
	if a.platform != b.platform {
		return a.platform < b.platform
	}
	return a.foffset < b.foffset
}

// collectSymbols implements Steps 2 and 3: it walks every platform's files,
// splits declarations into global symbols, and records each symbol's placement
// (file and offset) per platform.
func collectSymbols(pkgs []*PlatformPackage) (map[string]*symbol, error) {
	syms := map[string]*symbol{}

	get := func(key string) *symbol {
		s := syms[key]
		if s == nil {
			s = &symbol{key: key, byPlat: map[int]placement{}}
			syms[key] = s
		}
		return s
	}

	for pi, p := range pkgs {
		seen := map[string]bool{} // within one platform, a symbol is declared once
		platStr := platformString(p)
		for _, fname := range sortedFileNames(p.Files) {
			f := p.Files[fname]
			for di, decl := range f.Decls {
				items, err := splitDecl(decl, platStr, fname)
				if err != nil {
					return nil, err
				}
				for _, it := range items {
					if seen[it.key] {
						return nil, fmt.Errorf("merge: symbol %q declared twice on platform %q", it.key, platStr)
					}
					seen[it.key] = true
					off := offsetOf(f, it.node, di, it.specIdx)
					get(it.key).byPlat[pi] = placement{
						platform: pi,
						decl:     it.decl,
						file:     fname,
						offset:   off,
					}
				}
			}
		}
	}
	return syms, nil
}

// declItem is one global symbol extracted from a top-level declaration.
type declItem struct {
	key     string
	decl    ast.Decl // the declaration to emit for this symbol
	node    ast.Node // the node whose Pos() gives the offset (spec or decl)
	specIdx int      // index of the spec within its GenDecl, or -1 for whole-decl
}

// splitDecl turns one top-level declaration into its global symbols (Step 2).
// A var/type GenDecl with several specs yields one symbol per spec; a const
// block is always one symbol as a whole; a func/method is one symbol.
func splitDecl(decl ast.Decl, platStr, fname string) ([]declItem, error) {
	switch d := decl.(type) {
	case *ast.FuncDecl:
		if d.Name == nil || d.Name.Name == "_" {
			return nil, fmt.Errorf("merge: blank top-level func in %q (platform %q)", fname, platStr)
		}
		if d.Recv == nil && d.Name.Name == "init" {
			return nil, fmt.Errorf("merge: init function cannot be merged (%q, platform %q)", fname, platStr)
		}
		return []declItem{{key: funcKey(d), decl: d, node: d, specIdx: -1}}, nil

	case *ast.GenDecl:
		switch d.Tok {
		case token.CONST:
			// A const block is always one symbol, keyed by the first name.
			key, err := constBlockKey(d, platStr, fname)
			if err != nil {
				return nil, err
			}
			return []declItem{{key: key, decl: d, node: d, specIdx: -1}}, nil
		case token.TYPE, token.VAR:
			var items []declItem
			for si, spec := range d.Specs {
				key, node, err := specKey(d, spec, platStr, fname)
				if err != nil {
					return nil, err
				}
				// Wrap each spec in its own GenDecl so it can be placed and
				// emitted independently while keeping the doc comment.
				items = append(items, declItem{
					key:     key,
					decl:    singleSpecDecl(d, spec),
					node:    node,
					specIdx: si,
				})
			}
			return items, nil
		case token.IMPORT:
			// Import declarations are not global symbols; they are regenerated.
			return nil, nil
		}
	}
	return nil, nil
}

// funcKey returns the symbol key for a func or method (Step 2). For a method the
// key is RecvType.Method; pointer-ness of the receiver is not part of the key.
func funcKey(d *ast.FuncDecl) string {
	if d.Recv == nil || len(d.Recv.List) == 0 {
		return d.Name.Name
	}
	return recvTypeName(d.Recv.List[0].Type) + "." + d.Name.Name
}

// recvTypeName returns the base type name of a receiver, dropping pointer and
// type-parameter decoration.
func recvTypeName(t ast.Expr) string {
	switch e := t.(type) {
	case *ast.StarExpr:
		return recvTypeName(e.X)
	case *ast.Ident:
		return e.Name
	case *ast.IndexExpr: // generic receiver T[P]
		return recvTypeName(e.X)
	case *ast.IndexListExpr:
		return recvTypeName(e.X)
	}
	return ""
}

func constBlockKey(d *ast.GenDecl, platStr, fname string) (string, error) {
	for _, spec := range d.Specs {
		vs, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}
		for _, name := range vs.Names {
			if name.Name == "_" {
				continue
			}
			return "const:" + name.Name, nil
		}
	}
	return "", fmt.Errorf("merge: const block with no usable name in %q (platform %q)", fname, platStr)
}

func specKey(d *ast.GenDecl, spec ast.Spec, platStr, fname string) (string, ast.Node, error) {
	switch s := spec.(type) {
	case *ast.TypeSpec:
		if s.Name == nil || s.Name.Name == "_" {
			return "", nil, fmt.Errorf("merge: blank type in %q (platform %q)", fname, platStr)
		}
		return "type:" + s.Name.Name, s, nil
	case *ast.ValueSpec:
		// A var spec may declare several names; each is its own symbol.
		// Generated code uses one name per spec, but handle the general case by
		// keying on the first non-blank name and treating extra names as part of
		// the same spec node.
		for _, name := range s.Names {
			if name.Name == "_" {
				return "", nil, fmt.Errorf("merge: blank var in %q (platform %q)", fname, platStr)
			}
		}
		if len(s.Names) == 0 {
			return "", nil, fmt.Errorf("merge: var spec with no name in %q (platform %q)", fname, platStr)
		}
		return "var:" + s.Names[0].Name, s, nil
	}
	return "", nil, fmt.Errorf("merge: unsupported spec in %q (platform %q)", fname, platStr)
}

// singleSpecDecl builds a GenDecl carrying a single spec, inheriting the parent
// block's doc comment only when the parent was not parenthesized (so a shared
// doc on a parenthesized block is not duplicated across its specs).
func singleSpecDecl(parent *ast.GenDecl, spec ast.Spec) *ast.GenDecl {
	doc := specDoc(spec)
	if doc == nil && !parent.Lparen.IsValid() {
		// "var x int" / "type T int": the doc belongs to the GenDecl.
		doc = parent.Doc
	}
	return &ast.GenDecl{
		Doc:   doc,
		Tok:   parent.Tok,
		Specs: []ast.Spec{spec},
	}
}

func specDoc(spec ast.Spec) *ast.CommentGroup {
	switch s := spec.(type) {
	case *ast.TypeSpec:
		return s.Doc
	case *ast.ValueSpec:
		return s.Doc
	}
	return nil
}

// offsetOf computes the byte offset of a node within a file (Step 3). It uses
// node.Pos() - file.FileStart when positions are usable, and otherwise falls
// back to a synthetic offset derived from the declaration/spec index so the
// ordering stays consistent for files without positions.
func offsetOf(f *ast.File, node ast.Node, declIdx, specIdx int) int {
	if node != nil && node.Pos().IsValid() && f.FileStart.IsValid() {
		return int(node.Pos() - f.FileStart)
	}
	// Fallback: declIdx dominates, specIdx orders within a block. Scale declIdx
	// so specIdx never crosses declaration boundaries.
	if specIdx < 0 {
		specIdx = 0
	}
	return declIdx*(1<<20) + specIdx
}

// sortedFileNames returns the file names of a package in deterministic order.
func sortedFileNames(files map[string]*ast.File) []string {
	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// -----------------------------------------------------------------------------
// Step 3: home file and Offset of a symbol.

// homeFile returns the home file name of a symbol: if every platform declares it
// in the same file, that file; otherwise the lexicographically smallest file
// name among them.
func (s *symbol) homeFile() string {
	home := ""
	for _, pl := range s.byPlat {
		if home == "" || pl.file < home {
			home = pl.file
		}
	}
	return home
}

// sharedOffset returns the single Offset of a shared symbol: the smallest of its
// per-platform Offsets, together with the platform P* whose copy is used. P* is
// the first platform (in sorted order) that declares S in the home file.
func (s *symbol) sharedOffset(pkgs []*PlatformPackage) (offsetKey, int) {
	home := s.homeFile()
	var best offsetKey
	bestPlat := -1
	for pi := range pkgs {
		pl, ok := s.byPlat[pi]
		if !ok || pl.file != home {
			continue
		}
		ok2 := offsetKey{fname: pl.file, platform: platformString(pkgs[pi]), foffset: pl.offset}
		if bestPlat == -1 || ok2.less(best) {
			best = ok2
			bestPlat = pi
		}
	}
	return best, bestPlat
}

// platOffset returns the Offset of a symbol on one platform.
func (s *symbol) platOffset(pkgs []*PlatformPackage, pi int) offsetKey {
	pl := s.byPlat[pi]
	return offsetKey{fname: pl.file, platform: platformString(pkgs[pi]), foffset: pl.offset}
}
