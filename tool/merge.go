/*
 * Copyright (c) 2026 The XGo Authors (xgo.dev). All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package tool

import (
	"fmt"
	"go/ast"
	"go/token"
	"sort"
	"strconv"
	"strings"
)

// -----------------------------------------------------------------------------
// Automatic merging of per-platform packages.
//
// Proposal #1125 makes llcppg generate one independent Go package per
// $GOOS-$GOARCH platform. That leaves the consumer with N near-identical
// packages where, in practice, most declarations are the same on every
// platform. Merge combines the per-platform ASTs into one multi-platform Go
// package:
//
//   - Declarations that are identical on every platform are emitted once, in a
//     shared file.
//   - Declarations that differ on any platform (or exist on only some
//     platforms) are emitted per platform, in z_$(GOOS)_$(GOARCH).go files
//     guarded by build constraints.
//
// The merge works on go/ast values, so it is independent of how the platform
// packages were produced. It does not touch the filesystem, run clang, or
// format source code.
//
// See https://github.com/goplus/llcppg/issues/1129.

// PlatformPackage is the converted output of llcppg for one platform.
type PlatformPackage struct {
	Tags  []string             // $GOOS = Tags[0], $GOARCH = Tags[1]
	Files map[string]*ast.File // file name => AST
}

// platformFilePrefix is reserved for generated per-platform files. No input
// file name may start with it.
const platformFilePrefix = "z_"

// importPath is the full import path of a package, for example
// "github.com/goplus/lib/c".
type importPath string

// Merge merges per-platform packages into one multi-platform package.
// The result maps file name => AST.
//
// The input ASTs are treated as read-only. The only change Merge makes to them
// is the package-ref annotation in Ident.Obj.Data (Step 1). The output ASTs are
// new trees: the chosen declarations are copied, and every token.Pos in the
// copies is token.NoPos.
func Merge(pkgs []*PlatformPackage) (map[string]*ast.File, error) {
	if len(pkgs) == 0 {
		return nil, fmt.Errorf("merge: no platform packages given")
	}

	// Validate inputs and establish the common package name and the sorted
	// platform order used throughout the algorithm.
	pkgName, err := validate(pkgs)
	if err != nil {
		return nil, err
	}

	// Sort platforms by their $GOOS-$GOARCH string so the result is independent
	// of the order platforms were passed in.
	order := make([]int, len(pkgs))
	for i := range order {
		order[i] = i
	}
	sort.Slice(order, func(a, b int) bool {
		return platformString(pkgs[order[a]]) < platformString(pkgs[order[b]])
	})
	sorted := make([]*PlatformPackage, len(pkgs))
	for i, idx := range order {
		sorted[i] = pkgs[idx]
	}
	pkgs = sorted

	// Step 1: annotate extern package refs and remember local import names.
	localNames := map[importPath]string{}
	for _, p := range pkgs {
		for _, name := range sortedFileNames(p.Files) {
			annotatePackageRefs(p.Files[name], localNames)
		}
	}

	// Step 2 + 3: collect global symbols with their per-platform placement.
	syms, err := collectSymbols(pkgs)
	if err != nil {
		return nil, err
	}

	// Step 4 + 6: decide shared vs per-platform and assemble the output.
	return assemble(pkgs, pkgName, syms, localNames)
}

// -----------------------------------------------------------------------------
// Validation (Step 0).

func validate(pkgs []*PlatformPackage) (pkgName string, err error) {
	seen := map[string]bool{}
	for _, p := range pkgs {
		if len(p.Tags) != 2 || p.Tags[0] == "" || p.Tags[1] == "" {
			return "", fmt.Errorf("merge: each platform must have exactly two non-empty tags, got %v", p.Tags)
		}
		key := platformString(p)
		if seen[key] {
			return "", fmt.Errorf("merge: duplicate platform %q", key)
		}
		seen[key] = true

		for name, f := range p.Files {
			if strings.HasPrefix(name, platformFilePrefix) {
				return "", fmt.Errorf("merge: input file %q (platform %q) uses the reserved %q prefix",
					name, key, platformFilePrefix)
			}
			if f == nil || f.Name == nil {
				return "", fmt.Errorf("merge: file %q (platform %q) has no package clause", name, key)
			}
			if pkgName == "" {
				pkgName = f.Name.Name
			} else if f.Name.Name != pkgName {
				return "", fmt.Errorf("merge: package name mismatch: %q vs %q", pkgName, f.Name.Name)
			}
		}
	}
	if pkgName == "" {
		return "", fmt.Errorf("merge: no files to merge")
	}
	return pkgName, nil
}

// platformString returns the "$GOOS-$GOARCH" string of a platform package.
func platformString(p *PlatformPackage) string {
	return p.Tags[0] + "-" + p.Tags[1]
}

// -----------------------------------------------------------------------------
// Step 1: Make package references independent of the import name.

// annotatePackageRefs walks a file and marks every extern package ref
// (SelectorExpr whose X names an imported package) by recording the import path
// in X.Obj.Data. Local names are remembered in localNames so imports can be
// regenerated at the end; the first platform (in sorted order) to use a path
// supplies its local name.
func annotatePackageRefs(f *ast.File, localNames map[importPath]string) {
	// Build name => importPath for this file from its import declarations.
	// Dot imports and blank imports never produce package refs.
	imports := map[string]importPath{}
	for _, imp := range f.Imports {
		path, err := strconv.Unquote(imp.Path.Value)
		if err != nil {
			continue
		}
		var local string
		switch {
		case imp.Name == nil:
			local = pkgNameOf(path)
		case imp.Name.Name == "_" || imp.Name.Name == ".":
			continue // blank / dot import: never a package ref
		default:
			local = imp.Name.Name
		}
		imports[local] = importPath(path)
		// Remember the local name for import regeneration (Step 6). Because
		// platforms are processed in sorted order and we only record a path the
		// first time it is seen, the first platform's local name wins.
		if _, ok := localNames[importPath(path)]; !ok {
			localNames[importPath(path)] = local
		}
	}
	if len(imports) == 0 {
		return
	}

	// Track locally-declared names that shadow package names so a parameter or
	// field named like a package is not mistaken for a package ref.
	for _, decl := range f.Decls {
		annotateDecl(decl, imports, newScope(nil))
	}
}

// pkgNameOf derives the default package name from an import path (the last path
// element). This mirrors Go's own default when no explicit alias is given.
func pkgNameOf(path string) string {
	if i := strings.LastIndexByte(path, '/'); i >= 0 {
		return path[i+1:]
	}
	return path
}

// scope tracks names shadowing package names within a nested lexical context.
type scope struct {
	parent *scope
	names  map[string]bool
}

func newScope(parent *scope) *scope { return &scope{parent: parent, names: map[string]bool{}} }

func (s *scope) declare(name string) {
	if name != "" && name != "_" {
		s.names[name] = true
	}
}

func (s *scope) shadows(name string) bool {
	for cur := s; cur != nil; cur = cur.parent {
		if cur.names[name] {
			return true
		}
	}
	return false
}

// annotateDecl annotates package refs inside a single top-level declaration.
func annotateDecl(decl ast.Decl, imports map[string]importPath, sc *scope) {
	switch d := decl.(type) {
	case *ast.FuncDecl:
		fnScope := newScope(sc)
		declareFieldList(fnScope, d.Recv)
		if d.Type != nil {
			declareFieldList(fnScope, d.Type.TypeParams)
		}
		// Annotate the signature types before declaring params so a package ref
		// in a parameter type (func(c.Int)) is still recognized. Parameter and
		// result names shadow only within the body, which we also annotate.
		if d.Type != nil {
			annotateFieldListTypes(d.Type.Params, imports, fnScope)
			annotateFieldListTypes(d.Type.Results, imports, fnScope)
		}
		declareFieldList(fnScope, d.Type.Params)
		declareFieldList(fnScope, d.Type.Results)
		if d.Recv != nil {
			annotateFieldListTypes(d.Recv, imports, fnScope)
		}
		if d.Body != nil {
			annotateExpr(d.Body, imports, newScope(fnScope))
		}
	case *ast.GenDecl:
		for _, spec := range d.Specs {
			annotateSpec(spec, imports, sc)
		}
	}
}

func declareFieldList(sc *scope, fl *ast.FieldList) {
	if fl == nil {
		return
	}
	for _, field := range fl.List {
		for _, name := range field.Names {
			sc.declare(name.Name)
		}
	}
}

func annotateFieldListTypes(fl *ast.FieldList, imports map[string]importPath, sc *scope) {
	if fl == nil {
		return
	}
	for _, field := range fl.List {
		annotateExpr(field.Type, imports, sc)
	}
}

func annotateSpec(spec ast.Spec, imports map[string]importPath, sc *scope) {
	switch s := spec.(type) {
	case *ast.TypeSpec:
		typeScope := sc
		if s.TypeParams != nil {
			typeScope = newScope(sc)
			declareFieldList(typeScope, s.TypeParams)
			annotateFieldListTypes(s.TypeParams, imports, typeScope)
		}
		annotateExpr(s.Type, imports, typeScope)
	case *ast.ValueSpec:
		annotateExpr(s.Type, imports, sc)
		for _, v := range s.Values {
			annotateExpr(v, imports, sc)
		}
	}
}

// annotateExpr recursively walks an expression (or statement), marking extern
// package refs. It maintains lexical scopes so locally-declared names shadow
// package names.
func annotateExpr(node ast.Node, imports map[string]importPath, sc *scope) {
	if node == nil {
		return
	}
	switch n := node.(type) {
	case *ast.SelectorExpr:
		if id, ok := n.X.(*ast.Ident); ok {
			if path, isPkg := imports[id.Name]; isPkg && !sc.shadows(id.Name) {
				markPackageRef(id, path)
				return // X is a package name, not a sub-expression to recurse into
			}
		}
		annotateExpr(n.X, imports, sc)
	case *ast.Ident:
		// nothing to annotate on a bare identifier
	default:
		// Generic recursion for the remaining node kinds. We rely on
		// ast.Inspect for composite nodes but must handle scope for blocks and
		// function literals; those are handled explicitly below before falling
		// back to a structural walk.
		walkChildren(node, imports, sc)
	}
}

// walkChildren handles the container node kinds that matter for merge inputs,
// tracking scopes where new names are introduced.
func walkChildren(node ast.Node, imports map[string]importPath, sc *scope) {
	switch n := node.(type) {
	case *ast.FuncLit:
		litScope := newScope(sc)
		if n.Type != nil {
			annotateFieldListTypes(n.Type.Params, imports, litScope)
			annotateFieldListTypes(n.Type.Results, imports, litScope)
			declareFieldList(litScope, n.Type.Params)
			declareFieldList(litScope, n.Type.Results)
		}
		annotateExpr(n.Body, imports, newScope(litScope))
	case *ast.BlockStmt:
		blockScope := newScope(sc)
		for _, stmt := range n.List {
			annotateStmt(stmt, imports, blockScope)
		}
	case *ast.StructType:
		for _, field := range fields(n.Fields) {
			annotateExpr(field.Type, imports, sc)
		}
	case *ast.InterfaceType:
		for _, field := range fields(n.Methods) {
			annotateExpr(field.Type, imports, sc)
		}
	case *ast.FuncType:
		annotateFieldListTypes(n.Params, imports, sc)
		annotateFieldListTypes(n.Results, imports, sc)
	default:
		ast.Inspect(node, func(child ast.Node) bool {
			if child == nil || child == node {
				return true
			}
			switch child.(type) {
			case *ast.SelectorExpr, *ast.FuncLit, *ast.BlockStmt,
				*ast.StructType, *ast.InterfaceType, *ast.FuncType:
				annotateExpr(child, imports, sc)
				return false
			}
			return true
		})
	}
}

func annotateStmt(stmt ast.Stmt, imports map[string]importPath, sc *scope) {
	switch s := stmt.(type) {
	case *ast.DeclStmt:
		if gen, ok := s.Decl.(*ast.GenDecl); ok {
			for _, spec := range gen.Specs {
				if vs, ok := spec.(*ast.ValueSpec); ok {
					annotateExpr(vs.Type, imports, sc)
					for _, v := range vs.Values {
						annotateExpr(v, imports, sc)
					}
					for _, name := range vs.Names {
						sc.declare(name.Name)
					}
				}
			}
		}
	case *ast.AssignStmt:
		for _, rhs := range s.Rhs {
			annotateExpr(rhs, imports, sc)
		}
		if s.Tok == token.DEFINE {
			for _, lhs := range s.Lhs {
				if id, ok := lhs.(*ast.Ident); ok {
					sc.declare(id.Name)
				}
			}
		}
	case *ast.BlockStmt:
		walkChildren(s, imports, sc)
	default:
		if stmt != nil {
			annotateExpr(stmt, imports, sc)
		}
	}
}

func fields(fl *ast.FieldList) []*ast.Field {
	if fl == nil {
		return nil
	}
	return fl.List
}

// markPackageRef records an extern package ref on the identifier X of X.Sel.
// It is idempotent: re-marking with the same path leaves X.Obj unchanged.
func markPackageRef(x *ast.Ident, path importPath) {
	if x.Obj == nil {
		x.Obj = &ast.Object{Kind: ast.Pkg, Name: x.Name}
	}
	x.Obj.Kind = ast.Pkg
	x.Obj.Data = path
}

// packageRefOf returns the import path recorded on an identifier by Step 1, and
// whether it is an extern package ref.
func packageRefOf(x *ast.Ident) (importPath, bool) {
	if x == nil || x.Obj == nil || x.Obj.Kind != ast.Pkg {
		return "", false
	}
	path, ok := x.Obj.Data.(importPath)
	return path, ok
}
