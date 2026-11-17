package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

const XGoPackage = true

type StructStruct struct {
	A _llcppg_anon_0
	B _llcppg_anon_0
	C _llcppg_anon_0
}
type UnionStruct struct {
	_xgo_union [1]uint32
}
type StructUnion struct {
	A _llcppg_anon_2
	B _llcppg_anon_2
}
type UnionUnion struct {
	_xgo_union [1]uint64
}
type WrappedFields struct {
	A [2]_llcppg_anon_4
	B [2]_llcppg_anon_4
	C *_llcppg_anon_5
	D *_llcppg_anon_5
}
type WrappedMembers struct {
	_xgo_union [1]uint64
}
type DistinctDeclarations struct {
	A _llcppg_anon_8
	B _llcppg_anon_9
}
type PromotedFields struct {
	_xgo_union [2]uint32
}
type AliasFirst struct {
	A c.Int
	B _llcppg_anon_12
	C _llcppg_anon_12
}
type AliasMiddle struct {
	A _llcppg_anon_13
	B c.Int
	C _llcppg_anon_13
}
type PromotedUnionFields struct {
	_xgo_union [1]uint32
}
type NestedPromotedUnionFields struct {
	_xgo_union [1]uint32
}
type _llcppg_anon_0 struct {
	X c.Int
}
type _llcppg_anon_1 struct {
	X c.Int
}
type _llcppg_anon_2 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_3 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_4 struct {
	X c.Int
}
type _llcppg_anon_5 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_6 struct {
	X c.Int
}
type _llcppg_anon_7 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_8 struct {
	X c.Int
}
type _llcppg_anon_9 struct {
	X c.Int
}
type _llcppg_anon_10 struct {
	A _llcppg_anon_11
	B _llcppg_anon_11
}
type _llcppg_anon_11 struct {
	X c.Int
}
type _llcppg_anon_12 struct {
	X c.Int
}
type _llcppg_anon_13 struct {
	X c.Int
}
type _llcppg_anon_14 struct {
	X c.Int
}
type _llcppg_anon_15 struct {
	Value c.Int
}

func (p *UnionStruct) XGof_ref_a() *_llcppg_anon_1 {
	return (*_llcppg_anon_1)(unsafe.Pointer(p))
}
func (p *UnionStruct) XGof_ref_b() *_llcppg_anon_1 {
	return (*_llcppg_anon_1)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_2) XGof_ref_x() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_2) XGof_ref_y() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_3) XGof_ref_x() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_3) XGof_ref_y() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *UnionUnion) XGof_ref_a() *_llcppg_anon_3 {
	return (*_llcppg_anon_3)(unsafe.Pointer(p))
}
func (p *UnionUnion) XGof_ref_b() *_llcppg_anon_3 {
	return (*_llcppg_anon_3)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_5) XGof_ref_x() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_5) XGof_ref_y() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_7) XGof_ref_x() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_7) XGof_ref_y() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *WrappedMembers) XGof_ref_a() *[2]_llcppg_anon_6 {
	return (*[2]_llcppg_anon_6)(unsafe.Pointer(p))
}
func (p *WrappedMembers) XGof_ref_b() *[2]_llcppg_anon_6 {
	return (*[2]_llcppg_anon_6)(unsafe.Pointer(p))
}
func (p *WrappedMembers) XGof_ref_c() **_llcppg_anon_7 {
	return (**_llcppg_anon_7)(unsafe.Pointer(p))
}
func (p *WrappedMembers) XGof_ref_d() **_llcppg_anon_7 {
	return (**_llcppg_anon_7)(unsafe.Pointer(p))
}
func (p *PromotedFields) XGof_ref_a() *_llcppg_anon_11 {
	return &(*_llcppg_anon_10)(unsafe.Pointer(p)).A
}
func (p *PromotedFields) XGof_ref_b() *_llcppg_anon_11 {
	return &(*_llcppg_anon_10)(unsafe.Pointer(p)).B
}
func (p *PromotedUnionFields) XGof_ref_outer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *PromotedUnionFields) XGof_ref_x() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *PromotedUnionFields) XGof_ref_a() *_llcppg_anon_14 {
	return (*_llcppg_anon_14)(unsafe.Pointer(p))
}
func (p *PromotedUnionFields) XGof_ref_b() *_llcppg_anon_14 {
	return (*_llcppg_anon_14)(unsafe.Pointer(p))
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
func (p *NestedPromotedUnionFields) XGof_ref_a() *_llcppg_anon_15 {
	return (*_llcppg_anon_15)(unsafe.Pointer(p))
}
func (p *NestedPromotedUnionFields) XGof_ref_b() *_llcppg_anon_15 {
	return (*_llcppg_anon_15)(unsafe.Pointer(p))
}
func (p *NestedPromotedUnionFields) XGof_ref_afterW() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *NestedPromotedUnionFields) XGof_ref_tail() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
