package clang

import (
	"errors"
	"os/exec"
	"unsafe"

	"github.com/goplus/lib/c"
	"github.com/goplus/lib/c/clang"
)

const (
	LLGoPackage = "link: -L$(llvm-config --libdir) -lclang; -lclang"
)

type Config struct {
	File  string
	Temp  bool
	Args  []string
	IsCpp bool
	Index *clang.Index
}

type Visitor func(cursor, parent clang.Cursor) clang.ChildVisitResult

type InclusionVisitor func(included_file clang.File, inclusions []clang.SourceLocation)

const TEMP_FILE = "temp.h"

func CreateTranslationUnit(config *Config) (*clang.Index, *clang.TranslationUnit, error) {
	// default use the c/c++ standard of clang; c:gnu17 c++:gnu++17
	// https://clang.llvm.org/docs/CommandGuide/clang.html

	executableName := "clang"
	path, err := exec.LookPath(executableName)
	if err != nil {
		return nil, nil, err
	}

	allArgs := append(append([]string{path}, defaultArgs(config.IsCpp)...), config.Args...)
	cArgs := make([]*c.Char, len(allArgs))
	for i, arg := range allArgs {
		cArgs[i] = c.AllocaCStr(arg)
	}

	var index *clang.Index
	if config.Index != nil {
		index = config.Index
	} else {
		index = clang.CreateIndex(0, 0)
	}

	var unit *clang.TranslationUnit
	var code ErrorCode
	if config.Temp {
		content := c.AllocaCStr(config.File)
		tempFile := &clang.UnsavedFile{
			Filename: c.Str(TEMP_FILE),
			Contents: content,
			Length:   c.Ulong(c.Strlen(content)),
		}
		code = ParseTranslationUnit2FullArgv(index,
			tempFile.Filename,
			unsafe.SliceData(cArgs), c.Int(len(cArgs)),
			tempFile, 1,
			clang.DetailedPreprocessingRecord,
			&unit,
		)
	} else {

		cFile := c.AllocaCStr(config.File)
		code = ParseTranslationUnit2FullArgv(index,
			cFile,
			unsafe.SliceData(cArgs), c.Int(len(cArgs)),
			nil, 0,
			clang.DetailedPreprocessingRecord,
			&unit,
		)
	}

	if code != Error_Success {
		c.Printf(c.Str("code: %d\n"), code)
		return nil, nil, errors.New("failed to parse translation unit")
	}

	if unit == nil {
		return nil, nil, errors.New("failed to parse translation unit")
	}

	return index, unit, nil
}

func GetLocation(loc clang.SourceLocation) (file clang.File, line c.Uint, column c.Uint, offset c.Uint) {
	loc.SpellingLocation(&file, &line, &column, &offset)
	return
}

// Traverse up the semantic parents
func BuildScopingParts(cursor clang.Cursor) []string {
	var parts []string
	for cursor.IsNull() != 1 && cursor.Kind != clang.CursorTranslationUnit {
		name := cursor.String()
		qualified := c.GoString(name.CStr())
		parts = append([]string{qualified}, parts...)
		cursor = cursor.SemanticParent()
		name.Dispose()
	}
	return parts
}

func VisitChildren(cursor clang.Cursor, fn Visitor) c.Uint {
	return clang.VisitChildren(cursor, func(cursor, parent clang.Cursor, clientData unsafe.Pointer) clang.ChildVisitResult {
		cfn := *(*Visitor)(clientData)
		return cfn(cursor, parent)
	}, unsafe.Pointer(&fn))
}

func GetInclusions(unit *clang.TranslationUnit, visitor InclusionVisitor) {
	clang.GetInclusions(unit, func(inced clang.File, incin *clang.SourceLocation, incilen c.Uint, data c.Pointer) {
		ics := unsafe.Slice(incin, incilen)
		cfn := *(*InclusionVisitor)(data)
		cfn(inced, ics)
	}, unsafe.Pointer(&visitor))
}

func defaultArgs(isCpp bool) []string {
	args := []string{"-x", "c"}
	if isCpp {
		args = []string{"-x", "c++"}
	}
	return args
}

// CINDEX_LINKAGE CXTranslationUnit clang_parseTranslationUnit(CXIndex CIdx, const char *source_filename,
// 	const char *const *command_line_args,
// 	int num_command_line_args,
// 	struct CXUnsavedFile *unsaved_files,
// 	unsigned num_unsaved_files, unsigned options);

/**
 * Same as \c clang_parseTranslationUnit2, but returns
 * the \c CXTranslationUnit instead of an error code.  In case of an error this
 * routine returns a \c NULL \c CXTranslationUnit, without further detailed
 * error codes.
 */
//go:linkname ParseTranslationUnit C.clang_parseTranslationUnit
func ParseTranslationUnit(index *clang.Index, sourceFilename *c.Char, commandLineArgs **c.Char, numCommandLineArgs c.Int,
	unsavedFiles *clang.UnsavedFile, numUnsavedFiles c.Uint, options c.Uint) *clang.TranslationUnit

/**
 * Same as clang_parseTranslationUnit2 but requires a full command line
 * for \c command_line_args including argv[0]. This is useful if the standard
 * library paths are relative to the binary.
 */
// CINDEX_LINKAGE enum CXErrorCode
// clang_parseTranslationUnit2FullArgv(CXIndex CIdx, const char *source_filename, const char *const *command_line_args,
//                                     int num_command_line_args, struct CXUnsavedFile *unsaved_files,
//                                     unsigned num_unsaved_files, unsigned options, CXTranslationUnit *out_TU);

type ErrorCode int

const (
	Error_Success          = 0
	Error_Failure          = 1
	Error_Crashed          = 2
	Error_InvalidArguments = 3
	Error_ASTReadError     = 4
)

//go:linkname ParseTranslationUnit2FullArgv C.clang_parseTranslationUnit2FullArgv
func ParseTranslationUnit2FullArgv(index *clang.Index, sourceFilename *c.Char, commandLineArgs **c.Char, numCommandLineArgs c.Int,
	unsavedFiles *clang.UnsavedFile, numUnsavedFiles c.Uint, options c.Uint, out_TU **clang.TranslationUnit) ErrorCode
