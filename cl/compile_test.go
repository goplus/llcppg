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

package cl_test

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"maps"
	"os"
	"path/filepath"
	"slices"
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
	log.SetFlags(0)
}

// -----------------------------------------------------------------------------

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

func testGenGo(t *testing.T, pkg *gogen.Package, dir string, exp any) {
	var b bytes.Buffer
	err := pkg.WriteTo(&b)
	if err != nil {
		t.Fatal("gogen.WriteTo failed:", err)
	}
	testDiff(t, dir, "/out.go.txt", &b, exp)
}

func testFromDir(t *testing.T, sel, relDir string, srcLang cl.Language, genPubFile ...bool) {
	idx := clang.CreateIndex(1, 1)
	defer idx.Dispose()

	savePubFile := genPubFile != nil && genPubFile[0]
	cltest.TestFromDir(t, sel, relDir, func(t *testing.T, pkgDir string) {
		pkgDir, _ = filepath.Abs(pkgDir)

		conf, err := cltest.LoadConf(pkgDir + "/in.cfg")
		if err != nil && !os.IsNotExist(err) {
			t.Fatal(err)
		}
		srcFiles := conf.Files
		if len(srcFiles) == 0 {
			srcFiles = []string{"in.h"}
		}

		var fset *token.FileSet
		var goFiles map[string]*ast.File
		if len(conf.GoFiles) != 0 {
			fset = token.NewFileSet()
			goFiles, err = tool.LoadGoFiles(fset, pkgDir, conf.GoFiles, parser.AllErrors)
			if err != nil {
				t.Error("LoadGoFiles failed:", err)
				return
			}
		}

		var mod tool.Module
		if conf.LoadLibcPubFile {
			m, err := tool.LoadModuleFrom(pkgDir)
			if err != nil {
				t.Error("failed to load module:", err)
				return
			}
			mod = m
		}

		const pkgPrefix = "testcl/"
		rootDir, myPkgName := filepath.Split(pkgDir)
		imp := packages.NewImporter(nil, rootDir)
		lang := srcLang
		if myPkgName == "cxx" { // as a special case for C++ tests
			lang = cl.LanguageCXX
		}
		fakePkgName := "foo"
		if savePubFile { // use real package name
			fakePkgName = myPkgName
		}

		files := make([]cl.Source, len(srcFiles))
		for i, srcFile := range srcFiles {
			presumedFile := filepath.Join(pkgDir, srcFile)
			u, err := idx.ParseTranslationUnit(
				clang.DetailedPreprocessingRecord, presumedFile, "-x", cltest.LanguageOf(lang), "-fparse-all-comments")
			if err != nil {
				t.Fatal("ParseTranslationUnit failed:", err)
			}
			defer u.Dispose()
			files[i] = u
		}

		names := slices.Sorted(maps.Keys(conf.Define))
		defines := make([]cl.ConstDef, len(names))
		for i, name := range names {
			defines[i] = cl.ConstDef{
				Name:  name,
				Value: conf.Define[name],
			}
		}

		pkg, err := cl.NewPackage(pkgPrefix+myPkgName, fakePkgName, files, &cl.Config{
			Fset:              fset,
			Importer:          imp,
			Define:            defines,
			LLGoPackage:       conf.LLGoPackage,
			Language:          lang,
			GoFiles:           goFiles,
			WrapFileHeader:    conf.WrapFileHeader,
			NSPrefix:          conf.NSPrefix,
			FuncPrefix:        conf.FuncPrefix,
			EnumPrefix:        conf.EnumPrefix,
			Rename:            conf.Rename,
			FieldAlias:        conf.FieldAlias,
			NewCheck:          conf.NewCheck,
			MethodCheck:       conf.MethodCheck,
			Method:            conf.Method,
			TypeDef:           conf.TypeDef,
			TypeAlias:         conf.TypeAlias,
			TypeIgnore:        conf.TypeIgnore,
			TypePrefix:        conf.TypePrefix,
			TypeSuffix:        conf.TypeSuffix,
			TypeAbbr:          conf.TypeAbbr,
			MacroIgnore:       conf.MacroIgnore,
			MacroFunc:         conf.MacroFunc,
			MacroGroup:        conf.MacroGroup,
			MacroPrefix:       conf.MacroPrefix,
			CFlags:            conf.CFlags,
			EnumAlwaysUntyped: conf.EnumAlwaysUntyped,
			ForceCamelCase:    conf.ForceCamelCase,
			LoadLibcPubFile:   conf.LoadLibcPubFile,
			NameLookup:        nil,
			PubFileLookup: func(pkgPath string) (pubFile string, ok bool) {
				if name, ok := strings.CutPrefix(pkgPath, pkgPrefix); ok {
					return filepath.Join(rootDir, name, "llcppg.pub"), true
				} else if conf.LoadLibcPubFile {
					return mod.PubFileLookup(pkgPath)
				}
				return
			},
			PackageOf: func(headerFile string) (pkgPath string, ok bool) {
				dir := filepath.Dir(headerFile)
				tRootDir, tPkgName := filepath.Split(dir)
				if ok = tRootDir == rootDir; ok {
					pkgPath = pkgPrefix + tPkgName
				}
				return
			},
		})
		if err != nil {
			t.Error("cl.NewPackage:", err)
			return
		}
		exp, _ := os.ReadFile(pkgDir + "/out.go")
		testGenGo(t, pkg.Package, pkgDir, exp)
		wrapFile := "/wrap" + langExts[lang]
		wrap, _ := os.ReadFile(pkgDir + wrapFile)
		if pkg.Wrap != nil {
			testDiff(t, pkgDir, wrapFile+".txt", &pkg.Wrap.Content, wrap)
		}
		if savePubFile {
			pubExp, _ := os.ReadFile(pkgDir + "/llcppg.pub")
			pub := cl.PubFile(pkg.Public)
			testDiff(t, pkgDir, "/llcppg.pub.txt", bytes.NewBuffer(pub), pubExp)
		}
	})
}

var langExts = [...]string{
	cl.LanguageC:   ".c",
	cl.LanguageCXX: ".cpp",
}

func TestC(t *testing.T) {
	cl.SetDebug(cl.DbgFlagMacroFunc)
	defer cl.SetDebug(0)

	testFromDir(t, "", "./_testc", cl.LanguageC)
}

func TestCpp(t *testing.T) {
	testFromDir(t, "", "./_testcpp", cl.LanguageCXX)
}

func TestPP(t *testing.T) {
	testFromDir(t, "", "./_testpp", cl.LanguageCXX)
}

func TestPub(t *testing.T) {
	cl.SetDebug(cl.DbgFlagAll)
	defer cl.SetDebug(0)

	testFromDir(t, "", "./_testpub", cl.LanguageC, true)
}

// -----------------------------------------------------------------------------
