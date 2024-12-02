package convert

import (
	goast "go/ast"
	"strings"
)

const (
	TYPEC = "// llgo:type C"
)

func NewFuncDocComments(funcName string, goFuncName string) *goast.CommentGroup {
	fields := strings.FieldsFunc(goFuncName, func(r rune) bool {
		return r == '.'
	})
	txt := "//go:linkname " + goFuncName + " " + "C." + funcName
	if len(fields) > 1 {
		txt = "// llgo:link " + goFuncName + " " + "C." + funcName
	}
	comment := goast.Comment{Text: txt}
	commentGroup := goast.CommentGroup{List: []*goast.Comment{&comment}}
	return &commentGroup
}

func NewTypecDocComments() *goast.CommentGroup {
	return &goast.CommentGroup{
		List: []*goast.Comment{
			{Text: TYPEC},
		}}
}
