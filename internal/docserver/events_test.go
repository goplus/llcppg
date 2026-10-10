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
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// fakeWatcher drives the broker from a test without real file events: a value
// sent on its channel is one "reload now".
type fakeWatcher struct {
	ch chan struct{}
}

func newFakeWatcher() *fakeWatcher              { return &fakeWatcher{ch: make(chan struct{}, 1)} }
func (w *fakeWatcher) Changes() <-chan struct{} { return w.ch }
func (w *fakeWatcher) fire()                    { w.ch <- struct{}{} }

// openStream opens an SSE connection to the handler and returns the response
// and a buffered reader positioned to read event lines. It fails if the status
// is not 200.
func openStream(t *testing.T, srv *httptest.Server) (*http.Response, *bufio.Reader) {
	t.Helper()
	resp, err := http.Get(srv.URL + "/_events")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		t.Fatalf("GET /_events: status = %d, want 200", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "text/event-stream") {
		resp.Body.Close()
		t.Fatalf("Content-Type = %q, want text/event-stream", ct)
	}
	return resp, bufio.NewReader(resp.Body)
}

// readEvent reads lines until a blank line, returning the "event:" and "data:"
// values of one SSE message (either may be empty, e.g. a keep-alive comment is
// returned with event=="" and the comment preserved on the colon-prefixed line
// handled by the caller). It fails on timeout via the surrounding read deadline
// semantics of the httptest connection, so callers wrap it with a goroutine +
// select when a timeout matters.
func readEvent(r *bufio.Reader) (event, data, comment string, err error) {
	for {
		line, e := r.ReadString('\n')
		if e != nil {
			return event, data, comment, e
		}
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			return event, data, comment, nil
		}
		switch {
		case strings.HasPrefix(line, ":"):
			comment = strings.TrimSpace(line[1:])
		case strings.HasPrefix(line, "event:"):
			event = strings.TrimSpace(line[len("event:"):])
		case strings.HasPrefix(line, "data:"):
			data = strings.TrimSpace(line[len("data:"):])
		}
	}
}

// readEventAsync runs readEvent with a timeout so a hung stream fails the test
// rather than blocking it forever.
func readEventAsync(t *testing.T, r *bufio.Reader, timeout time.Duration) (event, data, comment string) {
	t.Helper()
	type res struct {
		event, data, comment string
		err                  error
	}
	done := make(chan res, 1)
	go func() {
		e, d, c, err := readEvent(r)
		done <- res{e, d, c, err}
	}()
	select {
	case r := <-done:
		if r.err != nil {
			t.Fatalf("reading SSE event: %v", r.err)
		}
		return r.event, r.data, r.comment
	case <-time.After(timeout):
		t.Fatal("timed out reading SSE event")
		return "", "", ""
	}
}

func TestEventsHelloThenReload(t *testing.T) {
	w := newFakeWatcher()
	b := newBroker("inst-1")
	go b.run(w)

	srv := httptest.NewServer(http.HandlerFunc(b.serveEvents(time.Hour, nil)))
	defer srv.Close()

	resp, r := openStream(t, srv)
	defer resp.Body.Close()

	// First message is hello with the instance id.
	ev, data, _ := readEventAsync(t, r, 2*time.Second)
	if ev != "hello" {
		t.Errorf("first event = %q, want hello", ev)
	}
	if data != "inst-1" {
		t.Errorf("hello data = %q, want inst-1", data)
	}

	// A change fans out a reload.
	w.fire()
	ev, data, _ = readEventAsync(t, r, 2*time.Second)
	if ev != "reload" {
		t.Errorf("event = %q, want reload", ev)
	}
}

func TestEventsKeepAlive(t *testing.T) {
	w := newFakeWatcher()
	b := newBroker("inst-1")
	go b.run(w)

	// A tiny keep-alive period so the comment arrives promptly.
	srv := httptest.NewServer(http.HandlerFunc(b.serveEvents(20*time.Millisecond, nil)))
	defer srv.Close()

	resp, r := openStream(t, srv)
	defer resp.Body.Close()

	// Skip hello.
	readEventAsync(t, r, 2*time.Second)

	// With no change fired, the next message must be a keep-alive comment.
	_, _, comment := readEventAsync(t, r, 2*time.Second)
	if !strings.Contains(comment, "keep-alive") {
		t.Errorf("expected a keep-alive comment, got comment=%q", comment)
	}
}

// TestBrokerDoesNotBlockOnStalledStream confirms a subscriber that never reads
// cannot block notify: the size-one channel drops the second pending reload.
func TestBrokerDoesNotBlockOnStalledStream(t *testing.T) {
	b := newBroker("inst-1")
	ch, ok := b.subscribe()
	if !ok {
		t.Fatal("subscribe should succeed")
	}
	defer b.unsubscribe(ch)

	done := make(chan struct{})
	go func() {
		// Many notifies with nobody draining ch must not block.
		for i := 0; i < 1000; i++ {
			b.notify()
		}
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("notify blocked on a stalled stream")
	}
	// Exactly one reload is pending.
	select {
	case <-ch:
	default:
		t.Fatal("expected one pending reload")
	}
	select {
	case <-ch:
		t.Fatal("expected only one pending reload")
	default:
	}
}

// TestEventsStreamCap confirms that once maxStreams connections are open,
// further ones get 503.
func TestEventsStreamCap(t *testing.T) {
	b := newBroker("inst-1")
	srv := httptest.NewServer(http.HandlerFunc(b.serveEvents(time.Hour, nil)))
	defer srv.Close()

	// Open maxStreams long-lived connections. Each must succeed and hold a slot.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	opened := make([]*http.Response, 0, maxStreams)
	for i := 0; i < maxStreams; i++ {
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, srv.URL+"/_events", nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("stream %d: %v", i, err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("stream %d: status = %d, want 200", i, resp.StatusCode)
		}
		opened = append(opened, resp)
	}
	defer func() {
		for _, resp := range opened {
			resp.Body.Close()
		}
	}()

	// Read the hello off each so the handler is parked in its select, holding
	// the subscription, before we test the cap.
	for _, resp := range opened {
		r := bufio.NewReader(resp.Body)
		readEventAsync(t, r, 2*time.Second)
	}

	// One more connection must be rejected with 503.
	resp, err := http.Get(srv.URL + "/_events")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("over-cap stream: status = %d, want 503", resp.StatusCode)
	}
}
