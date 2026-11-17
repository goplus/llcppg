package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

const XGoPackage = true

type EnumValueFields struct {
	A _llcppg_anon_0
	B _llcppg_anon_0
	C _llcppg_anon_0
}
type EnumConstFields struct {
	A _llcppg_anon_1
	B _llcppg_anon_1
}
type EnumValueMembers struct {
	_xgo_union [1]uint32
}
type EnumConstMembers struct {
	_xgo_union [1]uint32
}
type _llcppg_anon_0 c.Uint

const EnumValueFieldsValue _llcppg_anon_0 = 1

type _llcppg_anon_1 c.Uint

const EnumConstFieldsValue _llcppg_anon_1 = 1

type _llcppg_anon_2 c.Uint

const EnumValueMembersValue _llcppg_anon_2 = 1

type _llcppg_anon_3 c.Uint

const EnumConstMembersValue _llcppg_anon_3 = 1

func (p *EnumValueMembers) XGof_ref_a() *_llcppg_anon_2 {
	return (*_llcppg_anon_2)(unsafe.Pointer(p))
}
func (p *EnumValueMembers) XGof_ref_b() *_llcppg_anon_2 {
	return (*_llcppg_anon_2)(unsafe.Pointer(p))
}
func (p *EnumValueMembers) XGof_ref_c() *_llcppg_anon_2 {
	return (*_llcppg_anon_2)(unsafe.Pointer(p))
}
func (p *EnumConstMembers) XGof_ref_a() *_llcppg_anon_3 {
	return (*_llcppg_anon_3)(unsafe.Pointer(p))
}
func (p *EnumConstMembers) XGof_ref_b() *_llcppg_anon_3 {
	return (*_llcppg_anon_3)(unsafe.Pointer(p))
}
