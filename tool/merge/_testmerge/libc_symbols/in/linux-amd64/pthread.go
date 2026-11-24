package libc

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

//go:linkname PthreadCreate C.pthread_create
func PthreadCreate(__newthread *PthreadT, __attr *PthreadAttrT, __start_routine func(_llcppg_param1 unsafe.Pointer) unsafe.Pointer, __arg unsafe.Pointer) c.Int
