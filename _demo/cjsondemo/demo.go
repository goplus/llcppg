package main

import (
	"github.com/goplus/llgo/c"
	cjson "github.com/goplus/llpkg/cjson"
)

func main() {
	mod := cjson.CreateObject()
	mod.AddItemToObject(c.Str("hello"), cjson.CreateString(c.Str("llgo")))
	mod.AddItemToObject(c.Str("hello"), cjson.CreateString(c.Str("llcppg")))
	var b cjson.Bool = 1
	mod.AddItemToObject(c.Str("woman"), b.CreateBool())
	cstr := mod.PrintUnformatted()

	c.Printf(c.Str("%s\n"), cstr)
}
