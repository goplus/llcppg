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
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/goplus/gogen"
	"github.com/goplus/llcppg/cl"
	"github.com/goplus/llcppg/clang"
	"github.com/goplus/llcppg/tool"
	"github.com/qiniu/x/errors"
)

// -----------------------------------------------------------------------------

func runGen() {
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("usage: llcppg [-v -debug] [-d \"NAME=VALUE; ...\"] <dest-gopkg-dir> [<src-header-files-and-cfg-dir>]")
		fmt.Println("       llcppg -init [<module-name>]")
		fmt.Println("       llcppg -doc [-http addr] [-open=false] [-all] [-autoexit=false]")
		return
	}
	destDir := args[0]
	srcDir := "."
	if len(args) >= 2 {
		srcDir = args[1]
	}

	userDefines, err := parseDefines()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	log.SetFlags(0)
	if *debug {
		cl.SetDebug(cl.DbgFlagAll)
		tool.SetDebug(tool.DbgFlagAll)
		gogen.SetDebug(gogen.DbgFlagAll)
	} else if *verbose {
		cl.SetDebug(cl.DbgFlagMajorProc)
		tool.SetDebug(tool.DbgFlagSettings)
	}

	idx := clang.CreateIndex(1, 1)
	defer idx.Dispose()
	err = Gen(destDir, srcDir, idx, userDefines, *verbose || *debug)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// -----------------------------------------------------------------------------

// Gen generates Go files from the C/C++ header files and llcppg.cfg in the specified
// source directory and writes them to the destination directory. It uses the provided
// clang index for parsing. The defines are user constants (from the -d flag) injected
// into every generated package.
func Gen(destDir, srcDir string, index clang.Index, defines []cl.Define, verbose bool) (err error) {
	cfg, err := tool.LoadConf(filepath.Join(srcDir, "llcppg.cfg"))
	if err != nil {
		return
	}
	if strings.Contains(cfg.Name, "/") {
		err = errors.New("invalid package name: " + cfg.Name)
		return
	}
	mainPkgName := cfg.Name

	if !verbose {
		if cfg.IgnoreInline {
			cl.InlineFuncIgnore = cl.QuietIgnore // quiet ignore inline functions
		}
		if cfg.NoManglingIgnore {
			cl.NoManglingIgnore = cl.QuietIgnore // quiet ignore functions with no mangled symbol
		}
	}

	if len(cfg.Pkgs) == 0 {
		return genPkg(destDir, srcDir, mainPkgName, index, defines, &cfg)
	}

	errPkgCnt := 0
	for _, pkgSel := range cfg.Pkgs {
		cfgFile := filepath.Join(srcDir, "llcppg-"+pkgSel+".cfg")
		subConf, err := tool.LoadConf(cfgFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "LoadConf failed:", err)
			errPkgCnt++
			continue
		}
		if !strings.HasPrefix(subConf.Name, "/") {
			fmt.Fprintf(os.Stderr,
				"sub package %s: name must start with / but got %s\n", pkgSel, subConf.Name)
			errPkgCnt++
			continue
		}
		subConf.Apply(&cfg)
		errPkgCnt += convPkg(destDir, srcDir, mainPkgName, index, defines, &subConf)
	}
	if errPkgCnt > 0 {
		return fmt.Errorf("%d package(s) failed during generation", errPkgCnt)
	}
	return nil
}

func convPkg(destDir, srcDir, mainPkgName string, index clang.Index, defines []cl.Define, cfg *tool.Config) int {
	err := genPkg(destDir, srcDir, mainPkgName, index, defines, cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

func genPkg(destDir, srcDir, mainPkgName string, index clang.Index, defines []cl.Define, cfg *tool.Config) (err error) {
	pkgName := cfg.Name
	if pos := strings.LastIndex(pkgName, "/"); pos >= 0 {
		destDir = filepath.Join(destDir, pkgName[len(mainPkgName)+1:])
		pkgName = pkgName[pos+1:]
	}

	err = os.MkdirAll(destDir, 0755)
	if err != nil {
		return
	}

	pkg, _, err := cfg.NewPackage("", pkgName, srcDir, index, defines...)
	if err != nil {
		return
	}

	var errs errors.List
	var old = gogen.GeneratedHeader
	defer func() {
		gogen.GeneratedHeader = old
	}()
	gogen.GeneratedHeader = tool.GeneratedHeader
	pkg.ForEachFile(func(fname string, file *gogen.File) {
		if file.Empty() {
			return // skip empty Go files
		}
		goFile := filepath.Join(destDir, fname)
		e := pkg.WriteFile(goFile, fname)
		if e != nil {
			errs.Add(e)
		}
	})
	if err = errs.ToError(); err != nil {
		return
	}
	pub := cl.PubFile(pkg.Public)
	return os.WriteFile(filepath.Join(destDir, "llcppg.pub"), pub, 0644)
}

// -----------------------------------------------------------------------------
