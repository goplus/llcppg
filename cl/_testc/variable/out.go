package foo

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

//go:linkname Foo C.foo
var Foo c.Int

//go:linkname Name C.name
var Name *c.Char

//go:linkname X_count C._count
var X_count c.Uint
