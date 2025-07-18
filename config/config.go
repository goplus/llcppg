package config

import (
	"encoding/json"
	"fmt"

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
	Name   string `json:"name"`
	CFlags string `json:"cflags"`
	// NOTE(MeterosLiu): libs can be empty when we're in headerOnly mode
	Libs           string            `json:"libs,omitempty"`
	Include        []string          `json:"include"`
	TrimPrefixes   []string          `json:"trimPrefixes,omitempty"`
	Cplusplus      bool              `json:"cplusplus,omitempty"`
	Deps           []string          `json:"deps,omitempty"`
	KeepUnderScore bool              `json:"keepUnderScore,omitempty"`
	Impl           []ImplFiles       `json:"impl,omitempty"`
	Mix            bool              `json:"mix,omitempty"`
	SymMap         map[string]string `json:"symMap,omitempty"`
	TypeMap        map[string]string `json:"typeMap,omitempty"`
	StaticLib      bool              `json:"staticLib,omitempty"`
	HeaderOnly     bool              `json:"headerOnly,omitempty"`
	GenThird       bool              `json:"genThird,omitempty"`
}

// json middleware for validating
func (c *Config) UnmarshalJSON(data []byte) error {
	// create a new type here to avoid unmarshalling infinite loop.
	type newConfig Config

	var config newConfig
	err := json.Unmarshal(data, &config)

	if err != nil {
		return err
	}

	*c = Config(config)

	// do some check

	// when headeronly mode is disabled, libs must not be empty.
	if c.Libs == "" && !c.HeaderOnly {
		return fmt.Errorf("%w: libs must not be empty", ErrConfigError)
	}

	return nil
}

func NewDefault() *Config {
	return &Config{}
}

type SymbolInfo struct {
	Mangle string `json:"mangle"` // C++ Symbol
	CPP    string `json:"c++"`    // C++ function name
	Go     string `json:"go"`     // Go function name
}

// for better debug
func (s *SymbolInfo) String() string {
	if s == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Go: %s CPP: %s Mangle: %s", s.Go, s.CPP, s.Mangle)
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
