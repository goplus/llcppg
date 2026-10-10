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
	"strings"
	"testing"
)

// fixtureSrc is a small package that exercises constants, variables, a type
// with a constructor and a method, an undocumented symbol, and a comment with
// HTML-special characters.
const fixtureSrc = `// Package sqlite is a tiny fixture package.
package sqlite

// Version is the library version. The value is < 4 && > 2.
const Version = 3

// DefaultName is a package-level variable.
var DefaultName = "db"

// DB is a database handle.
type DB struct{}

// Open opens a DB. It returns a *DB.
func Open(name string) *DB { return &DB{} }

// Close closes the DB.
func (d *DB) Close() error { return nil }

func Query() {}
`

func renderFixture(t *testing.T, src string, allDecls bool) *pageData {
	t.Helper()
	dir := writeDir(t, map[string]string{"sqlite.go": src})
	p, err := load(dir, allDecls)
	if err != nil {
		t.Fatal(err)
	}
	return p.render()
}

func TestRenderGroupsAndConstructors(t *testing.T) {
	data := renderFixture(t, fixtureSrc, false)

	if data.Name != "sqlite" {
		t.Errorf("name = %q, want sqlite", data.Name)
	}
	if len(data.Consts) != 1 {
		t.Errorf("consts = %d, want 1", len(data.Consts))
	}
	if len(data.Vars) != 1 {
		t.Errorf("vars = %d, want 1", len(data.Vars))
	}

	// DB type should carry Open as a constructor and Close as a method; Query
	// is a free function (not attached to a type).
	var db *typeDoc
	for _, td := range data.Types {
		if td.Name == "DB" {
			db = td
		}
	}
	if db == nil {
		t.Fatal("type DB not found")
	}
	if len(db.Funcs) != 1 || db.Funcs[0].Name != "Open" {
		t.Errorf("DB constructors = %v, want [Open]", funcNames(db.Funcs))
	}
	if len(db.Methods) != 1 || db.Methods[0].Name != "Close" {
		t.Errorf("DB methods = %v, want [Close]", funcNames(db.Methods))
	}
	if len(data.Funcs) != 1 || data.Funcs[0].Name != "Query" {
		t.Errorf("free funcs = %v, want [Query]", funcNames(data.Funcs))
	}
}

func TestRenderEscapesHTML(t *testing.T) {
	data := renderFixture(t, fixtureSrc, false)
	// The Version doc comment contains "< 4 && > 2"; it must be escaped, not
	// rendered as raw markup.
	overviewAndConsts := string(data.Overview)
	for _, c := range data.Consts {
		overviewAndConsts += string(c.Doc)
	}
	if strings.Contains(overviewAndConsts, "< 4 && > 2") {
		t.Error("HTML-special characters were not escaped in the doc comment")
	}
	if !strings.Contains(overviewAndConsts, "&lt; 4 &amp;&amp; &gt; 2") {
		t.Errorf("expected escaped comment text, got: %q", overviewAndConsts)
	}
}

func TestRenderTypeLinking(t *testing.T) {
	data := renderFixture(t, fixtureSrc, false)
	var db *typeDoc
	for _, td := range data.Types {
		if td.Name == "DB" {
			db = td
		}
	}
	if db == nil {
		t.Fatal("type DB not found")
	}
	// Open's signature returns *DB, which should be linked to the DB anchor.
	decl := string(db.Funcs[0].Decl)
	if !strings.Contains(decl, `href="#DB"`) {
		t.Errorf("Open decl should link DB; got: %q", decl)
	}
}

func TestRenderSummary(t *testing.T) {
	data := renderFixture(t, fixtureSrc, false)
	s := data.Summary
	if s.Types != 1 {
		t.Errorf("summary types = %d, want 1", s.Types)
	}
	// Open and Query are exported funcs; Close is a method.
	if s.Funcs != 2 {
		t.Errorf("summary funcs = %d, want 2 (Open, Query)", s.Funcs)
	}
	if s.Methods != 1 {
		t.Errorf("summary methods = %d, want 1 (Close)", s.Methods)
	}
	// Query has no doc comment -> exactly one undocumented export.
	if s.Undocumented != 1 {
		t.Errorf("summary undocumented = %d, want 1 (Query)", s.Undocumented)
	}
}

func TestRenderNoPackageComment(t *testing.T) {
	data := renderFixture(t, "package nocomment\n\nfunc F() {}\n", false)
	if data.Overview != "" {
		t.Errorf("expected empty overview for a package with no comment, got %q", data.Overview)
	}
}

func funcNames(fs []*funcDoc) []string {
	var out []string
	for _, f := range fs {
		out = append(out, f.Name)
	}
	return out
}
