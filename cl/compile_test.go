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
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/goplus/gogen"
	"github.com/goplus/gogen/packages"
	"github.com/goplus/llcppg/cl"
	"github.com/goplus/llcppg/cl/cltest"
	"github.com/goplus/llcppg/clang"
	"github.com/qiniu/x/test"

	lc "github.com/goplus/llcppg/lib/clang"
)

func init() {
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
	cltest.TestFromDir(t, sel, relDir, func(t *testing.T, pkgDir string) {
		idx := clang.CreateIndex(0, 0)
		defer idx.Dispose()

		pkgDir, _ = filepath.Abs(pkgDir)
		conf, _ := cltest.LoadConf(pkgDir + "/in.cfg")
		srcFiles := conf.Files
		if len(srcFiles) == 0 {
			srcFiles = []string{"in.h"}
		}

		files := make([]cl.Source, len(srcFiles))
		options := lc.DefaultDiagnosticDisplayOptions()
		for i, srcFile := range srcFiles {
			presumedFile := filepath.Join(pkgDir, srcFile)
			u := idx.ParseTranslationUnit(
				clang.DetailedPreprocessingRecord, presumedFile, "-x", cltest.LanguageOf(lang))
			defer u.Dispose()
			files[i] = cl.Source{TU: u}
			u.VisitDiagnostics(func(diag clang.Diagnostic) {
				fmt.Fprintln(os.Stderr, diag.Format(options))
			})
		}

		const pkgPrefix = "testcl/"
		rootDir, myPkgName := filepath.Split(pkgDir)
		imp := packages.NewImporter(nil, rootDir)
		pkg, err := cl.NewPackage(pkgPrefix+myPkgName, "foo", files, &cl.Config{
			Importer:       imp,
			LLGoPackage:    conf.LLGoPackage,
			Language:       lang,
			WrapFileHeader: conf.WrapFileHeader,
			CFlags:         conf.CFlags,
			NameLookup:     nil,
			PubFileLookup: func(pkgPath string) (pubFile string, ok bool) {
				if name, ok := strings.CutPrefix(pkgPath, pkgPrefix); ok {
					return filepath.Join(rootDir, name, "llcppg.pub"), true
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
	testFromDir(t, "union_basic", "./_testc", cl.LanguageC)
}

func _TestCpp(t *testing.T) {
	testFromDir(t, "", "./_testcpp", cl.LanguageCXX)
}

func _TestPreprocessor(t *testing.T) {
	testFromDir(t, "", "./_testpp", cl.LanguageCXX)
}

// -----------------------------------------------------------------------------
