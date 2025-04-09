package main

import (
	"isl"

	"github.com/goplus/lib/c"
)

func main() {
	ctx := isl.Alloc()
	set1 := ctx.ReadFromStr(c.Str("{ [x] : 0 <= x <= 10 }"))
	set2 := ctx.ReadFromStr(c.Str("{ [x] : 5 <= x <= 15 }"))

	intersec := set1.Copy().Intersect(set2.Copy())

	set1.Dump()
	set2.Dump()
	intersec.Dump()

	set1.Free()
	set2.Free()
	intersec.Free()
	ctx.Free()
}
