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
func loadPubFile(pubfile string) (it iter.Seq2[string, string], n int, err error) {
	b, err := os.ReadFile(pubfile)
	if err != nil {
		return
	}

	text := string(b)
	lines := strings.Split(text, "\n")
	it = func(yield func(cName, goName string) bool) {
		for i, line := range lines {
			flds := strings.Fields(line)
			goName := ""
			switch len(flds) {
			case 1:
			case 2:
				goName = flds[1]
			case 0:
				continue
			default:
				err = fmt.Errorf("line %d: too many fields - %s\n", i+1, line)
				return
			}
			n++
			if !yield(flds[0], goName) {
				return
			}
		}
	}
	return
}

// -----------------------------------------------------------------------------

func savePubFile(file string, it iter.Seq2[string, string], n int) (err error) {
	if n == 0 {
		return
	}
	f, err := os.Create(file)
	if err != nil {
		return
	}
	defer f.Close()
	ret := make([]string, 0, n)
	for name, goName := range it {
		if goName == "" {
			ret = append(ret, name)
		} else {
			ret = append(ret, name+" "+goName)
		}
	}
	sort.Strings(ret)
	_, err = f.WriteString(strings.Join(ret, "\n"))
	return
}

// -----------------------------------------------------------------------------
