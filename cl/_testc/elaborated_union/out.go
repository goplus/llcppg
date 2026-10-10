package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

const XGoPackage = true

type OuterUnion struct {
	_xgo_union [2]uint64
}
type OuterWithLargeInner struct {
	_xgo_union [2]uint64
}

func (p *OuterUnion) XGof_ref_i() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *OuterUnion) XGof_ref_f() *c.Float {
	return (*c.Float)(unsafe.Pointer(p))
}
func (p *OuterUnion) XGof_ref_d() *[2]c.Double {
	return (*[2]c.Double)(unsafe.Pointer(p))
}

type _llcppg_anon_0 struct {
	_xgo_union [1]uint32
}

func (p *OuterUnion) XGof_ref_inner() *_llcppg_anon_0 {
	return (*_llcppg_anon_0)(unsafe.Pointer(p))
}
func (p *OuterWithLargeInner) XGof_ref_i() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}

type _llcppg_anon_1 struct {
	_xgo_union [2]uint64
}

func (p *OuterWithLargeInner) XGof_ref_inner() *_llcppg_anon_1 {
	return (*_llcppg_anon_1)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_0) XGof_ref_c() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_0) XGof_ref_s() *int16 {
	return (*int16)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_1) XGof_ref_d() *[2]c.Double {
	return (*[2]c.Double)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_1) XGof_ref_s() *int16 {
	return (*int16)(unsafe.Pointer(p))
}
