package gowrite

import (
	"bytes"
	"errors"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/goplus/gogen"
)

func TestWriteTo_EmptyFuncHasCompactBody(t *testing.T) {
	pkg := gogen.NewPackage("", "demo", nil)
	pkg.NewFunc(nil, "InitHooks", nil, nil, false).BodyStart(pkg).End()

	var buf bytes.Buffer
	if err := WriteTo(&buf, pkg, ""); err != nil {
		t.Fatalf("WriteTo failed: %v", err)
	}
	got := buf.String()
	want := (`package demo

func InitHooks() {
}
`)
	if got != want {
		t.Fatalf("unexpected output.\nwant:\n%s\ngot:\n%s", want, got)
	}
}

func TestWriteTo_NonEmptyFuncStaysMultiline(t *testing.T) {
	pkg := gogen.NewPackage("", "demo", nil)
	results := types.NewTuple(pkg.NewParam(token.NoPos, "", types.Typ[types.Int], false))
	pkg.NewFunc(nil, "RetZero", nil, results, false).BodyStart(pkg).Val(0).Return(1).End()

	var buf bytes.Buffer
	if err := WriteTo(&buf, pkg, ""); err != nil {
		t.Fatalf("WriteTo failed: %v", err)
	}
	got := buf.String()
	want := (`package demo

func RetZero() int {
	return 0
}
`)
	if got != want {
		t.Fatalf("unexpected output.\nwant:\n%s\ngot:\n%s", want, got)
	}
}

func TestAssignDeclAnchors_EmptyInterfaceStaysCompact(t *testing.T) {
	file := &ast.File{
		Name: ast.NewIdent("demo"),
		Decls: []ast.Decl{
			&ast.FuncDecl{
				Name: ast.NewIdent("Mprintf"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{ast.NewIdent("format")},
								Type:  ast.NewIdent("string"),
							},
							{
								Names: []*ast.Ident{ast.NewIdent("__llgo_va_list")},
								Type: &ast.Ellipsis{
									Elt: &ast.InterfaceType{Methods: &ast.FieldList{}},
								},
							},
						},
					},
				},
			},
		},
	}

	fset := token.NewFileSet()
	anchorDecls(fset, "demo.go", file)

	var buf bytes.Buffer
	if err := format.Node(&buf, fset, file); err != nil {
		t.Fatalf("format.Node failed: %v", err)
	}

	got := buf.String()
	want := (`package demo

func Mprintf(format string, __llgo_va_list ...interface{})
`)
	if got != want {
		t.Fatalf("unexpected output.\nwant:\n%s\ngot:\n%s", want, got)
	}
}

func TestAssignDeclAnchors_MultiLineDocAndLinkStaySeparated(t *testing.T) {
	file := &ast.File{
		Name: ast.NewIdent("demo"),
		Decls: []ast.Decl{
			&ast.FuncDecl{
				Doc: &ast.CommentGroup{
					List: []*ast.Comment{
						{Text: "/*\nExecuteFoo comment\n*/"},
						{Text: "//go:linkname CustomExecuteFoo2 C.ExecuteFoo2"},
					},
				},
				Name: ast.NewIdent("CustomExecuteFoo2"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
				},
			},
		},
	}

	fset := token.NewFileSet()
	anchorDecls(fset, "demo.go", file)

	var buf bytes.Buffer
	if err := format.Node(&buf, fset, file); err != nil {
		t.Fatalf("format.Node failed: %v", err)
	}

	got := buf.String()
	want := (`package demo

/*
ExecuteFoo comment
*/
//go:linkname CustomExecuteFoo2 C.ExecuteFoo2
func CustomExecuteFoo2()
`)
	if got != want {
		t.Fatalf("unexpected output.\nwant:\n%s\ngot:\n%s", want, got)
	}
}

func TestAssignDeclAnchors_FuncTypeDeclEllipsisInterfaceStaysCompact(t *testing.T) {
	file := &ast.File{
		Name: ast.NewIdent("demo"),
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: ast.NewIdent("Fn"),
						Type: &ast.FuncType{
							Params: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{ast.NewIdent("__llgo_va_list")},
										Type: &ast.Ellipsis{
											Elt: &ast.InterfaceType{Methods: &ast.FieldList{}},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	fset := token.NewFileSet()
	anchorDecls(fset, "demo.go", file)

	var buf bytes.Buffer
	if err := format.Node(&buf, fset, file); err != nil {
		t.Fatalf("format.Node failed: %v", err)
	}

	got := buf.String()
	want := (`package demo

type Fn func(__llgo_va_list ...interface{})
`)
	if got != want {
		t.Fatalf("unexpected output.\nwant:\n%s\ngot:\n%s", want, got)
	}
}

func TestWriteFile_WritesHeaderAndContent(t *testing.T) {
	pkg := gogen.NewPackage("", "demo", nil)
	pkg.NewFunc(nil, "InitHooks", nil, nil, false).BodyStart(pkg).End()

	oldHeader := gogen.GeneratedHeader
	gogen.GeneratedHeader = "// generated in test\n"
	t.Cleanup(func() { gogen.GeneratedHeader = oldHeader })

	outFile := filepath.Join(t.TempDir(), "demo.go")
	if err := WriteFile(pkg, outFile, ""); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	got := string(data)
	if !strings.HasPrefix(got, "// generated in test\n") {
		t.Fatalf("missing generated header in output:\n%s", got)
	}
	if !strings.Contains(got, "func InitHooks() {\n}\n") {
		t.Fatalf("missing function body in output:\n%s", got)
	}
}

func TestWriteTo_MissingNamedFileReturnsENOENT(t *testing.T) {
	pkg := gogen.NewPackage("", "demo", nil)
	var buf bytes.Buffer
	err := WriteTo(&buf, pkg, "missing.go")
	if !errors.Is(err, syscall.ENOENT) {
		t.Fatalf("want ENOENT, got: %v", err)
	}
}

func TestWriteFile_MissingNamedFileReturnsENOENT(t *testing.T) {
	pkg := gogen.NewPackage("", "demo", nil)
	outFile := filepath.Join(t.TempDir(), "out.go")
	err := WriteFile(pkg, outFile, "missing.go")
	if !errors.Is(err, syscall.ENOENT) {
		t.Fatalf("want ENOENT, got: %v", err)
	}
}

func TestAstFile_UsesBaseName(t *testing.T) {
	pkg := gogen.NewPackage("", "demo", nil)
	if _, err := pkg.SetCurFile("nested/path/custom.go", true); err != nil {
		t.Fatalf("SetCurFile failed: %v", err)
	}
	pkg.NewFunc(nil, "InitHooks", nil, nil, false).BodyStart(pkg).End()

	file, name, err := astFile(pkg, "nested/path/custom.go")
	if err != nil {
		t.Fatalf("astFile failed: %v", err)
	}
	if file == nil {
		t.Fatal("astFile returned nil file")
	}
	if name != "custom.go" {
		t.Fatalf("unexpected name: %q", name)
	}
}

func TestAnchorComments_SkipsNilCommentEntry(t *testing.T) {
	group := &ast.CommentGroup{
		List: []*ast.Comment{
			nil,
			{Text: "// doc"},
		},
	}

	next := 1
	newPos := func() token.Pos {
		p := token.Pos(next)
		next++
		return p
	}
	var skipped int
	anchorComments(group, newPos, func(n int) { skipped += n })

	if got := group.List[1].Slash; got != token.Pos(1) {
		t.Fatalf("unexpected slash pos: %v", got)
	}
	if skipped != 0 {
		t.Fatalf("unexpected skipped lines: %d", skipped)
	}
}

func TestNeedsEmptyIface_Cases(t *testing.T) {
	if needsEmptyIface(&ast.InterfaceType{Interface: token.Pos(1)}) {
		t.Fatal("interface with existing pos should not need anchor")
	}
	if needsEmptyIface(&ast.InterfaceType{}) {
		t.Fatal("nil methods should not need anchor")
	}
	if needsEmptyIface(&ast.InterfaceType{
		Methods: &ast.FieldList{List: []*ast.Field{{Type: ast.NewIdent("any")}}},
	}) {
		t.Fatal("non-empty methods should not need anchor")
	}
	if !needsEmptyIface(&ast.InterfaceType{
		Methods: &ast.FieldList{Opening: token.NoPos, Closing: token.NoPos},
	}) {
		t.Fatal("empty methods without braces pos should need anchor")
	}
	if needsEmptyIface(&ast.InterfaceType{
		Methods: &ast.FieldList{Opening: token.Pos(1), Closing: token.Pos(1)},
	}) {
		t.Fatal("empty methods with braces pos should not need anchor")
	}
}

func TestAnchorEmptyIfaceInFunc_SkipsAlreadyAnchored(t *testing.T) {
	ft := &ast.FuncType{
		Params: &ast.FieldList{
			List: []*ast.Field{
				{
					Type: &ast.Ellipsis{
						Elt: &ast.InterfaceType{
							Interface: token.Pos(99),
							Methods:   &ast.FieldList{},
						},
					},
				},
			},
		},
	}

	next := 1
	newPos := func() token.Pos {
		p := token.Pos(next)
		next++
		return p
	}
	anchorEmptyIfaceInFunc(ft, newPos)
	if next != 1 {
		t.Fatalf("newPos should not be consumed, got next=%d", next)
	}
}

func TestAnchorEmptyIfaceInTypeFuncs_SkipsNonFunctionTypeSpecs(t *testing.T) {
	fnIface := &ast.InterfaceType{Methods: &ast.FieldList{}}
	d := &ast.GenDecl{
		Specs: []ast.Spec{
			&ast.ValueSpec{Names: []*ast.Ident{ast.NewIdent("v")}},
			&ast.TypeSpec{Name: ast.NewIdent("S"), Type: &ast.StructType{}},
			&ast.TypeSpec{
				Name: ast.NewIdent("Fn"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{Type: &ast.Ellipsis{Elt: fnIface}},
						},
					},
				},
			},
		},
	}

	next := 1
	newPos := func() token.Pos {
		p := token.Pos(next)
		next++
		return p
	}
	anchorEmptyIfaceInTypeFuncs(d, newPos)

	if fnIface.Interface == token.NoPos {
		t.Fatal("function type interface should be anchored")
	}
}

func TestLineSpan_EmptyString(t *testing.T) {
	if got := lineSpan(""); got != 1 {
		t.Fatalf("unexpected line span: %d", got)
	}
}
