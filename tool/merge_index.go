package tool

import (
	"fmt"
	"go/ast"
	"go/token"
	"path"
	"slices"
	"strconv"
)

// Specs can share a declaration index, as in type (A int; B int).
type mergeDeclRef struct {
	platform  int
	file      string
	declIndex int
	specIndex int
}

type mergeSymbol struct {
	decl      mergeDeclRef
	receiver  ast.Expr
	signature *ast.FuncType
	aliases   []mergeDeclRef
}

type mergeIndex struct {
	files         map[int]map[string]*ast.File
	types         astMap
	symbols       map[string][]mergeSymbol
	values        map[string][]mergeDeclRef
	importRefs    map[*ast.Ident]importPath
	importAliases map[importPath][]string
}

func (m *mergeIndex) new() {
	if m.files == nil {
		m.files = make(map[int]map[string]*ast.File)
	}
	if m.symbols == nil {
		m.symbols = make(map[string][]mergeSymbol)
	}
	if m.values == nil {
		m.values = make(map[string][]mergeDeclRef)
	}
	if m.importRefs == nil {
		refs := make(map[*ast.Ident]importPath)
		m.importRefs = refs
		m.types.hasher.identName = func(ident *ast.Ident) (string, bool) {
			if path, ok := refs[ident]; ok {
				return string(path), true
			}
			return astIdentName(ident)
		}
	}
}

// Register import references before lookup so c.Int and libc.Int can share a bucket.
func (m *mergeIndex) collectPlatform(platform int, files map[string]*ast.File) error {
	m.files[platform] = files
	// A type in a.go and a function in b.go can both declare Open.
	seen := make(map[string]string)
	var collectErr error
	registerName := func(symbol, file string) bool {
		if symbol == "_" {
			collectErr = fmt.Errorf("platform %d: blank declaration in %q", platform, file)
			return false
		}
		if previous, ok := seen[symbol]; ok {
			collectErr = fmt.Errorf("platform %d: duplicate symbol %q in %q and %q", platform, symbol, previous, file)
			return false
		}
		seen[symbol] = file
		return true
	}
	for name, file := range files {
		aliases := make(map[string]importPath)
		for _, spec := range file.Imports {
			pkgPath, err := strconv.Unquote(spec.Path.Value)
			if err != nil {
				continue
			}
			// A different package name must be explicit, as in import lib "example.org/lib/v2".
			alias := path.Base(pkgPath)
			if spec.Name != nil {
				alias = spec.Name.Name
			}
			if alias != "." && alias != "_" {
				aliases[alias] = importPath(pkgPath)
			}
		}
		declIndex := 0
		ast.Inspect(file, func(node ast.Node) bool {
			if collectErr != nil {
				return false
			}
			if node == file {
				return true
			}
			decl, ok := node.(ast.Decl)
			if !ok {
				return false
			}
			ref := mergeDeclRef{platform: platform, file: name, declIndex: declIndex, specIndex: -1}
			declIndex++
			switch decl := decl.(type) {
			case *ast.GenDecl:
				switch decl.Tok {
				case token.TYPE:
					for i, spec := range decl.Specs {
						ref.specIndex = i
						typ := spec.(*ast.TypeSpec)
						if !registerName(typ.Name.Name, name) {
							return false
						}
						if len(aliases) != 0 {
							m.collectImportRefs(typ, typ.TypeParams, aliases)
						}
						if refs, ok := m.types.at(typ).(*[]mergeDeclRef); ok {
							if slices.Equal(m.genDeclDirectives((*refs)[0]), m.genDeclDirectives(ref)) {
								*refs = append(*refs, ref)
							}
						} else {
							m.types.set(typ, &[]mergeDeclRef{ref})
						}
					}
				case token.VAR, token.CONST:
					for i, spec := range decl.Specs {
						ref.specIndex = i
						value := spec.(*ast.ValueSpec)
						if len(aliases) != 0 {
							m.collectImportRefs(value, nil, aliases)
						}
						// Keep the shared initializer in the source spec, as in var A, B = pair().
						for _, ident := range value.Names {
							if !registerName(ident.Name, name) {
								return false
							}
							m.collectValue(ident.Name, decl, ref)
						}
					}
				}
			case *ast.FuncDecl:
				key := mergeFuncKey(decl)
				if key == "init" {
					collectErr = fmt.Errorf("platform %d: init function in %q", platform, name)
					return false
				}
				if decl.Name.Name == "_" {
					key = "_"
				}
				if !registerName(key, name) {
					return false
				}
				if len(aliases) != 0 {
					m.collectImportRefs(decl, decl.Type.TypeParams, aliases)
				}
				m.collectSymbol(decl, ref)
			}
			return false
		})
		if collectErr != nil {
			return collectErr
		}
	}
	return nil
}

func (m *mergeIndex) collectValue(name string, decl *ast.GenDecl, ref mergeDeclRef) {
	refs := m.values[name]
	if len(refs) != 0 {
		first := refs[0]
		previous := m.files[first.platform][first.file].Decls[first.declIndex].(*ast.GenDecl)
		if previous.Tok != decl.Tok ||
			!m.types.equal(previous.Specs[first.specIndex], decl.Specs[ref.specIndex]) ||
			!slices.Equal(m.genDeclDirectives(first), m.genDeclDirectives(ref)) {
			return
		}
	}
	m.values[name] = append(refs, ref)
}

func (m *mergeIndex) genDeclDirectives(ref mergeDeclRef) []string {
	decl := m.files[ref.platform][ref.file].Decls[ref.declIndex].(*ast.GenDecl)
	var comments commentgroup
	if ref.specIndex < 0 {
		comments.load(decl.Doc, nil)
		return comments.directives
	}
	var directives []string
	if len(decl.Specs) == 1 {
		// A directive above type Callback func() is attached to the GenDecl.
		comments.load(decl.Doc, nil)
		directives = comments.directives
	}
	var doc *ast.CommentGroup
	switch spec := decl.Specs[ref.specIndex].(type) {
	case *ast.TypeSpec:
		doc = spec.Doc
	case *ast.ValueSpec:
		doc = spec.Doc
	}
	comments.load(doc, nil)
	return append(directives, comments.directives...)
}

func (m *mergeIndex) collectSymbol(decl *ast.FuncDecl, ref mergeDeclRef) {
	key := mergeFuncKey(decl)
	signature := *decl.Type
	signature.Params = mergeSignatureFields(signature.Params)
	signature.Results = mergeSignatureFields(signature.Results)
	symbol := mergeSymbol{decl: ref, signature: &signature}
	if decl.Recv != nil {
		symbol.receiver = decl.Recv.List[0].Type
	}
	variants := m.symbols[key]
	for i := range variants {
		variant := &variants[i]
		if !m.types.equal(variant.receiver, symbol.receiver) || !m.types.equal(variant.signature, symbol.signature) {
			continue
		}
		first := variant.decl
		previous := m.files[first.platform][first.file].Decls[first.declIndex].(*ast.FuncDecl)
		var previousComments, comments commentgroup
		previousComments.load(previous.Doc, nil)
		comments.load(decl.Doc, nil)
		if !previousComments.equalDirectives(comments) {
			continue
		}
		variant.aliases = append(variant.aliases, ref)
		return
	}
	m.symbols[key] = append(variants, symbol)
}

// Expand grouped fields so (a, b int) and (int, int) compare by type.
func mergeSignatureFields(fields *ast.FieldList) *ast.FieldList {
	normalized := &ast.FieldList{}
	if count := fields.NumFields(); count != 0 {
		normalized.List = make([]*ast.Field, 0, count)
		for _, field := range fields.List {
			typ := &ast.Field{Type: field.Type}
			for range max(1, len(field.Names)) {
				normalized.List = append(normalized.List, typ)
			}
		}
	}
	return normalized
}

func (m *mergeIndex) collectImportRefs(decl ast.Node, params *ast.FieldList, aliases map[string]importPath) {
	ast.PreorderStack(decl, nil, func(node ast.Node, stack []ast.Node) bool {
		sel, ok := node.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		qualifier, ok := sel.X.(*ast.Ident)
		if !ok {
			return true
		}
		path, ok := aliases[qualifier.Name]
		if !ok || (qualifier.Obj != nil && qualifier.Obj.Kind != ast.Pkg && qualifier.Obj.Kind != ast.Bad) {
			return true
		}
		// Go 1.27 allows promoted field keys such as Box{Point.X: 1}.
		// llcppg does not emit them, so ignore ambiguous keys without a package binding.
		if _, imported := astIdentName(qualifier); !imported && len(stack) >= 2 {
			var key ast.Expr = sel
			parent := len(stack) - 1
			for parent > 0 {
				outer, ok := stack[parent].(*ast.SelectorExpr)
				if !ok || outer.X != key {
					break
				}
				key = outer
				parent--
			}
			if pair, ok := stack[parent].(*ast.KeyValueExpr); ok && parent > 0 && pair.Key == key {
				if lit, ok := stack[parent-1].(*ast.CompositeLit); ok {
					switch ast.Unparen(lit.Type).(type) {
					case *ast.ArrayType, *ast.MapType:
					default:
						return true
					}
				}
			}
		}
		if params != nil {
			for _, field := range params.List {
				for _, name := range field.Names {
					if name.Name == qualifier.Name {
						return true
					}
				}
			}
		}
		m.importRefs[qualifier] = path
		if m.importAliases == nil {
			m.importAliases = make(map[importPath][]string)
		}
		names := m.importAliases[path]
		if !slices.Contains(names, qualifier.Name) {
			m.importAliases[path] = append(names, qualifier.Name)
		}
		return true
	})
}

func mergeFuncKey(decl *ast.FuncDecl) string {
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
