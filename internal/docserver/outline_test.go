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
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

// outlineAnchorRE matches an href="#anchor" in the rendered outline tree.
var outlineAnchorRE = regexp.MustCompile(`href="#([^"]+)"`)

// TestOutlineBodyConsistency guards the contract that every outline anchor
// (other than the page-level section-* wrappers the shell adds) has a matching
// id in the body. Both come from the same TemplateData, so they must agree; a
// future template change that breaks this would be caught here.
func TestOutlineBodyConsistency(t *testing.T) {
	data := renderFixture(t, fixtureSrc, false)
	outline := string(data.Outline)
	body := string(data.Body)

	matches := outlineAnchorRE.FindAllStringSubmatch(outline, -1)
	if len(matches) == 0 {
		t.Fatalf("outline contains no anchors; outline:\n%s", outline)
	}
	for _, m := range matches {
		anchor := m[1]
		if strings.HasPrefix(anchor, "section-") {
			// section-documentation / section-sourcefiles are page-shell
			// anchors, not body ids.
			continue
		}
		if !strings.Contains(body, `id="`+anchor+`"`) {
			t.Errorf("outline anchor #%s has no matching id in the body", anchor)
		}
	}
}

// TestOutlineHasExpectedEntries checks the outline carries the standard
// pkg.go.dev sections for the fixture (Overview, Index, Functions, Types) and
// the per-symbol entries, so the sidebar tree is actually populated.
func TestOutlineHasExpectedEntries(t *testing.T) {
	data := renderFixture(t, fixtureSrc, false)
	outline := string(data.Outline)
	for _, want := range []string{
		`href="#pkg-overview"`,
		`href="#pkg-index"`,
		`href="#pkg-functions"`,
		`href="#pkg-types"`,
		`href="#Open"`,     // constructor under the DB type
		`href="#DB.Close"`, // method under the DB type
	} {
		if !strings.Contains(outline, want) {
			t.Errorf("outline missing %q; outline:\n%s", want, outline)
		}
	}
}

// TestPageHasOutlineTree checks the full page shell wires the outline into a
// tree the sidebar script can enhance, and pulls in outline.js.
func TestPageHasOutlineTree(t *testing.T) {
	dir := writeDir(t, map[string]string{"sqlite.go": fixtureSrc})
	p, err := load(dir, false)
	if err != nil {
		t.Fatal(err)
	}
	page := renderPage(t, p)

	for _, want := range []string{
		`class="go-Tree js-tree"`,       // the tree the script enhances
		`href="#section-documentation"`, // top-level Documentation entry
		`href="#pkg-overview"`,          // the nested outline is present
		`<script src="/static/outline.js"></script>`,
	} {
		if !strings.Contains(page, want) {
			t.Errorf("page missing %q", want)
		}
	}
}

func TestSourceFilesHelper(t *testing.T) {
	got := sourceFiles([]string{"a.go", "b.go"})
	if len(got) != 2 {
		t.Fatalf("sourceFiles returned %d entries, want 2", len(got))
	}
	if got[0].Name != "a.go" || got[0].URL != "/src/a.go" {
		t.Errorf("entry 0 = %+v, want {a.go /src/a.go}", got[0])
	}
	if got[1].URL != "/src/b.go" {
		t.Errorf("entry 1 URL = %q, want /src/b.go", got[1].URL)
	}
	// An empty list yields an empty (non-nil-dependent) slice so the template's
	// {{if .SourceFiles}} guard hides the section.
	if len(sourceFiles(nil)) != 0 {
		t.Error("sourceFiles(nil) should be empty")
	}
}

// TestPageSourceFilesSection checks that, for a multi-file package, the page
// renders the Source Files sidebar entry and the bottom section listing every
// selected file (sorted) with a link to its /src/ view that opens in a new tab.
func TestPageSourceFilesSection(t *testing.T) {
	dir := writeDir(t, map[string]string{
		"zeta.go":  "// Package multi is a fixture.\npackage multi\n\nfunc Z() {}\n",
		"alpha.go": "package multi\n\nfunc A() {}\n",
	})
	p, err := load(dir, false)
	if err != nil {
		t.Fatal(err)
	}
	// load sorts the selected files; the Source Files list follows that order.
	if len(p.Files) != 2 || p.Files[0] != "alpha.go" || p.Files[1] != "zeta.go" {
		t.Fatalf("p.Files = %v, want [alpha.go zeta.go]", p.Files)
	}

	page := renderPage(t, p)
	for _, want := range []string{
		`href="#section-sourcefiles"`,             // sidebar entry
		`id="section-sourcefiles"`,                // bottom section
		`<a href="/src/alpha.go" target="_blank"`, // first file, new tab
		`<a href="/src/zeta.go" target="_blank"`,  // second file
	} {
		if !strings.Contains(page, want) {
			t.Errorf("page missing %q", want)
		}
	}
	// alpha.go must appear before zeta.go in the rendered list.
	if strings.Index(page, "/src/alpha.go") > strings.Index(page, "/src/zeta.go") {
		t.Error("source files should be listed in sorted order")
	}
}

// TestPageNoSourceFilesWhenEmpty verifies the Source Files section and sidebar
// entry are omitted when the package has no files (the {{if .SourceFiles}}
// guard), rather than rendering an empty section.
func TestPageNoSourceFilesWhenEmpty(t *testing.T) {
	data := renderFixture(t, fixtureSrc, false)
	data.SourceFiles = nil // simulate the no-files case for the shell template
	var sb strings.Builder
	if err := pageTemplate.Execute(&sb, data); err != nil {
		t.Fatal(err)
	}
	page := sb.String()
	if strings.Contains(page, "section-sourcefiles") {
		t.Error("page should omit the Source Files section when there are no files")
	}
}

// TestHandlerServesOutlineScript checks outline.js is served from the embedded
// assets like the stylesheet, so the sidebar is enhanced on a live server.
func TestHandlerServesOutlineScript(t *testing.T) {
	dir := writeDir(t, map[string]string{"a.go": "package foo\n"})
	srv := httptest.NewServer(newHandler(dir, false, nil, nil))
	defer srv.Close()

	body := getBody(t, srv.URL+"/static/outline.js", http.StatusOK)
	if !strings.Contains(body, "js-tree") {
		t.Error("outline.js should be served from embedded assets")
	}
}

// TestOutlineScriptHasCollapseBehavior guards that outline.js still carries the
// collapse/expand wiring ported from pkgsite's tree.ts: items start collapsed
// (aria-expanded="false"), their child list is marked role="group", and the
// aria-level the CSS keys off is assigned. Without these the sidebar would be a
// flat always-open list again, which is exactly the behaviour the issue asked
// to change.
func TestOutlineScriptHasCollapseBehavior(t *testing.T) {
	dir := writeDir(t, map[string]string{"a.go": "package foo\n"})
	srv := httptest.NewServer(newHandler(dir, false, nil, nil))
	defer srv.Close()

	js := getBody(t, srv.URL+"/static/outline.js", http.StatusOK)
	for _, want := range []string{
		`"aria-expanded", "false"`, // collapsed by default
		`"role", "group"`,          // the child list the toggle hides/shows
		`"aria-owns"`,              // links a toggle to its group, like pkgsite
		`"aria-level"`,             // drives the pkgsite tree.css styling
		`"aria-selected", "true"`,  // scroll-spy selection
	} {
		if !strings.Contains(js, want) {
			t.Errorf("outline.js missing collapse/expand marker %q", want)
		}
	}
}

// TestStylesheetHasPkgsiteTreeRules guards that the copied pkgsite tree styling
// stays in page.css: the top-level entries are bold (font-weight 500, not the
// old link style), collapsed groups are hidden, and expanded ones are shown via
// the aria-expanded sibling rule. These are the rules the issue asked to copy
// from pkgsite.
func TestStylesheetHasPkgsiteTreeRules(t *testing.T) {
	dir := writeDir(t, map[string]string{"a.go": "package foo\n"})
	srv := httptest.NewServer(newHandler(dir, false, nil, nil))
	defer srv.Close()

	css := getBody(t, srv.URL+"/static/page.css", http.StatusOK)
	for _, want := range []string{
		".go-Tree a + ul {",                               // collapsed groups hidden
		"a[aria-expanded='true'] + ul[role='group']",      // expanded groups shown
		"font-weight: 500;",                               // bold top-level entries
		"a[aria-expanded='true'][aria-level='2']::before", // the rotating toggle
	} {
		if !strings.Contains(css, want) {
			t.Errorf("page.css missing pkgsite tree rule %q", want)
		}
	}
}
