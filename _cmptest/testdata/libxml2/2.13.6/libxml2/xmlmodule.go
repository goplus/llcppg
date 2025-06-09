package libxml2

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type X_xmlModule struct {
	Unused [8]uint8
}
type Module X_xmlModule
type ModulePtr *Module
type ModuleOption c.Int

const (
	MODULE_LAZY  ModuleOption = 1
	MODULE_LOCAL ModuleOption = 2
)

//go:linkname ModuleOpen C.xmlModuleOpen
func ModuleOpen(filename *c.Char, options c.Int) ModulePtr

//go:linkname ModuleSymbol C.xmlModuleSymbol
func ModuleSymbol(module ModulePtr, name *c.Char, result *c.Pointer) c.Int

//go:linkname ModuleClose C.xmlModuleClose
func ModuleClose(module ModulePtr) c.Int

//go:linkname ModuleFree C.xmlModuleFree
func ModuleFree(module ModulePtr) c.Int
