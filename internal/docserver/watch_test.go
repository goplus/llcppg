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
	"os"
	"path/filepath"
	"testing"
	"time"
)

// testQuiet is the quiet period used by the tests. It matches the production
// defaultQuiet (100ms): long enough that even a loaded CI runner delivers a
// burst's file-system events to the debouncer well within one window (macOS
// kqueue delivers a directory write via a re-scan whose latency is not under
// the test's control), short enough to keep the tests fast.
const testQuiet = 100 * time.Millisecond

// newTestWatcher starts a watcher on dir with the test quiet period and
// discards its stderr warnings. It stops when the returned cancel is called
// (also registered via t.Cleanup).
func newTestWatcher(t *testing.T, dir string) Watcher {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	w, err := newWatcher(ctx, dir, testQuiet, io.Discard)
	if err != nil {
		t.Fatalf("newWatcher: %v", err)
	}
	return w
}

// expectChange waits up to a generous multiple of the quiet period for exactly
// one change, then asserts no further change follows. It fails the test if no
// change arrives or a second one does.
func expectChange(t *testing.T, w Watcher) {
	t.Helper()
	select {
	case <-w.Changes():
	case <-time.After(2 * time.Second):
		t.Fatal("expected a change, got none")
	}
	// A single burst must produce exactly one change.
	select {
	case <-w.Changes():
		t.Fatal("expected exactly one change, got a second")
	case <-time.After(5 * testQuiet):
	}
}

// expectNoChange asserts that no change arrives within a window several times
// the quiet period.
func expectNoChange(t *testing.T, w Watcher) {
	t.Helper()
	select {
	case <-w.Changes():
		t.Fatal("expected no change, got one")
	case <-time.After(6 * testQuiet):
	}
}

func TestWatcherCreate(t *testing.T) {
	dir := t.TempDir()
	w := newTestWatcher(t, dir)
	writeFile(t, filepath.Join(dir, "a.go"), "package foo\n")
	expectChange(t, w)
}

func TestWatcherModify(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a.go")
	writeFile(t, path, "package foo\n")
	w := newTestWatcher(t, dir)
	writeFile(t, path, "package foo\n\nfunc Bar() {}\n")
	expectChange(t, w)
}

func TestWatcherDelete(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a.go")
	writeFile(t, path, "package foo\n")
	w := newTestWatcher(t, dir)
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	expectChange(t, w)
}

func TestWatcherRename(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "a.go")
	writeFile(t, src, "package foo\n")
	w := newTestWatcher(t, dir)
	if err := os.Rename(src, filepath.Join(dir, "b.go")); err != nil {
		t.Fatal(err)
	}
	expectChange(t, w)
}

// TestWatcherAtomicSave models an editor saving by writing a temporary file and
// renaming it over the target. Watching the directory (not the file) must still
// see the change as exactly one reload.
func TestWatcherAtomicSave(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "a.go")
	writeFile(t, target, "package foo\n")
	w := newTestWatcher(t, dir)

	tmp := filepath.Join(dir, ".a.go.tmp")
	writeFile(t, tmp, "package foo\n\nfunc Bar() {}\n")
	if err := os.Rename(tmp, target); err != nil {
		t.Fatal(err)
	}
	expectChange(t, w)
}

// TestWatcherBurst writes many files in quick succession. The debounce must
// collapse them into a single change delivered after the burst settles.
func TestWatcherBurst(t *testing.T) {
	dir := t.TempDir()
	w := newTestWatcher(t, dir)

	start := time.Now()
	for i := 0; i < 100; i++ {
		writeFile(t, filepath.Join(dir, "f"+itoa(i)+".go"), "package foo\n")
	}

	select {
	case <-w.Changes():
	case <-time.After(2 * time.Second):
		t.Fatal("expected one change after the burst, got none")
	}
	// The change must arrive only after the quiet period following the last
	// write, never before the writes even began settling.
	if elapsed := time.Since(start); elapsed < testQuiet {
		t.Errorf("change arrived after %v, want at least the quiet period %v", elapsed, testQuiet)
	}
	// And it must be a single change.
	select {
	case <-w.Changes():
		t.Fatal("burst produced more than one change")
	case <-time.After(5 * testQuiet):
	}
}

// TestWatcherMergesCloseWrites writes twice in quick succession, well within
// the quiet period; the two writes must merge into one change.
//
// The file is created before the watcher starts so both writes are plain
// modifications of an already-watched file. On macOS (kqueue) a create is
// delivered via a directory re-scan and the per-file watch is only registered
// afterwards, so a create-then-modify pair can surface as two events spaced
// further apart than the quiet period and split into two changes; starting from
// an existing file keeps this test about the debounce merging, not that race.
//
// The two writes are issued back-to-back with no intervening sleep: the merge
// only holds while both events reach the debouncer within one quiet window, and
// kqueue delivers a directory write via a re-scan whose latency the test cannot
// control, so spending none of the window on a deliberate sleep leaves the
// whole quiet period to absorb that delivery jitter.
func TestWatcherMergesCloseWrites(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a.go")
	writeFile(t, path, "package foo\n")
	w := newTestWatcher(t, dir)

	writeFile(t, path, "package foo\n\nfunc Bar() {}\n")
	writeFile(t, path, "package foo\n\nfunc Baz() {}\n")
	expectChange(t, w)
}

func TestWatcherIgnoresNonGo(t *testing.T) {
	dir := t.TempDir()
	w := newTestWatcher(t, dir)
	writeFile(t, filepath.Join(dir, "llcppg.cfg"), "{}\n")
	writeFile(t, filepath.Join(dir, "README.md"), "hi\n")
	expectNoChange(t, w)
}

// TestWatcherIgnoresSubdir confirms the watch is non-recursive: a .go file in a
// subdirectory created before the watcher starts must not trigger a change.
func TestWatcherIgnoresSubdir(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")
	if err := os.Mkdir(sub, 0755); err != nil {
		t.Fatal(err)
	}
	w := newTestWatcher(t, dir)
	writeFile(t, filepath.Join(sub, "a.go"), "package bar\n")
	expectNoChange(t, w)
}

func TestWatcherIgnoresChmod(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a.go")
	writeFile(t, path, "package foo\n")
	w := newTestWatcher(t, dir)
	if err := os.Chmod(path, 0600); err != nil {
		t.Fatal(err)
	}
	expectNoChange(t, w)
}

func TestNewWatcherErrorsOnMissingDir(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := newWatcher(ctx, filepath.Join(t.TempDir(), "nope"), testQuiet, io.Discard)
	if err == nil {
		t.Fatal("expected an error watching a non-existent directory")
	}
}

// TestWatcherStopsOnContextCancel confirms the Changes channel is closed when
// the context is cancelled, which is how the broker's run loop ends at
// shutdown.
func TestWatcherStopsOnContextCancel(t *testing.T) {
	dir := t.TempDir()
	ctx, cancel := context.WithCancel(context.Background())
	w, err := newWatcher(ctx, dir, testQuiet, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	cancel()
	select {
	case _, ok := <-w.Changes():
		if ok {
			// A spurious change is possible in theory; drain and wait for close.
			<-w.Changes()
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Changes channel was not closed after context cancel")
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

// itoa avoids a strconv import in a test file that otherwise needs none.
func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	var buf [20]byte
	pos := len(buf)
	for i > 0 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}
	return string(buf[pos:])
}
