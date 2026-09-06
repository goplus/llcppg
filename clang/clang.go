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
	"unsafe"

	"github.com/goplus/lib/c"
	"github.com/goplus/llcppg/lib/clang"
)

// -----------------------------------------------------------------------------

/**
 * An "index" that consists of a set of translation units that would
 * typically be linked together into an executable or library.
 */
type Index struct {
	impl *clang.Index
}

/**
 * CreateIndex provides a shared context for creating translation units.
 *
 * It provides two options:
 *
 * - excludeDeclarationsFromPCH: When non-zero, allows enumeration of "local"
 * declarations (when loading any new translation units). A "local" declaration
 * is one that belongs in the translation unit itself and not in a precompiled
 * header that was used by the translation unit. If zero, all declarations
 * will be enumerated.
 *
 * This process of creating the 'pch', loading it separately, and using it (via
 * -include-pch) allows 'excludeDeclsFromPCH' to remove redundant callbacks
 * (which gives the indexer the same performance benefit as the compiler).
 */
func CreateIndex(excludeDeclarationsFromPCH, displayDiagnostics int) Index {
	return Index{impl: clang.CreateIndex(c.Int(excludeDeclarationsFromPCH), c.Int(displayDiagnostics))}
}

/**
 * Destroy the given index.
 *
 * The index must not be destroyed until all of the translation units created
 * within that index have been destroyed.
 */
func (i Index) Dispose() {
	i.impl.Dispose()
}

// ParseTranslationUnit parses the given source file and returns the translation unit corresponding
// to that file.
func (i Index) ParseTranslationUnit(options uint, filename string, args ...string) TranslationUnit {
	cArgs := make([]*c.Char, len(args))
	for i, arg := range args {
		cArgs[i] = c.AllocaCStr(arg)
	}
	return TranslationUnit{
		impl: i.impl.ParseTranslationUnit(
			c.AllocaCStr(filename), unsafe.SliceData(cArgs), c.Int(len(cArgs)), nil, 0, c.Uint(options)),
	}
}

// -----------------------------------------------------------------------------

/**
 * A single translation unit, which resides in an index.
 */
type TranslationUnit struct {
	impl *clang.TranslationUnit
}

/**
 * Destroy the specified TranslationUnit object.
 */
func (u TranslationUnit) Dispose() {
	u.impl.Dispose()
}

// -----------------------------------------------------------------------------
