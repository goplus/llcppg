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

package docserver

import (
	"os"
	"path/filepath"
	"testing"
)

// writeDir creates a temporary directory populated with the given files
// (relative name -> contents) and returns its path.
func writeDir(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, content := range files {
		path := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		files   map[string]string
		wantErr bool
		// check runs when load succeeds.
		check func(t *testing.T, p *pkg)
	}{
		{
			name: "simple package",
			files: map[string]string{
				"a.go": "// Package foo does things.\npackage foo\n\n// A is a constant.\nconst A = 1\n",
			},
			check: func(t *testing.T, p *pkg) {
				if p.Doc.Name != "foo" {
					t.Errorf("name = %q, want foo", p.Doc.Name)
				}
			},
		},
		{
			name: "no go files",
			files: map[string]string{
				"readme.txt": "not go",
			},
			wantErr: true,
		},
		{
			name: "multiple packages",
			files: map[string]string{
				"a.go": "package foo\n",
				"b.go": "package bar\n",
			},
			wantErr: true,
		},
		{
			name: "test files are ignored",
			files: map[string]string{
				"a.go":      "package foo\n\nfunc F() {}\n",
				"a_test.go": "package foo\n\nimport \"testing\"\n\nfunc TestF(t *testing.T) {}\n",
			},
			check: func(t *testing.T, p *pkg) {
				if p.Doc.Name != "foo" {
					t.Errorf("name = %q, want foo", p.Doc.Name)
				}
			},
		},
		{
			name: "parse error is non-fatal",
			files: map[string]string{
				"good.go": "package foo\n\nfunc Good() {}\n",
				"bad.go":  "package foo\n\nfunc Bad( {\n", // syntax error
			},
			check: func(t *testing.T, p *pkg) {
				if len(p.ParseErrors) == 0 {
					t.Error("want parse errors recorded, got none")
				}
				// Good declaration should still be present.
				if len(p.Doc.Funcs) == 0 {
					t.Error("want Good func rendered despite the parse error")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := writeDir(t, tt.files)
			p, err := load(dir, false)
			if tt.wantErr {
				if err == nil {
					t.Fatal("want error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("load: %v", err)
			}
			if tt.check != nil {
				tt.check(t, p)
			}
		})
	}
}

func TestLoadImportPath(t *testing.T) {
	// go.mod at the directory root: import path is the module path.
	t.Run("module root", func(t *testing.T) {
		dir := writeDir(t, map[string]string{
			"go.mod": "module example.com/mypkg\n\ngo 1.27\n",
			"a.go":   "package mypkg\n",
		})
		p, err := load(dir, false)
		if err != nil {
			t.Fatal(err)
		}
		if p.ImportPath != "example.com/mypkg" {
			t.Errorf("import path = %q, want example.com/mypkg", p.ImportPath)
		}
	})

	// go.mod in a parent: import path is module path + relative subdir.
	t.Run("subdirectory", func(t *testing.T) {
		root := t.TempDir()
		if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/lib\n"), 0644); err != nil {
			t.Fatal(err)
		}
		sub := filepath.Join(root, "c", "sqlite")
		if err := os.MkdirAll(sub, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(sub, "a.go"), []byte("package sqlite\n"), 0644); err != nil {
			t.Fatal(err)
		}
		p, err := load(sub, false)
		if err != nil {
			t.Fatal(err)
		}
		if p.ImportPath != "example.com/lib/c/sqlite" {
			t.Errorf("import path = %q, want example.com/lib/c/sqlite", p.ImportPath)
		}
	})

	// No go.mod anywhere: fall back to the directory base name.
	t.Run("no go.mod", func(t *testing.T) {
		dir := writeDir(t, map[string]string{
			"a.go": "package foo\n",
		})
		p, err := load(dir, false)
		if err != nil {
			t.Fatal(err)
		}
		if p.ImportPath != filepath.Base(dir) {
			t.Errorf("import path = %q, want %q", p.ImportPath, filepath.Base(dir))
		}
	})
}

func TestModulePath(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"module example.com/foo\n", "example.com/foo"},
		{"// comment\nmodule example.com/bar\n\ngo 1.27\n", "example.com/bar"},
		{"module example.com/baz // inline comment\n", "example.com/baz"},
		{"module \"example.com/quoted\"\n", "example.com/quoted"},
		{"go 1.27\n", ""},
	}
	for _, tt := range tests {
		if got := modulePath([]byte(tt.in)); got != tt.want {
			t.Errorf("modulePath(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestAllDecls(t *testing.T) {
	files := map[string]string{
		"a.go": "package foo\n\n// Exported is public.\nfunc Exported() {}\n\n// unexported is private.\nfunc unexported() {}\n",
	}
	dir := writeDir(t, files)

	// Without AllDecls, only the exported func appears.
	p, err := load(dir, false)
	if err != nil {
		t.Fatal(err)
	}
	if n := len(p.Doc.Funcs); n != 1 {
		t.Errorf("without AllDecls: %d funcs, want 1", n)
	}

	// With AllDecls, the unexported func appears too.
	p, err = load(dir, true)
	if err != nil {
		t.Fatal(err)
	}
	if n := len(p.Doc.Funcs); n != 2 {
		t.Errorf("with AllDecls: %d funcs, want 2", n)
	}
}
