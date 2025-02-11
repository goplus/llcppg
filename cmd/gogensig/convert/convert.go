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

	Pkg *cppgtypes.Pkg
}

type Converter struct {
	Pkg    *cppgtypes.Pkg
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
		conf = &cppgtypes.Config{}
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

	for _, decl := range p.Pkg.File.Decls {
		switch decl := decl.(type) {
		case *ast.TypeDecl:
			p.setCurFile(decl.DeclBase.Loc.File)
			if err := p.GenPkg.NewTypeDecl(decl); err != nil {
				log.Printf("ConvertTypeDecl %s Fail: %s", decl.Name.Name, err.Error())
			}
		case *ast.EnumTypeDecl:
			p.setCurFile(decl.DeclBase.Loc.File)
			if err := p.GenPkg.NewEnumTypeDecl(decl); err != nil {
				log.Printf("ConvertEnumTyleDecl %s Fail: %s", decl.Name.Name, err.Error())
			}
		case *ast.TypedefDecl:
			p.setCurFile(decl.DeclBase.Loc.File)
			if err := p.GenPkg.NewTypedefDecl(decl); err != nil {
				log.Printf("ConvertTypedefDecl %s Fail: %s", decl.Name.Name, err.Error())
			}
		case *ast.FuncDecl:
			p.setCurFile(decl.DeclBase.Loc.File)
			if err := p.GenPkg.NewFuncDecl(decl); err != nil {
				log.Printf("ConvertFuncDecl %s Fail: %s", decl.Name.Name, err.Error())
			}
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

	if len(p.Pkg.File.Macros) != 0 {
		p.GenPkg.SetCurFile(&HeaderFile{
			File:         p.Conf.PkgName + "_autogen_macros",
			IsHeaderFile: false,
		})
		for _, macro := range p.Pkg.File.Macros {
			if err := p.GenPkg.NewMacro(macro); err != nil {
				log.Printf("NewMacro %s Fail: %s\n", macro.Name, err.Error())
			}
		}
		err = p.GenPkg.WriteMacrosFile()
		if err != nil {
			log.Printf("WriteMacrosFile: %v", err)
		}
	}
}

func (p *Converter) setCurFile(file string) {
	info := p.Pkg.FileMap[file]
	p.GenPkg.SetCurFile(Hfile(p.GenPkg, file, info.IncPath, info.IsSys))
}

func Hfile(pkg *Package, path string, incPath string, isSys bool) *HeaderFile {
	inPkgIncPath := false
	incPaths, notFounds, err := pkg.GetIncPaths()
	if len(notFounds) > 0 {
		log.Println("failed to find some include paths: \n", notFounds)
		if err != nil {
			log.Println("failed to get any include paths: \n", err.Error())
		}
	}
	for _, includePath := range incPaths {
		if includePath == path {
			inPkgIncPath = true
			break
		}
	}
	return &HeaderFile{
		File:         path,
		IncPath:      incPath,
		IsHeaderFile: true,
		InCurPkg:     inPkgIncPath,
		IsSys:        isSys,
	}
}
