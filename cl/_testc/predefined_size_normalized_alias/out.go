package foo

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type SizeT = uint64

//go:linkname Strlen C.strlen
func Strlen(_llcppg_param1 *c.Char) SizeT
