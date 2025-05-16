package convert

import (
	"log"
	"strings"

	"github.com/goplus/llcppg/ast"

	cfg "github.com/goplus/llcppg/cmd/gogensig/config"
	llconfig "github.com/goplus/llcppg/config"
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
	PkgName   string
	SymbFile  string // llcppg.symb.json
	CfgFile   string // llcppg.cfg
	OutputDir string

	Pkg *llconfig.Pkg
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
	Pkg    *llconfig.Pkg
	GenPkg *Package
	Conf   *Config
}

func NewConverter(config *Config) (*Converter, error) {
	symbTable, err := cfg.NewSymbolTable(config.SymbFile)
	if err != nil {
		if debugLog {
			log.Printf("Can't get llcppg.symb.json from %s Use empty table\n", config.SymbFile)
		}
		symbTable = cfg.CreateSymbolTable([]cfg.SymbolEntry{})
	}

	// todo: remove this
	conf, err := cfg.GetCppgCfgFromPath(config.CfgFile)
	if err != nil {
		if debugLog {
			log.Printf("Cant get llcppg.cfg from %s Use empty config\n", config.CfgFile)
		}
		conf = llconfig.NewDefault()
	}

	pkg, err := NewPackage(&PackageConfig{
		PkgBase: PkgBase{
			PkgPath: ".",
			Deps:    conf.Deps,
			Pubs:    conf.TypeMap,
		},
		Name:           config.PkgName,
		OutputDir:      config.OutputDir,
		SymbolTable:    symbTable,
		LibCommand:     conf.Libs,
		TrimPrefixes:   conf.TrimPrefixes,
		KeepUnderScore: conf.KeepUnderScore,
	})
	if err != nil {
		return nil, err
	}
	return &Converter{
		GenPkg: pkg,
		Pkg:    config.Pkg,
		Conf:   config,
	}, nil
}

// todo(zzy):throw error
func (p *Converter) Convert() {
	p.Process()
	p.Complete()
}

func (p *Converter) Process() {
	processDecl := func(file string, process func() error) {
		p.setCurFile(file)
		if err := process(); err != nil {
			log.Panicln(err)
		}
	}

	for _, macro := range p.Pkg.File.Macros {
		processDecl(macro.Loc.File, func() error {
			return p.GenPkg.NewMacro(macro)
		})
	}

	for _, decl := range p.Pkg.File.Decls {
		switch decl := decl.(type) {
		case *ast.TypeDecl:
			processDecl(decl.DeclBase.Loc.File, func() error {
				return p.GenPkg.NewTypeDecl(decl)
			})
		case *ast.EnumTypeDecl:
			processDecl(decl.DeclBase.Loc.File, func() error {
				return p.GenPkg.NewEnumTypeDecl(decl)
			})
		case *ast.TypedefDecl:
			processDecl(decl.DeclBase.Loc.File, func() error {
				return p.GenPkg.NewTypedefDecl(decl)
			})
		case *ast.FuncDecl:
			processDecl(decl.DeclBase.Loc.File, func() error {
				return p.GenPkg.NewFuncDecl(decl)
			})
		}
	}
}

func (p *Converter) Complete() {
	err := p.GenPkg.Complete()
	if err != nil {
		log.Panicf("Complete Fail: %v\n", err)
	}
}

func (p *Converter) setCurFile(file string) {
	info, exist := p.Pkg.FileMap[file]
	if !exist {
		var availableFiles []string
		for f := range p.Pkg.FileMap {
			availableFiles = append(availableFiles, f)
		}
		log.Panicf("File %q not found in FileMap. Available files:\n%s",
			file, strings.Join(availableFiles, "\n"))
	}
	p.GenPkg.SetCurFile(NewHeaderFile(file, info.FileType))
}
