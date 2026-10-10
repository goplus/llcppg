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

func TestRenderShellMetadata(t *testing.T) {
	data := renderFixture(t, fixtureSrc, false)

	if data.Name != "sqlite" {
		t.Errorf("name = %q, want sqlite", data.Name)
	}
	if data.Body == "" {
		t.Fatal("rendered body is empty")
	}
}

func TestRenderBodyContainsDeclarations(t *testing.T) {
	data := renderFixture(t, fixtureSrc, false)
	body := string(data.Body)

	// godoc's package template renders an index and the exported declarations.
	for _, want := range []string{"Version", "DefaultName", "DB", "Open", "Close", "Query"} {
		if !strings.Contains(body, want) {
			t.Errorf("rendered body missing %q", want)
		}
	}
}

func TestRenderEscapesHTML(t *testing.T) {
	data := renderFixture(t, fixtureSrc, false)
	body := string(data.Body)
	// The Version doc comment contains "< 4 && > 2"; godoc must escape it, not
	// emit it as raw markup.
	if strings.Contains(body, "< 4 && > 2") {
		t.Error("HTML-special characters were not escaped in the doc comment")
	}
	if !strings.Contains(body, "&lt; 4 &amp;&amp; &gt; 2") {
		t.Error("expected single-escaped comment text in body")
	}
	// The godoc markup itself (e.g. <p> tags) must be preserved, not escaped;
	// double-escaping would turn it into literal &lt;p&gt; text.
	if strings.Contains(body, "&lt;p&gt;") {
		t.Error("godoc markup was double-escaped")
	}
}

func TestRenderValidScriptLiterals(t *testing.T) {
	data := renderFixture(t, fixtureSrc, false)
	body := string(data.Body)
	// AnalysisData / CallGraph must be emitted as valid JS (null), never an
	// empty assignment that would be a syntax error.
	if strings.Contains(body, "document.ANALYSIS_DATA = ;") {
		t.Error("ANALYSIS_DATA rendered as an invalid empty assignment")
	}
	if !strings.Contains(body, "document.ANALYSIS_DATA = null;") {
		t.Error("ANALYSIS_DATA should be rendered as null")
	}
}

func TestRenderNoPackageComment(t *testing.T) {
	data := renderFixture(t, "package nocomment\n\nfunc F() {}\n", false)
	if data.Name != "nocomment" {
		t.Errorf("name = %q, want nocomment", data.Name)
	}
	if !strings.Contains(string(data.Body), "F") {
		t.Error("body should still render func F for a package with no comment")
	}
}
