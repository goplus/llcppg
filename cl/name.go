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
	"fmt"
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

type mthdCheck struct {
	prefix  string
	pattern []string
	result  string
}

func newMethodCheck(check string) (*mthdCheck, error) {
	pos := strings.IndexByte(check, ':')
	if pos < 0 {
		return nil, fmt.Errorf("invalid MethodCheck: %s", check)
	}
	pattern := strings.TrimSpace(check[:pos])
	parts := strings.Split(pattern, "*")
	result := strings.TrimSpace(check[pos+1:])
	return &mthdCheck{
		prefix:  parts[0],
		pattern: parts[1:],
		result:  result,
	}, nil
}

func newMethodChecks(check []string) ([]*mthdCheck, error) {
	ret := make([]*mthdCheck, len(check))
	for i, c := range check {
		m, err := newMethodCheck(c)
		if err != nil {
			return nil, err
		}
		ret[i] = m
	}
	return ret, nil
}

func (p *mthdCheck) check(name string) (cls string, ok bool) {
	name, ok = strings.CutPrefix(name, p.prefix)
	if !ok {
		return
	}
	n := 0
	match := make([]string, len(p.pattern))
	for i, p := range p.pattern {
		pos := strings.Index(name, p)
		if pos < 0 {
			return
		}
		n += pos
		match[i] = name[:pos]
		name = name[pos+len(p):]
	}
	cls = matchResult(p.result, match, n)
	ok = cls != ""
	return
}

func matchResult(result string, match []string, n int) string {
	b := make([]byte, 0, len(result)+n)
	for i := 0; i < len(result); i++ {
		if result[i] == '$' {
			if i+1 < len(result) {
				i++
				c := result[i]
				if c >= '1' && c <= '9' {
					if index := int(c - '1'); index < len(match) {
						b = append(b, match[index]...)
						continue
					}
				}
				b = append(b, '$', c)
				continue
			}
		}
		b = append(b, result[i])
	}
	return string(b)
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

func (p *pkgCtx) methodCheck(cName string) (cls string, mayClass bool) {
	if len(p.mthdCheck) == 0 {
		return "", true
	}
	for _, m := range p.mthdCheck {
		if cls, mayClass = m.check(cName); mayClass {
			return
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
			if v, ok := p.typeAbbr[typName]; ok { // Go type name => abbreviated name
				switch v := v.(type) {
				case string:
					typName = v
					typSuffix = []string{v}
				case []any:
					if n := len(v); n > 0 {
						abbrs := make([]string, n)
						for i, v := range v {
							abbrs[i] = v.(string)
						}
						typName = abbrs[n-1]
						typSuffix = abbrs
					}
				default:
					panic(fmt.Errorf("invalid TypeAbbr for %q: %v", typName, v))
				}
			} else {
				typName = cutSuffixes(typName, p.typeAbbrSuffix)
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
