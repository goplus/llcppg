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
type StructNamedEnumFields struct {
	Before c.Int
	Nested _llcppg_anon_8
	After  c.Int
}
type StructPromotedStructFields struct {
	Before c.Int
	_llcppg_anon_10
	After c.Int
}
type StructPromotedUnionFields struct {
	Before c.Int
	_llcppg_anon_12
	After c.Int
}
type StructPromotedEnumFields struct {
	Before c.Int
	_llcppg_anon_14
	After c.Int
}
type UnionNamedStructFields struct {
	Before c.Int
	Nested _llcppg_anon_16
	After  c.Int
}
type UnionNamedUnionFields struct {
	Before c.Int
	Nested _llcppg_anon_18
	After  c.Int
}
type UnionNamedEnumFields struct {
	Before c.Int
	Nested _llcppg_anon_20
	After  c.Int
}
type UnionPromotedStructFields struct {
	Before c.Int
	_llcppg_anon_22
	After c.Int
}
type UnionPromotedUnionFields struct {
	Before c.Int
	_llcppg_anon_24
	After c.Int
}
type UnionPromotedEnumFields struct {
	Before c.Int
	_llcppg_anon_26
	After c.Int
}
type StructNamedStructMembers struct {
	_xgo_union [2]uint32
}
type StructNamedUnionMembers struct {
	_xgo_union [2]uint64
}
type StructNamedEnumMembers struct {
	_xgo_union [2]uint32
}
type StructPromotedUnionMembers struct {
	_xgo_union [2]uint64
}
type StructPromotedEnumMembers struct {
	_xgo_union [2]uint32
}
type UnionNamedStructMembers struct {
	_xgo_union [1]uint32
}
type UnionNamedUnionMembers struct {
	_xgo_union [1]uint64
}
type UnionNamedEnumMembers struct {
	_xgo_union [1]uint32
}
type UnionPromotedUnionMembers struct {
	_xgo_union [1]uint64
}
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
type _llcppg_anon_9 c.Uint

const StructNamedEnumFieldsValue _llcppg_anon_9 = 1

type _llcppg_anon_10 struct {
	A _llcppg_anon_11
	B _llcppg_anon_11
}
type _llcppg_anon_11 struct {
	Value c.Int
}
type _llcppg_anon_12 struct {
	A _llcppg_anon_13
	B _llcppg_anon_13
}
type _llcppg_anon_13 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_14 struct {
	A _llcppg_anon_15
	B _llcppg_anon_15
}
type _llcppg_anon_15 c.Uint

const StructPromotedEnumFieldsValue _llcppg_anon_15 = 1

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
	_xgo_union [1]uint32
}
type _llcppg_anon_21 c.Uint

const UnionNamedEnumFieldsValue _llcppg_anon_21 = 1

type _llcppg_anon_22 struct {
	_xgo_union [1]uint32
}
type _llcppg_anon_23 struct {
	Value c.Int
}
type _llcppg_anon_24 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_25 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_26 struct {
	_xgo_union [1]uint32
}
type _llcppg_anon_27 c.Uint

const UnionPromotedEnumFieldsValue _llcppg_anon_27 = 1

type _llcppg_anon_28 struct {
	A _llcppg_anon_29
	B _llcppg_anon_29
}
type _llcppg_anon_29 struct {
	Value c.Int
}
type _llcppg_anon_30 struct {
	A _llcppg_anon_31
	B _llcppg_anon_31
}
type _llcppg_anon_31 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_32 struct {
	A _llcppg_anon_33
	B _llcppg_anon_33
}
type _llcppg_anon_33 c.Uint

const StructNamedEnumMembersValue _llcppg_anon_33 = 1

type _llcppg_anon_34 struct {
	A _llcppg_anon_35
	B _llcppg_anon_35
}
type _llcppg_anon_35 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_36 struct {
	A _llcppg_anon_37
	B _llcppg_anon_37
}
type _llcppg_anon_37 c.Uint

const StructPromotedEnumMembersValue _llcppg_anon_37 = 1

type _llcppg_anon_38 struct {
	_xgo_union [1]uint32
}
type _llcppg_anon_39 struct {
	Value c.Int
}
type _llcppg_anon_40 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_41 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_42 struct {
	_xgo_union [1]uint32
}
type _llcppg_anon_43 c.Uint

const UnionNamedEnumMembersValue _llcppg_anon_43 = 1

type _llcppg_anon_44 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_45 c.Uint

const UnionPromotedEnumMembersValue _llcppg_anon_45 = 1

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
func (p *_llcppg_anon_13) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_13) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
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
func (p *_llcppg_anon_20) XGof_ref_a() *_llcppg_anon_21 {
	return (*_llcppg_anon_21)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_20) XGof_ref_b() *_llcppg_anon_21 {
	return (*_llcppg_anon_21)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_22) XGof_ref_a() *_llcppg_anon_23 {
	return (*_llcppg_anon_23)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_22) XGof_ref_b() *_llcppg_anon_23 {
	return (*_llcppg_anon_23)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_25) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_25) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_24) XGof_ref_a() *_llcppg_anon_25 {
	return (*_llcppg_anon_25)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_24) XGof_ref_b() *_llcppg_anon_25 {
	return (*_llcppg_anon_25)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_26) XGof_ref_a() *_llcppg_anon_27 {
	return (*_llcppg_anon_27)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_26) XGof_ref_b() *_llcppg_anon_27 {
	return (*_llcppg_anon_27)(unsafe.Pointer(p))
}
func (p *StructNamedStructMembers) XGof_ref_before() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *StructNamedStructMembers) XGof_ref_nested() *_llcppg_anon_28 {
	return (*_llcppg_anon_28)(unsafe.Pointer(p))
}
func (p *StructNamedStructMembers) XGof_ref_after() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_31) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_31) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *StructNamedUnionMembers) XGof_ref_before() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *StructNamedUnionMembers) XGof_ref_nested() *_llcppg_anon_30 {
	return (*_llcppg_anon_30)(unsafe.Pointer(p))
}
func (p *StructNamedUnionMembers) XGof_ref_after() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *StructNamedEnumMembers) XGof_ref_before() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *StructNamedEnumMembers) XGof_ref_nested() *_llcppg_anon_32 {
	return (*_llcppg_anon_32)(unsafe.Pointer(p))
}
func (p *StructNamedEnumMembers) XGof_ref_after() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_35) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_35) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *StructPromotedUnionMembers) XGof_ref_before() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *StructPromotedUnionMembers) XGof_ref_a() *_llcppg_anon_35 {
	return &(*_llcppg_anon_34)(unsafe.Pointer(p)).A
}
func (p *StructPromotedUnionMembers) XGof_ref_b() *_llcppg_anon_35 {
	return &(*_llcppg_anon_34)(unsafe.Pointer(p)).B
}
func (p *StructPromotedUnionMembers) XGof_ref_after() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *StructPromotedEnumMembers) XGof_ref_before() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *StructPromotedEnumMembers) XGof_ref_a() *_llcppg_anon_37 {
	return &(*_llcppg_anon_36)(unsafe.Pointer(p)).A
}
func (p *StructPromotedEnumMembers) XGof_ref_b() *_llcppg_anon_37 {
	return &(*_llcppg_anon_36)(unsafe.Pointer(p)).B
}
func (p *StructPromotedEnumMembers) XGof_ref_after() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_38) XGof_ref_a() *_llcppg_anon_39 {
	return (*_llcppg_anon_39)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_38) XGof_ref_b() *_llcppg_anon_39 {
	return (*_llcppg_anon_39)(unsafe.Pointer(p))
}
func (p *UnionNamedStructMembers) XGof_ref_before() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *UnionNamedStructMembers) XGof_ref_nested() *_llcppg_anon_38 {
	return (*_llcppg_anon_38)(unsafe.Pointer(p))
}
func (p *UnionNamedStructMembers) XGof_ref_after() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_41) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_41) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_40) XGof_ref_a() *_llcppg_anon_41 {
	return (*_llcppg_anon_41)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_40) XGof_ref_b() *_llcppg_anon_41 {
	return (*_llcppg_anon_41)(unsafe.Pointer(p))
}
func (p *UnionNamedUnionMembers) XGof_ref_before() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *UnionNamedUnionMembers) XGof_ref_nested() *_llcppg_anon_40 {
	return (*_llcppg_anon_40)(unsafe.Pointer(p))
}
func (p *UnionNamedUnionMembers) XGof_ref_after() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_42) XGof_ref_a() *_llcppg_anon_43 {
	return (*_llcppg_anon_43)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_42) XGof_ref_b() *_llcppg_anon_43 {
	return (*_llcppg_anon_43)(unsafe.Pointer(p))
}
func (p *UnionNamedEnumMembers) XGof_ref_before() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *UnionNamedEnumMembers) XGof_ref_nested() *_llcppg_anon_42 {
	return (*_llcppg_anon_42)(unsafe.Pointer(p))
}
func (p *UnionNamedEnumMembers) XGof_ref_after() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_44) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_44) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *UnionPromotedUnionMembers) XGof_ref_before() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *UnionPromotedUnionMembers) XGof_ref_a() *_llcppg_anon_44 {
	return (*_llcppg_anon_44)(unsafe.Pointer(p))
}
func (p *UnionPromotedUnionMembers) XGof_ref_b() *_llcppg_anon_44 {
	return (*_llcppg_anon_44)(unsafe.Pointer(p))
}
func (p *UnionPromotedUnionMembers) XGof_ref_after() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *UnionPromotedEnumMembers) XGof_ref_before() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *UnionPromotedEnumMembers) XGof_ref_a() *_llcppg_anon_45 {
	return (*_llcppg_anon_45)(unsafe.Pointer(p))
}
func (p *UnionPromotedEnumMembers) XGof_ref_b() *_llcppg_anon_45 {
	return (*_llcppg_anon_45)(unsafe.Pointer(p))
}
func (p *UnionPromotedEnumMembers) XGof_ref_after() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
