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
	"log"
	"os/exec"
	"strings"
)

// -----------------------------------------------------------------------------

var appPath string

// Path returns the path to the clang binary.
func Path() string {
	if appPath == "" {
		appPath = findAppPath()
	}
	return appPath
}

func findAppPath() string {
	b, err := exec.Command("llgo", "env", "LLGO_CLANG").Output()
	if err != nil {
		log.Panicln(err)
	}
	return strings.TrimSpace(string(b))
}

// -----------------------------------------------------------------------------
