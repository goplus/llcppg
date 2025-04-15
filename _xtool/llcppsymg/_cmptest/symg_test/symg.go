package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goplus/llcppg/_xtool/llcppsymg/config"
	"github.com/goplus/llcppg/_xtool/llcppsymg/parse"
	"github.com/goplus/llcppg/_xtool/llcppsymg/symbol"
	"github.com/goplus/llcppg/llcppg"
	"github.com/goplus/llgo/xtool/nm"
)

func main() {
	TestParseHeaderFile()
}
func TestParseHeaderFile() {
	testCases := []struct {
		name         string
		path         string
		dylibSymbols []string
	}{
		{
			name: "inireader",
			path: "./inireader",
			dylibSymbols: []string{
				"ZN9INIReaderC1EPKc",
				"ZN9INIReaderC1EPKcl",
				"ZN9INIReaderD1Ev",
				"ZNK9INIReader10ParseErrorEv",
				"ZNK9INIReader3GetEPKcS1_S1_",
			},
		},
		{
			name: "lua",
			path: "./lua",
			dylibSymbols: []string{
				"lua_error",
				"lua_next",
				"lua_concat",
				"lua_stringtonumber",
			},
		},
		{
			name: "cjson",
			path: "./cjson",
			dylibSymbols: []string{
				"cJSON_Print",
				"cJSON_ParseWithLength",
				"cJSON_Delete",
				// mock multiple symbols
				"cJSON_Delete",
			},
		},
		{
			name: "isl",
			path: "./isl",
			dylibSymbols: []string{
				"isl_pw_qpolynomial_get_ctx",
			},
		},
		{
			name: "gpgerror",
			path: "./gpgerror",
			dylibSymbols: []string{
				"gpg_strsource",
				"gpg_strerror_r",
				"gpg_strerror",
			},
		},
	}

	for _, tc := range testCases {
		fmt.Printf("=== Test Case: %s ===\n", tc.name)
		projPath, err := filepath.Abs(tc.path)
		if err != nil {
			fmt.Println("Get Abs Path Error:", err)
		}
		cfgdata, err := os.ReadFile(filepath.Join(projPath, llcppg.LLCPPG_CFG))
		if err != nil {
			fmt.Println("Read Cfg File Error:", err)
		}
		cfg, err := config.GetConf(cfgdata)
		if err != nil {
			fmt.Println("Get Conf Error:", err)
		}
		if err != nil {
			fmt.Println("Read Symb File Error:", err)
		}

		cfg.CFlags = "-I" + projPath
		pkgHfileInfo := config.PkgHfileInfo(cfg.Config, []string{})
		headerSymbolMap, err := parse.ParseHeaderFile(pkgHfileInfo.CurPkgFiles(), cfg.TrimPrefixes, strings.Fields(cfg.CFlags), cfg.Cplusplus, false)
		if err != nil {
			fmt.Println("Error:", err)
		}
		if err != nil {
			fmt.Printf("Failed to create temp file: %v\n", err)
			return
		}

		// trim to nm symbols
		var dylibsymbs []*nm.Symbol
		for _, symb := range tc.dylibSymbols {
			dylibsymbs = append(dylibsymbs, &nm.Symbol{Name: symbol.AddSymbolPrefixUnder(symb, cfg.Cplusplus)})
		}
		// symbolData, err := symbol.GenerateSymTable(dylibsymbs, headerSymbolMap, filepath.Join(projPath, llcppg.LLCPPG_SYMB))
		symbolData, err := symbol.GenerateSymTable(dylibsymbs, headerSymbolMap)
		if err != nil {
			fmt.Println("Error:", err)
		}
		fmt.Println(string(symbolData))
	}
}
