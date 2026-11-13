package foo

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

//go:linkname Normal C.normal
func Normal()

//go:linkname X_leading C._leading
func X_leading()

//go:linkname Global C.global
var Global c.Int

//go:linkname X_leadingGlobal C._leading_global
var X_leadingGlobal c.Int
