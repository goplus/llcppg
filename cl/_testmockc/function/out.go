package foo

import "github.com/goplus/lib/c"

//go:linkname F C.f
func F(a c.Int) c.Uint

//go:linkname X_g C._g
func X_g()

//go:linkname Xprintf C.xprintf
func Xprintf(fmt *c.Char, __llgo_va_list ...any) c.Int
