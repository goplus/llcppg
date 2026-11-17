package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

const XGoPackage = true

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
	A InlineStructFieldsInlineStructFieldsType
	B *InlineStructFieldsInlineStructFieldsType
	C [2]InlineStructFieldsInlineStructFieldsType
}
type InlineUnionFields struct {
	A InlineUnionFieldsInlineUnionFieldsType
	B *InlineUnionFieldsInlineUnionFieldsType
	C [2]InlineUnionFieldsInlineUnionFieldsType
}
type InlineEnumFields struct {
	A InlineEnumFieldsInlineEnumFieldsType
	B *InlineEnumFieldsInlineEnumFieldsType
	C [2]InlineEnumFieldsInlineEnumFieldsType
}
type InlineStructMembers struct {
	_xgo_union [1]uint64
}
type InlineUnionMembers struct {
	_xgo_union [2]uint64
}
type InlineEnumMembers struct {
	_xgo_union [1]uint64
}
type StructAlias2 = StructAlias
type UnionAlias2 = UnionAlias
type EnumAlias2 = EnumAlias
type InlineStructFieldsInlineStructFieldsType struct {
	Value c.Int
}
type InlineUnionFieldsInlineUnionFieldsType struct {
	_xgo_union [1]uint64
}
type InlineEnumFieldsInlineEnumFieldsType c.Uint

const InlineEnumFields_InlineEnumFieldsTypeValue InlineEnumFieldsInlineEnumFieldsType = 1

type InlineStructMembersInlineStructMembersType struct {
	Value c.Int
}
type InlineUnionMembersInlineUnionMembersType struct {
	_xgo_union [1]uint64
}
type InlineEnumMembersInlineEnumMembersType c.Uint

const InlineEnumMembers_InlineEnumMembersTypeValue InlineEnumMembersInlineEnumMembersType = 1

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
func (p *InlineUnionFieldsInlineUnionFieldsType) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *InlineUnionFieldsInlineUnionFieldsType) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *InlineStructMembers) XGof_ref_a() *InlineStructMembersInlineStructMembersType {
	return (*InlineStructMembersInlineStructMembersType)(unsafe.Pointer(p))
}
func (p *InlineStructMembers) XGof_ref_b() **InlineStructMembersInlineStructMembersType {
	return (**InlineStructMembersInlineStructMembersType)(unsafe.Pointer(p))
}
func (p *InlineStructMembers) XGof_ref_c() *[2]InlineStructMembersInlineStructMembersType {
	return (*[2]InlineStructMembersInlineStructMembersType)(unsafe.Pointer(p))
}
func (p *InlineUnionMembersInlineUnionMembersType) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *InlineUnionMembersInlineUnionMembersType) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *InlineUnionMembers) XGof_ref_a() *InlineUnionMembersInlineUnionMembersType {
	return (*InlineUnionMembersInlineUnionMembersType)(unsafe.Pointer(p))
}
func (p *InlineUnionMembers) XGof_ref_b() **InlineUnionMembersInlineUnionMembersType {
	return (**InlineUnionMembersInlineUnionMembersType)(unsafe.Pointer(p))
}
func (p *InlineUnionMembers) XGof_ref_c() *[2]InlineUnionMembersInlineUnionMembersType {
	return (*[2]InlineUnionMembersInlineUnionMembersType)(unsafe.Pointer(p))
}
func (p *InlineEnumMembers) XGof_ref_a() *InlineEnumMembersInlineEnumMembersType {
	return (*InlineEnumMembersInlineEnumMembersType)(unsafe.Pointer(p))
}
func (p *InlineEnumMembers) XGof_ref_b() **InlineEnumMembersInlineEnumMembersType {
	return (**InlineEnumMembersInlineEnumMembersType)(unsafe.Pointer(p))
}
func (p *InlineEnumMembers) XGof_ref_c() *[2]InlineEnumMembersInlineEnumMembersType {
	return (*[2]InlineEnumMembersInlineEnumMembersType)(unsafe.Pointer(p))
}
