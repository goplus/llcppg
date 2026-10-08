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

// FilterFunc defines a function type used to filter files in a directory.
type FilterFunc func(FileEntry) bool

// ListFilterFiles returns a sequence of files in the specified directory that
// satisfy the filter function. If recursive is true, it includes files in
// subdirectories as well.
func ListFilterFiles(dir string, recursive bool, filter FilterFunc) iter.Seq2[FileEntry, error] {
	return func(yield func(FileEntry, error) bool) {
		for file, e := range ListFiles(dir, recursive) {
			if e != nil || filter(file) {
				if !yield(file, e) {
					return
				}
			}
		}
	}
}

// FilterPublicHeaderFile is a filter function that returns true if the given
// file is a public header file (i.e., it does not start with an underscore and
// has a header file extension).
func FilterPublicHeaderFile(file FileEntry) bool {
	name := file.Name()
	if name[0] != '_' {
		return IsHeaderFile(file.Path)
	}
	return false
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

// ListIncludes lists the include files for the specified header file. It returns a
// sequence of include file paths and a boolean indicating whether the file was found.
// The include file path is the resolved path when ok is true, and the raw include
// filename otherwise.
func ListIncludes(headerFile string, includeDirs []string) (includeFiles iter.Seq2[string, bool], err error) {
	b, err := os.ReadFile(headerFile)
	if err != nil {
		return
	}
	includeFiles = func(yield func(string, bool) bool) {
		headerDir := filepath.Dir(headerFile)
		for inc := range ScanIncludes(b) {
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
	for file, e := range ListFilterFiles(headerDir, recursive, FilterPublicHeaderFile) {
		if e != nil {
			return nil, e
		}
		headerFiles[file.Path] = false
	}
	return
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
