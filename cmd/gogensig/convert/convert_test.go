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

	bytes, err := config.SigfetchConfigTest(flagedCfgPath, confPath, cfg.Cplusplus)
	if err != nil {
		t.Fatal(err)
	}

	convertPkg, err := unmarshal.Pkg(bytes)
	if err != nil {
		t.Fatal(err)
	}

	cvt, err := convert.NewConverter(&convert.Config{
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
	cvt.Convert()

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
func TestNewConvert(t *testing.T) {
	_, err := convert.NewConverter(&convert.Config{
		PkgName:  "test",
		SymbFile: "",
		CfgFile:  "",
	})
	if err != nil {
		t.Fatal("NewAstConvert Fail")
	}
}

func TestNewConvertFail(t *testing.T) {
	_, err := convert.NewConverter(nil)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestNewConvertReadPubFail(t *testing.T) {
	_, err := convert.NewConverter(&convert.Config{
		CfgFile: "./testdata/cjson/llcppg.cfg",
		PubFile: "./testdata/invalidpub/llcppg.pub",
	})
	if err == nil {
		t.Fatal("no error")
	}
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
	err = config.RunCommand(outputDir, "go", "get", "github.com/goplus/llgo@v0.10.0")
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
