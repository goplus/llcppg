package config

import (
	"github.com/goplus/llcppg/ast"
)

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

// Config represents a configuration for the llcppg tool.
type Config struct {
	Name           string            `json:"name"`
	CFlags         string            `json:"cflags"`
	Libs           string            `json:"libs"`
	Include        []string          `json:"include"`
	TrimPrefixes   []string          `json:"trimPrefixes,omitempty"`
	Cplusplus      bool              `json:"cplusplus,omitempty"`
	Deps           []string          `json:"deps,omitempty"`
	KeepUnderScore bool              `json:"keepUnderScore,omitempty"`
	Impl           []ImplFiles       `json:"impl,omitempty"`
	Mix            bool              `json:"mix,omitempty"`
	SymMap         map[string]string `json:"symMap,omitempty"`
	TypeMap        map[string]string `json:"typeMap,omitempty"`
	LibStatic      bool              `json:"libstatic,omitempty"`
}

func NewDefault() *Config {
	return &Config{}
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
	Plat
)

type FileInfo struct {
	FileType FileType
}

type Pkg struct {
	File    *ast.File
	FileMap map[string]*FileInfo
}
