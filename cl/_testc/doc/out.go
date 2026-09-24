package foo

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

const LLGoPackage = "link: -L/path/foo -lfoo"

// A documented enum type.
type Color c.Int

const (
	Red   Color = 0
	Green Color = 1
	Blue  Color = 2
)

// A documented struct type.
type Point struct {
	X c.Int
	Y c.Int
}

// A documented typedef alias.
type MyInt = c.Int

// A documented function.
//
// It adds two integers.
//
//go:linkname Add C.add
func Add(a c.Int, b c.Int) c.Int

// A documented global variable.
//
//go:linkname Counter C.counter
var Counter c.Int
