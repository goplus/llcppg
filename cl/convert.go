package cl

import (
	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/cl/internal/convert"
	"github.com/goplus/llcppg/cl/nc"
	llconfig "github.com/goplus/llcppg/config"
)

const DbgFlagAll = convert.DbgFlagAll

func SetDebug(flag int) {
	convert.SetDebug(flag)
}

func ModInit(deps []string, outputDir string, modulePath string) error {
	return convert.ModInit(deps, outputDir, modulePath)
}

type ConvConfig struct {
	OutputDir string
	PkgPath   string
	PkgName   string
	Pkg       *ast.File
	FileMap   map[string]*llconfig.FileInfo
	ConvSym   func(name *ast.Object, mangleName string) (goName string, err error)
	NC        nc.NodeConverter

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
		NC:        config.NC,

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
