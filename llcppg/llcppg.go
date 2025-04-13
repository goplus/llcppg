package llcppg

import (
	"fmt"
	"strings"

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

func NewImplFiles() *ImplFiles {
	return &ImplFiles{Files: []string{}, Cond: Condition{OS: []string{}, Arch: []string{}}}
}

// Config represents a configuration for the llcppg tool.
type Config struct {
	Name           string            `json:"name"`
	CFlags         string            `json:"cflags"`
	Libs           string            `json:"libs"`
	Include        []string          `json:"include"`
	TrimPrefixes   []string          `json:"trimPrefixes"`
	Cplusplus      bool              `json:"cplusplus"`
	Deps           []string          `json:"deps"`
	KeepUnderScore bool              `json:"keepUnderScore"`
	Impl           []ImplFiles       `json:"impl"`
	Mix            bool              `json:"mix"`
	SymMap         map[string]string `json:"symMap"`
	TypeMap        map[string]string `json:"typeMap"`
}

func NewDefaultConfig() *Config {
	cfg := &Config{SymMap: make(map[string]string), TypeMap: make(map[string]string), Impl: []ImplFiles{*NewImplFiles()}}
	return cfg
}

func (p *Config) String() string {
	strBuilder := strings.Builder{}
	strBuilder.WriteString("Name:" + p.Name + "\n")
	strBuilder.WriteString("CFlags:" + p.CFlags + "\n")
	strBuilder.WriteString("Libs:" + p.Libs + "\n")
	strBuilder.WriteString("Include:" + fmt.Sprintln(p.Include))
	strBuilder.WriteString("TrimPrefixes:" + fmt.Sprintln(p.TrimPrefixes))
	strBuilder.WriteString("Cplusplus:" + fmt.Sprintln(p.Cplusplus))
	strBuilder.WriteString("SymMap:" + fmt.Sprintln(p.SymMap))
	return strBuilder.String()
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

type FileInfo struct {
	FileType FileType
}

type Pkg struct {
	File    *ast.File
	FileMap map[string]*FileInfo
}
