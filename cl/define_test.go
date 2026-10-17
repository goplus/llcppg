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

package cl

import (
	"bytes"
	"go/constant"
	"go/token"
	"go/types"
	"strings"
	"testing"

	"github.com/goplus/gogen"
)

// -----------------------------------------------------------------------------

func TestParseDefineArgs(t *testing.T) {
	t.Run("single entry", func(t *testing.T) {
		defs, err := ParseDefineArgs([]string{"NAME=VALUE"})
		if err != nil {
			t.Fatal(err)
		}
		if len(defs) != 1 || defs[0].Name != "NAME" || defs[0].Value != "VALUE" {
			t.Fatalf("got %+v", defs)
		}
		if defs[0].Source != "-d" || defs[0].Group != 0 {
			t.Fatalf("unexpected source/group: %+v", defs[0])
		}
	})

	t.Run("multiple entries, trimming", func(t *testing.T) {
		defs, err := ParseDefineArgs([]string{" A = 1 ;  B = hello world "})
		if err != nil {
			t.Fatal(err)
		}
		want := []Define{{Name: "A", Value: "1"}, {Name: "B", Value: "hello world"}}
		if len(defs) != len(want) {
			t.Fatalf("len = %d, want %d (%+v)", len(defs), len(want), defs)
		}
		for i, w := range want {
			if defs[i].Name != w.Name || defs[i].Value != w.Value {
				t.Errorf("defs[%d] = %+v, want name=%q value=%q", i, defs[i], w.Name, w.Value)
			}
		}
	})

	t.Run("value may contain =", func(t *testing.T) {
		defs, err := ParseDefineArgs([]string{"K=a=b=c"})
		if err != nil {
			t.Fatal(err)
		}
		if defs[0].Value != "a=b=c" {
			t.Fatalf("value = %q, want %q", defs[0].Value, "a=b=c")
		}
	})

	t.Run("multiple -d flags get distinct groups", func(t *testing.T) {
		defs, err := ParseDefineArgs([]string{"A=1", "B=2"})
		if err != nil {
			t.Fatal(err)
		}
		if defs[0].Group != 0 || defs[1].Group != 1 {
			t.Fatalf("groups = %d,%d; want 0,1", defs[0].Group, defs[1].Group)
		}
	})

	t.Run("empty entries tolerated", func(t *testing.T) {
		defs, err := ParseDefineArgs([]string{"A=1;;"})
		if err != nil {
			t.Fatal(err)
		}
		if len(defs) != 1 {
			t.Fatalf("len = %d, want 1", len(defs))
		}
	})

	t.Run("missing =", func(t *testing.T) {
		if _, err := ParseDefineArgs([]string{"NOEQUALS"}); err == nil {
			t.Fatal("expected error for missing '='")
		}
	})

	t.Run("invalid identifier", func(t *testing.T) {
		for _, bad := range []string{"1abc=1", "a-b=1", "=1", "func=1"} {
			if _, err := ParseDefineArgs([]string{bad}); err == nil {
				t.Errorf("expected error for %q", bad)
			}
		}
	})
}

// -----------------------------------------------------------------------------

func TestSoleRef(t *testing.T) {
	tests := []struct {
		raw    string
		name   string
		wantOk bool
	}{
		{"$(VERSION)", "VERSION", true},
		{"v$(VERSION)", "", false},
		{"$(A)$(B)", "", false},
		{"$(A).$(B)", "", false},
		{"100", "", false},
		{"$(A)x", "", false},
		{"x$(A)", "", false},
	}
	for _, tt := range tests {
		name, ok := soleRef(tt.raw)
		if ok != tt.wantOk || name != tt.name {
			t.Errorf("soleRef(%q) = (%q,%v), want (%q,%v)", tt.raw, name, ok, tt.name, tt.wantOk)
		}
	}
}

func TestLiteralValue(t *testing.T) {
	tests := []struct {
		raw  string
		kind constant.Kind
		str  string // ExactString form
	}{
		{"100", constant.Int, "100"},
		{"0x10", constant.Int, "16"},
		{"0b101", constant.Int, "5"},
		{"-1", constant.Int, "-1"},
		{"3.14", constant.Float, "157/50"},
		{"1e9", constant.Float, "1000000000"},
		{"v3.14.8", constant.String, `"v3.14.8"`},
		{"stable", constant.String, `"stable"`},
		{"", constant.String, `""`},
	}
	for _, tt := range tests {
		v := literalValue(tt.raw)
		if v.Kind() != tt.kind {
			t.Errorf("literalValue(%q).Kind() = %v, want %v", tt.raw, v.Kind(), tt.kind)
			continue
		}
		if got := v.ExactString(); got != tt.str {
			t.Errorf("literalValue(%q) = %s, want %s", tt.raw, got, tt.str)
		}
	}
}

// -----------------------------------------------------------------------------

// newTestPkgCtx builds a minimal pkgCtx backed by a real gogen package, with
// the given seed constants inserted into the package scope so that $(name)
// references and name-conflict checks can be exercised without clang.
func newTestPkgCtx(seed map[string]any) *pkgCtx {
	pkg := gogen.NewPackage("", "foo", &gogen.Config{
		Fset:          token.NewFileSet(),
		DefaultGoFile: "llcppg.i.go",
	})
	pkg.SetRedeclarable(true)
	if len(seed) > 0 {
		scope := pkg.Types.Scope()
		for name, val := range seed {
			scope.Insert(types.NewConst(token.NoPos, pkg.Types, name, types.Typ[types.UntypedInt], toConst(val)))
		}
	}
	return &pkgCtx{pkg: pkg, cb: pkg.CB(), fset: pkg.Fset, defaultGoFile: "llcppg.i.go"}
}

func toConst(v any) constant.Value {
	switch v := v.(type) {
	case int:
		return constant.MakeInt64(int64(v))
	case string:
		return constant.MakeString(v)
	case float64:
		return constant.MakeFloat64(v)
	default:
		panic("unsupported seed value")
	}
}

func genSource(t *testing.T, ctx *pkgCtx) string {
	t.Helper()
	var b bytes.Buffer
	if err := ctx.pkg.WriteTo(&b); err != nil {
		t.Fatalf("WriteTo: %v", err)
	}
	return b.String()
}

func lookupConst(ctx *pkgCtx, name string) (*types.Const, bool) {
	o := ctx.pkg.Types.Scope().Lookup(name)
	c, ok := o.(*types.Const)
	return c, ok
}

// readableConst renders a constant value in a human-readable form: numbers in
// decimal (0.5, 3) and strings as their raw content (no quotes).
func readableConst(v constant.Value) string {
	switch v.Kind() {
	case constant.String:
		return constant.StringVal(v)
	case constant.Float:
		return v.String()
	default:
		return v.ExactString()
	}
}

func TestApplyDefines_Typing(t *testing.T) {
	ctx := newTestPkgCtx(nil)
	ctx.applyDefines([]Define{
		{Name: "MaxRetries", Value: "3", Source: "-d"},
		{Name: "Scale", Value: "0.5", Source: "-d"},
		{Name: "Channel", Value: "stable", Source: "-d"},
		{Name: "BuildID", Value: "`20261017`", Source: "-d"},
	})
	if ctx.errCnt != 0 {
		t.Fatalf("errCnt = %d, want 0", ctx.errCnt)
	}
	// Readable value per kind: int "3", float "0.5", string "stable".
	checks := []struct {
		name string
		kind constant.Kind
		str  string
	}{
		{"MaxRetries", constant.Int, "3"},
		{"Scale", constant.Float, "0.5"},
		{"Channel", constant.String, "stable"},
		{"BuildID", constant.String, "20261017"},
	}
	for _, c := range checks {
		cst, ok := lookupConst(ctx, c.name)
		if !ok {
			t.Errorf("const %s not found", c.name)
			continue
		}
		if cst.Val().Kind() != c.kind {
			t.Errorf("%s kind = %v, want %v", c.name, cst.Val().Kind(), c.kind)
		}
		if got := readableConst(cst.Val()); got != c.str {
			t.Errorf("%s = %s, want %s", c.name, got, c.str)
		}
	}
	// The constants are emitted together in a single const block in declaration
	// order in the default file.
	src := normalizeWS(genSource(t, ctx))
	for _, want := range []string{"MaxRetries = 3", "Scale = 0.5", `Channel = "stable"`, `BuildID = "20261017"`} {
		if !strings.Contains(src, want) {
			t.Errorf("generated source missing %q\n%s", want, src)
		}
	}
	if n := strings.Count(src, "const ("); n != 1 {
		t.Errorf("expected a single const block, got %d\n%s", n, src)
	}
}

// normalizeWS collapses runs of spaces/tabs so gogen's column alignment does
// not interfere with substring assertions.
func normalizeWS(s string) string {
	lines := strings.Split(s, "\n")
	for i, ln := range lines {
		lines[i] = strings.Join(strings.Fields(ln), " ")
	}
	return strings.Join(lines, "\n")
}

func TestApplyDefines_References(t *testing.T) {
	ctx := newTestPkgCtx(map[string]any{
		"VERSION":          "3.14.8",
		"PY_MAJOR_VERSION": 3,
	})
	ctx.applyDefines([]Define{
		{Name: "LLGoLibVersion", Value: "v$(VERSION)", Source: "-d"},
		{Name: "ABIMajor", Value: "$(PY_MAJOR_VERSION)", Source: "-d"},
		{Name: "ABIMajorString", Value: "`$(PY_MAJOR_VERSION)`", Source: "-d"},
		{Name: "VerCopy", Value: "$(VERSION)", Source: "-d"},
		{Name: "Combo", Value: "$(PY_MAJOR_VERSION).$(VERSION)", Source: "-d"},
	})
	if ctx.errCnt != 0 {
		t.Fatalf("errCnt = %d, want 0", ctx.errCnt)
	}
	checks := []struct {
		name string
		kind constant.Kind
		str  string
	}{
		{"LLGoLibVersion", constant.String, `"v3.14.8"`}, // interpolation -> string
		{"ABIMajor", constant.Int, "3"},                  // pure reference -> inherits int
		{"ABIMajorString", constant.String, `"3"`},       // backtick -> string of int
		{"VerCopy", constant.String, `"3.14.8"`},         // pure reference -> inherits string
		{"Combo", constant.String, `"3.3.14.8"`},         // multi-ref interpolation
	}
	for _, c := range checks {
		cst, ok := lookupConst(ctx, c.name)
		if !ok {
			t.Errorf("const %s not found", c.name)
			continue
		}
		if cst.Val().Kind() != c.kind {
			t.Errorf("%s kind = %v, want %v", c.name, cst.Val().Kind(), c.kind)
		}
		if cst.Val().String() != c.str {
			t.Errorf("%s = %s, want %s", c.name, cst.Val().String(), c.str)
		}
	}
}

func TestApplyDefines_RefBetweenDefines(t *testing.T) {
	ctx := newTestPkgCtx(nil)
	// Forward reference: A refers to B which is defined later.
	ctx.applyDefines([]Define{
		{Name: "A", Value: "x$(B)", Source: "-d"},
		{Name: "B", Value: "42", Source: "-d"},
	})
	if ctx.errCnt != 0 {
		t.Fatalf("errCnt = %d, want 0", ctx.errCnt)
	}
	if cst, ok := lookupConst(ctx, "A"); !ok || cst.Val().String() != `"x42"` {
		t.Errorf("A = %v (ok=%v), want \"x42\"", cst, ok)
	}
}

func TestApplyDefines_DollarEscapes(t *testing.T) {
	ctx := newTestPkgCtx(nil)
	ctx.applyDefines([]Define{
		{Name: "Price", Value: "$$5", Source: "-d"}, // $$ -> $, string
		{Name: "Raw", Value: "a$b", Source: "-d"},   // lone $ kept, string
	})
	if ctx.errCnt != 0 {
		t.Fatalf("errCnt = %d, want 0", ctx.errCnt)
	}
	if cst, _ := lookupConst(ctx, "Price"); cst.Val().String() != `"$5"` {
		t.Errorf("Price = %s, want \"$5\"", cst.Val().String())
	}
	if cst, _ := lookupConst(ctx, "Raw"); cst.Val().String() != `"a$b"` {
		t.Errorf("Raw = %s, want \"a$b\"", cst.Val().String())
	}
}

func TestApplyDefines_Errors(t *testing.T) {
	tests := []struct {
		name string
		seed map[string]any
		defs []Define
	}{
		{
			name: "duplicate in same group",
			defs: []Define{
				{Name: "A", Value: "1", Source: "-d", Group: 0},
				{Name: "A", Value: "2", Source: "-d", Group: 0},
			},
		},
		{
			name: "conflict with existing symbol",
			seed: map[string]any{"Existing": 1},
			defs: []Define{{Name: "Existing", Value: "2", Source: "-d"}},
		},
		{
			name: "unresolved reference",
			defs: []Define{{Name: "A", Value: "$(MISSING)", Source: "-d"}},
		},
		{
			name: "reference cycle",
			defs: []Define{
				{Name: "A", Value: "$(B)", Source: "-d"},
				{Name: "B", Value: "$(A)", Source: "-d"},
			},
		},
		{
			name: "unterminated backtick",
			defs: []Define{{Name: "A", Value: "`abc", Source: "-d"}},
		},
		{
			name: "unterminated reference",
			defs: []Define{{Name: "A", Value: "x$(B", Source: "-d"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := newTestPkgCtx(tt.seed)
			ctx.applyDefines(tt.defs)
			if ctx.errCnt == 0 {
				t.Errorf("expected an error, got none")
			}
		})
	}
}

func TestApplyDefines_RefToNonConst(t *testing.T) {
	pkg := gogen.NewPackage("", "foo", &gogen.Config{
		Fset:          token.NewFileSet(),
		DefaultGoFile: "llcppg.i.go",
	})
	pkg.SetRedeclarable(true)
	// Insert a type name (not a constant) into scope.
	scope := pkg.Types.Scope()
	tn := types.NewTypeName(token.NoPos, pkg.Types, "SomeType", nil)
	types.NewNamed(tn, types.Typ[types.Int], nil)
	scope.Insert(tn)

	ctx := &pkgCtx{pkg: pkg, cb: pkg.CB(), fset: pkg.Fset, defaultGoFile: "llcppg.i.go"}
	ctx.applyDefines([]Define{{Name: "A", Value: "$(SomeType)", Source: "-d"}})
	if ctx.errCnt == 0 {
		t.Error("expected error referencing a non-constant")
	}
}

func TestApplyDefines_LaterGroupOverrides(t *testing.T) {
	ctx := newTestPkgCtx(nil)
	// Two -d groups: later group wins, keeping first slot's position.
	ctx.applyDefines([]Define{
		{Name: "Channel", Value: "stable", Source: "-d", Group: 0},
		{Name: "Channel", Value: "beta", Source: "-d", Group: 1},
	})
	if ctx.errCnt != 0 {
		t.Fatalf("errCnt = %d, want 0", ctx.errCnt)
	}
	if cst, _ := lookupConst(ctx, "Channel"); cst.Val().String() != `"beta"` {
		t.Errorf("Channel = %s, want \"beta\"", cst.Val().String())
	}
}

// -----------------------------------------------------------------------------
