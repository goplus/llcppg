package visitor

import (
	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/llcppg"
)

type DocVisitor interface {
	VisitStart(path string, incPath string, isSys bool, fileType llcppg.FileType)
	Visit(node ast.Node)
	VisitFuncDecl(funcDecl *ast.FuncDecl)
	VisitDone(path string)
	VisitStruct(structName *ast.Ident, fields *ast.FieldList, typeDecl *ast.TypeDecl)
	//VisitClass(className *ast.Ident, fields *ast.FieldList, typeDecl *ast.TypeDecl)
	//VisitMethod(className *ast.Ident, method *ast.FuncDecl, typeDecl *ast.TypeDecl)
	VisitUnion(unionName *ast.Ident, fields *ast.FieldList, typeDecl *ast.TypeDecl)
	VisitEnumTypeDecl(enumTypeDecl *ast.EnumTypeDecl)
	VisitTypedefDecl(typedefDecl *ast.TypedefDecl)
	VisitMacro(macro *ast.Macro)
}

type DocVisitorList struct {
	VisitorList []DocVisitor
}

func NewDocVisitorList(visitorList []DocVisitor) *DocVisitorList {
	return &DocVisitorList{VisitorList: visitorList}
}

func (p *DocVisitorList) Visit(node ast.Node, path string, incPath string, isSys bool, fileType llcppg.FileType) bool {
	for _, v := range p.VisitorList {
		v.VisitStart(path, incPath, isSys, fileType)
		v.Visit(node)
		v.VisitDone(path)
	}
	return true
}
