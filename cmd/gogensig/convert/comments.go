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
	if doc != nil && doc.List != nil {
		for _, comment := range doc.List {
			goDoc.List = append(goDoc.List,
				&goast.Comment{
					Slash: token.NoPos, Text: comment.Text,
				},
			)
		}
	}
	return goDoc
}
