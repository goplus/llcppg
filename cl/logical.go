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
	"go/types"

	"github.com/goplus/llcppg/clang"
)

// -----------------------------------------------------------------------------

// Logical classes.
//
// MethodCheck can tell which logical type a C function belongs to (for example
// PyList_GetItem belongs to List even though its receiver is the physical base
// class Object). Instead of attaching every such function to the physical base
// class - which pollutes the base API surface and invites method-name conflicts
// (PyList_GetItem and PyDict_GetItem would both want to be Item) - llcppg emits
// each method on its own logical class that embeds the physical base class, plus
// a conversion method on the base class named As<Class> that reinterprets the
// pointer.
//
//	type List struct {
//		Object
//	}
//
//	func (o *Object) AsList() *List { return (*List)(unsafe.Pointer(o)) }
//
//	// llgo:link (*List).Item C.PyList_GetItem
//	func (self *List) Item(index *Object) *Object { return nil }
//
// A logical class is generated lazily: only when at least one function actually
// resolves to it. If the logical type is the physical type itself (for example
// PyObject_IsTrue resolves to Object), the method stays on Object and no class
// or AsObject is generated. See https://github.com/goplus/llcppg/issues/948.

const asMethodPrefix = "As"

type logicalClass struct {
	named *types.Named // the logical class type (List)
	base  *types.Named // the physical base class it embeds (Object)
}

// logicalClassOf returns the logical class named goName whose base is the
// physical type base, creating it (and its As<Class> conversion method on base)
// on first use. It is lazy and deterministic: the type and the conversion method
// are emitted in the order functions first resolve to the logical class.
//
// If a package-level type with the same name already exists (generated from the
// headers or provided through TypeAlias), it is reused instead of generating a
// second one, and the As<Class> conversion is still added when it is missing.
func (p *pkgCtx) logicalClassOf(decl clang.Cursor, goName string, base *types.Named) *types.Named {
	if lc, ok := p.logicals[goName]; ok {
		return lc.named
	}

	named := p.newLogicalType(decl, goName, base)
	lc := &logicalClass{named: named, base: base}
	p.logicals[goName] = lc

	// genAsMethod runs once per logical class (subsequent resolutions hit the
	// cache above). It is attempted whether the type was freshly emitted or an
	// existing package-level type was reused, so a reused type still gets its
	// As<Class> conversion; genAsMethod itself keeps any pre-existing method of
	// that name and reports a diagnostic.
	p.genAsMethod(decl, lc)
	return named
}

// newLogicalType returns the named type for the logical class, reusing an
// existing package-level type when one is present and otherwise emitting a fresh
// "type <goName> struct { <base> }".
func (p *pkgCtx) newLogicalType(decl clang.Cursor, goName string, base *types.Named) *types.Named {
	if o := p.pkg.Types.Scope().Lookup(goName); o != nil {
		if tn, ok := o.(*types.TypeName); ok {
			if n, ok := tn.Type().(*types.Named); ok {
				return n
			}
		}
	}

	pkg := p.pkg
	pkgTypes := pkg.Types
	typDecl := newType(p, decl, "", goName)
	embed := types.NewField(goNodePos(p, decl), pkgTypes, base.Obj().Name(), base, true)
	typDecl.InitType(pkg, types.NewStruct([]*types.Var{embed}, nil))
	return typDecl.Type()
}

// genAsMethod emits the conversion method on the base class, for example
//
//	func (o *Object) AsList() *List { return (*List)(unsafe.Pointer(o)) }
//
// The conversion only reinterprets the pointer; it does not call into C and does
// not check the object's real type. If a method of the same name already exists
// on the base class, the existing method is kept and a diagnostic is reported.
func (p *pkgCtx) genAsMethod(decl clang.Cursor, lc *logicalClass) {
	name := asMethodPrefix + lc.named.Obj().Name()
	if pos, _, exists := findMember(lc.base, name); exists {
		p.errorf(decl, "%s redeclared in this block\n\t%v: other declaration of %s", name, p.position(pos), name)
		return
	}

	pkg := p.pkg
	pkgTypes := pkg.Types
	recvType := types.NewPointer(lc.base)
	retType := types.NewPointer(lc.named)
	recv := types.NewParam(goNodePos(p, decl), pkgTypes, "o", recvType)
	results := types.NewTuple(types.NewParam(token.NoPos, pkgTypes, "", retType))
	sig := types.NewSignatureType(recv, nil, nil, nil, results, false)

	f, err := pkg.NewFuncWith(goNodePos(p, decl), name, sig, nil)
	if err != nil {
		p.panicf(decl, "logical class %s: genAsMethod failed - %v", lc.named.Obj().Name(), err)
	}
	cb := f.BodyStart(pkg)
	// return (*List)(unsafe.Pointer(o))
	cb.Typ(retType).
		Typ(p.unsafePointer()).VarVal("o").
		Call(1).
		Call(1).
		Return(1).End()
}

// -----------------------------------------------------------------------------
