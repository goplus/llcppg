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
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watcher reports that the .go files in a directory have changed. Each value
// received on Changes means "reload now"; bursts of file-system events are
// already coalesced by the quiet-period debounce, so the consumer never needs
// to debounce again.
type Watcher interface {
	Changes() <-chan struct{}
}

// defaultQuiet is the trailing-edge debounce window used by Run. A regeneration
// writes many files over a short span; waiting for a quiet period after the
// last event yields one reload on a consistent state rather than many reloads
// on half-written packages. It is a parameter of NewWatcher only so tests can
// shorten it; there is no command-line flag.
const defaultQuiet = 100 * time.Millisecond

// fsWatcher is the fsnotify-backed Watcher. It watches the directory (not the
// individual files) so that atomic saves (write-temp-then-rename) and
// delete/recreate during regeneration keep working: a watch on a single file
// would be lost when the file is replaced, while a directory watch sees the
// create/rename of the new file.
type fsWatcher struct {
	changes chan struct{}
}

func (w *fsWatcher) Changes() <-chan struct{} { return w.changes }

// NewWatcher watches the .go files directly inside dir (non-recursively). A
// change is announced on the returned Watcher's Changes channel once no further
// relevant event has arrived for the quiet period.
//
// It returns an error only when the underlying watcher cannot be created or dir
// cannot be added; the caller treats that as "no live reload" and keeps
// serving. Once started, the watcher runs until ctx is done, at which point it
// closes the fsnotify watcher and the Changes channel.
func NewWatcher(ctx context.Context, dir string, quiet time.Duration) (Watcher, error) {
	return newWatcher(ctx, dir, quiet, os.Stderr)
}

// newWatcher is NewWatcher with an explicit log sink so tests can capture or
// discard the stderr warnings.
func newWatcher(ctx context.Context, dir string, quiet time.Duration, logw io.Writer) (Watcher, error) {
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("create file watcher: %w", err)
	}
	if err := fsw.Add(dir); err != nil {
		fsw.Close()
		return nil, fmt.Errorf("watch %s: %w", dir, err)
	}

	w := &fsWatcher{changes: make(chan struct{}, 1)}
	go w.loop(ctx, fsw, quiet, logw)
	return w, nil
}

// loop drains fsnotify's channels, applies the .go-name and operation filters,
// and runs the trailing-edge debounce. It owns the fsnotify watcher and closes
// it (and the Changes channel) when ctx is done.
func (w *fsWatcher) loop(ctx context.Context, fsw *fsnotify.Watcher, quiet time.Duration, logw io.Writer) {
	defer close(w.changes)
	defer fsw.Close()

	// A single reusable timer implements the quiet period. It is created
	// stopped; a relevant event (re)starts it, and only its firing announces a
	// change. Using one timer keeps the "last event of a burst wins" semantics
	// without allocating per event.
	timer := time.NewTimer(quiet)
	if !timer.Stop() {
		<-timer.C
	}
	armed := false
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case ev, ok := <-fsw.Events:
			if !ok {
				return
			}
			if !relevant(ev) {
				continue
			}
			// Trailing edge: (re)start the quiet period on every relevant
			// event, so a burst collapses into the single firing that follows
			// the last event.
			if armed && !timer.Stop() {
				// Timer already fired but we have not consumed it yet; drain so
				// the reset below starts a fresh window.
				select {
				case <-timer.C:
				default:
				}
			}
			timer.Reset(quiet)
			armed = true

		case err, ok := <-fsw.Errors:
			if !ok {
				return
			}
			// An event-queue overflow means some events were lost, so we cannot
			// know whether a .go file changed. Treat it as a change: it is
			// better to reload once too often than to miss an update.
			if errors.Is(err, fsnotify.ErrEventOverflow) {
				timer.Reset(quiet)
				armed = true
				continue
			}
			// Other errors (including the watched directory being removed) do
			// not stop the server; log and keep going so manual refresh still
			// works.
			fmt.Fprintf(logw, "llcppg: warning: file watcher error: %v\n", err)

		case <-timer.C:
			armed = false
			announce(w.changes)
		}
	}
}

// relevant reports whether a file-system event should feed the debouncer: the
// entry's name must end in ".go" and the operation must include a create,
// write, remove or rename. A chmod-only event changes nothing the page shows
// and is ignored.
func relevant(ev fsnotify.Event) bool {
	if !strings.HasSuffix(filepath.Base(ev.Name), ".go") {
		return false
	}
	return ev.Op&(fsnotify.Create|fsnotify.Write|fsnotify.Remove|fsnotify.Rename) != 0
}

// announce performs a non-blocking send on ch. The channel is buffered to one;
// if a reload is already pending, a second is dropped, because one pending
// reload is all a consumer needs. This keeps the debounce goroutine from ever
// blocking on a slow consumer.
func announce(ch chan<- struct{}) {
	select {
	case ch <- struct{}{}:
	default:
	}
}
