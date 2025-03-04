package convert_test

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/goplus/llcppg/_xtool/llcppsymg/args"
	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/cmd/gogensig/config"
	"github.com/goplus/llcppg/cmd/gogensig/convert"
	"github.com/goplus/llcppg/cmd/gogensig/dbg"
	"github.com/goplus/llcppg/cmd/gogensig/unmarshal"
	"github.com/goplus/llcppg/llcppg"
	ctoken "github.com/goplus/llcppg/token"
	"github.com/goplus/llgo/xtool/env"
)

func init() {
	dbg.SetDebugAll()
}

func TestFromTestdata(t *testing.T) {
	testFromDir(t, "./_testdata", false)
}

func TestSysToPkg(t *testing.T) {
	name := "_systopkg"
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal("Getwd failed:", err)
	}
	testFrom(t, name, path.Join(dir, "_testdata", name), false, func(t *testing.T, pkg *convert.Package, cvt *convert.Converter) {
		// check FileMap's info is right
		inFileMap := func(file string) {
			_, ok := cvt.Pkg.FileMap[file]
			if !ok {
				t.Fatal("File not found in FileMap:", file)
			}
		}
		for _, decl := range cvt.Pkg.File.Decls {
			switch decl := decl.(type) {
			case *ast.TypeDecl:
				inFileMap(decl.DeclBase.Loc.File)
			case *ast.EnumTypeDecl:
				inFileMap(decl.DeclBase.Loc.File)
			case *ast.TypedefDecl:
				inFileMap(decl.DeclBase.Loc.File)
			case *ast.FuncDecl:
				inFileMap(decl.DeclBase.Loc.File)
			}
		}
	})
}

func TestDepPkg(t *testing.T) {
	name := "_depcjson"
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal("Getwd failed:", err)
	}
	const hfileDir = "hfile"
	testDataDir := "testdata"

	buildTestPath := func(components ...string) string {
		return path.Join(components...)
	}
	inc := func(cfgPath string, incFlag string) func() {
		origContent, err := os.ReadFile(cfgPath)
		_ = origContent
		if err != nil {
			t.Fatal(err)
		}

		cfg, err := config.GetCppgCfgFromPath(cfgPath)
		if err != nil {
			t.Fatal(err)
		}
		cfg.CFlags = cfg.CFlags + incFlag
		err = config.CreateJSONFile(cfgPath, cfg)
		if err != nil {
			t.Fatal(err)
		}

		return func() {
			if err := os.WriteFile(cfgPath, origContent, 0644); err != nil {
				t.Fatal("Failed to restore original config:", err)
			}
		}
	}

	depcjson := path.Join(dir, "_testdata", name)
	depcjsonConf := path.Join(depcjson, "conf", "llcppg.cfg")

	thirdDepPath := buildTestPath(dir, testDataDir, "thirddep")
	thirdDep2Path := buildTestPath(dir, testDataDir, "thirddep2")
	thirdDep3Path := buildTestPath(dir, testDataDir, "thirddep3")
	basicDepPath := buildTestPath(dir, testDataDir, "basicdep")

	thirdDepHFile := buildTestPath(thirdDepPath, hfileDir)
	thirdDep2HFile := buildTestPath(thirdDep2Path, hfileDir)
	thirdDep3HFile := buildTestPath(thirdDep3Path, hfileDir)
	basicDepHFile := buildTestPath(basicDepPath, hfileDir)

	cleanups := []func(){
		inc(depcjsonConf, fmt.Sprintf(" -I%s -I%s -I%s -I%s", thirdDepHFile, thirdDep2HFile, thirdDep3HFile, basicDepHFile)),
		inc(buildTestPath(thirdDepPath, "llcppg.cfg"), fmt.Sprintf(" -I%s -I%s", thirdDepHFile, basicDepHFile)),
		inc(buildTestPath(thirdDep2Path, "llcppg.cfg"), fmt.Sprintf(" -I%s -I%s -I%s", thirdDep2HFile, thirdDepHFile, basicDepHFile)),
		inc(buildTestPath(basicDepPath, "llcppg.cfg"), fmt.Sprintf(" -I%s", basicDepHFile)),
	}

	for i := len(cleanups) - 1; i >= 0; i-- {
		defer cleanups[i]()
	}

	testFrom(t, name, depcjson, false, nil)
}

func testFromDir(t *testing.T, relDir string, gen bool) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal("Getwd failed:", err)
	}
	dir = path.Join(dir, relDir)
	fis, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal("ReadDir failed:", err)
	}
	for _, fi := range fis {
		name := fi.Name()
		if strings.HasPrefix(name, "_") || strings.HasPrefix(name, ".") {
			continue
		}
		t.Run(name, func(t *testing.T) {
			testFrom(t, name, dir+"/"+name, gen, nil)
		})
	}
}

func testFrom(t *testing.T, name, dir string, gen bool, validateFunc func(t *testing.T, pkg *convert.Package, converter *convert.Converter)) {
	confPath := filepath.Join(dir, "conf")
	cfgPath := filepath.Join(confPath, "llcppg.cfg")
	symbPath := filepath.Join(confPath, "llcppg.symb.json")
	pubPath := filepath.Join(confPath, "llcppg.pub")
	expect := filepath.Join(dir, "gogensig.expect")
	var expectContent []byte
	if !gen {
		var err error
		expectContent, err = os.ReadFile(expect)
		if err != nil {
			t.Fatal(expectContent)
		}
	}

	cfg, err := config.GetCppgCfgFromPath(cfgPath)
	if err != nil {
		t.Fatal(err)
	}

	// origin cflags + test deps folder cflags,because the test deps 's cflags is depend on machine
	if cfg.CFlags != "" {
		cfg.CFlags = env.ExpandEnv(cfg.CFlags)
	}

	cfg.CFlags = " -I" + filepath.Join(dir, "hfile") + " " + cfg.CFlags
	flagedCfgPath, err := config.CreateTmpJSONFile(args.LLCPPG_CFG, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(flagedCfgPath)

	if err != nil {
		t.Fatal(err)
	}

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(originalWd); err != nil {
			t.Fatal(err)
		}
	}()
	outputDir, err := ModInit(name)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(outputDir)

	bytes, err := config.SigfetchConfig(flagedCfgPath, confPath, cfg.Cplusplus)
	if err != nil {
		t.Fatal(err)
	}

	convertPkg, err := unmarshal.Pkg(bytes)
	if err != nil {
		t.Fatal(err)
	}

	cvt, err := convert.NewConverter(&convert.ConverterConfig{
		PkgName:   name,
		SymbFile:  symbPath,
		CfgFile:   flagedCfgPath,
		OutputDir: outputDir,
		PubFile:   pubPath,
		Pkg:       convertPkg,
	})

	if err != nil {
		t.Fatal(err)
	}
	cvt.Process()

	var res strings.Builder

	outDir, err := os.ReadDir(outputDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, fi := range outDir {
		if strings.HasSuffix(fi.Name(), "go.mod") || strings.HasSuffix(fi.Name(), "go.sum") || strings.HasSuffix(fi.Name(), "llcppg.pub") {
			continue
		} else {
			content, err := os.ReadFile(filepath.Join(outputDir, fi.Name()))
			if err != nil {
				t.Fatal(err)
			}
			res.WriteString(fmt.Sprintf("===== %s =====\n", fi.Name()))
			res.Write(content)
			res.WriteString("\n")
		}
	}

	pub, err := os.ReadFile(filepath.Join(outputDir, "llcppg.pub"))
	if err == nil {
		res.WriteString("===== llcppg.pub =====\n")
		res.Write(pub)
	}

	if gen {
		if err := os.WriteFile(expect, []byte(res.String()), 0644); err != nil {
			t.Fatal(err)
		}
	} else {
		expect := string(expectContent)
		got := res.String()
		if strings.TrimSpace(expect) != strings.TrimSpace(got) {
			t.Errorf("does not match expected.\nExpected:\n%s\nGot:\n%s", expect, got)
		}
	}

	if validateFunc != nil {
		validateFunc(t, cvt.GenPkg, cvt)
	}
}

// ===========================error
func TestNewAstConvert(t *testing.T) {
	_, err := convert.NewAstConvert(&convert.Config{
		PkgName:  "test",
		SymbFile: "",
		CfgFile:  "",
	})
	if err != nil {
		t.Fatal("NewAstConvert Fail")
	}
}

func TestNewAstConvertFail(t *testing.T) {
	_, err := convert.NewAstConvert(nil)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestVisitDone(t *testing.T) {
	pkg, err := convert.NewAstConvert(&convert.Config{
		PkgName:  "test",
		SymbFile: "",
		CfgFile:  "",
	})
	if err != nil {
		t.Fatal("NewAstConvert Fail")
	}
	pkg.SetVisitDone(func(pkg *convert.Package, incPath string) {
		if incPath != "test.h" {
			t.Fatal("doc path error")
		}
	})
	pkg.VisitDone("test.h")
}

func TestVisitFail(t *testing.T) {
	converter, err := convert.NewAstConvert(&convert.Config{
		PkgName:  "test",
		SymbFile: "",
		CfgFile:  "",
	})
	if err != nil {
		t.Fatal("NewAstConvert Fail")
	}

	converter.VisitStart("/path/to/temp.h", llcppg.Inter)

	// expect type
	converter.VisitTypedefDecl(&ast.TypedefDecl{
		Name: &ast.Ident{Name: "NormalType"},
		Type: &ast.BuiltinType{Kind: ast.Int},
	})

	// not appear in output,because expect error
	converter.VisitTypedefDecl(&ast.TypedefDecl{
		Name: &ast.Ident{Name: "Foo"},
		Type: nil,
	})

	errRecordType := &ast.RecordType{
		Tag: ast.Struct,
		Fields: &ast.FieldList{
			List: []*ast.Field{
				{Type: &ast.BuiltinType{Kind: ast.Int, Flags: ast.Double}},
			},
		},
	}
	// error field type for struct
	converter.VisitStruct(&ast.Ident{Name: "Foo"}, nil, &ast.TypeDecl{
		Name: &ast.Ident{Name: "Foo"},
		Type: errRecordType,
	})

	// error field type for anonymous struct
	converter.VisitStruct(&ast.Ident{Name: "Foo"}, nil, &ast.TypeDecl{
		Name: nil,
		Type: errRecordType,
	})

	converter.VisitStruct(&ast.Ident{Name: "Union (unnamed at /usr/local/Cellar/msgpack/6.0.2/include/msgpack/object.h:75:9)"}, nil, &ast.TypeDecl{
		Name: &ast.Ident{Name: "Union (unnamed at /usr/local/Cellar/msgpack/6.0.2/include/msgpack/object.h:75:9)"},
		Type: errRecordType,
	})

	converter.VisitEnumTypeDecl(&ast.EnumTypeDecl{
		Name: &ast.Ident{Name: "NormalType"},
		Type: &ast.EnumType{},
	})

	// error enum item for anonymous enum
	converter.VisitEnumTypeDecl(&ast.EnumTypeDecl{
		Name: nil,
		Type: &ast.EnumType{
			Items: []*ast.EnumItem{
				{Name: &ast.Ident{Name: "Item1"}},
			},
		},
	})

	converter.VisitFuncDecl(&ast.FuncDecl{
		Name: &ast.Ident{Name: "Foo"},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{Type: &ast.BuiltinType{Kind: ast.Int, Flags: ast.Double}},
				},
			},
		},
	})

	converter.VisitMacro(&ast.Macro{
		Name: "Foo",
		Tokens: []*ast.Token{
			{Token: ctoken.IDENT, Lit: "Foo"},
			{Token: ctoken.LITERAL, Lit: "1"},
		},
	})
	// not appear in output

	buf, err := converter.Pkg.WriteToBuffer("temp.go")
	if err != nil {
		t.Fatalf("WriteTo failed: %v", err)
	}

	expectedOutput :=
		`
package test

import (
	"github.com/goplus/llgo/c"
	_ "unsafe"
)

type NormalType c.Int
type Foo struct {
	Unused [8]uint8
}
`
	if strings.TrimSpace(expectedOutput) != strings.TrimSpace(buf.String()) {
		t.Errorf("does not match expected.\nExpected:\n%s\nGot:\n%s", expectedOutput, buf.String())
	}
}

func TestWritePkgFilesFail(t *testing.T) {
	tempDir, err := os.MkdirTemp(dir, "test_package_write_unwritable")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	converter, err := convert.NewAstConvert(&convert.Config{
		PkgName:   "test",
		SymbFile:  "",
		CfgFile:   "",
		OutputDir: tempDir,
	})
	if err != nil {
		t.Fatal("NewAstConvert Fail")
	}
	err = os.Chmod(tempDir, 0555)
	defer func() {
		if err := os.Chmod(tempDir, 0755); err != nil {
			t.Fatalf("Failed to change directory permissions: %v", err)
		}
	}()
	if err != nil {
		t.Fatalf("Failed to change directory permissions: %v", err)
	}
	converter.VisitStart("/path/to/test.h", llcppg.Inter)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but got: %v", r)
		}
	}()
	converter.WritePkgFiles()
}

func ModInit(name string) (string, error) {
	tempDir, err := os.MkdirTemp("", "gogensig-test")
	if err != nil {
		return "", err
	}
	outputDir := filepath.Join(tempDir, name)
	err = os.MkdirAll(outputDir, 0744)
	if err != nil {
		return "", err
	}
	projectRoot, err := filepath.Abs("../../../")
	if err != nil {
		return "", err
	}
	if err := os.Chdir(outputDir); err != nil {
		return "", err
	}

	err = config.RunCommand(outputDir, "go", "mod", "init", name)
	if err != nil {
		return "", err
	}
	err = config.RunCommand(outputDir, "go", "get", "github.com/goplus/llgo@main")
	if err != nil {
		return "", err
	}
	err = config.RunCommand(outputDir, "go", "get", "github.com/goplus/llcppg")
	if err != nil {
		return "", err
	}
	err = config.RunCommand(outputDir, "go", "mod", "edit", "-replace", "github.com/goplus/llcppg="+projectRoot)
	if err != nil {
		return "", err
	}
	return outputDir, nil
}
