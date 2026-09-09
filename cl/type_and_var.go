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

package cl

import (
	"go/types"

	lc "github.com/goplus/llcppg/lib/clang"
)

// -----------------------------------------------------------------------------

const (
	flagIsParam = 1 << iota
	flagIsStructField
	flagIsExtern
	flagIsTypedef
	flagRetType
)

func toType(ctx *blockCtx, typ lc.Type, flags int) types.Type {
	panic("todo")
}

// -----------------------------------------------------------------------------

func avoidKeyword(name *string) {
	switch *name {
	case "map", "type", "range", "chan", "var", "func", "go", "select",
		"defer", "package", "import", "interface", "fallthrough":
		*name += "_"
	}
}

// -----------------------------------------------------------------------------
