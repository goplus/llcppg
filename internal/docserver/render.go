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
	"context"
	"html/template"

	"github.com/goplus/llcppg/internal/godoc"
)

// pageData is the data passed to the page shell template (assets/page.html).
// The package body itself is rendered by the vendored internal/godoc package so
// the output matches the familiar pkg.go.dev layout; this struct only carries
// the shell metadata and that rendered body.
type pageData struct {
	ImportPath  string
	Name        string
	Dir         string
	ParseErrors []string

	// LiveReload is true when the server started a file watcher, so the page
	// should pull in the live-reload client script. A degraded server (no
	// watcher) leaves it false, and the page is exactly as it was before live
	// reload existed.
	LiveReload bool

	// Body is the package documentation rendered by internal/godoc.
	Body template.HTML

	// Outline is the documentation outline (Overview, Index, Constants,
	// Variables, Functions, Types with their constructors and methods, Notes)
	// rendered by internal/godoc. It is the same safe HTML dochtml already
	// produces alongside Body; the page shell drops it into the sidebar tree.
	// Its anchors (#pkg-overview, #{Type}, #{Type}.{Method}, …) agree with the
	// ids in Body by construction, because both are rendered from the same
	// TemplateData.
	Outline template.HTML

	// SourceFiles lists the Go files go/build selected for the package, each
	// linking to the existing /src/ source view. It drives both the "Source
	// Files" sidebar entry and the section at the bottom of the page.
	SourceFiles []sourceFile
}

// sourceFile is one entry in the "Source Files" list: the file's base name and
// the URL of its source view.
type sourceFile struct {
	Name string
	URL  string
}

// sourceFiles turns the package's selected file names into links to the
// existing /src/ source view. The names come from go/build (already sorted,
// _test.go excluded), and each file name is a bare base name validated by the
// /src/ handler, so the link resolves to exactly the file llcppg read.
func sourceFiles(names []string) []sourceFile {
	out := make([]sourceFile, 0, len(names))
	for _, n := range names {
		out = append(out, sourceFile{Name: n, URL: srcURLPrefix + pathEscapeSlashes(n)})
	}
	return out
}

// render builds the page-shell data for p, rendering the package body with the
// vendored internal/godoc package (a trimmed copy of x/pkgsite's godoc).
// liveReload selects whether the page includes the live-reload client script.
func (p *pkg) render(liveReload bool) *pageData {
	data := &pageData{
		ImportPath:  p.ImportPath,
		Name:        p.Doc.Name,
		Dir:         p.Dir,
		ParseErrors: p.ParseErrors,
		LiveReload:  liveReload,
		SourceFiles: sourceFiles(p.Files),
	}

	// load() already built p.Doc with go/doc; render that directly. We must not
	// rebuild the doc from the raw ASTs here, because go/doc.NewFromFiles
	// consumes the files' comment associations, so a second pass would drop all
	// the doc comments.
	parts, err := godoc.RenderDocLinked(context.Background(), p.FileSet, p.Doc, godoc.LinkOptions{
		SourceLinkFunc: p.sourceLinkFunc(),
	})
	if err != nil {
		// Rendering the body failed; surface it inside the body rather than
		// dropping the whole page. The shell (import path, parse errors) is
		// still useful on its own.
		data.Body = template.HTML("<pre class=\"render-error\">doc render error: " +
			template.HTMLEscapeString(err.Error()) + "</pre>")
		return data
	}

	// parts.Body is a safehtml.HTML produced by godoc's safe templates; its
	// String is known to be safe HTML, so promoting it to template.HTML for the
	// shell template is correct (no double-escaping). parts.Outline is produced
	// by the same safe templates and gets the same treatment.
	data.Body = template.HTML(parts.Body.String())
	data.Outline = template.HTML(parts.Outline.String())
	return data
}
