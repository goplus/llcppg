package foo

import "github.com/goplus/lib/c"

type Foo struct {
	U      _llcppg_anon_0
	Shorts Shorts
}
type Shorts struct {
	S  int16
	Us uint16
}
type _llcppg_anon_0 struct {
	X c.Float
	Y c.Float
}
type Bar struct {
	_llcppg_anon_1
}
type _llcppg_anon_1 struct {
	S  int16
	Us uint16
}
