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

package clang

import (
	"strings"
)

// -----------------------------------------------------------------------------

/**
 * Retrieve a name for the entity referenced by this cursor.
 */
func Mangling(fn Cursor) string {
	manglingName := GoStringAndDispose(fn.Mangling())
	if manglingName == "" {
		return ""
	}
	info := fn.TranslationUnit().TargetInfo()
	triple := GoStringAndDispose(info.Triple())
	info.Dispose()
	// TargetInfo supplies a normalized arch-vendor-os[-environment] triple.
	_, rest, _ := strings.Cut(triple, "-")
	_, targetOS, _ := strings.Cut(rest, "-")
	for _, os := range []string{"darwin", "macosx", "ios", "tvos", "watchos", "driverkit", "xros"} {
		if strings.HasPrefix(targetOS, os) {
			// Darwin adds one underscore: _normal becomes normal, while
			// __leading keeps the underscore from the source declaration.
			return strings.TrimPrefix(manglingName, "_")
		}
	}
	return manglingName
}

// -----------------------------------------------------------------------------
