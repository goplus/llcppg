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
	"path/filepath"
	"testing"
)

func TestRelSourceFile(t *testing.T) {
	dir := filepath.FromSlash("/pkg/clang")
	p := &pkg{Dir: dir}

	ok := []struct{ in, want string }{
		{filepath.Join(dir, "clang.go"), "clang.go"},
		{filepath.Join(dir, "a_b.go"), "a_b.go"},
	}
	for _, c := range ok {
		got, valid := p.relSourceFile(c.in)
		if !valid || got != c.want {
			t.Errorf("relSourceFile(%q) = (%q, %v), want (%q, true)", c.in, got, valid, c.want)
		}
	}

	bad := []string{
		filepath.Join(dir, "notes.txt"),           // not a .go file
		filepath.Join(dir, "sub", "x.go"),         // nested, not a direct child
		filepath.Join(dir, "..", "other", "x.go"), // escapes the package dir
		dir,               // the directory itself
		filepath.Dir(dir), // a parent directory
	}
	for _, in := range bad {
		if got, valid := p.relSourceFile(in); valid {
			t.Errorf("relSourceFile(%q) = (%q, true), want rejected", in, got)
		}
	}
}

func TestSourceLinkFuncNilNode(t *testing.T) {
	p := &pkg{Dir: filepath.FromSlash("/pkg/clang")}
	if url := p.sourceLinkFunc()(nil); url != "" {
		t.Errorf("sourceLinkFunc()(nil) = %q, want empty", url)
	}
}
