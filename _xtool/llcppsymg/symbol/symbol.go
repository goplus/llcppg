package symbol

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
	"unsafe"

	"github.com/goplus/lib/c"
	"github.com/goplus/llcppg/_xtool/llcppsymg/config"
	"github.com/goplus/llcppg/_xtool/llcppsymg/config/cfgparse"
	"github.com/goplus/llcppg/_xtool/llcppsymg/dbg"
	"github.com/goplus/llcppg/_xtool/llcppsymg/parse"
	"github.com/goplus/llcppg/_xtool/llcppsymg/syspath"
	"github.com/goplus/llcppg/llcppg"
	llgoc "github.com/goplus/llgo/c"
	"github.com/goplus/llgo/xtool/nm"
	"github.com/goplus/llpkg/cjson"
)

// ParseDylibSymbols parses symbols from dynamic libraries specified in the lib string.
// It handles multiple libraries (e.g., -L/opt/homebrew/lib -llua -lm) and returns
// symbols if at least one library is successfully parsed. Errors from inaccessible
// libraries (like standard libs) are logged as warnings.
//
// Returns symbols and nil error if any symbols are found, or nil and error if none found.
func ParseDylibSymbols(lib string) ([]*nm.Symbol, error) {
	if dbg.GetDebugSymbol() {
		fmt.Println("ParseDylibSymbols:from", lib)
	}
	sysPaths := syspath.GetLibPaths()
	if dbg.GetDebugSymbol() {
		fmt.Println("ParseDylibSymbols:sysPaths", sysPaths)
	}

	lbs := cfgparse.ParseLibs(lib)
	if dbg.GetDebugSymbol() {
		fmt.Println("ParseDylibSymbols:LibConfig Parse To")
		fmt.Println("libs.Names: ", lbs.Names)
		fmt.Println("libs.Paths: ", lbs.Paths)
	}
	dylibPaths, notFounds, err := lbs.GenDylibPaths(sysPaths)
	if err != nil {
		return nil, fmt.Errorf("failed to generate some dylib paths: %v", err)
	}

	if dbg.GetDebugSymbol() {
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
			if dbg.GetDebugSymbol() {
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
		if dbg.GetDebugSymbol() {
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
func GetCommonSymbols(dylibSymbols []*nm.Symbol, headerSymbols map[string]*parse.SymbolInfo) []*llcppg.SymbolInfo {
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

func ReadExistingSymbolTable(fileName string) (map[string]llcppg.SymbolInfo, bool) {
	if _, err := os.Stat(fileName); err != nil {
		return nil, false
	}

	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, false
	}

	parsedJSON := config.ParseBytes(data)
	if parsedJSON == nil {
		return nil, false
	}

	existingSymbols := make(map[string]llcppg.SymbolInfo)
	arraySize := parsedJSON.GetArraySize()

	for i := 0; i < int(arraySize); i++ {
		item := parsedJSON.GetArrayItem(llgoc.Int(i))
		symbol := llcppg.SymbolInfo{
			Mangle: config.GetStringItem(item, "mangle", ""),
			CPP:    config.GetStringItem(item, "c++", ""),
			Go:     config.GetStringItem(item, "go", ""),
		}
		existingSymbols[symbol.Mangle] = symbol
	}

	return existingSymbols, true
}

func GenSymbolTableData(commonSymbols []*llcppg.SymbolInfo, existingSymbols map[string]llcppg.SymbolInfo) ([]byte, error) {
	if len(existingSymbols) > 0 {
		if dbg.GetDebugSymbol() {
			fmt.Println("GenSymbolTableData:generate symbol table with exist symbol table")
		}
		for i := range commonSymbols {
			if existingSymbol, exists := existingSymbols[commonSymbols[i].Mangle]; exists && commonSymbols[i].Go != existingSymbol.Go {
				if dbg.GetDebugSymbol() {
					fmt.Println("symbol", commonSymbols[i].Mangle, "already exist, use exist symbol", existingSymbol.Go)
				}
				commonSymbols[i].Go = existingSymbol.Go
			} else {
				if dbg.GetDebugSymbol() {
					fmt.Println("new symbol", commonSymbols[i].Mangle, "-", commonSymbols[i].CPP, "-", commonSymbols[i].Go)
				}
			}
		}
	} else {
		if dbg.GetDebugSymbol() {
			fmt.Println("GenSymbolTableData:generate symbol table without symbol table")
			for _, symbol := range commonSymbols {
				fmt.Println("new symbol", symbol.Mangle, "-", symbol.CPP, "-", symbol.Go)
			}
		}
	}

	root := cjson.Array()
	defer root.Delete()

	for _, symbol := range commonSymbols {
		item := cjson.Object()
		item.SetItem(c.Str("mangle"), cjson.String(c.AllocaCStr(symbol.Mangle)))
		item.SetItem(c.Str("c++"), cjson.String(c.AllocaCStr(symbol.CPP)))
		item.SetItem(c.Str("go"), cjson.String(c.AllocaCStr(symbol.Go)))
		root.AddItem(item)
	}

	cStr := root.Print()
	if cStr == nil {
		return nil, errors.New("symbol table is empty")
	}
	defer c.Free(unsafe.Pointer(cStr))
	result := []byte(c.GoString(cStr))
	return result, nil
}

func GenerateAndUpdateSymbolTable(symbols []*nm.Symbol, headerInfos map[string]*parse.SymbolInfo, symbFile string) ([]byte, error) {
	commonSymbols := GetCommonSymbols(symbols, headerInfos)
	if dbg.GetDebugSymbol() {
		fmt.Println("GenerateAndUpdateSymbolTable:", len(commonSymbols), "common symbols")
	}

	existSymbols, exist := ReadExistingSymbolTable(symbFile)
	if exist && dbg.GetDebugSymbol() {
		fmt.Println("GenerateAndUpdateSymbolTable:current path have exist symbol table", symbFile)
	}

	symbolData, err := GenSymbolTableData(commonSymbols, existSymbols)
	if err != nil {
		return nil, err
	}

	return symbolData, nil
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
