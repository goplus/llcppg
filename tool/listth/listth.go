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

package listth

import (
	"bytes"
	"iter"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// -----------------------------------------------------------------------------

// Include represents a C/C++ include directive.
type Include struct {
	Filename string
	Quote    bool
}

// Search searches for the include file in the specified directories.
func (p *Include) Search(workDir string, searchDirs []string) (found string, ok bool) {
	if p.Quote {
		if found, ok = fileFound(p.Filename, workDir); ok {
			return
		}
	}
	for _, dir := range searchDirs {
		if found, ok = fileFound(p.Filename, dir); ok {
			return
		}
	}
	return
}

func fileFound(file, dir string) (string, bool) {
	found := filepath.Join(dir, file)
	_, err := os.Stat(found)
	return found, err == nil
}

// -----------------------------------------------------------------------------

var include = []byte("#include")

// LoadIncludes loads the include directives from the specified header file.
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

func listIncludeFiles(headerFile string, includeDirs []string) (includeFiles iter.Seq[string], err error) {
	includes, err := LoadIncludes(headerFile)
	if err != nil {
		return
	}
	includeFiles = func(yield func(string) bool) {
		headerDir := filepath.Dir(headerFile)
		for inc := range includes {
			includeFile, ok := inc.Search(headerDir, includeDirs)
			if !ok {
				log.Printf("[WARN] include file not found: %s", inc.Filename)
				continue
			}
			if !yield(includeFile) {
				return
			}
		}
	}
	return
}

func calcHeaderDeps(headerFiles map[string]bool, headerDir string, includeDirs []string) {
	for headerFile, included := range headerFiles {
		if included {
			continue
		}
		includeFiles, err := listIncludeFiles(headerFile, includeDirs)
		if err != nil {
			log.Panicln("[FATAL] searchIncludeFiles failed:", err)
		}
		for includeFile := range includeFiles {
			if strings.HasPrefix(includeFile, headerDir) {
				headerFiles[includeFile] = true
			}
		}
	}
}

func collectHeaders(headerDir string, recursive bool) (headerFiles map[string]bool, err error) {
	if recursive {
		panic("recursive not implemented")
	}
	fis, err := os.ReadDir(headerDir)
	if err != nil {
		return
	}
	headerFiles = make(map[string]bool, len(fis))
	headerDir += string(os.PathSeparator)
	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}
		name := fi.Name()
		headerFile := filepath.Join(headerDir, name)
		headerFiles[headerFile] = false
	}
	return
}

// TopHeaders returns the top-level header files in the specified directory.
func TopHeaders(headerDir string, recursive bool, includeDirs []string) (topHeaders []string, err error) {
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
			topHeaders = append(topHeaders, headerFile[len(headerDir):])
		}
	}
	return
}

// -----------------------------------------------------------------------------
