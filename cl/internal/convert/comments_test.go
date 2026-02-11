package convert_test

import (
	goast "go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/goplus/gogen"
	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/cl/internal/convert"
)

func TestNewCommentGroupFromC(t *testing.T) {
	tests := []struct {
		name     string
		doc      *ast.CommentGroup
		expected []string
	}{
		{
			name:     "nil_input",
			doc:      nil,
			expected: nil,
		},
		{
			name:     "empty_list",
			doc:      &ast.CommentGroup{List: nil},
			expected: nil,
		},
		{
			name: "single_block_comment",
			doc: &ast.CommentGroup{
				List: []*ast.Comment{
					{Text: "/* This is a block comment */"},
				},
			},
			expected: []string{"/* This is a block comment */"},
		},
		{
			name: "multi_line_block_comment",
			doc: &ast.CommentGroup{
				List: []*ast.Comment{
					{Text: "/* Create an iterator for traversing a domain\n   The domain NULL denotes the default domain */"},
				},
			},
			expected: []string{"/* Create an iterator for traversing a domain\n   The domain NULL denotes the default domain */"},
		},
		{
			name: "legacy_split_block_two_nodes",
			doc: &ast.CommentGroup{
				List: []*ast.Comment{
					{Text: "/* Create an iterator for traversing a domain\n"},
					{Text: "   The domain NULL denotes the default domain.  */\n"},
				},
			},
			expected: []string{"/* Create an iterator for traversing a domain\n   The domain NULL denotes the default domain.  */"},
		},
		{
			name: "legacy_split_block_three_nodes",
			doc: &ast.CommentGroup{
				List: []*ast.Comment{
					{Text: "/* Return the error string for ERR in the user-supplied buffer BUF of\n"},
					{Text: " * size BUFLEN.  This function is, in contrast to gpg_strerror,\n"},
					{Text: " * thread-safe.  */\n"},
				},
			},
			expected: []string{"/* Return the error string for ERR in the user-supplied buffer BUF of\n * size BUFLEN.  This function is, in contrast to gpg_strerror,\n * thread-safe.  */"},
		},
		{
			name: "legacy_javadoc_style",
			doc: &ast.CommentGroup{
				List: []*ast.Comment{
					{Text: "/**\n"},
					{Text: "Foo comment\n"},
					{Text: "*/\n"},
				},
			},
			expected: []string{"/**\nFoo comment\n*/"},
		},
		{
			name: "single_line_comment",
			doc: &ast.CommentGroup{
				List: []*ast.Comment{
					{Text: "// Foo comment"},
				},
			},
			expected: []string{"// Foo comment"},
		},
		{
			name: "multiple_line_comments",
			doc: &ast.CommentGroup{
				List: []*ast.Comment{
					{Text: "// First line\n"},
					{Text: "// Second line\n"},
				},
			},
			expected: []string{"// First line", "// Second line"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goDoc := convert.NewCommentGroupFromC(tt.doc)
			if len(goDoc.List) != len(tt.expected) {
				t.Fatalf("expected %d comments, got %d", len(tt.expected), len(goDoc.List))
			}
			for i, want := range tt.expected {
				if goDoc.List[i].Text != want {
					t.Errorf("comment[%d]: got %q, want %q", i, goDoc.List[i].Text, want)
				}
			}
		})
	}
}

// TestMultiLineBlockCommentWithGogen verifies that multi-line block comments
// generate valid Go code when processed through gogen.
func TestMultiLineBlockCommentWithGogen(t *testing.T) {
	tests := []struct {
		name    string
		comment string
	}{
		{
			name:    "two_line_block",
			comment: "/* Create an iterator for traversing a domain\n   The domain NULL denotes the default domain */",
		},
		{
			name:    "multi_line_with_asterisks",
			comment: "/* Return a pointer to a string containing a description of the error\n * code in the error value ERR.  This function is not thread-safe.  */",
		},
		{
			name:    "single_line_block",
			comment: "/* A simple comment */",
		},
		{
			name:    "javadoc_style",
			comment: "/**\nFoo comment\n*/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg := gogen.NewPackage("", "demo", nil)

			commentGroup := &goast.CommentGroup{
				List: []*goast.Comment{
					{Slash: token.NoPos, Text: tt.comment},
				},
			}

			fn := pkg.NewFunc(nil, "ExampleFunction", nil, nil, false)
			fn.SetComments(pkg, commentGroup)
			fn.BodyStart(pkg).End()

			var buf strings.Builder
			err := gogen.WriteTo(&buf, pkg, "")
			if err != nil {
				t.Fatalf("gogen.WriteTo error: %v", err)
			}

			code := buf.String()
			fset := token.NewFileSet()
			_, err = parser.ParseFile(fset, "test.go", code, parser.ParseComments)
			if err != nil {
				t.Errorf("generated invalid Go code: %v\nCode:\n%s", err, code)
			}
		})
	}
}
