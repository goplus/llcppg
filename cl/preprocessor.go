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
	"log"

	"github.com/goplus/gogen"
	"github.com/goplus/llcppg/clang"
	lc "github.com/llarhub/clang-c"
)

// -----------------------------------------------------------------------------

func loadInclude(ctx *pkgCtx, decl clang.Cursor) {
	file := decl.IncludedFile()
	if file == clang.InvalidFile {
		return
	}
	includeFile := clang.FileName(file) // already absolute path
	if debugCompileDecl {
		log.Println("include", includeFile)
	}

	if _, ok := ctx.lastSeen[includeFile]; ok {
		return // already loaded
	}
	ctx.thisSeen[includeFile] = none{}

	pkgPath, ok := ctx.pkgOf(includeFile)
	if !ok {
		log.Println("[WARN] loadInclude: package not found for", includeFile)
		return
	}
	if ctx.pkg.Types.Path() != pkgPath {
		ctx.importPkg(pkgPath)
	}
}

// -----------------------------------------------------------------------------

func emitMacroFunc(ctx *pkgCtx, decl clang.Cursor, macroName string, paramNames []string, tu clang.TranslationUnit, body []lc.Token) {
	if ctx.isMacroIgnored(macroName) {
		if debugCompileDecl {
			log.Println("macroFunc", macroName, paramNames, "- ignored")
		}
		return
	}
	if len(body) >= 3 { // fn()
		if v, ok := evalFuncCallExpr(ctx, tu, body); ok {
			if debugMacroFunc {
				log.Println("macroFunc", macroName, paramNames)
			}
			doCompileAutoMacroFunc(ctx, decl, macroName, paramNames, v)
		}
	}
}

// -----------------------------------------------------------------------------

func emitMacroConst(defs *gogen.ConstDefs, ctx *pkgCtx, decl clang.Cursor, docRaw string, grouped bool) {
	origName := clang.String(decl)
	if ctx.isMacroIgnored(origName) {
		if debugCompileDecl {
			log.Println("macro", origName, "- ignored")
		}
		return
	}

	tu := clang.TU(decl)
	tokens, dispose := tu.Tokenize(decl.Extent())
	defer dispose()
	if debugCompileDecl {
		log.Println("macro", origName, "-", len(tokens), "tokens")
	}
	if len(tokens) > 1 {
		if v, ok := evalConstExpr(ctx, tu, tokens[1:]); ok && v != nil { // not NULL
			if oldv, ok := ctx.macroVals[origName]; ok {
				if oldv == v {
					return // allow redefinition with the same value
				}
			}
			at := defs.NewPos()
			if grouped {
				at.Doc = doDocCommentGroup(toLineComments(docRaw))
			} else {
				defs.SetComments(docCommentGroup(decl))
			}
			name := ctx.macroName(origName)
			defs.NewAt(at, func(cb *gogen.CodeBuilder) int {
				cb.Val(v)
				return 1
			}, 0, goNodePos(ctx, decl), nil, name)
			ctx.macroVals[origName] = v
		}
	}
}

// -----------------------------------------------------------------------------

func loadMacroConst(ctx *pkgCtx, defs *gogen.ConstDefs, decl clang.Cursor) {
	emitMacroConst(defs, ctx, decl, "", false)
}

func loadMacroConstGroup(ctx *pkgCtx, defs *gogen.ConstDefs, items []clang.Cursor) {
	firstRaw := clang.RawCommentsOf(items[0])
	lastRaw := clang.RawCommentsOf(items[len(items)-1])

	// Resolve the group's doc comment from the first macro's preceding comment
	// (priority) or the last macro's trailing comment.
	var docRaw string
	if firstRaw.Preceding != "" {
		docRaw, firstRaw.Preceding = firstRaw.Preceding, ""
	} else if lastRaw.Trailing != "" {
		docRaw, lastRaw.Trailing = lastRaw.Trailing, ""
	}

	if docRaw != "" {
		doc := doDocCommentGroup(toLineComments(docRaw))
		defs.SetComments(doc)
	}

	for i, m := range items {
		var raw clang.RawComments
		if i == 0 {
			raw = firstRaw
		} else if i == len(items)-1 {
			raw = lastRaw
		} else {
			raw = clang.RawCommentsOf(m)
		}
		emitMacroConst(defs, ctx, m, raw.Text(), true)
	}
}

// -----------------------------------------------------------------------------
