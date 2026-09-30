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

package pputil // preprocessor utility

import (
	"bytes"
	"iter"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// -----------------------------------------------------------------------------

// FileEntry represents a file entry in a directory.
type FileEntry struct {
	Path string
	os.DirEntry
}

// ListFiles returns a sequence of files in the specified directory. If recursive is
// true, it includes files in subdirectories as well.
func ListFiles(dir string, recursive bool) iter.Seq2[FileEntry, error] {
	return func(yield func(FileEntry, error) bool) {
		e := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				if !recursive && path != dir {
					return filepath.SkipDir
				}
				return nil
			}
			if !yield(FileEntry{Path: path, DirEntry: d}, nil) {
				return filepath.SkipAll
			}
			return nil
		})
		if e != nil {
			yield(FileEntry{}, e)
		}
	}
}

// -----------------------------------------------------------------------------

// Include represents a C/C++ include directive.
type Include struct {
	Filename string
	Quote    bool
}

// Search searches for the include file in the specified directories. Note: `found` is
// the resolved path only when ok is true, and the raw include name otherwise.
func (p *Include) Search(workDir string, searchDirs []string) (found string, ok bool) {
	filename := p.Filename
	if p.Quote {
		if found, ok = fileFound(filename, workDir); ok {
			return
		}
	}
	for _, dir := range searchDirs {
		if found, ok = fileFound(filename, dir); ok {
			return
		}
	}
	return filename, false
}

func fileFound(file, dir string) (string, bool) {
	found := filepath.Join(dir, file)
	_, err := os.Stat(found)
	return found, err == nil
}

// -----------------------------------------------------------------------------

// LoadIncludes loads the include directives from the specified header file. It
// returns a sequence of Include objects.
//
// It parses the file in a single pass, correctly handling `#include` directives
// with spaces or tabs between `#` and `include`, and ignoring any `#include`
// occurrences found within comments or string/character literals.
func LoadIncludes(headerFile string) (includes iter.Seq[Include], err error) {
	b, err := os.ReadFile(headerFile)
	if err != nil {
		return
	}
	includes = func(yield func(Include) bool) {
		scanIncludes(b, yield)
	}
	return
}

var includeKw = []byte("include")

// scanIncludes scans b for `#include` directives, yielding each one. It skips
// line comments, block comments, and string/character literals so that
// `#include` tokens appearing inside them are not mistaken for directives.
//
// A directive is recognized only when `#` is the first non-blank, non-comment
// token on a line, optionally followed by spaces or tabs, then `include`, then
// the header name in `"..."` or `<...>`. Comments are treated as whitespace, so
// a same-line block comment before the directive (e.g. `/* c */ #include <x.h>`)
// does not suppress it.
func scanIncludes(b []byte, yield func(Include) bool) {
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

// ListIncludes lists the include files for the specified header file. It returns a
// sequence of include file paths and a boolean indicating whether the file was found.
// The include file path is the resolved path when ok is true, and the raw include
// filename otherwise.
func ListIncludes(headerFile string, includeDirs []string) (includeFiles iter.Seq2[string, bool], err error) {
	includes, err := LoadIncludes(headerFile)
	if err != nil {
		return
	}
	includeFiles = func(yield func(string, bool) bool) {
		headerDir := filepath.Dir(headerFile)
		for inc := range includes {
			includeFile, ok := inc.Search(headerDir, includeDirs)
			if !yield(includeFile, ok) {
				return
			}
		}
	}
	return
}

// -----------------------------------------------------------------------------

func calcHeaderDeps(headerFiles map[string]bool, headerDir string, includeDirs []string) {
	for headerFile := range headerFiles {
		includeFiles, err := ListIncludes(headerFile, includeDirs)
		if err != nil {
			log.Panicln("[FATAL] ListIncludes failed:", err)
		}
		for includeFile, ok := range includeFiles {
			if !ok {
				log.Println("[WARN] include file not found:", includeFile)
			}
			if strings.HasPrefix(includeFile, headerDir) {
				headerFiles[includeFile] = true
			}
		}
	}
}

func collectHeaders(headerDir string, recursive bool) (headerFiles map[string]bool, err error) {
	headerFiles = make(map[string]bool)
	for file, e := range ListFiles(headerDir, recursive) {
		if e != nil {
			return nil, e
		}
		name := file.Name()
		if name[0] != '_' && IsHeaderFile(name) {
			headerFiles[file.Path] = false
		}
	}
	return
}

// IsHeaderFile checks if the given file name has a header file extension.
func IsHeaderFile(name string) bool {
	ext := filepath.Ext(name)
	switch ext {
	case ".h", ".hpp", ".hh", ".hxx":
		return true
	default:
		return false
	}
}

// TopHeaders returns the top-level header files in the specified directory.
func TopHeaders(headerDir string, recursive, rel bool, includeDirs []string) (topHeaders []string, err error) {
	headerDir, err = filepath.Abs(headerDir)
	if err != nil {
		return
	}

	headerFiles, err := collectHeaders(headerDir, recursive)
	if err != nil {
		return
	}

	headerDir += string(os.PathSeparator)
	calcHeaderDeps(headerFiles, headerDir, includeDirs)
	topHeaders = make([]string, 0, len(headerFiles))
	for headerFile, included := range headerFiles {
		if !included {
			if rel {
				headerFile = headerFile[len(headerDir):]
			}
			topHeaders = append(topHeaders, headerFile)
		}
	}
	sort.Strings(topHeaders)
	return
}

// -----------------------------------------------------------------------------
