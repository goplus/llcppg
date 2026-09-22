package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

const XGoPackage = true

type Foo struct {
	U      _llcppg_union_0
	Shorts Foo_Shorts
}
type Foo_Shorts struct {
	_xgo_union [1]uint16
}
type _llcppg_union_0 struct {
	_xgo_union [1]float32
}
type Bar struct {
	_llcppg_union_1
}
type _llcppg_union_1 struct {
	_xgo_union [1]uint16
}

func (p *Foo_Shorts) XGof_ref_s() *int16 {
	return (*int16)(unsafe.Pointer(p))
}
func (p *Foo_Shorts) XGof_ref_us() *uint16 {
	return (*uint16)(unsafe.Pointer(p))
}
func (p *_llcppg_union_0) XGof_ref_x() *c.Float {
	return (*c.Float)(unsafe.Pointer(p))
}
func (p *_llcppg_union_0) XGof_ref_y() *c.Float {
	return (*c.Float)(unsafe.Pointer(p))
}
func (p *_llcppg_union_1) XGof_ref_s() *int16 {
	return (*int16)(unsafe.Pointer(p))
}
func (p *_llcppg_union_1) XGof_ref_us() *uint16 {
	return (*uint16)(unsafe.Pointer(p))
}
