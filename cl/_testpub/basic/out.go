package basic

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

const XGoPackage = true

type Point struct {
	X c.Int
	Y c.Int
}
type Record struct {
	Origin Point
	Count  Length
}
type Status c.Uint

const (
	OK    Status = 0
	ERROR Status = 1
)

type Value struct {
	_xgo_union [1]uint32
}
type Handle struct {
}
type Length = c.Ulong
type HandlePtr = *Handle

// llgo:type C
type Callback = func(_llcppg_param1 Length)

func (p *Value) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *Value) XGof_ref_real() *c.Float {
	return (*c.Float)(unsafe.Pointer(p))
}
