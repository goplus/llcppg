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
	"log"
	"os"
	"path/filepath"
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
	cl.SetDebug(cl.DbgFlagAll)
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

func testFromDir(t *testing.T, sel, relDir string, lang cl.Language) {
	idx := clang.CreateIndex(1, 1)
	defer idx.Dispose()

	cltest.TestFromDir(t, sel, relDir, func(t *testing.T, pkgDir string) {
		pkgDir, _ = filepath.Abs(pkgDir)

		conf, _ := cltest.LoadConf(pkgDir + "/in.cfg")
		srcFiles := conf.Files
		if len(srcFiles) == 0 {
			srcFiles = []string{"in.h"}
		}

		var mod tool.Module
		if conf.LoadLibcPubFile {
			m, err := tool.LoadModuleFrom(pkgDir)
			if err != nil {
				t.Errorf("failed to load module: %v", err)
				return
			}
			mod = m
		}

		files := make([]cl.Source, len(srcFiles))
		for i, srcFile := range srcFiles {
			presumedFile := filepath.Join(pkgDir, srcFile)
			u, err := idx.ParseTranslationUnit(
				clang.DetailedPreprocessingRecord, presumedFile, "-x", cltest.LanguageOf(lang))
			if err != nil {
				t.Fatal("ParseTranslationUnit failed:", err)
			}
			defer u.Dispose()
			files[i] = u
		}

		const pkgPrefix = "testcl/"
		rootDir, myPkgName := filepath.Split(pkgDir)
		imp := packages.NewImporter(nil, rootDir)
		pkg, err := cl.NewPackage(pkgPrefix+myPkgName, "foo", files, &cl.Config{
			Importer:        imp,
			LLGoPackage:     conf.LLGoPackage,
			Language:        lang,
			WrapFileHeader:  conf.WrapFileHeader,
			TypeAlias:       conf.TypeAlias,
			MacroIgnore:     conf.MacroIgnore,
			CFlags:          conf.CFlags,
			LoadLibcPubFile: conf.LoadLibcPubFile,
			DontKeepDoc:     !conf.KeepDoc,
			NameLookup:      nil,
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
	})
}

var langExts = [...]string{
	cl.LanguageC:   ".c",
	cl.LanguageCXX: ".cpp",
}

func TestC(t *testing.T) {
	testFromDir(t, "", "./_testc", cl.LanguageC)
}

func TestCpp(t *testing.T) {
	testFromDir(t, "", "./_testcpp", cl.LanguageCXX)
}

func TestPreprocessor(t *testing.T) {
	testFromDir(t, "", "./_testpp", cl.LanguageCXX)
}

// -----------------------------------------------------------------------------
