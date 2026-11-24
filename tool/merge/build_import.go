package merge

import (
	"go/ast"
	"go/token"
	"maps"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

// buildImports adds the imports used by file, with one name per import path.
func (p *mergeIndex) buildImports(file *ast.File) {
	names := make(map[importPath]string)
	replacements := make(map[ast.Node]ast.Node)
	linkname := false
	ast.PreorderStack(file, nil, func(node ast.Node, stack []ast.Node) bool {
		switch node := node.(type) {
		case *ast.Ident:
			pkgPath, ok := p.importRefs[node]
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
			for _, directive := range parseComments(node).directives {
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
