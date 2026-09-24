package CXDiagnostic

import (
	"clang/CXSourceLocation"
	"clang/CXString"
	"github.com/goplus/lib/c"
	_ "unsafe"
)

const LLGoPackage = "link: -L$(llvm-config --libdir) -lclang; -lclang"

// Describes the severity of a particular diagnostic.
type DiagnosticSeverity c.Int

const (
	Diagnostic_Ignored DiagnosticSeverity = 0
	Diagnostic_Note    DiagnosticSeverity = 1
	Diagnostic_Warning DiagnosticSeverity = 2
	Diagnostic_Error   DiagnosticSeverity = 3
	Diagnostic_Fatal   DiagnosticSeverity = 4
)

// A single diagnostic, containing the diagnostic's severity,
// location, text, source ranges, and fix-it hints.
type Diagnostic uintptr

// A group of CXDiagnostics.
type DiagnosticSet uintptr

// Describes the kind of error that occurred (if any) in a call to
// \c clang_loadDiagnostics.
type LoadDiag_Error c.Int

const (
	LoadDiag_None        LoadDiag_Error = 0
	LoadDiag_Unknown     LoadDiag_Error = 1
	LoadDiag_CannotLoad  LoadDiag_Error = 2
	LoadDiag_InvalidFile LoadDiag_Error = 3
)

// Options to control the display of diagnostics.
//
// The values in this enum are meant to be combined to customize the
// behavior of \c clang_formatDiagnostic().
type DiagnosticDisplayOptions c.Int

const (
	Diagnostic_DisplaySourceLocation DiagnosticDisplayOptions = 1
	Diagnostic_DisplayColumn         DiagnosticDisplayOptions = 2
	Diagnostic_DisplaySourceRanges   DiagnosticDisplayOptions = 4
	Diagnostic_DisplayOption         DiagnosticDisplayOptions = 8
	Diagnostic_DisplayCategoryId     DiagnosticDisplayOptions = 16
	Diagnostic_DisplayCategoryName   DiagnosticDisplayOptions = 32
)

// Determine the number of diagnostics in a CXDiagnosticSet.
//
// llgo:link DiagnosticSet.NumDiagnosticsInSet C.clang_getNumDiagnosticsInSet
func (Diags DiagnosticSet) NumDiagnosticsInSet() c.Uint {
	return 0
}

// Retrieve a diagnostic associated with the given CXDiagnosticSet.
//
// \param Diags the CXDiagnosticSet to query.
// \param Index the zero-based diagnostic number to retrieve.
//
// \returns the requested diagnostic. This diagnostic must be freed
// via a call to \c clang_disposeDiagnostic().
//
// llgo:link DiagnosticSet.DiagnosticInSet C.clang_getDiagnosticInSet
func (Diags DiagnosticSet) DiagnosticInSet(Index c.Uint) Diagnostic {
	return 0
}

// Deserialize a set of diagnostics from a Clang diagnostics bitcode
// file.
//
// \param file The name of the file to deserialize.
// \param error A pointer to a enum value recording if there was a problem
//        deserializing the diagnostics.
// \param errorString A pointer to a CXString for recording the error string
//        if the file was not successfully loaded.
//
// \returns A loaded CXDiagnosticSet if successful, and NULL otherwise.  These
// diagnostics should be released using clang_disposeDiagnosticSet().
//
//go:linkname LoadDiagnostics C.clang_loadDiagnostics
func LoadDiagnostics(file *c.Char, error *LoadDiag_Error, errorString *CXString.String) DiagnosticSet

// Release a CXDiagnosticSet and all of its contained diagnostics.
//
// llgo:link DiagnosticSet.Dispose C.clang_disposeDiagnosticSet
func (Diags DiagnosticSet) Dispose() {
}

// Retrieve the child diagnostics of a CXDiagnostic.
//
// This CXDiagnosticSet does not need to be released by
// clang_disposeDiagnosticSet.
//
// llgo:link Diagnostic.ChildDiagnostics C.clang_getChildDiagnostics
func (D Diagnostic) ChildDiagnostics() DiagnosticSet {
	return 0
}

// Destroy a diagnostic.
//
// llgo:link Diagnostic.Dispose C.clang_disposeDiagnostic
func (Diagnostic Diagnostic) Dispose() {
}

// Format the given diagnostic in a manner that is suitable for display.
//
// This routine will format the given diagnostic to a string, rendering
// the diagnostic according to the various options given. The
// \c clang_defaultDiagnosticDisplayOptions() function returns the set of
// options that most closely mimics the behavior of the clang compiler.
//
// \param Diagnostic The diagnostic to print.
//
// \param Options A set of options that control the diagnostic display,
// created by combining \c CXDiagnosticDisplayOptions values.
//
// \returns A new string containing for formatted diagnostic.
//
// llgo:link Diagnostic.Format C.clang_formatDiagnostic
func (Diagnostic Diagnostic) Format(Options c.Uint) CXString.String {
	return CXString.String{}
}

// Retrieve the set of display options most similar to the
// default behavior of the clang compiler.
//
// \returns A set of display options suitable for use with \c
// clang_formatDiagnostic().
//
//go:linkname DefaultDiagnosticDisplayOptions C.clang_defaultDiagnosticDisplayOptions
func DefaultDiagnosticDisplayOptions() c.Uint

// Determine the severity of the given diagnostic.
//
// llgo:link Diagnostic.Severity C.clang_getDiagnosticSeverity
func (_llcppg_param1 Diagnostic) Severity() DiagnosticSeverity {
	return 0
}

// Retrieve the source location of the given diagnostic.
//
// This location is where Clang would print the caret ('^') when
// displaying the diagnostic on the command line.
//
// llgo:link Diagnostic.Location C.clang_getDiagnosticLocation
func (_llcppg_param1 Diagnostic) Location() CXSourceLocation.SourceLocation {
	return CXSourceLocation.SourceLocation{}
}

// Retrieve the text of the given diagnostic.
//
// llgo:link Diagnostic.Spelling C.clang_getDiagnosticSpelling
func (_llcppg_param1 Diagnostic) Spelling() CXString.String {
	return CXString.String{}
}

// Retrieve the name of the command-line option that enabled this
// diagnostic.
//
// \param Diag The diagnostic to be queried.
//
// \param Disable If non-NULL, will be set to the option that disables this
// diagnostic (if any).
//
// \returns A string that contains the command-line option used to enable this
// warning, such as "-Wconversion" or "-pedantic".
//
// llgo:link Diagnostic.Option C.clang_getDiagnosticOption
func (Diag Diagnostic) Option(Disable *CXString.String) CXString.String {
	return CXString.String{}
}

// Retrieve the category number for this diagnostic.
//
// Diagnostics can be categorized into groups along with other, related
// diagnostics (e.g., diagnostics under the same warning flag). This routine
// retrieves the category number for the given diagnostic.
//
// \returns The number of the category that contains this diagnostic, or zero
// if this diagnostic is uncategorized.
//
// llgo:link Diagnostic.Category C.clang_getDiagnosticCategory
func (_llcppg_param1 Diagnostic) Category() c.Uint {
	return 0
}

// Retrieve the name of a particular diagnostic category.  This
//  is now deprecated.  Use clang_getDiagnosticCategoryText()
//  instead.
//
// \param Category A diagnostic category number, as returned by
// \c clang_getDiagnosticCategory().
//
// \returns The name of the given diagnostic category.
//
//go:linkname GetDiagnosticCategoryName C.clang_getDiagnosticCategoryName
func GetDiagnosticCategoryName(Category c.Uint) CXString.String

// Retrieve the diagnostic category text for a given diagnostic.
//
// \returns The text of the given diagnostic category.
//
// llgo:link Diagnostic.CategoryText C.clang_getDiagnosticCategoryText
func (_llcppg_param1 Diagnostic) CategoryText() CXString.String {
	return CXString.String{}
}

// Determine the number of source ranges associated with the given
// diagnostic.
//
// llgo:link Diagnostic.NumRanges C.clang_getDiagnosticNumRanges
func (_llcppg_param1 Diagnostic) NumRanges() c.Uint {
	return 0
}

// Retrieve a source range associated with the diagnostic.
//
// A diagnostic's source ranges highlight important elements in the source
// code. On the command line, Clang displays source ranges by
// underlining them with '~' characters.
//
// \param Diagnostic the diagnostic whose range is being extracted.
//
// \param Range the zero-based index specifying which range to
//
// \returns the requested source range.
//
// llgo:link Diagnostic.Range C.clang_getDiagnosticRange
func (Diagnostic Diagnostic) Range(Range c.Uint) CXSourceLocation.SourceRange {
	return CXSourceLocation.SourceRange{}
}

// Determine the number of fix-it hints associated with the
// given diagnostic.
//
// llgo:link Diagnostic.NumFixIts C.clang_getDiagnosticNumFixIts
func (Diagnostic Diagnostic) NumFixIts() c.Uint {
	return 0
}

// Retrieve the replacement information for a given fix-it.
//
// Fix-its are described in terms of a source range whose contents
// should be replaced by a string. This approach generalizes over
// three kinds of operations: removal of source code (the range covers
// the code to be removed and the replacement string is empty),
// replacement of source code (the range covers the code to be
// replaced and the replacement string provides the new code), and
// insertion (both the start and end of the range point at the
// insertion location, and the replacement string provides the text to
// insert).
//
// \param Diagnostic The diagnostic whose fix-its are being queried.
//
// \param FixIt The zero-based index of the fix-it.
//
// \param ReplacementRange The source range whose contents will be
// replaced with the returned replacement string. Note that source
// ranges are half-open ranges [a, b), so the source code should be
// replaced from a and up to (but not including) b.
//
// \returns A string containing text that should be replace the source
// code indicated by the \c ReplacementRange.
//
// llgo:link Diagnostic.FixIt C.clang_getDiagnosticFixIt
func (Diagnostic Diagnostic) FixIt(FixIt c.Uint, ReplacementRange *CXSourceLocation.SourceRange) CXString.String {
	return CXString.String{}
}
