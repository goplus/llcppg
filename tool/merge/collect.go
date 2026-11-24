package merge

import (
	"go/ast"
	"go/token"
	"slices"
)

// collect groups the declarations of all platforms.
func (p *mergeIndex) collect() {
	for platform, pkg := range p.files.pkgs {
		for name, file := range pkg.Files {
			p.collectFile(declRef{platform: platform, file: name}, file)
		}
	}
}

// collectFile registers import references before grouping, so the hasher
// can put c.Int and libc.Int in one bucket.
func (p *mergeIndex) collectFile(ref declRef, file *ast.File) {
	imports := collectImports(file)
	for i, decl := range file.Decls {
		ref.declIndex = i
		switch decl := decl.(type) {
		case *ast.GenDecl:
			switch decl.Tok {
			case token.TYPE:
				p.collectTypes(decl, ref, imports)
			case token.VAR, token.CONST:
				p.collectValues(decl, ref, imports)
			}
		case *ast.FuncDecl:
			ref.specIndex = -1
			p.collectSymbol(decl, ref, imports)
		}
	}
}

func (p *mergeIndex) collectTypes(decl *ast.GenDecl, ref declRef, imports map[string]importPath) {
	for i, spec := range decl.Specs {
		ref.specIndex = i
		typ := spec.(*ast.TypeSpec)
		p.collectImportRefs(typ, typ.TypeParams, imports)
		refs, ok := p.types.at(typ).(*[]declRef)
		if !ok {
			p.types.set(typ, &[]declRef{ref})
		} else if slices.Equal(p.directives((*refs)[0]), p.directives(ref)) {
			*refs = append(*refs, ref)
		}
	}
}

// collectValues groups each name separately, but the shared initializer
// stays in the source spec, as in var A, B = pair().
func (p *mergeIndex) collectValues(decl *ast.GenDecl, ref declRef, imports map[string]importPath) {
	for i, spec := range decl.Specs {
		ref.specIndex = i
		value := spec.(*ast.ValueSpec)
		p.collectImportRefs(value, nil, imports)
		for _, name := range value.Names {
			p.collectValue(name.Name, decl, ref)
		}
	}
}

func (p *mergeIndex) collectValue(name string, decl *ast.GenDecl, ref declRef) {
	refs := p.values[name]
	if len(refs) != 0 {
		first := refs[0]
		if p.files.genDecl(first).Tok != decl.Tok ||
			!p.types.equal(p.files.spec(first), decl.Specs[ref.specIndex]) ||
			!slices.Equal(p.directives(first), p.directives(ref)) {
			return
		}
	}
	p.values[name] = append(refs, ref)
}

func (p *mergeIndex) collectSymbol(decl *ast.FuncDecl, ref declRef, imports map[string]importPath) {
	p.collectImportRefs(decl, decl.Type.TypeParams, imports)
	key := funcKey(decl)
	groups := p.symbols[key]
	for i, refs := range groups {
		if p.equalFunc(p.files.funcDecl(refs[0]), decl) && slices.Equal(p.directives(refs[0]), p.directives(ref)) {
			groups[i] = append(refs, ref)
			return
		}
	}
	p.symbols[key] = append(groups, []declRef{ref})
}

// equalFunc compares the receiver and signature types of two functions.
func (p *mergeIndex) equalFunc(a, b *ast.FuncDecl) bool {
	return p.types.equal(recvType(a), recvType(b)) &&
		p.types.equal(a.Type.TypeParams, b.Type.TypeParams) &&
		p.equalFields(a.Type.Params, b.Type.Params) &&
		p.equalFields(a.Type.Results, b.Type.Results)
}

// equalFields compares field types in order, so (a, b int) equals (int, int).
func (p *mergeIndex) equalFields(a, b *ast.FieldList) bool {
	if a.NumFields() != b.NumFields() {
		return false
	}
	if a.NumFields() == 0 {
		return true
	}
	i, n := 0, 0 // field and name index in b
	for _, field := range a.List {
		for range max(1, len(field.Names)) {
			other := b.List[i]
			if !p.types.equal(field.Type, other.Type) {
				return false
			}
			if n++; n == max(1, len(other.Names)) {
				i, n = i+1, 0
			}
		}
	}
	return true
}

func recvType(decl *ast.FuncDecl) ast.Expr {
	if decl.Recv == nil {
		return nil
	}
	return decl.Recv.List[0].Type
}

// directives returns the directives of a declaration, such as //go:linkname.
// Declarations only match when their directives are equal.
func (p *mergeIndex) directives(ref declRef) []string {
	switch decl := p.files.decl(ref).(type) {
	case *ast.FuncDecl:
		return parseComments(decl.Doc).directives
	case *ast.GenDecl:
		var directives []string
		if len(decl.Specs) == 1 {
			// A directive above type Callback func() is attached to the GenDecl.
			directives = parseComments(decl.Doc).directives
		}
		doc, _ := specComments(decl.Specs[ref.specIndex])
		return append(directives, parseComments(doc).directives...)
	}
	return nil
}

// funcKey returns the name of a function, or T.M for a method.
func funcKey(decl *ast.FuncDecl) string {
	name := decl.Name.Name
	if decl.Recv == nil {
		return name
	}
	recv := ast.Unparen(decl.Recv.List[0].Type)
	if ptr, ok := recv.(*ast.StarExpr); ok {
		recv = ast.Unparen(ptr.X)
	}
	switch typ := recv.(type) {
	case *ast.IndexExpr:
		recv = typ.X
	case *ast.IndexListExpr:
		recv = typ.X
	}
	return recv.(*ast.Ident).Name + "." + name
}
