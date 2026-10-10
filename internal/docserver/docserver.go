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

// Package docserver implements `llcppg -doc`: it loads the Go package in a
// directory, renders its documentation to HTML with the vendored
// github.com/goplus/llcppg/internal/godoc package (a trimmed copy of
// golang.org/x/pkgsite's godoc renderer), and serves it on the loopback
// interface so the author can read a freshly generated binding the way a
// consumer would.
//
// It is deliberately not named "doc" to avoid confusion with the standard
// library. The package loading and the HTTP server use only the Go standard
// library; the HTML rendering is delegated to internal/godoc so the output
// matches the familiar pkg.go.dev layout. Unlike the cl/tool packages it has no
// cgo/libclang dependency, so it can be built and tested with plain `go`.
package docserver

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

// Options configures Run.
type Options struct {
	// Addr is the TCP listen address. Empty means DefaultAddr
	// ("127.0.0.1:0"), which lets the OS pick a free loopback port.
	Addr string

	// OpenBrowser launches the default browser at the served URL once the
	// listener is ready. A failed launch is reported as a warning only.
	OpenBrowser bool

	// AllDecls includes unexported declarations (doc.AllDecls). This helps
	// when debugging generated bindings that carry compiler-facing helper
	// types and padding fields.
	AllDecls bool

	// Stdout receives the one-line startup message ("serving documentation
	// ... at http://..."). Empty means os.Stdout.
	Stdout io.Writer

	// Stderr receives warnings (for example a non-loopback bind or a failed
	// browser launch). Empty means os.Stderr.
	Stderr io.Writer
}

// DefaultAddr is the default listen address: a loopback host and an
// OS-assigned port.
const DefaultAddr = "127.0.0.1:0"

func (o *Options) addr() string {
	if o.Addr != "" {
		return o.Addr
	}
	return DefaultAddr
}

func (o *Options) stdout() io.Writer {
	if o.Stdout != nil {
		return o.Stdout
	}
	return os.Stdout
}

func (o *Options) stderr() io.Writer {
	if o.Stderr != nil {
		return o.Stderr
	}
	return os.Stderr
}

// Run loads the package in dir and serves its documentation until ctx is done.
//
// It validates up front that dir contains exactly one importable package for
// the current build context; a directory with no Go files or with more than one
// package is a usage error reported before the server starts. Once that initial
// load succeeds the listener is opened, the URL is printed, the browser is
// optionally launched, and each request to "/" re-reads the directory so the
// "tweak config → regenerate → refresh" loop stays tight. Later parse errors do
// not stop the server: the page shows them and renders whatever parsed.
//
// Run blocks until ctx is cancelled (typically on SIGINT) and then shuts the
// server down gracefully. It only ever reads dir (and go.mod files above it);
// it writes nothing and never executes package code.
func Run(ctx context.Context, dir string, opts Options) error {
	// Load once up front so usage errors (no Go files, multiple packages) are
	// reported before we bind a port. We keep the result only to report the
	// import path in the startup message; each request reloads independently so
	// edits are picked up without a restart.
	p, err := load(dir, opts.AllDecls)
	if err != nil {
		return err
	}

	addr := opts.addr()
	if host := hostOf(addr); host != "" && !isLoopback(host) {
		fmt.Fprintf(opts.stderr(), "llcppg: warning: %s is not a loopback address; the documentation server will be reachable from other hosts\n", host)
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("llcppg: cannot listen on %s: %w", addr, err)
	}

	url := "http://" + ln.Addr().String()

	// Start the file watcher for live reload. If it cannot be created (for
	// example the OS inotify limit is reached), we fall back to a degraded
	// server with no live reload: b stays nil, the page carries no client
	// script and /_events is not registered. The documentation itself never
	// depends on the watcher.
	var b *broker
	w, werr := NewWatcher(ctx, p.Dir, defaultQuiet)
	if werr != nil {
		fmt.Fprintf(opts.stderr(), "llcppg: warning: live reload disabled: %v\n", werr)
	} else {
		b = newBroker(instanceID())
		go b.run(w)
	}

	handler := newHandler(dir, opts.AllDecls, b, ctx.Done())

	// ReadHeaderTimeout bounds how long a client may take to send request
	// headers. It is harmless on the loopback default and prevents a trivial
	// slowloris when an explicit non-loopback address is used.
	srv := &http.Server{Handler: handler, ReadHeaderTimeout: 10 * time.Second}

	// Report the URL before launching the browser so the author always has it,
	// even if the launch fails or there is no browser (SSH, container, CI).
	fmt.Fprintf(opts.stdout(), "llcppg: serving documentation for %s at %s\n", p.ImportPath, url)
	if b != nil {
		fmt.Fprintf(opts.stdout(), "llcppg: watching *.go in %s (live reload on)\n", p.Dir)
	}

	if opts.OpenBrowser {
		if err := openBrowser(url); err != nil {
			fmt.Fprintf(opts.stderr(), "llcppg: warning: could not open a browser (%v); open %s manually\n", err, url)
		}
	}

	// Serve in the background so we can watch ctx for a shutdown signal.
	serveErr := make(chan error, 1)
	go func() {
		serveErr <- srv.Serve(ln)
	}()

	select {
	case <-ctx.Done():
		// Graceful shutdown; a background context keeps Shutdown from being
		// cancelled by the same signal that triggered it. Open event streams
		// are closed through the done channel (ctx.Done(), passed to
		// serveEvents) rather than their request contexts, which Shutdown does
		// not cancel; the watcher goroutine stops when ctx is done. So shutdown
		// does not wait on idle tabs.
		return srv.Shutdown(context.Background())
	case err := <-serveErr:
		if err == http.ErrServerClosed {
			return nil
		}
		return err
	}
}

// instanceID returns a short identifier that is stable for the life of this
// process and (with very high probability) differs from a previous run on the
// same address. It is sent in the SSE hello message so a reconnecting page can
// tell a brief network blip (same id, do nothing) from a server restart
// (different id, reload once). Start time plus the process id is enough: it
// does not need to be unpredictable, only different across restarts.
func instanceID() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), os.Getpid())
}

// hostOf returns the host part of a "host:port" address, or the whole string
// when it has no port. An empty host (as in ":0") returns "".
func hostOf(addr string) string {
	if host, _, err := net.SplitHostPort(addr); err == nil {
		return host
	}
	return addr
}

// isLoopback reports whether host names the loopback interface, either by the
// usual literal names or by a loopback IP literal.
func isLoopback(host string) bool {
	switch strings.ToLower(host) {
	case "localhost", "":
		return true
	}
	if ip := net.ParseIP(host); ip != nil {
		return ip.IsLoopback()
	}
	return false
}
