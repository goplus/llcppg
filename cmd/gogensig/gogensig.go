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

	"github.com/goplus/llcppg/_xtool/llcppsymg/args"
	"github.com/goplus/llcppg/cmd/gogensig/config"
	"github.com/goplus/llcppg/cmd/gogensig/convert"
	"github.com/goplus/llcppg/cmd/gogensig/convert/filesetprocessor"
	"github.com/goplus/llcppg/cmd/gogensig/dbg"
	"github.com/goplus/llcppg/cmd/gogensig/unmarshal"
)

func main() {

	ags, remainArgs := args.ParseArgs(os.Args[1:], "-", nil)

	if ags.Help {
		printUsage()
		return
	}

	if ags.Verbose {
		dbg.SetDebugAll()
	}

	var cfgFile string
	for i := 0; i < len(remainArgs); i++ {
		arg := remainArgs[i]
		if strings.HasPrefix(arg, "-cfg=") {
			cfgFile = args.StringArg(arg, args.LLCPPG_CFG)
		}
	}
	if cfgFile == "" {
		cfgFile = args.LLCPPG_CFG
	}

	conf, err := config.GetCppgCfgFromPath(cfgFile)
	check(err)
	wd, err := os.Getwd()
	check(err)

	err = prepareEnv(wd, conf.Name, conf.Deps)
	check(err)

	data, err := filesetprocessor.ReadSigfetchFile(filepath.Join(wd, ags.CfgFile))
	check(err)

	convertPkg, err := unmarshal.Pkg(data)
	check(err)

	cvt, err := convert.NewConverter(&convert.ConverterConfig{
		PkgName:  conf.Name,
		SymbFile: filepath.Join(wd, args.LLCPPG_SYMB),
		CfgFile:  filepath.Join(wd, cfgFile),
		PubFile:  filepath.Join(wd, args.LLCPPG_PUB),
		Pkg:      convertPkg,
	})
	if err != nil {
		check(err)
	}
	cvt.Process()
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func prepareEnv(wd, pkg string, deps []string) error {
	dir := filepath.Join(wd, pkg)

	err := os.MkdirAll(dir, 0744)
	if err != nil {
		return err
	}

	err = os.Chdir(pkg)
	if err != nil {
		return err
	}

	err = config.RunCommand(dir, "go", "mod", "init", pkg)
	if err != nil {
		return err
	}

	for _, dep := range deps {
		_, std := convert.IsDepStd(dep)
		if std {
			continue
		}
		err := config.RunCommand(dir, "go", "get", dep)
		if err != nil {
			return err
		}
	}

	return config.RunCommand(dir, "go", "get", "github.com/goplus/llgo@main")
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "Usage: gogensig [-v|-cfg] [sigfetch-file]")
}
