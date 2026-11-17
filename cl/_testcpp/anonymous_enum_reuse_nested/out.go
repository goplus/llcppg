package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

const XGoPackage = true

type StructPromotedEnumFieldsType c.Uint

const StructPromotedEnumFieldsValue StructPromotedEnumFieldsType = 1

type StructPromotedEnumFields struct {
	Before c.Int
	_llcppg_anon_0
	After c.Int
}
type UnionPromotedEnumFieldsType c.Uint

const UnionPromotedEnumFieldsValue UnionPromotedEnumFieldsType = 1

type UnionPromotedEnumFields struct {
	Before c.Int
	_llcppg_anon_1
	After c.Int
}
type StructPromotedEnumMembersType c.Uint

const StructPromotedEnumMembersValue StructPromotedEnumMembersType = 1

type StructPromotedEnumMembers struct {
	_xgo_union [2]uint32
}
type UnionPromotedEnumMembersType c.Uint

const UnionPromotedEnumMembersValue UnionPromotedEnumMembersType = 1

type UnionPromotedEnumMembers struct {
	_xgo_union [1]uint32
}
type _llcppg_anon_0 struct {
	A StructPromotedEnumFieldsType
	B StructPromotedEnumFieldsType
}
type _llcppg_anon_1 struct {
	_xgo_union [1]uint32
}
type _llcppg_anon_2 struct {
	A StructPromotedEnumMembersType
	B StructPromotedEnumMembersType
}

func (p *_llcppg_anon_1) XGof_ref_a() *UnionPromotedEnumFieldsType {
	return (*UnionPromotedEnumFieldsType)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_1) XGof_ref_b() *UnionPromotedEnumFieldsType {
	return (*UnionPromotedEnumFieldsType)(unsafe.Pointer(p))
}
func (p *StructPromotedEnumMembers) XGof_ref_before() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *StructPromotedEnumMembers) XGof_ref_a() *StructPromotedEnumMembersType {
	return &(*_llcppg_anon_2)(unsafe.Pointer(p)).A
}
func (p *StructPromotedEnumMembers) XGof_ref_b() *StructPromotedEnumMembersType {
	return &(*_llcppg_anon_2)(unsafe.Pointer(p)).B
}
func (p *StructPromotedEnumMembers) XGof_ref_after() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *UnionPromotedEnumMembers) XGof_ref_before() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *UnionPromotedEnumMembers) XGof_ref_a() *UnionPromotedEnumMembersType {
	return (*UnionPromotedEnumMembersType)(unsafe.Pointer(p))
}
func (p *UnionPromotedEnumMembers) XGof_ref_b() *UnionPromotedEnumMembersType {
	return (*UnionPromotedEnumMembersType)(unsafe.Pointer(p))
}
func (p *UnionPromotedEnumMembers) XGof_ref_after() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
