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
	"log"

	"github.com/goplus/llcppg/clang"
)

// -----------------------------------------------------------------------------

func loadTypedef(ctx *pkgCtx, decl clang.Cursor, ns string) {
	pkg := ctx.pkg
	pkgTypes := pkg.Types
	origName := nameWithNS(clang.String(decl), ns)
	underlying := decl.TypedefDeclUnderlyingType()
	if debugCompileDecl {
		log.Println("typedef", origName, "-", clang.String(underlying))
	}
	name := ctx.typeName(origName, true)
	tunder := toType(ctx, pkgTypes, underlying, flagIsTypeDef)
	if tn, ok := tunder.(*types.Named); ok {
		if o := tn.Obj(); o.Pkg() == pkgTypes && o.Name() == name {
			return // already defined
		}
	}
	var obj *types.TypeName
	var typDefs = pkg.NewTypeDefs()
	var cName = clang.String(decl.Type())
	if _, ok := ctx.classes[cName]; ok {
		if tunder == types.Typ[types.UnsafePointer] {
			tunder = types.Typ[types.Uintptr] // unsafe.Pointer => uintptr
		}
		t := typDefs.NewType(name).InitType(pkg, tunder)
		obj = t.Obj()
	} else {
		t := typDefs.AliasType(name, tunder).(*types.Alias)
		obj = t.Obj()
	}
	ctx.types[cName] = obj
}

// -----------------------------------------------------------------------------
