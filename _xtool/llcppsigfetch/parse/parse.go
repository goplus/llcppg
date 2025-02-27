package parse

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/goplus/llcppg/_xtool/llcppsigfetch/dbg"
	"github.com/goplus/llcppg/_xtool/llcppsymg/clangutils"
	"github.com/goplus/llcppg/_xtool/llcppsymg/config"
	"github.com/goplus/llcppg/llcppg"
	"github.com/goplus/llgo/c/cjson"
)

// temp to avoid call clang in llcppsigfetch,will cause hang
var ClangSearchPath []string
var ClangResourceDir string

type Context struct {
	FileSet []*llcppg.FileEntry
	*ContextConfig
}

type ContextConfig struct {
	Conf        *llcppg.Config
	IncFlags    []string
	PkgFileInfo *config.PkgHfilesInfo
}

func NewContext(cfg *ContextConfig) *Context {
	return &Context{
		FileSet:       make([]*llcppg.FileEntry, 0),
		ContextConfig: cfg,
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

func (p *Context) parseFile(path string) ([]*llcppg.FileEntry, error) {
	if dbg.GetDebugParse() {
		fmt.Fprintln(os.Stderr, "parseFile: path", path)
		fmt.Fprintln(os.Stderr, p.PkgFileInfo != nil)
	}
	converter, err := NewConverter(&clangutils.Config{
		File:  path,
		Temp:  false,
		IsCpp: p.Conf.Cplusplus,
		Args:  strings.Fields(p.Conf.CFlags),
	}, p.PkgFileInfo)
	if err != nil {
		return nil, errors.New("failed to create converter " + path)
	}
	defer converter.Dispose()

	files, err := converter.Convert()

	if err != nil {
		return nil, err
	}

	return files, nil
}

type ParseConfig struct {
	Conf                  *llcppg.Config
	CombinedFile          string
	PreprocessedFile      string
	IncedPreprocessedFile string
	OutputFile            bool
}

func Do(cfg *ParseConfig) (*llcppg.Pkg, error) {
	if err := createTempIfNoExist(&cfg.CombinedFile, cfg.Conf.Name+"*.h"); err != nil {
		return nil, err
	}
	if err := createTempIfNoExist(&cfg.PreprocessedFile, cfg.Conf.Name+"*.i"); err != nil {
		return nil, err
	}
	if err := createTempIfNoExist(&cfg.IncedPreprocessedFile, cfg.Conf.Name+"*.i"); err != nil {
		return nil, err
	}

	if dbg.GetDebugParse() {
		fmt.Fprintln(os.Stderr, "Do: combinedFile", cfg.CombinedFile)
		fmt.Fprintln(os.Stderr, "Do: preprocessedFile", cfg.PreprocessedFile)
		fmt.Fprintln(os.Stderr, "Do: incedPreprocessedFile", cfg.IncedPreprocessedFile)
	}

	// compose includes to a combined file
	err := clangutils.ComposeIncludes(cfg.Conf.Include, cfg.CombinedFile)
	if err != nil {
		return nil, err
	}

	// prepare clang flags to preprocess the combined file
	clangFlags := strings.Fields(cfg.Conf.CFlags)
	clangFlags = append(clangFlags, "-C")  // keep comment
	clangFlags = append(clangFlags, "-dD") // keep macro

	err = clangutils.Preprocess(&clangutils.PreprocessConfig{
		File:    cfg.CombinedFile,
		IsCpp:   cfg.Conf.Cplusplus,
		Args:    clangFlags,
		OutFile: cfg.PreprocessedFile,
	})
	if err != nil {
		return nil, err
	}

	// preprocess the combined file to get the include paths
	incFlags := strings.Fields(cfg.Conf.CFlags)
	incFlags = append(incFlags, "-dI")
	err = clangutils.Preprocess(&clangutils.PreprocessConfig{
		File:    cfg.CombinedFile,
		IsCpp:   cfg.Conf.Cplusplus,
		Args:    incFlags,
		OutFile: cfg.IncedPreprocessedFile,
	})
	if err != nil {
		return nil, err
	}

	// parse the preprocessed file & prepare the same flag as clang's flag
	libclangFlags := []string{}
	if ClangResourceDir != "" {
		libclangFlags = append(libclangFlags, "-resource-dir="+ClangResourceDir, "-I"+path.Join(ClangResourceDir, "include"))
	}
	libclangFlags = append(libclangFlags, strings.Fields(cfg.Conf.CFlags)...)
	converter, err := NewConverterX(
		&Config{
			IncPreprocessedFile: cfg.IncedPreprocessedFile,
			Cfg: &clangutils.Config{
				File:  cfg.PreprocessedFile,
				IsCpp: cfg.Conf.Cplusplus,
				Args:  libclangFlags,
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

func createTempIfNoExist(filename *string, pattern string) error {
	if *filename != "" {
		return nil
	}
	f, err := os.CreateTemp("", pattern)
	if err != nil {
		return err
	}
	*filename = f.Name()
	return nil
}
