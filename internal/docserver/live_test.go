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
	"bufio"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestLiveReloadEndToEnd starts a real watcher and broker behind the handler,
// edits a source file, and verifies that a reload arrives on /_events and that
// the next GET / reflects the edit.
func TestLiveReloadEndToEnd(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "a.go")
	writeFile(t, src, "package foo\n\nfunc First() {}\n")

	w := newTestWatcher(t, dir)
	b := newBroker("inst-1")
	go b.run(w)

	srv := httptest.NewServer(newHandler(dir, false, b, nil))
	defer srv.Close()

	// The page must advertise live reload.
	body := getBody(t, srv.URL+"/", http.StatusOK)
	if !strings.Contains(body, "/static/live.js") {
		t.Error("live-reload page should include the client script")
	}

	resp, r := openStream(t, srv)
	defer resp.Body.Close()
	// Skip hello.
	readEventAsync(t, r, 2*time.Second)

	// Edit the source; a reload must fan out.
	writeFile(t, src, "package foo\n\nfunc Second() {}\n")
	ev, _, _ := readEventAsync(t, r, 2*time.Second)
	if ev != "reload" {
		t.Fatalf("event = %q, want reload", ev)
	}

	// A fresh GET / must now show the edited content.
	body = getBody(t, srv.URL+"/", http.StatusOK)
	if !strings.Contains(body, `id="Second"`) {
		t.Error("reloaded page should reflect the edited source (Second)")
	}
	if strings.Contains(body, `id="First"`) {
		t.Error("reloaded page should no longer show the removed First")
	}
}

// TestDegradedModeNoLiveReload confirms that when the server has no broker
// (watcher could not start), GET / omits the client script and /_events
// returns 404.
func TestDegradedModeNoLiveReload(t *testing.T) {
	dir := writeDir(t, map[string]string{"a.go": "package foo\n\nfunc Bar() {}\n"})
	srv := httptest.NewServer(newHandler(dir, false, nil, nil))
	defer srv.Close()

	body := getBody(t, srv.URL+"/", http.StatusOK)
	if strings.Contains(body, "/static/live.js") {
		t.Error("degraded page must not include the live-reload script")
	}
	if !strings.Contains(body, `id="Bar"`) {
		t.Error("degraded page should still render the package")
	}

	resp, err := http.Get(srv.URL + "/_events")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("GET /_events in degraded mode: status = %d, want 404", resp.StatusCode)
	}
}

// TestRestartReloadsOnDifferentID models a page reconnecting after the server
// restarted: a hello with a different instance id must tell the client to
// reload, while the same id must not. The server side only emits the id; the
// comparison lives in the client script, so this test asserts the contract the
// script relies on: the hello data equals the broker's id.
func TestRestartHelloCarriesInstanceID(t *testing.T) {
	b1 := newBroker("inst-A")
	srv1 := httptest.NewServer(http.HandlerFunc(b1.serveEvents(time.Hour, nil)))
	defer srv1.Close()

	resp1, r1 := openStream(t, srv1)
	defer resp1.Body.Close()
	ev, data, _ := readEventAsync(t, r1, 2*time.Second)
	if ev != "hello" || data != "inst-A" {
		t.Fatalf("first hello = (%q,%q), want (hello,inst-A)", ev, data)
	}

	// A "restarted" server with a different id.
	b2 := newBroker("inst-B")
	srv2 := httptest.NewServer(http.HandlerFunc(b2.serveEvents(time.Hour, nil)))
	defer srv2.Close()

	resp2, r2 := openStream(t, srv2)
	defer resp2.Body.Close()
	ev, data, _ = readEventAsync(t, r2, 2*time.Second)
	if ev != "hello" || data != "inst-B" {
		t.Fatalf("second hello = (%q,%q), want (hello,inst-B)", ev, data)
	}
	if data == "inst-A" {
		t.Error("restarted server must present a different instance id")
	}
}

// TestInstanceIDChangesAcrossCalls confirms instanceID produces a different
// value on a later call, which is what makes restart detection work.
func TestInstanceIDChangesAcrossCalls(t *testing.T) {
	a := instanceID()
	time.Sleep(2 * time.Millisecond)
	bID := instanceID()
	if a == bID {
		t.Errorf("instanceID returned the same value twice: %q", a)
	}
}

// TestRunLiveReloadIntegration starts Run on a loopback port, connects to
// /_events, edits a source file, and verifies both a reload event and that the
// startup output announced live reload.
func TestRunLiveReloadIntegration(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "a.go")
	writeFile(t, src, "// Package foo is a fixture.\npackage foo\n\nfunc First() {}\n")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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

	url := waitForURL(t, &out)

	// The startup output should announce live reload.
	if !strings.Contains(out.String(), "live reload on") {
		t.Errorf("startup output should mention live reload; got: %q", out.String())
	}

	// Connect to the event stream and skip hello.
	resp, err := http.Get(url + "/_events")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET /_events: status = %d, want 200", resp.StatusCode)
	}
	r := bufio.NewReader(resp.Body)
	readEventAsync(t, r, 2*time.Second)

	// Edit the source and expect a reload.
	writeFile(t, src, "// Package foo is a fixture.\npackage foo\n\nfunc Second() {}\n")
	ev, _, _ := readEventAsync(t, r, 3*time.Second)
	if ev != "reload" {
		t.Fatalf("event = %q, want reload", ev)
	}

	// A fresh GET / must reflect the edit.
	body := getBody(t, url+"/", http.StatusOK)
	if !strings.Contains(body, `id="Second"`) {
		t.Error("page after reload should reflect the edited source")
	}

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
