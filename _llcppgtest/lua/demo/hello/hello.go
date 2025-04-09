package main

import (
	"lua"

	"github.com/goplus/lib/c"
)

func main() {
	L := lua.Newstate__1()
	defer L.Close()
	L.Openlibs()
	if res := L.Loadstring(c.Str("print('hello world')")); res != lua.OK {
		panic("error")
	}
	// MULTRET -> -1
	if res := L.Pcallk(c.Int(0), c.Int(-1), c.Int(0), 0, nil); res != 0 {
		panic("error")
	}
}
