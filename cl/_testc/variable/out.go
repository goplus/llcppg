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

// A const-qualified variable is an immutable var: it maps to a Go var linked
// to the same C symbol, exactly like a mutable var.
//
//go:linkname MaxValue C.maxValue
var MaxValue c.Int

//go:linkname Ratio C.ratio
var Ratio c.Double
