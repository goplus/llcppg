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

func funcUSR(decl clang.Cursor) string {
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
		if decl.Kind == lc.Cursor_TranslationUnit {
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

func cBaseName(decl clang.Cursor) string {
	return trimTypeTag(clang.String(decl))
}

func cTypeName(typ lc.Type) string {
	return cNameOf(typ.Declaration())
}

// -----------------------------------------------------------------------------

func (p *pkgCtx) globalName(name string, trimPrefixs []string) string {
	if v, ok := p.rename[name]; ok {
		return v // special case
	}
	return p.cstyleToGo(rmPrefix(name, trimPrefixs), true)
}

func (p *pkgCtx) fieldName(name string, public bool) string {
	return p.cstyleToGo(name, public)
}

func (p *pkgCtx) varName(name string) string {
	return p.globalName(name, p.varPrefix)
}

func (p *pkgCtx) macroName(name string) string {
	return p.globalName(name, p.macroPrefix)
}

func (p *pkgCtx) enumvalName(name, ns string) string {
	if ns == "" {
		return p.globalName(name, p.enumPrefix)
	}
	ns = p.globalName(ns, p.nsPrefix)
	name = p.cstyleToGo(name, false)
	if ns != "" {
		name = ns + "_" + name
	}
	if v, ok := p.rename[name]; ok {
		return v // special case
	}
	return name
}

func (p *pkgCtx) typeName(cName string, _ bool) string {
	if v, ok := p.rename[cName]; ok {
		return v // special case
	}
	cName = rmPrefix(cName, p.nsPrefix)
	cName = rmPrefix(cName, p.typePrefix)
	cName = rmSuffix(cName, p.typeSuffix)
	return p.cstyleToGo(cName, true)
}

func (p *pkgCtx) funcName(name string, order int, typName, typCName string, global, _ bool) string {
	if v, ok := p.rename[name]; ok {
		return v // special case
	}
	name = rmPrefix(name, p.nsPrefix)
	if global {
		name = rmPrefix(name, p.fnPrefix)
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
		// don't remove type name suffix for a method
		typName = ""
	}
	if !strings.HasPrefix(name, "XGo_") { // avoid rewriting XGo_xxx names
		name = p.cstyleToGo(name, true)
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
				typName = rmSuffix(typName, p.typeAbbrSuffix)
				typSuffix = []string{typName}
			}
			name = rmSuffix(name, typSuffix)
			name = cutMethodPrefix(name, typName)
		}
	}
	if order >= 0 {
		name = name + "__" + strconv.FormatInt(int64(order), 36)
	}
	return name
}

func (p *pkgCtx) cstyleToGo(cName string, public bool) string {
	rename := p.rename
	parts, hasNS := cNameSplit(cName)
	if !hasNS && isAllUpperStart(parts) {
		if parts[0] == "" && public {
			return "X" + cName
		}
		return cName
	}
	lastEndWithUpper := false
	for i := 0; i < len(parts); i++ {
		part := parts[i]
		if part == "" {
			if i == 0 && public {
				parts[i] = "X_"
				i++ // skip next part
			} else {
				parts[i] = "_"
			}
			continue
		}
		if v, ok := rename[part]; ok && v != "" {
			part = v
		} else if i > 0 || public {
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
	return strings.Join(parts, "")
}

// -----------------------------------------------------------------------------

func isAllUpperStart(parts []string) bool {
	for _, part := range parts {
		if part != "" {
			if r := part[0]; 'a' <= r && r <= 'z' {
				return false
			}
		}
	}
	return true
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

func rmPrefix(name string, prefix []string) string {
	for _, pfx := range prefix {
		if strings.HasPrefix(name, pfx) {
			return name[len(pfx):]
		}
	}
	return name
}

func rmSuffix(name string, suffix []string) string {
	for _, sfx := range suffix {
		if strings.HasSuffix(name, sfx) {
			return name[:len(name)-len(sfx)]
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
