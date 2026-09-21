package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

const LLGoPackage = "link: -L/path/foo -lfoo"

type Value struct {
	_xgo_union [1]uint64
}
type Variant struct {
	Kind c.Int
	Val  Value
}
type Vec2 struct {
	_xgo_union [2]float64
}

// int i
func (p *Value) XGo_union_i() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}

// double d
func (p *Value) XGo_union_d() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}

// char * s
func (p *Value) XGo_union_s() **c.Char {
	return (**c.Char)(unsafe.Pointer(p))
}

// unsigned char[8] raw
func (p *Value) XGo_union_raw() *[8]uint8 {
	return (*[8]uint8)(unsafe.Pointer(p))
}

// double[2] xy
func (p *Vec2) XGo_union_xy() *[2]c.Double {
	return (*[2]c.Double)(unsafe.Pointer(p))
}

// double v
func (p *Vec2) XGo_union_v() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
