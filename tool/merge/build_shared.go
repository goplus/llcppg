package merge

import (
	"cmp"
	"go/ast"
	"iter"
	"slices"
)

// sharedGroup is a group with one declaration from every platform.
// It is printed from the chosen declaration, the one of the first file name, then the first platform.
type sharedGroup struct {
	declRef           // chosen declaration
	refs    []declRef // all declarations of the group
}

// sharedDecl is a shared group materialized for output.
// Its comments are reduced to the ones all platforms agree on.
type sharedDecl struct {
	decl     ast.Decl
	spec     ast.Spec          // nil for a function
	blockDoc *ast.CommentGroup // doc of a grouped declaration, as in const (...)
}

// buildShared returns the shared declarations by source file name.
// Groups stay references until they are printed, one source file at a time.
func (p *mergeIndex) buildShared() map[string][]ast.Decl {
	var groups []sharedGroup
	add := func(refs []declRef) {
		if p.shared(refs) {
			groups = append(groups, sharedGroup{declRef: chooseRef(refs), refs: refs})
		}
	}
	p.types.iterate(func(_ ast.Node, value any) {
		add(*value.(*[]declRef))
	})
	for _, symbols := range p.symbols {
		for _, refs := range symbols {
			add(refs)
		}
	}
	for _, refs := range p.values {
		add(refs)
	}

	// Files contribute in platform order, each in source order.
	slices.SortFunc(groups, func(a, b sharedGroup) int {
		return cmp.Or(cmp.Compare(a.file, b.file), cmp.Compare(a.platform, b.platform),
			cmp.Compare(a.declIndex, b.declIndex), cmp.Compare(a.specIndex, b.specIndex))
	})
	files := make(map[string][]ast.Decl)
	for len(groups) != 0 {
		first := groups[0]
		n := 1
		for n < len(groups) && groups[n].file == first.file && groups[n].platform == first.platform {
			n++
		}
		files[first.file] = append(files[first.file], p.buildSharedFile(groups[:n])...)
		groups = groups[n:]
	}
	return files
}

// buildSharedFile prints the shared groups chosen from one source file.
// Specs from one declaration stay together, as in type (A int; B int).
func (p *mergeIndex) buildSharedFile(groups []sharedGroup) []ast.Decl {
	var list []ast.Decl
	for i, group := range groups {
		shared := p.shareDecl(group)
		if shared.spec == nil {
			list = append(list, shared.decl)
			continue
		}
		blockDoc := declBlockDoc(shared.decl)
		if i > 0 && groups[i-1].declIndex == group.declIndex {
			gen := list[len(list)-1].(*ast.GenDecl)
			gen.Specs = append(gen.Specs, shared.spec)
			// Include each member's source title, since another platform can split const (A = 1; B = 2).
			if blockDoc != nil && gen.Doc != nil && gen.Doc != shared.blockDoc {
				gen.Doc = commonDoc(gen.Doc, slices.Values([]*ast.CommentGroup{shared.blockDoc}))
			}
			continue
		}
		gen := *shared.decl.(*ast.GenDecl)
		if blockDoc != nil {
			gen.Doc = shared.blockDoc
		}
		gen.Specs = []ast.Spec{shared.spec}
		list = append(list, &gen)
	}
	sortDecls(p.files.file(groups[0].declRef), list)
	return list
}

func (p *mergeIndex) shareDecl(group sharedGroup) sharedDecl {
	decl := p.files.decl(group.declRef)
	doc, line, onSpec := declComments(decl, group.specIndex)
	sharedDoc := commonDoc(doc, p.peerComments(group, declDoc))
	sharedLine := commonLine(line, p.peerComments(group, declLine))
	sharedBlockDoc := commonDoc(declBlockDoc(decl), p.peerComments(group, blockDoc))
	shared := sharedDecl{decl: decl, blockDoc: sharedBlockDoc}
	if !onSpec && sharedDoc != doc {
		shared.decl = withDoc(decl, sharedDoc)
	}
	if group.specIndex < 0 {
		return shared
	}
	spec := p.files.spec(group.declRef)
	specDoc, _ := specComments(spec)
	if onSpec {
		specDoc = sharedDoc
	}
	switch spec := spec.(type) {
	case *ast.TypeSpec:
		shared.spec = p.shareTypeSpec(spec, group.refs, specDoc, sharedLine)
	case *ast.ValueSpec:
		shared.spec = shareValueSpec(spec, specDoc, sharedLine)
	}
	return shared
}

// peerComments yields a comment of each declaration in group other than the chosen one.
// It reads the input on demand rather than gathering the comments.
func (p *mergeIndex) peerComments(group sharedGroup, comment func(decl ast.Decl, specIndex int) *ast.CommentGroup) iter.Seq[*ast.CommentGroup] {
	return func(yield func(*ast.CommentGroup) bool) {
		for _, ref := range group.refs {
			if ref != group.declRef && !yield(comment(p.files.decl(ref), ref.specIndex)) {
				return
			}
		}
	}
}

func (p *mergeIndex) shareTypeSpec(spec *ast.TypeSpec, refs []declRef, doc, line *ast.CommentGroup) *ast.TypeSpec {
	typ := spec.Type
	if st, ok := ast.Unparen(typ).(*ast.StructType); ok && st.Fields != nil && len(refs) > 1 {
		typ = p.shareFieldComments(typ, refs)
	}
	if typ == spec.Type && doc == spec.Doc && line == spec.Comment {
		return spec
	}
	shared := *spec
	shared.Doc, shared.Type, shared.Comment = doc, typ, line
	return &shared
}

func shareValueSpec(spec *ast.ValueSpec, doc, line *ast.CommentGroup) *ast.ValueSpec {
	if doc == spec.Doc && line == spec.Comment {
		return spec
	}
	shared := *spec
	shared.Doc, shared.Comment = doc, line
	return &shared
}

// shareFieldComments reduces the field comments of a struct type to those of all refs.
// Type matching already aligns fields in order, including an embedding such as struct { T }.
func (p *mergeIndex) shareFieldComments(typ ast.Expr, refs []declRef) ast.Expr {
	st := ast.Unparen(typ).(*ast.StructType)
	var fields []*ast.Field
	for i, field := range st.Fields.List {
		if field.Doc == nil && field.Comment == nil {
			continue
		}
		doc := commonDoc(field.Doc, p.fieldComments(refs, i, fieldDoc))
		line := commonLine(field.Comment, p.fieldComments(refs, i, fieldLine))
		if doc == field.Doc && line == field.Comment {
			continue
		}
		if fields == nil {
			fields = slices.Clone(st.Fields.List)
		}
		shared := *field
		shared.Doc, shared.Comment = doc, line
		fields[i] = &shared
	}
	if fields == nil {
		return typ
	}
	list := *st.Fields
	list.List = fields
	shared := *st
	shared.Fields = &list
	// Copy enclosing parentheses rather than modify the input, as in type T (struct{ X int }).
	result := typ
	inner := &result
	for {
		paren, ok := (*inner).(*ast.ParenExpr)
		if !ok {
			break
		}
		copied := *paren
		*inner = &copied
		inner = &copied.X
	}
	*inner = &shared
	return result
}

// fieldComments yields a comment of field i in the struct type of each ref.
func (p *mergeIndex) fieldComments(refs []declRef, i int, comment func(*ast.Field) *ast.CommentGroup) iter.Seq[*ast.CommentGroup] {
	return func(yield func(*ast.CommentGroup) bool) {
		for _, ref := range refs {
			st := ast.Unparen(p.files.spec(ref).(*ast.TypeSpec).Type).(*ast.StructType)
			if !yield(comment(st.Fields.List[i])) {
				return
			}
		}
	}
}

// chooseRef picks the declaration of the first file name, then the first platform.
func chooseRef(refs []declRef) declRef {
	return slices.MinFunc(refs, func(a, b declRef) int {
		return cmp.Or(cmp.Compare(a.file, b.file), cmp.Compare(a.platform, b.platform))
	})
}

func withDoc(decl ast.Decl, doc *ast.CommentGroup) ast.Decl {
	switch decl := decl.(type) {
	case *ast.FuncDecl:
		shared := *decl
		shared.Doc = doc
		return &shared
	case *ast.GenDecl:
		shared := *decl
		shared.Doc = doc
		return &shared
	}
	return decl
}

// declComments returns the comments of a function or spec, and whether doc belongs to the spec.
// A doc above const A = 1 belongs to A, but a doc above const (...) belongs to the block.
func declComments(decl ast.Decl, specIndex int) (doc, line *ast.CommentGroup, onSpec bool) {
	switch decl := decl.(type) {
	case *ast.FuncDecl:
		return decl.Doc, nil, false
	case *ast.GenDecl:
		doc, line = specComments(decl.Specs[specIndex])
		if decl.Doc != nil && len(decl.Specs) == 1 && !decl.Lparen.IsValid() {
			return decl.Doc, line, false
		}
		return doc, line, true
	}
	return nil, nil, false
}

func declDoc(decl ast.Decl, specIndex int) *ast.CommentGroup {
	doc, _, _ := declComments(decl, specIndex)
	return doc
}

func declLine(decl ast.Decl, specIndex int) *ast.CommentGroup {
	_, line, _ := declComments(decl, specIndex)
	return line
}

func blockDoc(decl ast.Decl, _ int) *ast.CommentGroup {
	return declBlockDoc(decl)
}

func fieldDoc(field *ast.Field) *ast.CommentGroup  { return field.Doc }
func fieldLine(field *ast.Field) *ast.CommentGroup { return field.Comment }

// declBlockDoc returns the doc of a grouped declaration, as in const (...).
func declBlockDoc(decl ast.Decl) *ast.CommentGroup {
	if decl, ok := decl.(*ast.GenDecl); ok && (decl.Lparen.IsValid() || len(decl.Specs) > 1) {
		return decl.Doc
	}
	return nil
}
