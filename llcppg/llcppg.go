package llcppg

import "github.com/goplus/llcppg/ast"

const LLCPPG_CFG = "llcppg.cfg"
const LLCPPG_SYMB = "llcppg.symb.json"
const LLCPPG_SIGFETCH = "llcppg.sigfetch.json"
const LLCPPG_PUB = "llcppg.pub"

type Condition struct {
	OS   []string `json:"os"`
	Arch []string `json:"arch"`
}

type ImplFiles struct {
	Files []string  `json:"files"`
	Cond  Condition `json:"cond"`
}

func NewImplFiles() *ImplFiles {
	return &ImplFiles{Files: []string{}, Cond: Condition{OS: []string{}, Arch: []string{}}}
}

// Config represents a configuration for the llcppg tool.
type Config struct {
	Name           string      `json:"name"`
	CFlags         string      `json:"cflags"`
	Libs           string      `json:"libs"`
	Include        []string    `json:"include"`
	TrimPrefixes   []string    `json:"trimPrefixes"`
	Cplusplus      bool        `json:"cplusplus"`
	Deps           []string    `json:"deps"`
	KeepUnderScore bool        `json:"keepUnderScore"`
	Impl           []ImplFiles `json:"impl"`
	Mix            bool        `json:"mix"`
}

func NewDefaultConfig() *Config {
	return &Config{Impl: []ImplFiles{*NewImplFiles()}}
}

type SymbolInfo struct {
	Mangle string `json:"mangle"` // C++ Symbol
	CPP    string `json:"c++"`    // C++ function name
	Go     string `json:"go"`     // Go function name
}

type FileType uint

const (
	Inter FileType = iota + 1
	Impl
	Third
)

type FileEntry struct {
	Path     string
	FileType FileType
	Doc      *ast.File
}

type FileInfo struct {
	Typ FileType
}

type Pkg struct {
	File    *ast.File
	FileMap map[string]*FileInfo
}
