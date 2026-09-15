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
	"go/token"

	"github.com/goplus/gogen"
	"github.com/goplus/llcppg/clang"
)

// -----------------------------------------------------------------------------

type wrapFile struct {
}

var langExts = [...]string{
	LanguageC:   ".c",
	LanguageCXX: ".cpp",
}

func newWrapFile(ctx *blockCtx) *wrapFile {
	ext := langExts[ctx.lang]
	filename := "_wrap/" + ctx.pkg.Types.Name() + ext
	llgoFiles := ctx.cflags + ": " + filename
	ctx.llgo.New(func(cb *gogen.CodeBuilder) int {
		cb.Val(llgoFiles)
		return 1
	}, 0, token.NoPos, nil, "LLGoFiles")
	return &wrapFile{}
}

func wrapInlineFunc(ctx *blockCtx, manglingName, origName string, fn clang.Cursor) string {
	if ctx.wrap == nil {
		ctx.wrap = newWrapFile(ctx)
	}
	wrapName := "_llcppg_" + manglingName
	// TODO(xsw): wrap inline func
	_ = fn
	_ = origName
	return wrapName
}

// -----------------------------------------------------------------------------
