package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

const XGoPackage = true

type PromotedFields struct {
	_xgo_union [2]uint32
}
type PromotedUnionFields struct {
	_xgo_union [1]uint32
}
type NestedPromotedUnionFields struct {
	_xgo_union [1]uint32
}
type StructNamedStructFields struct {
	Before c.Int
	Nested _llcppg_anon_4
	After  c.Int
}
type StructNamedUnionFields struct {
	Before c.Int
	Nested _llcppg_anon_6
	After  c.Int
}
type StructPromotedStructFields struct {
	Before c.Int
	_llcppg_anon_8
	After c.Int
}
type StructPromotedUnionFields struct {
	Before c.Int
	_llcppg_anon_10
	After c.Int
}
type UnionNamedStructFields struct {
	Before c.Int
	Nested _llcppg_anon_12
	After  c.Int
}
type UnionNamedUnionFields struct {
	Before c.Int
	Nested _llcppg_anon_14
	After  c.Int
}
type UnionPromotedStructFields struct {
	Before c.Int
	_llcppg_anon_16
	After c.Int
}
type UnionPromotedUnionFields struct {
	Before c.Int
	_llcppg_anon_18
	After c.Int
}
type StructNamedStructMembers struct {
	_xgo_union [2]uint32
}
type StructNamedUnionMembers struct {
	_xgo_union [2]uint64
}
type StructPromotedUnionMembers struct {
	_xgo_union [2]uint64
}
type UnionNamedStructMembers struct {
	_xgo_union [1]uint32
}
type UnionNamedUnionMembers struct {
	_xgo_union [1]uint64
}
type UnionPromotedUnionMembers struct {
	_xgo_union [1]uint64
}
type StructPromotedEnumFieldsType c.Uint

const StructPromotedEnumFieldsValue StructPromotedEnumFieldsType = 1

type StructPromotedEnumFields struct {
	Before c.Int
	_llcppg_anon_31
	After c.Int
}
type UnionPromotedEnumFieldsType c.Uint

const UnionPromotedEnumFieldsValue UnionPromotedEnumFieldsType = 1

type UnionPromotedEnumFields struct {
	Before c.Int
	_llcppg_anon_32
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
	A _llcppg_anon_1
	B _llcppg_anon_1
}
type _llcppg_anon_1 struct {
	X c.Int
}
type _llcppg_anon_2 struct {
	X c.Int
}
type _llcppg_anon_3 struct {
	Value c.Int
}
type _llcppg_anon_4 struct {
	A _llcppg_anon_5
	B _llcppg_anon_5
}
type _llcppg_anon_5 struct {
	Value c.Int
}
type _llcppg_anon_6 struct {
	A _llcppg_anon_7
	B _llcppg_anon_7
}
type _llcppg_anon_7 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_8 struct {
	A _llcppg_anon_9
	B _llcppg_anon_9
}
type _llcppg_anon_9 struct {
	Value c.Int
}
type _llcppg_anon_10 struct {
	A _llcppg_anon_11
	B _llcppg_anon_11
}
type _llcppg_anon_11 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_12 struct {
	_xgo_union [1]uint32
}
type _llcppg_anon_13 struct {
	Value c.Int
}
type _llcppg_anon_14 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_15 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_16 struct {
	_xgo_union [1]uint32
}
type _llcppg_anon_17 struct {
	Value c.Int
}
type _llcppg_anon_18 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_19 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_20 struct {
	A _llcppg_anon_21
	B _llcppg_anon_21
}
type _llcppg_anon_21 struct {
	Value c.Int
}
type _llcppg_anon_22 struct {
	A _llcppg_anon_23
	B _llcppg_anon_23
}
type _llcppg_anon_23 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_24 struct {
	A _llcppg_anon_25
	B _llcppg_anon_25
}
type _llcppg_anon_25 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_26 struct {
	_xgo_union [1]uint32
}
type _llcppg_anon_27 struct {
	Value c.Int
}
type _llcppg_anon_28 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_29 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_30 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_31 struct {
	A StructPromotedEnumFieldsType
	B StructPromotedEnumFieldsType
}
type _llcppg_anon_32 struct {
	_xgo_union [1]uint32
}
type _llcppg_anon_33 struct {
	A StructPromotedEnumMembersType
	B StructPromotedEnumMembersType
}

func (p *PromotedFields) XGof_ref_a() *_llcppg_anon_1 {
	return &(*_llcppg_anon_0)(unsafe.Pointer(p)).A
}
func (p *PromotedFields) XGof_ref_b() *_llcppg_anon_1 {
	return &(*_llcppg_anon_0)(unsafe.Pointer(p)).B
}
func (p *PromotedUnionFields) XGof_ref_outer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *PromotedUnionFields) XGof_ref_x() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *PromotedUnionFields) XGof_ref_a() *_llcppg_anon_2 {
	return (*_llcppg_anon_2)(unsafe.Pointer(p))
}
func (p *PromotedUnionFields) XGof_ref_b() *_llcppg_anon_2 {
	return (*_llcppg_anon_2)(unsafe.Pointer(p))
}
func (p *PromotedUnionFields) XGof_ref_tail() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *NestedPromotedUnionFields) XGof_ref_outer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *NestedPromotedUnionFields) XGof_ref_x() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *NestedPromotedUnionFields) XGof_ref_y() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *NestedPromotedUnionFields) XGof_ref_a() *_llcppg_anon_3 {
	return (*_llcppg_anon_3)(unsafe.Pointer(p))
}
func (p *NestedPromotedUnionFields) XGof_ref_b() *_llcppg_anon_3 {
	return (*_llcppg_anon_3)(unsafe.Pointer(p))
}
func (p *NestedPromotedUnionFields) XGof_ref_afterW() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *NestedPromotedUnionFields) XGof_ref_tail() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_7) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_7) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_11) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_11) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_12) XGof_ref_a() *_llcppg_anon_13 {
	return (*_llcppg_anon_13)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_12) XGof_ref_b() *_llcppg_anon_13 {
	return (*_llcppg_anon_13)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_15) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_15) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_14) XGof_ref_a() *_llcppg_anon_15 {
	return (*_llcppg_anon_15)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_14) XGof_ref_b() *_llcppg_anon_15 {
	return (*_llcppg_anon_15)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_16) XGof_ref_a() *_llcppg_anon_17 {
	return (*_llcppg_anon_17)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_16) XGof_ref_b() *_llcppg_anon_17 {
	return (*_llcppg_anon_17)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_19) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_19) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_18) XGof_ref_a() *_llcppg_anon_19 {
	return (*_llcppg_anon_19)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_18) XGof_ref_b() *_llcppg_anon_19 {
	return (*_llcppg_anon_19)(unsafe.Pointer(p))
}
func (p *StructNamedStructMembers) XGof_ref_before() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *StructNamedStructMembers) XGof_ref_nested() *_llcppg_anon_20 {
	return (*_llcppg_anon_20)(unsafe.Pointer(p))
}
func (p *StructNamedStructMembers) XGof_ref_after() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_23) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_23) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *StructNamedUnionMembers) XGof_ref_before() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *StructNamedUnionMembers) XGof_ref_nested() *_llcppg_anon_22 {
	return (*_llcppg_anon_22)(unsafe.Pointer(p))
}
func (p *StructNamedUnionMembers) XGof_ref_after() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_25) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_25) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *StructPromotedUnionMembers) XGof_ref_before() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *StructPromotedUnionMembers) XGof_ref_a() *_llcppg_anon_25 {
	return &(*_llcppg_anon_24)(unsafe.Pointer(p)).A
}
func (p *StructPromotedUnionMembers) XGof_ref_b() *_llcppg_anon_25 {
	return &(*_llcppg_anon_24)(unsafe.Pointer(p)).B
}
func (p *StructPromotedUnionMembers) XGof_ref_after() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_26) XGof_ref_a() *_llcppg_anon_27 {
	return (*_llcppg_anon_27)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_26) XGof_ref_b() *_llcppg_anon_27 {
	return (*_llcppg_anon_27)(unsafe.Pointer(p))
}
func (p *UnionNamedStructMembers) XGof_ref_before() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *UnionNamedStructMembers) XGof_ref_nested() *_llcppg_anon_26 {
	return (*_llcppg_anon_26)(unsafe.Pointer(p))
}
func (p *UnionNamedStructMembers) XGof_ref_after() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_29) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_29) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_28) XGof_ref_a() *_llcppg_anon_29 {
	return (*_llcppg_anon_29)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_28) XGof_ref_b() *_llcppg_anon_29 {
	return (*_llcppg_anon_29)(unsafe.Pointer(p))
}
func (p *UnionNamedUnionMembers) XGof_ref_before() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *UnionNamedUnionMembers) XGof_ref_nested() *_llcppg_anon_28 {
	return (*_llcppg_anon_28)(unsafe.Pointer(p))
}
func (p *UnionNamedUnionMembers) XGof_ref_after() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_30) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_30) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *UnionPromotedUnionMembers) XGof_ref_before() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *UnionPromotedUnionMembers) XGof_ref_a() *_llcppg_anon_30 {
	return (*_llcppg_anon_30)(unsafe.Pointer(p))
}
func (p *UnionPromotedUnionMembers) XGof_ref_b() *_llcppg_anon_30 {
	return (*_llcppg_anon_30)(unsafe.Pointer(p))
}
func (p *UnionPromotedUnionMembers) XGof_ref_after() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_32) XGof_ref_a() *UnionPromotedEnumFieldsType {
	return (*UnionPromotedEnumFieldsType)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_32) XGof_ref_b() *UnionPromotedEnumFieldsType {
	return (*UnionPromotedEnumFieldsType)(unsafe.Pointer(p))
}
func (p *StructPromotedEnumMembers) XGof_ref_before() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *StructPromotedEnumMembers) XGof_ref_a() *StructPromotedEnumMembersType {
	return &(*_llcppg_anon_33)(unsafe.Pointer(p)).A
}
func (p *StructPromotedEnumMembers) XGof_ref_b() *StructPromotedEnumMembersType {
	return &(*_llcppg_anon_33)(unsafe.Pointer(p)).B
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
