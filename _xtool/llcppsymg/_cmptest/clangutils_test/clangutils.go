package main

import (
	"fmt"
	"os"

	"github.com/goplus/llcppg/_xtool/llcppsymg/clangutils"
	"github.com/goplus/llgo/c"
	"github.com/goplus/llgo/c/clang"
)

func main() {
	TestClangUtil()
	TestComposeIncludes()
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
		index, unit, err := clangutils.CreateTranslationUnit(config)
		if err != nil {
			fmt.Printf("CreateTranslationUnit failed: %v\n", err)
			continue
		}

		fmt.Println("CreateTranslationUnit succeeded")

		cursor := unit.Cursor()

		clangutils.VisitChildren(cursor, func(cursor, parent clang.Cursor) clang.ChildVisitResult {
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
			case clang.CursorNamespace:
				namespaceName := cursor.String()
				fmt.Printf("Namespace: %s\n", c.GoString(namespaceName.CStr()))
				namespaceName.Dispose()
			}
			return clang.ChildVisit_Recurse
		})
		index.Dispose()
		unit.Dispose()

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
