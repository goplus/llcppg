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
	"strconv"
	"strings"

	"github.com/goplus/llcppg/clang"
	lc "github.com/llarhub/clang-c"
)

// -----------------------------------------------------------------------------

// field = "fieldName fieldType"
// fieldType:
//   - "pkg.Object"
//   - ".Object"
//   - "*pkg.Object"
//   - "*.Object"
//   - "[8].X__SharedPtrEmplace_Storage_Data[_Tp, _Alloc]"
func parseField(ctx *pkgCtx, field string, scope *scopeCtx) (fldName string, fldType types.Type, embedded, found bool) {
	fldTy := ""
	pos := strings.IndexAny(field, " \t")
	embedded = pos < 0
	if embedded {
		fldName = field
		fldTy = field
	} else {
		fldName = field[:pos]
		fldTy = strings.TrimLeft(field[pos+1:], " \t")
	}
	fldType, found = parseType(ctx, fldTy, scope)
	if found && embedded {
		t := fldType
		if tp, ok := t.(*types.Pointer); ok {
			t = tp.Elem()
		}
		switch t := t.(type) {
		case *types.Named:
			fldName = t.Obj().Name()
		case *types.Alias:
			fldName = t.Obj().Name()
		case *types.Basic:
			fldName = t.Name()
		default:
			found = false
		}
	}
	return
}

func parseType(ctx *pkgCtx, typ string, scope *scopeCtx) (types.Type, bool) {
	if typ == "" {
		return nil, false
	}
	switch typ[0] {
	case '*':
		if elem, ok := parseType(ctx, typ[1:], scope); ok {
			return types.NewPointer(elem), true
		}
	case '[':
		typ = typ[1:]
		if end := strings.IndexByte(typ, ']'); end >= 0 {
			idx := typ[:end]
			if elem, ok := parseType(ctx, typ[end+1:], scope); ok {
				if n, err := strconv.ParseInt(idx, 10, 64); err == nil {
					return types.NewArray(elem, n), true
				}
			}
		}
	default:
		return ctx.goTypeOf(typ, scope)
	}
	return nil, false
}

// -----------------------------------------------------------------------------

func loadVar(ctx *pkgCtx, decl clang.Cursor) {
	ctx.addCompileUnit(func(ctx *pkgCtx) {
		compileVar(ctx, decl)
	})
}

func compileVar(ctx *pkgCtx, decl clang.Cursor) {
	ns := cNS(decl)
	cName := cNameWithNS(clang.String(decl), ns)
	if feats := ctx.nsFeats(ns); feats&featAllIgnore != 0 {
		ctx.ignoref(feats, decl, "var %s: parent is ignored", cName)
		return
	}
	if varHasInitExpr(ctx, decl) {
		ctx.ignoref(featQuietIgnore, decl, "var %s: has initialized expression, ignored", cName)
		return
	}

	manglingName := clang.Mangling(decl)
	if _, ok := ctx.nameLookup(manglingName); !ok {
		ctx.ignoref(featExplicitIgnore, decl, "var %s: symbol not found in lib files, ignored", cName)
		return
	}

	if debugCompileDecl {
		ctx.logf(decl, "var %s: %s", cName, clang.String(decl.Type()))
	}

	pkg := ctx.pkg
	pkgTypes := pkg.Types

	feats := 0
	typ := toTypeEx(ctx, pkgTypes, decl.Type(), flagIsVarDef, &feats, nil)
	if feats&featAllIgnore != 0 {
		ctx.ignoref(feats, decl, "var %s: unsupported type, ignored", cName)
		return
	}

	goName := ctx.varName(cName)
	ctx.forceImportUnsafe()
	defs := pkg.NewVarDefs(pkgTypes.Scope()).SetComments(
		ctx.directiveComments(decl, "\n//go:linkname "+goName+" C."+manglingName))
	defs.New(goNodePos(ctx, decl), typ, goName)
}

func varHasInitExpr(ctx *pkgCtx, v clang.Cursor) (hasInitExpr bool) {
	clang.VisitChildren(v, func(decl, parent clang.Cursor) clang.ChildVisitResult {
		switch decl.Kind {
		case lc.Cursor_InitListExpr, lc.Cursor_CallExpr, lc.Cursor_DeclRefExpr,
			lc.Cursor_UnaryOperator, lc.Cursor_BinaryOperator, lc.Cursor_UnexposedExpr,
			lc.Cursor_CXXStaticCastExpr, lc.Cursor_CXXBoolLiteralExpr,
			lc.Cursor_StringLiteral, lc.Cursor_FloatingLiteral, lc.Cursor_CharacterLiteral:
			// An initializer expression (an initialized global, which the
			// generator intentionally omits). StringLiteral covers
			// `char text[] = "hi"`, whose string-completed array would
			// otherwise reach the panic below.
			hasInitExpr = true
			return clang.Break
		case lc.Cursor_TypeRef, lc.Cursor_UnaryExpr, lc.Cursor_ParmDecl,
			lc.Cursor_IntegerLiteral:
		default:
			if lc.Cursor_FirstAttr > decl.Kind || decl.Kind > lc.Cursor_LastAttr {
				ctx.panicf(decl, "varHasInitExpr: unknown kind - %v", decl.Kind)
			}
		}
		return clang.Continue
	})
	return
}

// -----------------------------------------------------------------------------
