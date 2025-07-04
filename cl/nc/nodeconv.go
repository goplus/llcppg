package nc

import (
	"errors"

	"github.com/goplus/llcppg/ast"
)

var (
	// ErrSkip is used to skip the node
	ErrSkip = errors.New("skip this node")
)

type Condition struct {
	OS   string // OS,like darwin,linux,windows
	Arch string // Architecture,like amd64,arm64
}

type GoFile struct {
	FileName  string     // Go file name,like cJSON.go ini_darwin_amd64.go
	Condition *Condition // Condition for the given file,if no condition,it is nil
}

type NodeConverter interface {
	ConvDecl(file string, decl ast.Decl) (goName string, goFile *GoFile, err error)
	ConvMacro(file string, macro *ast.Macro) (goName string, goFile *GoFile, err error)
	ConvEnumItem(decl *ast.EnumTypeDecl, item *ast.EnumItem) (goName string, err error)
	ConvTagExpr(cname string) string
	Lookup(name string) (locFile string, ok bool)
	IsPublic(cname string) bool
}
