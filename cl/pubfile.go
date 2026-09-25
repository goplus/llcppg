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
			kind := flds[0][0]
			cName := flds[1]
			goName := ""
			if kind == 'T' && len(flds) > 2 {
				igo := 2
				if isTypeTag(flds[1]) {
					cName = flds[1] + " " + flds[2]
					igo = 3
				}
				if igo < len(flds) {
					goName = flds[igo]
				}
			} else if len(flds) > 3 {
				tooFewOrManyFields(i, "many", line)
			}
			if !yield(Entry{Kind: kind, Name: cName, GoName: goName}) {
				return
			}
		}
	}
	return
}

func isTypeTag(tag string) bool {
	switch tag {
	case "enum", "struct", "union", "class":
		return true
	}
	return false
}

func tooFewOrManyFields(i int, fewOrMany, line string) {
	panic(fmt.Errorf("line %d: too %s fields - %s", i+1, fewOrMany, line))
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
