package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

const XGoPackage = true

// Named union documentation.
type Value struct {
	_xgo_union [1]uint32
}

// Typedef union documentation.
type Result struct {
	_xgo_union [1]uint32
}

// Inline union type documentation.
type InlineUnion struct {
	_xgo_union [1]uint32
}

// Empty union documentation.
type Empty struct {
}

func (p *Value) XGof_ref_number() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *Value) XGof_ref_real() *c.Float {
	return (*c.Float)(unsafe.Pointer(p))
}
func (p *Result) XGof_ref_code() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *Result) XGof_ref_value() *c.Float {
	return (*c.Float)(unsafe.Pointer(p))
}
func (p *InlineUnion) XGof_ref_value() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
