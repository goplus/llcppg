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

package cl_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/goplus/llcppg/cl"
	"github.com/goplus/llcppg/cl/cltest"
	"github.com/goplus/llcppg/clang"
)

// TestDefineEndToEnd drives the full -d pipeline: a header with real macros is
// parsed by libclang, converted by cl.NewPackage, and the Define constants are
// injected after the load phase. It verifies value typing and that $(name)
// resolves against the macro-derived Go names in the generated package.
func TestDefineEndToEnd(t *testing.T) {
	const header = `
#define VERSION "3.14.8"
#define PY_MAJOR_VERSION 3
`
	dir := t.TempDir()
	hdr := filepath.Join(dir, "in.h")
	if err := os.WriteFile(hdr, []byte(header), 0644); err != nil {
		t.Fatal(err)
	}

	idx := clang.CreateIndex(1, 1)
	defer idx.Dispose()
	u, err := idx.ParseTranslationUnit(
		clang.DetailedPreprocessingRecord, hdr, "-x", cltest.LanguageOf(cl.LanguageC), "-fparse-all-comments")
	if err != nil {
		t.Fatal("ParseTranslationUnit failed:", err)
	}
	defer u.Dispose()

	pkg, err := cl.NewPackage("testcl/foo", "foo", []cl.Source{u}, &cl.Config{
		Language:      cl.LanguageC,
		DefaultGoFile: "llcppg.i.go",
		Define: []cl.Define{
			{Name: "LLGoLibVersion", Value: "v$(VERSION)", Source: "-d"},
			{Name: "ABIMajor", Value: "$(PY_MAJOR_VERSION)", Source: "-d"},
			{Name: "ABIMajorString", Value: "`$(PY_MAJOR_VERSION)`", Source: "-d"},
			{Name: "Channel", Value: "stable", Source: "-d"},
		},
		PubFileLookup: func(pkgPath string) (string, bool) { return "", false },
		PackageOf:     func(headerFile string) (string, bool) { return "testcl/foo", true },
		GoFileOf:      func(headerFile string) (string, bool) { return "llcppg.i.go", true },
	})
	if err != nil {
		t.Fatal("NewPackage failed:", err)
	}

	var b bytes.Buffer
	if err := pkg.WriteTo(&b); err != nil {
		t.Fatal("WriteTo failed:", err)
	}
	src := b.String()

	for _, want := range []string{
		`LLGoLibVersion = "v3.14.8"`,
		`ABIMajor = 3`,
		`ABIMajorString = "3"`,
		`Channel = "stable"`,
	} {
		if !containsNormalized(src, want) {
			t.Errorf("generated source missing %q\n---\n%s", want, src)
		}
	}
}

// containsNormalized reports whether src contains want after collapsing runs of
// spaces/tabs on each line (gogen aligns const blocks with padding).
func containsNormalized(src, want string) bool {
	norm := func(s string) string {
		lines := strings.Split(s, "\n")
		for i, ln := range lines {
			lines[i] = strings.Join(strings.Fields(ln), " ")
		}
		return strings.Join(lines, "\n")
	}
	return strings.Contains(norm(src), want)
}
