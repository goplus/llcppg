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
	"os"
	"path"
	"strings"
	"testing"
)

// -----------------------------------------------------------------------------

func DoTestFromDir(t *testing.T, sel, relDir string, testFunc func(t *testing.T, pkgDir string)) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal("Getwd failed:", err)
	}
	dir = path.Join(dir, relDir)
	fis, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal("ReadDir failed:", err)
	}
	for _, fi := range fis {
		name := fi.Name()
		if strings.HasPrefix(name, "_") {
			continue
		}
		t.Run(name, func(t *testing.T) {
			pkgDir := dir + "/" + name
			if sel != "" && !strings.Contains(pkgDir, sel) {
				return
			}
			testFunc(t, pkgDir)
		})
	}
}

// -----------------------------------------------------------------------------
/*
func testFromDir(t *testing.T, sel, relDir string) {
	DoTestFromDir(t, sel, relDir, func(t *testing.T, pkgDir string) {
	})
}
*/
// -----------------------------------------------------------------------------
