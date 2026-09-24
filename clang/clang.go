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
	"iter"
	"unsafe"

	"github.com/goplus/lib/c"
	"github.com/goplus/llcppg/lib/clang"
)

// -----------------------------------------------------------------------------

func GoStringAndDispose(str clang.String) string {
	s := c.GoString(str.CStr())
	str.Dispose()
	return s
}

type stringer interface {
	String() clang.String
}

// String returns the Go string of a value whose String() returns a clang String.
func String[T stringer](v T) string {
	return GoStringAndDispose(v.String())
}

// -----------------------------------------------------------------------------

/**
 * An "index" that consists of a set of translation units that would
 * typically be linked together into an executable or library.
 */
type Index struct {
	*clang.Index
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
	return Index{
		Index: clang.CreateIndex(c.Int(excludeDeclarationsFromPCH), c.Int(displayDiagnostics)),
	}
}

/**
 * Flags that control the creation of translation units.
 *
 * The enumerators in this enumeration type are meant to be bitwise
 * ORed together to specify which options should be used when
 * constructing the translation unit.
 */
const (
	/**
	 * Used to indicate that the parser should construct a "detailed"
	 * preprocessing record, including all macro definitions and instantiations.
	 *
	 * Constructing a detailed preprocessing record requires more memory
	 * and time to parse, since the information contained in the record
	 * is usually not retained. However, it can be useful for
	 * applications that require more detailed information about the
	 * behavior of the preprocessor.
	 */
	DetailedPreprocessingRecord = clang.DetailedPreprocessingRecord
)

// ParseTranslationUnit parses the given source file and returns the translation unit corresponding
// to that file.
func (i Index) ParseTranslationUnit(options uint, filename string, args ...string) TranslationUnit {
	cArgs := make([]*c.Char, len(args))
	for i, arg := range args {
		cArgs[i] = c.AllocaCStr(arg)
	}
	return TranslationUnit{
		TranslationUnit: i.Index.ParseTranslationUnit(
			c.AllocaCStr(filename), unsafe.SliceData(cArgs), c.Int(len(cArgs)), nil, 0, c.Uint(options)),
	}
}

// ParseTranslationUnits parses the given source files and returns the translation units corresponding
// to those files.
func (i Index) ParseTranslationUnits(options uint, filenames []string, args ...string) iter.Seq[TranslationUnit] {
	cArgs := make([]*c.Char, len(args))
	for i, arg := range args {
		cArgs[i] = c.AllocaCStr(arg)
	}
	return func(yield func(TranslationUnit) bool) {
		for _, filename := range filenames {
			u := TranslationUnit{
				TranslationUnit: i.Index.ParseTranslationUnit(
					c.AllocaCStr(filename), unsafe.SliceData(cArgs), c.Int(len(cArgs)), nil, 0, c.Uint(options)),
			}
			if !yield(u) {
				return
			}
		}
	}
}

// -----------------------------------------------------------------------------

/**
 * A particular source file that is part of a translation unit.
 */
type File = clang.File

const (
	// NULL file handle (invalid/absent File)
	InvalidFile = File(0)
)

/**
 * Retrieve the name of a particular source file.
 */
func FileName(f clang.File) string {
	return GoStringAndDispose(f.FileName())
}

// -----------------------------------------------------------------------------

// Diagnostic represents a diagnostic message, such as a compiler warning or error.
type Diagnostic struct {
	*clang.Diagnostic
}

/**
 * Returns a string that describes the diagnostic.
 */
func (e Diagnostic) String() string {
	return GoStringAndDispose(e.Diagnostic.String())
}

/**
 * Returns the category text for the given diagnostic.
 */
func (e Diagnostic) CategoryText() string {
	return GoStringAndDispose(e.Diagnostic.CategoryText())
}

/**
 * Format the given diagnostic according to the specified display options.
 */
func (e Diagnostic) Format(options clang.DiagnosticDisplayOptions) string {
	return GoStringAndDispose(e.Diagnostic.Format(options))
}

// -----------------------------------------------------------------------------

/**
 * A single translation unit, which resides in an index.
 */
type TranslationUnit struct {
	*clang.TranslationUnit
}

// File returns the File object corresponding to the given filename in the translation unit.
func (u TranslationUnit) File(filename string) File {
	return u.TranslationUnit.File(c.AllocaCStr(filename))
}

// FileContents returns the contents of the specified file in the translation unit.
func (u TranslationUnit) FileContents(file File) []byte {
	var size c.SizeT
	data := u.TranslationUnit.FileContents(file, &size)
	return unsafe.Slice((*byte)(unsafe.Pointer(data)), int(size))
}

// Tokenize tokenizes the source code described by the given source range
// into raw lexical tokens. Call dispose() to free the memory allocated for
// the tokens after use.
func (u TranslationUnit) Tokenize(extent clang.SourceRange) (ret []clang.Token, dispose func()) {
	var tokens *clang.Token
	var numTokens c.Uint
	u.TranslationUnit.Tokenize(extent, &tokens, &numTokens)
	ret = unsafe.Slice(tokens, int(numTokens))
	dispose = func() {
		u.TranslationUnit.DisposeTokens(tokens, numTokens)
	}
	return
}

/**
 * Determine the spelling of the given token.
 *
 * The spelling of a token is the textual representation of that token, e.g.,
 * the text of an identifier or keyword.
 */
func (u TranslationUnit) Token(tok clang.Token) string {
	return GoStringAndDispose(u.TranslationUnit.Token(tok))
}

/**
 * Retrieve the diagnostic associated with the given index in the translation unit.
 */
func (u TranslationUnit) Diagnostic(index c.Uint) (ret Diagnostic) {
	return Diagnostic{u.TranslationUnit.Diagnostic(index)}
}

// VisitDiagnostics visits all diagnostics in the translation unit and invokes the
// given function for each diagnostic.
func (u TranslationUnit) VisitDiagnostics(fn func(diag Diagnostic)) {
	numDiagnostics := u.NumDiagnostics()
	for i := range numDiagnostics {
		diag := u.Diagnostic(i)
		fn(diag)
		diag.Dispose()
	}
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
 * Identifies a specific source location within a translation
 * unit.
 *
 * Use clang_getExpansionLocation() or clang_getSpellingLocation()
 * to map a source location to a particular file, line, and column.
 */
type SourceLocation = clang.SourceLocation

// PresumedFile returns the presumed file name for the given source location.
func PresumedFile(loc SourceLocation) string {
	var filename clang.String
	loc.PresumedLocation(&filename, nil, nil)
	return GoStringAndDispose(filename)
}

/**
 * Retrieve the display name for the entity referenced by this cursor.
 *
 * The display name contains extra information that helps identify the cursor,
 * such as the parameters of a function or template or the arguments of a
 * class template specialization.
 */
func DisplayName(entity clang.Cursor) string {
	return GoStringAndDispose(entity.DisplayName())
}

// RawComment returns the raw documentation comment associated with the given
// cursor, including the comment markers (e.g. "/** ... */" or "///"). It returns
// an empty string if the cursor has no associated comment.
func RawComment(entity clang.Cursor) string {
	return GoStringAndDispose(entity.RawCommentText())
}

/**
 * Retrieve the translation unit that a cursor originated from.
 */
func TU(c Cursor) (ret TranslationUnit) {
	return TranslationUnit{TranslationUnit: c.TU()}
}

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

// OverriddenCursors returns the set of methods that the given method cursor
// immediately overrides, computed by libclang via clang_getOverriddenCursors.
//
// For a C++ virtual member function this is the set of virtual member functions
// with the same signature declared in its (immediate) base classes; with
// multiple inheritance a method can override several. libclang only reports the
// immediate overridden methods, so walking the override graph transitively
// requires calling this again on each result.
//
// The array libclang allocates is released before returning, so the returned
// slice is a freshly copied, caller-owned []Cursor (empty when there are none).
func OverriddenCursors(cur Cursor) (ret []Cursor, dispose func()) {
	var overridden *Cursor
	var num c.Uint
	cur.OverriddenCursors(&overridden, &num)
	ret = unsafe.Slice(overridden, int(num))
	dispose = func() {
		overridden.DisposeOverriddenCursors()
	}
	return
}

// -----------------------------------------------------------------------------
