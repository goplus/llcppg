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

package tool_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/goplus/llcppg/clang"
	"github.com/goplus/llcppg/tool"
)

func TestNewPackage(t *testing.T) {
	idx := clang.CreateIndex(1, 1)
	cfg := tool.Config{}
	_, _, err := cfg.NewPackage("", "foo", ".", ".", idx)
	if err == nil || err.Error() != `invalid language: ""` {
		t.Fatal("cfg.NewPackage:", err)
	}

	cfg = tool.Config{
		Language: "c",
		Dir:      ".",
		GoFiles:  []string{"unknown.go"},
	}
	_, _, err = cfg.NewPackage("", "foo", ".", ".", idx)
	if err == nil || err.Error() != `open unknown.go: no such file or directory` {
		t.Fatal("cfg.NewPackage:", err)
	}
}

func TestPubImportNamingConfiguration(t *testing.T) {
	pkgDir, err := filepath.Abs("_testc/clang-c-22.1.8")
	if err != nil {
		t.Fatal(err)
	}
	headerDir := filepath.Join(pkgDir, "include/clang-c") + string(os.PathSeparator)
	idx := clang.CreateIndex(1, 1)
	defer idx.Dispose()

	for _, tt := range []struct {
		name  string
		setup func(*tool.Config)
	}{
		{"type_suffix", func(cfg *tool.Config) { cfg.TypeSuffix = []string{"_t"} }},
		{"rename", func(cfg *tool.Config) {
			cfg.Rename["size_t"] = "LocalSize"
			cfg.Rename["time_t"] = "LocalTime"
		}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := tool.LoadConf(filepath.Join(pkgDir, "llcppg.cfg"))
			if err != nil {
				t.Fatal(err)
			}
			tt.setup(&cfg)
			// The existing golden must retain File.Time() cstdlib.TimeT
			// and File.Name() CXString.String with either importer setting.
			testSingleFile(t, idx, pkgDir, headerDir, headerDir+"CXFile.h", &cfg)
		})
	}
}
