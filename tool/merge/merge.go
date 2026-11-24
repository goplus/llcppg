package merge

import (
	"cmp"
	"fmt"
	"go/ast"
	"slices"
	"strings"
)

// PlatformPackage is the converted output of llcppg for one platform.
type PlatformPackage struct {
	Tags  []string             // GOOS and GOARCH, such as linux and amd64.
	Files map[string]*ast.File // File name to AST.
}

func (p *PlatformPackage) platform() string {
	return p.Tags[0] + "-" + p.Tags[1]
}

// Merge merges platform packages into a map from output file name to AST.
// For example, Go declarations common to linux/amd64 and darwin/arm64 share a file.
func Merge(pkgs []*PlatformPackage) (map[string]*ast.File, error) {
	pkgs, err := sortPackages(pkgs)
	if err != nil {
		return nil, err
	}
	pkgName, err := checkFiles(pkgs)
	if err != nil {
		return nil, err
	}
	index := newMergeIndex(pkgs)
	index.collect()
	return index.build(pkgName), nil
}

// sortPackages returns pkgs in platform order, such as darwin-arm64 before linux-amd64.
func sortPackages(pkgs []*PlatformPackage) ([]*PlatformPackage, error) {
	for i, pkg := range pkgs {
		if pkg == nil || len(pkg.Tags) != 2 || pkg.Tags[0] == "" || pkg.Tags[1] == "" {
			return nil, fmt.Errorf("platform %d: expected GOOS and GOARCH tags", i)
		}
	}
	sorted := slices.Clone(pkgs)
	slices.SortFunc(sorted, func(a, b *PlatformPackage) int {
		return cmp.Compare(a.platform(), b.platform())
	})
	for i := 1; i < len(sorted); i++ {
		if slices.Equal(sorted[i-1].Tags, sorted[i].Tags) {
			return nil, fmt.Errorf("duplicate platform %q", sorted[i].platform())
		}
	}
	return sorted, nil
}

// checkFiles returns the package name shared by all files.
func checkFiles(pkgs []*PlatformPackage) (pkgName string, err error) {
	for _, pkg := range pkgs {
		platform := pkg.platform()
		for name, file := range pkg.Files {
			if strings.HasPrefix(name, "z_") {
				return "", fmt.Errorf("platform %s: reserved file name %q", platform, name)
			}
			if file == nil || file.Name == nil || file.Name.Name == "" {
				return "", fmt.Errorf("platform %s: missing package name in %q", platform, name)
			}
			if pkgName == "" {
				pkgName = file.Name.Name
			}
			if file.Name.Name != pkgName {
				return "", fmt.Errorf("platform %s: package %q in %q does not match %q", platform, file.Name.Name, name, pkgName)
			}
		}
		if err := checkDecls(pkg.Files); err != nil {
			return "", fmt.Errorf("platform %s: %w", platform, err)
		}
	}
	return pkgName, nil
}

// checkDecls rejects declarations that cannot be matched by name across platforms.
// Names share one namespace, so a type in a.go and a function in b.go cannot both declare Open.
func checkDecls(files map[string]*ast.File) error {
	seen := make(map[string]string) // declared name => file name
	add := func(name, file string) error {
		if name == "_" {
			return fmt.Errorf("blank declaration in %q", file)
		}
		if previous, ok := seen[name]; ok {
			return fmt.Errorf("duplicate symbol %q in %q and %q", name, previous, file)
		}
		seen[name] = file
		return nil
	}
	for fname, file := range files {
		for _, decl := range file.Decls {
			switch decl := decl.(type) {
			case *ast.GenDecl:
				for _, spec := range decl.Specs {
					switch spec := spec.(type) {
					case *ast.TypeSpec:
						if err := add(spec.Name.Name, fname); err != nil {
							return err
						}
					case *ast.ValueSpec:
						for _, name := range spec.Names {
							if err := add(name.Name, fname); err != nil {
								return err
							}
						}
					}
				}
			case *ast.FuncDecl:
				key := funcKey(decl)
				if key == "init" {
					return fmt.Errorf("init function in %q", fname)
				}
				if decl.Name.Name == "_" {
					key = "_"
				}
				if err := add(key, fname); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
