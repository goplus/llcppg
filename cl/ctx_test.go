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
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"strings"
	"testing"

	"github.com/goplus/gogen"
	"github.com/goplus/lib/c"
	"github.com/goplus/llcppg/clang"
)

func newPkgCtx() *pkgCtx {
	return &pkgCtx{
		pkg:   gogen.NewPackage("", "foo", nil),
		fset:  token.NewFileSet(),
		files: make(map[clang.File]file),
	}
}

func Test_loadGoImportSpec(t *testing.T) {
	ctx := newPkgCtx()
	pkgTypes := ctx.pkg.Types
	loadGoImportSpec(ctx, pkgTypes, pkgTypes.Scope(), &ast.ImportSpec{
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: `"fmt`,
		},
	})
	if ctx.errCnt != 1 {
		t.Fatalf("expected 1 error, got %d", ctx.errCnt)
	}
	loadGoImportSpec(ctx, pkgTypes, pkgTypes.Scope(), &ast.ImportSpec{
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: `"unknown"`,
		},
	})
	if ctx.errCnt != 2 {
		t.Fatalf("expected 2 errors, got %d", ctx.errCnt)
	}
}

func Test_insertObj(t *testing.T) {
	ctx := newPkgCtx()
	pkgTypes := ctx.pkg.Types
	scope := pkgTypes.Scope()
	insertObj(ctx, scope, types.NewConst(0, pkgTypes, "bar", nil, constant.MakeInt64(0)))
	insertObj(ctx, scope, types.NewConst(0, pkgTypes, "bar", nil, constant.MakeInt64(1)))
	if ctx.errCnt != 1 {
		t.Fatalf("expected 1 error, got %d", ctx.errCnt)
	}
}

func Test_npErrorf(t *testing.T) {
	defer func() {
		if e := recover(); e == nil {
			t.Fatal("Test_npErrorf: no error?")
		}
	}()
	exit = func(code int) {
		panic(code)
	}
	ctx := &pkgCtx{failFast: 1}
	ctx.npErrorf("Test_npErrorf")
}

func lineCommentsText(raw string) string {
	comments := toLineComments(raw)
	lines := make([]string, len(comments))
	for i, c := range comments {
		lines[i] = c.Text
	}
	return strings.Join(lines, "\n")
}

func TestToLineComments(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want string
	}{
		{
			name: "empty",
			raw:  "",
			want: "",
		},
		{
			name: "line_comment",
			raw:  "/// A documented enum type.",
			want: "// A documented enum type.",
		},
		{
			name: "bang_line_comment",
			raw:  "//! A bang documented item.",
			want: "// A bang documented item.",
		},
		{
			name: "plain_line_comment",
			raw:  "// plain comment",
			want: "// plain comment",
		},
		{
			name: "single_line_block",
			raw:  "/* simple */",
			want: "// simple",
		},
		{
			name: "doxygen_block",
			raw:  "/**\n * A documented function.\n *\n * It adds two integers.\n */",
			want: "// A documented function.\n//\n// It adds two integers.",
		},
		{
			// A "/***" opener must not leak a stray "*" into a "// *" line.
			name: "triple_star_opener",
			raw:  "/***\n * Determine whether the given cursor represents a preprocessing\n * element.\n */",
			want: "// Determine whether the given cursor represents a preprocessing\n// element.",
		},
		{
			// Interior indentation after the "* " decoration is preserved.
			name: "triple_star_preserves_indent",
			raw:  "/***\n * Determine whether the given cursor represents a currently\n *  unexposed piece of the AST.\n */",
			want: "// Determine whether the given cursor represents a currently\n//  unexposed piece of the AST.",
		},
		{
			// Banner-style separators must collapse to blank lines, never
			// emit rows of asterisks or a trailing "**/".
			name: "banner_block",
			raw:  "/**************\nSymbols and macros to supply platform-independent interfaces.\n**************/",
			want: "// Symbols and macros to supply platform-independent interfaces.",
		},
		{
			name: "block_with_star_prefix",
			raw:  "/* uintptr_t is the C9X name for a type such that a\n * void* can be cast to uintptr_t and back.\n */",
			want: "// uintptr_t is the C9X name for a type such that a\n// void* can be cast to uintptr_t and back.",
		},
		{
			name: "block_no_star_prefix",
			raw:  "/* line one\nline two */",
			want: "// line one\n// line two",
		},
		{
			name: "whitespace_only_block",
			raw:  "/**\n *\n */",
			want: "",
		},
		{
			// libclang concatenates a plain banner and the following doc
			// block into one raw comment; neither the interior "*/" nor the
			// interior "/**" may leak into the output.
			name: "multi_block_plain_then_doc",
			raw:  "/* Declarations */\n  /**\n   * A declaration whose specific kind is not exposed.\n   */",
			want: "// Declarations\n//\n// A declaration whose specific kind is not exposed.",
		},
		{
			name: "inline_block_comment",
			raw:  "/* Decl references */",
			want: "// Decl references",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lineCommentsText(tt.raw)
			if got != tt.want {
				t.Errorf("toLineComments(%q):\n got:\n%s\nwant:\n%s", tt.raw, got, tt.want)
			}
		})
	}
}

func Test_typeNameOf(t *testing.T) {
	if name := typeNameOf(types.Typ[types.Int]); name != "" {
		t.Fatal("Test_typeNameOf: expected empty string for types.Int")
	}
}

func Test_isMethodCall(t *testing.T) {
	if new(funcCallExpr).isMethodCall("foo") {
		t.Fatal("Test_isMethodCall: expected false")
	}
}

func Test_nodeInterp(t *testing.T) {
	func() {
		defer func() {
			if e := recover(); e != nil {
				t.Logf("Test_nodeInterp: panic?")
			}
		}()
		var ni nodeInterp
		if v := ni.LoadExpr(nil); v != "" {
			t.Logf("Test_nodeInterp: expected empty string for nil node, got %q", v)
		}
	}()
	const content = "hello"
	var ni nodeInterp
	var fp = clang.File(1)
	ctx := newPkgCtx()
	node := newNode(ctx, content, fp, 1, 2)
	if v := ni.LoadExpr(node); v != "e" {
		t.Fatal("Test_nodeInterp: expected 'e', got ", v)
	}
	node = newNode(ctx, content, fp, 4, 6)
	if v := ni.LoadExpr(node); v != "o" {
		t.Fatal("Test_nodeInterp: expected 'o', got ", v)
	}
	if _, _, ok := ctx.fileOffset(0); ok {
		t.Fatal("Test_nodeInterp: expected fileOffset(0) to be not ok")
	}
}

func newNode(ctx *pkgCtx, content string, fp clang.File, pos, end c.Uint) ast.Node {
	base := getFileBase(ctx, content, fp)
	return &node{pos: token.Pos(int(pos) + base), end: token.Pos(int(end) + base), ctx: ctx}
}

func getFileBase(p *pkgCtx, content string, fp clang.File) int {
	f, ok := p.files[fp]
	if !ok {
		src := []byte(content)
		tf := p.fset.AddFile("foo.h", -1, len(src))
		tf.SetLinesForContent(src)
		f = file{tf.Base(), src}
		p.files[fp] = f
	}
	return f.base
}
