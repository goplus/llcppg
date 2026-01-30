/*
 * Copyright (c) 2025 The XGo Authors (xgo.dev). All rights reserved.
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

package tool

import (
	"os"
	"path/filepath"

	"github.com/goplus/gogen"
	"github.com/goplus/llcppg/cl"
	"github.com/goplus/llcppg/cl/nc"
	"github.com/goplus/llcppg/parser"
	"github.com/qiniu/x/errors"
)

// Config is the configuration for generating Go code from C++ header files.
type Config struct {
	PkgPath string
	PkgName string

	Deps []string // dependent packages
	Libs string   // $(pkg-config --libs xxx)

	HeaderFiles []string

	Parser parser.Config
	NC     nc.NodeConverter
}

// GenGo generates Go code from C++ header files.
func GenGo(outDir, buildDir string, conf *Config) (err error) {
	allFile := filepath.Join(buildDir, "llcppg_all.hpp")
	err = composeIncludes(allFile, conf.HeaderFiles)
	if err != nil {
		return
	}
	f, err := parser.ParseFile(nil, allFile, "", &conf.Parser)
	if err != nil {
		return
	}
	pkg, err := cl.Convert(&cl.ConvConfig{
		OutputDir: outDir,
		PkgPath:   conf.PkgPath,
		PkgName:   conf.PkgName,

		Pkg:  f,
		NC:   conf.NC,
		Deps: conf.Deps,
		Libs: conf.Libs,
	})
	if err != nil {
		return
	}
	var errs errors.List
	pkg.ForEachFile(func(fname string, _ *gogen.File) {
		outFile := filepath.Join(outDir, fname)
		e := pkg.WriteFile(outFile, fname)
		if e != nil {
			errs.Add(e)
		}
	})
	return errs.ToError()
}

// composeIncludes create a include list file
// #include "file1.h"
// #include "file2.h"
func composeIncludes(outFile string, files []string) error {
	str := make([]byte, 0, len(files)*20)
	for _, file := range files {
		str = append(str, "#include \""...)
		str = append(str, file...)
		str = append(str, "\"\n"...)
	}
	return os.WriteFile(outFile, str, 0644)
}
