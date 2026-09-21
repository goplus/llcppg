package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

type Value struct {
	_xgo_union [1]uint64
}
type Keyword struct {
	_xgo_union [1]uint32
}
type Doubles struct {
	_xgo_union [1]float64
}

// int i
func (p *Value) XGof_ref_i() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}

// double d
func (p *Value) XGof_ref_d() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}

// char * p
func (p *Value) XGof_ref_p() **c.Char {
	return (**c.Char)(unsafe.Pointer(p))
}

// int type
func (p *Keyword) XGof_ref_type() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}

// unsigned int func
func (p *Keyword) XGof_ref_func() *c.Uint {
	return (*c.Uint)(unsafe.Pointer(p))
}

// double x
func (p *Doubles) XGof_ref_x() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}

// double y
func (p *Doubles) XGof_ref_y() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
