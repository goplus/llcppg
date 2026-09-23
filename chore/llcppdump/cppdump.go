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
	"github.com/goplus/llcppg/tool"

	lc "github.com/goplus/llcppg/lib/clang"
)

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
	tool.Dump(root, "", filepath.Dir(filename))
}
