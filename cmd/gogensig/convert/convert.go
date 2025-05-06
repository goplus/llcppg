package convert

import (
	"errors"
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
	PubFile   string // llcppg.pub
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
	if config == nil {
		return nil, errors.New("config is nil")
	}
	symbTable, err := cfg.NewSymbolTable(config.SymbFile)
	if err != nil {
		if debugLog {
			log.Printf("Can't get llcppg.symb.json from %s Use empty table\n", config.SymbFile)
		}
		symbTable = cfg.CreateSymbolTable([]cfg.SymbolEntry{})
	}

	conf, err := cfg.GetCppgCfgFromPath(config.CfgFile)
	if err != nil {
		if debugLog {
			log.Printf("Cant get llcppg.cfg from %s Use empty config\n", config.CfgFile)
		}
		conf = llconfig.NewDefault()
	}

	pkg := NewPackage(&PackageConfig{
		PkgBase: PkgBase{
			PkgPath:  ".",
			CppgConf: conf,
			Pubs:     conf.TypeMap,
		},
		Name:        config.PkgName,
		OutputDir:   config.OutputDir,
		SymbolTable: symbTable,
	})
	return &Converter{
		GenPkg: pkg,
		Pkg:    config.Pkg,
		Conf:   config,
	}, nil
}

func (p *Converter) Convert() {
	p.Process()
	p.Write()
	p.Fmt()
	p.Tidy()
}

func (p *Converter) Process() {
	processDecl := func(file string, name *ast.Ident, declType string, process func() error) {
		var declName string
		if name != nil {
			declName = name.Name
		} else {
			declName = "<anonymous>"
		}
		p.setCurFile(file)
		if err := process(); err != nil {
			log.Printf("Convert%s %s Fail: %s", declType, declName, err.Error())
		}
	}

	for _, macro := range p.Pkg.File.Macros {
		processDecl(macro.Loc.File, &ast.Ident{Name: macro.Name}, "Macro", func() error {
			return p.GenPkg.NewMacro(macro)
		})
	}

	for _, decl := range p.Pkg.File.Decls {
		switch decl := decl.(type) {
		case *ast.TypeDecl:
			processDecl(decl.DeclBase.Loc.File, decl.Name, "TypeDecl", func() error {
				return p.GenPkg.NewTypeDecl(decl)
			})
		case *ast.EnumTypeDecl:
			processDecl(decl.DeclBase.Loc.File, decl.Name, "EnumTypeDecl", func() error {
				return p.GenPkg.NewEnumTypeDecl(decl)
			})
		case *ast.TypedefDecl:
			processDecl(decl.DeclBase.Loc.File, decl.Name, "TypedefDecl", func() error {
				return p.GenPkg.NewTypedefDecl(decl)
			})
		case *ast.FuncDecl:
			processDecl(decl.DeclBase.Loc.File, decl.Name, "FuncDecl", func() error {
				return p.GenPkg.NewFuncDecl(decl)
			})
		}
	}
}

func (p *Converter) Write() {
	err := p.GenPkg.WritePkgFiles()
	if err != nil {
		log.Panicf("WritePkgFiles: %v\n", err)
	}
	err = p.GenPkg.WritePubFile()
	if err != nil {
		log.Panicf("WritePubFile: %v\n", err)
	}
	_, err = p.GenPkg.WriteLinkFile()
	if err != nil {
		log.Panicf("WriteLinkFile: %v\n", err)
	}
}

func (p *Converter) Fmt() {
	err := cfg.RunCommand(p.Conf.OutputDir, "go", "fmt", ".")
	if err != nil {
		log.Panicf("go fmt: %v\n", err)
	}
}

func (p *Converter) Tidy() {
	err := cfg.RunCommand(p.Conf.OutputDir, "go", "mod", "tidy")
	if err != nil {
		log.Panicf("go mod tidy: %v\n", err)
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
