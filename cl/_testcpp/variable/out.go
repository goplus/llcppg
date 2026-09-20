package foo

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type Config struct {
}

//go:linkname Total C.total
var Total c.Int

//go:linkname Bar_value C._ZN3bar5valueE
var Bar_value c.Int

//go:linkname Config_count C._ZN6Config5countE
var Config_count c.Int

//go:linkname Config_name C._ZN6Config4nameE
var Config_name *c.Char
