package foo

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

const (
	LLGoPackage = "link: -L/path/foo -lfoo"
	LLGoFiles   = "-I/path/foo/include: _wrap/llcppg.c"
)

//go:linkname F C.f
func F(a c.Int) c.Uint

//go:linkname X_g C._llcppg__g
func X_g()
