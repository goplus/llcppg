package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

//go:linkname Sort C.sort
func Sort(a unsafe.Pointer, b unsafe.Pointer, elementSize c.Int, count c.Int, cmp func(_llcppg_param1 unsafe.Pointer, _llcppg_param2 unsafe.Pointer) c.Int) c.Int

//go:linkname G C.g
func G()
