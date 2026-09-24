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
	"testing"
)

func TestToLineComments(t *testing.T) {
	got := func(raw string) []string {
		cs := toLineComments(raw)
		out := make([]string, len(cs))
		for i, c := range cs {
			out[i] = c.Text
		}
		return out
	}
	eq := func(name string, a, b []string) {
		if len(a) != len(b) {
			t.Fatalf("%s: len mismatch: got %d %q, want %d %q", name, len(a), a, len(b), b)
		}
		for i := range a {
			if a[i] != b[i] {
				t.Fatalf("%s: line %d: got %q, want %q", name, i, a[i], b[i])
			}
		}
	}

	eq("empty", got(""), nil)
	eq("blank", got("   \n  \n"), nil)

	eq("single line //", got("// hello world"), []string{"// hello world"})
	eq("doxygen ///", got("/// a typedef alias"), []string{"// a typedef alias"})
	eq("doxygen //!", got("//! bang doc"), []string{"// bang doc"})

	eq("multi line //",
		got("// line one\n// line two"),
		[]string{"// line one", "// line two"})

	eq("block /* */",
		got("/* a brief doc */"),
		[]string{"// a brief doc"})

	eq("javadoc block",
		got("/**\n * A documented function.\n *\n * It adds two integers.\n */"),
		[]string{"// A documented function.", "//", "// It adds two integers."})

	eq("block trims blank edges",
		got("/*\n\n content \n\n */"),
		[]string{"// content"})
}
