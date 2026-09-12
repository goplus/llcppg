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
	"go/token"
	"go/types"
	"log"
	"os"

	"github.com/goplus/gogen"
	"github.com/goplus/gogen/packages"
)

func main() {
	pkg := gogen.NewPackage("", "foo", &gogen.Config{
		Importer:        packages.NewImporter(nil),
		LoadNamed:       nil,
		HandleErr:       nil,
		NewBuiltin:      nil,
		NodeInterpreter: nil,
		CanImplicitCast: nil,
		DefaultGoFile:   "",
	})
	pkg.SetRedeclarable(true)
	sig := types.NewSignatureType(nil, nil, nil, nil, nil, false)
	_, err := pkg.NewFuncWith(token.NoPos, "g", sig, nil)
	if err != nil {
		log.Panicln("compileFunc:", "g", err)
	}
	err = pkg.WriteTo(os.Stdout)
	if err != nil {
		log.Panicln("gogen.WriteTo failed:", err)
	}
}
