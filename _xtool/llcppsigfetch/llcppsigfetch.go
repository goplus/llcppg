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

	"github.com/goplus/llcppg/_xtool/llcppsigfetch/parse"
	args "github.com/goplus/llcppg/_xtool/llcppsymg/tool/arg"
	clangutils "github.com/goplus/llcppg/_xtool/llcppsymg/tool/clang"
	"github.com/goplus/llcppg/_xtool/llcppsymg/tool/config"
	llcppg "github.com/goplus/llcppg/config"
)

func main() {
	ags, remainArgs := args.ParseArgs(os.Args[1:], llcppg.LLCPPG_CFG, map[string]bool{
		"--extract": true,
	})

	if ags.Help {
		printUsage()
		return
	}
	if ags.Verbose {
		parse.SetDebug(parse.DbgFlagAll)
	}
	extract := false
	out := false

	var extractFile string
	isTemp := false
	isCpp := true
	otherArgs := []string{}

	for i := 0; i < len(remainArgs); i++ {
		arg := remainArgs[i]
		switch {
		case arg == "--extract":
			extract = true
			if i+1 < len(remainArgs) && !strings.HasPrefix(remainArgs[i+1], "-") {
				extractFile = remainArgs[i+1]
				i++
			} else {
				fmt.Fprintln(os.Stderr, "Error: --extract requires a valid file argument")
				printUsage()
				os.Exit(1)
			}
		case strings.HasPrefix(arg, "-out="):
			out = args.BoolArg(arg, false)
		case strings.HasPrefix(arg, "-temp="):
			isTemp = args.BoolArg(arg, false)
		case strings.HasPrefix(arg, "-cpp="):
			isCpp = args.BoolArg(arg, true)
		case strings.HasPrefix(arg, "-ClangResourceDir="):
			// temp to avoid call clang  in llcppsigfetch,will cause hang
			parse.ClangResourceDir = args.StringArg(arg, "")
		case strings.HasPrefix(arg, "-ClangSearchPath="):
			// temp to avoid call clang  in llcppsigfetch,will cause hang
			parse.ClangSearchPath = strings.Split(args.StringArg(arg, ""), ",")
		default:
			otherArgs = append(otherArgs, arg)
		}
	}

	parseConfig := &parse.Config{
		Exec: parse.OutputPkg,
		Out:  out,
	}

	if extract {
		conf, err := buildExtractConfig(extractFile, isTemp, isCpp, otherArgs)
		check(err)
		parseConfig.Conf = conf
	} else {
		conf, err := config.GetConf(ags.UseStdin, ags.CfgFile)
		check(err)
		parseConfig.Conf = conf
	}

	err := parse.Do(parseConfig)
	check(err)
}

func buildExtractConfig(extractFile string, isTemp bool, isCpp bool, otherArgs []string) (conf *llcppg.Config, err error) {
	var file string
	cflags := otherArgs
	if isTemp {
		temp, err := os.Create(clangutils.TEMP_FILE)
		if err != nil {
			panic(err)
		}
		defer temp.Close()
		defer os.Remove(file)
		temp.Write([]byte(extractFile))
		file = temp.Name()
		cflags = append(cflags, "-I"+filepath.Dir(file))
	} else {
		file = extractFile
	}
	return &llcppg.Config{
		Include:   []string{file},
		CFlags:    strings.Join(cflags, ""),
		Cplusplus: isCpp,
	}, nil
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  llcppsigfetch [-v] [-out=<bool>] [config_file]")
	fmt.Println("  OR")
	fmt.Println("  llcppsigfetch --extract <file> [-out=<bool>] [-temp=<bool>] [-cpp=<bool>] [-v] [args...]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  [config_file]: Path to the configuration file (use '-' for stdin)")
	fmt.Println("                   If not provided, uses default 'llcppg.cfg'")
	fmt.Println("  -out=<bool>:     Optional. Set to 'true' to output results to a file,")
	fmt.Println("                   'false' (default) to output to stdout")
	fmt.Println("                   This option can be used with both modes")
	fmt.Println("")
	fmt.Println("  --extract:       Extract information from a single file")
	fmt.Println("    <file>:        Path to the file to process, or file content if -temp=true")
	fmt.Println("    -temp=<bool>:  Optional. Set to 'true' if <file> contains file content,")
	fmt.Println("                   'false' (default) if it's a file path")
	fmt.Println("    -cpp=<bool>:   Optional. Set to 'true' if the language is C++ (default: true)")
	fmt.Println("                   If not present, <file> is a file path")
	fmt.Println("    [args]:        Optional additional arguments")
	fmt.Println("                   Default for C++: -x c++")
	fmt.Println("                   Default for C: -x c")
	fmt.Println("")
	fmt.Println("  --help, -h:      Show this help message")
	fmt.Println("")
	fmt.Println("Note: The two usage modes are mutually exclusive. Use either [<config_file>] OR --extract, not both.")
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
