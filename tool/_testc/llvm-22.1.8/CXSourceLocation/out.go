package CXSourceLocation

import (
	"clang/CXFile"
	"clang/CXString"
	"github.com/goplus/lib/c"
	"unsafe"
)

const LLGoPackage = "link: -L$(llvm-config --libdir) -lclang; -lclang"

type SourceLocation struct {
	PtrData [2]unsafe.Pointer
	IntData c.Uint
}
type SourceRange struct {
	PtrData      [2]unsafe.Pointer
	BeginIntData c.Uint
	EndIntData   c.Uint
}
type SourceRangeList struct {
	Count  c.Uint
	Ranges *SourceRange
}

//go:linkname GetNullLocation C.clang_getNullLocation
func GetNullLocation() SourceLocation

//go:linkname EqualLocations C.clang_equalLocations
func EqualLocations(loc1 SourceLocation, loc2 SourceLocation) c.Uint

//go:linkname IsBeforeInTranslationUnit C.clang_isBeforeInTranslationUnit
func IsBeforeInTranslationUnit(loc1 SourceLocation, loc2 SourceLocation) c.Uint

// llgo:link SourceLocation.LocationIsInSystemHeader C.clang_Location_isInSystemHeader
func (location SourceLocation) LocationIsInSystemHeader() c.Int {
	return 0
}

// llgo:link SourceLocation.LocationIsFromMainFile C.clang_Location_isFromMainFile
func (location SourceLocation) LocationIsFromMainFile() c.Int {
	return 0
}

//go:linkname GetNullRange C.clang_getNullRange
func GetNullRange() SourceRange

//go:linkname GetRange C.clang_getRange
func GetRange(begin SourceLocation, end SourceLocation) SourceRange

//go:linkname EqualRanges C.clang_equalRanges
func EqualRanges(range1 SourceRange, range2 SourceRange) c.Uint

// llgo:link SourceRange.RangeIsNull C.clang_Range_isNull
func (range_ SourceRange) RangeIsNull() c.Int {
	return 0
}

// llgo:link SourceLocation.ExpansionLocation C.clang_getExpansionLocation
func (location SourceLocation) ExpansionLocation(file *CXFile.File, line *c.Uint, column *c.Uint, offset *c.Uint) {
}

// llgo:link SourceLocation.PresumedLocation C.clang_getPresumedLocation
func (location SourceLocation) PresumedLocation(filename *CXString.String, line *c.Uint, column *c.Uint) {
}

// llgo:link SourceLocation.InstantiationLocation C.clang_getInstantiationLocation
func (location SourceLocation) InstantiationLocation(file *CXFile.File, line *c.Uint, column *c.Uint, offset *c.Uint) {
}

// llgo:link SourceLocation.SpellingLocation C.clang_getSpellingLocation
func (location SourceLocation) SpellingLocation(file *CXFile.File, line *c.Uint, column *c.Uint, offset *c.Uint) {
}

// llgo:link SourceLocation.FileLocation C.clang_getFileLocation
func (location SourceLocation) FileLocation(file *CXFile.File, line *c.Uint, column *c.Uint, offset *c.Uint) {
}

// llgo:link SourceRange.RangeStart C.clang_getRangeStart
func (range_ SourceRange) RangeStart() SourceLocation {
	return SourceLocation{}
}

// llgo:link SourceRange.RangeEnd C.clang_getRangeEnd
func (range_ SourceRange) RangeEnd() SourceLocation {
	return SourceLocation{}
}

// llgo:link (*SourceRangeList).Dispose C.clang_disposeSourceRangeList
func (ranges *SourceRangeList) Dispose() {
}
