package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

const XGoPackage = true

type DirectCallback struct {
	_xgo_union [1]uint64
}

//llgo:cdecl
type EmbeddedCallback struct {
	_llcppg_anon_0
}
type PromotedStructCallback struct {
	_xgo_union [2]uint64
}
type NestedUnionCallback struct {
	_xgo_union [2]uint64
}

//llgo:cdecl
type Callback = func(_llcppg_param1 c.Int) c.Int
type _llcppg_anon_0 struct {
	_xgo_union [1]uint64
}

//llgo:cdecl
type _llcppg_anon_1 struct {
	Tag      c.Int
	Callback func(_llcppg_param1 c.Int) c.Int
}

//llgo:cdecl
type _llcppg_anon_2 struct {
	Tag c.Int
	_llcppg_anon_3
}
type _llcppg_anon_3 struct {
	_xgo_union [1]uint64
}

//llgo:cdecl
type _llcppg_anon_4 func(_llcppg_param1 c.Int) c.Int

func (p *DirectCallback) XGof_ref_callback() *_llcppg_anon_4 {
	return (*_llcppg_anon_4)(unsafe.Pointer(p))
}

//llgo:cdecl
type _llcppg_anon_5 func(_llcppg_param1 c.Int)

func (p *DirectCallback) XGof_ref_notify() *_llcppg_anon_5 {
	return (*_llcppg_anon_5)(unsafe.Pointer(p))
}
func (p *DirectCallback) XGof_ref_named() *Callback {
	return (*Callback)(unsafe.Pointer(p))
}
func (p *DirectCallback) XGof_ref_value() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}

//llgo:cdecl
type _llcppg_anon_6 func(_llcppg_param1 c.Int) c.Int

func (p *_llcppg_anon_0) XGof_ref_callback() *_llcppg_anon_6 {
	return (*_llcppg_anon_6)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_0) XGof_ref_value() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *PromotedStructCallback) XGof_ref_tag() *c.Int {
	return &(*_llcppg_anon_1)(unsafe.Pointer(p)).Tag
}

//llgo:cdecl
type _llcppg_anon_7 func(_llcppg_param1 c.Int) c.Int

func (p *PromotedStructCallback) XGof_ref_callback() *_llcppg_anon_7 {
	return (*_llcppg_anon_7)(&(*_llcppg_anon_1)(unsafe.Pointer(p)).Callback)
}

//llgo:cdecl
type _llcppg_anon_8 func(_llcppg_param1 c.Int) c.Int

func (p *_llcppg_anon_3) XGof_ref_callback() *_llcppg_anon_8 {
	return (*_llcppg_anon_8)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_3) XGof_ref_value() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *NestedUnionCallback) XGof_ref_tag() *c.Int {
	return &(*_llcppg_anon_2)(unsafe.Pointer(p)).Tag
}
func (p *NestedUnionCallback) XGof_ref_callback() *_llcppg_anon_8 {
	return (*_llcppg_anon_2)(unsafe.Pointer(p)).XGof_ref_callback()
}
func (p *NestedUnionCallback) XGof_ref_value() *c.Int {
	return (*_llcppg_anon_2)(unsafe.Pointer(p)).XGof_ref_value()
}
