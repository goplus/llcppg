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
	"testing"

	"github.com/goplus/gogen"
	"github.com/goplus/gogen/packages"
	"github.com/goplus/llcppg/cl"
	"github.com/goplus/llcppg/cl/cltest"
	"github.com/goplus/llcppg/clang"
	"github.com/qiniu/x/test"
)

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
	log.Println("==> testGenGo", dir)
	var b bytes.Buffer
	err := pkg.WriteTo(&b)
	if err != nil {
		t.Fatal("gogen.WriteTo failed:", err)
	}
	log.Println("==> testGenGo", dir, "len:", b.Len())
	testDiff(t, dir, "/result.txt", &b, exp)
}

func testFromDir(t *testing.T, sel, relDir, lang string) {
	cltest.TestFromDir(t, sel, relDir, func(t *testing.T, pkgDir string) {
		idx := clang.CreateIndex(0, 0)
		defer idx.Dispose()

		filename := pkgDir + "/in.h"
		u := idx.ParseTranslationUnit(0, filename, "-x", lang)
		defer u.Dispose()

		imp := packages.NewImporter(nil)
		file := u.File(filename)
		pkg, err := cl.NewPackage("", "foo", cl.Source{TU: u, Handle: file}, &cl.Config{
			Importer:   imp,
			NameLookup: cltest.MockNameLookup,
		})
		if err != nil {
			t.Error("cl.NewPackage:", err)
			return
		}
		exp, _ := os.ReadFile(pkgDir + "/out.go")
		testGenGo(t, pkg.Package, pkgDir, exp)
	})
}

func _TestMockC(t *testing.T) {
	cl.SetDebug(cl.DbgFlagAll)
	testFromDir(t, "", "./_testmockc", "c")
}

// -----------------------------------------------------------------------------
