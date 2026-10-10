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
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/goplus/llcppg/internal/docserver"
)

var (
	docHTTP = flag.String("http", docserver.DefaultAddr, "listen address for -doc (loopback only unless an explicit non-loopback host is given)")
	docOpen = flag.Bool("open", true, "launch the default browser for -doc; use -open=false to only print the URL")
	docAll  = flag.Bool("all", false, "include unexported declarations")
)

func main() {
	flag.Parse()
	runDoc()
}

// runDoc serves the documentation of the Go package in the current directory
// (`llcppg -doc`). It blocks until interrupted (SIGINT), then shuts down
// gracefully. Usage errors (no Go files, more than one package, a bad listen
// address) are reported and exit non-zero.
func runDoc() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	err := docserver.Run(ctx, ".", docserver.Options{
		Addr:        *docHTTP,
		OpenBrowser: *docOpen,
		AllDecls:    *docAll,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
