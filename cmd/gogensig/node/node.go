package node

import (
	"fmt"
	"strings"

	"github.com/goplus/llcppg/_xtool/llcppsymg/tool/name"
	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/cl"
	"github.com/goplus/llcppg/cmd/gogensig/config"
	llconfig "github.com/goplus/llcppg/config"
)

// todo(zzy):a temp abstract,for cl/convert test & gogensig

type NodeConverter struct {
	symbols *cl.ProcessSymbol
	conf    *NodeConverterConfig
}

type NodeConverterConfig struct {
	PkgName      string
	SymbTable    *config.SymbolTable
	FileMap      map[string]*llconfig.FileInfo
	TrimPrefixes []string
	TypeMap      map[string]string

	// todo(zzy):remove this field
	Symbols *cl.ProcessSymbol
}

func NewNodeConverter(cfg *NodeConverterConfig) *NodeConverter {
	var symbols *cl.ProcessSymbol
	if cfg.Symbols == nil {
		symbols = cl.NewProcessSymbol()
	} else {
		symbols = cfg.Symbols
	}
	return &NodeConverter{
		symbols: symbols,
		conf:    cfg,
	}
}

func (c *NodeConverter) ConvDecl(decl ast.Decl) (goName, goFile string, err error) {
	return "", "", nil
}

func (c *NodeConverter) ConvEnumItem(decl *ast.EnumTypeDecl, item *ast.EnumItem) (goName, goFile string, err error) {
	return "", "", nil
}

func (p *NodeConverter) ConvMacro(macro *ast.Macro) (goName, goFile string, err error) {
	node := cl.NewNode(macro.Name, Macro)
	goName, goFile, err = p.Register(macro.Loc, node, p.constName)
	if err != nil {
		return
	}
	return
}

func (p *NodeConverter) Register(loc *ast.Location, node cl.Node, nameMethod NameMethod) (goName string, goFile string, err error) {
	goFile, err = p.goFile(loc.File)
	if err != nil {
		return
	}
	pubName, exist := p.symbols.Lookup(node)
	if exist {
		return pubName, goFile, nil
	}
	goName, _ = p.GetUniqueName(node, nameMethod)
	return goName, goFile, nil
}

type NameMethod func(name string) string

func (p *NodeConverter) goFile(file string) (string, error) {
	info, ok := p.conf.FileMap[file]
	if !ok {
		var availableFiles []string
		for f := range p.conf.FileMap {
			availableFiles = append(availableFiles, f)
		}
		return "", fmt.Errorf("file %q not found in FileMap. Available files:\n%s",
			file, strings.Join(availableFiles, "\n"))
	}
	switch info.FileType {
	case llconfig.Inter:
		return name.HeaderFileToGo(file), nil
	case llconfig.Impl:
		return p.conf.PkgName + "_autogen.go", nil
	default:
		return "", cl.ErrSkip
	}
}

// GetUniqueName generates a unique public name for a given node using the provided name transformation method.
// It ensures the generated name doesn't conflict with existing names by adding a numeric suffix if needed.
//
// Parameters:
//   - node: The node containing the original name to be transformed
//   - nameMethod: Function used to transform the original name (e.g., declName, constName)
//
// Returns:
//   - pubName: The generated unique public name
//   - changed: Whether the generated name differs from the original name
func (p *NodeConverter) GetUniqueName(node cl.Node, nameMethod NameMethod) (pubName string, changed bool) {
	pubName = nameMethod(node.Name())
	uniquePubName := p.symbols.Register(node, pubName)
	return uniquePubName, uniquePubName != node.Name()
}

// which is define in llcppg.cfg/typeMap
func (p *NodeConverter) definedName(name string) (string, bool) {
	definedName, ok := p.conf.TypeMap[name]
	if ok {
		if definedName == "" {
			return name, true
		}
		return definedName, true
	}
	return name, false
}

// transformName handles identifier name conversion following these rules:
// 1. First checks if the name exists in predefined mapping (in typeMap of llcppg.cfg)
// 2. If not in predefined mapping, applies the transform function
// 3. Before applying the transform function, removes specified prefixes (obtained via trimPrefixes)
//
// Parameters:
//   - name: Original C/C++ identifier name
//   - transform: Name transformation function (like names.PubName or names.ExportName)
//
// Returns:
//   - Transformed identifier name
func (p *NodeConverter) transformName(cname string, transform NameMethod) string {
	if definedName, ok := p.definedName(cname); ok {
		return definedName
	}
	return transform(name.RemovePrefixedName(cname, p.conf.TrimPrefixes))
}

// func (p *NodeConverter) declName(cname string) string {
// 	return p.transformName(cname, name.PubName)
// }

func (p *NodeConverter) constName(cname string) string {
	return p.transformName(cname, name.ExportName)
}

const (
	FuncDecl cl.NodeKind = iota + 1
	TypeDecl
	TypedefDecl
	EnumTypeDecl
	EnumItem
	Macro
)
