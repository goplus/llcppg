package foo

import "github.com/goplus/lib/c"

type Foo struct {
	U      _llcppg_struct_0
	Shorts FooShorts
}
type FooShorts struct {
	S  int16
	Us uint16
}
type _llcppg_struct_0 struct {
	X c.Float
	Y c.Float
}
type Bar struct {
	_llcppg_struct_1
}
type _llcppg_struct_1 struct {
	S  int16
	Us uint16
}
