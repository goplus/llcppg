// Minimal reproduction for issue #618: Multi-line block comment bug with gogen
//
// This demonstrates how incorrectly structured ast.Comment nodes cause
// invalid Go code generation with gogen v1.20.2.
//
// Run: go run reproduce_comment_bug.go

package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"strings"

	"github.com/goplus/gogen"
)

// This creates BUGGY comments - splits block comment across multiple nodes
func createBuggyComments() *ast.CommentGroup {
	// This is what llcppg's buggy ParseComment() does:
	// It splits multi-line block comments by newline
	return &ast.CommentGroup{
		List: []*ast.Comment{
			{Text: "/* Create an iterator for traversing a domain\n"},
			{Text: "   The domain NULL denotes the default domain */\n"},
		},
	}
}

// This creates CORRECT comments - keeps block comment as single node
func createCorrectComments() *ast.CommentGroup {
	// This is what should be done:
	// Keep the entire block comment as a single ast.Comment node
	return &ast.CommentGroup{
		List: []*ast.Comment{
			{Text: "/* Create an iterator for traversing a domain\n   The domain NULL denotes the default domain */\n"},
		},
	}
}

func generateCode(comments *ast.CommentGroup, label string) {
	fmt.Println(label)
	fmt.Println(strings.Repeat("=", len(label)))

	pkg := gogen.NewPackage("", "test", nil)

	// Create a type declaration with the comment
	tyDecl := pkg.NewTypeDefs()
	spec := tyDecl.NewType("ExampleType")
	spec.InitType(pkg, types.Typ[types.Int])

	// Attach the comment as documentation
	if comments != nil {
		spec.SetComments(pkg, comments)
	}

	// Get the generated AST and try to format it
	file := pkg.ASTFile()

	var buf bytes.Buffer
	fset := token.NewFileSet()
	err := format.Node(&buf, fset, file)
	if err != nil {
		fmt.Printf("❌ ERROR: %v\n\n", err)
		fmt.Println("This demonstrates the bug with gogen v1.20.2!")
	} else {
		fmt.Println("✅ Generated Go code:")
		fmt.Println(buf.String())
	}
	fmt.Println()
}

func main() {
	fmt.Println("=== Gogen Comment Bug Reproduction (Issue #618) ===\n")

	fmt.Println("This reproduces the bug WITHOUT involving llcppg - pure gogen usage.\n")

	// Test 1: Buggy comment structure (what llcppg currently generates)
	fmt.Println("Test 1: BUGGY comment structure")
	fmt.Println("--------------------------------")
	fmt.Println("Comment split into 2 nodes:")
	buggy := createBuggyComments()
	for i, c := range buggy.List {
		fmt.Printf("  [%d]: %q\n", i, c.Text)
	}
	fmt.Println()

	generateCode(buggy, "Generated Go code (will have malformed comments):")

	// Test 2: Correct comment structure
	fmt.Println("\nTest 2: CORRECT comment structure")
	fmt.Println("----------------------------------")
	fmt.Println("Comment as single node:")
	correct := createCorrectComments()
	for i, c := range correct.List {
		fmt.Printf("  [%d]: %q\n", i, c.Text)
	}
	fmt.Println()

	generateCode(correct, "Generated Go code (properly formatted comments):")

	fmt.Println("=== Explanation ===")
	fmt.Println("Go's ast.Comment documentation states each Comment represents")
	fmt.Println("'a single //-style or /*-style comment' - emphasis on SINGLE.")
	fmt.Println()
	fmt.Println("Block comments /* ... */ must be ONE ast.Comment node, not split")
	fmt.Println("across multiple nodes. When gogen v1.20.2 prints malformed comment")
	fmt.Println("structures, it generates invalid Go syntax.")
	fmt.Println()
	fmt.Println("This is the root cause of test failures in:")
	fmt.Println("  - TestFromTestdata/gettext")
	fmt.Println("  - TestFromTestdata/gpgerror")
	fmt.Println("  - TestFromTestdata/issue507")
	fmt.Println("  - TestFromTestdata/keepcomment")
}
