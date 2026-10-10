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
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
)

//go:embed assets/page.html assets/style.css
var assets embed.FS

// pageTemplate is parsed once from the embedded template.
var pageTemplate = template.Must(template.ParseFS(assets, "assets/page.html"))

// newHandler builds the HTTP handler. "GET /" loads and renders the package in
// dir on every request so edits are picked up without a restart; "GET
// /static/..." serves the embedded assets. Nothing on the file system is
// served, and package code is never executed.
func newHandler(dir string, allDecls bool) http.Handler {
	mux := http.NewServeMux()

	// Serve embedded static assets (the stylesheet). The embed root has an
	// "assets/" prefix; strip the URL "/static/" prefix and re-root onto it.
	static, _ := fs.Sub(assets, "assets")
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(static))))

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		p, err := load(dir, allDecls)
		if err != nil {
			// A load error here means the directory stopped being a single
			// importable package since start-up (for example all files were
			// removed mid-edit). Report it rather than rendering a blank page.
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := pageTemplate.Execute(w, p.render()); err != nil {
			// The header may already be written; there is nothing more we can
			// do except surface the error for diagnosis.
			fmt.Fprintf(w, "\n<!-- template error: %s -->\n", template.HTMLEscapeString(err.Error()))
		}
	})

	return mux
}
