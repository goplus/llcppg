package merge

import "go/ast"

// declRef locates a declaration in the input files.
// Specs can share a declaration index, as in type (A int; B int).
type declRef struct {
	platform  int
	file      string
	declIndex int
	specIndex int // -1 for a function
}

// files gives access to the input files of the platform packages in Merge order.
type files struct {
	pkgs []*PlatformPackage
}

func (p files) len() int {
	return len(p.pkgs)
}

func (p files) decl(ref declRef) ast.Decl {
	return p.pkgs[ref.platform].Files[ref.file].Decls[ref.declIndex]
}

func (p files) genDecl(ref declRef) *ast.GenDecl {
	return p.decl(ref).(*ast.GenDecl)
}

func (p files) spec(ref declRef) ast.Spec {
	return p.genDecl(ref).Specs[ref.specIndex]
}

func (p files) file(ref declRef) *ast.File {
	return p.pkgs[ref.platform].Files[ref.file]
}

func (p files) funcDecl(ref declRef) *ast.FuncDecl {
	return p.decl(ref).(*ast.FuncDecl)
}
