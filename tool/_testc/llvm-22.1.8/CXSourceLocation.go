package clang

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

// Identifies a specific source location within a translation
// unit.
//
// Use clang_getExpansionLocation() or clang_getSpellingLocation()
// to map a source location to a particular file, line, and column.
type SourceLocation struct {
	PtrData [2]unsafe.Pointer
	IntData c.Uint
}

// Identifies a half-open character range in the source code.
//
// Use clang_getRangeStart() and clang_getRangeEnd() to retrieve the
// starting and end locations from a source range, respectively.
type SourceRange struct {
	PtrData      [2]unsafe.Pointer
	BeginIntData c.Uint
	EndIntData   c.Uint
}

// Identifies an array of ranges.
type SourceRangeList struct {
	Count  c.Uint
	Ranges *SourceRange
}

// Retrieve a NULL (invalid) source location.
//
//go:linkname GetNullLocation C.clang_getNullLocation
func GetNullLocation() SourceLocation

// Determine whether two source locations, which must refer into
// the same translation unit, refer to exactly the same point in the source
// code.
//
// \returns non-zero if the source locations refer to the same location, zero
// if they refer to different locations.
//
//go:linkname EqualLocations C.clang_equalLocations
func EqualLocations(loc1 SourceLocation, loc2 SourceLocation) c.Uint

// Determine for two source locations if the first comes
// strictly before the second one in the source code.
//
// \returns non-zero if the first source location comes
// strictly before the second one, zero otherwise.
//
//go:linkname IsBeforeInTranslationUnit C.clang_isBeforeInTranslationUnit
func IsBeforeInTranslationUnit(loc1 SourceLocation, loc2 SourceLocation) c.Uint

// Returns non-zero if the given source location is in a system header.
//
// llgo:link SourceLocation.IsInSystemHeader C.clang_Location_isInSystemHeader
func (location SourceLocation) IsInSystemHeader() c.Int {
	return 0
}

// Returns non-zero if the given source location is in the main file of
// the corresponding translation unit.
//
// llgo:link SourceLocation.IsFromMainFile C.clang_Location_isFromMainFile
func (location SourceLocation) IsFromMainFile() c.Int {
	return 0
}

// Retrieve a NULL (invalid) source range.
//
//go:linkname GetNullRange C.clang_getNullRange
func GetNullRange() SourceRange

// Retrieve a source range given the beginning and ending source
// locations.
//
//go:linkname GetRange C.clang_getRange
func GetRange(begin SourceLocation, end SourceLocation) SourceRange

// Determine whether two ranges are equivalent.
//
// \returns non-zero if the ranges are the same, zero if they differ.
//
//go:linkname EqualRanges C.clang_equalRanges
func EqualRanges(range1 SourceRange, range2 SourceRange) c.Uint

// Returns non-zero if \p range is null.
//
// llgo:link SourceRange.IsNull C.clang_Range_isNull
func (range_ SourceRange) IsNull() c.Int {
	return 0
}

// Retrieve the file, line, column, and offset represented by
// the given source location.
//
// If the location refers into a macro expansion, retrieves the
// location of the macro expansion.
//
// \param location the location within a source file that will be decomposed
// into its parts.
//
// \param file [out] if non-NULL, will be set to the file to which the given
// source location points.
//
// \param line [out] if non-NULL, will be set to the line to which the given
// source location points.
//
// \param column [out] if non-NULL, will be set to the column to which the given
// source location points.
//
// \param offset [out] if non-NULL, will be set to the offset into the
// buffer to which the given source location points.
//
// llgo:link SourceLocation.Expansion C.clang_getExpansionLocation
func (location SourceLocation) Expansion(file *File, line *c.Uint, column *c.Uint, offset *c.Uint) {
}

// Retrieve the file, line and column represented by the given source
// location, as specified in a # line directive.
//
// Example: given the following source code in a file somefile.c
//
// \code
// #123 "dummy.c" 1
//
// static int func(void)
// {
//     return 0;
// }
// \endcode
//
// the location information returned by this function would be
//
// File: dummy.c Line: 124 Column: 12
//
// whereas clang_getExpansionLocation would have returned
//
// File: somefile.c Line: 3 Column: 12
//
// \param location the location within a source file that will be decomposed
// into its parts.
//
// \param filename [out] if non-NULL, will be set to the filename of the
// source location. Note that filenames returned will be for "virtual" files,
// which don't necessarily exist on the machine running clang - e.g. when
// parsing preprocessed output obtained from a different environment. If
// a non-NULL value is passed in, remember to dispose of the returned value
// using \c clang_disposeString() once you've finished with it. For an invalid
// source location, an empty string is returned.
//
// \param line [out] if non-NULL, will be set to the line number of the
// source location. For an invalid source location, zero is returned.
//
// \param column [out] if non-NULL, will be set to the column number of the
// source location. For an invalid source location, zero is returned.
//
// llgo:link SourceLocation.Presumed C.clang_getPresumedLocation
func (location SourceLocation) Presumed(filename *String, line *c.Uint, column *c.Uint) {
}

// Legacy API to retrieve the file, line, column, and offset represented
// by the given source location.
//
// This interface has been replaced by the newer interface
// #clang_getExpansionLocation(). See that interface's documentation for
// details.
//
// llgo:link SourceLocation.Instantiation C.clang_getInstantiationLocation
func (location SourceLocation) Instantiation(file *File, line *c.Uint, column *c.Uint, offset *c.Uint) {
}

// Retrieve the file, line, column, and offset represented by
// the given source location.
//
// If the location refers into a macro instantiation, return where the
// location was originally spelled in the source file.
//
// \param location the location within a source file that will be decomposed
// into its parts.
//
// \param file [out] if non-NULL, will be set to the file to which the given
// source location points.
//
// \param line [out] if non-NULL, will be set to the line to which the given
// source location points.
//
// \param column [out] if non-NULL, will be set to the column to which the given
// source location points.
//
// \param offset [out] if non-NULL, will be set to the offset into the
// buffer to which the given source location points.
//
// llgo:link SourceLocation.Spelling C.clang_getSpellingLocation
func (location SourceLocation) Spelling(file *File, line *c.Uint, column *c.Uint, offset *c.Uint) {
}

// Retrieve the file, line, column, and offset represented by
// the given source location.
//
// If the location refers into a macro expansion, return where the macro was
// expanded or where the macro argument was written, if the location points at
// a macro argument.
//
// \param location the location within a source file that will be decomposed
// into its parts.
//
// \param file [out] if non-NULL, will be set to the file to which the given
// source location points.
//
// \param line [out] if non-NULL, will be set to the line to which the given
// source location points.
//
// \param column [out] if non-NULL, will be set to the column to which the given
// source location points.
//
// \param offset [out] if non-NULL, will be set to the offset into the
// buffer to which the given source location points.
//
// llgo:link SourceLocation.File C.clang_getFileLocation
func (location SourceLocation) File(file *File, line *c.Uint, column *c.Uint, offset *c.Uint) {
}

// Retrieve a source location representing the first character within a
// source range.
//
// llgo:link SourceRange.Start C.clang_getRangeStart
func (range_ SourceRange) Start() SourceLocation {
	return SourceLocation{}
}

// Retrieve a source location representing the last character within a
// source range.
//
// llgo:link SourceRange.End C.clang_getRangeEnd
func (range_ SourceRange) End() SourceLocation {
	return SourceLocation{}
}

// Destroy the given \c CXSourceRangeList.
//
// llgo:link (*SourceRangeList).Dispose C.clang_disposeSourceRangeList
func (ranges *SourceRangeList) Dispose() {
}
