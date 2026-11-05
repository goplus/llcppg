package foo

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type SizeT = uint64

//go:linkname Strlen C.strlen
func Strlen(__s *c.Char) SizeT
