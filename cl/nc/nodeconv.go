package nc

import (
	"errors"

	"github.com/goplus/llcppg/ast"
)

var (
	// ErrSkip is used to skip the node
	ErrSkip = errors.New("skip this node")
)

type NodeConverter interface {
	ConvDecl(file string, decl ast.Decl) (goName, goFile string, err error)
	ConvMacro(file string, macro *ast.Macro) (goName, goFile string, err error)
	ConvEnumItem(decl *ast.EnumTypeDecl, item *ast.EnumItem) (goName string, err error)
	ConvTagExpr(cname string) string
	Lookup(name string) (locFile string, ok bool)
	IsPublic(cname string) bool
}
