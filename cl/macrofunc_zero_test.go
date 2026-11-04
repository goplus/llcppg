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
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"testing"

	"github.com/goplus/llcppg/clang"
	lc "github.com/llarhub/clang-c"
)

func TestCZeroArgSource(t *testing.T) {
	const src = `typedef int number_t;
void target(int *value);
#define PLAIN() target(0)
#define PARENS() target(((0)))
#define NEGATIVE() target(-0)
#define INT_CAST() target((number_t)0)
#define FOLDED_ZERO() target(1 - 1)
#define FOLDED_NONZERO() target(2 - 1)
#define FLOAT_ZERO() target(0.0)
#define NONZERO() target(256)
`
	filename := filepath.Join(t.TempDir(), "in.h")
	if err := os.WriteFile(filename, []byte(src), 0o600); err != nil {
		t.Fatal(err)
	}
	idx := clang.CreateIndex(1, 1)
	defer idx.Dispose()
	tu, err := idx.ParseTranslationUnit(clang.DetailedPreprocessingRecord, filename, "-x", "c")
	if err != nil {
		t.Fatal(err)
	}
	defer tu.Dispose()

	macros := make(map[string]clang.Cursor)
	var target clang.Cursor
	clang.VisitChildren(tu.Cursor(), func(cur, parent clang.Cursor) clang.ChildVisitResult {
		switch cur.Kind {
		case lc.Cursor_FunctionDecl:
			if clang.String(cur) == "target" {
				target = cur
			}
		case lc.Cursor_MacroDefinition:
			macros[clang.String(cur)] = cur
		}
		return clang.Continue
	})
	if target.IsNull() != 0 {
		t.Fatal("target declaration not found")
	}

	ctx := newPkgCtx()
	ctx.types = map[string]typeObj{
		"number_t": {TypeName: types.NewTypeName(token.NoPos, nil, "number_t", types.Typ[types.Int32])},
	}
	goType := types.NewPointer(types.Typ[types.Int32])
	cType := target.Argument(0).Type()
	for _, tt := range []struct {
		name string
		want bool
	}{
		{"PLAIN", true},
		{"PARENS", true},
		{"NEGATIVE", true},
		{"INT_CAST", true},
		{"FOLDED_ZERO", true},
		{"FOLDED_NONZERO", false},
		{"FLOAT_ZERO", false},
		{"NONZERO", false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			cur, ok := macros[tt.name]
			if !ok {
				t.Fatalf("macro %s not found", tt.name)
			}
			tokens, dispose := tu.Tokenize(cur.Extent())
			defer dispose()
			_, next, variadic, ok := macroParamNames(tu, tokens)
			if !ok || variadic {
				t.Fatal("invalid macro parameters")
			}
			call, ok := evalFuncCallExpr(ctx, tu, tokens[next:])
			if !ok || len(call.args) != 1 {
				t.Fatal("invalid macro call")
			}
			if got := isZeroArg(call.args[0], cType, goType); got != tt.want {
				t.Errorf("isZeroArg = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("CXXReference", func(t *testing.T) {
		filename := filepath.Join(t.TempDir(), "reference.h")
		if err := os.WriteFile(filename, []byte("void target_ref(const int &value);\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		cppTU, err := idx.ParseTranslationUnit(clang.DetailedPreprocessingRecord, filename, "-x", "c++")
		if err != nil {
			t.Fatal(err)
		}
		defer cppTU.Dispose()

		var ref clang.Cursor
		clang.VisitChildren(cppTU.Cursor(), func(cur, parent clang.Cursor) clang.ChildVisitResult {
			if cur.Kind == lc.Cursor_FunctionDecl && clang.String(cur) == "target_ref" {
				ref = cur
			}
			return clang.Continue
		})
		if ref.IsNull() != 0 {
			t.Fatal("reference declaration not found")
		}
		if isZeroArg(funcCallArg{fca_Const, 0}, ref.Argument(0).Type(), goType) {
			t.Fatal("a C++ reference must not become a nil pointer")
		}
	})
}
