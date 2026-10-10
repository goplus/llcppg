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
	}

	// load() already built p.Doc with go/doc; render that directly. We must not
	// rebuild the doc from the raw ASTs here, because go/doc.NewFromFiles
	// consumes the files' comment associations, so a second pass would drop all
	// the doc comments.
	parts, err := godoc.RenderDocLinked(context.Background(), p.FileSet, p.Doc, godoc.LinkOptions{
		SourceLinkFunc: p.sourceLinkFunc(),
		FileLinkFunc:   p.fileLinkFunc(),
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
	// shell template is correct (no double-escaping).
	data.Body = template.HTML(parts.Body.String())
	return data
}
