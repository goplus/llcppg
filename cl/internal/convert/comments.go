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
	if doc == nil || len(doc.List) == 0 {
		return goDoc
	}
	if strings.HasPrefix(doc.List[0].Text, "/*") {
		// Block comment: reassemble into a single Go ast.Comment node.
		// This handles both correctly formed single-node block comments
		// and legacy split block comments from older llcppsigfetch builds.
		var b strings.Builder
		for _, comment := range doc.List {
			b.WriteString(comment.Text)
		}
		text := strings.TrimRight(b.String(), "\n")
		goDoc.List = []*goast.Comment{{Slash: token.NoPos, Text: text}}
	} else {
		// Line comments: each comment is a separate node.
		for _, comment := range doc.List {
			text := strings.TrimRight(comment.Text, "\n")
			if text != "" {
				goDoc.List = append(goDoc.List, &goast.Comment{
					Slash: token.NoPos, Text: text,
				})
			}
		}
	}
	return goDoc
}

// hasBlockComment reports whether the comment group contains a block comment (/* ... */).
func hasBlockComment(doc *goast.CommentGroup) bool {
	for _, c := range doc.List {
		if strings.HasPrefix(c.Text, "/*") {
			return true
		}
	}
	return false
}
