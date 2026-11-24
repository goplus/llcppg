// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// From golang.org/x/sys/unix/internal/mkmerge at v0.47.0.

package merge

import "go/ast"

// filter selects declarations by their original indexes.
// For const (A = 1; B = 2), the spec index lets the output keep just B.
func filter(f *ast.File, keep func(declIndex, specIndex int) bool) []ast.Decl {
	decls := make([]ast.Decl, 0, len(f.Decls))
	for i, decl := range f.Decls {
		switch decl := decl.(type) {
		case *ast.GenDecl:
			filtered := *decl
			filtered.Specs = make([]ast.Spec, 0, len(decl.Specs))
			for j, spec := range decl.Specs {
				if keep(i, j) {
					filtered.Specs = append(filtered.Specs, spec)
				}
			}
			if len(filtered.Specs) > 0 {
				decls = append(decls, &filtered)
			}
		case *ast.FuncDecl:
			if keep(i, -1) {
				decls = append(decls, decl)
			}
		}
	}
	return decls
}

/*
Copyright 2009 The Go Authors.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google LLC nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/
