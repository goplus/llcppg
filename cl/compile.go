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

	"github.com/goplus/gogen"
	"github.com/goplus/lib/c"
	"github.com/goplus/llcppg/clang"
)

// -----------------------------------------------------------------------------

type PkgInfo struct {
}

// Package represents a generated Go package.
type Package struct {
	*gogen.Package
	pi *PkgInfo
}

// Reused specifies to reuse the Package instance between processing multiple C/C++ header files.
type Reused struct {
	pkg Package
}

// -----------------------------------------------------------------------------

// Config specifies the configuration for compiling C/C++ header files.
type Config struct {
	// An Importer resolves import paths to Packages.
	Importer types.Importer

	// Include specifies include searching directories.
	Include []string

	// Reused specifies to reuse the Package instance between processing multiple C/C++ header files.
	*Reused
}

// -----------------------------------------------------------------------------

// Source represents a C/C++ header to compile.
type Source struct {
	clang.TranslationUnit
	PresumedFile *c.Char
}

// -----------------------------------------------------------------------------

// NewPackage creates a new Package instance for the specified package path and name, using
// the provided Source and Config.
func NewPackage(pkgPath, pkgName string, src Source, conf *Config) (pkg Package, err error) {
	panic("todo")
}

// -----------------------------------------------------------------------------
