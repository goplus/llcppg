package convert

import (
	"errors"
	"log"
	"strings"

	"github.com/goplus/llcppg/ast"
	cfg "github.com/goplus/llcppg/cmd/gogensig/config"
	"github.com/goplus/llcppg/cmd/gogensig/dbg"
	"github.com/goplus/llcppg/cmd/gogensig/visitor"
	cppgtypes "github.com/goplus/llcppg/types"
)

type AstConvert struct {
	*visitor.BaseDocVisitor
	Pkg       *Package
	visitDone func(pkg *Package, incPath string)
}

type AstConvertConfig struct {
	PkgName   string
	SymbFile  string // llcppg.symb.json
	CfgFile   string // llcppg.cfg
	PubFile   string // llcppg.pub
	OutputDir string
}

func NewAstConvert(config *AstConvertConfig) (*AstConvert, error) {
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

func (p *AstConvert) VisitStart(path string, incPath string, isSys bool) {
	inPkgIncPath := false
	incPaths, notFounds, err := p.Pkg.GetIncPaths()
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
	p.Pkg.SetCurFile(&HeaderFile{
		File:         path,
		IncPath:      incPath,
		IsHeaderFile: true,
		InCurPkg:     inPkgIncPath,
		IsSys:        isSys,
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
