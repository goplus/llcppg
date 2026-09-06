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

type stringer interface {
	String() clang.String
}

// String returns the string representation of the given stringer.
func String[T stringer](v T) string {
	return clang.GoString(v.String())
}

// -----------------------------------------------------------------------------

/**
 * An "index" that consists of a set of translation units that would
 * typically be linked together into an executable or library.
 */
type Index struct {
	impl *clang.Index
}

/**
 * Provides a shared context for creating translation units.
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

/**
 * Retrieve the cursor that represents the given translation unit.
 *
 * The translation unit cursor can be used to start traversing the
 * various declarations within the given translation unit.
 */
func (u TranslationUnit) Cursor() Cursor {
	return u.impl.Cursor()
}

// -----------------------------------------------------------------------------

/**
 * A cursor representing some element in the abstract syntax tree for
 * a translation unit.
 *
 * The cursor abstraction unifies the different kinds of entities in a
 * program--declaration, statements, expressions, references to declarations,
 * etc.--under a single "cursor" abstraction with a common set of operations.
 * Common operation for a cursor include: getting the physical location in
 * a source file where the cursor points, getting the name associated with a
 * cursor, and retrieving cursors for any child nodes of a particular cursor.
 *
 * Cursors can be produced in two specific ways.
 * clang_getTranslationUnitCursor() produces a cursor for a translation unit,
 * from which one can use clang_visitChildren() to explore the rest of the
 * translation unit. clang_getCursor() maps from a physical source location
 * to the entity that resides at that location, allowing one to map from the
 * source code into the AST.
 */
type Cursor = clang.Cursor

/**
 * Describes how the traversal of the children of a particular
 * cursor should proceed after visiting a particular child cursor.
 */
type ChildVisitResult = clang.ChildVisitResult

const (
	/**
	 * Terminates the cursor traversal.
	 */
	Break ChildVisitResult = clang.ChildVisit_Break

	/**
	 * Continues the cursor traversal with the next sibling of
	 * the cursor just visited, without visiting its children.
	 */
	Continue ChildVisitResult = clang.ChildVisit_Continue

	/**
	 * Recursively traverse the children of this cursor, using
	 * the same visitor and client data.
	 */
	Recurse ChildVisitResult = clang.ChildVisit_Recurse
)

// VisitChildren visits the children of a particular cursor, invoking the given
// visitor function for each child cursor. The traversal may be recursive,
// depending on the return value of the visitor function.
func VisitChildren(root Cursor, fn func(cur, parent Cursor) ChildVisitResult) uint {
	return uint(clang.VisitChildren(
		root, func(cur, parent Cursor, param clang.ClientData) ChildVisitResult {
			return c.GoClosure[func(cur, parent Cursor) ChildVisitResult](param)(cur, parent)
		}, c.ClosureData(fn)))
}

// -----------------------------------------------------------------------------
