// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package godoc renders Go package documentation to HTML in the pkg.go.dev
// style.
//
// It is a trimmed, vendored copy of golang.org/x/pkgsite/internal/godoc,
// reduced to just what llcppg's local documentation server needs: build a
// Package from already-parsed *ast.File values and render it to HTML. The
// pkgsite original additionally encodes/decodes ASTs for database storage and
// extracts module/version/source metadata; none of that applies to llcppg,
// which renders a single, freshly generated package straight from disk, so the
// codec, symbol extraction, and source-linking machinery are omitted.
//
// This code (and the HTML templates under dochtml/doc) is BSD-licensed Go
// Authors code; see the LICENSE file in this directory.
package godoc

import (
	"go/ast"
	"go/doc"
	"go/token"
	"strings"

	"github.com/goplus/llcppg/internal/dochtml"
)

// ErrTooLarge is returned (wrapped) when rendered documentation exceeds the
// configured size limit.
var ErrTooLarge = dochtml.ErrTooLarge

// ModuleInfo re-exports dochtml.ModuleInfo for callers that still construct it.
// llcppg does not populate it, but keeping the alias mirrors the pkgsite API.
type ModuleInfo = dochtml.ModuleInfo

// A Package contains package-level information needed to render Go
// documentation.
type Package struct {
	Fset  *token.FileSet
	Files []*File
	// docMode is the doc.Mode used when building documentation. Setting
	// doc.AllDecls includes unexported declarations.
	docMode doc.Mode

	renderCalled bool
}

// A File contains everything needed about a source file to render
// documentation.
type File struct {
	Name string // file pathname
	AST  *ast.File
}

// NewPackage returns a new Package with the given file set. If allDecls is
// true, unexported declarations are included (doc.AllDecls).
func NewPackage(fset *token.FileSet, allDecls bool) *Package {
	p := &Package{Fset: fset}
	if allDecls {
		p.docMode |= doc.AllDecls
	}
	return p
}

// AddFile adds a file to the Package. After it returns, the contents of the
// ast.File are unsuitable for anything other than the methods of this package.
func (p *Package) AddFile(f *ast.File, removeNodes bool) {
	filename := p.Fset.Position(f.Package).Filename
	// Don't trim anything from a test file or one in a XXX_test package; it
	// may be part of a playable example.
	if removeNodes && !strings.HasSuffix(filename, "_test.go") && !strings.HasSuffix(f.Name.Name, "_test") {
		removeUnusedASTNodes(f)
	}
	p.Files = append(p.Files, &File{
		Name: filename,
		AST:  f,
	})
}

// removeUnusedASTNodes removes parts of the AST not needed for documentation.
// It doesn't remove unexported consts, vars or types, although it probably
// could.
func removeUnusedASTNodes(pf *ast.File) {
	var decls []ast.Decl
	for _, d := range pf.Decls {
		if f, ok := d.(*ast.FuncDecl); ok {
			// Remove all unexported functions and function bodies.
			if f.Name == nil || !ast.IsExported(f.Name.Name) {
				continue
			}
			// Remove the function body, unless it's an example.
			// The doc contains example bodies.
			if !strings.HasPrefix(f.Name.Name, "Example") {
				f.Body = nil
			}
		}
		decls = append(decls, d)
	}
	// Don't remove pf.Comments; they may contain Notes.
	pf.Decls = decls
}
