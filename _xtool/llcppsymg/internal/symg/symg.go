package symg

import (
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"

	"github.com/goplus/llcppg/_xtool/internal/clangtool"
	"github.com/goplus/llcppg/_xtool/internal/header"
	"github.com/goplus/llcppg/_xtool/internal/ld"
	"github.com/goplus/llcppg/_xtool/internal/symbol"
	llcppg "github.com/goplus/llcppg/config"
	"github.com/goplus/llgo/xtool/nm"
)

type dbgFlags = int

var (
	dbgSymbol bool
)

const (
	DbgSymbol  dbgFlags = 1 << iota
	DbgFlagAll          = DbgSymbol
)

func SetDebug(flags dbgFlags) {
	dbgSymbol = (flags & DbgSymbol) != 0
}

type Config struct {
	Libs         string
	CFlags       string
	Includes     []string
	Mix          bool
	TrimPrefixes []string
	SymMap       map[string]string
	IsCpp        bool
	HeaderOnly   bool
	LibMode      LibMode
}

func Do(conf *Config) (symbolTable []*llcppg.SymbolInfo, err error) {
	pkgHfiles := header.PkgHfileInfo(&header.Config{
		Includes: conf.Includes,
		Args:     strings.Fields(conf.CFlags),
		Mix:      conf.Mix,
	})
	if dbgSymbol {
		fmt.Println("interfaces", pkgHfiles.Inters)
		fmt.Println("implements", pkgHfiles.Impls)
		fmt.Println("thirdhfile", pkgHfiles.Thirds)
	}

	tempFile, err := os.CreateTemp("", "combine*.h")
	if err != nil {
		return
	}
	defer os.Remove(tempFile.Name())
	err = clangtool.ComposeIncludes(conf.Includes, tempFile.Name())
	if err != nil {
		return
	}

	headerInfos, err := ParseHeaderFile(
		tempFile.Name(),
		pkgHfiles.CurPkgFiles(),
		conf.TrimPrefixes,
		strings.Fields(conf.CFlags),
		conf.SymMap, conf.IsCpp,
	)
	if err != nil {
		return
	}

	if conf.HeaderOnly {
		symbolTable = headerInfos.ToSymbolTable()
	} else {
		var symbols []*nm.Symbol
		symbols, err = FetchSymbols(conf.Libs, conf.LibMode)
		if err != nil {
			return
		}
		symbolTable = GetCommonSymbols(symbols, headerInfos)
	}

	sort.Slice(symbolTable, func(i, j int) bool {
		return symbolTable[i].Mangle < symbolTable[j].Mangle
	})
	return
}

// fetchSymbols parses symbols from dynamic libraries specified in the lib string.
// It handles multiple libraries (e.g., -L/opt/homebrew/lib -llua -lm) and returns
// symbols if at least one library is successfully parsed. Errors from inaccessible
// libraries (like standard libs) are logged as warnings.
//
// Returns symbols and nil error if any symbols are found, or nil and error if none found.
func FetchSymbols(lib string, mode LibMode) ([]*nm.Symbol, error) {
	if dbgSymbol {
		fmt.Println("fetchSymbols:from", lib)
	}
	sysPaths := ld.GetLibSearchPaths()
	if dbgSymbol {
		fmt.Println("fetchSymbols:sysPaths", sysPaths)
	}

	lbs := ParseLibs(lib)
	if dbgSymbol {
		fmt.Println("fetchSymbols:LibConfig Parse To")
		fmt.Println("libs.Names: ", lbs.Names)
		fmt.Println("libs.Paths: ", lbs.Paths)
	}

	libFiles, notFounds, err := lbs.Files(sysPaths, mode)
	if err != nil {
		return nil, fmt.Errorf("failed to generate some lib paths: %v", err)
	}

	if dbgSymbol {
		fmt.Println("fetchSymbols:libFiles", libFiles)
		if len(notFounds) > 0 {
			fmt.Println("fetchSymbols:not found libname", notFounds)
		} else {
			fmt.Println("fetchSymbols:every library is found")
		}
	}

	var symbols []*nm.Symbol
	var parseErrors []string

	for _, libFile := range libFiles {
		args := []string{"-g"}
		if runtime.GOOS == "linux" && mode == symbol.ModeDynamic {
			args = append(args, "-D")
		}

		files, err := nm.New("llvm-nm").List(libFile, args...)
		if err != nil {
			parseErrors = append(parseErrors, fmt.Sprintf("fetchSymbols:Failed to list symbols in lib %s: %v", libFile, err))
			continue
		}

		for _, file := range files {
			symbols = append(symbols, file.Symbols...)
		}
	}

	if len(symbols) > 0 {
		if dbgSymbol {
			if len(parseErrors) > 0 {
				fmt.Printf("fetchSymbols:Some libraries could not be parsed: %v\n", parseErrors)
			}
			fmt.Println("fetchSymbols:", len(symbols), "symbols")
		}
		return symbols, nil
	}

	return nil, fmt.Errorf("no symbols found in any lib. Errors: %v", parseErrors)
}

// todo(zzy):only public for test,when llgo test support private package test,this function should be private
// GetCommonSymbols finds the intersection of symbols from the library symbol table and the symbols parsed from header files.
// It returns a list of symbols that can be externally linked.
func GetCommonSymbols(syms []*nm.Symbol, headerSymbols HeaderSymbols) []*llcppg.SymbolInfo {
	var commonSymbols []*llcppg.SymbolInfo
	processedSymbols := make(map[string]bool)

	for _, sym := range syms {
		symName := sym.Name
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

	return commonSymbols
}
