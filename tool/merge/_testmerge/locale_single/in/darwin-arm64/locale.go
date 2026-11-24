package libc

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

//go:linkname Setlocale C.setlocale
func Setlocale(_llcppg_param1 c.Int, _llcppg_param2 *c.Char) *c.Char
