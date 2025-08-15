package cltest

import (
	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/cl/nc"
	"github.com/goplus/llcppg/cl/nc/ncimpl"
	llcppg "github.com/goplus/llcppg/config"
)

func NewConvSym(syms ...llcppg.SymbolInfo) func(name *ast.Object, mangleName string) (goName string, err error) {
	return fromSymbTable(llcppg.NewSymTable(syms))
}

func GetConvSym(symbFile string) func(name *ast.Object, mangleName string) (goName string, err error) {
	if symbFile == "" {
		panic("symbol file not set")
	}
	symbTable, err := llcppg.GetSymTableFromFile(symbFile)
	if err != nil {
		// NOTE(xsw): not a good idea, but make sense in test cases
		return NewConvSym()
	}
	return fromSymbTable(symbTable)
}

func fromSymbTable(symbTable *llcppg.SymTable) func(name *ast.Object, mangleName string) (goName string, err error) {
	return func(name *ast.Object, mangleName string) (goName string, err error) {
		item, err := symbTable.LookupSymbol(mangleName)
		if err != nil {
			return
		}
		return item.Go, nil
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
		GenThird:       cfg.GenThird,
	}
}
