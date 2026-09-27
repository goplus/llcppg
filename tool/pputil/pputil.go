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

// Search searches for the include file in the specified directories.
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

var include = []byte("#include")

// LoadIncludes loads the include directives from the specified header file. It
// returns a sequence of Include objects.
func LoadIncludes(headerFile string) (includes iter.Seq[Include], err error) {
	b, err := os.ReadFile(headerFile)
	if err != nil {
		return
	}
	includes = func(yield func(Include) bool) {
		for {
			pos := bytes.Index(b, include)
			if pos < 0 {
				break
			}
			b = b[pos+len(include):]
			b = bytes.TrimLeft(b, " \t")
			if len(b) == 0 {
				break
			}
			quote := b[0]
			if quote != '"' {
				if quote != '<' {
					continue
				}
				quote = '>'
			}
			b = b[1:]
			pos = bytes.IndexByte(b, quote)
			if pos < 0 {
				continue
			}
			fname := string(b[:pos])
			b = b[pos+1:]
			if !yield(Include{Filename: fname, Quote: quote == '"'}) {
				return
			}
		}
	}
	return
}

// -----------------------------------------------------------------------------

// ListIncludes lists the include files for the specified header file. It returns a
// sequence of include file paths and a boolean indicating whether the file was found.
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
