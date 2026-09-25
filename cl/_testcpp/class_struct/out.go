package foo

import "github.com/goplus/lib/c"

type Foo struct {
	u      _llcppg_struct_0
	shorts FooShorts
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
	_llcppg_struct_2
	_llcppg_struct_3
}
type _llcppg_struct_1 struct {
	S  int16
	Us uint16
}
type _llcppg_struct_2 struct {
	C  c.Char
	Uc uint8
}
type _llcppg_struct_3 struct {
	X c.Float
	Y c.Float
}
