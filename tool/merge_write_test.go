package tool

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/qiniu/x/test"
)

func writeMergedFile(w io.Writer, fset *token.FileSet, input *ast.File) error {
	var source bytes.Buffer
	if input.Doc != nil {
		for _, comment := range input.Doc.List {
			fmt.Fprintln(&source, comment.Text)
		}
		source.WriteByte('\n')
	}
	fmt.Fprintln(&source, "package", input.Name.Name)

	// Print each declaration with its source positions so //go:linkname Dlopen C.dlopen stays with Dlopen.
	for _, decl := range input.Decls {
		source.WriteByte('\n')
		if err := format.Node(&source, fset, decl); err != nil {
			return err
		}
		source.WriteByte('\n')
	}
	formatted, err := format.Source(source.Bytes())
	if err != nil {
		return err
	}
	_, err = w.Write(formatted)
	return err
}

func TestMergeGolden(t *testing.T) {
	const root = "_testmerge"
	cases, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, fixture := range cases {
		t.Run(fixture.Name(), func(t *testing.T) {
			dir := filepath.Join(root, fixture.Name())
			inDir := filepath.Join(dir, "in")
			platforms, err := os.ReadDir(inDir)
			if err != nil {
				t.Fatal(err)
			}
			fset := token.NewFileSet()
			pkgs := make([]*PlatformPackage, 0, len(platforms))
			// Supply linux before darwin so the golden checks Merge's platform ordering.
			for _, platform := range slices.Backward(platforms) {
				platformDir := filepath.Join(inDir, platform.Name())
				paths, err := filepath.Glob(filepath.Join(platformDir, "*.go"))
				if err != nil {
					t.Fatal(err)
				}
				if len(paths) == 0 {
					t.Fatalf("no Go inputs in %s", platformDir)
				}
				for i, path := range paths {
					paths[i] = filepath.Base(path)
				}
				files, err := LoadGoFiles(fset, platformDir, paths, parser.ParseComments)
				if err != nil {
					t.Fatal(err)
				}
				pkgs = append(pkgs, &PlatformPackage{Tags: strings.Split(platform.Name(), "-"), Files: files})
			}
			merged, err := Merge(pkgs)
			if err != nil {
				t.Fatal(err)
			}

			outDir := filepath.Join(dir, "out")
			for _, name := range slices.Sorted(maps.Keys(merged)) {
				golden := filepath.Join(outDir, strings.TrimSuffix(name, ".go")+".merge.go")
				want, err := os.ReadFile(golden)
				if err != nil && !os.IsNotExist(err) {
					t.Fatal(err)
				}
				var got bytes.Buffer
				if err := writeMergedFile(&got, fset, merged[name]); err != nil {
					t.Fatalf("write %s: %v", name, err)
				}
				test.Diff(t, golden+".txt", got.Bytes(), want)
			}
			expected, err := filepath.Glob(filepath.Join(outDir, "*.merge.go"))
			if err != nil {
				t.Fatal(err)
			}
			for _, path := range expected {
				name := strings.TrimSuffix(filepath.Base(path), ".merge.go") + ".go"
				if merged[name] == nil {
					t.Errorf("missing merged output %s", name)
				}
			}
		})
	}
}
