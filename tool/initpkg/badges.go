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

package initpkg

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// -----------------------------------------------------------------------------

// badgeOrg is the GitHub organization that owns the binding repositories. It is
// fixed text in the badge URLs, matching the template repository's organization.
const badgeOrg = "llarhub"

// badgeLines returns the four standard badge lines for module, in order: GoDoc,
// release, LLGo, and XGo. The module name is filled into the GoDoc and release
// badges (two uses each); the LLGo and XGo badges are fixed text.
func badgeLines(module string) []string {
	return []string{
		fmt.Sprintf("[![GoDoc](https://pkg.go.dev/badge/github.com/%s/%s.svg)](https://pkg.go.dev/github.com/%s/%s)", badgeOrg, module, badgeOrg, module),
		fmt.Sprintf("[![GitHub release](https://img.shields.io/github/v/tag/%s/%s.svg?label=release)](https://github.com/%s/%s/releases)", badgeOrg, module, badgeOrg, module),
		"[![LLGo](https://img.shields.io/badge/powered_by-LLGo-green.svg)](https://github.com/xgo-dev/llgo)",
		"[![XGo](https://img.shields.io/badge/project-XGo-blue.svg)](https://github.com/goplus/xgo)",
	}
}

// badgeMarker is the substring whose presence means the badges are already in
// the README (for example because the template already carries them), so the
// block must not be inserted a second time.
func badgeMarker(module string) string {
	return fmt.Sprintf("https://pkg.go.dev/badge/github.com/%s/%s.svg", badgeOrg, module)
}

// addReadmeBadges inserts the four standard badges after the title of the
// README.md at the root of dir, filling in module. It is part of the main
// branch phase of Init and runs after the template files are in place and
// before the commit, so the change lands in the same commit as the rest of the
// main branch.
//
// The step is cosmetic and never fails the command: it only edits README.md,
// and any problem (a missing title, an unreadable or unwritable file) results in
// a warning on errOut and an unchanged README rather than an error. A missing
// README.md is not even a warning; the step simply does nothing. When the README
// already carries the GoDoc badge for module, the block is not added again.
func addReadmeBadges(dir, module string, errOut io.Writer) {
	path := filepath.Join(dir, "README.md")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// No README on main: nothing to do, and no warning.
			return
		}
		fmt.Fprintf(errOut, "llcppg: warning: could not read %q to add the standard badges (%v); add them by hand if you want them\n", "README.md", err)
		return
	}
	content := string(data)

	// Already present (e.g. the template carries them): do not add twice.
	if strings.Contains(content, badgeMarker(module)) {
		return
	}

	updated, ok := insertBadges(content, module)
	if !ok {
		fmt.Fprintf(errOut, "llcppg: warning: README.md has no recognizable title; skipped adding the standard badges. Add them by hand if you want them.\n")
		return
	}

	info, err := os.Stat(path)
	mode := os.FileMode(0644)
	if err == nil {
		mode = info.Mode().Perm()
	}
	if err := os.WriteFile(path, []byte(updated), mode); err != nil {
		fmt.Fprintf(errOut, "llcppg: warning: could not write %q to add the standard badges (%v); add them by hand if you want them\n", "README.md", err)
	}
}

// -----------------------------------------------------------------------------
// README editing as text (not as a Markdown tree): the step only finds one line
// and inserts text after it, so the rest of the README stays byte-for-byte what
// the template provided.

// readmeLine is one physical line of the README: its text without the line
// terminator, and the terminator itself ("\r\n", "\n", or "" for a last line
// that has no trailing newline). Keeping the terminator per line lets the
// reconstruction preserve the file's existing line endings and trailing-newline
// state exactly.
type readmeLine struct {
	text string
	term string
}

// insertBadges inserts the badge block after the README's title and returns the
// new content. ok is false when no title is found, in which case content must be
// left unchanged. The module name fills the GoDoc and release badges.
func insertBadges(content, module string) (string, bool) {
	lines := splitLines(content)
	titleIdx, ok := findTitleLine(lines)
	if !ok {
		return "", false
	}

	// The inserted lines follow the file's existing convention: CRLF if the file
	// uses it anywhere, LF otherwise.
	nl := "\n"
	if strings.Contains(content, "\r\n") {
		nl = "\r\n"
	}
	trailingNewline := len(lines) > 0 && lines[len(lines)-1].term != ""

	var b strings.Builder
	// Everything up to and including the title line, with original terminators.
	// If the title was the very last line without a newline, it now needs one
	// because the badges follow it.
	for idx := 0; idx <= titleIdx; idx++ {
		b.WriteString(lines[idx].text)
		term := lines[idx].term
		if idx == titleIdx && term == "" {
			term = nl
		}
		b.WriteString(term)
	}

	// Reuse the blank lines that already followed the title (collapsing any run
	// of them to a single blank line) rather than doubling them.
	k := titleIdx + 1
	for k < len(lines) && strings.TrimSpace(lines[k].text) == "" {
		k++
	}
	rest := lines[k:]

	// One blank line, the four badges, then one blank line before whatever
	// followed the title (if anything did).
	b.WriteString(nl)
	badges := badgeLines(module)
	for idx, bl := range badges {
		b.WriteString(bl)
		if idx < len(badges)-1 {
			b.WriteString(nl)
			continue
		}
		// Last badge line.
		if len(rest) > 0 {
			b.WriteString(nl) // end the badge line
			b.WriteString(nl) // the blank line before the rest
		} else if trailingNewline {
			// Nothing follows: match the file's original trailing-newline state.
			b.WriteString(nl)
		}
	}
	for _, ln := range rest {
		b.WriteString(ln.text)
		b.WriteString(ln.term)
	}
	return b.String(), true
}

// splitLines breaks content into physical lines, recording each line's terminator
// so the exact bytes (line endings and a possibly-absent final newline) can be
// reproduced.
func splitLines(content string) []readmeLine {
	var lines []readmeLine
	for i := 0; i < len(content); {
		j := strings.IndexByte(content[i:], '\n')
		if j < 0 {
			lines = append(lines, readmeLine{text: content[i:], term: ""})
			break
		}
		end := i + j
		if end > i && content[end-1] == '\r' {
			lines = append(lines, readmeLine{text: content[i : end-1], term: "\r\n"})
		} else {
			lines = append(lines, readmeLine{text: content[i:end], term: "\n"})
		}
		i = end + 1
	}
	return lines
}

// findTitleLine returns the index of the line after which the badge block is
// inserted: the line of the first level-1 heading (ATX "# Title" style) or the
// "===" underline of the first setext-style level-1 heading, whichever comes
// first. It returns ok false when no level-1 heading is found.
//
// Content before the title is skipped over: a YAML front matter block at the
// very top (delimited by "---" lines), HTML comments, and any other text are
// allowed before the title. Fenced code blocks (``` or ~~~) are tracked so a
// "# comment" line inside one is not mistaken for a heading, and a setext
// underline inside one is not mistaken for a title. Only the two Markdown
// heading forms are recognized; a title written as raw HTML (for example an
// <h1> tag) is not detected.
func findTitleLine(lines []readmeLine) (int, bool) {
	i := 0
	// Skip a YAML front matter block only when it is the very first line.
	if len(lines) > 0 && strings.TrimRight(lines[0].text, " \t") == "---" {
		for j := 1; j < len(lines); j++ {
			if strings.TrimRight(lines[j].text, " \t") == "---" {
				i = j + 1
				break
			}
		}
	}

	inFence := false
	var fenceMarker byte // '`' or '~'
	for ; i < len(lines); i++ {
		text := lines[i].text
		if marker, ok := fenceDelimiter(text); ok {
			if !inFence {
				inFence = true
				fenceMarker = marker
			} else if marker == fenceMarker {
				inFence = false
			}
			continue
		}
		if inFence {
			continue
		}
		if isATXLevel1(text) {
			return i, true
		}
		// A setext level-1 heading: a non-blank, non-heading line immediately
		// followed by a line of only '=' characters. The title text may span
		// several physical lines before the underline; the badge block still goes
		// after the underline. The underline cannot be the first line of content.
		if isSetextUnderline(text, '=') && i > 0 && hasSetextTitleText(lines, i) {
			return i, true
		}
	}
	return 0, false
}

// hasSetextTitleText reports whether the line(s) immediately before the setext
// underline at underlineIdx form valid heading text: at least one non-blank line
// that is not itself a fence delimiter or an ATX heading.
func hasSetextTitleText(lines []readmeLine, underlineIdx int) bool {
	prev := lines[underlineIdx-1].text
	if strings.TrimSpace(prev) == "" {
		return false
	}
	if _, ok := fenceDelimiter(prev); ok {
		return false
	}
	if isATXLevel1(prev) {
		return false
	}
	return true
}

// isATXLevel1 reports whether text is an ATX level-1 heading: up to three
// leading spaces, a single '#', then either end-of-line or a space (an optional
// closing run of '#' is allowed, as in "# Title #"). "##" or deeper is not a
// level-1 heading.
func isATXLevel1(text string) bool {
	s := text
	spaces := 0
	for len(s) > 0 && s[0] == ' ' && spaces < 3 {
		s = s[1:]
		spaces++
	}
	if len(s) == 0 || s[0] != '#' {
		return false
	}
	s = s[1:]
	// A second '#' makes it level-2 or deeper.
	if len(s) > 0 && s[0] == '#' {
		return false
	}
	// After the single '#', the line must be empty or start with a space/tab.
	if len(s) == 0 {
		return true
	}
	return s[0] == ' ' || s[0] == '\t'
}

// isSetextUnderline reports whether text is a setext underline made only of the
// given character c ('=' for level 1, '-' for level 2), with up to three leading
// spaces and optional trailing spaces, and at least one c.
func isSetextUnderline(text string, c byte) bool {
	s := strings.TrimRight(text, " \t")
	spaces := 0
	for len(s) > 0 && s[0] == ' ' && spaces < 3 {
		s = s[1:]
		spaces++
	}
	if len(s) == 0 {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] != c {
			return false
		}
	}
	return true
}

// fenceDelimiter reports whether text opens or closes a fenced code block and,
// if so, which marker it uses ('`' or '~'). A fence is at least three identical
// marker characters with up to three leading spaces; an info string after the
// run (as in "```go") is allowed on an opening fence.
func fenceDelimiter(text string) (byte, bool) {
	s := text
	spaces := 0
	for len(s) > 0 && s[0] == ' ' && spaces < 3 {
		s = s[1:]
		spaces++
	}
	if len(s) < 3 {
		return 0, false
	}
	marker := s[0]
	if marker != '`' && marker != '~' {
		return 0, false
	}
	run := 0
	for run < len(s) && s[run] == marker {
		run++
	}
	if run < 3 {
		return 0, false
	}
	return marker, true
}

// -----------------------------------------------------------------------------
