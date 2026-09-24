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
	"iter"
	"os"
	"sort"
	"strings"
)

// -----------------------------------------------------------------------------

// loadPubFile loads a public file and returns an iterator that yields each public entry.
//
// T cName
// T cName goName
// T enum cName
// T enum cName goName
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
			if len(flds) == 0 {
				continue
			}
			kind := flds[0][0]
			goName := ""
			switch len(flds) {
			case 2: // T cName
			case 3:
				if kind == 'T' && flds[1] == "enum" {
					// T enum cName
					flds[1] = "enum " + flds[2]
				} else {
					// T cName goName
					goName = flds[2]
				}
			case 4:
				if kind == 'T' && flds[1] == "enum" {
					// T enum cName goName
					flds[1] = "enum " + flds[2]
					goName = flds[3]
				} else {
					err = fmt.Errorf("line %d: too few/many fields - %s\n", i+1, line)
					return
				}
			default:
				err = fmt.Errorf("line %d: too few/many fields - %s\n", i+1, line)
				return
			}
			if !yield(Entry{Kind: kind, Name: flds[1], GoName: goName}) {
				return
			}
		}
	}
	return
}

// -----------------------------------------------------------------------------

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

// -----------------------------------------------------------------------------
