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
type TypeDeclarations struct {
	_xgo_union [1]uint64
}
type MutexMutexState struct {
	Lock c.Int
}
type NamedMembersNamedStruct struct {
	Value c.Int
	Next  *NamedMembersNamedStruct
}
type NamedMembersNamedUnion struct {
	_xgo_union [1]uint64
}
type TypeDeclarationsEntry struct {
	Value c.Int
}
type TypeDeclarationsDetail struct {
	Code int16
}
type TypeDeclarationsNumber struct {
	_xgo_union [1]uint64
}

func (p *Mutex) XGof_ref_data() *MutexMutexState {
	return (*MutexMutexState)(unsafe.Pointer(p))
}
func (p *Mutex) XGof_ref_align() *c.Long {
	return (*c.Long)(unsafe.Pointer(p))
}
func (p *NamedMembersNamedUnion) XGof_ref_number() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *NamedMembersNamedUnion) XGof_ref_real() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *NamedMembers) XGof_ref_record() *NamedMembersNamedStruct {
	return (*NamedMembersNamedStruct)(unsafe.Pointer(p))
}
func (p *NamedMembers) XGof_ref_pointer() **NamedMembersNamedStruct {
	return (**NamedMembersNamedStruct)(unsafe.Pointer(p))
}
func (p *NamedMembers) XGof_ref_array() *[2]NamedMembersNamedStruct {
	return (*[2]NamedMembersNamedStruct)(unsafe.Pointer(p))
}
func (p *NamedMembers) XGof_ref_value() *NamedMembersNamedUnion {
	return (*NamedMembersNamedUnion)(unsafe.Pointer(p))
}
func (p *NamedMembers) XGof_ref_union_pointer() **NamedMembersNamedUnion {
	return (**NamedMembersNamedUnion)(unsafe.Pointer(p))
}

//go:linkname TypeDeclarationsCount C._ZN16TypeDeclarations5countE
var TypeDeclarationsCount c.Int

func (p *TypeDeclarationsNumber) XGof_ref_number() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *TypeDeclarationsNumber) XGof_ref_real() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *TypeDeclarations) XGof_ref_entry() *TypeDeclarationsEntry {
	return (*TypeDeclarationsEntry)(unsafe.Pointer(p))
}
func (p *TypeDeclarations) XGof_ref_detail() *TypeDeclarationsDetail {
	return (*TypeDeclarationsDetail)(unsafe.Pointer(p))
}
func (p *TypeDeclarations) XGof_ref_number() *TypeDeclarationsNumber {
	return (*TypeDeclarationsNumber)(unsafe.Pointer(p))
}
