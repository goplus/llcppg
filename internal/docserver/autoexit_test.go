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
	"net/url"
	"strings"
	"testing"
	"time"
)

// autoExitServer starts an httptest server whose handler runs in auto-exit
// mode. It returns the server plus a channel that is closed when the handler's
// cancel is called, so a test can wait for the auto-exit to fire. The handler's
// host is set to the test server's own host so the bye beacon's Host check
// passes for same-origin requests.
func autoExitServer(t *testing.T, dir string) (*httptest.Server, <-chan struct{}, *safeBuffer) {
	t.Helper()
	cancelled := make(chan struct{})
	var out safeBuffer
	// The handler needs the listen host for the Host check, but httptest picks
	// the port only after NewServer. Resolve it from the URL afterwards and set
	// it on the server the handler closes over.
	s := &server{dir: dir, autoExit: true, out: &out}
	s.cancel = func() {
		select {
		case <-cancelled:
		default:
			close(cancelled)
		}
	}
	srv := httptest.NewServer(s.mux())
	t.Cleanup(srv.Close)
	s.host = strings.TrimPrefix(srv.URL, "http://")
	return srv, cancelled, &out
}

// postBye sends POST /_doc/bye to srv with the given Origin and Host overrides.
// An empty origin omits the header; an empty host leaves the client default
// (the server's own host, which is what a same-origin browser sends).
func postBye(t *testing.T, srv *httptest.Server, origin, host string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, srv.URL+byeURL, nil)
	if err != nil {
		t.Fatal(err)
	}
	if origin != "" {
		req.Header.Set("Origin", origin)
	}
	if host != "" {
		req.Host = host
	}
	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	return resp
}

// TestAutoExitFirstPageMonitored verifies that only the first GET / claims the
// monitored page: it sets the sessionStorage marker, later pages do not.
func TestAutoExitFirstPageMonitored(t *testing.T) {
	dir := writeDir(t, map[string]string{"a.go": "package foo\n\nfunc Bar() {}\n"})
	srv, _, _ := autoExitServer(t, dir)

	first := getBody(t, srv.URL+"/", http.StatusOK)
	if !strings.Contains(first, `sessionStorage.setItem("llcppg-doc"`) {
		t.Error("first page should carry the monitor script")
	}
	// html/template renders the Monitor bool as a JS literal, padding it with
	// spaces in a script context (if ( true ) ...), so match on the literal.
	if !monitorTrue(first) {
		t.Errorf("first page should render Monitor=true; body did not contain the true marker")
	}

	second := getBody(t, srv.URL+"/", http.StatusOK)
	if monitorTrue(second) {
		t.Error("second page should render Monitor=false, so it does not set the marker")
	}
	// Both pages still carry the pagehide beacon wiring (a reloaded monitored
	// tab keeps its marker and must still beacon), so the script block is
	// present on every page in auto-exit mode.
	if !strings.Contains(second, "/_doc/bye") {
		t.Error("auto-exit pages should reference the bye endpoint")
	}
}

// TestAutoExitDisabledNoScript verifies that with auto-exit off the page has no
// monitor script and /_doc/bye is not registered.
func TestAutoExitDisabledNoScript(t *testing.T) {
	dir := writeDir(t, map[string]string{"a.go": "package foo\n\nfunc Bar() {}\n"})
	srv := httptest.NewServer(newHandler(dir, false, nil, nil, false, "", nil, nil))
	defer srv.Close()

	body := getBody(t, srv.URL+"/", http.StatusOK)
	if strings.Contains(body, "llcppg-doc") || strings.Contains(body, "/_doc/bye") {
		t.Error("page must not carry the auto-exit script when auto-exit is off")
	}

	resp, err := http.Post(srv.URL+byeURL, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("POST %s with auto-exit off: status = %d, want 404", byeURL, resp.StatusCode)
	}
}

// TestAutoExitByeShutsDown verifies that a valid POST /_doc/bye shuts the
// server down after the grace period.
func TestAutoExitByeShutsDown(t *testing.T) {
	orig := exitGrace
	exitGrace = 50 * time.Millisecond
	defer func() { exitGrace = orig }()

	dir := writeDir(t, map[string]string{"a.go": "package foo\n"})
	srv, cancelled, out := autoExitServer(t, dir)

	resp := postBye(t, srv, "", "")
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("POST bye: status = %d, want 200", resp.StatusCode)
	}

	select {
	case <-cancelled:
	case <-time.After(2 * time.Second):
		t.Fatal("auto-exit did not fire within 2s of bye")
	}
	if !strings.Contains(out.String(), "documentation page closed") {
		t.Errorf("shutdown notice not printed; got %q", out.String())
	}
}

// TestAutoExitReloadCancels verifies that a GET / arriving within the grace
// period after bye cancels the pending exit, so a reload does not shut the
// server down.
func TestAutoExitReloadCancels(t *testing.T) {
	orig := exitGrace
	exitGrace = 200 * time.Millisecond
	defer func() { exitGrace = orig }()

	dir := writeDir(t, map[string]string{"a.go": "package foo\n"})
	srv, cancelled, _ := autoExitServer(t, dir)

	// Claim the page, then simulate a reload: bye from the old page followed
	// immediately by a GET / from the new page.
	getBody(t, srv.URL+"/", http.StatusOK)
	resp := postBye(t, srv, "", "")
	resp.Body.Close()
	getBody(t, srv.URL+"/", http.StatusOK)

	select {
	case <-cancelled:
		t.Fatal("server exited despite a reload within the grace period")
	case <-time.After(exitGrace + 150*time.Millisecond):
		// Good: the GET / cancelled the timer.
	}
}

// TestAutoExitByeRejectsBadRequests verifies that a GET, a mismatched Origin,
// or a wrong Host does not start the exit timer.
func TestAutoExitByeRejectsBadRequests(t *testing.T) {
	orig := exitGrace
	exitGrace = 50 * time.Millisecond
	defer func() { exitGrace = orig }()

	dir := writeDir(t, map[string]string{"a.go": "package foo\n"})

	t.Run("GET not POST", func(t *testing.T) {
		srv, cancelled, _ := autoExitServer(t, dir)
		resp, err := http.Get(srv.URL + byeURL)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		assertNoExit(t, cancelled)
	})

	t.Run("mismatched Origin", func(t *testing.T) {
		srv, cancelled, _ := autoExitServer(t, dir)
		resp := postBye(t, srv, "http://evil.example", "")
		resp.Body.Close()
		assertNoExit(t, cancelled)
	})

	t.Run("wrong Host", func(t *testing.T) {
		srv, cancelled, _ := autoExitServer(t, dir)
		resp := postBye(t, srv, "", "other-host:1234")
		resp.Body.Close()
		assertNoExit(t, cancelled)
	})

	t.Run("matching Origin still exits", func(t *testing.T) {
		srv, cancelled, _ := autoExitServer(t, dir)
		resp := postBye(t, srv, srv.URL, "")
		resp.Body.Close()
		select {
		case <-cancelled:
		case <-time.After(2 * time.Second):
			t.Fatal("a same-origin bye should start the exit timer")
		}
	})
}

// monitorTrue reports whether body's monitor script sets the sessionStorage
// marker unconditionally, i.e. the Monitor bool rendered as true. html/template
// pads the substituted literal with spaces in a script context, so the check
// collapses whitespace before matching.
func monitorTrue(body string) bool {
	collapsed := strings.Join(strings.Fields(body), " ")
	return strings.Contains(collapsed, `if ( true ) sessionStorage.setItem("llcppg-doc"`)
}

// TestAutoExitReloadBeaconGuard verifies the client-side wiring that keeps a
// live reload from being mistaken for a close (issue #1074): the auto-exit page
// must guard the bye beacon on window.__llcppgReloading, and live.js must set
// that flag before every location.reload(). Without both halves an automatic
// update would beacon bye and start the server's exit timer.
func TestAutoExitReloadBeaconGuard(t *testing.T) {
	dir := writeDir(t, map[string]string{"a.go": "package foo\n\nfunc Bar() {}\n"})
	// Live reload and auto-exit on together is the default `llcppg -doc` setup
	// and the only configuration where the two scripts interact.
	b := newBroker("inst-guard")
	srv := httptest.NewServer(newHandler(dir, false, b, nil, true, "", nil, nil))
	defer srv.Close()

	body := getBody(t, srv.URL+"/", http.StatusOK)

	// The page must pull in both the live-reload client and the auto-exit script.
	if !strings.Contains(body, "/static/live.js") {
		t.Error("auto-exit + live-reload page should include the live-reload client")
	}
	collapsed := strings.Join(strings.Fields(body), " ")
	if !strings.Contains(collapsed, "if (window.__llcppgReloading) return;") {
		t.Error("pagehide handler must skip the bye beacon while a live reload is in progress")
	}

	// live.js must set the flag before reloading so the guard above has something
	// to check.
	js := getBody(t, srv.URL+"/static/live.js", http.StatusOK)
	collapsedJS := strings.Join(strings.Fields(js), " ")
	if !strings.Contains(collapsedJS, "window.__llcppgReloading = true;") {
		t.Error("live.js must set window.__llcppgReloading before location.reload()")
	}
	if !strings.Contains(collapsedJS, "location.reload();") {
		t.Error("live.js should still reload the page on an update")
	}
}

// assertNoExit fails if cancelled fires within a window comfortably longer than
// the shortened exitGrace.
func assertNoExit(t *testing.T, cancelled <-chan struct{}) {
	t.Helper()
	select {
	case <-cancelled:
		t.Fatal("rejected bye should not start the exit timer")
	case <-time.After(exitGrace + 100*time.Millisecond):
	}
}

// TestRunAutoExitEndToEnd drives the full Run path with auto-exit on: the first
// page is monitored, and a bye beacon makes Run return after the grace period.
func TestRunAutoExitEndToEnd(t *testing.T) {
	orig := exitGrace
	exitGrace = 50 * time.Millisecond
	defer func() { exitGrace = orig }()

	dir := writeDir(t, map[string]string{
		"a.go": "// Package foo is a fixture.\npackage foo\n\nfunc Bar() {}\n",
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var out safeBuffer
	done := make(chan error, 1)
	go func() {
		done <- Run(ctx, dir, Options{
			Addr:        "127.0.0.1:0",
			OpenBrowser: false,
			AutoExit:    true,
			Stdout:      &out,
			Stderr:      io.Discard,
		})
	}()

	baseURL := waitForURL(t, &out)
	if !strings.Contains(out.String(), "will exit when the opened page is closed") {
		t.Errorf("startup output should mention auto-exit; got %q", out.String())
	}

	// Claim the monitored page.
	body := getBody(t, baseURL+"/", http.StatusOK)
	if !monitorTrue(body) {
		t.Error("the first page served by Run should be monitored")
	}

	// Beacon bye with a same-origin request; Run must return after the grace.
	u, err := url.Parse(baseURL)
	if err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest(http.MethodPost, baseURL+byeURL, nil)
	req.Host = u.Host
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("Run returned %v, want nil after auto-exit", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Run did not auto-exit within 5s of bye")
	}
	if !strings.Contains(out.String(), "documentation page closed") {
		t.Errorf("Run should print the shutdown notice; got %q", out.String())
	}
}
