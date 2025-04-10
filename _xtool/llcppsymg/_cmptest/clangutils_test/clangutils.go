package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/goplus/lib/c"
	"github.com/goplus/lib/c/clang"
	"github.com/goplus/llcppg/_xtool/llcppsymg/clangutils"
)

func main() {
	TestClangUtil()
	TestComposeIncludes()
	TestPreprocess()
	TestComment()
}

func TestClangUtil() {
	testCases := []struct {
		name    string
		content string
		isTemp  bool
		isCpp   bool
	}{
		{
			name: "C Header File",
			content: `
				int test_function(int a, int b);
				void another_function(void);
			`,
			isTemp: false,
			isCpp:  false,
		},
		{
			name: "C++ Temp File",
			content: `
				class TestClass {
				public:
					void test_method();
					static int static_method(float f);
				};
				
				namespace TestNamespace {
					void namespaced_function();
				}
			`,
			isTemp: true,
			isCpp:  true,
		},
	}

	for _, tc := range testCases {
		fmt.Printf("=== Test Case: %s ===\n", tc.name)

		var filePath string
		var tempFile *os.File
		if tc.isTemp {
			filePath = tc.content
		} else {
			var err error
			tempFile, err = os.CreateTemp("", "test_*.h")
			if err != nil {
				fmt.Printf("Failed to create temporary file: %v\n", err)
				continue
			}

			_, err = tempFile.Write([]byte(tc.content))
			if err != nil {
				fmt.Printf("Failed to write to temporary file: %v\n", err)
				tempFile.Close()
				os.Remove(tempFile.Name())
				continue
			}
			tempFile.Close()
			filePath = tempFile.Name()
		}

		config := &clangutils.Config{
			File:  filePath,
			Temp:  tc.isTemp,
			IsCpp: tc.isCpp,
		}
		outputInfoFromTranslationUnit(config, func(cursor, parent clang.Cursor) clang.ChildVisitResult {
			switch cursor.Kind {
			case clang.CursorFunctionDecl, clang.CursorCXXMethod:
				funcName := cursor.String()
				fmt.Printf("Function/Method: %s\n", c.GoString(funcName.CStr()))
				parts := clangutils.BuildScopingParts(cursor)
				fmt.Printf("Scoping parts: %v\n", parts)
				funcName.Dispose()
			case clang.CursorClassDecl:
				className := cursor.String()
				fmt.Printf("Class: %s\n", c.GoString(className.CStr()))
				className.Dispose()
				return clang.ChildVisit_Recurse
			case clang.CursorNamespace:
				namespaceName := cursor.String()
				fmt.Printf("Namespace: %s\n", c.GoString(namespaceName.CStr()))
				namespaceName.Dispose()
				return clang.ChildVisit_Recurse
			}
			return clang.ChildVisit_Continue
		})

		if !tc.isTemp && tempFile != nil {
			os.Remove(tempFile.Name())
		}

		fmt.Println()
	}
}

func TestComposeIncludes() {
	fmt.Println("=== Test ComposeIncludes ===")
	testCases := []struct {
		name  string
		files []string
	}{
		{
			name:  "One file",
			files: []string{"file1.h"},
		},
		{
			name:  "Two files",
			files: []string{"file1.h", "file2.h"},
		},
		{
			name:  "Empty files",
			files: []string{},
		},
	}
	for _, tc := range testCases {
		outfile, err := os.CreateTemp("", "compose_*.h")
		if err != nil {
			panic(err)
		}
		err = clangutils.ComposeIncludes(tc.files, outfile.Name())
		if err != nil {
			panic(err)
		}
		content, err := os.ReadFile(outfile.Name())
		if err != nil {
			panic(err)
		}
		fmt.Println(string(content))
		outfile.Close()
		os.Remove(outfile.Name())
	}
}

func TestPreprocess() {
	fmt.Println("=== TestPreprocess ===")
	outfile, err := os.CreateTemp("", "compose_*.h")
	if err != nil {
		panic(err)
	}
	absPath, err := filepath.Abs("./hfile")
	if err != nil {
		panic(err)
	}
	clangutils.ComposeIncludes([]string{"main.h", "compat.h"}, outfile.Name())

	efile, err := os.CreateTemp("", "temp_*.i")
	if err != nil {
		panic(err)
	}
	defer os.Remove(efile.Name())

	cfg := &clangutils.PreprocessConfig{
		File:    outfile.Name(),
		IsCpp:   true,
		Args:    []string{"-I" + absPath},
		OutFile: efile.Name(),
	}
	err = clangutils.Preprocess(cfg)
	if err != nil {
		panic(err)
	}
	config := &clangutils.Config{
		File:  efile.Name(),
		Temp:  false,
		IsCpp: false,
	}
	outputInfoFromTranslationUnit(config, func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		switch cursor.Kind {
		case clang.CursorEnumDecl, clang.CursorStructDecl, clang.CursorUnionDecl, clang.CursorTypedefDecl:
			declName := cursor.String()
			var filename clang.String
			var line, column c.Uint
			cursor.Location().PresumedLocation(&filename, &line, &column)
			fmt.Printf("TypeKind: %d Name: %s\n", cursor.Kind, c.GoString(declName.CStr()))
			fmt.Printf("Location: %s:%d:%d\n", path.Base(c.GoString(filename.CStr())), line, column)
			declName.Dispose()
		}
		return clang.ChildVisit_Continue
	})
	outfile.Close()
	os.Remove(outfile.Name())
}

func TestComment() {
	fmt.Println("=== TestComment ===")
	config := &clangutils.Config{
		File: `
		#include <comment.h>
		`,
		Temp:  true,
		IsCpp: false,
		Args:  []string{"-I./hfile", "-E", "-fparse-all-comments"},
	}
	index, unit, err := clangutils.CreateTranslationUnit(config)
	if err != nil {
		fmt.Printf("CreateTranslationUnit failed: %v\n", err)
	}

	cursor := unit.Cursor()

	clangutils.VisitChildren(cursor, func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		if cursor.Kind != clang.CursorMacroDefinition && cursor.Kind != clang.CursorInclusionDirective {
			fmt.Println("cursor", clang.GoString(cursor.String()), "rawComment:", clang.GoString(cursor.RawCommentText()))
			commentRange := cursor.CommentRange()
			cursorRange := cursor.Extent()
			fmt.Printf("commentRange %d:%d -> %d:%d\n", commentRange.RangeStart().Line(), commentRange.RangeStart().Column(), commentRange.RangeEnd().Line(), commentRange.RangeEnd().Column())
			fmt.Printf("cursorRange %d:%d -> %d:%d\n", cursorRange.RangeStart().Line(), cursorRange.RangeStart().Column(), cursorRange.RangeEnd().Line(), cursorRange.RangeEnd().Column())
			fmt.Println("--------------------------------")
		}
		return clang.ChildVisit_Recurse
	})

	index.Dispose()
	unit.Dispose()
}

func outputInfoFromTranslationUnit(config *clangutils.Config, visitFunc func(cursor, parent clang.Cursor) clang.ChildVisitResult) {
	index, unit, err := clangutils.CreateTranslationUnit(config)
	if err != nil {
		panic(err)
	}

	fmt.Println("CreateTranslationUnit succeeded")

	cursor := unit.Cursor()

	clangutils.VisitChildren(cursor, visitFunc)
	index.Dispose()
	unit.Dispose()
}
