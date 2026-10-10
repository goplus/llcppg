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
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// safeBuffer is a goroutine-safe string buffer so the test can read the
// server's startup output while Run writes to it from another goroutine.
type safeBuffer struct {
	mu  sync.Mutex
	buf strings.Builder
}

func (b *safeBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *safeBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

func TestHandlerServesPackage(t *testing.T) {
	dir := writeDir(t, map[string]string{
		"a.go": "// Package foo is a fixture.\npackage foo\n\n// Bar does nothing.\nfunc Bar() {}\n",
	})
	srv := httptest.NewServer(newHandler(dir, false, nil, nil))
	defer srv.Close()

	body := getBody(t, srv.URL+"/", http.StatusOK)
	if !strings.Contains(body, `id="Bar"`) {
		t.Error("page should contain an anchor for Bar")
	}
	if !strings.Contains(body, "Package foo is a fixture.") {
		t.Error("page should contain the package overview")
	}
}

func TestHandlerServesStatic(t *testing.T) {
	dir := writeDir(t, map[string]string{"a.go": "package foo\n"})
	srv := httptest.NewServer(newHandler(dir, false, nil, nil))
	defer srv.Close()

	// llcppg's stylesheet, which styles both the shell chrome and the
	// godoc-rendered body's Documentation-* classes.
	body := getBody(t, srv.URL+"/static/page.css", http.StatusOK)
	if !strings.Contains(body, "parse-errors") {
		t.Error("shell stylesheet should be served from embedded assets")
	}
	if !strings.Contains(body, ".Documentation") {
		t.Error("stylesheet should also style the godoc body (.Documentation)")
	}

	// The renderer no longer depends on x/tools/godoc's bundled assets, so the
	// old /godoc/ path must not be served.
	resp, err := http.Get(srv.URL + "/godoc/style.css")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("GET /godoc/style.css: status = %d, want 404 (no x/tools assets)", resp.StatusCode)
	}
}

func TestHandlerServesSource(t *testing.T) {
	dir := writeDir(t, map[string]string{
		"a.go": "// Package foo is a fixture.\npackage foo\n\n// Bar does nothing.\nfunc Bar() {}\n",
	})
	srv := httptest.NewServer(newHandler(dir, false, nil, nil))
	defer srv.Close()

	// The documentation links symbols to /src/<file>#L<line>; that endpoint
	// must serve the file with per-line anchors so the fragment resolves.
	body := getBody(t, srv.URL+"/src/a.go", http.StatusOK)
	if !strings.Contains(body, "func Bar() {}") {
		t.Error("source view should contain the file's text")
	}
	if !strings.Contains(body, `id="L5"`) {
		t.Error("source view should carry a per-line anchor the #L fragment targets")
	}

	// Raw source text must be HTML-escaped, not emitted as markup.
	escDir := writeDir(t, map[string]string{"b.go": "package foo\n\nvar X = \"<b>\"\n"})
	esrv := httptest.NewServer(newHandler(escDir, false, nil, nil))
	defer esrv.Close()
	ebody := getBody(t, esrv.URL+"/src/b.go", http.StatusOK)
	if strings.Contains(ebody, `"<b>"`) || !strings.Contains(ebody, "&lt;b&gt;") {
		t.Error("source view must HTML-escape the file contents")
	}
}

func TestHandlerSourceRejectsBadPaths(t *testing.T) {
	dir := writeDir(t, map[string]string{
		"a.go":      "package foo\n",
		"notes.txt": "secret\n",
	})
	srv := httptest.NewServer(newHandler(dir, false, nil, nil))
	defer srv.Close()

	// Only .go files directly in the package directory are served; a non-.go
	// file, a nested path, and a bare prefix all 404.
	for _, p := range []string{
		"/src/notes.txt",
		"/src/sub/x.go",
		"/src/",
	} {
		resp, err := http.Get(srv.URL + p)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("GET %s: status = %d, want 404", p, resp.StatusCode)
		}
	}
}

func TestHandlerSourceRejectsOversizeFile(t *testing.T) {
	dir := writeDir(t, map[string]string{
		"big.go":   "package foo\n",
		"small.go": "package foo\n",
	})
	srv := httptest.NewServer(newHandler(dir, false, nil, nil))
	defer srv.Close()

	// Lower the cap so a tiny file counts as "oversize" without writing a huge
	// fixture; restore it afterwards.
	orig := maxSourceFileSize
	maxSourceFileSize = 1
	defer func() { maxSourceFileSize = orig }()

	// A file over the limit is declined (404); a file within it still serves.
	resp, err := http.Get(srv.URL + "/src/big.go")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("oversize file: status = %d, want 404", resp.StatusCode)
	}

	maxSourceFileSize = orig
	getBody(t, srv.URL+"/src/small.go", http.StatusOK)
}

func TestHandlerLiveRefresh(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "a.go")
	if err := os.WriteFile(src, []byte("package foo\n\nfunc First() {}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(newHandler(dir, false, nil, nil))
	defer srv.Close()

	body := getBody(t, srv.URL+"/", http.StatusOK)
	if !strings.Contains(body, `id="First"`) {
		t.Fatal("first load should show First")
	}

	// Modify the source; the next request must reflect it without a restart.
	if err := os.WriteFile(src, []byte("package foo\n\nfunc Second() {}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	body = getBody(t, srv.URL+"/", http.StatusOK)
	if !strings.Contains(body, `id="Second"`) {
		t.Error("second load should reflect the edited source (Second)")
	}
	if strings.Contains(body, `id="First"`) {
		t.Error("second load should no longer show the removed First")
	}
}

func TestHandlerNotFound(t *testing.T) {
	dir := writeDir(t, map[string]string{"a.go": "package foo\n"})
	srv := httptest.NewServer(newHandler(dir, false, nil, nil))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/does-not-exist")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", resp.StatusCode)
	}
}

func TestRunRejectsEmptyDir(t *testing.T) {
	dir := t.TempDir() // no Go files
	err := Run(context.Background(), dir, Options{
		OpenBrowser: false,
		Stdout:      io.Discard,
		Stderr:      io.Discard,
	})
	if err == nil {
		t.Fatal("Run should reject a directory with no Go files")
	}
}

func TestRunServesAndShutsDown(t *testing.T) {
	dir := writeDir(t, map[string]string{
		"a.go": "// Package foo is a fixture.\npackage foo\n\nfunc Bar() {}\n",
	})

	ctx, cancel := context.WithCancel(context.Background())
	var out safeBuffer
	done := make(chan error, 1)
	go func() {
		done <- Run(ctx, dir, Options{
			Addr:        "127.0.0.1:0",
			OpenBrowser: false,
			Stdout:      &out,
			Stderr:      io.Discard,
		})
	}()

	// Wait until the server has printed its URL, then fetch it.
	url := waitForURL(t, &out)
	getBody(t, url, http.StatusOK)

	// Cancelling the context should shut the server down and return nil.
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Errorf("Run returned %v, want nil after shutdown", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Run did not shut down within 5s")
	}
}

// waitForURL polls out until the startup line "... at http://..." appears and
// returns the URL.
func waitForURL(t *testing.T, out interface{ String() string }) string {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		s := out.String()
		if i := strings.Index(s, "http://"); i >= 0 {
			url := strings.TrimSpace(s[i:])
			if nl := strings.IndexAny(url, " \n"); nl >= 0 {
				url = url[:nl]
			}
			return url
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("server did not print a URL; output so far: %q", out.String())
	return ""
}

func getBody(t *testing.T, url string, wantStatus int) string {
	t.Helper()
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != wantStatus {
		t.Fatalf("GET %s: status = %d, want %d; body: %s", url, resp.StatusCode, wantStatus, data)
	}
	return string(data)
}
