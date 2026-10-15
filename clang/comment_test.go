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

package clang

import (
	"strings"
	"testing"
)

// TestPrecedingComment exercises the source-scanning that recovers a macro's
// doc comment, which libclang does not associate with macro cursors. The marker
// "@" in each source marks the byte MacroRawComment would scan back from (the
// start of the "#define" line); it is stripped before the call.
func TestPrecedingComment(t *testing.T) {
	tests := []struct {
		name string
		src  string // "@" marks the scan start (start of the macro line)
		want string
	}{
		{
			name: "line_comment_above",
			src:  "/// Maximum size limit.\n@#define MAX_SIZE 100\n",
			want: "/// Maximum size limit.",
		},
		{
			name: "block_comment_above",
			src:  "/** Initialize for deflation. */\n@#define z_deflateInit(a, b) x\n",
			want: "/** Initialize for deflation. */",
		},
		{
			name: "multiline_block_above",
			src:  "/*\n * Init for deflation.\n */\n@#define Z 1\n",
			want: "/*\n * Init for deflation.\n */",
		},
		{
			name: "consecutive_line_comments",
			src:  "// line one\n// line two\n@#define X 1\n",
			want: "// line one\n// line two",
		},
		{
			name: "blank_line_detaches",
			src:  "/// detached doc\n\n@#define X 1\n",
			want: "",
		},
		{
			name: "blank_line_detaches_only_first",
			src:  "// a\n\n// b\n@#define X 1\n",
			want: "// b",
		},
		{
			name: "no_comment",
			src:  "int x;\n@#define X 1\n",
			want: "",
		},
		{
			name: "code_above_comment_adjacent",
			src:  "int y;\n/// doc\n@#define X 1\n",
			want: "/// doc",
		},
		{
			name: "trailing_comment_on_code_line",
			src:  "int z; // trailing\n@#define X 1\n",
			want: "",
		},
		{
			name: "comment_at_file_start",
			src:  "/// doc\n@#define X 1",
			want: "/// doc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			at := strings.IndexByte(tt.src, '@')
			if at < 0 {
				t.Fatalf("test source missing '@' marker")
			}
			src := tt.src[:at] + tt.src[at+1:]
			got := precedingComment([]byte(src), at)
			if got != tt.want {
				t.Errorf("precedingComment:\n got: %q\nwant: %q", got, tt.want)
			}
		})
	}
}

// TestFollowingComment exercises the source-scanning that recovers a macro's
// trailing comment (issue #1056: for macros, a trailing comment on the line
// directly below the #define is used when there is no leading comment). The
// marker "@" in each source marks the byte MacroRawComment scans forward from
// (the end of the macro's extent); it is stripped before the call.
func TestFollowingComment(t *testing.T) {
	tests := []struct {
		name string
		src  string // "@" marks the scan start (end of the macro extent)
		want string
	}{
		{
			name: "line_comment_below",
			src:  "#define X 1@\n/// Maximum size limit.\n",
			want: "/// Maximum size limit.",
		},
		{
			name: "block_comment_below",
			src:  "#define Z 1@\n/** Init for deflation. */\n",
			want: "/** Init for deflation. */",
		},
		{
			name: "multiline_block_below",
			src:  "#define Z 1@\n/*\n * Init for deflation.\n */\n",
			want: "/*\n * Init for deflation.\n */",
		},
		{
			name: "consecutive_line_comments",
			src:  "#define X 1@\n// line one\n// line two\nint y;\n",
			want: "// line one\n// line two",
		},
		{
			name: "blank_line_detaches",
			src:  "#define X 1@\n\n/// detached doc\n",
			want: "",
		},
		{
			name: "inline_comment_not_trailing",
			src:  "#define X 1@ // inline\n",
			want: "",
		},
		{
			name: "code_below_no_comment",
			src:  "#define X 1@\nint y;\n",
			want: "",
		},
		{
			name: "next_decl_comment_detached_by_blank",
			src:  "#define X 1@\n\n/// doc for next\n#define Y 2\n",
			want: "",
		},
		{
			name: "comment_at_file_end_no_newline",
			src:  "#define X 1@\n/// doc",
			want: "/// doc",
		},
		{
			name: "no_comment_at_eof",
			src:  "#define X 1@",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			at := strings.IndexByte(tt.src, '@')
			if at < 0 {
				t.Fatalf("test source missing '@' marker")
			}
			src := tt.src[:at] + tt.src[at+1:]
			got := followingComment([]byte(src), at)
			if got != tt.want {
				t.Errorf("followingComment:\n got: %q\nwant: %q", got, tt.want)
			}
		})
	}
}
