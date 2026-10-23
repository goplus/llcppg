package foo

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type Object struct {
	V c.Int
}
type ModuleObject struct {
	Object
}

// llgo:link (*Object).IsTrue C.PyObject_IsTrue
func (o *Object) IsTrue() c.Int {
	return 0
}

// llgo:link (*ModuleObject).GetDict C.PyModule_GetDict
func (self *ModuleObject) GetDict() *Object {
	return nil
}

// llgo:link (*ModuleObject).GetName C.PyModule_GetName
func (self *ModuleObject) GetName() *c.Char {
	return nil
}

//go:linkname NewModuleObject C.PyModule_New
func NewModuleObject(name *c.Char) *Object
