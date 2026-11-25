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

package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMergeCLI(t *testing.T) {
	const root = "../../tool/merge/_testmerge"
	for _, fixture := range []string{"dlfcn", "locale_single", "comment_edges", "invalid_package_names"} {
		t.Run(fixture, func(t *testing.T) {
			destDir := t.TempDir()
			inDir := filepath.Join(root, fixture, "in")
			outDir := filepath.Join(root, fixture, "out")
			mergeErr := mergePkg(destDir, inDir)
			if want, err := os.ReadFile(filepath.Join(outDir, "error")); err == nil {
				if mergeErr == nil || !strings.Contains(mergeErr.Error(), strings.TrimSpace(string(want))) {
					t.Fatalf("merge error = %v, want %s", mergeErr, want)
				}
				return
			} else if !os.IsNotExist(err) {
				t.Fatal(err)
			}
			if mergeErr != nil {
				t.Fatal(mergeErr)
			}

			goldens, err := filepath.Glob(filepath.Join(outDir, "*.merge.go"))
			if err != nil {
				t.Fatal(err)
			}
			files, err := os.ReadDir(destDir)
			if err != nil {
				t.Fatal(err)
			}
			if len(files) != len(goldens) {
				t.Errorf("got %d files, want %d", len(files), len(goldens))
			}
			for _, golden := range goldens {
				name := strings.TrimSuffix(filepath.Base(golden), ".merge.go") + ".go"
				want, err := os.ReadFile(golden)
				if err != nil {
					t.Fatal(err)
				}
				got, err := os.ReadFile(filepath.Join(destDir, name))
				if err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(got, want) {
					t.Errorf("%s: unexpected output\n%s\nwant:\n%s", name, got, want)
				}
			}
		})
	}
}
