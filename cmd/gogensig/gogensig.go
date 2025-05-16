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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goplus/gogen"
	args "github.com/goplus/llcppg/_xtool/llcppsymg/tool/arg"
	"github.com/goplus/llcppg/cl"
	"github.com/goplus/llcppg/cmd/gogensig/config"
	"github.com/goplus/llcppg/cmd/gogensig/unmarshal"
	llcppg "github.com/goplus/llcppg/config"
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
	for i := 0; i < len(remainArgs); i++ {
		arg := remainArgs[i]
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

	conf, err := config.GetCppgCfgFromPath(cfgFile)
	check(err)
	wd, err := os.Getwd()
	check(err)

	outputDir := filepath.Join(wd, conf.Name)

	err = prepareEnv(outputDir, conf.Deps, modulePath)
	check(err)

	data, err := config.ReadSigfetchFile(filepath.Join(wd, ags.CfgFile))
	check(err)

	convertPkg, err := unmarshal.Pkg(data)
	check(err)

	cvt, err := cl.Convert(&cl.ConvConfig{
		PkgName:  conf.Name,
		SymbFile: filepath.Join(wd, llcppg.LLCPPG_SYMB),
		CfgFile:  filepath.Join(wd, cfgFile),
		Pkg:      convertPkg,
	})
	if err != nil {
		check(err)
	}

	err = config.WritePubFile(filepath.Join(outputDir, llcppg.LLCPPG_PUB), cvt.GenPkg.Pubs)
	check(err)

	err = writePkg(cvt.GenPkg.Pkg(), outputDir)
	check(err)

	err = config.RunCommand(outputDir, "go", "fmt", ".")
	check(err)

	err = config.RunCommand(outputDir, "go", "mod", "tidy")
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

func printUsage() {
	fmt.Fprintln(os.Stderr, "Usage: gogensig [-v|-cfg|-mod] [sigfetch-file]")
}
