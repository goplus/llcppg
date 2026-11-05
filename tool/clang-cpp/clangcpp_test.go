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

import (
	"testing"
	_ "unsafe"

	"github.com/goplus/lib/c"
)

//go:linkname tokenName C._ZN5clang3tok12getTokenNameENS0_9TokenKindE
func tokenName(kind uint16) *c.Char

func TestLink(t *testing.T) {
	if got := c.GoString(tokenName(0)); got != "unknown" {
		t.Fatalf("tokenName(0) = %q, want %q", got, "unknown")
	}
}
