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

const (
	anonPrefix = "_llcppg_anon_"
)

func (p *pkgCtx) nextAnonName() string {
	name := anonPrefix + strconv.Itoa(p.anonSeq)
	p.anonSeq++
	return name
}

func objUSR(decl clang.Cursor) string {
	return clang.USR(decl)
}

// -----------------------------------------------------------------------------

func cNameSplit(cName string) (parts []string, hasNS bool) {
	for {
		pos := strings.IndexAny(cName, "_:")
		if pos < 0 {
			parts = append(parts, cName)
			return
		}
		parts = append(parts, cName[:pos])
		if cName[pos] == ':' && len(cName) > pos+1 && cName[pos+1] == ':' {
			hasNS = true
			cName = cName[pos+2:]
		} else {
			cName = cName[pos+1:]
		}
	}
}

func cTypeName(typ lc.Type) string {
	return cNameOf(typ.Declaration())
}

func cBaseName(decl clang.Cursor) string {
	return trimTypeTag(clang.String(decl))
}

func cNameWithNS(name, ns string) string {
	if ns == "" {
		return name
	}
	return ns + "::" + name
}

func cNameOf(decl clang.Cursor) string {
	cName := cBaseName(decl)
	for {
		decl = decl.SemanticParent()
		if decl.IsNull() != 0 || decl.Kind == lc.Cursor_TranslationUnit {
			break
		}
		cName = cBaseName(decl) + "::" + cName
	}
	return cName
}

func cNS(decl clang.Cursor) (ns string) {
	for {
		decl = decl.SemanticParent()
		if decl.IsNull() != 0 || decl.Kind == lc.Cursor_TranslationUnit {
			return
		}
		name := cBaseName(decl)
		if ns == "" {
			ns = name
		} else {
			ns = name + "::" + ns
		}
	}
}

func (p *pkgCtx) cNS(decl clang.Cursor) (ns string) {
	if p.lang != LanguageC {
		ns = cNS(decl) // C doesn't have namespaces
	}
	return
}

// -----------------------------------------------------------------------------

func (p *pkgCtx) globalName(name string, trimPrefixs []string) string {
	if v, ok := p.rename[name]; ok {
		return v // special case
	}
	name, underscoreStart := rmPrefixAndUnderscoreStart(name, trimPrefixs)
	return p.cstyleToGo(name, underscoreStart, true)
}

func (p *pkgCtx) localName(cName string, public bool) string {
	name, underscoreStart := checkUnderscoreStart(cName)
	return p.cstyleToGo(name, underscoreStart, public)
}

func (p *pkgCtx) fieldName(name string, public bool) string {
	if name == "" {
		return "_"
	}
	return p.localName(name, public)
}

func (p *pkgCtx) macroName(name string) string {
	return p.globalName(name, p.macroPrefix)
}

func (p *pkgCtx) enumvalName(name, ns string) string {
	if ns == "" {
		return p.globalName(name, p.enumPrefix)
	}
	ns = p.globalName(ns, p.nsPrefix)
	name = p.localName(name, false)
	if ns != "" {
		name = ns + "_" + name
	}
	if v, ok := p.rename[name]; ok {
		return v // special case
	}
	return name
}

func (p *pkgCtx) varName(cName string) string {
	if v, ok := p.rename[cName]; ok {
		return v // special case
	}
	name := rmPrefix(cName, p.nsPrefix)
	name, underscoreStart := rmPrefixAndUnderscoreStart(name, p.varPrefix)
	return p.cstyleToGo(name, underscoreStart, true)
}

func (p *pkgCtx) typeName(cName string, _ bool) string {
	if v, ok := p.rename[cName]; ok {
		return v // special case
	}
	name := rmPrefix(cName, p.nsPrefix)
	name, underscoreStart := rmPrefixAndUnderscoreStart(name, p.typePrefix)
	name = cutSuffixes(name, p.typeSuffix)
	return p.cstyleToGo(name, underscoreStart, true)
}

// mayCreator sets true even if cls is not found, indicating that a creator might exist.
func (p *pkgCtx) creatorCheck(pkgTypes *types.Package, cName string) (creator string, obj types.Object, mayMethod, mayCreator bool) {
	rule, mayCreator := p.newCheck.match(cName, true)
	if !mayCreator {
		return
	}
	pos := strings.IndexByte(rule, ' ')
	if pos <= 0 {
		return
	}
	scope := pkgTypes.Scope()
	cls := rule[:pos]
	obj = scope.Lookup(cls)
	if obj != nil {
		creator = strings.TrimLeft(rule[pos+1:], " \t")
		creator, mayMethod = strings.CutPrefix(creator, ".")
	}
	return
}

func (p *pkgCtx) methodCheck(pkgTypes *types.Package, cName string) (cls string, obj types.Object, mayClass bool) {
	if len(p.mthdCheck) == 0 {
		return "", nil, true // allow tryToMethod
	}
	classes, _ := p.mthdCheck.matchList(cName, false)
	for _, class := range classes {
		if o := pkgTypes.Scope().Lookup(class); o != nil {
			return class, o, true
		}
	}
	return
}

func (p *pkgCtx) funcName(cName string, order int, typName, typCName string, global, _ bool) string {
	if v, ok := p.rename[cName]; ok {
		return v // special case
	}
	name := rmPrefix(cName, p.nsPrefix)
	underscoreStart := false
	if global {
		name, underscoreStart = rmPrefixAndUnderscoreStart(name, p.fnPrefix)
		if typCName != "" {
			// remove typCName prefix & suffix
			if before, ok := strings.CutSuffix(name, typCName); ok {
				name = strings.TrimSuffix(before, "_")
			} else if after, ok := strings.CutPrefix(name, typCName+"_"); ok {
				name = after
			} else {
				name = strings.TrimPrefix(name, typName+"_")
			}
		}
	} else {
		name, underscoreStart = checkUnderscoreStart(name)
		typName = "" // don't remove type name suffix for a method
	}
	if underscoreStart || !strings.HasPrefix(name, "XGo_") { // avoid rewriting XGo_xxx names
		name = p.cstyleToGo(name, underscoreStart, true)
		if typName != "" {
			var typSuffix []string
			if abbrs, _ := p.typeAbbr.matchList(typName, true); len(abbrs) > 0 {
				typName = abbrs[len(abbrs)-1]
				typSuffix = abbrs
			} else {
				typSuffix = []string{typName}
			}
			name = cutSuffixes(name, typSuffix)
			name = cutMethodPrefix(name, typName)
		}
	}
	if order >= 0 {
		name = name + "__" + strconv.FormatInt(int64(order), 36)
	}
	return name
}

func (p *pkgCtx) cstyleToGo(name string, underscoreStart, public bool) string {
	rename := p.rename
	parts, hasNS := cNameSplit(name)
	if p.shouldKeepCStyle(parts, hasNS) {
		return goNameOf(name, underscoreStart, public)
	}
	lastEndWithUpper := false
	for i := 0; i < len(parts); i++ {
		part := parts[i]
		if part == "" {
			parts[i] = "_"
			continue
		}
		if v, ok := rename[part]; ok && v != "" {
			part = v
		} else if i > 0 || public && !underscoreStart {
			if c := part[0]; 'a' <= c && c <= 'z' {
				c -= 'a' - 'A'
				part = string(c) + part[1:]
			}
		}
		c := part[len(part)-1]
		endWithUpper := 'A' <= c && c <= 'Z'
		if lastEndWithUpper && endWithUpper {
			part = "_" + part
		}
		lastEndWithUpper = endWithUpper
		parts[i] = part
	}
	return goNameOf(strings.Join(parts, ""), underscoreStart, public)
}

func (p *pkgCtx) shouldKeepCStyle(parts []string, hasNS bool) bool {
	if !hasNS && isAllUpperStart(parts) {
		if p.forceCamelCase {
			return isAnyUpperEnd(parts)
		}
		return true
	}
	return false
}

// -----------------------------------------------------------------------------

func isAllUpperStart(parts []string) bool {
	for _, part := range parts {
		if part != "" {
			if c := part[0]; 'a' <= c && c <= 'z' {
				return false
			}
		}
	}
	return true
}

func isAnyUpperEnd(parts []string) bool {
	for _, part := range parts {
		if part != "" {
			if c := part[len(part)-1]; 'A' <= c && c <= 'Z' {
				return true
			}
		}
	}
	return false
}

func goNameOf(name string, underscoreStart, public bool) string {
	if underscoreStart {
		if public {
			return "X_" + name
		}
		return "_" + name
	}
	return name
}

func cutMethodPrefix(name, objName string) string {
	name = cutPrefix(name, objName)
	name = cutPrefix(name, "Get")
	return cutPrefix(name, objName)
}

func cutPrefix(name, prefix string) string {
	after, ok := strings.CutPrefix(name, prefix)
	if ok && after != "" {
		if c := after[0]; 'A' <= c && c <= 'Z' {
			return after
		}
	}
	return name
}

func cutSuffixes(name string, suffix []string) string {
	for _, sfx := range suffix {
		if strings.HasSuffix(name, sfx) {
			if name == sfx {
				return name
			}
			return name[:len(name)-len(sfx)]
		}
	}
	return name
}

func rmPrefixAndUnderscoreStart(cName string, prefix []string) (name string, underscoreStart bool) {
	name, underscoreStart = checkUnderscoreStart(cName)
	name = rmPrefix(name, prefix)
	return
}

func checkUnderscoreStart(cName string) (name string, underscoreStart bool) {
	name = cName
	if len(name) > 0 && name[0] == '_' {
		name = name[1:]
		underscoreStart = true
	}
	return
}

func rmPrefix(name string, prefix []string) string {
	for _, pfx := range prefix {
		if strings.HasPrefix(name, pfx) {
			return name[len(pfx):]
		}
	}
	return name
}

func contains(v string, names []string) bool {
	for _, name := range names {
		if name == v {
			return true
		}
	}
	return false
}

// -----------------------------------------------------------------------------
