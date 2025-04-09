package clangutils

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unsafe"

	"github.com/goplus/lib/c"
	"github.com/goplus/lib/c/clang"
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
	allArgs := append(defaultArgs(config.IsCpp), config.Args...)

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
			clang.DetailedPreprocessingRecord,
		)

	} else {
		cFile := c.AllocaCStr(config.File)
		unit = index.ParseTranslationUnit(
			cFile,
			unsafe.SliceData(cArgs), c.Int(len(cArgs)),
			nil, 0,
			clang.DetailedPreprocessingRecord,
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

// ComposeIncludes create Include list
// #include <file1.h>
// #include <file2.h>
func ComposeIncludes(files []string, outfile string) error {
	var str string
	for _, file := range files {
		str += ("#include <" + file + ">\n")
	}
	return os.WriteFile(outfile, []byte(str), 0644)
}

func defaultArgs(isCpp bool) []string {
	args := []string{"-x", "c"}
	if isCpp {
		args = []string{"-x", "c++"}
	}
	return args
}

type PreprocessConfig struct {
	File    string
	IsCpp   bool
	Args    []string
	OutFile string
}

func Preprocess(cfg *PreprocessConfig) error {
	args := []string{"-E"}
	args = append(args, defaultArgs(cfg.IsCpp)...)
	args = append(args, cfg.Args...)
	args = append(args, cfg.File)
	args = append(args, "-o", cfg.OutFile)
	cmd := exec.Command("clang", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func GetIncludePaths(isCpp bool) []string {
	args := []string{"-E", "-v"}
	args = append(args, defaultArgs(isCpp)...)
	args = append(args, "/dev/null")
	cmd := exec.Command("clang", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	return ParseClangIncOutput(string(output))
}

func ParseClangIncOutput(output string) []string {
	var paths []string
	start := strings.Index(output, "#include <...> search starts here:")
	end := strings.Index(output, "End of search list.")
	if start == -1 || end == -1 {
		return paths
	}
	content := output[start:end]
	lines := strings.Split(content, "\n")
	for _, line := range lines[1:] {
		for _, item := range strings.Fields(line) {
			if path := strings.TrimSpace(item); filepath.IsAbs(path) {
				paths = append(paths, path)
			}
		}
	}
	return paths
}
