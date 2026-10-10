// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package derrors provides the error-wrapping helper used by the vendored
// godoc rendering code.
//
// It is a trimmed copy of github.com/goplus/llcppg/internal/godoc/internal/derrors containing
// only Wrap, which is all the rendering path needs. The full pkgsite package
// also carries HTTP status mapping and crash reporting that llcppg's local
// documentation server does not use.
package derrors

import "fmt"

// Wrap adds context to the error and allows unwrapping the result to recover
// the original error.
//
// Example:
//
//	defer derrors.Wrap(&err, "copy(%s, %s)", src, dst)
//
// See Add for an equivalent function that does not allow
// the result to be unwrapped.
func Wrap(errp *error, format string, args ...any) {
	if *errp != nil {
		*errp = fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), *errp)
	}
}
