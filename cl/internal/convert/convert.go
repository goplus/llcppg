package convert

import (
	"log"

	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/cl/nc"
	cfg "github.com/goplus/llcppg/cmd/gogensig/config"
)

type dbgFlags = int

var debugLog bool

const (
	DbgLog     dbgFlags = 1 << iota
	DbgFlagAll          = DbgLog
)

func SetDebug(dbgFlags dbgFlags) {
	debugLog = (dbgFlags & DbgLog) != 0
}

type Config struct {
	OutputDir string
	PkgPath   string
	PkgName   string
	Pkg       *ast.File
	NC        nc.NodeConverter

	Deps []string // dependent packages
	Libs string
}

// if modulePath is not empty, init the module by modulePath
func ModInit(deps []string, outputDir string, modulePath string) error {
	var err error
	if modulePath != "" {
		err = cfg.RunCommand(outputDir, "go", "mod", "init", modulePath)
		if err != nil {
			return err
		}
	}

	loadDeps := []string{"github.com/goplus/lib@v0.2.0"}

	for _, dep := range deps {
		_, std := IsDepStd(dep)
		if !std {
			loadDeps = append(loadDeps, dep)
		}
	}
	for _, dep := range loadDeps {
		err = cfg.RunCommand(outputDir, "go", "get", dep)
		if err != nil {
			return err
		}
	}
	return nil
}

type Converter struct {
	Pkg    *ast.File
	GenPkg *Package
	Conf   *Config
	NC     nc.NodeConverter
}

func NewConverter(config *Config) (*Converter, error) {
	pkg, err := NewPackage(config.NC, &PackageConfig{
		PkgBase: PkgBase{
			PkgPath: config.PkgPath,
			Deps:    config.Deps,
			Pubs:    make(map[string]string),
		},
		Name:       config.PkgName,
		OutputDir:  config.OutputDir,
		LibCommand: config.Libs,
	})
	if err != nil {
		return nil, err
	}
	return &Converter{
		Pkg:    config.Pkg,
		GenPkg: pkg,
		Conf:   config,
		NC:     config.NC,
	}, nil
}

// todo(zzy):throw error
func (p *Converter) Convert() {
	p.Process()
	p.Complete()
}

func (p *Converter) Process() {
	pnc := p.NC
	ctx := p.GenPkg
	for _, macro := range p.Pkg.Macros {
		goName, goFile, err := pnc.ConvMacro(macro.Loc.File, macro)
		if err != nil {
			if err == nc.ErrSkip {
				continue
			}
			log.Panicln("ConvMacro:", err)
		}
		ctx.setGoFile(goFile)
		err = ctx.NewMacro(goName, macro)
		if err != nil {
			log.Panicln("NewMacro:", err)
		}
	}

	for _, decl := range p.Pkg.Decls {
		obj := ast.ObjectOf(decl)
		goName, goFile, err := pnc.ConvDecl(obj.Loc.File, decl)
		if err != nil {
			if err == nc.ErrSkip {
				continue
			}
			log.Panicln("ConvDecl:", err)
		}
		ctx.setGoFile(goFile)
		switch decl := decl.(type) {
		case *ast.TypeDecl:
			err = ctx.NewTypeDecl(goName, decl, pnc)
		case *ast.EnumTypeDecl:
			err = ctx.NewEnumTypeDecl(goName, decl, pnc)
		case *ast.TypedefDecl:
			err = ctx.NewTypedefDecl(goName, decl, pnc)
		case *ast.FuncDecl:
			err = ctx.NewFuncDecl(goName, decl)
		}
		if err != nil {
			log.Panicln(err)
		}
	}
}

func (p *Converter) Complete() {
	err := p.GenPkg.Complete()
	if err != nil {
		log.Panicf("Complete Fail: %v\n", err)
	}
}
