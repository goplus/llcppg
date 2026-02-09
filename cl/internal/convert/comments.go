package convert

import (
	goast "go/ast"
	"go/token"
	"strings"

	"github.com/goplus/llcppg/ast"
)

const (
	TYPEC = "// llgo:type C"
)

func NewFuncDocComment(funcName string, goFuncName string) *goast.Comment {
	fields := strings.FieldsFunc(goFuncName, func(r rune) bool {
		return r == '.'
	})
	txt := "//go:linkname " + goFuncName + " " + "C." + funcName
	if len(fields) > 1 {
		txt = "// llgo:link " + goFuncName + " " + "C." + funcName
	}
	return &goast.Comment{Text: txt}
}

func NewTypecDocComment() *goast.Comment {
	return &goast.Comment{Text: TYPEC}
}

func NewCommentGroup(comments ...*goast.Comment) *goast.CommentGroup {
	return &goast.CommentGroup{List: comments}
}

func NewCommentGroupFromC(doc *ast.CommentGroup) *goast.CommentGroup {
	goDoc := &goast.CommentGroup{}
	if doc == nil || doc.List == nil {
		return goDoc
	}

	// Process comments, merging multi-line block comments into single nodes.
	// Go's ast.Comment requires block comments (/* ... */) to be a single node,
	// but line comments (// ...) should be separate nodes per line.
	i := 0
	for i < len(doc.List) {
		comment := doc.List[i]
		text := strings.TrimRight(comment.Text, "\n")

		// Check if this is the start of a block comment
		if strings.HasPrefix(text, "/*") {
			// If the block comment is complete (contains */), add as single node
			if strings.Contains(text, "*/") {
				goDoc.List = append(goDoc.List, &goast.Comment{
					Slash: token.NoPos, Text: text,
				})
				i++
				continue
			}

			// Multi-line block comment: merge all lines until we find */
			var lines []string
			lines = append(lines, text)
			i++

			for i < len(doc.List) {
				nextText := strings.TrimRight(doc.List[i].Text, "\n")
				lines = append(lines, nextText)
				i++
				if strings.Contains(nextText, "*/") {
					break
				}
			}

			// Join all lines with newlines to form complete block comment
			mergedComment := strings.Join(lines, "\n")
			goDoc.List = append(goDoc.List, &goast.Comment{
				Slash: token.NoPos, Text: mergedComment,
			})
		} else {
			// Line comment or other: add as-is (without trailing newline)
			goDoc.List = append(goDoc.List, &goast.Comment{
				Slash: token.NoPos, Text: text,
			})
			i++
		}
	}
	return goDoc
}
