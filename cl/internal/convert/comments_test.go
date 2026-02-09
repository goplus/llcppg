package convert_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/goplus/gogen"
	llcppgast "github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/cl/internal/convert"
)

// TestCommentParsing tests the comment parsing and conversion logic.
// These atomic tests verify that comments are correctly converted from
// llcppg AST to Go AST following Go's ast.Comment specification:
// - For /* */ block comments: entire block must be a single Comment node
// - For // line comments: each line is a separate Comment node

func TestBlockCommentSingleNode(t *testing.T) {
	// A multi-line block comment should be a single ast.Comment node
	blockComment := "/* Create an iterator for traversing a domain\n   The domain NULL denotes the default domain */"

	// Create llcppg AST with single Comment node (correct approach)
	llcppgDoc := &llcppgast.CommentGroup{
		List: []*llcppgast.Comment{
			{Text: blockComment},
		},
	}

	// Convert to Go AST
	goDoc := convert.NewCommentGroupFromC(llcppgDoc)

	// Verify single Comment node
	if len(goDoc.List) != 1 {
		t.Errorf("Expected 1 Comment node, got %d", len(goDoc.List))
	}
	if goDoc.List[0].Text != blockComment {
		t.Errorf("Comment text mismatch.\nExpected: %q\nGot: %q", blockComment, goDoc.List[0].Text)
	}

	// Verify the generated code is valid Go
	assertValidGoCode(t, goDoc, "BlockCommentSingleNode")
}

func TestBlockCommentSplitByNewlines_Invalid(t *testing.T) {
	// This test demonstrates the WRONG approach that causes the bug
	// Splitting a block comment by newlines creates invalid Go code with gogen v1.20.2
	blockComment := "/* Create an iterator for traversing a domain\n   The domain NULL denotes the default domain */"

	// WRONG: Split block comment by newlines (old behavior)
	lines := strings.Split(blockComment, "\n")
	var llcppgComments []*llcppgast.Comment
	for _, line := range lines {
		llcppgComments = append(llcppgComments, &llcppgast.Comment{Text: line + "\n"})
	}
	llcppgDoc := &llcppgast.CommentGroup{List: llcppgComments}

	// Convert to Go AST
	goDoc := convert.NewCommentGroupFromC(llcppgDoc)

	// This produces multiple Comment nodes (wrong!)
	if len(goDoc.List) == 1 {
		t.Skip("If this passes with 1 node, the test setup may be wrong")
	}

	// With gogen v1.20.2, this should produce invalid Go code
	// Skip validation as we expect this to be invalid
	t.Logf("Created %d Comment nodes (wrong approach, may cause invalid Go code)", len(goDoc.List))
}

func TestLineCommentsSplitByNewlines(t *testing.T) {
	// Line comments should be split by newlines, one Comment per line
	// This is the correct behavior for // style comments

	// Create llcppg AST with separate Comment nodes for each line
	llcppgDoc := &llcppgast.CommentGroup{
		List: []*llcppgast.Comment{
			{Text: "// Line 1"},
			{Text: "// Line 2"},
			{Text: "// Line 3"},
		},
	}

	// Convert to Go AST
	goDoc := convert.NewCommentGroupFromC(llcppgDoc)

	// Verify three Comment nodes
	if len(goDoc.List) != 3 {
		t.Errorf("Expected 3 Comment nodes, got %d", len(goDoc.List))
	}

	// Verify the generated code is valid Go
	assertValidGoCode(t, goDoc, "LineCommentsSplit")
}

func TestSingleLineBlockComment(t *testing.T) {
	// A single-line block comment should also be a single Comment node
	blockComment := "/* Single line block comment */"

	llcppgDoc := &llcppgast.CommentGroup{
		List: []*llcppgast.Comment{
			{Text: blockComment},
		},
	}

	goDoc := convert.NewCommentGroupFromC(llcppgDoc)

	if len(goDoc.List) != 1 {
		t.Errorf("Expected 1 Comment node, got %d", len(goDoc.List))
	}
	if goDoc.List[0].Text != blockComment {
		t.Errorf("Comment text mismatch.\nExpected: %q\nGot: %q", blockComment, goDoc.List[0].Text)
	}

	assertValidGoCode(t, goDoc, "SingleLineBlockComment")
}

func TestEmptyCommentGroup(t *testing.T) {
	// Test nil and empty comment groups
	t.Run("nil", func(t *testing.T) {
		goDoc := convert.NewCommentGroupFromC(nil)
		if goDoc == nil {
			t.Error("Expected non-nil CommentGroup")
		}
		if len(goDoc.List) != 0 {
			t.Errorf("Expected 0 Comment nodes, got %d", len(goDoc.List))
		}
	})

	t.Run("empty_list", func(t *testing.T) {
		llcppgDoc := &llcppgast.CommentGroup{List: nil}
		goDoc := convert.NewCommentGroupFromC(llcppgDoc)
		if goDoc == nil {
			t.Error("Expected non-nil CommentGroup")
		}
		if len(goDoc.List) != 0 {
			t.Errorf("Expected 0 Comment nodes, got %d", len(goDoc.List))
		}
	})
}

func TestMixedComments(t *testing.T) {
	// Test a group with both block and line comments
	// In practice, this would be separate groups, but test the conversion anyway
	llcppgDoc := &llcppgast.CommentGroup{
		List: []*llcppgast.Comment{
			{Text: "/* Block comment */"},
			{Text: "// Line comment"},
		},
	}

	goDoc := convert.NewCommentGroupFromC(llcppgDoc)

	if len(goDoc.List) != 2 {
		t.Errorf("Expected 2 Comment nodes, got %d", len(goDoc.List))
	}

	assertValidGoCode(t, goDoc, "MixedComments")
}

func TestBlockCommentWithSpecialChars(t *testing.T) {
	// Test block comments with special characters
	testCases := []struct {
		name    string
		comment string
	}{
		{"asterisks", "/* Comment with * asterisks * inside */"},
		{"slashes", "/* Comment with / slashes / inside */"},
		{"newlines_and_tabs", "/* Comment with\n\ttabs and\n\tnewlines */"},
		{"unicode", "/* Unicode: 中文, 日本語, émojis 🎉 */"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			llcppgDoc := &llcppgast.CommentGroup{
				List: []*llcppgast.Comment{
					{Text: tc.comment},
				},
			}

			goDoc := convert.NewCommentGroupFromC(llcppgDoc)

			if len(goDoc.List) != 1 {
				t.Errorf("Expected 1 Comment node, got %d", len(goDoc.List))
			}
			if goDoc.List[0].Text != tc.comment {
				t.Errorf("Comment text mismatch.\nExpected: %q\nGot: %q", tc.comment, goDoc.List[0].Text)
			}
		})
	}
}

// assertValidGoCode generates Go code with the comment and validates it
func assertValidGoCode(t *testing.T, commentGroup *ast.CommentGroup, funcName string) {
	t.Helper()

	pkg := gogen.NewPackage("", "demo", nil)
	fn := pkg.NewFunc(nil, funcName, nil, nil, false)
	fn.SetComments(pkg, commentGroup)
	fn.BodyStart(pkg).End()

	var buf strings.Builder
	err := gogen.WriteTo(&buf, pkg, "")
	if err != nil {
		t.Fatalf("gogen.WriteTo failed: %v", err)
	}

	code := buf.String()
	t.Logf("Generated code:\n%s", code)

	fset := token.NewFileSet()
	_, err = parser.ParseFile(fset, "generated.go", code, parser.ParseComments)
	if err != nil {
		t.Fatalf("Generated code is invalid Go: %v\nCode:\n%s", err, code)
	}
}
