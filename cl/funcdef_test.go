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

package cl_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/goplus/llcppg/cl"
	"github.com/goplus/llcppg/clang"
)

// funcDefHeader is a small C API used by the FuncDef validation tests: an object
// type plus a few functions whose first parameter is a pointer to it.
const funcDefHeader = `
typedef struct {
    int v;
} PyObject;

PyObject *PyModule_New(const char *name);
PyObject *PyModule_GetDict(PyObject *);
int PyObject_HasNoArg(void);
`

// runFuncDef parses funcDefHeader and generates a package with the given TypeDef and
// FuncDef config, returning whether generation reported an error.
func runFuncDef(t *testing.T, typeDef map[string][]cl.TypeDef, funcDef map[string][]cl.FuncDef) error {
	t.Helper()
	dir := t.TempDir()
	hfile := filepath.Join(dir, "in.h")
	if err := os.WriteFile(hfile, []byte(funcDefHeader), 0600); err != nil {
		t.Fatal(err)
	}

	idx := clang.CreateIndex(1, 1)
	defer idx.Dispose()
	u, err := idx.ParseTranslationUnit(clang.DetailedPreprocessingRecord, hfile, "-x", "c")
	if err != nil {
		t.Fatal("ParseTranslationUnit:", err)
	}
	defer u.Dispose()

	_, err = cl.NewPackage("testcl/funcdef", "foo", []cl.Source{u}, &cl.Config{
		Language:   cl.LanguageC,
		TypePrefix: []string{"Py"},
		TypeDef:    typeDef,
		FuncDef:    funcDef,
	})
	return err
}

func TestFuncDefErrors(t *testing.T) {
	typeDef := map[string][]cl.TypeDef{
		"": {{Name: "ModuleObject", Fields: []string{".Object"}}},
	}
	tests := []struct {
		name    string
		typeDef map[string][]cl.TypeDef
		funcDef map[string][]cl.FuncDef
		wantErr bool
	}{
		{
			name:    "ok package func",
			funcDef: map[string][]cl.FuncDef{"": {{Name: "NewModule", Link: "PyModule_New"}}},
		},
		{
			name:    "ok method on TypeDef type",
			typeDef: typeDef,
			funcDef: map[string][]cl.FuncDef{"": {{Name: "(*ModuleObject).GetDict", Link: "PyModule_GetDict"}}},
		},
		{
			name:    "missing Name",
			funcDef: map[string][]cl.FuncDef{"": {{Link: "PyModule_New"}}},
			wantErr: true,
		},
		{
			name:    "missing Link",
			funcDef: map[string][]cl.FuncDef{"": {{Name: "NewModule"}}},
			wantErr: true,
		},
		{
			name:    "invalid Name form",
			funcDef: map[string][]cl.FuncDef{"": {{Name: "(*ModuleObject)GetDict", Link: "PyModule_GetDict"}}},
			wantErr: true,
		},
		{
			name:    "invalid Go identifier",
			funcDef: map[string][]cl.FuncDef{"": {{Name: "New Module", Link: "PyModule_New"}}},
			wantErr: true,
		},
		{
			name:    "Link symbol not found",
			funcDef: map[string][]cl.FuncDef{"": {{Name: "NewModule", Link: "PyModule_DoesNotExist"}}},
			wantErr: true,
		},
		{
			name:    "receiver type does not exist",
			funcDef: map[string][]cl.FuncDef{"": {{Name: "(*NoSuchType).GetDict", Link: "PyModule_GetDict"}}},
			wantErr: true,
		},
		{
			name:    "method on func without parameters",
			typeDef: typeDef,
			funcDef: map[string][]cl.FuncDef{"": {{Name: "(*ModuleObject).Bad", Link: "PyObject_HasNoArg"}}},
			wantErr: true,
		},
		{
			name: "duplicate Go name",
			funcDef: map[string][]cl.FuncDef{"": {
				{Name: "NewModule", Link: "PyModule_New"},
				{Name: "NewModule", Link: "PyModule_GetDict"},
			}},
			wantErr: true,
		},
		{
			// PyObject_HasNoArg is emitted by default as a package-level function
			// (it has no pointer-to-type first parameter) under the Go name
			// PyObject_HasNoArg. A FuncDef that reuses that name for a different
			// symbol must be rejected.
			name:    "name collides with default function",
			funcDef: map[string][]cl.FuncDef{"": {{Name: "PyObject_HasNoArg", Link: "PyModule_New"}}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			td := tt.typeDef
			if td == nil {
				td = map[string][]cl.TypeDef{}
			}
			err := runFuncDef(t, td, tt.funcDef)
			if (err != nil) != tt.wantErr {
				t.Fatalf("runFuncDef() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}
