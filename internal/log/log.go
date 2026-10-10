// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package log provides the two logging helpers the vendored godoc rendering
// code calls.
//
// It is a trimmed replacement for github.com/goplus/llcppg/internal/godoc/internal/log. pkgsite's
// version talks to Stackdriver and an experiment framework; llcppg only needs a
// place for the renderer's occasional warning or error line to go, so these
// write to the standard library logger. The context argument is accepted for
// source compatibility with the pkgsite call sites and is otherwise unused.
package log

import (
	"context"
	stdlog "log"
)

// Warningf logs a formatted warning.
func Warningf(ctx context.Context, format string, args ...any) {
	stdlog.Printf("WARNING: "+format, args...)
}

// Errorf logs a formatted error.
func Errorf(ctx context.Context, format string, args ...any) {
	stdlog.Printf("ERROR: "+format, args...)
}
