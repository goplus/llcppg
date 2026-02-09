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
// with gogen's printer. This test demonstrates the current comment parsing approach
// used in llcppg (_xtool/internal/parser/parser.go:ParseComment).
//
// Background:
//   - Go's ast.Comment documentation states each Comment node represents
//     "a single //-style or /*-style comment"
//   - For /* */ comments, the entire block should be a single Comment node
//   - llcppg currently splits all comments by newlines, creating multiple Comment nodes
//
// This test should PASS with gogen v1.19.7 but FAIL with gogen v1.20.2
// due to stricter comment handling in the updated printer.
func TestMultiLineBlockComment(t *testing.T) {
	// This is a typical multi-line C block comment from a header file
	// Example from gettext: /* Create an iterator for traversing a domain
	//                          The domain NULL denotes the default domain */
	rawComment := "/* Create an iterator for traversing a domain\n   The domain NULL denotes the default domain */"

	t.Run("current_approach_split_by_newlines", func(t *testing.T) {
		// Current approach in llcppg: split by newlines (as in ParseComment)
		// This mimics _xtool/internal/parser/parser.go:195-202
		pkg := gogen.NewPackage("", "demo", nil)

		lines := strings.Split(rawComment, "\n")
		var commentList []*ast.Comment
		for _, line := range lines {
			commentList = append(commentList, &ast.Comment{Text: line + "\n"})
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

	t.Run("correct_approach_single_comment_node", func(t *testing.T) {
		// Correct approach: keep block comment as single Comment node
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
}
