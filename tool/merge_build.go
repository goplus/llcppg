package tool

import (
	"cmp"
	"go/ast"
	"go/scanner"
	"go/token"
	"maps"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

func (m *mergeIndex) buildFiles(pkgName string, tags [][]string) map[string]*ast.File {
	platformFiles := make([]string, len(tags))
	for i, platform := range tags {
		platformFiles[i] = "z_" + platform[0] + "_" + platform[1] + ".go"
	}

	var headerComments []string
	if header := []byte(GeneratedHeader); len(header) != 0 {
		var scan scanner.Scanner
		scan.Init(token.NewFileSet().AddFile("", -1, len(header)), header, nil, scanner.ScanComments)
		for _, tok, text := scan.Scan(); tok == token.COMMENT; _, tok, text = scan.Scan() {
			headerComments = append(headerComments, text)
		}
	}

	decls := m.buildDecls(platformFiles)
	files := make(map[string]*ast.File, len(decls))
	for name, list := range decls {
		file := &ast.File{Name: ast.NewIdent(pkgName), Decls: list}
		m.addImports(file)
		if len(headerComments) != 0 {
			comments := make([]*ast.Comment, len(headerComments))
			for i, text := range headerComments {
				comments[i] = &ast.Comment{Text: text}
			}
			file.Doc = &ast.CommentGroup{List: comments}
		}
		files[name] = file
	}
	for i, name := range platformFiles {
		if file := files[name]; file != nil {
			// The printer separates //go:build linux && amd64 from the package clause.
			if file.Doc == nil {
				file.Doc = &ast.CommentGroup{}
			}
			file.Doc.List = slices.Insert(file.Doc.List, 0, &ast.Comment{
				Text: "//go:build " + tags[i][0] + " && " + tags[i][1],
			})
		}
	}
	return files
}

func (m *mergeIndex) addImports(file *ast.File) {
	names := make(map[importPath]string)
	replacements := make(map[ast.Node]ast.Node)
	linkname := false
	ast.PreorderStack(file, nil, func(node ast.Node, stack []ast.Node) bool {
		switch node := node.(type) {
		case *ast.Ident:
			pkgPath, ok := m.importRefs[node]
			if !ok {
				return true
			}
			// Declarations from earlier platforms come first, so c.Int can set the name for libc.Int.
			name, ok := names[pkgPath]
			if !ok {
				names[pkgPath] = node.Name
			} else if node.Name != name {
				replacements[node] = &ast.Ident{NamePos: node.NamePos, Name: name, Obj: node.Obj}
				for _, parent := range slices.Backward(stack) {
					if _, ok := replacements[parent]; ok {
						break
					}
					replacements[parent] = nil
				}
			}
		case *ast.CommentGroup:
			var comments commentgroup
			comments.load(node, nil)
			for _, directive := range comments.directives {
				if strings.HasPrefix(directive, "go:linkname ") {
					linkname = true
					break
				}
			}
		}
		return true
	})
	// A //go:linkname Foo C.foo directive needs a blank import when there is no unsafe.Pointer use.
	if _, ok := names["unsafe"]; !ok && linkname {
		names["unsafe"] = "_"
	}
	if len(names) == 0 {
		return
	}
	if len(replacements) != 0 {
		*file = *copyImportRefs(reflect.ValueOf(file), replacements).Interface().(*ast.File)
	}
	// TODO: Resolve alias collisions when source files land in one output file.
	// For example, c "example.org/a" and c "example.org/b" cannot coexist.
	decl := &ast.GenDecl{Tok: token.IMPORT}
	for _, pkgPath := range slices.Sorted(maps.Keys(names)) {
		spec := &ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: strconv.Quote(string(pkgPath))}}
		if name := names[pkgPath]; pkgPath != "unsafe" || name != "unsafe" {
			spec.Name = ast.NewIdent(name)
		}
		decl.Specs = append(decl.Specs, spec)
		file.Imports = append(file.Imports, spec)
	}
	file.Decls = append([]ast.Decl{decl}, file.Decls...)
}

// Copy the ancestors of a renamed reference so, for example, X libc.Int in the input retains its name.
// A nil replacement marks an ancestor, and an invalid return value means the input can be reused.
func copyImportRefs(value reflect.Value, replacements map[ast.Node]ast.Node) reflect.Value {
	switch value.Kind() {
	case reflect.Interface:
		if !value.IsNil() {
			return copyImportRefs(value.Elem(), replacements)
		}
	case reflect.Pointer:
		node, ok := value.Interface().(ast.Node)
		if !ok {
			break
		}
		replacement, ok := replacements[node]
		if !ok {
			break
		}
		if replacement != nil {
			return reflect.ValueOf(replacement)
		}
		source := value.Elem()
		result := reflect.New(source.Type())
		result.Elem().Set(source)
		for i := range source.NumField() {
			if field := copyImportRefs(source.Field(i), replacements); field.IsValid() {
				result.Elem().Field(i).Set(field)
			}
		}
		return result
	case reflect.Slice:
		var result reflect.Value
		for i := range value.Len() {
			if item := copyImportRefs(value.Index(i), replacements); item.IsValid() {
				if !result.IsValid() {
					result = reflect.MakeSlice(value.Type(), value.Len(), value.Len())
					reflect.Copy(result, value)
				}
				result.Index(i).Set(item)
			}
		}
		return result
	}
	return reflect.Value{}
}

// buildDecls uses platform file names in tag order, such as z_darwin_arm64.go before z_linux_amd64.go.
func (m *mergeIndex) buildDecls(platformFiles []string) map[string][]ast.Decl {
	type representative struct {
		mergeDeclRef
		decl     ast.Decl
		spec     ast.Spec
		blockDoc *ast.CommentGroup
	}
	platformCount := len(m.files)
	var representatives []representative
	selectShared := func(first mergeDeclRef, aliases []mergeDeclRef) {
		if len(aliases)+1 != platformCount {
			return
		}
		chosen := first
		for _, alias := range aliases {
			if alias.file < chosen.file || (alias.file == chosen.file && alias.platform < chosen.platform) {
				chosen = alias
			}
		}
		decl := m.files[chosen.platform][chosen.file].Decls[chosen.declIndex]
		selected := representative{mergeDeclRef: chosen, decl: decl}
		if chosen.specIndex >= 0 {
			selected.spec = decl.(*ast.GenDecl).Specs[chosen.specIndex]
		}
		doc, line, onSpec := sharedDeclComments(decl, chosen.specIndex)
		blockDoc := declBlockDoc(decl)
		var comments, block, other commentgroup
		comments.load(doc, line)
		block.load(blockDoc, nil)
		original, originalBlock := comments, block
		for i := 0; i <= len(aliases) && (len(comments.doc) > 0 || comments.line != nil || len(block.doc) > 0); i++ {
			ref := first
			if i > 0 {
				ref = aliases[i-1]
			}
			if ref == chosen {
				continue
			}
			otherDecl := m.files[ref.platform][ref.file].Decls[ref.declIndex]
			if len(comments.doc) > 0 || comments.line != nil {
				otherDoc, otherLine, _ := sharedDeclComments(otherDecl, ref.specIndex)
				other.load(otherDoc, otherLine)
				comments.intersect(other)
			}
			if len(block.doc) > 0 {
				other.load(declBlockDoc(otherDecl), nil)
				block.intersect(other)
			}
		}
		if !block.equalDoc(originalBlock) {
			blockDoc = block.withDirectives(blockDoc)
		}
		selected.blockDoc = blockDoc
		mergedDoc := doc
		if !comments.equalDoc(original) {
			mergedDoc = comments.withDirectives(doc)
		}
		if !onSpec && mergedDoc != doc {
			switch decl := decl.(type) {
			case *ast.FuncDecl:
				shared := *decl
				shared.Doc = mergedDoc
				selected.decl = &shared
			case *ast.GenDecl:
				shared := *decl
				shared.Doc = mergedDoc
				selected.decl = &shared
			}
		}
		switch spec := selected.spec.(type) {
		case *ast.TypeSpec:
			sharedType := spec.Type
			if st, ok := ast.Unparen(spec.Type).(*ast.StructType); ok && st.Fields != nil && len(aliases) > 0 {
				peers := make([]*ast.StructType, len(aliases)+1)
				for i := range peers {
					ref := first
					if i > 0 {
						ref = aliases[i-1]
					}
					peerDecl := m.files[ref.platform][ref.file].Decls[ref.declIndex].(*ast.GenDecl)
					peerSpec := peerDecl.Specs[ref.specIndex].(*ast.TypeSpec)
					peers[i] = ast.Unparen(peerSpec.Type).(*ast.StructType)
				}
				sharedType = intersectStructFieldComments(spec.Type, peers)
			}
			if sharedType != spec.Type || spec.Comment != comments.line || (onSpec && spec.Doc != mergedDoc) {
				shared := *spec
				if onSpec {
					shared.Doc = mergedDoc
				}
				shared.Type = sharedType
				shared.Comment = comments.line
				selected.spec = &shared
			}
		case *ast.ValueSpec:
			if spec.Comment != comments.line || (onSpec && spec.Doc != mergedDoc) {
				shared := *spec
				if onSpec {
					shared.Doc = mergedDoc
				}
				shared.Comment = comments.line
				selected.spec = &shared
			}
		}
		representatives = append(representatives, selected)
	}
	m.types.iterate(func(_ ast.Node, value any) {
		refs := *value.(*[]mergeDeclRef)
		selectShared(refs[0], refs[1:])
	})
	for _, variants := range m.symbols {
		for _, symbol := range variants {
			selectShared(symbol.decl, symbol.aliases)
		}
	}
	for _, refs := range m.values {
		selectShared(refs[0], refs[1:])
	}

	outputs := make(map[string][]ast.Decl)
	slices.SortFunc(representatives, func(a, b representative) int {
		return cmp.Or(cmp.Compare(a.file, b.file), cmp.Compare(a.platform, b.platform),
			cmp.Compare(a.declIndex, b.declIndex), cmp.Compare(a.specIndex, b.specIndex))
	})
	for i := 0; i < len(representatives); {
		source := representatives[i]
		file := m.files[source.platform][source.file]
		var decls []ast.Decl
		lastDeclIndex := -1
		for ; i < len(representatives); i++ {
			ref := representatives[i]
			if ref.file != source.file || ref.platform != source.platform {
				break
			}
			decl := ref.decl
			if ref.specIndex >= 0 {
				original := decl.(*ast.GenDecl)
				blockDoc := declBlockDoc(original)
				// Keep specs from one declaration together, as in type (A int; B int).
				if ref.declIndex == lastDeclIndex {
					group := decls[len(decls)-1].(*ast.GenDecl)
					group.Specs = append(group.Specs, ref.spec)
					// Include each member's source title, since another platform can split const (A = 1; B = 2).
					if blockDoc != nil && group.Doc != nil && group.Doc != ref.blockDoc {
						var comments, other commentgroup
						comments.load(group.Doc, nil)
						original := comments
						other.load(ref.blockDoc, nil)
						comments.intersect(other)
						if !comments.equalDoc(original) {
							group.Doc = comments.withDirectives(group.Doc)
						}
					}
					continue
				}
				group := *original
				if blockDoc != nil {
					group.Doc = ref.blockDoc
				}
				group.Specs = []ast.Spec{ref.spec}
				decl = &group
			}
			decls = append(decls, decl)
			lastDeclIndex = ref.declIndex
		}
		sortMergeDecls(file, decls)
		outputs[source.file] = append(outputs[source.file], decls...)
	}
	for platform, output := range platformFiles {
		files := m.files[platform]
		for _, name := range slices.Sorted(maps.Keys(files)) {
			file := files[name]
			decls := filter(file, func(i, j int) bool {
				ref := mergeDeclRef{platform: platform, file: name, declIndex: i, specIndex: j}
				switch decl := file.Decls[i].(type) {
				case *ast.FuncDecl:
					for _, group := range m.symbols[mergeFuncKey(decl)] {
						if len(group.aliases)+1 == platformCount && (group.decl == ref || slices.Contains(group.aliases, ref)) {
							return false
						}
					}
				case *ast.GenDecl:
					if decl.Tok == token.IMPORT || len(decl.Specs) == 0 {
						return false
					}
					switch decl.Tok {
					case token.TYPE:
						refs, ok := m.types.at(decl.Specs[j]).(*[]mergeDeclRef)
						return !ok || len(*refs) != platformCount || !slices.Contains(*refs, ref)
					case token.VAR, token.CONST:
						for _, ident := range decl.Specs[j].(*ast.ValueSpec).Names {
							refs := m.values[ident.Name]
							if len(refs) == platformCount && slices.Contains(refs, ref) {
								return false
							}
						}
					}
				}
				return true
			})
			if len(decls) > 0 {
				sortMergeDecls(file, decls)
				outputs[output] = append(outputs[output], decls...)
			}
		}
	}
	return outputs
}

func intersectStructFieldComments(typ ast.Expr, peers []*ast.StructType) ast.Expr {
	st := ast.Unparen(typ).(*ast.StructType)
	var fields []*ast.Field
	for index, field := range st.Fields.List {
		if field.Doc == nil && field.Comment == nil {
			continue
		}
		var comments, other commentgroup
		comments.load(field.Doc, field.Comment)
		original := comments
		// Type matching already aligns fields in order, including an embedding such as struct { T }.
		for i := 0; i < len(peers) && (len(comments.doc) > 0 || comments.line != nil); i++ {
			peer := peers[i].Fields.List[index]
			other.load(peer.Doc, peer.Comment)
			comments.intersect(other)
		}
		doc := field.Doc
		if !comments.equalDoc(original) {
			doc = comments.withDirectives(doc)
		}
		if doc == field.Doc && comments.line == field.Comment {
			continue
		}
		if fields == nil {
			fields = slices.Clone(st.Fields.List)
		}
		shared := *field
		shared.Doc, shared.Comment = doc, comments.line
		fields[index] = &shared
	}
	if fields == nil {
		return typ
	}
	list := *st.Fields
	list.List = fields
	shared := *st
	shared.Fields = &list
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

func sharedDeclComments(decl ast.Decl, specIndex int) (doc, line *ast.CommentGroup, onSpec bool) {
	switch decl := decl.(type) {
	case *ast.FuncDecl:
		return decl.Doc, nil, false
	case *ast.GenDecl:
		switch spec := decl.Specs[specIndex].(type) {
		case *ast.TypeSpec:
			doc, line = spec.Doc, spec.Comment
		case *ast.ValueSpec:
			doc, line = spec.Doc, spec.Comment
		default:
			return nil, nil, false
		}
		// A doc above const A = 1 belongs to A. A doc above const (...) belongs to the block.
		if decl.Doc != nil && len(decl.Specs) == 1 && !decl.Lparen.IsValid() {
			return decl.Doc, line, false
		}
		return doc, line, true
	}
	return nil, nil, false
}

func declBlockDoc(decl ast.Decl) *ast.CommentGroup {
	if decl, ok := decl.(*ast.GenDecl); ok && (decl.Lparen.IsValid() || len(decl.Specs) > 1) {
		return decl.Doc
	}
	return nil
}

// sortMergeDecls keeps input order when a file lacks positions, as with gogen's ASTFile.
func sortMergeDecls(file *ast.File, decls []ast.Decl) {
	if !file.FileStart.IsValid() {
		return
	}
	valid := func(pos token.Pos) bool {
		return pos >= file.FileStart && (!file.FileEnd.IsValid() || pos < file.FileEnd)
	}
	for _, decl := range file.Decls {
		if gen, ok := decl.(*ast.GenDecl); ok {
			switch gen.Tok {
			case token.IMPORT:
				continue
			case token.TYPE, token.VAR, token.CONST:
				for _, spec := range gen.Specs {
					if !valid(spec.Pos()) {
						return
					}
				}
				continue
			}
		}
		if !valid(decl.Pos()) {
			return
		}
	}
	for _, decl := range decls {
		if gen, ok := decl.(*ast.GenDecl); ok && (gen.Tok == token.TYPE || gen.Tok == token.VAR || gen.Tok == token.CONST) {
			slices.SortStableFunc(gen.Specs, func(a, b ast.Spec) int {
				return cmp.Compare(a.Pos(), b.Pos())
			})
		}
	}
	position := func(decl ast.Decl) token.Pos {
		if gen, ok := decl.(*ast.GenDecl); ok && (gen.Tok == token.TYPE || gen.Tok == token.VAR || gen.Tok == token.CONST) {
			return gen.Specs[0].Pos()
		}
		return decl.Pos()
	}
	slices.SortStableFunc(decls, func(a, b ast.Decl) int {
		return cmp.Compare(position(a), position(b))
	})
}
