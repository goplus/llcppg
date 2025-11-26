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
	"encoding/json"
	"fmt"
	"os"

	"github.com/goplus/llcppg/_xtool/internal/symbol"
	"github.com/goplus/llcppg/_xtool/llcppsymg/internal/symg"
	llcppg "github.com/goplus/llcppg/config"
	args "github.com/goplus/llcppg/internal/arg"
)

func main() {
	ags, _ := args.ParseArgs(os.Args[1:], llcppg.LLCPPG_CFG, nil)

	if ags.Help {
		printUsage()
		return
	}

	var err error
	var conf llcppg.Config

	if ags.UseStdin {
		conf, err = llcppg.GetConfFromStdin()
	} else {
		conf, err = llcppg.GetConfFromFile(ags.CfgFile)
	}
	check(err)

	if ags.Verbose {
		symg.SetDebug(symg.DbgFlagAll)
		if ags.UseStdin {
			fmt.Println("Config From Stdin")
		} else {
			fmt.Println("Config From File", ags.CfgFile)
		}
		fmt.Println("Name:", conf.Name)
		fmt.Println("CFlags:", conf.CFlags)
		fmt.Println("Libs:", conf.Libs)
		fmt.Println("Include:", conf.Include)
		fmt.Println("TrimPrefixes:", conf.TrimPrefixes)
		fmt.Println("Cplusplus:", conf.Cplusplus)
		fmt.Println("SymMap:", conf.SymMap)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to parse config file:", ags.CfgFile)
	}

	libMode := symbol.ModeDynamic
	if conf.StaticLib {
		libMode = symbol.ModeStatic
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
	check(err)

	jsonData, err := json.MarshalIndent(symbolTable, "", "  ")
	check(err)

	err = os.WriteFile(llcppg.LLCPPG_SYMB, jsonData, 0644)
	check(err)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "Usage: llcppsymg [-v] [config-file]")
}
