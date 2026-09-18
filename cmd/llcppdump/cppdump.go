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

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/goplus/llcppg/clang"
	lc "github.com/goplus/llcppg/lib/clang"
)

func dump(node clang.Cursor, ns, dir string) {
	clang.VisitChildren(node, func(cur, parent clang.Cursor) clang.ChildVisitResult {
		at := clang.PresumedFile(cur.Location())
		if filepath.Dir(at) != dir {
			return clang.Continue
		}
		kind := cur.Kind
		if kind == lc.CursorCXXAccessSpecifier {
			log.Println("==>", kind, "CXXAccessSpecifier", cur.CXXAccessSpecifier())
			return clang.Continue
		}
		name := ns + clang.String(cur)
		log.Println("==>", kind, clang.String(kind), name, typeOf(cur.Type()))
		switch kind {
		case lc.CursorFunctionDecl, lc.CursorCXXMethod, lc.CursorConstructor, lc.CursorDestructor:
		case lc.CursorClassDecl, lc.CursorNamespace:
			dump(cur, name+"::", dir)
		}
		return clang.Continue
	})
}

func typeOf(t lc.Type) string {
	switch t.Kind {
	case lc.TypeElaborated:
		return typeOf(t.NamedType())
	default:
		return clang.String(t)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: llcppdump <source-file> [<language>]")
		return
	}

	idx := clang.CreateIndex(0, 0)
	defer idx.Dispose()

	filename := os.Args[1]
	lang := "c++"
	if len(os.Args) > 2 {
		lang = strings.ToLower(os.Args[2])
	}
	filename, _ = filepath.Abs(filename)
	log.Println("==> dump", filename, "as", lang)
	u := idx.ParseTranslationUnit(clang.DetailedPreprocessingRecord, filename, "-x", lang)
	defer u.Dispose()

	options := lc.DefaultDiagnosticDisplayOptions()
	u.VisitDiagnostics(func(diag clang.Diagnostic) {
		fmt.Fprintln(os.Stderr, diag.Format(options))
	})

	root := u.Cursor()
	dump(root, "", filepath.Dir(filename))
}
