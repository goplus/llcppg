package symg

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"

	"github.com/goplus/llcppg/_xtool/internal/config"
	"github.com/goplus/llcppg/_xtool/internal/ld"
	"github.com/goplus/llcppg/_xtool/llcppsymg/internal/flag"
	llcppg "github.com/goplus/llcppg/config"
	"github.com/goplus/llgo/xtool/nm"
)

type dbgFlags = int

var (
	dbgSymbol        bool
	dbgParseIsMethod bool
)

const (
	DbgSymbol        dbgFlags = 1 << iota
	DbgParseIsMethod          //print parse.go isMethod debug log info
	DbgFlagAll       = DbgSymbol | DbgParseIsMethod
)

func SetDebug(flags dbgFlags) {
	dbgSymbol = (flags & DbgSymbol) != 0
	dbgParseIsMethod = (flags & DbgParseIsMethod) != 0
}

type Config struct {
	Libs         string
	CFlags       string
	Includes     []string
	Mix          bool
	TrimPrefixes []string
	SymMap       map[string]string
	IsCpp        bool
}

func Do(conf *Config) error {
	symbols, err := ParseDylibSymbols(conf.Libs)
	if err != nil {
		return err
	}

	pkgHfiles := config.PkgHfileInfo(conf.Includes, strings.Fields(conf.CFlags), conf.Mix)
	if dbgSymbol {
		fmt.Println("interfaces", pkgHfiles.Inters)
		fmt.Println("implements", pkgHfiles.Impls)
		fmt.Println("thirdhfile", pkgHfiles.Thirds)
	}

	headerInfos, err := ParseHeaderFile(pkgHfiles.CurPkgFiles(), conf.TrimPrefixes, strings.Fields(conf.CFlags), conf.SymMap, conf.IsCpp, false)
	if err != nil {
		return err
	}

	commonSymbols := GetCommonSymbols(symbols, headerInfos)

	jsonData, err := json.MarshalIndent(commonSymbols, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(llcppg.LLCPPG_SYMB, jsonData, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// ParseDylibSymbols parses symbols from dynamic libraries specified in the lib string.
// It handles multiple libraries (e.g., -L/opt/homebrew/lib -llua -lm) and returns
// symbols if at least one library is successfully parsed. Errors from inaccessible
// libraries (like standard libs) are logged as warnings.
//
// Returns symbols and nil error if any symbols are found, or nil and error if none found.
func ParseDylibSymbols(lib string) ([]*nm.Symbol, error) {
	if dbgSymbol {
		fmt.Println("ParseDylibSymbols:from", lib)
	}
	sysPaths := ld.GetLibSearchPaths()
	if dbgSymbol {
		fmt.Println("ParseDylibSymbols:sysPaths", sysPaths)
	}

	lbs := flag.ParseLibs(lib)
	if dbgSymbol {
		fmt.Println("ParseDylibSymbols:LibConfig Parse To")
		fmt.Println("libs.Names: ", lbs.Names)
		fmt.Println("libs.Paths: ", lbs.Paths)
	}
	dylibPaths, notFounds, err := lbs.GenDylibPaths(sysPaths)
	if err != nil {
		return nil, fmt.Errorf("failed to generate some dylib paths: %v", err)
	}

	if dbgSymbol {
		fmt.Println("ParseDylibSymbols:dylibPaths", dylibPaths)
		if len(notFounds) > 0 {
			fmt.Println("ParseDylibSymbols:not found libname", notFounds)
		} else {
			fmt.Println("ParseDylibSymbols:every library is found")
		}
	}

	var symbols []*nm.Symbol
	var parseErrors []string

	for _, dylibPath := range dylibPaths {
		if _, err := os.Stat(dylibPath); err != nil {
			if dbgSymbol {
				fmt.Printf("ParseDylibSymbols:Failed to access dylib %s: %v\n", dylibPath, err)
			}
			continue
		}

		args := []string{}
		if runtime.GOOS == "linux" {
			args = append(args, "-D")
		}

		files, err := nm.New("").List(dylibPath, args...)
		if err != nil {
			parseErrors = append(parseErrors, fmt.Sprintf("ParseDylibSymbols:Failed to list symbols in dylib %s: %v", dylibPath, err))
			continue
		}

		for _, file := range files {
			symbols = append(symbols, file.Symbols...)
		}
	}

	if len(symbols) > 0 {
		if dbgSymbol {
			if len(parseErrors) > 0 {
				fmt.Printf("ParseDylibSymbols:Some libraries could not be parsed: %v\n", parseErrors)
			}
			fmt.Println("ParseDylibSymbols:", len(symbols), "symbols")
		}
		return symbols, nil
	}

	return nil, fmt.Errorf("no symbols found in any dylib. Errors: %v", parseErrors)
}

// finds the intersection of symbols from the dynamic library's symbol table and the symbols parsed from header files.
// It returns a list of symbols that can be externally linked.
func GetCommonSymbols(dylibSymbols []*nm.Symbol, headerSymbols map[string]*SymbolInfo) []*llcppg.SymbolInfo {
	var commonSymbols []*llcppg.SymbolInfo
	processedSymbols := make(map[string]bool)

	for _, dylibSym := range dylibSymbols {
		symName := dylibSym.Name
		if runtime.GOOS == "darwin" {
			symName = strings.TrimPrefix(symName, "_")
		}
		if _, ok := processedSymbols[symName]; ok {
			continue
		}
		if symInfo, ok := headerSymbols[symName]; ok {
			symbolInfo := &llcppg.SymbolInfo{
				Mangle: symName,
				CPP:    symInfo.ProtoName,
				Go:     symInfo.GoName,
			}
			commonSymbols = append(commonSymbols, symbolInfo)
			processedSymbols[symName] = true
		}
	}

	sort.Slice(commonSymbols, func(i, j int) bool {
		return commonSymbols[i].Mangle < commonSymbols[j].Mangle
	})

	return commonSymbols
}

// For mutiple os test,the nm output's symbol name is different.
func AddSymbolPrefixUnder(name string, isCpp bool) string {
	prefix := ""
	if runtime.GOOS == "darwin" {
		prefix = prefix + "_"
	}
	if isCpp {
		prefix = prefix + "_"
	}
	return prefix + name
}
