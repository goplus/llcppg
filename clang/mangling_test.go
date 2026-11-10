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
	"os"
	"path/filepath"
	"reflect"
	"testing"

	lc "github.com/llarhub/clang-c"
)

func TestMangling(t *testing.T) {
	header := filepath.Join(t.TempDir(), "in.h")
	if err := os.WriteFile(header, []byte(`
void normal(void);
void _leading(void);
extern int global;
extern int _leading_global;
`), 0600); err != nil {
		t.Fatal(err)
	}
	for _, language := range []string{"c", "c++"} {
		t.Run(language, func(t *testing.T) {
			want := map[string]string{
				"normal":          "normal",
				"_leading":        "_leading",
				"global":          "global",
				"_leading_global": "_leading_global",
			}
			if language == "c++" {
				want["normal"] = "_Z6normalv"
				want["_leading"] = "_Z8_leadingv"
			}
			for _, target := range []string{
				"x86_64-apple-darwin",
				"arm64-apple-darwin",
				"x86_64-unknown-linux-gnu",
				"aarch64-unknown-linux-gnu",
				"aarch64-unknown-darwin",
				"aarch64-apple-macosx14.0",
				"aarch64-apple-macos14.0",
				"aarch64-apple-ios17.0",
				"aarch64-apple-tvos17.0",
				"armv7k-apple-watchos5.0",
				"aarch64-apple-xros1.0",
				"aarch64-apple-visionos1.0",
				"x86_64-apple-driverkit20.0",
				"aarch64",
			} {
				t.Run(target, func(t *testing.T) {
					index := CreateIndex(1, 1)
					defer index.Dispose()
					unit, err := index.ParseTranslationUnit(0, header, "-x", language, "-target", target)
					if err != nil {
						t.Fatal(err)
					}
					defer unit.Dispose()
					if n := unit.NumDiagnostics(); n != 0 {
						t.Fatalf("unexpected Clang diagnostics: %d", n)
					}
					got := make(map[string]string)
					VisitChildren(unit.Cursor(), func(cursor, parent Cursor) ChildVisitResult {
						switch cursor.Kind {
						case lc.Cursor_FunctionDecl, lc.Cursor_VarDecl:
							got[String(cursor)] = Mangling(cursor)
						}
						return Continue
					})
					if !reflect.DeepEqual(got, want) {
						t.Errorf("Mangling() = %v; want %v", got, want)
					}
				})
			}
		})
	}
}
