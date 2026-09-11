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
	"strings"

	"github.com/goplus/lib/c"
	"github.com/goplus/llcppg/clang"
	lc "github.com/goplus/llcppg/lib/clang"
)

func dump(node clang.Cursor, ns string, presumedFile *c.Char) {
	clang.VisitChildren(node, func(cur, parent clang.Cursor) clang.ChildVisitResult {
		if presumedFile != nil {
			loc := cur.Location()
			at := clang.PresumedFile(loc)
			cmpf := c.Strcmp(at.CStr(), presumedFile)
			at.Dispose()
			if cmpf != 0 {
				return clang.Continue
			}
		}
		kind := cur.Kind
		if kind == lc.CursorCXXAccessSpecifier {
			log.Println("==>", kind, "CXXAccessSpecifier", cur.CXXAccessSpecifier())
			return clang.Continue
		}
		name := ns + clang.String(cur)
		log.Println("==>", kind, clang.String(kind), name)
		switch kind {
		case lc.CursorFunctionDecl, lc.CursorCXXMethod, lc.CursorConstructor, lc.CursorDestructor:
		case lc.CursorClassDecl, lc.CursorNamespace:
			dump(cur, name+"::", presumedFile)
		}
		return clang.Continue
	})
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
	u := idx.ParseTranslationUnit(0, filename, "-x", lang)
	defer u.Dispose()

	usys := u.Underlying()
	spelling := usys.Spelling()
	defer spelling.Dispose()
	log.Println("==> TranslationUnit", c.GoString(spelling.CStr()))

	file := usys.File(spelling.CStr())
	loc := usys.GetLocationForOffset(file, 2)
	presumedFile := clang.PresumedFile(loc)
	defer presumedFile.Dispose()
	log.Println("==> PresumedFile", c.GoString(presumedFile.CStr()))

	root := u.Cursor()
	dump(root, "", presumedFile.CStr())
}
