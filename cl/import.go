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
	"go/types"
	"iter"
	"log"
	"os"
	"strings"
)

// -----------------------------------------------------------------------------

type entryKind = byte

const (
	entryType       entryKind = 'T' // type
	entryFunc       entryKind = 'f' // func
	entryVar        entryKind = 'v' // var
	entryDependency entryKind = 'D' // dependency
)

// Entry represents a C/C++ name and its corresponding Go name.
type Entry struct {
	Name   string
	GoName string    // optional
	Tag    typeTag   // only valid for Kind == 'T'
	Kind   entryKind // 'T' (type), 'f' (func), 'v' (var), 'D' (dependency)
}

// -----------------------------------------------------------------------------

// loadPubFile loads a public file and returns an iterator that yields each public entry.
//
// D depPkgPath
// T cName
// T cName goName
// T <tag> cName
// T <tag> cName goName
func loadPubFile(pubfile string) (it iter.Seq[Entry], err error) {
	b, err := os.ReadFile(pubfile)
	if err != nil {
		return
	}

	text := string(b)
	lines := strings.Split(text, "\n")
	it = func(yield func(Entry) bool) {
		for i, line := range lines {
			flds := strings.Fields(line)
			switch len(flds) {
			case 0:
				continue
			case 1:
				tooFewOrManyFields(i, "few", line)
			}
			kindFld := flds[0]
			if len(kindFld) != 1 {
				panic(fmt.Errorf("line %d: invalid kind - %s", i+1, kindFld))
			}
			kind := kindFld[0]
			cName := flds[1]
			goName := ""
			tag := typeTag(0)
			if kind == entryType && len(flds) > 2 {
				igo := 2
				if tag = getTypeTag(flds[1]); tag != 0 {
					cName, igo = flds[2], 3
				}
				if igo < len(flds) {
					goName = flds[igo]
				}
			} else if len(flds) > 3 {
				tooFewOrManyFields(i, "many", line)
			}
			if !yield(Entry{Kind: kind, Name: cName, GoName: goName, Tag: tag}) {
				return
			}
		}
	}
	return
}

func getTypeTag(tag string) typeTag {
	switch tag {
	case "enum":
		return tagEnum
	case "struct":
		return tagStruct
	case "class":
		return tagClass
	case "union":
		return tagUnion
	}
	return 0
}

func tooFewOrManyFields(i int, fewOrMany, line string) {
	panic(fmt.Errorf("line %d: too %s fields - %s", i+1, fewOrMany, line))
}

// -----------------------------------------------------------------------------
/*
func savePubFile(file string, it iter.Seq[Entry], n int) (err error) {
	if n == 0 {
		return
	}
	f, err := os.Create(file)
	if err != nil {
		return
	}
	defer f.Close()
	ret := make([]string, 0, n)
	for e := range it {
		line := string(rune(e.Kind)) + " " + e.Name
		if e.GoName != "" {
			line += " " + e.GoName
		}
		ret = append(ret, line)
	}
	sort.Strings(ret)
	_, err = f.WriteString(strings.Join(ret, "\n"))
	return
}
*/
// -----------------------------------------------------------------------------

func (p *pkgCtx) forceImportUnsafe() {
	p.pkg.ForceImport("unsafe")
}

func (p *pkgCtx) importPkg(pkgPath string) {
	if _, ok := p.impPkgs[pkgPath]; ok {
		return // imported already
	}
	p.impPkgs[pkgPath] = none{}

	pubFile, ok := p.pubLookup(pkgPath)
	if !ok {
		log.Panicln("[ERROR] pubFile not found for", pkgPath)
	}

	if debugCompileDecl {
		log.Println("==> importPkg", pkgPath)
	}

	pkg := p.pkg.Import(pkgPath)
	entries, err := loadPubFile(pubFile)
	if err != nil {
		if os.IsNotExist(err) {
			return // ignore missing pub file
		}
		log.Panicln("[ERROR] loadPubFile failed:", err)
	}

	scope := pkg.Types.Scope()
	for e := range entries {
		switch e.Kind {
		case entryDependency:
			p.importPkg(e.Name)

		case entryType:
			if e.GoName == "" {
				e.GoName = p.typeName(strings.ReplaceAll(e.Name, "::", "_"), true)
			}
			if o := scope.Lookup(e.GoName); o != nil {
				if t, ok := o.(*types.TypeName); ok {
					p.types[e.Name] = typeObj{t, 0}
					if e.Tag != 0 {
						p.types[tagStrvals[e.Tag]+e.Name] = typeObj{t, 0}
					}
				}
			}

		default:
			panic("importPkg: unsupport - " + e.Name)
		}
	}
}

// -----------------------------------------------------------------------------
