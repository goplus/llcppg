package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

const XGoPackage = true

type StructValue struct {
	Value c.Int
}
type StructTypedefFields struct {
	A StructValue
	B StructPointer
	C StructArray
}
type StructTypedefMembers struct {
	_xgo_union [1]uint64
}
type UnionValue struct {
	_xgo_union [1]uint64
}
type UnionTypedefFields struct {
	A UnionValue
	B UnionPointer
	C UnionArray
}
type UnionTypedefMembers struct {
	_xgo_union [2]uint64
}
type EnumValue c.Uint

const EnumValueValue EnumValue = 1

type EnumTypedefFields struct {
	A EnumValue
	B EnumPointer
	C EnumArray
}
type EnumTypedefMembers struct {
	_xgo_union [1]uint64
}
type StructPointer = *StructValue
type StructArray = [2]StructValue
type UnionPointer = *UnionValue
type UnionArray = [2]UnionValue
type EnumPointer = *EnumValue
type EnumArray = [2]EnumValue

func (p *StructTypedefMembers) XGof_ref_a() *StructValue {
	return (*StructValue)(unsafe.Pointer(p))
}
func (p *StructTypedefMembers) XGof_ref_b() *StructPointer {
	return (*StructPointer)(unsafe.Pointer(p))
}
func (p *StructTypedefMembers) XGof_ref_c() *StructArray {
	return (*StructArray)(unsafe.Pointer(p))
}
func (p *UnionValue) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *UnionValue) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *UnionTypedefMembers) XGof_ref_a() *UnionValue {
	return (*UnionValue)(unsafe.Pointer(p))
}
func (p *UnionTypedefMembers) XGof_ref_b() *UnionPointer {
	return (*UnionPointer)(unsafe.Pointer(p))
}
func (p *UnionTypedefMembers) XGof_ref_c() *UnionArray {
	return (*UnionArray)(unsafe.Pointer(p))
}
func (p *EnumTypedefMembers) XGof_ref_a() *EnumValue {
	return (*EnumValue)(unsafe.Pointer(p))
}
func (p *EnumTypedefMembers) XGof_ref_b() *EnumPointer {
	return (*EnumPointer)(unsafe.Pointer(p))
}
func (p *EnumTypedefMembers) XGof_ref_c() *EnumArray {
	return (*EnumArray)(unsafe.Pointer(p))
}
