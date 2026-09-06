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
	"github.com/goplus/llcppg/clang"
)

func dump(c clang.Cursor) {
	clang.VisitChildren(c, func(cur, parent clang.Cursor) clang.ChildVisitResult {
		return clang.Continue
	})
}

func main() {
	idx := clang.CreateIndex(0, 0)
	defer idx.Dispose()

	u := idx.ParseTranslationUnit(0, "")
	defer u.Dispose()

	dump(u.Cursor())
}
