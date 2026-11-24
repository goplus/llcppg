package libc

import (
	"unsafe"
)

//go:linkname Calloc C.calloc
func Calloc(__count SizeT, __size SizeT) unsafe.Pointer
