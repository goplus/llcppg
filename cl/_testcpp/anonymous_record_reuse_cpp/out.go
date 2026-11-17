package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

const XGoPackage = true

type ReferenceFields struct {
	A *_llcppg_anon_0
	B *_llcppg_anon_0
}
type RvalueReferenceFields struct {
	A *_llcppg_anon_1
	B *_llcppg_anon_1
}
type ScopedAliasFields struct {
	A ScopedAliasFieldsValue
	B ScopedAliasFieldsPointer
	C [2]ScopedAliasFieldsValue
}
type ScopedAliasMembers struct {
	_xgo_union [1]uint64
}
type ScopedTypedefFields struct {
	A ScopedTypedefFieldsValue
	B ScopedTypedefFieldsPointer
}
type ScopedTypedefMembers struct {
	_xgo_union [1]uint64
}
type ScopedEnumFields struct {
	A ScopedEnumFieldsKind
	B ScopedEnumFieldsKind
}
type ScopedEnumMembers struct {
	_xgo_union [1]uint32
}
type LeftValue struct {
	Value c.Int
}
type RightValue struct {
	Value c.Int
}
type NamespaceFields struct {
	A LeftValue
	B RightValue
}
type _llcppg_anon_0 struct {
	Value c.Int
}
type _llcppg_anon_1 struct {
	Value c.Int
}
type ScopedAliasFieldsValue struct {
	Value c.Int
}
type ScopedAliasFieldsPointer = *ScopedAliasFieldsValue
type ScopedAliasMembersValue struct {
	Value c.Int
}
type ScopedAliasMembersPointer = *ScopedAliasMembersValue
type ScopedTypedefFieldsValue struct {
	Value c.Int
}
type ScopedTypedefFieldsPointer = *ScopedTypedefFieldsValue
type ScopedTypedefMembersValue struct {
	Value c.Int
}
type ScopedTypedefMembersPointer = *ScopedTypedefMembersValue
type ScopedEnumFieldsKind c.Uint

const ScopedEnumFieldsKind_Value ScopedEnumFieldsKind = 1

type ScopedEnumMembersKind c.Uint

const ScopedEnumMembersKind_Value ScopedEnumMembersKind = 1

func (p *ScopedAliasMembers) XGof_ref_a() *ScopedAliasMembersValue {
	return (*ScopedAliasMembersValue)(unsafe.Pointer(p))
}
func (p *ScopedAliasMembers) XGof_ref_b() *ScopedAliasMembersPointer {
	return (*ScopedAliasMembersPointer)(unsafe.Pointer(p))
}
func (p *ScopedAliasMembers) XGof_ref_c() *[2]ScopedAliasMembersValue {
	return (*[2]ScopedAliasMembersValue)(unsafe.Pointer(p))
}
func (p *ScopedTypedefMembers) XGof_ref_a() *ScopedTypedefMembersValue {
	return (*ScopedTypedefMembersValue)(unsafe.Pointer(p))
}
func (p *ScopedTypedefMembers) XGof_ref_b() *ScopedTypedefMembersPointer {
	return (*ScopedTypedefMembersPointer)(unsafe.Pointer(p))
}
func (p *ScopedEnumMembers) XGof_ref_a() *ScopedEnumMembersKind {
	return (*ScopedEnumMembersKind)(unsafe.Pointer(p))
}
func (p *ScopedEnumMembers) XGof_ref_b() *ScopedEnumMembersKind {
	return (*ScopedEnumMembersKind)(unsafe.Pointer(p))
}
