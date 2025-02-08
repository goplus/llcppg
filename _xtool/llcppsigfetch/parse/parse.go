package parse

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/goplus/llcppg/_xtool/llcppsigfetch/dbg"
	"github.com/goplus/llcppg/_xtool/llcppsymg/clangutils"
	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/types"
	"github.com/goplus/llgo/c/cjson"
)

type Context struct {
	FileSet []*ast.FileEntry
	*ContextConfig
}

type ContextConfig struct {
	Conf     *types.Config
	IncFlags []string
}

func NewContext(cfg *ContextConfig) *Context {
	return &Context{
		FileSet: make([]*ast.FileEntry, 0),
		ContextConfig: &ContextConfig{
			Conf:     cfg.Conf,
			IncFlags: cfg.IncFlags,
		},
	}
}

func (p *Context) Output() *cjson.JSON {
	return MarshalFileSet(p.FileSet)
}

// ProcessFiles processes the given files and adds them to the context
func (p *Context) ProcessFiles(files []string) error {
	if dbg.GetDebugParse() {
		fmt.Fprintln(os.Stderr, "ProcessFiles: files", files, "isCpp", p.Conf.Cplusplus)
	}
	for _, file := range files {
		if err := p.processFile(file); err != nil {
			return err
		}
	}
	return nil
}

// parse file and add it to the context,avoid duplicate parsing
func (p *Context) processFile(path string) error {
	if dbg.GetDebugParse() {
		fmt.Fprintln(os.Stderr, "processFile: path", path)
	}
	for _, entry := range p.FileSet {
		if entry.Path == path {
			if dbg.GetDebugParse() {
				fmt.Fprintln(os.Stderr, "processFile: already parsed", path)
			}
			return nil
		}
	}
	parsedFiles, err := p.parseFile(path)
	if err != nil {
		return errors.New("failed to parse file: " + path)
	}

	p.FileSet = append(p.FileSet, parsedFiles...)
	return nil
}

func (p *Context) parseFile(path string) ([]*ast.FileEntry, error) {
	if dbg.GetDebugParse() {
		fmt.Fprintln(os.Stderr, "parseFile: path", path)
	}
	converter, err := NewConverter(&clangutils.Config{
		File:  path,
		Temp:  false,
		IsCpp: p.Conf.Cplusplus,
		Args:  p.IncFlags,
	})
	if err != nil {
		return nil, errors.New("failed to create converter " + path)
	}
	defer converter.Dispose()

	files, err := converter.Convert()

	// the entry file is the first file in the files list
	entryFile := files[0]
	if entryFile.IncPath != "" {
		return nil, errors.New("entry file " + entryFile.Path + " has include path " + entryFile.IncPath)
	}

	for _, include := range p.Conf.Include {
		if strings.Contains(entryFile.Path, include) {
			entryFile.IncPath = include
			break
		}
	}

	if entryFile.IncPath == "" {
		return nil, errors.New("entry file " + entryFile.Path + " is not in include list")
	}

	if err != nil {
		return nil, err
	}

	return files, nil
}

type ParseConfig struct {
	Conf             *types.Config
	CombinedFile     string
	PreprocessedFile string
	OutputFile       bool
}

func Do(cfg *ParseConfig) (*types.Pkg, error) {
	if cfg.CombinedFile == "" {
		combinedFile, err := os.CreateTemp("", cfg.Conf.Name+"*.h")
		if err != nil {
			return nil, err
		}
		defer combinedFile.Close()
		cfg.CombinedFile = combinedFile.Name()
	}

	if cfg.PreprocessedFile == "" {
		preprocessedFile, err := os.CreateTemp("", cfg.Conf.Name+"*.i")
		if err != nil {
			return nil, err
		}
		defer preprocessedFile.Close()
		cfg.PreprocessedFile = preprocessedFile.Name()
	}

	if dbg.GetDebugParse() {
		fmt.Fprintln(os.Stderr, "Do: combinedFile", cfg.CombinedFile)
		fmt.Fprintln(os.Stderr, "Do: preprocessedFile", cfg.PreprocessedFile)
	}
	err := clangutils.ComposeIncludes(cfg.Conf.Include, cfg.CombinedFile)
	if err != nil {
		return nil, err
	}

	flags := strings.Fields(cfg.Conf.CFlags)
	flags = append(flags, "-nobuiltininc") // to avoid libclang & clang different search path
	flags = append(flags, "-C")            // keep comment
	flags = append(flags, "-dD")           // keep macro

	err = clangutils.Preprocess(&clangutils.PreprocessConfig{
		File:    cfg.CombinedFile,
		IsCpp:   cfg.Conf.Cplusplus,
		Args:    flags,
		OutFile: cfg.PreprocessedFile,
	})
	if err != nil {
		return nil, err
	}

	converter, err := NewConverterX(
		&Config{
			CombinedFile: cfg.CombinedFile,
			Cfg: &clangutils.Config{
				File:  cfg.PreprocessedFile,
				IsCpp: cfg.Conf.Cplusplus,
				Args:  strings.Fields(cfg.Conf.CFlags),
			},
		})
	if err != nil {
		return nil, err
	}
	pkg, err := converter.ConvertX()
	if err != nil {
		return nil, err
	}

	return pkg, nil
}
