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
	"bytes"
	"html/template"
	"sync"
	texttemplate "text/template"

	"golang.org/x/tools/godoc"
	"golang.org/x/tools/godoc/static"
	"golang.org/x/tools/godoc/vfs/mapfs"
)

// pageData is the data passed to the page shell template (assets/page.html).
// The package body itself is rendered by x/tools/godoc so the output matches
// the familiar pkg.go.dev / godoc layout; this struct only carries the shell
// metadata and that rendered body.
type pageData struct {
	ImportPath  string
	Name        string
	Dir         string
	ParseErrors []string

	// Body is the package documentation rendered by godoc's package template.
	Body template.HTML
}

// pres is the godoc Presentation used purely for its template FuncMap and node
// formatting. It is backed by an empty in-memory filesystem because we never
// use godoc's corpus, indexing, search, or file serving — only its package
// template and the AST/comment formatting helpers that template needs.
//
// Building it is cheap but not free, so it is created once and reused; the
// Presentation is only read during rendering (template execution), which the
// http.Server already serializes per request goroutine without sharing it.
var (
	presOnce sync.Once
	pres     *godoc.Presentation
	pkgTmpl  *texttemplate.Template
)

func initPres() {
	corpus := godoc.NewCorpus(mapfs.New(map[string]string{}))
	pres = godoc.NewPresentation(corpus)
	// Link declarations and identifiers to their in-page anchors, matching
	// godoc's own behaviour.
	pres.DeclLinks = true
	// godoc's package template is a text/template (not html/template): its
	// FuncMap helpers such as comment_html and node_html already return escaped,
	// safe HTML, so parsing it with html/template would double-escape the
	// output. We mirror godoc's own choice here.
	pkgTmpl = texttemplate.Must(texttemplate.New("package.html").
		Funcs(pres.FuncMap()).
		Parse(static.Files["package.html"]))
}

// render builds the page-shell data for p, rendering the package body with
// x/tools/godoc's package template.
func (p *pkg) render() *pageData {
	presOnce.Do(initPres)

	dp := p.Doc
	info := &godoc.PageInfo{
		Dirname:  p.Dir,
		FSet:     p.FileSet,
		PDoc:     dp,
		Examples: dp.Examples,
		Notes:    dp.Notes,
		// AnalysisData and CallGraph must be valid JavaScript literals because
		// the template emits them into a <script> block. We do not run godoc's
		// analysis, so render them as JSON null; the "null" string also makes
		// the template skip the call-graph section.
		AnalysisData: template.JS("null"),
		CallGraph:    template.JS("null"),
	}

	data := &pageData{
		ImportPath:  dp.ImportPath,
		Name:        dp.Name,
		Dir:         p.Dir,
		ParseErrors: p.ParseErrors,
	}

	var buf bytes.Buffer
	if err := pkgTmpl.Execute(&buf, info); err != nil {
		// Rendering the body failed; surface it inside the body rather than
		// dropping the whole page. The shell (import path, parse errors) is
		// still useful on its own.
		data.Body = template.HTML("<pre class=\"render-error\">doc render error: " +
			template.HTMLEscapeString(err.Error()) + "</pre>")
		return data
	}
	data.Body = template.HTML(buf.String())
	return data
}
