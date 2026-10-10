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

//go:embed assets/page.html assets/page.css assets/live.js assets/outline.js
var assets embed.FS

// pageTemplate is parsed once from the embedded shell template.
var pageTemplate = template.Must(template.ParseFS(assets, "assets/page.html"))

// newHandler builds the HTTP handler. "GET /" loads and renders the package in
// dir on every request so edits are picked up without a restart; "GET
// /static/..." serves llcppg's own assets (the page.css stylesheet and, when
// live reload is on, live.js); "GET /src/..." serves the package source so
// documentation links can point at it. Nothing else on the file system is
// served, and package code is never executed. dir is the absolute package
// directory resolved by load(), so the source handler and the symbol links
// generated during rendering agree on the same path.
//
// When b is non-nil the server is in live-reload mode: it registers the
// "GET /_events" Server-Sent Events route and the rendered page pulls in the
// live-reload client script. When b is nil the server is degraded (the file
// watcher could not start); /_events is not registered, so it returns 404, and
// the page carries no live-reload script.
//
// done is the server's lifetime signal passed to the event handler so open
// streams close on shutdown (see broker.serveEvents); it is ignored when b is
// nil.
func newHandler(dir string, allDecls bool, b *broker, done <-chan struct{}) http.Handler {
	mux := http.NewServeMux()

	// Serve llcppg's assets (page.css, outline.js, live.js). The embed root has
	// an "assets/" prefix; strip the URL "/static/" prefix and re-root onto it.
	static, _ := fs.Sub(assets, "assets")
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(static))))

	// Serve the package source so symbol and "View Source" links from the
	// documentation can jump to the exact line of a declaration. Only .go files
	// in dir are served (see serveSource); package code is never executed.
	mux.HandleFunc(srcURLPrefix, serveSource(dir))

	liveReload := b != nil
	if liveReload {
		mux.HandleFunc("/_events", b.serveEvents(keepAlivePeriod, done))
	}

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
		if err := pageTemplate.Execute(w, p.render(liveReload)); err != nil {
			// The header may already be written; there is nothing more we can
			// do except surface the error for diagnosis.
			fmt.Fprintf(w, "\n<!-- template error: %s -->\n", template.HTMLEscapeString(err.Error()))
		}
	})

	return mux
}
