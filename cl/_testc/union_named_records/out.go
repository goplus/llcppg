package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

const XGoPackage = true

type Mutex struct {
	_xgo_union [1]uint64
}
type NamedMembers struct {
	_xgo_union [4]uint64
}
type PromotedStruct struct {
	_xgo_union [2]uint64
}
type PromotedUnion struct {
	_xgo_union [1]uint64
}
type MutexState struct {
	Lock c.Int
}
type NamedStruct struct {
	Value c.Int
	Next  *NamedStruct
}
type NamedUnion struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_0 struct {
	Tag  c.Int
	View ViewState
	Item ViewValue
}
type ViewState struct {
	First  int16
	Second int16
}
type ViewValue struct {
	_xgo_union [1]uint64
}
type PromotedState struct {
	Code c.Int
}

func (p *Mutex) XGof_ref_data() *MutexState {
	return (*MutexState)(unsafe.Pointer(p))
}
func (p *Mutex) XGof_ref_align() *c.Long {
	return (*c.Long)(unsafe.Pointer(p))
}
func (p *NamedUnion) XGof_ref_number() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *NamedUnion) XGof_ref_real() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *NamedMembers) XGof_ref_record() *NamedStruct {
	return (*NamedStruct)(unsafe.Pointer(p))
}
func (p *NamedMembers) XGof_ref_pointer() **NamedStruct {
	return (**NamedStruct)(unsafe.Pointer(p))
}
func (p *NamedMembers) XGof_ref_array() *[2]NamedStruct {
	return (*[2]NamedStruct)(unsafe.Pointer(p))
}
func (p *NamedMembers) XGof_ref_value() *NamedUnion {
	return (*NamedUnion)(unsafe.Pointer(p))
}
func (p *NamedMembers) XGof_ref_union_pointer() **NamedUnion {
	return (**NamedUnion)(unsafe.Pointer(p))
}
func (p *ViewValue) XGof_ref_code() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *ViewValue) XGof_ref_real() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *PromotedStruct) XGof_ref_tag() *c.Int {
	return &(*_llcppg_anon_0)(unsafe.Pointer(p)).Tag
}
func (p *PromotedStruct) XGof_ref_view() *ViewState {
	return &(*_llcppg_anon_0)(unsafe.Pointer(p)).View
}
func (p *PromotedStruct) XGof_ref_item() *ViewValue {
	return &(*_llcppg_anon_0)(unsafe.Pointer(p)).Item
}
func (p *PromotedStruct) XGof_ref_align() *c.Long {
	return (*c.Long)(unsafe.Pointer(p))
}
func (p *PromotedUnion) XGof_ref_value() *PromotedState {
	return (*PromotedState)(unsafe.Pointer(p))
}
func (p *PromotedUnion) XGof_ref_align() *c.Long {
	return (*c.Long)(unsafe.Pointer(p))
}
