package merge

import (
	"go/ast"
	"path"
	"strconv"
)

// collectImports maps the package names of a file to their import paths.
// A different package name must be explicit, as in import lib "example.org/lib/v2".
func collectImports(file *ast.File) map[string]importPath {
	imports := make(map[string]importPath, len(file.Imports))
	for _, spec := range file.Imports {
		pkgPath, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			continue
		}
		name := path.Base(pkgPath)
		if spec.Name != nil {
			name = spec.Name.Name
		}
		if name != "." && name != "_" {
			imports[name] = importPath(pkgPath)
		}
	}
	return imports
}

// collectImportRefs records the import path of each package qualifier in node, such as c in c.Int.
// Type parameters in params can hide a package name, as in func F[c any]().
func (p *mergeIndex) collectImportRefs(node ast.Node, params *ast.FieldList, imports map[string]importPath) {
	if len(imports) == 0 {
		return
	}
	ast.PreorderStack(node, nil, func(node ast.Node, stack []ast.Node) bool {
		sel, ok := node.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		qualifier, ok := sel.X.(*ast.Ident)
		if !ok {
			return true
		}
		path, ok := imports[qualifier.Name]
		if !ok || (qualifier.Obj != nil && qualifier.Obj.Kind != ast.Pkg && qualifier.Obj.Kind != ast.Bad) {
			return true
		}
		if _, imported := astIdentName(qualifier); !imported && isFieldKey(sel, stack) {
			return true
		}
		if !declaresName(params, qualifier.Name) {
			p.importRefs[qualifier] = path
		}
		return true
	})
}

// isFieldKey reports whether sel is a key in a struct literal.
// Go 1.27 allows promoted field keys such as Box{Point.X: 1}.
// llcppg does not emit them, so ignore ambiguous keys without a package binding.
func isFieldKey(sel *ast.SelectorExpr, stack []ast.Node) bool {
	if len(stack) < 2 {
		return false
	}
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
	pair, ok := stack[parent].(*ast.KeyValueExpr)
	if !ok || parent == 0 || pair.Key != key {
		return false
	}
	lit, ok := stack[parent-1].(*ast.CompositeLit)
	if !ok {
		return false
	}
	switch ast.Unparen(lit.Type).(type) {
	case *ast.ArrayType, *ast.MapType:
		return false
	}
	return true
}

func declaresName(params *ast.FieldList, name string) bool {
	if params == nil {
		return false
	}
	for _, field := range params.List {
		for _, ident := range field.Names {
			if ident.Name == name {
				return true
			}
		}
	}
	return false
}
