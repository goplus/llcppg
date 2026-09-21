package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

const XGoPackage = true

type Value struct {
	_xgo_union [1]uint64
}
type Keyword struct {
	_xgo_union [1]uint32
}
type Doubles struct {
	_xgo_union [1]float64
}
type Bytes struct {
	_xgo_union [1]uint8
}
type Shorts struct {
	_xgo_union [1]uint16
}
type Floats struct {
	_xgo_union [1]float32
}
type BigFloat struct {
	_xgo_union [4]uint64
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

// char c
func (p *Bytes) XGof_ref_c() *c.Char {
	return (*c.Char)(unsafe.Pointer(p))
}

// unsigned char uc
func (p *Bytes) XGof_ref_uc() *uint8 {
	return (*uint8)(unsafe.Pointer(p))
}

// short s
func (p *Shorts) XGof_ref_s() *int16 {
	return (*int16)(unsafe.Pointer(p))
}

// unsigned short us
func (p *Shorts) XGof_ref_us() *uint16 {
	return (*uint16)(unsafe.Pointer(p))
}

// float x
func (p *Floats) XGof_ref_x() *c.Float {
	return (*c.Float)(unsafe.Pointer(p))
}

// float y
func (p *Floats) XGof_ref_y() *c.Float {
	return (*c.Float)(unsafe.Pointer(p))
}

// int i
func (p *BigFloat) XGof_ref_i() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}

// double[4] arr
func (p *BigFloat) XGof_ref_arr() *[4]c.Double {
	return (*[4]c.Double)(unsafe.Pointer(p))
}
