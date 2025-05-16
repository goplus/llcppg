package main

import (
	"github.com/goplus/lib/c"
	cjson "github.com/goplus/llpkg/cjson"
)

func main() {
	mod := cjson.Object()
	mod.AddItemToObjectCS(c.Str("hello"), cjson.String(c.Str("llgo")))
	mod.AddItemToObjectCS(c.Str("hello"), cjson.String(c.Str("llcppg")))
	var b cjson.Bool = 1
	mod.AddItemToObjectCS(c.Str("woman"), cjson.CreateBool(b))
	cstr := mod.CStr()

	c.Printf(c.Str("%s\n"), cstr)
}
