package clang_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/goplus/lib/c/clang"
	clangutils "github.com/goplus/llcppg/_xtool/internal/clang"
)

func TestClangUtil(t *testing.T) {
	testCases := []struct {
		name    string
		content string
		isTemp  bool
		isCpp   bool
		expect  string
	}{
		{
			name: "C Header File",
			content: `
				int test_function(int a, int b);
				void another_function(void);
			`,
			isTemp: false,
			isCpp:  false,
			expect: `
Function/Method: test_function
Scoping parts: test_function
Function/Method: another_function
Scoping parts: another_function
`,
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
			expect: `
Class: TestClass
Function/Method: test_method
Scoping parts: TestClass,test_method
Function/Method: static_method
Scoping parts: TestClass,static_method
Namespace: TestNamespace
Function/Method: namespaced_function
Scoping parts: TestNamespace,namespaced_function
			`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var filePath string
			var tempFile *os.File
			if tc.isTemp {
				filePath = tc.content
			} else {
				var err error
				tempFile, err = os.CreateTemp("", "test_*.h")
				if err != nil {
					t.Fatalf("Failed to create temporary file: %w\n", err)
					return
				}
				defer tempFile.Close()

				_, err = tempFile.Write([]byte(tc.content))
				if err != nil {
					t.Fatalf("Failed to write to temporary file: %w\n", err)
				}
				defer os.Remove(tempFile.Name())
				filePath = tempFile.Name()
			}

			config := &clangutils.Config{
				File:  filePath,
				Temp:  tc.isTemp,
				IsCpp: tc.isCpp,
			}

			var str strings.Builder

			visit(config, func(cursor, parent clang.Cursor) clang.ChildVisitResult {
				switch cursor.Kind {
				case clang.CursorFunctionDecl, clang.CursorCXXMethod:
					str.WriteString("Function/Method: ")
					str.WriteString(clang.GoString(cursor.String()))
					str.WriteString("\n")
					parts := clangutils.BuildScopingParts(cursor)
					str.WriteString("Scoping parts: ")
					str.WriteString(strings.Join(parts, ","))
					str.WriteString("\n")
				case clang.CursorClassDecl:
					str.WriteString("Class: ")
					str.WriteString(clang.GoString(cursor.String()))
					str.WriteString("\n")
					return clang.ChildVisit_Recurse
				case clang.CursorNamespace:
					str.WriteString("Namespace: ")
					str.WriteString(clang.GoString(cursor.String()))
					str.WriteString("\n")
					return clang.ChildVisit_Recurse
				}
				return clang.ChildVisit_Continue
			})
			compareOutput(t, tc.expect, str.String())
		})
	}
}

func TestComment(t *testing.T) {
	config := &clangutils.Config{
		File: `
		#include <comment.h>
		`,
		Temp:  true,
		IsCpp: false,
		Args:  []string{"-I./testdata/hfile", "-E", "-fparse-all-comments"},
	}

	var str strings.Builder

	visit(config, func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		if cursor.Kind != clang.CursorMacroDefinition && cursor.Kind != clang.CursorInclusionDirective {
			str.WriteString("cursor ")
			str.WriteString(clang.GoString(cursor.String()))
			str.WriteString("rawComment: ")
			str.WriteString(clang.GoString(cursor.RawCommentText()))
			str.WriteString("\n")
			commentRange := cursor.CommentRange()
			cursorRange := cursor.Extent()
			str.WriteString("commentRange ")
			str.WriteString(fmt.Sprintf("%d:%d -> %d:%d\n", commentRange.RangeStart().Line(), commentRange.RangeStart().Column(), commentRange.RangeEnd().Line(), commentRange.RangeEnd().Column()))
			str.WriteString("cursorRange ")
			str.WriteString(fmt.Sprintf("%d:%d -> %d:%d\n", cursorRange.RangeStart().Line(), cursorRange.RangeStart().Column(), cursorRange.RangeEnd().Line(), cursorRange.RangeEnd().Column()))
			str.WriteString("--------------------------------\n")
		}
		return clang.ChildVisit_Recurse
	})

	expect := `
cursor FoorawComment: // doc
commentRange 1:1 -> 1:7
cursorRange 2:1 -> 8:2
--------------------------------
cursor xrawComment: // doc
commentRange 3:5 -> 3:11
cursorRange 4:5 -> 4:10
--------------------------------
cursor yrawComment: // comment
commentRange 5:12 -> 5:22
cursorRange 5:5 -> 5:10
--------------------------------
cursor zrawComment: // comment
commentRange 7:12 -> 7:22
cursorRange 7:5 -> 7:10
--------------------------------
cursor foorawComment: // doc
commentRange 10:1 -> 10:7
cursorRange 11:1 -> 11:11
--------------------------------`

	compareOutput(t, expect, str.String())
}

func visit(config *clangutils.Config, visitFunc func(cursor, parent clang.Cursor) clang.ChildVisitResult) {
	index, unit, err := clangutils.CreateTranslationUnit(config)
	if err != nil {
		panic(err)
	}
	cursor := unit.Cursor()
	clangutils.VisitChildren(cursor, visitFunc)
	index.Dispose()
	unit.Dispose()
}

func compareOutput(t *testing.T, expected, actual string) {
	expected = strings.TrimSpace(expected)
	actual = strings.TrimSpace(actual)
	if expected != actual {
		t.Fatalf("Test failed: expected \n%s \ngot \n%s", expected, actual)
	}
}
