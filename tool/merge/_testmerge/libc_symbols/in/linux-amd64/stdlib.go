package libc

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

//go:linkname Calloc C.calloc
func Calloc(__nmemb SizeT, __size SizeT) unsafe.Pointer

//go:linkname Realpath C.realpath
func Realpath(__name *c.Char, __resolved *c.Char) *c.Char
