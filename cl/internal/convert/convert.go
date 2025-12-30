package convert

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/cl/nc"
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
		err = runCommand(outputDir, "go", "mod", "init", modulePath)
		if err != nil {
			return err
		}
	}

	loadDeps := []string{"github.com/goplus/lib@v0.3.1"}

	for _, dep := range deps {
		_, std := IsDepStd(dep)
		if !std {
			loadDeps = append(loadDeps, dep)
		}
	}
	for _, dep := range loadDeps {
		err = runCommand(outputDir, "go", "get", dep)
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

func (p *Converter) Convert() error {
	err := p.Process()
	if err != nil {
		return err
	}
	return p.Complete()
}

func (p *Converter) Process() error {
	pnc := p.NC
	ctx := p.GenPkg
	for _, macro := range p.Pkg.Macros {
		goName, goFile, err := pnc.ConvMacro(macro.Loc.File, macro)
		if err != nil {
			if err == nc.ErrSkip {
				continue
			}
			return fmt.Errorf("ConvMacro: %w", err)
		}
		ctx.setGoFile(goFile)
		err = ctx.NewMacro(goName, macro)
		if err != nil {
			return err
		}
	}

	for _, decl := range p.Pkg.Decls {
		obj := ast.ObjectOf(decl)
		goName, goFile, err := pnc.ConvDecl(obj.Loc.File, decl)
		if err != nil {
			if err == nc.ErrSkip {
				continue
			}
			return fmt.Errorf("ConvDecl: %w", err)
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
			return err
		}
	}
	return nil
}

func (p *Converter) Complete() error {
	err := p.GenPkg.Complete()
	if err != nil {
		return fmt.Errorf("Complete Fail: %w", err)
	}
	return nil
}

func runCommand(dir, cmdName string, args ...string) error {
	execCmd := exec.Command(cmdName, args...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Dir = dir
	return execCmd.Run()
}
