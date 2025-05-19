package cl

import (
	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/cl/internal/convert"
	llconfig "github.com/goplus/llcppg/config"
)

var ErrSkip = convert.ErrSkip

const DbgFlagAll = convert.DbgFlagAll

func SetDebug(flag int) {
	convert.SetDebug(flag)
}

func ModInit(deps []string, outputDir string, modulePath string) error {
	return convert.ModInit(deps, outputDir, modulePath)
}

type NodeConverter = convert.NodeConverter

type ConvConfig struct {
	OutputDir string
	PkgPath   string
	PkgName   string
	Pkg       *ast.File
	FileMap   map[string]*llconfig.FileInfo
	ConvSym   func(name *ast.Object, mangleName string) (goName string, err error)
	NodeConv  NodeConverter

	// CfgFile   string // llcppg.cfg
	TypeMap        map[string]string // llcppg.pub
	Deps           []string          // dependent packages
	TrimPrefixes   []string
	Libs           string
	KeepUnderScore bool
}

func Convert(config *ConvConfig) (pkg Package, err error) {
	cvt, err := convert.NewConverter(&convert.Config{
		OutputDir: config.OutputDir,
		PkgPath:   config.PkgPath,
		PkgName:   config.PkgName,
		Pkg:       config.Pkg,
		FileMap:   config.FileMap,
		ConvSym:   config.ConvSym,
		NodeConv:  config.NodeConv,

		TypeMap:        config.TypeMap,
		Deps:           config.Deps,
		TrimPrefixes:   config.TrimPrefixes,
		Libs:           config.Libs,
		KeepUnderScore: config.KeepUnderScore,
	})
	if err != nil {
		return
	}
	cvt.Convert()
	gp := cvt.GenPkg
	return Package{gp.Pkg(), gp.PkgInfo}, nil
}
