package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

const XGoPackage = true

// A simple C union: three members sharing storage. The alignment is 8 (double),
// the size is 8, so the storage is [1]uint64 (mixed int/float/pointer, so the
// element type is an unsigned integer, not a float). Each accessible member
// gets a XGof_ref_<member> accessor returning a typed pointer into the storage.
type Value struct {
	_xgo_union [1]uint64
}

// A member named like a Go keyword is kept verbatim: the accessor is
// XGof_ref_type, not renamed.
type Keyword struct {
	_xgo_union [1]uint32
}

// An all-float union: every scalar leaf is a double of width 8, so the storage
// element is float64 to keep by-value ABI classification correct.
type Doubles struct {
	_xgo_union [1]float64
}

// Alignment 1: every member is byte-sized, so the storage element is uint8.
type Bytes struct {
	_xgo_union [1]uint8
}

// Alignment 2: the widest member is a short, so the storage element is uint16.
type Shorts struct {
	_xgo_union [1]uint16
}

// All-float, width 4: every scalar leaf is a float of width 4, so the storage
// element is float32.
type Floats struct {
	_xgo_union [1]float32
}

// A non-float leaf defeats the all-float rule even when the alignment is 8: the
// int member is not a float, so the storage falls back to an unsigned integer
// element ([4]uint64 for the 32-byte double[4]) rather than [4]float64. The
// array member also exercises an accessor whose type is a pointer to an array.
type BigFloat struct {
	_xgo_union [4]uint64
}
type Empty struct {
}
type Opaque struct {
}

func (p *Value) XGof_ref_i() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *Value) XGof_ref_d() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *Value) XGof_ref_p() **c.Char {
	return (**c.Char)(unsafe.Pointer(p))
}
func (p *Keyword) XGof_ref_type() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *Keyword) XGof_ref_func() *c.Uint {
	return (*c.Uint)(unsafe.Pointer(p))
}
func (p *Doubles) XGof_ref_x() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *Doubles) XGof_ref_y() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *Bytes) XGof_ref_c() *c.Char {
	return (*c.Char)(unsafe.Pointer(p))
}
func (p *Bytes) XGof_ref_uc() *uint8 {
	return (*uint8)(unsafe.Pointer(p))
}
func (p *Shorts) XGof_ref_s() *int16 {
	return (*int16)(unsafe.Pointer(p))
}
func (p *Shorts) XGof_ref_us() *uint16 {
	return (*uint16)(unsafe.Pointer(p))
}
func (p *Floats) XGof_ref_x() *c.Float {
	return (*c.Float)(unsafe.Pointer(p))
}
func (p *Floats) XGof_ref_y() *c.Float {
	return (*c.Float)(unsafe.Pointer(p))
}
func (p *BigFloat) XGof_ref_i() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *BigFloat) XGof_ref_arr() *[4]c.Double {
	return (*[4]c.Double)(unsafe.Pointer(p))
}
