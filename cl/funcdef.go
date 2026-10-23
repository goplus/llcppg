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
	"log"
	"strings"

	"github.com/goplus/llcppg/clang"
)

// -----------------------------------------------------------------------------

// funcDef is a parsed FuncDef entry: the target file, the resolved Go name parts
// and the C symbol it links to.
type funcDef struct {
	fname  string // target Go file name
	name   string // raw FuncDef.Name, used in error messages
	link   string // C symbol name (FuncDef.Link)
	recv   string // receiver Go type name, empty for a package-level function
	method string // method name, empty for a package-level function
	isPtr  bool   // true for a pointer receiver "(*Type).Method"
}

// loadFuncDefs registers the FuncDef entries. Each entry's Link symbol is claimed
// immediately so the default conversion, NewCheck and MethodCheck skip it. The actual
// declaration is emitted after the default functions are compiled and after the
// receiver types (from the headers or TypeDef) exist, so a FuncDef entry can detect a
// name that is already present in the package and only supplies missing symbols.
// See https://github.com/goplus/llcppg/issues/1083.
func loadFuncDefs(ctx *pkgCtx, funcDefs map[string][]FuncDef) {
	if len(funcDefs) == 0 {
		return
	}
	pkg := ctx.pkg
	for fname, defs := range funcDefs {
		for _, raw := range defs {
			def, ok := parseFuncDef(ctx, fname, raw)
			if !ok {
				continue
			}
			// Claim the C symbol up front so the default emission does not also bind
			// it under another Go name (this runs before the compile phase).
			ctx.funcDefLinks[def.link] = none{}

			pkg.SetCurFile(fname, true)
			// Defer through a load unit so the compile unit is appended after the
			// default functions' compile units: FuncDef emission then runs last and
			// can see every symbol already in the package.
			ctx.addLoadUnit(func(ctx *pkgCtx) {
				ctx.addCompileUnit(func(ctx *pkgCtx) {
					compileFuncDef(ctx, def)
				})
			})
		}
	}
}

// parseFuncDef validates a raw FuncDef entry and splits its Name into the three
// accepted forms: "Ident", "(*Type).Method" or "(Type).Method".
func parseFuncDef(ctx *pkgCtx, fname string, raw FuncDef) (def funcDef, ok bool) {
	def = funcDef{fname: fname, name: raw.Name, link: raw.Link}
	if raw.Name == "" {
		ctx.npErrorf("FuncDef %q: Name is required", fname)
		return def, false
	}
	if raw.Link == "" {
		ctx.npErrorf("FuncDef %q: entry %q: Link is required", fname, raw.Name)
		return def, false
	}
	if raw.Name[0] == '(' {
		// (*Type).Method or (Type).Method
		end := strings.IndexByte(raw.Name, ')')
		if end < 0 {
			ctx.npErrorf("FuncDef %q: entry %q: invalid Name", fname, raw.Name)
			return def, false
		}
		recv := raw.Name[1:end]
		rest, found := strings.CutPrefix(raw.Name[end+1:], ".")
		if !found || rest == "" {
			ctx.npErrorf("FuncDef %q: entry %q: invalid Name", fname, raw.Name)
			return def, false
		}
		if recv != "" && recv[0] == '*' {
			def.isPtr = true
			recv = recv[1:]
		}
		if !isGoIdent(recv) || !isGoIdent(rest) {
			ctx.npErrorf("FuncDef %q: entry %q: invalid Go identifier in Name", fname, raw.Name)
			return def, false
		}
		def.recv = recv
		def.method = rest
	} else {
		if !isGoIdent(raw.Name) {
			ctx.npErrorf("FuncDef %q: entry %q: invalid Go identifier in Name", fname, raw.Name)
			return def, false
		}
	}
	return def, true
}

// compileFuncDef emits the Go declaration for a FuncDef entry: a package-level
// function with a //go:linkname directive, or a method with a receiver whose type
// is taken from Name. The signature is derived from the C declaration of def.link.
func compileFuncDef(ctx *pkgCtx, def funcDef) {
	if debugCompileDecl {
		log.Println("==> funcDef", def.name, "link", def.link)
	}
	fnCur, ok := ctx.lookupFuncDecl(def.link)
	if !ok {
		ctx.npErrorf("FuncDef %q: entry %q: Link symbol %q not found or is not a function", def.fname, def.name, def.link)
		return
	}

	pkg := ctx.pkg
	pkgTypes := pkg.Types

	manglingName := clang.Mangling(fnCur)
	if manglingName == "" {
		ctx.npErrorf("FuncDef %q: entry %q: Link symbol %q has no mangled symbol", def.fname, def.name, def.link)
		return
	}

	fn := clangFunc{Cursor: fnCur}
	feats := 0
	params, variadic := newParamsEx(ctx, pkgTypes, fn, &feats, nil)
	results := toFuncResults(ctx, pkgTypes, fn.ResultType(), &feats, nil)
	if feats&featAllIgnore != 0 {
		ctx.npErrorf("FuncDef %q: entry %q: Link symbol %q has unsupported features", def.fname, def.name, def.link)
		return
	}

	if fnCur.IsFunctionInlined() != 0 {
		if ctx.cflags == "" {
			ctx.npErrorf("FuncDef %q: entry %q: Link symbol %q is an inline function but no CFlags in config", def.fname, def.name, def.link)
			return
		}
		manglingName = wrapInlineFunc(ctx, manglingName, fn, nil)
	}

	var recv *types.Var
	var fnName, nameInPkg string
	if def.recv == "" {
		// package-level function
		fnName = def.name
		if o := pkgTypes.Scope().Lookup(fnName); o != nil {
			ctx.npErrorf("FuncDef %q: entry %q: %s already exists in the package", def.fname, def.name, fnName)
			return
		}
		nameInPkg = fnName
	} else {
		recvNamed, ok := funcDefRecv(ctx, pkgTypes, def)
		if !ok {
			return
		}
		if len(params) == 0 {
			ctx.npErrorf("FuncDef %q: entry %q: Link symbol %q has no parameters, cannot be a method", def.fname, def.name, def.link)
			return
		}
		if !funcDefRecvCompatible(params[0].Type(), recvNamed, def.isPtr) {
			ctx.npErrorf("FuncDef %q: entry %q: receiver type %s is not compatible with the first parameter of %q", def.fname, def.name, def.recv, def.link)
			return
		}
		if pos, _, exists := findMember(recvNamed, def.method); exists {
			ctx.npErrorf("FuncDef %q: entry %q: method %s already exists on %s\n\t%v: other declaration", def.fname, def.name, def.method, def.recv, ctx.position(pos))
			return
		}
		var recvType types.Type = recvNamed
		if def.isPtr {
			recvType = types.NewPointer(recvNamed)
		}
		recv = types.NewParam(token.NoPos, pkgTypes, c2goMethodRecvName, recvType)
		params = params[1:]
		fnName = def.method
		if def.isPtr {
			nameInPkg = "(*" + def.recv + ")." + fnName
		} else {
			nameInPkg = def.recv + "." + fnName
		}
	}

	sig := types.NewSignatureType(recv, nil, nil, types.NewTuple(params...), results, variadic)
	f, err := pkg.NewFuncWith(fn.NodePos(ctx), fnName, sig, nil)
	if err != nil {
		ctx.npErrorf("FuncDef %q: entry %q: %v", def.fname, def.name, err)
		return
	}

	if recv == nil {
		ctx.forceImportUnsafe()
		f.SetComments(pkg, ctx.directiveCommentsEx(fn,
			"\n//go:linkname "+nameInPkg+" C."+manglingName))
	} else {
		f.SetComments(pkg, ctx.directiveCommentsEx(fn,
			"\n// llgo:link "+nameInPkg+" C."+manglingName))
		cb := f.BodyStart(pkg)
		if n := results.Len(); n > 0 {
			retType := results.At(0).Type()
			if types.Identical(recv.Type(), retType) {
				cb.Val(recv)
			} else {
				cb.ZeroLit(retType)
			}
			cb.Return(1)
		}
		cb.End()
	}
}

// funcDefRecv resolves the receiver type named by a FuncDef method entry. The type
// must exist in the generated package.
func funcDefRecv(ctx *pkgCtx, pkgTypes *types.Package, def funcDef) (*types.Named, bool) {
	o := pkgTypes.Scope().Lookup(def.recv)
	if o == nil {
		ctx.npErrorf("FuncDef %q: entry %q: receiver type %s does not exist in the package", def.fname, def.name, def.recv)
		return nil, false
	}
	tn, ok := o.(*types.TypeName)
	if !ok {
		ctx.npErrorf("FuncDef %q: entry %q: receiver %s is not a type", def.fname, def.name, def.recv)
		return nil, false
	}
	named, ok := types.Unalias(tn.Type()).(*types.Named)
	if !ok {
		ctx.npErrorf("FuncDef %q: entry %q: receiver type %s cannot have methods", def.fname, def.name, def.recv)
		return nil, false
	}
	return named, true
}

// funcDefRecvCompatible reports whether the receiver type is compatible with the C
// function's first parameter type: the receiver is the first parameter type itself,
// or a type that embeds it as its first field. A pointer receiver is matched against
// a pointer first parameter, and a value receiver against a value first parameter.
func funcDefRecvCompatible(first types.Type, recv *types.Named, isPtr bool) bool {
	if isPtr {
		p, ok := first.(*types.Pointer)
		if !ok {
			return false
		}
		first = p.Elem()
	} else if _, ok := first.(*types.Pointer); ok {
		return false
	}
	base, ok := types.Unalias(first).(*types.Named)
	if !ok {
		return false
	}
	return isTypeFromBase(recv, base)
}

// isGoIdent reports whether s is a valid Go identifier.
func isGoIdent(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '_':
		case 'a' <= c && c <= 'z', 'A' <= c && c <= 'Z':
		case '0' <= c && c <= '9':
			if i == 0 {
				return false
			}
		default:
			return false
		}
	}
	return !token.Lookup(s).IsKeyword()
}

// -----------------------------------------------------------------------------
