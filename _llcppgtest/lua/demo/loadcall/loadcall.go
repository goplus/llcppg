package main

import (
	_ "unsafe"

	"lua"

	"github.com/goplus/lib/c"
)

func main() {
	L := lua.Newstate__1()
	defer L.Close()

	L.Openlibs()
	if res := L.Loadstring(c.Str("print('hello world')")); res != lua.OK {
		println("error")
	}
	if res := L.Pcallk(0, 0, 0, 0, nil); res != lua.OK {
		println("error")
	}

}

/* Expected output:
hello world
*/
