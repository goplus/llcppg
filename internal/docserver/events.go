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
	"net/http"
	"sync"
	"time"
)

// maxStreams caps how many Server-Sent Events connections the broker holds at
// once. The traffic per stream is tiny, but the cap stops a runaway script (or
// a page in a reload loop) from accumulating goroutines without bound; further
// connections get 503 and the browser's EventSource retries later.
const maxStreams = 32

// keepAlivePeriod is how often the handler writes an SSE comment line. It keeps
// proxies and idle-connection timeouts from closing an otherwise silent stream.
// It is a field on the handler only so tests can shorten it.
const keepAlivePeriod = 15 * time.Second

// broker owns the set of connected SSE streams and fans a reload out to all of
// them. The watcher calls notify once per settled burst; subscribers each hold
// a buffered channel of size one, so a slow or stalled browser tab can never
// block the broker: if its channel already holds a pending reload, a second is
// dropped, which is harmless because one reload fully refreshes the page.
type broker struct {
	// id identifies this server instance. It is sent in the hello message so a
	// reconnecting page can tell "same server, brief blip" (do nothing) from
	// "server restarted" (reload once). It is chosen at start-up and never
	// changes for the life of the process.
	id string

	mu   sync.Mutex
	subs map[chan struct{}]struct{}
}

// newBroker creates an empty broker with the given instance id.
func newBroker(id string) *broker {
	return &broker{id: id, subs: make(map[chan struct{}]struct{})}
}

// subscribe registers a new stream and returns its channel together with the
// current subscriber count after registration. The caller must call
// unsubscribe when the stream ends. The returned ok is false when the stream
// cap is already reached, in which case nothing is registered.
func (b *broker) subscribe() (ch chan struct{}, ok bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.subs) >= maxStreams {
		return nil, false
	}
	ch = make(chan struct{}, 1)
	b.subs[ch] = struct{}{}
	return ch, true
}

// unsubscribe removes a stream registered by subscribe. It is safe to call once
// per subscribe, including after the broker has fanned out reloads.
func (b *broker) unsubscribe(ch chan struct{}) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.subs, ch)
}

// notify sends a reload to every connected stream. Delivery is non-blocking per
// stream (see broker doc), so a single stalled consumer cannot delay the others
// or the watcher goroutine that calls this.
func (b *broker) notify() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for ch := range b.subs {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}

// run forwards every change from w into a fan-out notify. It returns when the
// Changes channel is closed (the watcher's context was cancelled), so it ends
// as part of normal shutdown.
func (b *broker) run(w Watcher) {
	for range w.Changes() {
		b.notify()
	}
}

// serveEvents is the "GET /_events" handler. Each connected page holds one such
// request open. It writes a hello message carrying the instance id, then a
// reload message whenever the broker fans one out, and a keep-alive comment
// every keepAlive period. It returns (closing the stream) when the client goes
// away, the request context is cancelled, or done is closed on server shutdown.
//
// done is the server's lifetime signal: http.Server.Shutdown does not cancel
// in-flight request contexts, so without this an open SSE stream would keep
// Shutdown blocking until the browser disconnected. Closing done lets every
// stream return promptly so shutdown does not wait on idle tabs.
func (b *broker) serveEvents(keepAlive time.Duration, done <-chan struct{}) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			// Without streaming support we cannot deliver events incrementally.
			http.Error(w, "streaming unsupported", http.StatusInternalServerError)
			return
		}

		ch, ok := b.subscribe()
		if !ok {
			// Too many open streams. The browser's EventSource will retry, so a
			// transient 503 simply delays this page's live reload rather than
			// breaking it.
			http.Error(w, "too many event streams", http.StatusServiceUnavailable)
			return
		}
		defer b.unsubscribe(ch)

		h := w.Header()
		h.Set("Content-Type", "text/event-stream")
		h.Set("Cache-Control", "no-cache")
		h.Set("Connection", "keep-alive")
		w.WriteHeader(http.StatusOK)

		// hello carries the instance id so a reconnecting page can detect a
		// restart. Flush it immediately so EventSource.onopen fires and the page
		// marks itself "live".
		fmt.Fprintf(w, "event: hello\ndata: %s\n\n", b.id)
		flusher.Flush()

		ticker := time.NewTicker(keepAlive)
		defer ticker.Stop()

		ctx := req.Context()
		for {
			select {
			case <-ctx.Done():
				return
			case <-done:
				return
			case <-ch:
				fmt.Fprint(w, "event: reload\ndata: 1\n\n")
				flusher.Flush()
			case <-ticker.C:
				// A comment line (": ...") is ignored by EventSource but keeps
				// the connection and any proxies between from timing out.
				fmt.Fprint(w, ": keep-alive\n\n")
				flusher.Flush()
			}
		}
	}
}
