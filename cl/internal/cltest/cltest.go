package cltest

import (
	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/cl/nc"
	"github.com/goplus/llcppg/cl/nc/ncimpl"
	cfg "github.com/goplus/llcppg/cmd/gogensig/config"
	llcppg "github.com/goplus/llcppg/config"
)

type SymbolEntry = cfg.SymbolEntry

func NewConvSym(syms ...SymbolEntry) func(name *ast.Object, mangleName string) (goName string, err error) {
	return fromSymbTable(cfg.CreateSymbolTable(syms))
}

func GetConvSym(symbFile string) func(name *ast.Object, mangleName string) (goName string, err error) {
	if symbFile == "" {
		panic("symbol file not set")
	}
	symbTable, err := cfg.NewSymbolTable(symbFile)
	if err != nil {
		// NOTE(xsw): not a good idea, but make sense in test cases
		return NewConvSym()
	}
	return fromSymbTable(symbTable)
}

func fromSymbTable(symbTable *cfg.SymbolTable) func(name *ast.Object, mangleName string) (goName string, err error) {
	return func(name *ast.Object, mangleName string) (goName string, err error) {
		item, err := symbTable.LookupSymbol(mangleName)
		if err != nil {
			return
		}
		return item.GoName, nil
	}
}

func NC(cfg *llcppg.Config, fileMap map[string]*llcppg.FileInfo, convSym func(name *ast.Object, mangleName string) (goName string, err error)) nc.NodeConverter {
	return &ncimpl.Converter{
		PkgName:        cfg.Name,
		Pubs:           cfg.TypeMap,
		FileMap:        fileMap,
		ConvSym:        convSym,
		TrimPrefixes:   cfg.TrimPrefixes,
		KeepUnderScore: cfg.KeepUnderScore,
	}
}
