package types

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

const XGoPackage = true

type Macro struct {
	X c.Int
}
type MacroAnon struct {
	X c.Int
}
type Point struct {
	X c.Int
	Y c.Int
}
type Anon struct {
	X c.Int
}
type Status c.Uint

const B_OK Status = 0

type Value struct {
	_xgo_union [1]uint32
}
type Container struct {
	Inner _llcppg_anon_0
}
type Cdbl struct {
	V complex128
}
type Handle struct {
}
type Size = c.Ulong
type Count = c.Int
type Alias = Point
type HandlePtr = *Handle
type _llcppg_anon_0 struct {
	X c.Int
}

// llgo:type C
type Callback = func(_llcppg_param1 c.Int)
type Configured = c.Long
type private = c.Int
type Cflt = complex64

func (p *Value) XGof_ref_n() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *Value) XGof_ref_f() *c.Float {
	return (*c.Float)(unsafe.Pointer(p))
}
