package foo

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type Object struct {
	V c.Int
}
type SizeT = c.Ulong
type ModuleObject struct {
	Object
	*byte
	int
	SizeT
}

//go:linkname ModuleNewObject C.PyModule_NewObject
func ModuleNewObject(name *Object) *ModuleObject

//go:linkname ModuleNew C.PyModule_New
func ModuleNew(name *c.Char) *ModuleObject

// llgo:link (*ModuleObject).Dict C.PyModule_GetDict
func (_llcppg_param1 *ModuleObject) Dict() *Object {
	return nil
}

// llgo:link (*ModuleObject).NameObject C.PyModule_GetNameObject
func (_llcppg_param1 *ModuleObject) NameObject() *Object {
	return nil
}

// llgo:link (*ModuleObject).Name C.PyModule_GetName
func (_llcppg_param1 *ModuleObject) Name() *c.Char {
	return nil
}
