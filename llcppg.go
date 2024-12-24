/*
 * Copyright (c) 2024 The GoPlus Authors (goplus.org). All rights reserved.
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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/goplus/llcppg/_xtool/llcppsymg/args"
	"github.com/goplus/llcppg/types"
	"github.com/goplus/llgo/xtool/env"
)

var verbose bool

func command(name string, args ...string) *exec.Cmd {
	if verbose {
		args = append([]string{"-v"}, args...)
	}
	return exec.Command(name, args...)
}

func llcppsymg(conf []byte) error {
	cmd := command("llcppsymg", "-")
	cmd.Stdin = bytes.NewReader(conf)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func llcppsigfetch(conf []byte, out *io.PipeWriter) {
	cmd := command("llcppsigfetch", "-")
	cmd.Stdin = bytes.NewReader(conf)
	cmd.Stdout = out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	check(err)
	out.Close()
}

func gogensig(in io.Reader, cfg string) error {
	cmd := command("gogensig", "-", "-cfg="+cfg)
	cmd.Stdin = in
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	ags, _ := args.ParseArgs(os.Args[1:], args.LLCPPG_CFG, nil)
	if ags.Help {
		fmt.Fprintln(os.Stderr, "Usage: llcppg [config-file] [-v]")
		return
	}
	verbose = ags.Verbose

	f, err := os.Open(ags.CfgFile)
	check(err)
	defer f.Close()

	var conf types.Config
	json.NewDecoder(f).Decode(&conf)
	conf.CFlags = env.ExpandEnv(conf.CFlags)
	conf.Libs = env.ExpandEnv(conf.Libs)

	b, err := json.MarshalIndent(&conf, "", "  ")
	check(err)

	err = llcppsymg(b)
	check(err)

	r, w := io.Pipe()
	go llcppsigfetch(b, w)

	err = gogensig(r, ags.CfgFile)
	check(err)
}

func check(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1) // trigger github action failure
	}
}
