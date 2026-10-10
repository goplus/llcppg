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
	"html/template"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// srcURLPrefix is the URL path under which package source files are served.
// It mirrors pkg.go.dev's behaviour of letting a reader jump from a symbol to
// the exact line of its definition, but serves the files llcppg just read from
// disk rather than a remote VCS.
const srcURLPrefix = "/src/"

// sourceLinkFunc returns a SourceLinkFunc for the package loaded into p: given a
// declaration's AST node it returns a "/src/<file>#L<line>" URL pointing at the
// line the declaration starts on. It returns "" (no link) when the node has no
// position or lies outside the package directory, so a symbol is never linked
// to a file the source handler would refuse to serve.
func (p *pkg) sourceLinkFunc() func(ast.Node) string {
	return func(n ast.Node) string {
		if n == nil {
			return ""
		}
		pos := p.FileSet.Position(n.Pos())
		if !pos.IsValid() || pos.Filename == "" {
			return ""
		}
		rel, ok := p.relSourceFile(pos.Filename)
		if !ok {
			return ""
		}
		return fmt.Sprintf("%s%s#L%d", srcURLPrefix, pathEscapeSlashes(rel), pos.Line)
	}
}

// fileLinkFunc returns a FileLinkFunc for p: given the base name of a package
// .go file it returns the "/src/<file>" URL that serves it. The name comes from
// go/doc, which uses the file's base name; it is linked only when the file is
// one the source handler will serve.
func (p *pkg) fileLinkFunc() func(string) string {
	return func(name string) string {
		rel, ok := p.relSourceFile(filepath.Join(p.Dir, name))
		if !ok {
			return ""
		}
		return srcURLPrefix + pathEscapeSlashes(rel)
	}
}

// relSourceFile validates that absPath is a .go file located directly in the
// package directory and returns its path relative to that directory (a bare
// base name). It rejects anything outside p.Dir, so neither the link functions
// nor a crafted request can escape the package directory.
func (p *pkg) relSourceFile(absPath string) (string, bool) {
	clean := filepath.Clean(absPath)
	rel, err := filepath.Rel(p.Dir, clean)
	if err != nil {
		return "", false
	}
	rel = filepath.ToSlash(rel)
	// Reject parent-directory escapes and nested paths: llcppg documents a
	// single directory, so every source file is a direct child of p.Dir.
	if rel == "." || rel == "" || strings.HasPrefix(rel, "../") || strings.Contains(rel, "/") {
		return "", false
	}
	if !strings.HasSuffix(rel, ".go") {
		return "", false
	}
	return rel, true
}

// pathEscapeSlashes percent-encodes a relative path for use in a URL path,
// leaving the path separators intact. Package source file names are plain base
// names today, but encoding keeps the URL well-formed if one ever contains a
// character that is unsafe in a path segment.
func pathEscapeSlashes(rel string) string {
	segs := strings.Split(rel, "/")
	for i, s := range segs {
		segs[i] = url.PathEscape(s)
	}
	return strings.Join(segs, "/")
}

// serveSource serves a package source file as a syntax-free, line-numbered HTML
// page whose lines carry "L<n>" anchors, so a "/src/<file>#L<line>" link from
// the documentation scrolls to the referenced declaration. It reloads the file
// on every request to match the rest of the server, and only ever serves .go
// files that live directly in dir.
func serveSource(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		name := strings.TrimPrefix(req.URL.Path, srcURLPrefix)
		if name == "" {
			http.NotFound(w, req)
			return
		}
		// req.URL.Path is already percent-decoded by net/http; validate it as a
		// bare base name under dir before touching the file system.
		p := &pkg{Dir: dir}
		rel, ok := p.relSourceFile(filepath.Join(dir, name))
		if !ok {
			http.NotFound(w, req)
			return
		}
		data, err := os.ReadFile(filepath.Join(dir, rel))
		if err != nil {
			http.NotFound(w, req)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := sourceTemplate.Execute(w, newSourceData(rel, data)); err != nil {
			fmt.Fprintf(w, "\n<!-- template error: %s -->\n", template.HTMLEscapeString(err.Error()))
		}
	}
}

// sourceData is the data passed to sourceTemplate.
type sourceData struct {
	Name  string       // the file's base name, for the title and heading
	Lines []sourceLine // the file split into anchored lines
}

// sourceLine is one line of a source file, numbered from 1.
type sourceLine struct {
	Num  int
	Text string
}

func newSourceData(name string, data []byte) *sourceData {
	// Trim a single trailing newline so a file ending in "\n" does not render a
	// spurious empty final line.
	text := strings.TrimSuffix(string(data), "\n")
	raw := strings.Split(text, "\n")
	lines := make([]sourceLine, len(raw))
	for i, s := range raw {
		lines[i] = sourceLine{Num: i + 1, Text: s}
	}
	return &sourceData{Name: path.Base(name), Lines: lines}
}

// sourceTemplate renders a source file with one anchored, numbered line per
// row. html/template escapes every line's text, so the raw source cannot inject
// markup.
var sourceTemplate = template.Must(template.New("source").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>{{.Name}} - llcppg doc source</title>
<link rel="stylesheet" href="/static/page.css">
</head>
<body>
<div id="page">
<div class="container">
<h1 class="src-title">{{.Name}}</h1>
<pre class="src-code">{{range .Lines}}<span id="L{{.Num}}" class="src-line"><a class="src-lineNum" href="#L{{.Num}}">{{.Num}}</a>{{.Text}}
</span>{{end}}</pre>
</div>
</div>
</body>
</html>
`))
