package convert_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/goplus/gogen"
	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/cl/internal/cltest"
	"github.com/goplus/llcppg/cl/internal/convert"
	"github.com/goplus/llcppg/cmd/gogensig/unmarshal"
	llcppg "github.com/goplus/llcppg/config"
	"github.com/goplus/llgo/xtool/env"
)

func init() {
	convert.SetDebug(convert.DbgFlagAll)
}

func TestFromTestdata(t *testing.T) {
	testFromDir(t, "./_testdata", false)
}

func TestDepWithVersion(t *testing.T) {
	name := "_depwithversion"
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal("Getwd failed:", err)
	}
	testFrom(t, path.Join(dir, "_testdata", name), false, func(t *testing.T, pkg *convert.Package, cvt *convert.Converter) {
		modFile := filepath.Join(cvt.Conf.OutputDir, "go.mod")
		modContent, err := os.ReadFile(modFile)
		if err != nil {
			t.Fatal("Read go.mod failed:", err)
		}
		if !strings.Contains(string(modContent), "libxml2 v1.0.1") {
			t.Fatal(string(modContent), "\ngo.mod does not contain libxml2 v1.0.1")
		}
		if !strings.Contains(string(modContent), "zlib v1.0.1") {
			t.Fatal(string(modContent), "\ngo.mod does not contain zlib v1.0.1")
		}
	})
}

func TestSysToPkg(t *testing.T) {
	name := "_systopkg"
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal("Getwd failed:", err)
	}
	testFrom(t, path.Join(dir, "_testdata", name), false, func(t *testing.T, pkg *convert.Package, cvt *convert.Converter) {
		// check FileMap's info is right
		inFileMap := func(file string) {
			/* TODO(xsw): remove this
			_, ok := cvt.FileMap[file]
			if !ok {
				t.Fatal("File not found in FileMap:", file)
			}
			*/
		}
		for _, decl := range cvt.Pkg.Decls {
			switch decl := decl.(type) {
			case *ast.TypeDecl:
				inFileMap(decl.Object.Loc.File)
			case *ast.EnumTypeDecl:
				inFileMap(decl.Object.Loc.File)
			case *ast.TypedefDecl:
				inFileMap(decl.Object.Loc.File)
			case *ast.FuncDecl:
				inFileMap(decl.Object.Loc.File)
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

		cfg, err := llcppg.GetConfFromFile(cfgPath)
		if err != nil {
			t.Fatal(err, cfgPath)
		}
		cfg.CFlags = cfg.CFlags + incFlag
		err = createJSONFile(cfgPath, cfg)
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
	depcjsonConf := path.Join(depcjson, "conf", llcppg.LLCPPG_CFG)

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
		inc(buildTestPath(thirdDepPath, llcppg.LLCPPG_CFG), fmt.Sprintf(" -I%s -I%s", thirdDepHFile, basicDepHFile)),
		inc(buildTestPath(thirdDep2Path, llcppg.LLCPPG_CFG), fmt.Sprintf(" -I%s -I%s -I%s", thirdDep2HFile, thirdDepHFile, basicDepHFile)),
		inc(buildTestPath(basicDepPath, llcppg.LLCPPG_CFG), fmt.Sprintf(" -I%s", basicDepHFile)),
	}

	for i := len(cleanups) - 1; i >= 0; i-- {
		defer cleanups[i]()
	}

	testFrom(t, depcjson, false, nil)
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
		if strings.HasPrefix(name, "_") {
			continue
		}
		t.Run(name, func(t *testing.T) {
			testFrom(t, dir+"/"+name, gen, nil)
		})
	}
}

func testFrom(t *testing.T, dir string, gen bool, validateFunc func(t *testing.T, pkg *convert.Package, converter *convert.Converter)) {
	confPath := filepath.Join(dir, "conf")
	cfgPath := filepath.Join(confPath, llcppg.LLCPPG_CFG)
	symbPath := filepath.Join(confPath, llcppg.LLCPPG_SYMB)
	expect := filepath.Join(dir, "gogensig.expect")
	var expectContent []byte
	if !gen {
		var err error
		expectContent, err = os.ReadFile(expect)
		if err != nil {
			t.Fatal(expectContent)
		}
	}

	cfg, err := llcppg.GetConfFromFile(cfgPath)
	if err != nil {
		t.Fatal(err)
	}

	// origin cflags + test deps folder cflags,because the test deps 's cflags is depend on machine
	if cfg.CFlags != "" {
		cfg.CFlags = env.ExpandEnv(cfg.CFlags)
	}

	cfg.CFlags = " -I" + filepath.Join(dir, "hfile") + " " + cfg.CFlags
	flagedCfgPath, err := createTmpJSONFile(llcppg.LLCPPG_CFG, cfg)
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
	outputDir, err := prepareEnv(cfg.Name, cfg.Deps)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(outputDir)

	bytes, err := callSigfetch(flagedCfgPath, confPath)
	if err != nil {
		t.Fatal(err)
	}

	convertPkg, err := unmarshal.Pkg(bytes)
	if err != nil {
		t.Fatal(err)
	}

	cvt, err := convert.NewConverter(&convert.Config{
		PkgPath:   ".",
		PkgName:   cfg.Name,
		OutputDir: outputDir,
		Pkg:       convertPkg.File,
		NC:        cltest.NC(&cfg, convertPkg.FileMap, cltest.GetConvSym(symbPath)),
		Deps:      cfg.Deps,
		Libs:      cfg.Libs,
		GenThird:  true,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = cvt.Convert()
	if err != nil {
		t.Fatal(err)
	}

	pkg := cvt.GenPkg.Pkg()
	pkg.ForEachFile(func(fname string, _ *gogen.File) {
		if fname != "" { // gogen default fname
			outFile := filepath.Join(outputDir, fname)
			e := pkg.WriteFile(outFile, fname)
			if e != nil {
				t.Fatal(e)
			}
		}
	})

	// todo:reuse same write logic
	err = llcppg.WritePubFile(filepath.Join(outputDir, llcppg.LLCPPG_PUB), cvt.GenPkg.Pubs)
	if err != nil {
		t.Fatal(err)
	}

	err = runCommand(outputDir, "go", "fmt", ".")
	if err != nil {
		t.Fatal(err)
	}

	err = runCommand(outputDir, "go", "mod", "tidy")
	if err != nil {
		t.Fatal(err)
	}

	var res strings.Builder

	outDir, err := os.ReadDir(outputDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, fi := range outDir {
		if strings.HasSuffix(fi.Name(), "go.mod") || strings.HasSuffix(fi.Name(), "go.sum") || strings.HasSuffix(fi.Name(), llcppg.LLCPPG_PUB) {
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

	pub, err := os.ReadFile(filepath.Join(outputDir, llcppg.LLCPPG_PUB))
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

func TestNewConvert(t *testing.T) {
	cfg := &llcppg.Config{
		Libs: "${pkg-config --libs xxx}",
	}

	_, err := convert.NewConverter(&convert.Config{
		PkgPath: ".",
		PkgName: "test",
		NC:      cltest.NC(cfg, nil, cltest.NewConvSym()),
		Deps:    cfg.Deps,
		Libs:    cfg.Libs,
	})
	if err != nil {
		t.Fatal("NewAstConvert Fail")
	}
}

func TestModInitFail(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gogensig-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	invalidMod := "github.com/user!name/project"

	t.Run("mod init fail", func(t *testing.T) {
		err = convert.ModInit([]string{}, tempDir, invalidMod)
		if err == nil {
			t.Fatal("no error")
		}
	})
	t.Run("dep get fail", func(t *testing.T) {
		err = convert.ModInit([]string{invalidMod}, tempDir, "")
		if err == nil {
			t.Fatal("no error")
		}
	})
}

func prepareEnv(name string, deps []string) (string, error) {
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

	// with the same module to import the internal/testdata/pkg
	err = runCommand(outputDir, "go", "mod", "init", "github.com/goplus/llcppg/cl/internal/"+name)
	if err != nil {
		return "", err
	}

	err = runCommand(outputDir, "go", "get", "github.com/goplus/llcppg")
	if err != nil {
		return "", err
	}
	err = runCommand(outputDir, "go", "mod", "edit", "-replace", "github.com/goplus/llcppg="+projectRoot)
	if err != nil {
		return "", err
	}

	err = convert.ModInit(deps, outputDir, "")
	if err != nil {
		return "", err
	}
	return outputDir, nil
}

func runCommand(dir, cmdName string, args ...string) error {
	execCmd := exec.Command(cmdName, args...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Dir = dir
	return execCmd.Run()
}

func createTmpJSONFile(filename string, data any) (string, error) {
	filePath := filepath.Join(os.TempDir(), filename)
	err := createJSONFile(filePath, data)
	return filePath, err
}

func createJSONFile(filepath string, data any) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func callSigfetch(configFile string, dir string) ([]byte, error) {
	cmd := exec.Command("llcppsigfetch", "-v", configFile)
	cmd.Dir = dir

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
