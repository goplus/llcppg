package foo

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

const XGoPackage = true
const LLGoPackage = "link: -L/path/foo -lfoo"

//go:linkname F__1 C._Z1fi
func F__1(a c.Int) c.Uint

//go:linkname F__0 C._Z1fv
func F__0()
