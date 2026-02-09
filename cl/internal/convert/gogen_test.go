package convert_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/goplus/gogen"
)

// TestMultiLineBlockComment tests the behavior of multi-line block comments
// with gogen's printer. This test verifies the FIXED comment parsing approach
// used in llcppg (_xtool/internal/parser/parser.go:ParseComment).
//
// Background:
//   - Go's ast.Comment documentation states each Comment node represents
//     "a single //-style or /*-style comment"
//   - For /* */ comments, the entire block should be a single Comment node
//   - The fix keeps block comments as single nodes instead of splitting by newlines
//
// This test should PASS with both gogen v1.19.7 and v1.20.2.
func TestMultiLineBlockComment(t *testing.T) {
	// This is a typical multi-line C block comment from a header file
	// Example from gettext: /* Create an iterator for traversing a domain
	//                          The domain NULL denotes the default domain */
	rawComment := "/* Create an iterator for traversing a domain\n   The domain NULL denotes the default domain */"

	t.Run("fixed_approach_single_comment_node", func(t *testing.T) {
		// Fixed approach: keep block comment as single Comment node
		// This is how ParseComment now works after the fix
		pkg := gogen.NewPackage("", "demo", nil)

		commentList := []*ast.Comment{
			{
				Slash: token.NoPos,
				Text:  rawComment,
			},
		}
		commentGroup := &ast.CommentGroup{List: commentList}

		fn := pkg.NewFunc(nil, "ExampleFunction", nil, nil, false)
		fn.SetComments(pkg, commentGroup)
		fn.BodyStart(pkg).End()

		var buf strings.Builder
		err := gogen.WriteTo(&buf, pkg, "")
		if err != nil {
			t.Fatalf("gogen.WriteTo failed: %v", err)
		}

		code := buf.String()
		t.Logf("Generated code:\n%s", code)

		// Validate the generated code is valid Go
		fset := token.NewFileSet()
		_, err = parser.ParseFile(fset, "generated.go", code, parser.ParseComments)
		if err != nil {
			t.Fatalf("Generated code is invalid Go: %v\nCode:\n%s", err, code)
		}
	})

	t.Run("line_comments_split_by_newlines", func(t *testing.T) {
		// Line comments should be split by newlines (one Comment per line)
		// This is the correct behavior for // style comments
		pkg := gogen.NewPackage("", "demo", nil)

		commentList := []*ast.Comment{
			{Slash: token.NoPos, Text: "// Line 1 of comment"},
			{Slash: token.NoPos, Text: "// Line 2 of comment"},
		}
		commentGroup := &ast.CommentGroup{List: commentList}

		fn := pkg.NewFunc(nil, "ExampleFunction", nil, nil, false)
		fn.SetComments(pkg, commentGroup)
		fn.BodyStart(pkg).End()

		var buf strings.Builder
		err := gogen.WriteTo(&buf, pkg, "")
		if err != nil {
			t.Fatalf("gogen.WriteTo failed: %v", err)
		}

		code := buf.String()
		t.Logf("Generated code:\n%s", code)

		// Validate the generated code is valid Go
		fset := token.NewFileSet()
		_, err = parser.ParseFile(fset, "generated.go", code, parser.ParseComments)
		if err != nil {
			t.Fatalf("Generated code is invalid Go: %v\nCode:\n%s", err, code)
		}
	})
}

// TestParseCommentLogic tests the logic that should be implemented in ParseComment
// This simulates the fixed behavior of _xtool/internal/parser/parser.go:ParseComment
func TestParseCommentLogic(t *testing.T) {
	// parseComment simulates the fixed ParseComment function
	parseComment := func(rawComment string) []*ast.Comment {
		// Block comment (/* ... */) - keep as single Comment node
		if strings.HasPrefix(rawComment, "/*") {
			return []*ast.Comment{{Text: rawComment}}
		}

		// Line comments (// ...) - split by newlines
		var comments []*ast.Comment
		lines := strings.Split(rawComment, "\n")
		for _, line := range lines {
			line = strings.TrimRight(line, "\r")
			if line != "" {
				comments = append(comments, &ast.Comment{Text: line})
			}
		}
		return comments
	}

	testCases := []struct {
		name          string
		input         string
		expectedCount int
	}{
		{
			name:          "single_line_block_comment",
			input:         "/* Single line block */",
			expectedCount: 1,
		},
		{
			name:          "multi_line_block_comment",
			input:         "/* Line 1\n   Line 2\n   Line 3 */",
			expectedCount: 1,
		},
		{
			name:          "single_line_comment",
			input:         "// Single line",
			expectedCount: 1,
		},
		{
			name:          "multi_line_line_comments",
			input:         "// Line 1\n// Line 2\n// Line 3",
			expectedCount: 3,
		},
		{
			name:          "block_comment_with_asterisks",
			input:         "/* Comment with * asterisks * inside */",
			expectedCount: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			comments := parseComment(tc.input)
			if len(comments) != tc.expectedCount {
				t.Errorf("Expected %d comments, got %d", tc.expectedCount, len(comments))
			}

			// For block comments, verify the text is preserved
			if strings.HasPrefix(tc.input, "/*") {
				if len(comments) > 0 && comments[0].Text != tc.input {
					t.Errorf("Block comment text mismatch.\nExpected: %q\nGot: %q", tc.input, comments[0].Text)
				}
			}
		})
	}
}
