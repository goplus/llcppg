package cl

import (
	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/cl/internal/convert"
	"github.com/goplus/llcppg/cl/nc"
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
	NC        nc.NodeConverter

	Deps []string // dependent packages
	Libs string   // $(pkg-config --libs xxx)
}

func Convert(config *ConvConfig) (pkg Package, err error) {
	cvt, err := convert.NewConverter(&convert.Config{
		OutputDir: config.OutputDir,
		PkgPath:   config.PkgPath,
		PkgName:   config.PkgName,
		Pkg:       config.Pkg,
		NC:        config.NC,
		Deps:      config.Deps,
		Libs:      config.Libs,
	})
	if err != nil {
		return
	}
	cvt.Convert()
	gp := cvt.GenPkg
	return Package{gp.Pkg(), gp.PkgInfo}, nil
}
