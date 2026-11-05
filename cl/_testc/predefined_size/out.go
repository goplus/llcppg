package foo

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

//go:linkname Strlen C.strlen
func Strlen(_llcppg_param1 *c.Char) c.SizeT

//go:linkname OrdinarySize C.ordinary_size
func OrdinarySize() c.Ulong
