package convert

import (
	"errors"
	"log"
	"strings"

	"github.com/goplus/llcppg/ast"
	cfg "github.com/goplus/llcppg/cmd/gogensig/config"
	"github.com/goplus/llcppg/cmd/gogensig/dbg"
	"github.com/goplus/llcppg/cmd/gogensig/visitor"
	"github.com/goplus/llcppg/llcppg"
)

type AstConvert struct {
	*visitor.BaseDocVisitor
	Pkg       *Package
	visitDone func(pkg *Package, incPath string)
}

type Config struct {
	PkgName      string
	SigfetchFile string
	SymbFile     string // llcppg.symb.json
	CfgFile      string // llcppg.cfg
	PubFile      string // llcppg.pub
	OutputDir    string
}

func NewAstConvert(config *Config) (*AstConvert, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}
	p := new(AstConvert)
	p.BaseDocVisitor = visitor.NewBaseDocVisitor(p)
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
	p.Pkg = pkg
	return p, nil
}

func (p *AstConvert) SetVisitDone(fn func(pkg *Package, incPath string)) {
	p.visitDone = fn
}

func (p *AstConvert) WriteLinkFile() {
	p.Pkg.WriteLinkFile()
}

func (p *AstConvert) WritePubFile() {
	p.Pkg.WritePubFile()
}

func (p *AstConvert) VisitFuncDecl(funcDecl *ast.FuncDecl) {
	err := p.Pkg.NewFuncDecl(funcDecl)
	if err != nil {
		if dbg.GetDebugError() {
			log.Printf("NewFuncDecl %s Fail: %s\n", funcDecl.Name.Name, err.Error())
		}
	}
}

func (p *AstConvert) VisitMacro(macro *ast.Macro) {
	err := p.Pkg.NewMacro(macro)
	if err != nil {
		log.Printf("NewMacro %s Fail: %s\n", macro.Name, err.Error())
	}
}

/*
//TODO
func (p *AstConvert) VisitClass(className *ast.Ident, fields *ast.FieldList, typeDecl *ast.TypeDecl) {
	fmt.Printf("visit class %s\n", className.Name)
	p.pkg.NewTypeDecl(typeDecl)
}

func (p *AstConvert) VisitMethod(className *ast.Ident, method *ast.FuncDecl, typeDecl *ast.TypeDecl) {
	fmt.Printf("visit method %s of %s\n", method.Name.Name, className.Name)
}*/

func (p *AstConvert) VisitStruct(structName *ast.Ident, fields *ast.FieldList, typeDecl *ast.TypeDecl) {
	// https://github.com/goplus/llcppg/issues/66 ignore unexpected struct name
	// Union (unnamed at /usr/local/Cellar/msgpack/6.0.2/include/msgpack/object.h:75:9)
	if strings.ContainsAny(structName.Name, ":\\/") {
		if dbg.GetDebugLog() {
			log.Println("structName", structName.Name, "ignored to convert")
		}
		return
	}
	err := p.Pkg.NewTypeDecl(typeDecl)
	if typeDecl.Name == nil {
		log.Printf("NewTypeDecl anonymous struct skipped")
	}
	if err != nil {
		if name := typeDecl.Name; name != nil {
			log.Printf("NewTypeDecl %s Fail: %s\n", name.Name, err.Error())
		}
	}
}

func (p *AstConvert) VisitUnion(unionName *ast.Ident, fields *ast.FieldList, typeDecl *ast.TypeDecl) {
	p.VisitStruct(unionName, fields, typeDecl)
}

func (p *AstConvert) VisitEnumTypeDecl(enumTypeDecl *ast.EnumTypeDecl) {
	err := p.Pkg.NewEnumTypeDecl(enumTypeDecl)
	if err != nil {
		if name := enumTypeDecl.Name; name != nil {
			log.Printf("NewEnumTypeDecl %s Fail: %s\n", name.Name, err.Error())
		} else {
			log.Printf("NewEnumTypeDecl anonymous Fail: %s\n", err.Error())
		}
	}
}

func (p *AstConvert) VisitTypedefDecl(typedefDecl *ast.TypedefDecl) {
	err := p.Pkg.NewTypedefDecl(typedefDecl)
	if err != nil {
		log.Printf("NewTypedefDecl %s Fail: %s\n", typedefDecl.Name.Name, err.Error())
	}
}

func (p *AstConvert) VisitStart(path string, fileType llcppg.FileType) {
	p.Pkg.SetCurFile(&HeaderFile{
		File:     path,
		FileType: fileType,
	})
}

func (p *AstConvert) VisitDone(incPath string) {
	if p.visitDone != nil {
		p.visitDone(p.Pkg, incPath)
	}
}

func (p *AstConvert) WritePkgFiles() {
	err := p.Pkg.WritePkgFiles()
	if err != nil {
		log.Panicf("WritePkgFiles: %v", err)
	}
}

type ConverterConfig struct {
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
	Conf   *ConverterConfig
}

func NewConverter(config *ConverterConfig) (*Converter, error) {
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
		conf = &llcppg.Config{}
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

func (p *Converter) Process() {
	processDecl := func(file string, name *ast.Ident, declType string, process func() error) {
		var declName string
		if name != nil {
			declName = name.Name
		} else {
			declName = "<anonymous>"
		}
		if !p.setCurFile(file) {
			return
		}
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

	err := p.GenPkg.WritePkgFiles()
	if err != nil {
		log.Printf("WritePkgFiles: %v", err)
	}
	_, err = p.GenPkg.WriteLinkFile()
	if err != nil {
		log.Printf("WriteLinkFile: %v", err)
	}
	err = p.GenPkg.WritePubFile()
	if err != nil {
		log.Printf("WritePubFile: %v", err)
	}
}

func (p *Converter) setCurFile(file string) bool {
	info, exist := p.Pkg.FileMap[file]
	if !exist {
		var availableFiles []string
		for f := range p.Pkg.FileMap {
			availableFiles = append(availableFiles, f)
		}
		log.Fatalf("File %q not found in FileMap. Available files:\n%s",
			file, strings.Join(availableFiles, "\n"))
	}
	p.GenPkg.SetCurFile(NewHeaderFile(file, info.FileType))
	return true
}
