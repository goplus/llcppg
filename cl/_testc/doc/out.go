package foo

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

const LLGoPackage = "link: -L/path/foo -lfoo"

// A documented enum type.
type Color c.Int

const (
// The red primary color.
	Red Color = 0
// The green primary color.
	Green Color = 1
// The blue primary color.
	Blue Color = 2
)

// A documented struct type.
type Point struct {
	X c.Int
	Y c.Int
}

// A documented typedef alias.
type MyInt = c.Int

// Maximum size limit.
const MAX_SIZE = 100

// Default timeout value.
const TIMEOUT = 30

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
