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
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/goplus/gogen"
	"github.com/goplus/llcppg/_xtool/parse"
	"github.com/goplus/llcppg/_xtool/symg"
	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/cl"
	"github.com/goplus/llcppg/cl/nc/ncimpl"
	llcppg "github.com/goplus/llcppg/config"
	"github.com/goplus/llcppg/internal/gowrite"
	"github.com/goplus/llgo/xtool/env"
	"github.com/qiniu/x/errors"

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

func buildSymbolTable(conf *llcppg.Config, v verboseFlags) error {
	if (v & VerboseSymg) != 0 {
		symg.SetDebug(symg.DbgFlagAll)
	}
	libMode := symg.ModeDynamic
	if conf.StaticLib {
		libMode = symg.ModeStatic
	}
	symbolTable, err := symg.Do(&symg.Config{
		Libs:         conf.Libs,
		CFlags:       conf.CFlags,
		Includes:     conf.Include,
		Mix:          conf.Mix,
		TrimPrefixes: conf.TrimPrefixes,
		SymMap:       conf.SymMap,
		IsCpp:        conf.Cplusplus,
		HeaderOnly:   conf.HeaderOnly,
		LibMode:      libMode,
	})
	if err != nil {
		return err
	}
	jsonData, err := json.MarshalIndent(symbolTable, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(llcppg.LLCPPG_SYMB, jsonData, os.ModePerm)
}

func parseHeaders(conf *llcppg.Config, v verboseFlags) (*llcppg.Pkg, error) {
	if (v & VerboseSigfetch) != 0 {
		parse.SetDebug(parse.DbgFlagAll)
	}
	var pkg *llcppg.Pkg
	parseCfg := &parse.Config{
		Conf: conf,
		Exec: func(_ *parse.Config, p *llcppg.Pkg) {
			pkg = p
		},
	}
	if err := parse.Do(parseCfg); err != nil {
		return nil, err
	}
	return pkg, nil
}

// gengo converts C header AST information into corresponding LLGo bindings.
func gengo(conf *llcppg.Config, in *llcppg.Pkg, modulePath string, v verboseFlags) error {
	if (v & VerboseGogen) != 0 {
		cl.SetDebug(cl.DbgFlagAll)
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	outputDir := filepath.Join(wd, conf.Name)
	if err := prepareEnv(outputDir, conf.Deps, modulePath); err != nil {
		return err
	}
	if err := os.Chdir(outputDir); err != nil {
		return err
	}
	defer func() { _ = os.Chdir(wd) }()
	symbFile := filepath.Join(wd, llcppg.LLCPPG_SYMB)
	symbTable, err := llcppg.GetSymTableFromFile(symbFile)
	if err != nil {
		return err
	}
	pkg, err := cl.Convert(&cl.ConvConfig{
		OutputDir: outputDir,
		PkgName:   conf.Name,
		Pkg:       in.File,
		NC: &ncimpl.Converter{
			PkgName: conf.Name,
			Pubs:    conf.TypeMap,
			ConvSym: func(name *ast.Object, mangleName string) (goName string, err error) {
				item, err := symbTable.LookupSymbol(mangleName)
				if err != nil {
					return "", err
				}
				return item.Go, nil
			},
			FileMap:        in.FileMap,
			TrimPrefixes:   conf.TrimPrefixes,
			KeepUnderScore: conf.KeepUnderScore,
		},
		Deps: conf.Deps,
		Libs: conf.Libs,
	})
	if err != nil {
		return err
	}
	if err := llcppg.WritePubFile(filepath.Join(outputDir, llcppg.LLCPPG_PUB), pkg.Pubs); err != nil {
		return err
	}
	if err := writePkg(pkg.Package, outputDir); err != nil {
		return err
	}
	if err := runCommand(outputDir, "go", "fmt", "."); err != nil {
		return err
	}
	return runCommand(outputDir, "go", "mod", "tidy")
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
	flag.BoolVar(&symbGen, "symbgen", false, "Only use llcppsymg to generate llcppg.sym.json")
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
	err = json.NewDecoder(f).Decode(&conf)
	check(err)

	// Keep original libs expression for generated LLGoPackage, but use expanded
	// values for symbol extraction and header parsing passes.
	rawLibs := conf.Libs
	conf.CFlags = env.ExpandEnv(conf.CFlags)
	conf.Libs = env.ExpandEnv(conf.Libs)

	if mode&ModeSymbGen != 0 {
		// Pass 1: build native symbol table.
		err = buildSymbolTable(&conf, verbose)
		check(err)
	}

	if mode&ModeCodegen != 0 {
		// Pass 2: parse headers into AST package.
		pkg, err := parseHeaders(&conf, verbose)
		check(err)
		// Pass 3: convert C header AST information into corresponding LLGo bindings.
		codegenConf := conf
		codegenConf.Libs = rawLibs
		err = gengo(&codegenConf, pkg, modulePath, verbose)
		check(err)
	}
}

func prepareEnv(outputDir string, deps []string, modulePath string) error {
	if err := os.MkdirAll(outputDir, 0744); err != nil {
		return err
	}
	return cl.ModInit(deps, outputDir, modulePath)
}

func writePkg(pkg *gogen.Package, outDir string) error {
	var errs errors.List
	pkg.ForEachFile(func(fname string, _ *gogen.File) {
		if fname != "" {
			outFile := filepath.Join(outDir, fname)
			if err := gowrite.WriteFile(pkg, outFile, fname); err != nil {
				errs.Add(err)
			}
		}
	})
	return errs.ToError()
}

func runCommand(dir, cmdName string, args ...string) error {
	execCmd := exec.Command(cmdName, args...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Dir = dir
	return execCmd.Run()
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
