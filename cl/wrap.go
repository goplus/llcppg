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
	"bytes"
	"go/token"

	"github.com/goplus/gogen"
	"github.com/goplus/lib/c"
	"github.com/goplus/llcppg/clang"
	lc "github.com/goplus/llcppg/lib/clang"
)

// -----------------------------------------------------------------------------

type WrapFile struct {
	Filename string
	Content  bytes.Buffer
}

var langExts = [...]string{
	LanguageC:   ".c",
	LanguageCXX: ".cpp",
}

func newWrapFile(ctx *pkgCtx) *WrapFile {
	ext := langExts[ctx.lang]
	filename := "_wrap/llcppg" + ext
	llgoFiles := ctx.cflags + ": " + filename
	ctx.llgo.New(func(cb *gogen.CodeBuilder) int {
		cb.Val(llgoFiles)
		return 1
	}, 0, token.NoPos, nil, "LLGoFiles")
	return &WrapFile{Filename: filename}
}

func wrapInlineFunc(ctx *pkgCtx, manglingName string, fn clang.Cursor, cls *classCtx, isStatic bool) string {
	first := ctx.wrap == nil
	if first {
		ctx.wrap = newWrapFile(ctx)
	}
	w := &ctx.wrap.Content
	if first {
		w.WriteString(ctx.wrapFileHeader)
	}
	w.WriteByte('\n')
	wrapName := "_llcppg_" + manglingName
	writeFunc(w, wrapName, fn, cls, isStatic, ctx.lang)
	return wrapName
}

// -----------------------------------------------------------------------------

type writerT = bytes.Buffer

func writeFunc(b *writerT, name string, fn clang.Cursor, cls *classCtx, isStatic bool, lang Language) {
	var call writerT
	writeFuncProto(b, &call, name, fn, cls, isStatic, lang)
	b.WriteString(" {\n")
	b.Write(call.Bytes())
	b.WriteString("}\n")
}

func writeFuncProto(out, call *writerT, name string, fn clang.Cursor, cls *classCtx, isStatic bool, lang Language) {
	var b writerT
	b.WriteString(name)
	b.WriteByte('(')
	call.WriteByte('\t')
	retType := fn.ResultType()
	if retType.Kind != lc.TypeVoid {
		call.WriteString("return ")
	}
	// A static method has no "this" receiver, but the call must be qualified
	// with the enclosing class name (ClassName::method).
	hasThis := cls != nil && !isStatic
	notFirst := hasThis
	if hasThis {
		b.WriteString(clang.String(cls.decl.Type()))
		b.WriteString("* this")
		call.WriteString("this->")
	} else if isStatic && cls != nil {
		call.WriteString(clang.String(cls.decl.Type()))
		call.WriteString("::")
	}
	call.WriteString(clang.String(fn))
	call.WriteByte('(')
	for i := range c.Uint(fn.NumArguments()) {
		if i > 0 {
			call.WriteString(", ")
		}
		if notFirst {
			b.WriteString(", ")
		} else {
			notFirst = true
		}
		arg := fn.Argument(i)
		argName := clang.String(arg)
		writeParam(&b, arg.Type(), argName)
		call.WriteString(argName)
	}
	b.WriteByte(')')
	call.WriteString(");\n")
	if lang == LanguageCXX {
		out.WriteString(`extern "C" `)
	}
	writeParam(out, retType, b.String())
}

func writeParam(b *writerT, typ lc.Type, name string) {
	tderef, lvl := deref(typ)
	if tderef.Kind == lc.TypeFunctionProto {
		writeFuncParam(b, tderef, lvl, name)
		return
	}
	b.WriteString(clang.String(typ))
	if typ.Kind != lc.TypePointer {
		b.WriteByte(' ')
	}
	b.WriteString(name)
}

func writeFuncParam(out *writerT, fn lc.Type, lvl int, name string) {
	var b writerT
	b.WriteByte('(')
	for range lvl {
		b.WriteByte('*')
	}
	b.WriteString(name)
	b.WriteByte(')')
	b.WriteString("(")
	for i := range c.Uint(fn.NumArgTypes()) {
		if i > 0 {
			b.WriteString(", ")
		}
		arg := fn.ArgType(i)
		writeParam(&b, arg, "")
	}
	b.WriteString(")")
	writeParam(out, fn.ResultType(), b.String())
}

func deref(typ lc.Type) (lc.Type, int) {
	n := 0
	for typ.Kind == lc.TypePointer {
		typ = typ.PointeeType()
		n++
	}
	return typ, n
}

// -----------------------------------------------------------------------------
