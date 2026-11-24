package libc

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

//go:linkname Realpath C.realpath$DARWIN_EXTSN
func Realpath(_llcppg_param1 *c.Char, _llcppg_param2 *c.Char) *c.Char
