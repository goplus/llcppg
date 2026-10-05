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
	"reflect"
	"testing"
)

func TestToLineComments(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want []string
	}{
		{
			name: "asterisk banner delimiters",
			raw: `/**************************************************************************
Symbols and macros.

Please preserve this description.
**************************************************************************/`,
			want: []string{
				"// Symbols and macros.",
				"//",
				"// Please preserve this description.",
			},
		},
		{
			name: "bang block comment",
			raw:  "/*! A documented declaration. */",
			want: []string{"// A documented declaration."},
		},
		{
			name: "doxygen block comment",
			raw:  "/**\n * A documented declaration.\n */",
			want: []string{"// A documented declaration."},
		},
		{
			name: "doxygen line comment",
			raw:  "/// A documented declaration.",
			want: []string{"// A documented declaration."},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotComments := toLineComments(tt.raw)
			got := make([]string, len(gotComments))
			for i, comment := range gotComments {
				got[i] = comment.Text
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toLineComments() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
