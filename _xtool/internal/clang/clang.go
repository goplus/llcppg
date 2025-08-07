package clang

import (
	"errors"
	"path/filepath"
	"unsafe"

	"github.com/goplus/lib/c"
	"github.com/goplus/llcppg/_xtool/internal/clangtool"
	clang "github.com/goplus/llcppg/_xtool/internal/libclang"
)

type Config struct {
	File    string
	Temp    bool
	Args    []string
	IsCpp   bool
	Index   *clang.Index
	Options c.Uint
}

type Visitor func(cursor, parent clang.Cursor) clang.ChildVisitResult

type InclusionVisitor func(included_file clang.File, inclusions []clang.SourceLocation)

const TEMP_FILE = "temp.h"

func CreateTranslationUnit(config *Config) (*clang.Index, *clang.TranslationUnit, error) {
	// default use the C/C++ standard of clang; c:gnu17 C++:gnu++17
	// https://clang.llvm.org/docs/CommandGuide/clang.html
	allArgs := clangtool.WithSysRoot(append(defaultArgs(config.IsCpp), config.Args...))

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

	if config.Temp {
		content := c.AllocaCStr(config.File)
		tempFile := &clang.UnsavedFile{
			Filename: c.Str(TEMP_FILE),
			Contents: content,
			Length:   c.Ulong(c.Strlen(content)),
		}

		unit = index.ParseTranslationUnit(
			tempFile.Filename,
			unsafe.SliceData(cArgs), c.Int(len(cArgs)),
			tempFile, 1,
			config.Options,
		)

	} else {
		cFile := c.AllocaCStr(config.File)
		unit = index.ParseTranslationUnit(
			cFile,
			unsafe.SliceData(cArgs), c.Int(len(cArgs)),
			nil, 0,
			config.Options,
		)
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

func GetPresumedLocation(loc clang.SourceLocation) (fileGo string, line c.Uint, column c.Uint) {
	var file clang.String
	loc.PresumedLocation(&file, &line, &column)
	fileGo = filepath.Clean(clang.GoString(file))
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
