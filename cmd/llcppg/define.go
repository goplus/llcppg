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

package main

import (
	"flag"
	"fmt"
	"go/token"
	"strconv"
	"strings"

	"github.com/goplus/llcppg/cl"
)

// -----------------------------------------------------------------------------

// multiFlag collects a repeatable string flag's values in the order given.
type multiFlag []string

func (m *multiFlag) String() string { return strings.Join(*m, "; ") }

func (m *multiFlag) Set(v string) error {
	*m = append(*m, v)
	return nil
}

var (
	defines multiFlag // -d NAME=VALUE, repeatable
)

func init() {
	flag.Var(&defines, "d", "inject constants into the generated package (may be given more than once): -d NAME1=VALUE1 -d NAME2=VALUE2 ...")
}

// -----------------------------------------------------------------------------

func parseDefines() ([]cl.Define, error) {
	args := []string(defines)
	defs := make([]cl.Define, len(args))
	for i, arg := range args {
		def, e := parseDefine(arg)
		if e != nil {
			return nil, e
		}
		defs[i] = def
	}
	return defs, nil
}

func parseDefine(entry string) (def cl.Define, err error) {
	pos := strings.IndexByte(entry, '=')
	if pos < 0 {
		err = fmt.Errorf("invalid define %q: missing '='", entry)
		return
	}
	def.Name = strings.TrimSpace(entry[:pos])
	if !token.IsIdentifier(def.Name) || token.IsKeyword(def.Name) {
		err = fmt.Errorf("invalid define %q: %q is not a valid Go identifier", entry, def.Name)
		return
	}
	value := strings.TrimSpace(entry[pos+1:])
	if len(value) > 0 && value[0] == '\'' && value[len(value)-1] == '\'' {
		def.Value = value[1 : len(value)-1]
	} else if ival, e := strconv.ParseInt(value, 0, 64); e == nil {
		def.Value = ival
	} else if fval, e := strconv.ParseFloat(value, 64); e == nil {
		def.Value = fval
	} else {
		def.Value = value
	}
	return
}

// -----------------------------------------------------------------------------
