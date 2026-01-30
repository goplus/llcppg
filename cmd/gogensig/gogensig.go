/*
 * Copyright (c) 2024 The XGo Authors (xgo.dev). All rights reserved.
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
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/goplus/gogen"
	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/cl"
	"github.com/goplus/llcppg/cl/nc/ncimpl"
	"github.com/goplus/llcppg/cmd/gogensig/unmarshal"
	llcppg "github.com/goplus/llcppg/config"
	args "github.com/goplus/llcppg/internal/arg"
	"github.com/qiniu/x/errors"
)

func main() {
	ags, remainArgs := args.ParseArgs(os.Args[1:], "-", nil)

	if ags.Help {
		printUsage()
		return
	}

	if ags.Verbose {
		cl.SetDebug(cl.DbgFlagAll)
	}

	var cfgFile string
	var modulePath string
	for _, arg := range remainArgs {
		if strings.HasPrefix(arg, "-cfg=") {
			cfgFile = args.StringArg(arg, llcppg.LLCPPG_CFG)
		}
		if strings.HasPrefix(arg, "-mod=") {
			modulePath = args.StringArg(arg, "")
		}
	}
	if cfgFile == "" {
		cfgFile = llcppg.LLCPPG_CFG
	}
	conf, err := llcppg.GetConfFromFile(cfgFile)
	check(err)
	wd, err := os.Getwd()
	check(err)

	outputDir := filepath.Join(wd, conf.Name)

	err = prepareEnv(outputDir, conf.Deps, modulePath)
	check(err)

	data, err := readSigfetchFile(filepath.Join(wd, ags.CfgFile))
	check(err)

	convertPkg, err := unmarshal.Pkg(data)
	check(err)

	symbFile := filepath.Join(wd, llcppg.LLCPPG_SYMB)
	symbTable, err := llcppg.GetSymTableFromFile(symbFile)
	check(err)

	pkg, err := cl.Convert(&cl.ConvConfig{
		PkgName: conf.Name,
		Pkg:     convertPkg.File,
		NC: &ncimpl.Converter{
			PkgName: conf.Name,
			Pubs:    conf.TypeMap,
			ConvSym: func(name *ast.Object, mangleName string) (goName string, err error) {
				item, err := symbTable.LookupSymbol(mangleName)
				if err != nil {
					return
				}
				return item.Go, nil
			},
			FileMap:        convertPkg.FileMap,
			TrimPrefixes:   conf.TrimPrefixes,
			KeepUnderScore: conf.KeepUnderScore,
		},
		Deps: conf.Deps,
		Libs: conf.Libs,
	})
	check(err)

	err = llcppg.WritePubFile(filepath.Join(outputDir, llcppg.LLCPPG_PUB), pkg.Pubs)
	check(err)

	err = writePkg(pkg.Package, outputDir)
	check(err)

	err = runCommand(outputDir, "go", "fmt", ".")
	check(err)

	err = runCommand(outputDir, "go", "mod", "tidy")
	check(err)
}

// Write all files in the package to the output directory
func writePkg(pkg *gogen.Package, outDir string) error {
	var errs errors.List
	pkg.ForEachFile(func(fname string, _ *gogen.File) {
		if fname != "" { // gogen default fname
			outFile := filepath.Join(outDir, fname)
			e := pkg.WriteFile(outFile, fname)
			if e != nil {
				errs.Add(e)
			}
		}
	})
	return errs.ToError()
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func prepareEnv(outputDir string, deps []string, modulePath string) error {
	err := os.MkdirAll(outputDir, 0744)
	if err != nil {
		return err
	}

	err = os.Chdir(outputDir)
	if err != nil {
		return err
	}

	return cl.ModInit(deps, outputDir, modulePath)
}

func runCommand(dir, cmdName string, args ...string) error {
	execCmd := exec.Command(cmdName, args...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Dir = dir
	return execCmd.Run()
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "Usage: gogensig [-v|-cfg|-mod] [sigfetch-file]")
}

func readSigfetchFile(sigfetchFile string) ([]byte, error) {
	_, file := filepath.Split(sigfetchFile)
	if file == "-" {
		return io.ReadAll(os.Stdin)
	}
	return os.ReadFile(sigfetchFile)
}
