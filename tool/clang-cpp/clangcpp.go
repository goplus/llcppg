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

package clangcpp

import "github.com/goplus/lib/c"

const LLGoPackage = "link: -L$(llvm-config --libdir) -lclang-cpp; -lclang-cpp"

// Type represents an opaque clang::Type. A *Type must point to a Clang-owned
// object whose AST context is still alive.
type Type struct{}

// GetTypeClassName returns Clang's static class name, such as "PredefinedSugar".
// The caller must not free the returned pointer.
//
//llgo:link (*Type).GetTypeClassName C._ZNK5clang4Type16getTypeClassNameEv
func (t *Type) GetTypeClassName() *c.Char {
	return nil
}
