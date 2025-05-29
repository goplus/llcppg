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
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"

	llcppg "github.com/goplus/llcppg/config"
	"github.com/goplus/llgo/xtool/env"

	// import to make it linked in go.mod
	_ "github.com/goplus/lib/c"
)

type modeFlags int

const (
	ModeCodegen modeFlags = 1 << iota
	ModeSymbGen
	ModeAll = ModeCodegen | ModeSymbGen
)

type verboseFlags int

const (
	VerboseSymg verboseFlags = 1 << iota
	VerboseSigfetch
	VerboseGogen
	VerboseAll = VerboseSymg | VerboseSigfetch | VerboseGogen
)

type CommandOptions struct {
	Name    string
	Args    []string
	Verbose bool
}

func command(opts CommandOptions) *exec.Cmd {
	args := opts.Args
	if opts.Verbose {
		args = append([]string{"-v"}, args...)
	}
	return exec.Command(opts.Name, args...)
}

func llcppsymg(conf []byte, v verboseFlags) error {
	cmd := command(CommandOptions{
		Name:    "llcppsymg",
		Args:    []string{"-"},
		Verbose: (v & VerboseSymg) != 0,
	})
	cmd.Stdin = bytes.NewReader(conf)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func llcppsigfetch(conf []byte, v verboseFlags, out *io.PipeWriter) {
	cmd := command(CommandOptions{
		Name:    "llcppsigfetch",
		Args:    []string{"-"},
		Verbose: (v & VerboseSigfetch) != 0,
	})
	cmd.Stdin = bytes.NewReader(conf)
	cmd.Stdout = out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	check(err)
	out.Close()
}

func gogensig(in io.Reader, cfg string, modulePath string, v verboseFlags) error {
	cmd := command(CommandOptions{
		Name:    "gogensig",
		Args:    []string{"-", "-cfg=" + cfg, "-mod=" + modulePath},
		Verbose: (v & VerboseGogen) != 0,
	})
	cmd.Stdin = in
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	var symbGen, codeGen, help bool
	var vSymg, vSigfetch, vGogen, vAll bool
	var modulePath string
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: llcppg [-v|-vfetch|-vsymg|-vgogen] [-symbgen] [-codegen] [-h|--help] [config-file]")
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
	}
	flag.BoolVar(&vAll, "v", false, "Enable verbose output")
	flag.BoolVar(&vSigfetch, "vfetch", false, "Enable verbose of llcppsigfetch")
	flag.BoolVar(&vSymg, "vsymg", false, "Enable verbose of llcppsymg")
	flag.BoolVar(&vGogen, "vgogen", false, "Enable verbose of gogensig")
	flag.BoolVar(&symbGen, "symbgen", false, "Only use llcppsymg to generate llcppg.symb.json")
	flag.BoolVar(&codeGen, "codegen", false, "Only use (llcppsigfetch & gogensig) to generate go code binding")
	flag.BoolVar(&help, "h", false, "Display help information")
	flag.BoolVar(&help, "help", false, "Display help information")
	flag.StringVar(&modulePath, "mod", "", "The module path of the generated code,if not set,will not init a new module")
	flag.Parse()

	verbose := verboseFlags(0)
	mode := ModeAll
	if vAll {
		verbose = VerboseAll
		mode = ModeAll
	}
	if vSigfetch {
		verbose |= VerboseSigfetch
	}
	if vGogen {
		verbose |= VerboseGogen
	}
	if vSymg {
		verbose |= VerboseSymg
	}

	if codeGen {
		mode = ModeCodegen
	}
	if symbGen {
		mode = ModeSymbGen
	}

	if help {
		flag.Usage()
		return
	}

	remainArgs := flag.Args()

	var cfgFile string
	if len(remainArgs) > 0 {
		cfgFile = remainArgs[0]
	} else {
		cfgFile = llcppg.LLCPPG_CFG
	}

	do(cfgFile, mode, verbose, modulePath)
}

func do(cfgFile string, mode modeFlags, verbose verboseFlags, modulePath string) {
	f, err := os.Open(cfgFile)
	check(err)
	defer f.Close()

	var conf llcppg.Config
	json.NewDecoder(f).Decode(&conf)
	conf.CFlags = env.ExpandEnv(conf.CFlags)
	conf.Libs = env.ExpandEnv(conf.Libs)

	b, err := json.MarshalIndent(&conf, "", "  ")
	check(err)

	if mode&ModeSymbGen != 0 {
		err = llcppsymg(b, verbose)
		check(err)
	}

	if mode&ModeCodegen != 0 {
		r, w := io.Pipe()
		go llcppsigfetch(b, verbose, w)

		err = gogensig(r, cfgFile, modulePath, verbose)
		check(err)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
