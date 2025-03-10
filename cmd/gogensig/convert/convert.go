package convert

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/goplus/llcppg/ast"
	cfg "github.com/goplus/llcppg/cmd/gogensig/config"
	"github.com/goplus/llcppg/cmd/gogensig/dbg"
	"github.com/goplus/llcppg/llcppg"
)

type Config struct {
	PkgName   string
	SymbFile  string // llcppg.symb.json
	CfgFile   string // llcppg.cfg
	PubFile   string // llcppg.pub
	OutputDir string

	Pkg *llcppg.Pkg
}

type Converter struct {
	Pkg    *llcppg.Pkg
	GenPkg *Package
	Conf   *Config
}

func NewConverter(config *Config) (*Converter, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}
	symbTable, err := cfg.NewSymbolTable(config.SymbFile)
	if err != nil {
		if dbg.GetDebugError() {
			log.Printf("Can't get llcppg.symb.json from %s Use empty table\n", config.SymbFile)
		}
		symbTable = cfg.CreateSymbolTable([]cfg.SymbolEntry{})
	}

	conf, err := cfg.GetCppgCfgFromPath(config.CfgFile)
	if err != nil {
		if dbg.GetDebugError() {
			log.Printf("Cant get llcppg.cfg from %s Use empty config\n", config.CfgFile)
		}
		conf = llcppg.NewDefaultConfig()
	}

	pubs, err := cfg.GetPubFromPath(config.PubFile)
	if err != nil {
		return nil, err
	}
	pkg := NewPackage(&PackageConfig{
		PkgBase: PkgBase{
			PkgPath:  ".",
			CppgConf: conf,
			Pubs:     pubs,
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
	cmd := exec.Command("go", "fmt", ".")
	cmd.Dir = p.Conf.OutputDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Panicf("go fmt: %v\n", err)
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
