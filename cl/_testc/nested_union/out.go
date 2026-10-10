package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

const XGoPackage = true

type Packet struct {
	Kind c.Int
	Data _llcppg_anon_0
}
type _llcppg_anon_0 struct {
	_xgo_union [1]uint32
}
type Anonymous struct {
	_xgo_union [1]uint32
}
type Recursive struct {
	_xgo_union [1]uint32
}
type Mixed struct {
	_xgo_union [1]uint32
}
type Envelope struct {
	Kind c.Int
	_llcppg_anon_1
}
type _llcppg_anon_1 struct {
	_xgo_union [1]uint32
}
type Tagged struct {
	_xgo_union [1]uint32
}
type Alias = Tagged
type References struct {
	_xgo_union [1]uint32
}
type Holder struct {
	Value Alias
}

func (p *_llcppg_anon_0) XGof_ref_i() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}

type _llcppg_anon_2 struct {
	_xgo_union [1]uint32
}

func (p *_llcppg_anon_0) XGof_ref_inner() *_llcppg_anon_2 {
	return (*_llcppg_anon_2)(unsafe.Pointer(p))
}
func (p *Anonymous) XGof_ref_i() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *Anonymous) XGof_ref_s() *int16 {
	return (*int16)(unsafe.Pointer(p))
}
func (p *Anonymous) XGof_ref_f() *c.Float {
	return (*c.Float)(unsafe.Pointer(p))
}
func (p *Recursive) XGof_ref_i() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *Recursive) XGof_ref_s() *int16 {
	return (*int16)(unsafe.Pointer(p))
}
func (p *Recursive) XGof_ref_f() *c.Float {
	return (*c.Float)(unsafe.Pointer(p))
}
func (p *Mixed) XGof_ref_i() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *Mixed) XGof_ref_s() *int16 {
	return (*int16)(unsafe.Pointer(p))
}

type _llcppg_anon_3 struct {
	_xgo_union [1]uint32
}

func (p *Mixed) XGof_ref_inner() *_llcppg_anon_3 {
	return (*_llcppg_anon_3)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_1) XGof_ref_i() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_1) XGof_ref_s() *int16 {
	return (*int16)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_1) XGof_ref_f() *c.Float {
	return (*c.Float)(unsafe.Pointer(p))
}
func (p *Tagged) XGof_ref_i() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *Tagged) XGof_ref_f() *c.Float {
	return (*c.Float)(unsafe.Pointer(p))
}
func (p *References) XGof_ref_tagged() *Tagged {
	return (*Tagged)(unsafe.Pointer(p))
}
func (p *References) XGof_ref_alias() *Alias {
	return (*Alias)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_2) XGof_ref_x() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_2) XGof_ref_s() *int16 {
	return (*int16)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_3) XGof_ref_n() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_3) XGof_ref_f() *c.Float {
	return (*c.Float)(unsafe.Pointer(p))
}
