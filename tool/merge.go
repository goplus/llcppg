package tool

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

// Merge merges platform packages into a map from output file name to AST.
// For example, Go declarations common to linux/amd64 and darwin/arm64 share a file.
func Merge(pkgs []*PlatformPackage) (map[string]*ast.File, error) {
	for i, pkg := range pkgs {
		if pkg == nil || len(pkg.Tags) != 2 || pkg.Tags[0] == "" || pkg.Tags[1] == "" {
			return nil, fmt.Errorf("platform %d: expected GOOS and GOARCH tags", i)
		}
	}
	ordered := slices.Clone(pkgs)
	slices.SortFunc(ordered, func(a, b *PlatformPackage) int {
		return cmp.Compare(a.Tags[0]+"-"+a.Tags[1], b.Tags[0]+"-"+b.Tags[1])
	})

	var pkgName string
	tags := make([][]string, len(ordered))
	for i, pkg := range ordered {
		platform := pkg.Tags[0] + "-" + pkg.Tags[1]
		if i > 0 && slices.Equal(tags[i-1], pkg.Tags) {
			return nil, fmt.Errorf("duplicate platform %q", platform)
		}
		tags[i] = pkg.Tags
		for name, file := range pkg.Files {
			if strings.HasPrefix(name, "z_") {
				return nil, fmt.Errorf("platform %s: reserved file name %q", platform, name)
			}
			if file == nil || file.Name == nil || file.Name.Name == "" {
				return nil, fmt.Errorf("platform %s: missing package name in %q", platform, name)
			}
			if pkgName == "" {
				pkgName = file.Name.Name
			}
			if file.Name.Name != pkgName {
				return nil, fmt.Errorf("platform %s: package %q in %q does not match %q", platform, file.Name.Name, name, pkgName)
			}
		}
	}

	var index mergeIndex
	index.new()
	for i, pkg := range ordered {
		if err := index.collectPlatform(i, pkg.Files); err != nil {
			return nil, err
		}
	}
	return index.buildFiles(pkgName, tags), nil
}
