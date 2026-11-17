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
type StructPromotedStructFields struct {
	Before c.Int
	_llcppg_anon_4
	After c.Int
}
type StructPromotedUnionFields struct {
	Before c.Int
	_llcppg_anon_6
	After c.Int
}
type UnionNamedStructFields struct {
	Before c.Int
	Nested _llcppg_anon_8
	After  c.Int
}
type UnionNamedUnionFields struct {
	Before c.Int
	Nested _llcppg_anon_10
	After  c.Int
}
type UnionPromotedStructFields struct {
	Before c.Int
	_llcppg_anon_12
	After c.Int
}
type UnionPromotedUnionFields struct {
	Before c.Int
	_llcppg_anon_14
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
	_xgo_union [1]uint32
}
type _llcppg_anon_9 struct {
	Value c.Int
}
type _llcppg_anon_10 struct {
	_xgo_union [1]uint64
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
	A _llcppg_anon_17
	B _llcppg_anon_17
}
type _llcppg_anon_17 struct {
	Value c.Int
}
type _llcppg_anon_18 struct {
	A _llcppg_anon_19
	B _llcppg_anon_19
}
type _llcppg_anon_19 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_20 struct {
	A _llcppg_anon_21
	B _llcppg_anon_21
}
type _llcppg_anon_21 struct {
	_xgo_union [1]uint64
}
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
	_xgo_union [1]uint64
}

func (p *_llcppg_anon_3) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_3) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_7) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_7) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_8) XGof_ref_a() *_llcppg_anon_9 {
	return (*_llcppg_anon_9)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_8) XGof_ref_b() *_llcppg_anon_9 {
	return (*_llcppg_anon_9)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_11) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_11) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_10) XGof_ref_a() *_llcppg_anon_11 {
	return (*_llcppg_anon_11)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_10) XGof_ref_b() *_llcppg_anon_11 {
	return (*_llcppg_anon_11)(unsafe.Pointer(p))
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
func (p *StructNamedStructMembers) XGof_ref_before() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *StructNamedStructMembers) XGof_ref_nested() *_llcppg_anon_16 {
	return (*_llcppg_anon_16)(unsafe.Pointer(p))
}
func (p *StructNamedStructMembers) XGof_ref_after() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_19) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_19) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *StructNamedUnionMembers) XGof_ref_before() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *StructNamedUnionMembers) XGof_ref_nested() *_llcppg_anon_18 {
	return (*_llcppg_anon_18)(unsafe.Pointer(p))
}
func (p *StructNamedUnionMembers) XGof_ref_after() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_21) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_21) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *StructPromotedUnionMembers) XGof_ref_before() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *StructPromotedUnionMembers) XGof_ref_a() *_llcppg_anon_21 {
	return &(*_llcppg_anon_20)(unsafe.Pointer(p)).A
}
func (p *StructPromotedUnionMembers) XGof_ref_b() *_llcppg_anon_21 {
	return &(*_llcppg_anon_20)(unsafe.Pointer(p)).B
}
func (p *StructPromotedUnionMembers) XGof_ref_after() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_22) XGof_ref_a() *_llcppg_anon_23 {
	return (*_llcppg_anon_23)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_22) XGof_ref_b() *_llcppg_anon_23 {
	return (*_llcppg_anon_23)(unsafe.Pointer(p))
}
func (p *UnionNamedStructMembers) XGof_ref_before() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *UnionNamedStructMembers) XGof_ref_nested() *_llcppg_anon_22 {
	return (*_llcppg_anon_22)(unsafe.Pointer(p))
}
func (p *UnionNamedStructMembers) XGof_ref_after() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
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
func (p *UnionNamedUnionMembers) XGof_ref_before() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *UnionNamedUnionMembers) XGof_ref_nested() *_llcppg_anon_24 {
	return (*_llcppg_anon_24)(unsafe.Pointer(p))
}
func (p *UnionNamedUnionMembers) XGof_ref_after() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_26) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_26) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *UnionPromotedUnionMembers) XGof_ref_before() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *UnionPromotedUnionMembers) XGof_ref_a() *_llcppg_anon_26 {
	return (*_llcppg_anon_26)(unsafe.Pointer(p))
}
func (p *UnionPromotedUnionMembers) XGof_ref_b() *_llcppg_anon_26 {
	return (*_llcppg_anon_26)(unsafe.Pointer(p))
}
func (p *UnionPromotedUnionMembers) XGof_ref_after() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
