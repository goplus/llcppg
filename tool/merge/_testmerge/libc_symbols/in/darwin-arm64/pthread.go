package libc

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

//go:linkname PthreadCreate C.pthread_create
func PthreadCreate(_llcppg_param1 *PthreadT, _llcppg_param2 *PthreadAttrT, _llcppg_param3 func(_llcppg_param1 unsafe.Pointer) unsafe.Pointer, _llcppg_param4 unsafe.Pointer) c.Int
