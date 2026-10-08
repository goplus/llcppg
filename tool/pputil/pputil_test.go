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

package pputil

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/goplus/llcppg/cl/cltest"
	"github.com/qiniu/x/test"
)

func scanAll(src string) []Include {
	var got []Include
	doScanIncludes([]byte(src), func(inc Include) bool {
		got = append(got, inc)
		return true
	})
	return got
}

func TestScanIncludes(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want []Include
	}{
		{
			name: "basic quote and angle",
			src:  "#include \"foo.h\"\n#include <bar.h>\n",
			want: []Include{
				{Filename: "foo.h", Quote: true},
				{Filename: "bar.h", Quote: false},
			},
		},
		{
			name: "spaces and tabs after hash",
			src:  "#  include \"a.h\"\n#\tinclude <b.h>\n#   \t include \"c.h\"\n",
			want: []Include{
				{Filename: "a.h", Quote: true},
				{Filename: "b.h", Quote: false},
				{Filename: "c.h", Quote: true},
			},
		},
		{
			name: "leading whitespace before hash",
			src:  "   #include \"a.h\"\n\t#include <b.h>\n",
			want: []Include{
				{Filename: "a.h", Quote: true},
				{Filename: "b.h", Quote: false},
			},
		},
		{
			name: "ignore line comment",
			src:  "// #include \"skip.h\"\n#include \"keep.h\"\nint x; // #include <also_skip.h>\n",
			want: []Include{
				{Filename: "keep.h", Quote: true},
			},
		},
		{
			name: "ignore block comment single line",
			src:  "/* #include \"skip.h\" */\n#include \"keep.h\"\n",
			want: []Include{
				{Filename: "keep.h", Quote: true},
			},
		},
		{
			name: "ignore block comment multi line",
			src:  "/*\n#include \"skip1.h\"\n#include <skip2.h>\n*/\n#include \"keep.h\"\n",
			want: []Include{
				{Filename: "keep.h", Quote: true},
			},
		},
		{
			name: "same-line block comment before include",
			src:  "/* c */ #include <x.h>\n/* lead */#include \"y.h\"\n",
			want: []Include{
				{Filename: "x.h", Quote: false},
				{Filename: "y.h", Quote: true},
			},
		},
		{
			name: "code before same-line block comment then include is ignored",
			src:  "int x; /* c */ #include <skip.h>\n#include \"keep.h\"\n",
			want: []Include{
				{Filename: "keep.h", Quote: true},
			},
		},
		{
			name: "ignore string literal",
			src:  "const char *s = \"#include <fake.h>\";\n#include \"real.h\"\n",
			want: []Include{
				{Filename: "real.h", Quote: true},
			},
		},
		{
			name: "ignore char literal and escapes",
			src:  "char c = '\"';\nchar *q = \"a\\\"b #include <fake.h>\";\n#include \"real.h\"\n",
			want: []Include{
				{Filename: "real.h", Quote: true},
			},
		},
		{
			name: "hash not at line start is ignored",
			src:  "int x = 1; #include \"skip.h\"\n#include \"keep.h\"\n",
			want: []Include{
				{Filename: "keep.h", Quote: true},
			},
		},
		{
			name: "no separator between include and name",
			src:  "#include\"a.h\"\n#include<b.h>\n",
			want: nil,
		},
		{
			name: "other directives ignored",
			src:  "#ifndef FOO\n#define FOO\n#include \"a.h\"\n#endif\n",
			want: []Include{
				{Filename: "a.h", Quote: true},
			},
		},
		{
			name: "unterminated header name ignored",
			src:  "#include \"a.h\n#include <b.h>\n",
			want: []Include{
				{Filename: "b.h", Quote: false},
			},
		},
		{
			name: "no newline at eof",
			src:  "#include \"a.h\"",
			want: []Include{
				{Filename: "a.h", Quote: true},
			},
		},
		{
			name: "include keyword substring not matched",
			src:  "#includes \"a.h\"\n#include_next <b.h>\n",
			want: nil,
		},
		{
			name: "empty input",
			src:  "",
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := scanAll(tt.src)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("scanIncludes(%q) = %v, want %v", tt.src, got, tt.want)
			}
		})
	}
}

func TestScanIncludesEarlyStop(t *testing.T) {
	src := "#include \"a.h\"\n#include \"b.h\"\n#include \"c.h\"\n"
	var got []Include
	doScanIncludes([]byte(src), func(inc Include) bool {
		got = append(got, inc)
		return len(got) < 2 // stop after two
	})
	want := []Include{
		{Filename: "a.h", Quote: true},
		{Filename: "b.h", Quote: true},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("early stop = %v, want %v", got, want)
	}
}

func TestLoadIncludes(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test.h")
	src := "// #include \"comment.h\"\n" +
		"#  include \"spaced.h\"\n" +
		"const char *s = \"#include <literal.h>\";\n" +
		"#include <sys.h>\n"
	if err := os.WriteFile(file, []byte(src), 0644); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	var got []Include
	for inc := range ScanIncludes(b) {
		got = append(got, inc)
	}
	want := []Include{
		{Filename: "spaced.h", Quote: true},
		{Filename: "sys.h", Quote: false},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("LoadIncludes = %v, want %v", got, want)
	}
}

func TestListIncludesError(t *testing.T) {
	_, err := ListIncludes(filepath.Join(t.TempDir(), "does-not-exist.h"), nil)
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestListFiles(t *testing.T) {
	files := ListFiles("./_testth/basic", false)
	for f, e := range files {
		t.Log(f, e)
		break
	}
}

func testTopHeaders(t *testing.T, dir string) {
	topHeaders, err := TopHeaders(dir, true, true, []string{"./_include"})
	if err != nil {
		t.Fatal(err)
	}
	topHeaderList := strings.Join(topHeaders, "\n")
	exp, _ := os.ReadFile(dir + "/out.expect")
	if test.Diff(t, dir+"/out.expect.txt", []byte(topHeaderList), exp) {
		t.Error(dir, ": unexpect result")
	}
}

func TestTopHeaders(t *testing.T) {
	cltest.TestFromDir(t, "", "./_testth", testTopHeaders)
}

func TestErrTopHeaders(t *testing.T) {
	_, err := TopHeaders("./unknown-dir", true, true, nil)
	if err == nil {
		t.Fatal("TestErrTopHeaders: expected error for unknown directory, got nil")
	}
}

func TestErr_calcHeaderDeps(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for unknown directory, got nil")
		}
	}()
	calcHeaderDeps(map[string]bool{"unkown.h": true}, "./unknown-dir", nil)
}
