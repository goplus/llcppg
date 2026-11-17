package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

const XGoPackage = true

type AliasFirst struct {
	A c.Int
	B _llcppg_anon_0
	C _llcppg_anon_0
}
type AliasMiddle struct {
	A _llcppg_anon_1
	B c.Int
	C _llcppg_anon_1
}
type TaggedStruct struct {
	Value c.Int
}
type TaggedUnion struct {
	_xgo_union [1]uint64
}
type TaggedEnum c.Uint

const TaggedEnumValue TaggedEnum = 1

type StructAlias struct {
	Value c.Int
}
type UnionAlias struct {
	_xgo_union [1]uint64
}
type EnumAlias c.Uint

const EnumAliasValue EnumAlias = 1

type StructAliasFields struct {
	NamedA   TaggedStruct
	NamedB   *TaggedStruct
	AliasA   StructAlias
	AliasB   StructAlias2
	PointerA *StructAlias
	PointerB *StructAlias
	ArrayA   [2]StructAlias
	ArrayB   [3]StructAlias
}
type UnionAliasFields struct {
	NamedA   TaggedUnion
	NamedB   *TaggedUnion
	AliasA   UnionAlias
	AliasB   UnionAlias2
	PointerA *UnionAlias
	PointerB *UnionAlias
	ArrayA   [2]UnionAlias
	ArrayB   [3]UnionAlias
}
type EnumAliasFields struct {
	NamedA   TaggedEnum
	NamedB   *TaggedEnum
	AliasA   EnumAlias
	AliasB   EnumAlias2
	PointerA *EnumAlias
	PointerB *EnumAlias
	ArrayA   [2]EnumAlias
	ArrayB   [3]EnumAlias
}
type StructAliasMembers struct {
	_xgo_union [2]uint64
}
type UnionAliasMembers struct {
	_xgo_union [3]uint64
}
type EnumAliasMembers struct {
	_xgo_union [2]uint64
}
type InlineStructFields struct {
	A InlineStructFieldsType
	B *InlineStructFieldsType
	C [2]InlineStructFieldsType
}
type InlineUnionFields struct {
	A InlineUnionFieldsType
	B *InlineUnionFieldsType
	C [2]InlineUnionFieldsType
}
type InlineStructMembers struct {
	_xgo_union [1]uint64
}
type InlineUnionMembers struct {
	_xgo_union [2]uint64
}
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
type _llcppg_anon_0 struct {
	X c.Int
}
type _llcppg_anon_1 struct {
	X c.Int
}
type StructAlias2 = StructAlias
type UnionAlias2 = UnionAlias
type EnumAlias2 = EnumAlias
type InlineStructFieldsType struct {
	Value c.Int
}
type InlineUnionFieldsType struct {
	_xgo_union [1]uint64
}
type InlineStructMembersType struct {
	Value c.Int
}
type InlineUnionMembersType struct {
	_xgo_union [1]uint64
}
type StructPointer = *StructValue
type StructArray = [2]StructValue
type UnionPointer = *UnionValue
type UnionArray = [2]UnionValue
type EnumPointer = *EnumValue
type EnumArray = [2]EnumValue

func (p *TaggedUnion) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *TaggedUnion) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *UnionAlias) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *UnionAlias) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *StructAliasMembers) XGof_ref_named_a() *TaggedStruct {
	return (*TaggedStruct)(unsafe.Pointer(p))
}
func (p *StructAliasMembers) XGof_ref_named_b() **TaggedStruct {
	return (**TaggedStruct)(unsafe.Pointer(p))
}
func (p *StructAliasMembers) XGof_ref_alias_a() *StructAlias {
	return (*StructAlias)(unsafe.Pointer(p))
}
func (p *StructAliasMembers) XGof_ref_alias_b() *StructAlias2 {
	return (*StructAlias2)(unsafe.Pointer(p))
}
func (p *StructAliasMembers) XGof_ref_pointer_a() **StructAlias {
	return (**StructAlias)(unsafe.Pointer(p))
}
func (p *StructAliasMembers) XGof_ref_pointer_b() **StructAlias {
	return (**StructAlias)(unsafe.Pointer(p))
}
func (p *StructAliasMembers) XGof_ref_array_a() *[2]StructAlias {
	return (*[2]StructAlias)(unsafe.Pointer(p))
}
func (p *StructAliasMembers) XGof_ref_array_b() *[3]StructAlias {
	return (*[3]StructAlias)(unsafe.Pointer(p))
}
func (p *UnionAliasMembers) XGof_ref_named_a() *TaggedUnion {
	return (*TaggedUnion)(unsafe.Pointer(p))
}
func (p *UnionAliasMembers) XGof_ref_named_b() **TaggedUnion {
	return (**TaggedUnion)(unsafe.Pointer(p))
}
func (p *UnionAliasMembers) XGof_ref_alias_a() *UnionAlias {
	return (*UnionAlias)(unsafe.Pointer(p))
}
func (p *UnionAliasMembers) XGof_ref_alias_b() *UnionAlias2 {
	return (*UnionAlias2)(unsafe.Pointer(p))
}
func (p *UnionAliasMembers) XGof_ref_pointer_a() **UnionAlias {
	return (**UnionAlias)(unsafe.Pointer(p))
}
func (p *UnionAliasMembers) XGof_ref_pointer_b() **UnionAlias {
	return (**UnionAlias)(unsafe.Pointer(p))
}
func (p *UnionAliasMembers) XGof_ref_array_a() *[2]UnionAlias {
	return (*[2]UnionAlias)(unsafe.Pointer(p))
}
func (p *UnionAliasMembers) XGof_ref_array_b() *[3]UnionAlias {
	return (*[3]UnionAlias)(unsafe.Pointer(p))
}
func (p *EnumAliasMembers) XGof_ref_named_a() *TaggedEnum {
	return (*TaggedEnum)(unsafe.Pointer(p))
}
func (p *EnumAliasMembers) XGof_ref_named_b() **TaggedEnum {
	return (**TaggedEnum)(unsafe.Pointer(p))
}
func (p *EnumAliasMembers) XGof_ref_alias_a() *EnumAlias {
	return (*EnumAlias)(unsafe.Pointer(p))
}
func (p *EnumAliasMembers) XGof_ref_alias_b() *EnumAlias2 {
	return (*EnumAlias2)(unsafe.Pointer(p))
}
func (p *EnumAliasMembers) XGof_ref_pointer_a() **EnumAlias {
	return (**EnumAlias)(unsafe.Pointer(p))
}
func (p *EnumAliasMembers) XGof_ref_pointer_b() **EnumAlias {
	return (**EnumAlias)(unsafe.Pointer(p))
}
func (p *EnumAliasMembers) XGof_ref_array_a() *[2]EnumAlias {
	return (*[2]EnumAlias)(unsafe.Pointer(p))
}
func (p *EnumAliasMembers) XGof_ref_array_b() *[3]EnumAlias {
	return (*[3]EnumAlias)(unsafe.Pointer(p))
}
func (p *InlineUnionFieldsType) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *InlineUnionFieldsType) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *InlineStructMembers) XGof_ref_a() *InlineStructMembersType {
	return (*InlineStructMembersType)(unsafe.Pointer(p))
}
func (p *InlineStructMembers) XGof_ref_b() **InlineStructMembersType {
	return (**InlineStructMembersType)(unsafe.Pointer(p))
}
func (p *InlineStructMembers) XGof_ref_c() *[2]InlineStructMembersType {
	return (*[2]InlineStructMembersType)(unsafe.Pointer(p))
}
func (p *InlineUnionMembersType) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *InlineUnionMembersType) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *InlineUnionMembers) XGof_ref_a() *InlineUnionMembersType {
	return (*InlineUnionMembersType)(unsafe.Pointer(p))
}
func (p *InlineUnionMembers) XGof_ref_b() **InlineUnionMembersType {
	return (**InlineUnionMembersType)(unsafe.Pointer(p))
}
func (p *InlineUnionMembers) XGof_ref_c() *[2]InlineUnionMembersType {
	return (*[2]InlineUnionMembersType)(unsafe.Pointer(p))
}
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
