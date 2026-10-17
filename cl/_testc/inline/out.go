package foo

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

const (
	LLGoPackage    = "link: -L/path/foo -lfoo"
	LLGoFiles      = "-I/path/foo/include: _wrap/llcppg.c"
	LLGoLibVersion = "v1.0.2"
	X_FOO          = 1.2
)

//go:linkname Add C._llcppg_add
func Add(a c.Int, b c.Int) c.Int

//go:linkname Mul C._llcppg_mul
func Mul(a c.Int, b c.Int) c.Int

//go:linkname G C.g
func G()
