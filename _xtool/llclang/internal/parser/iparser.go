package parser

import (
	"fmt"
	"strconv"
	"unsafe"

	"github.com/goplus/lib/c"
	cparser "github.com/goplus/llcppg/_xtool/internal/parser"
	"github.com/goplus/llcppg/parser"
	"github.com/goplus/llpkg/cjson"
)

type Mode = parser.Mode

func RunParseIntermediateFile(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("RunParseIntermediateFile: lack of argument")
	}
	filename := args[0]
	var mode Mode
	if modeStr := args[1]; modeStr != "" {
		modeInt, err := strconv.Atoi(modeStr)
		if err != nil {
			return fmt.Errorf("RunParseIntermediateFile: invalid mode %s: %w", modeStr, err)
		}
		mode = Mode(modeInt)
	}
	err := parseIntermediateFile(filename, mode)
	if err != nil {
		return err
	}
	return nil
}

// parses an intermediate file (*.i) and output the corresponding AST to stdout.
func parseIntermediateFile(filename string, mode Mode) error {
	isCpp := mode&parser.ParseC == 0

	args := []string{}
	if mode&parser.ParseAllComments != 0 {
		args = append(args, "-fparse-all-comments")
	}

	file, err := cparser.Do(&cparser.ConverterConfig{
		File:  filename,
		IsCpp: isCpp,
		Args:  args,
	})
	if err != nil {
		return fmt.Errorf("parseIntermediateFile: %w", err)
	}
	json := cparser.MarshalASTFile(file)
	str := json.Print()
	defer cjson.FreeCStr(unsafe.Pointer(str))
	defer json.Delete()
	c.Printf(c.Str("%s"), str)
	return nil
}
