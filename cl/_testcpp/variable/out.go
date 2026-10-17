package foo

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type Config struct {
}

//go:linkname Total C.total
var Total c.Int

// A const-qualified variable is an immutable var: it maps to a Go var linked
// to the same C++ symbol, exactly like a mutable var.
//
//go:linkname MaxValue C.maxValue
var MaxValue c.Int

//go:linkname BarValue C._ZN3bar5valueE
var BarValue c.Int

//go:linkname BarRatio C._ZN3bar5ratioE
var BarRatio c.Double

//go:linkname ConfigCount C._ZN6Config5countE
var ConfigCount c.Int

//go:linkname ConfigName C._ZN6Config4nameE
var ConfigName *c.Char
