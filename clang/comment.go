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

package clang

import (
	"github.com/goplus/lib/c"
	"github.com/llarhub/clang-c"
)

// -----------------------------------------------------------------------------

// RawComment returns the raw documentation comment associated with the given
// cursor, including the comment markers (e.g. "/** ... */" or "///"). It returns
// an empty string if the cursor has no associated comment.
func RawComment(entity clang.Cursor) string {
	raw := RawCommentsOf(entity)
	return raw.Text()
}

// -----------------------------------------------------------------------------

// RawComments represents the raw documentation comments associated with a declaration,
// including preceding, trailing, and inline comments.
type RawComments struct {
	// A preceding comment counts as the declaration's leading doc only when it sits
	// on the line(s) directly above the declaration with no blank line in between;
	// a blank line detaches it. The string keeps the original comment markers so it
	// can be fed to the same comment-normalizing logic as RawComment. It can be ""
	// when there is no such comment.
	Preceding string `json:"preceding,omitempty"`

	// A trailing comment is considered part of the declaration's documentation
	// only when it appears on the line immediately below the declaration with no
	// blank line in between.
	Trailing string `json:"trailing,omitempty"`

	// An inline comment is considered part of the declaration's documentation
	// only when it appears on the same line as the declaration.
	Inline string `json:"inline,omitempty"`
}

func (raw *RawComments) Text() string {
	if raw.Preceding != "" {
		return raw.Preceding
	}
	if raw.Inline != "" {
		return raw.Inline
	}
	return raw.Trailing
}

// RawCommentsOf returns the raw documentation comments associated with the given
// cursor, including the preceding, trailing, and inline comments. Each comment
// retains its original markers (e.g., "/** ... */" or "///"). If a particular
// type of comment is not present, the corresponding field in the returned
// RawComments struct will be an empty string.
func RawCommentsOf(entity clang.Cursor) (ret RawComments) {
	var file File
	var offset c.Uint
	rg := entity.Extent()
	rg.Start().Expansion(&file, nil, nil, &offset)
	if file == InvalidFile {
		return
	}
	src := TU(entity).FileContents(file)
	// A macro cursor's extent starts at the macro name, i.e. after the
	// "#define " keyword, so a comment on the preceding line is not directly
	// behind it. Back up to the start of the macro's own line first, then scan
	// for an adjacent comment above it.
	ret.Preceding = precedingRawComment(src, lineStart(src, int(offset)))
	var end c.Uint
	rg.End().Expansion(nil, nil, nil, &end)
	isTrailing, raw := followingRawComment(src, int(end))
	if isTrailing == 0 {
		ret.Inline = raw
	} else {
		ret.Trailing = raw
	}
	return
}

func precedingRawComment(src []byte, start int) string {
	if start > len(src) {
		start = len(src)
	}
	// Walk back over the inline whitespace that sits between the comment and
	// the declaration. Stop at the first non-space byte or the newline that
	// ends the comment's own line; a second newline means a blank line gap,
	// which detaches the comment.
	i := start
	newlines := 0
	for i > 0 {
		c := src[i-1]
		switch c {
		case ' ', '\t', '\r':
			i--
			continue
		case '\n':
			if newlines++; newlines > 1 {
				return ""
			}
			i--
			continue
		}
		break
	}
	if i == 0 {
		return ""
	}

	// A block comment ends with "*/".
	if i >= 2 && src[i-1] == '/' && src[i-2] == '*' {
		for j := i - 2; j >= 1; j-- {
			if src[j-1] == '/' && src[j] == '*' {
				if blank(src, j-1) {
					return string(src[j-1 : i])
				}
				break
			}
		}
		return ""
	}

	// Otherwise collect a run of adjacent single-line comments ("//", "///",
	// "//!"). The scan is already positioned at the end of the line holding the
	// nearest comment (newlines <= 1 above guarantees adjacency).
	lineEnd := i
	var first int = -1
	for {
		ls := lineStart(src, lineEnd)
		content := trimLeadingSpace(src[ls:lineEnd])
		if !hasPrefix(content, "//") {
			break
		}
		first = ls
		if ls == 0 {
			break
		}
		lineEnd = ls - 1 // step onto the newline that ends the previous line
		if lineEnd > 0 && src[lineEnd-1] == '\r' {
			lineEnd--
		}
	}
	if first < 0 {
		return ""
	}
	return string(trimLeadingSpace(src[first:i]))
}

func followingRawComment(src []byte, end int) (newlines int, raw string) {
	if end < 0 {
		return
	}
	// Walk forward over the remainder of the declaration's line, then over at
	// most one newline onto the line below. A second newline means a blank line
	// gap, which detaches any comment.
	i := end
	for i < len(src) {
		switch src[i] {
		case ' ', '\t', '\r', ';':
			i++
			continue
		case '\n':
			if newlines++; newlines > 1 {
				return
			}
			i++
			continue
		}
		break
	}
	if i >= len(src) {
		return
	}

	// A block comment starts with "/*".
	if i+1 < len(src) && src[i] == '/' && src[i+1] == '*' {
		for j := i + 2; j+1 < len(src); j++ {
			if src[j] == '*' && src[j+1] == '/' {
				raw = string(src[i : j+2])
				return
			}
		}
		return
	}

	// Otherwise collect a run of adjacent single-line comments ("//", "///",
	// "//!"). Each must begin a line with only whitespace before it.
	if !hasPrefix(src[i:], "//") {
		return
	}
	first := i
	last := i
	for i < len(src) {
		le := lineEndAt(src, i)
		if !hasPrefix(trimLeadingSpace(src[i:le]), "//") {
			break
		}
		last = le
		// Advance to the start of the next line: skip the trimmed "\r" (if any)
		// and the "\n". If there is no newline, we are at EOF.
		nl := le
		for nl < len(src) && src[nl] != '\n' {
			nl++
		}
		if nl >= len(src) {
			break
		}
		i = nl + 1
	}
	raw = string(src[first:last])
	return
}

// lineEndAt returns the index of the newline (or len(src)) that terminates the
// line containing pos, trimming a trailing "\r" so CRLF endings are excluded.
func lineEndAt(src []byte, pos int) int {
	for pos < len(src) && src[pos] != '\n' {
		pos++
	}
	if pos > 0 && src[pos-1] == '\r' {
		pos--
	}
	return pos
}

// lineStart returns the index of the first byte of the line containing pos.
func lineStart(src []byte, pos int) int {
	for pos > 0 && src[pos-1] != '\n' {
		pos--
	}
	return pos
}

func blank(src []byte, pos int) bool {
	for pos > 0 {
		switch src[pos-1] {
		case '\n':
			return true
		case ' ', '\t':
			pos--
		default:
			return false
		}
	}
	return true
}

// trimLeadingSpace drops leading spaces and tabs from b without copying.
func trimLeadingSpace(b []byte) []byte {
	i := 0
	for i < len(b) && (b[i] == ' ' || b[i] == '\t') {
		i++
	}
	return b[i:]
}

// hasPrefix reports whether b starts with the ASCII prefix p.
func hasPrefix(b []byte, p string) bool {
	if len(b) < len(p) {
		return false
	}
	for i := 0; i < len(p); i++ {
		if b[i] != p[i] {
			return false
		}
	}
	return true
}

// -----------------------------------------------------------------------------
