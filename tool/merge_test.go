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

package tool

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"sort"
	"strings"
	"testing"
)

// parseFiles parses a set of named sources into a map[name]*ast.File. Comments
// are kept so directives survive, and object resolution is left on so Step 1 has
// imports to resolve.
func parseFiles(t *testing.T, srcs map[string]string) map[string]*ast.File {
	t.Helper()
	files := map[string]*ast.File{}
	fset := token.NewFileSet()
	for name, src := range srcs {
		f, err := parser.ParseFile(fset, name, src, parser.ParseComments)
		if err != nil {
			t.Fatalf("parse %s: %v", name, err)
		}
		files[name] = f
	}
	return files
}

// printFile renders a merged *ast.File to Go source. The merged ASTs carry
// token.NoPos everywhere, so format.Node supplies synthetic positions; this
// mirrors what a real writer would do (see Non-Goals in the proposal).
func printFile(t *testing.T, f *ast.File) string {
	t.Helper()
	var sb strings.Builder
	fset := token.NewFileSet()
	if err := format.Node(&sb, fset, f); err != nil {
		t.Fatalf("format: %v", err)
	}
	return sb.String()
}

// mergedNames returns the sorted output file names.
func mergedNames(m map[string]*ast.File) []string {
	names := make([]string, 0, len(m))
	for n := range m {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

func TestMergeWorkedExample(t *testing.T) {
	darwin := parseFiles(t, map[string]string{
		"types.go": `package libc

import "github.com/goplus/lib/c"

type Mode_t uint16

//go:linkname Getpid C.getpid
func Getpid() c.Int
`,
	})
	linux := parseFiles(t, map[string]string{
		"libc_types.go": `package libc

import libc "github.com/goplus/lib/c"

type Mode_t uint32

//go:linkname Getpid C.getpid
func Getpid() libc.Int
`,
	})

	merged, err := Merge([]*PlatformPackage{
		{Tags: []string{"darwin", "arm64"}, Files: darwin},
		{Tags: []string{"linux", "amd64"}, Files: linux},
	})
	if err != nil {
		t.Fatalf("Merge: %v", err)
	}

	// Getpid is identical (c.Int == libc.Int == github.com/goplus/lib/c.Int), so
	// it is shared and lands in the smallest home file name (libc_types.go).
	// Mode_t differs (uint16 vs uint32), so it is per-platform.
	want := []string{"libc_types.go", "z_darwin_arm64.go", "z_linux_amd64.go"}
	if got := mergedNames(merged); !equalStrs(got, want) {
		t.Fatalf("output files = %v, want %v", got, want)
	}

	shared := printFile(t, merged["libc_types.go"])
	if !strings.Contains(shared, "func Getpid()") {
		t.Errorf("shared file missing Getpid:\n%s", shared)
	}
	if !strings.Contains(shared, "//go:linkname Getpid C.getpid") {
		t.Errorf("shared file missing linkname directive:\n%s", shared)
	}
	// A linkname directive forces import _ "unsafe".
	if !strings.Contains(shared, `_ "unsafe"`) {
		t.Errorf("shared file missing blank unsafe import:\n%s", shared)
	}
	// The import path is preserved; the local name comes from the first platform
	// in sorted order (darwin => alias "c").
	if !strings.Contains(shared, `"github.com/goplus/lib/c"`) {
		t.Errorf("shared file missing lib/c import:\n%s", shared)
	}
	if strings.Contains(shared, "Mode_t") {
		t.Errorf("Mode_t must not be shared:\n%s", shared)
	}

	darwinOut := printFile(t, merged["z_darwin_arm64.go"])
	if !strings.Contains(darwinOut, "//go:build darwin && arm64") {
		t.Errorf("darwin file missing build constraint:\n%s", darwinOut)
	}
	if !strings.Contains(darwinOut, "type Mode_t uint16") {
		t.Errorf("darwin file missing uint16 Mode_t:\n%s", darwinOut)
	}
	linuxOut := printFile(t, merged["z_linux_amd64.go"])
	if !strings.Contains(linuxOut, "//go:build linux && amd64") {
		t.Errorf("linux file missing build constraint:\n%s", linuxOut)
	}
	if !strings.Contains(linuxOut, "type Mode_t uint32") {
		t.Errorf("linux file missing uint32 Mode_t:\n%s", linuxOut)
	}
}

func TestMergeDirectiveDifferenceSplits(t *testing.T) {
	darwin := parseFiles(t, map[string]string{
		"a.go": `package p

//go:linkname Getpid C.getpid
func Getpid() int32
`,
	})
	linux := parseFiles(t, map[string]string{
		"a.go": `package p

//go:linkname Getpid C._getpid
func Getpid() int32
`,
	})
	merged, err := Merge([]*PlatformPackage{
		{Tags: []string{"darwin", "arm64"}, Files: darwin},
		{Tags: []string{"linux", "amd64"}, Files: linux},
	})
	if err != nil {
		t.Fatalf("Merge: %v", err)
	}
	// The directive difference alone makes Getpid platform-specific.
	if _, ok := merged["a.go"]; ok {
		if strings.Contains(printFile(t, merged["a.go"]), "Getpid") {
			t.Errorf("Getpid must not be shared when directives differ")
		}
	}
	for _, n := range []string{"z_darwin_arm64.go", "z_linux_amd64.go"} {
		if _, ok := merged[n]; !ok {
			t.Errorf("expected per-platform file %s", n)
		}
	}
}

func TestMergeSinglePlatformNoBuildConstraint(t *testing.T) {
	only := parseFiles(t, map[string]string{
		"x.go": `package p

type T int

const ( A = iota; B; C )
`,
	})
	merged, err := Merge([]*PlatformPackage{
		{Tags: []string{"linux", "amd64"}, Files: only},
	})
	if err != nil {
		t.Fatalf("Merge: %v", err)
	}
	// With one platform everything is shared; no z_ file is produced.
	for n := range merged {
		if strings.HasPrefix(n, "z_") {
			t.Errorf("single-platform merge produced a z_ file: %s", n)
		}
	}
	out := printFile(t, merged["x.go"])
	if !strings.Contains(out, "type T int") || !strings.Contains(out, "A = iota") {
		t.Errorf("unexpected single-platform output:\n%s", out)
	}
	if strings.Contains(out, "go:build") {
		t.Errorf("single-platform output must have no build constraint:\n%s", out)
	}
}

func TestMergeConstBlockAtomic(t *testing.T) {
	darwin := parseFiles(t, map[string]string{
		"e.go": `package p

const ( EAGAIN = 35; EINVAL = 22 )
`,
	})
	linux := parseFiles(t, map[string]string{
		"e.go": `package p

const ( EAGAIN = 11; EINVAL = 22 )
`,
	})
	merged, err := Merge([]*PlatformPackage{
		{Tags: []string{"darwin", "arm64"}, Files: darwin},
		{Tags: []string{"linux", "amd64"}, Files: linux},
	})
	if err != nil {
		t.Fatalf("Merge: %v", err)
	}
	// One member differs, so the whole block is duplicated per platform.
	d := printFile(t, merged["z_darwin_arm64.go"])
	l := printFile(t, merged["z_linux_amd64.go"])
	if !strings.Contains(d, "EAGAIN = 35") || !strings.Contains(d, "EINVAL = 22") {
		t.Errorf("darwin const block wrong:\n%s", d)
	}
	if !strings.Contains(l, "EAGAIN = 11") || !strings.Contains(l, "EINVAL = 22") {
		t.Errorf("linux const block wrong:\n%s", l)
	}
}

func TestMergeDeterministicOrder(t *testing.T) {
	mk := func() []*PlatformPackage {
		d := parseFiles(t, map[string]string{"t.go": "package p\ntype A int\ntype B int\n"})
		l := parseFiles(t, map[string]string{"t.go": "package p\ntype A int\ntype B int\n"})
		return []*PlatformPackage{
			{Tags: []string{"linux", "amd64"}, Files: l},
			{Tags: []string{"darwin", "arm64"}, Files: d},
		}
	}
	m1, err := Merge(mk())
	if err != nil {
		t.Fatal(err)
	}
	// Reversed input order must give the same result.
	rev := mk()
	rev[0], rev[1] = rev[1], rev[0]
	m2, err := Merge(rev)
	if err != nil {
		t.Fatal(err)
	}
	if printFile(t, m1["t.go"]) != printFile(t, m2["t.go"]) {
		t.Errorf("result depends on input order:\n%s\n---\n%s",
			printFile(t, m1["t.go"]), printFile(t, m2["t.go"]))
	}
}

func TestMergeErrors(t *testing.T) {
	good := parseFiles(t, map[string]string{"a.go": "package p\ntype T int\n"})
	cases := []struct {
		name string
		in   []*PlatformPackage
	}{
		{"bad tags", []*PlatformPackage{{Tags: []string{"linux"}, Files: good}}},
		{"empty tag", []*PlatformPackage{{Tags: []string{"linux", ""}, Files: good}}},
		{"duplicate platform", []*PlatformPackage{
			{Tags: []string{"linux", "amd64"}, Files: parseFiles(t, map[string]string{"a.go": "package p\ntype T int\n"})},
			{Tags: []string{"linux", "amd64"}, Files: parseFiles(t, map[string]string{"a.go": "package p\ntype U int\n"})},
		}},
		{"package mismatch", []*PlatformPackage{
			{Tags: []string{"linux", "amd64"}, Files: parseFiles(t, map[string]string{"a.go": "package p\ntype T int\n"})},
			{Tags: []string{"darwin", "arm64"}, Files: parseFiles(t, map[string]string{"a.go": "package q\ntype T int\n"})},
		}},
		{"reserved z_ name", []*PlatformPackage{
			{Tags: []string{"linux", "amd64"}, Files: parseFiles(t, map[string]string{"z_foo.go": "package p\ntype T int\n"})},
		}},
		{"duplicate symbol", []*PlatformPackage{
			{Tags: []string{"linux", "amd64"}, Files: parseFiles(t, map[string]string{
				"a.go": "package p\ntype T int\n",
				"b.go": "package p\ntype T int\n",
			})},
		}},
		{"init function", []*PlatformPackage{
			{Tags: []string{"linux", "amd64"}, Files: parseFiles(t, map[string]string{"a.go": "package p\nfunc init() {}\n"})},
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := Merge(tc.in); err == nil {
				t.Errorf("expected error for %q, got nil", tc.name)
			}
		})
	}
}

func TestMergeHomeFileSmallestName(t *testing.T) {
	darwin := parseFiles(t, map[string]string{
		"sys_types.go": "package p\ntype Mode_t uint32\n",
	})
	linux := parseFiles(t, map[string]string{
		"types.go": "package p\ntype Mode_t uint32\n",
	})
	merged, err := Merge([]*PlatformPackage{
		{Tags: []string{"darwin", "arm64"}, Files: darwin},
		{Tags: []string{"linux", "amd64"}, Files: linux},
	})
	if err != nil {
		t.Fatalf("Merge: %v", err)
	}
	// Identical type, different file names => home is the smallest name.
	if _, ok := merged["sys_types.go"]; !ok {
		t.Errorf("expected shared symbol in sys_types.go, got files %v", mergedNames(merged))
	}
	if _, ok := merged["types.go"]; ok {
		t.Errorf("types.go must not be a home file (sys_types.go is smaller)")
	}
}

func TestMergeMissingOnOnePlatformIsPerPlatform(t *testing.T) {
	darwin := parseFiles(t, map[string]string{
		"a.go": "package p\ntype Only int\ntype Shared int\n",
	})
	linux := parseFiles(t, map[string]string{
		"a.go": "package p\ntype Shared int\n",
	})
	merged, err := Merge([]*PlatformPackage{
		{Tags: []string{"darwin", "arm64"}, Files: darwin},
		{Tags: []string{"linux", "amd64"}, Files: linux},
	})
	if err != nil {
		t.Fatalf("Merge: %v", err)
	}
	shared := printFile(t, merged["a.go"])
	if !strings.Contains(shared, "type Shared int") {
		t.Errorf("Shared should be in a.go:\n%s", shared)
	}
	if strings.Contains(shared, "Only") {
		t.Errorf("Only exists on one platform, must not be shared:\n%s", shared)
	}
	d := printFile(t, merged["z_darwin_arm64.go"])
	if !strings.Contains(d, "type Only int") {
		t.Errorf("Only should be in z_darwin_arm64.go:\n%s", d)
	}
}

// TestMergeShadowedPackageName checks that a parameter named like an imported
// package is not mistaken for a package ref (Step 1 shadowing rule). The two
// platforms spell the parameter differently but with the same declared type, so
// the function is still shared.
func TestMergeShadowedPackageName(t *testing.T) {
	darwin := parseFiles(t, map[string]string{
		"a.go": `package p

import "github.com/goplus/lib/c"

func F(c int) c.Int
`,
	})
	linux := parseFiles(t, map[string]string{
		"a.go": `package p

import "github.com/goplus/lib/c"

func F(x int) c.Int
`,
	})
	// Note: "c int" shadows package c only within F's signature/body; the result
	// type c.Int still refers to the package. Parameter names are ignored, so F
	// is identical and shared.
	merged, err := Merge([]*PlatformPackage{
		{Tags: []string{"darwin", "arm64"}, Files: darwin},
		{Tags: []string{"linux", "amd64"}, Files: linux},
	})
	if err != nil {
		t.Fatalf("Merge: %v", err)
	}
	for n := range merged {
		if strings.HasPrefix(n, "z_") {
			t.Errorf("F should be shared, got per-platform file %s", n)
		}
	}
}

// TestMergeStructFieldsCompared checks that struct fields are compared by name
// and tag (field order matters), so a field rename makes the type per-platform.
func TestMergeStructFieldsCompared(t *testing.T) {
	darwin := parseFiles(t, map[string]string{
		"a.go": "package p\ntype S struct{ A int }\n",
	})
	linux := parseFiles(t, map[string]string{
		"a.go": "package p\ntype S struct{ B int }\n",
	})
	merged, err := Merge([]*PlatformPackage{
		{Tags: []string{"darwin", "arm64"}, Files: darwin},
		{Tags: []string{"linux", "amd64"}, Files: linux},
	})
	if err != nil {
		t.Fatalf("Merge: %v", err)
	}
	if _, ok := merged["z_darwin_arm64.go"]; !ok {
		t.Errorf("struct with different field name must be per-platform, got %v", mergedNames(merged))
	}
}

// TestMergeDifferentAliasesStillShared mirrors the worked example: the same
// import path under different local aliases is still equal (identity is the
// import path, not the local name).
func TestMergeDifferentAliasesStillShared(t *testing.T) {
	darwin := parseFiles(t, map[string]string{
		"a.go": `package p

import "github.com/goplus/lib/c"

var X c.Int
`,
	})
	linux := parseFiles(t, map[string]string{
		"a.go": `package p

import libc "github.com/goplus/lib/c"

var X libc.Int
`,
	})
	merged, err := Merge([]*PlatformPackage{
		{Tags: []string{"darwin", "arm64"}, Files: darwin},
		{Tags: []string{"linux", "amd64"}, Files: linux},
	})
	if err != nil {
		t.Fatalf("Merge: %v", err)
	}
	for n := range merged {
		if strings.HasPrefix(n, "z_") {
			t.Errorf("X should be shared despite alias difference, got %s", n)
		}
	}
	out := printFile(t, merged["a.go"])
	if !strings.Contains(out, `"github.com/goplus/lib/c"`) {
		t.Errorf("merged import path wrong:\n%s", out)
	}
}

func equalStrs(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
