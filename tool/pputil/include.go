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

package pputil

import (
	"bytes"
	"iter"
)

// -----------------------------------------------------------------------------

// ScanIncludes scans the given content for `#include` directives and returns a
// sequence of Include objects.
//
// It parses the file in a single pass, correctly handling `#include` directives
// with spaces or tabs between `#` and `include`, and ignoring any `#include`
// occurrences found within comments or string/character literals.
func ScanIncludes(content []byte) iter.Seq[Include] {
	return func(yield func(Include) bool) {
		doScanIncludes(content, yield)
	}
}

var includeKw = []byte("include")

// doScanIncludes scans b for `#include` directives, yielding each one. It skips
// line comments, block comments, and string/character literals so that
// `#include` tokens appearing inside them are not mistaken for directives.
//
// A directive is recognized only when `#` is the first non-blank, non-comment
// token on a line, optionally followed by spaces or tabs, then `include`, then
// the header name in `"..."` or `<...>`. Comments are treated as whitespace, so
// a same-line block comment before the directive (e.g. `/* c */ #include <x.h>`)
// does not suppress it.
func doScanIncludes(b []byte, yield func(Include) bool) {
	atLineStart := true // no non-blank character seen yet on the current line
	for i := 0; i < len(b); {
		c := b[i]
		switch {
		case c == '\n':
			atLineStart = true
			i++
		case c == ' ' || c == '\t' || c == '\r' || c == '\f' || c == '\v':
			i++ // leading blanks keep us at line start
		case c == '/' && i+1 < len(b) && b[i+1] == '/':
			// line comment: skip to end of line
			i += 2
			for i < len(b) && b[i] != '\n' {
				i++
			}
		case c == '/' && i+1 < len(b) && b[i+1] == '*':
			// block comment: skip to closing */. Comments count as
			// whitespace, so atLineStart is left unchanged.
			i += 2
			for i < len(b) && !(b[i] == '*' && i+1 < len(b) && b[i+1] == '/') {
				i++
			}
			i += 2
		case c == '"' || c == '\'':
			i = skipLiteral(b, i, c)
			atLineStart = false
		case c == '#' && atLineStart:
			if next, inc, ok := parseInclude(b, i); ok {
				if !yield(inc) {
					return
				}
				i = next
			} else {
				i++
			}
			atLineStart = false
		default:
			atLineStart = false
			i++
		}
	}
}

// skipLiteral returns the index just past a string (") or character (') literal
// that starts at b[i] == quote, honoring backslash escapes.
func skipLiteral(b []byte, i int, quote byte) int {
	i++ // opening quote
	for i < len(b) {
		switch b[i] {
		case '\\':
			i += 2 // skip the escaped character
		case quote:
			return i + 1
		case '\n':
			return i // unterminated literal; stop at end of line
		default:
			i++
		}
	}
	return i
}

// parseInclude tries to parse an `#include "..."` or `#include <...>` directive
// starting at b[i] == '#'. On success it returns the index just past the
// directive, the parsed Include, and true.
func parseInclude(b []byte, i int) (next int, inc Include, ok bool) {
	j := i + 1 // past '#'
	for j < len(b) && (b[j] == ' ' || b[j] == '\t') {
		j++
	}
	if !bytes.HasPrefix(b[j:], includeKw) {
		return
	}
	j += len(includeKw)
	// require a blank between `include` and the header name
	start := j
	for j < len(b) && (b[j] == ' ' || b[j] == '\t') {
		j++
	}
	if j == start || j >= len(b) {
		return
	}
	open := b[j]
	var closer byte
	switch open {
	case '"':
		closer = '"'
	case '<':
		closer = '>'
	default:
		return
	}
	j++
	end := j
	for end < len(b) && b[end] != closer && b[end] != '\n' {
		end++
	}
	if end >= len(b) || b[end] != closer {
		return
	}
	return end + 1, Include{Filename: string(b[j:end]), Quote: open == '"'}, true
}

// -----------------------------------------------------------------------------
