package merge

import "go/ast"

// mergeIndex groups declarations that are equal across platforms.
// A group is shared when it has one declaration from every platform.
//
// Groups only hold declRefs into the input files, so the index stays small
// for large packages; nodes are looked up through files when needed.
type mergeIndex struct {
	files      files
	types      astMap                    // type spec => *[]declRef
	symbols    map[string][][]declRef    // function key, such as T.M => groups of equal signatures
	values     map[string][]declRef      // const or var name => declarations
	importRefs map[*ast.Ident]importPath // package qualifier => import path
}

func newMergeIndex(pkgs []*PlatformPackage) *mergeIndex {
	p := &mergeIndex{
		files:      files{pkgs: pkgs},
		symbols:    make(map[string][][]declRef),
		values:     make(map[string][]declRef),
		importRefs: make(map[*ast.Ident]importPath),
	}
	p.types.hasher.identName = p.identName
	return p
}

// identName resolves a package qualifier so c.Int and libc.Int can compare equal.
func (p *mergeIndex) identName(ident *ast.Ident) (name string, imported bool) {
	if path, ok := p.importRefs[ident]; ok {
		return string(path), true
	}
	return astIdentName(ident)
}

// shared reports whether refs holds one declaration from every platform.
func (p *mergeIndex) shared(refs []declRef) bool {
	return len(refs) == p.files.len()
}
