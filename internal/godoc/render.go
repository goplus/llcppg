// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package godoc

import (
	"context"
	"errors"
	"fmt"
	"go/ast"
	"go/doc"
	"go/token"
	"strings"

	"github.com/google/safehtml/template"
	"github.com/goplus/llcppg/internal/derrors"
	"github.com/goplus/llcppg/internal/dochtml"
	"github.com/goplus/llcppg/internal/importer"
)

const (
	megabyte             = 1000 * 1000
	maxImportsPerPackage = 5000

	// Exported for tests.
	DocTooLargeReplacement = `<p>Documentation is too large to display.</p>`
)

// MaxDocumentationHTML is a limit on the rendered documentation HTML size.
//
// The current limit of is based on the largest packages that
// pkg.go.dev has encountered. See https://golang.org/issue/40576.
//
// It is a variable for testing.
var MaxDocumentationHTML = 40 * megabyte

// DocPackage computes and returns a doc.Package for the given import path.
//
// It is a trimmed version of the pkgsite original: llcppg always renders a
// single, locally parsed package, so the module/version and standard-library
// ("builtin") special cases of pkgsite's DocPackage are not needed.
//
// doc.NewFromFiles consumes the comment associations of the ASTs it is given,
// so DocPackage (and therefore Render) must be called at most once for a given
// set of files; a caller that already holds a *doc.Package should use the
// package-level RenderDoc instead.
func (p *Package) DocPackage(importPath string) (_ *doc.Package, err error) {
	defer derrors.Wrap(&err, "DocPackage(%q)", importPath)

	var allGoFiles []*ast.File
	nonTestFiles := map[string]*ast.File{}
	for _, f := range p.Files {
		allGoFiles = append(allGoFiles, f.AST)
		if !strings.HasSuffix(f.Name, "_test.go") {
			nonTestFiles[f.Name] = f.AST
		}
	}
	// Call ast.NewPackage for side-effects to populate objects. In Go 1.25+
	// doc.NewFromFiles will not cause the objects to be populated.
	//lint:ignore SA1019 We had a preexisting dependency on ast.Object.
	ast.NewPackage(p.Fset, nonTestFiles, importer.SimpleImporter, nil)
	d, err := doc.NewFromFiles(p.Fset, allGoFiles, importPath, p.docMode)
	if err != nil {
		return nil, fmt.Errorf("doc.NewFromFiles: %v", err)
	}

	// Process package imports.
	if len(d.Imports) > maxImportsPerPackage {
		return nil, fmt.Errorf("%d imports found package %q; exceeds limit %d for maxImportsPerPackage", len(d.Imports), importPath, maxImportsPerPackage)
	}
	return d, nil
}

// Render renders the documentation for the package under importPath.
//
// It is a trimmed version of the pkgsite original: there is no module
// versioning or "since version" information, because llcppg serves a single
// freshly generated package from disk rather than a versioned module from a
// database. Source links are left disabled (RenderDoc's defaults); a caller
// that can serve the package source should use RenderDocLinked instead.
// Rendering destroys p's AST; do not call any methods of p after it returns.
func (p *Package) Render(ctx context.Context, importPath string) (_ *dochtml.Parts, err error) {
	p.renderCalled = true

	d, err := p.DocPackage(importPath)
	if err != nil {
		return nil, err
	}
	return RenderDoc(ctx, p.Fset, d)
}

// LinkOptions configures how RenderDocLinked turns declarations and files into
// links back to their source. A nil func (or one returning the empty string)
// disables the corresponding link, matching RenderDoc's behaviour.
type LinkOptions struct {
	// SourceLinkFunc returns the URL to link a declaration's source to, given
	// the declaration's AST node. Returning "" leaves the symbol unlinked.
	SourceLinkFunc func(ast.Node) string
	// FileLinkFunc returns the URL to link a .go file name to. Returning ""
	// leaves the file unlinked.
	FileLinkFunc func(file string) string
}

// RenderDoc renders an already-built doc.Package to HTML with source linking
// disabled. See RenderDocLinked for the variant that links symbols to their
// source. fset must be the file set the package was parsed with.
func RenderDoc(ctx context.Context, fset *token.FileSet, d *doc.Package) (_ *dochtml.Parts, err error) {
	return RenderDocLinked(ctx, fset, d, LinkOptions{})
}

// RenderDocLinked renders an already-built doc.Package to HTML, linking symbols
// and files to their source via opts. It is useful for callers (like llcppg's
// docserver) that construct the doc.Package themselves with go/doc and want the
// pkg.go.dev-style HTML rendering, without handing the raw ASTs to this
// package. fset must be the file set the package was parsed with.
func RenderDocLinked(ctx context.Context, fset *token.FileSet, d *doc.Package, opts LinkOptions) (_ *dochtml.Parts, err error) {
	sourceLink := opts.SourceLinkFunc
	if sourceLink == nil {
		sourceLink = func(ast.Node) string { return "" }
	}
	fileLink := opts.FileLinkFunc
	if fileLink == nil {
		fileLink = func(string) string { return "" }
	}
	renderOpts := dochtml.RenderOptions{
		FileLinkFunc:     fileLink,
		SourceLinkFunc:   sourceLink,
		SinceVersionFunc: func(string) string { return "" },
		Limit:            int64(MaxDocumentationHTML),
	}
	parts, err := dochtml.Render(ctx, fset, d, renderOpts)
	if errors.Is(err, ErrTooLarge) {
		return &dochtml.Parts{Body: template.MustParseAndExecuteToHTML(DocTooLargeReplacement)}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("dochtml.Render: %v", err)
	}
	return parts, nil
}
