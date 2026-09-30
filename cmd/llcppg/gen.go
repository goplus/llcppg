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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goplus/gogen"
	"github.com/goplus/llcppg/clang"
	"github.com/goplus/llcppg/tool"
	"github.com/qiniu/x/errors"
)

// -----------------------------------------------------------------------------

// Gen generates Go files from the C/C++ header files and llcppg.cfg in the specified
// source directory and writes them to the destination directory. It uses the provided
// clang index for parsing.
func Gen(destDir, srcDir string, index clang.Index) (err error) {
	cfg, err := tool.LoadConf(filepath.Join(srcDir, "llcppg.cfg"))
	if err != nil {
		return
	}
	if strings.Contains(cfg.Name, "/") {
		err = errors.New("invalid package name: " + cfg.Name)
		return
	}
	mainPkgName := cfg.Name

	if len(cfg.Pkgs) == 0 {
		return genPkg(destDir, srcDir, mainPkgName, index, &cfg)
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
		subConf.Apply(&cfg)
		errPkgCnt += convPkg(destDir, srcDir, mainPkgName, index, &subConf)
	}
	if errPkgCnt > 0 {
		return fmt.Errorf("%d package(s) failed during generation", errPkgCnt)
	}
	return nil
}

func convPkg(destDir, srcDir, mainPkgName string, index clang.Index, cfg *tool.Config) int {
	err := genPkg(destDir, srcDir, mainPkgName, index, cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

func genPkg(destDir, srcDir, mainPkgName string, index clang.Index, cfg *tool.Config) (err error) {
	pkgName := cfg.Name
	if pos := strings.LastIndex(pkgName, "/"); pos >= 0 {
		destDir = filepath.Join(destDir, pkgName[len(mainPkgName)+1:])
		pkgName = pkgName[pos+1:]
	}

	err = os.MkdirAll(destDir, 0755)
	if err != nil {
		return
	}

	pkg, _, err := cfg.NewPackage("", pkgName, srcDir, index)
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
	return errs.ToError()
}

// -----------------------------------------------------------------------------
