package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

const XGoPackage = true

type StructNamedStructFields struct {
	Before c.Int
	Nested _llcppg_anon_0
	After  c.Int
}
type StructNamedUnionFields struct {
	Before c.Int
	Nested _llcppg_anon_2
	After  c.Int
}
type StructNamedEnumFields struct {
	Before c.Int
	Nested _llcppg_anon_4
	After  c.Int
}
type StructPromotedStructFields struct {
	Before c.Int
	_llcppg_anon_6
	After c.Int
}
type StructPromotedUnionFields struct {
	Before c.Int
	_llcppg_anon_8
	After c.Int
}
type StructPromotedEnumFields struct {
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
type UnionNamedEnumFields struct {
	Before c.Int
	Nested _llcppg_anon_16
	After  c.Int
}
type UnionPromotedStructFields struct {
	Before c.Int
	_llcppg_anon_18
	After c.Int
}
type UnionPromotedUnionFields struct {
	Before c.Int
	_llcppg_anon_20
	After c.Int
}
type UnionPromotedEnumFields struct {
	Before c.Int
	_llcppg_anon_22
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
	Value c.Int
}
type _llcppg_anon_2 struct {
	A _llcppg_anon_3
	B _llcppg_anon_3
}
type _llcppg_anon_3 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_4 struct {
	A _llcppg_anon_5
	B _llcppg_anon_5
}
type _llcppg_anon_5 c.Uint

const StructNamedEnumFieldsValue _llcppg_anon_5 = 1

type _llcppg_anon_6 struct {
	A _llcppg_anon_7
	B _llcppg_anon_7
}
type _llcppg_anon_7 struct {
	Value c.Int
}
type _llcppg_anon_8 struct {
	A _llcppg_anon_9
	B _llcppg_anon_9
}
type _llcppg_anon_9 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_10 struct {
	A _llcppg_anon_11
	B _llcppg_anon_11
}
type _llcppg_anon_11 c.Uint

const StructPromotedEnumFieldsValue _llcppg_anon_11 = 1

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
type _llcppg_anon_17 c.Uint

const UnionNamedEnumFieldsValue _llcppg_anon_17 = 1

type _llcppg_anon_18 struct {
	_xgo_union [1]uint32
}
type _llcppg_anon_19 struct {
	Value c.Int
}
type _llcppg_anon_20 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_21 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_22 struct {
	_xgo_union [1]uint32
}
type _llcppg_anon_23 c.Uint

const UnionPromotedEnumFieldsValue _llcppg_anon_23 = 1

type _llcppg_anon_24 struct {
	A _llcppg_anon_25
	B _llcppg_anon_25
}
type _llcppg_anon_25 struct {
	Value c.Int
}
type _llcppg_anon_26 struct {
	A _llcppg_anon_27
	B _llcppg_anon_27
}
type _llcppg_anon_27 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_28 struct {
	A _llcppg_anon_29
	B _llcppg_anon_29
}
type _llcppg_anon_29 c.Uint

const StructNamedEnumMembersValue _llcppg_anon_29 = 1

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

const StructPromotedEnumMembersValue _llcppg_anon_33 = 1

type _llcppg_anon_34 struct {
	_xgo_union [1]uint32
}
type _llcppg_anon_35 struct {
	Value c.Int
}
type _llcppg_anon_36 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_37 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_38 struct {
	_xgo_union [1]uint32
}
type _llcppg_anon_39 c.Uint

const UnionNamedEnumMembersValue _llcppg_anon_39 = 1

type _llcppg_anon_40 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_41 c.Uint

const UnionPromotedEnumMembersValue _llcppg_anon_41 = 1

func (p *_llcppg_anon_3) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_3) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_9) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_9) XGof_ref_decimal() *c.Double {
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
func (p *_llcppg_anon_18) XGof_ref_a() *_llcppg_anon_19 {
	return (*_llcppg_anon_19)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_18) XGof_ref_b() *_llcppg_anon_19 {
	return (*_llcppg_anon_19)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_21) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_21) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
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
func (p *StructNamedStructMembers) XGof_ref_before() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *StructNamedStructMembers) XGof_ref_nested() *_llcppg_anon_24 {
	return (*_llcppg_anon_24)(unsafe.Pointer(p))
}
func (p *StructNamedStructMembers) XGof_ref_after() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_27) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_27) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *StructNamedUnionMembers) XGof_ref_before() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *StructNamedUnionMembers) XGof_ref_nested() *_llcppg_anon_26 {
	return (*_llcppg_anon_26)(unsafe.Pointer(p))
}
func (p *StructNamedUnionMembers) XGof_ref_after() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *StructNamedEnumMembers) XGof_ref_before() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *StructNamedEnumMembers) XGof_ref_nested() *_llcppg_anon_28 {
	return (*_llcppg_anon_28)(unsafe.Pointer(p))
}
func (p *StructNamedEnumMembers) XGof_ref_after() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_31) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_31) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *StructPromotedUnionMembers) XGof_ref_before() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *StructPromotedUnionMembers) XGof_ref_a() *_llcppg_anon_31 {
	return &(*_llcppg_anon_30)(unsafe.Pointer(p)).A
}
func (p *StructPromotedUnionMembers) XGof_ref_b() *_llcppg_anon_31 {
	return &(*_llcppg_anon_30)(unsafe.Pointer(p)).B
}
func (p *StructPromotedUnionMembers) XGof_ref_after() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *StructPromotedEnumMembers) XGof_ref_before() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *StructPromotedEnumMembers) XGof_ref_a() *_llcppg_anon_33 {
	return &(*_llcppg_anon_32)(unsafe.Pointer(p)).A
}
func (p *StructPromotedEnumMembers) XGof_ref_b() *_llcppg_anon_33 {
	return &(*_llcppg_anon_32)(unsafe.Pointer(p)).B
}
func (p *StructPromotedEnumMembers) XGof_ref_after() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_34) XGof_ref_a() *_llcppg_anon_35 {
	return (*_llcppg_anon_35)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_34) XGof_ref_b() *_llcppg_anon_35 {
	return (*_llcppg_anon_35)(unsafe.Pointer(p))
}
func (p *UnionNamedStructMembers) XGof_ref_before() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *UnionNamedStructMembers) XGof_ref_nested() *_llcppg_anon_34 {
	return (*_llcppg_anon_34)(unsafe.Pointer(p))
}
func (p *UnionNamedStructMembers) XGof_ref_after() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_37) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_37) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_36) XGof_ref_a() *_llcppg_anon_37 {
	return (*_llcppg_anon_37)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_36) XGof_ref_b() *_llcppg_anon_37 {
	return (*_llcppg_anon_37)(unsafe.Pointer(p))
}
func (p *UnionNamedUnionMembers) XGof_ref_before() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *UnionNamedUnionMembers) XGof_ref_nested() *_llcppg_anon_36 {
	return (*_llcppg_anon_36)(unsafe.Pointer(p))
}
func (p *UnionNamedUnionMembers) XGof_ref_after() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_38) XGof_ref_a() *_llcppg_anon_39 {
	return (*_llcppg_anon_39)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_38) XGof_ref_b() *_llcppg_anon_39 {
	return (*_llcppg_anon_39)(unsafe.Pointer(p))
}
func (p *UnionNamedEnumMembers) XGof_ref_before() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *UnionNamedEnumMembers) XGof_ref_nested() *_llcppg_anon_38 {
	return (*_llcppg_anon_38)(unsafe.Pointer(p))
}
func (p *UnionNamedEnumMembers) XGof_ref_after() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_40) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_40) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *UnionPromotedUnionMembers) XGof_ref_before() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *UnionPromotedUnionMembers) XGof_ref_a() *_llcppg_anon_40 {
	return (*_llcppg_anon_40)(unsafe.Pointer(p))
}
func (p *UnionPromotedUnionMembers) XGof_ref_b() *_llcppg_anon_40 {
	return (*_llcppg_anon_40)(unsafe.Pointer(p))
}
func (p *UnionPromotedUnionMembers) XGof_ref_after() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *UnionPromotedEnumMembers) XGof_ref_before() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *UnionPromotedEnumMembers) XGof_ref_a() *_llcppg_anon_41 {
	return (*_llcppg_anon_41)(unsafe.Pointer(p))
}
func (p *UnionPromotedEnumMembers) XGof_ref_b() *_llcppg_anon_41 {
	return (*_llcppg_anon_41)(unsafe.Pointer(p))
}
func (p *UnionPromotedEnumMembers) XGof_ref_after() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
