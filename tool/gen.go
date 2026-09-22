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

package tool

import (
	"github.com/goplus/llcppg/cl"
	"github.com/goplus/llcppg/clang"
)

// -----------------------------------------------------------------------------

type Config struct {
	LLGoPackage    string `json:"LLGoPackage"`
	WrapFileHeader string `json:"WrapFileHeader"`
	CFlags         string `json:"CFlags"`
	Dir            string `json:"Dir"` // dir or dir/... (recursive)
}

// -----------------------------------------------------------------------------

// LoadSources parses the given source files and returns the translation units corresponding
// to those files.
func LoadSources(index clang.Index, headerFiles []string, args ...string) (files []cl.Source) {
	return index.ParseTranslationUnits(clang.DetailedPreprocessingRecord, headerFiles, args...)
}

// -----------------------------------------------------------------------------
