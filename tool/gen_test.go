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

package tool_test

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/goplus/gogen"
	"github.com/goplus/gogen/packages"
	"github.com/goplus/llcppg/cl"
	"github.com/goplus/llcppg/cl/cltest"
	"github.com/goplus/llcppg/clang"
	"github.com/goplus/llcppg/tool"
	"github.com/qiniu/x/test"
)

func init() {
	cl.SetDebug(cl.DbgFlagAll)
	tool.SetDebug(tool.DbgFlagAll)
}

func testDiff(t *testing.T, dir string, outfname string, b *bytes.Buffer, exp any) {
	if expected, ok := exp.(string); ok {
		result := b.String()
		if result != expected {
			t.Errorf("\nResult:\n%s\nExpected:\n%s\n", result, expected)
		}
	} else if test.Diff(t, dir+outfname, b.Bytes(), exp.([]byte)) {
		t.Error(dir, ": unexpect result")
	}
}

func testGenGo(t *testing.T, pkg *gogen.Package, dir, fname string, exp any) {
	var b bytes.Buffer
	if fname != "" && tool.GeneratedHeader != "" {
		b.WriteString(tool.GeneratedHeader)
	}
	err := pkg.WriteTo(&b, fname)
	if err != nil {
		t.Fatal("gogen.WriteTo failed:", err)
	}
	var outfname string
	if fname == "" {
		outfname = "/out.go.txt"
	} else {
		outfname = "/" + fname + ".txt"
	}
	testDiff(t, dir, outfname, &b, exp)
}

func testSingleFile(t *testing.T, idx clang.Index, pkgDir, headerDir, headerFile string, conf *tool.Config) {
	myPkgName := headerFile[len(headerDir) : len(headerFile)-2]
	if myPkgName == "Index" || myPkgName == "DenseMapInfo" && runtime.GOOS != "darwin" {
		log.Println("==> only test Index on macOS")
		return
	}

	log.Println("============== testSingleFile: package", myPkgName, "==============")

	lang, ok := conf.Lang()
	if !ok {
		t.Errorf("invalid language: %q", conf.Language)
		return
	}

	mod, err := tool.LoadModuleFrom(pkgDir)
	if err != nil {
		t.Errorf("failed to load module: %v", err)
		return
	}

	incDirs := make([]string, 1, 6)
	incDirs[0] = pkgDir + "/include"
	for _, dep := range conf.Deps {
		dep = "github.com/llarhub/" + dep
		incDir, err := mod.IncludeDir(dep)
		if err != nil {
			t.Errorf("failed to get include dir for dependency %q: %v", dep, err)
			return
		}
		incDirs = append(incDirs, incDir)
	}
	incDirs = append(incDirs, conf.StdlibDirs(pkgDir)...)
	log.Println("==> includeDirs:", incDirs)

	flags := tool.ParseFlags(incDirs, conf.Language)
	u, err := idx.ParseTranslationUnit(clang.DetailedPreprocessingRecord, headerFile, flags...)
	if err != nil {
		t.Error("ParseTranslationUnit failed:", err)
		return
	}
	defer u.Dispose()

	files := []cl.Source{u}
	imp := packages.NewImporter(nil, headerDir)
	pkgPrefix := conf.Name
	if pos := strings.Index(pkgPrefix, "/"); pos > 0 {
		pkgPrefix = pkgPrefix[:pos+1]
	} else {
		pkgPrefix += "/"
	}
	pkg, err := cl.NewPackage(pkgPrefix+myPkgName, myPkgName, files, &cl.Config{
		Importer:    imp,
		LLGoPackage: conf.LLGoPackage,
		Language:    lang,
		Class:       conf.Class,
		NonClass:    conf.NonClass,
		FuncPrefix:  conf.FuncPrefix,
		EnumPrefix:  conf.EnumPrefix,
		MacroPrefix: conf.MacroPrefix,
		TypePrefix:  conf.TypePrefix,
		TypeAbbr:    conf.TypeAbbr,
		Rename:      conf.Rename,
		TypeIgnore:  conf.TypeIgnore,
		NSIgnore:    conf.NSIgnore,
		NameLookup:  nil,
		PubFileLookup: func(pkgPath string) (pubFile string, ok bool) {
			if name, ok := strings.CutPrefix(pkgPath, pkgPrefix); ok {
				return filepath.Join(pkgDir, name, "llcppg.pub"), true
			}
			return mod.PubFileLookup(pkgPath)
		},
		PackageOf: func(headerFile string) (pkgPath string, ok bool) {
			tRootDir, tPkgName := filepath.Split(headerFile)
			if ok = tRootDir == headerDir; ok {
				pkgPath = pkgPrefix + tPkgName[:len(tPkgName)-2]
			} else if ok = tRootDir == pkgDir+"/cstdlib/"; ok {
				pkgPath = pkgPrefix + "cstdlib"
			} else if strings.Contains(headerFile, "github.com/llarhub/libcxx/c") {
				pkgPath, ok = "github.com/llarhub/libcxx", true
			}
			return
		},
	})
	if err != nil {
		t.Error("cl.NewPackage:", err)
		return
	}
	pkgDir = filepath.Join(pkgDir, myPkgName)
	os.Mkdir(pkgDir, 0755)
	exp, _ := os.ReadFile(pkgDir + "/out.go")
	testGenGo(t, pkg.Package, pkgDir, "", exp)
}

func testFromDir(t *testing.T, sel, relDir string, single bool, subPkg ...string) {
	dirSel := sel
	if single {
		dirSel = ""
	}
	cltest.TestFromDir(t, dirSel, relDir, func(t *testing.T, pkgDir string) {
		idx := clang.CreateIndex(1, 1)
		defer idx.Dispose()

		pkgDir, _ = filepath.Abs(pkgDir)
		conf, err := tool.LoadConf(pkgDir + "/llcppg.cfg")
		if err != nil {
			log.Fatal("LoadConf failed:", err)
		}

		var pkgSel string
		if len(subPkg) > 0 {
			pkgSel = subPkg[0]
		} else {
			switch len(conf.Pkgs) {
			case 0:
			case 1:
				pkgSel = conf.Pkgs[0]
			default:
				t.Fatal("conf.Pkgs can't be multi-packages for testing")
			}
		}
		if pkgSel != "" {
			cfgFile := pkgDir + "/llcppg-" + pkgSel + ".cfg"
			subConf, err := tool.LoadConf(cfgFile)
			if err != nil {
				t.Fatal("LoadConf failed:", err)
			}
			subConf.Apply(&conf)
			conf = subConf
		}

		if single {
			headerDir := filepath.Join(pkgDir, conf.Dir)
			fis, err := os.ReadDir(headerDir)
			if err != nil {
				log.Fatal("os.ReadDir failed:", err)
			}
			headerDir += string(os.PathSeparator)
			for _, fi := range fis {
				if fi.IsDir() {
					continue
				}
				name := fi.Name()
				if !strings.HasSuffix(name, ".h") {
					continue
				}
				pkgName := name[:len(name)-2]
				if sel != "" && pkgName != sel {
					continue
				}
				headerFile := headerDir + name
				t.Run(pkgName, func(t *testing.T) {
					testSingleFile(t, idx, pkgDir, headerDir, headerFile, &conf)
				})
			}
			return
		}

		if runtime.GOOS != "darwin" {
			log.Println("==> only test on macOS") // TODO(xsw): cross-platform
			return
		}

		pkg, lang, err := conf.NewPackage("", pkgDir, idx)
		if err != nil {
			t.Error("conf.NewPackage:", err)
			return
		}
		pkg.ForEachFile(func(fname string, file *gogen.File) {
			if file.Empty() {
				return // skip empty Go files
			}
			exp, _ := os.ReadFile(pkgDir + "/" + fname)
			testGenGo(t, pkg.Package, pkgDir, fname, exp)
		})
		wrapFile := "/wrap" + langExts[lang]
		wrap, _ := os.ReadFile(pkgDir + wrapFile)
		if pkg.Wrap != nil {
			testDiff(t, pkgDir, wrapFile+".txt", &pkg.Wrap.Content, wrap)
		}
	})
}

var langExts = [...]string{
	cl.LanguageC:   ".c",
	cl.LanguageCXX: ".cpp",
}

func TestC(t *testing.T) {
	testFromDir(t, "", "./_testc", false)
}

func TestSingleC(t *testing.T) {
	testFromDir(t, "", "./_testc", true)
}

func TestLLVM_AMDGPUAddrSpace(t *testing.T) {
	testFromDir(t, "AMDGPUAddrSpace", "./_testcpp", true, "system")
}

func TestLLVM_AMDHSAKernelDescriptor(t *testing.T) {
	testFromDir(t, "AMDHSAKernelDescriptor", "./_testcpp", true, "system")
}

func TestLLVM_AtomicOrdering(t *testing.T) {
	testFromDir(t, "AtomicOrdering", "./_testcpp", true, "system")
}

func TestLLVM_Atomic(t *testing.T) {
	testFromDir(t, "Atomic", "./_testcpp", true, "system")
}

func TestLLVM_Compiler(t *testing.T) {
	testFromDir(t, "Compiler", "./_testcpp", true, "system")
}

func TestLLVM_DenseMapInfo(t *testing.T) {
	testFromDir(t, "DenseMapInfo", "./_testcpp", true, "adt")
}

func _TestLLVM_String(t *testing.T) {
	testFromDir(t, "StringRef", "./_testcpp", true, "adt")
}
