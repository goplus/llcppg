package merge

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

	"github.com/goplus/llcppg/tool"
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
				files, err := tool.LoadGoFiles(fset, platformDir, paths, parser.ParseComments)
				if err != nil {
					t.Fatal(err)
				}
				pkgs = append(pkgs, &PlatformPackage{Tags: strings.Split(platform.Name(), "-"), Files: files})
			}
			outDir := filepath.Join(dir, "out")
			merged, mergeErr := Merge(pkgs)
			if expected, err := os.ReadFile(filepath.Join(outDir, "error")); err == nil {
				// Duplicate diagnostics may list a.go and b.go in either order.
				want := strings.TrimSpace(string(expected))
				if want == "" || mergeErr == nil || !strings.Contains(mergeErr.Error(), want) {
					t.Fatalf("Merge error = %v, want an error containing %q", mergeErr, want)
				}
				return
			} else if !os.IsNotExist(err) {
				t.Fatal(err)
			}
			if mergeErr != nil {
				t.Fatal(mergeErr)
			}

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

func TestMergeInvalidInput(t *testing.T) {
	linux := []string{"linux", "amd64"}
	tests := []struct {
		name      string
		packages  []*PlatformPackage
		wantError string
	}{
		{"duplicate_platform", []*PlatformPackage{{Tags: linux}, {Tags: linux}}, `duplicate platform "linux-amd64"`},
		{"missing_package_name", []*PlatformPackage{{
			Tags:  linux,
			Files: map[string]*ast.File{"input.go": {}},
		}}, `missing package name in "input.go"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Merge(tt.packages)
			if err == nil || !strings.Contains(err.Error(), tt.wantError) {
				t.Fatalf("Merge error = %v, want an error containing %q", err, tt.wantError)
			}
		})
	}
}

func TestMergeKeepsOrderWithoutFileStart(t *testing.T) {
	// gogen's ASTFile omits FileStart, so a positioned First stays before an unpositioned Second.
	file := &ast.File{
		Name: ast.NewIdent("example"),
		Decls: []ast.Decl{
			&ast.FuncDecl{Name: ast.NewIdent("First"), Type: &ast.FuncType{Func: 10, Params: &ast.FieldList{}}},
			&ast.FuncDecl{Name: ast.NewIdent("Second"), Type: &ast.FuncType{Params: &ast.FieldList{}}},
		},
	}
	merged, err := Merge([]*PlatformPackage{{
		Tags:  []string{"linux", "amd64"},
		Files: map[string]*ast.File{"input.go": file},
	}})
	if err != nil {
		t.Fatal(err)
	}
	out := merged["input.go"]
	if out == nil {
		t.Fatal("missing merged input.go")
	}
	var names []string
	for _, decl := range out.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			t.Fatalf("unexpected declaration %T", decl)
		}
		names = append(names, fn.Name.Name)
	}
	if want := []string{"First", "Second"}; !slices.Equal(names, want) {
		t.Fatalf("declaration order = %v, want %v", names, want)
	}
}
