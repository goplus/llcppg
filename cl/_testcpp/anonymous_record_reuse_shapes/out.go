package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

const XGoPackage = true

type StructPointerFields struct {
	A *_llcppg_anon_0
	B *_llcppg_anon_0
}
type StructDoublePointerFields struct {
	A **_llcppg_anon_1
	B **_llcppg_anon_1
}
type StructArrayFields struct {
	A [2]_llcppg_anon_2
	B [2]_llcppg_anon_2
}
type StructMatrixFields struct {
	A [2][3]_llcppg_anon_3
	B [2][3]_llcppg_anon_3
}
type StructPointerArrayFields struct {
	A [2]*_llcppg_anon_4
	B [2]*_llcppg_anon_4
}
type StructArrayPointerFields struct {
	A *[2]_llcppg_anon_5
	B *[2]_llcppg_anon_5
}
type StructConstFields struct {
	A _llcppg_anon_6
	B _llcppg_anon_6
}
type UnionDoublePointerFields struct {
	A **_llcppg_anon_7
	B **_llcppg_anon_7
}
type UnionArrayFields struct {
	A [2]_llcppg_anon_8
	B [2]_llcppg_anon_8
}
type UnionMatrixFields struct {
	A [2][3]_llcppg_anon_9
	B [2][3]_llcppg_anon_9
}
type UnionPointerArrayFields struct {
	A [2]*_llcppg_anon_10
	B [2]*_llcppg_anon_10
}
type UnionArrayPointerFields struct {
	A *[2]_llcppg_anon_11
	B *[2]_llcppg_anon_11
}
type UnionConstFields struct {
	A _llcppg_anon_12
	B _llcppg_anon_12
}
type CallbackRecordFields struct {
	A _llcppg_anon_13
	B _llcppg_anon_13
}
type StructPointerMembers struct {
	_xgo_union [1]uint64
}
type StructDoublePointerMembers struct {
	_xgo_union [1]uint64
}
type StructArrayMembers struct {
	_xgo_union [2]uint32
}
type StructMatrixMembers struct {
	_xgo_union [6]uint32
}
type StructPointerArrayMembers struct {
	_xgo_union [2]uint64
}
type StructArrayPointerMembers struct {
	_xgo_union [1]uint64
}
type StructConstMembers struct {
	_xgo_union [1]uint32
}
type UnionDoublePointerMembers struct {
	_xgo_union [1]uint64
}
type UnionArrayMembers struct {
	_xgo_union [2]uint64
}
type UnionMatrixMembers struct {
	_xgo_union [6]uint64
}
type UnionPointerArrayMembers struct {
	_xgo_union [2]uint64
}
type UnionArrayPointerMembers struct {
	_xgo_union [1]uint64
}
type UnionConstMembers struct {
	_xgo_union [1]uint64
}
type CallbackRecordMembers struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_0 struct {
	Value c.Int
}
type _llcppg_anon_1 struct {
	Value c.Int
}
type _llcppg_anon_2 struct {
	Value c.Int
}
type _llcppg_anon_3 struct {
	Value c.Int
}
type _llcppg_anon_4 struct {
	Value c.Int
}
type _llcppg_anon_5 struct {
	Value c.Int
}
type _llcppg_anon_6 struct {
	Value c.Int
}
type _llcppg_anon_7 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_8 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_9 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_10 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_11 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_12 struct {
	_xgo_union [1]uint64
}

//llgo:cdecl
type _llcppg_anon_13 struct {
	Call func(_llcppg_param1 c.Int)
}
type _llcppg_anon_14 struct {
	Value c.Int
}
type _llcppg_anon_15 struct {
	Value c.Int
}
type _llcppg_anon_16 struct {
	Value c.Int
}
type _llcppg_anon_17 struct {
	Value c.Int
}
type _llcppg_anon_18 struct {
	Value c.Int
}
type _llcppg_anon_19 struct {
	Value c.Int
}
type _llcppg_anon_20 struct {
	Value c.Int
}
type _llcppg_anon_21 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_22 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_23 struct {
	_xgo_union [1]uint64
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

//llgo:cdecl
type _llcppg_anon_27 struct {
	Call func(_llcppg_param1 c.Int)
}

func (p *_llcppg_anon_7) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_7) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_8) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_8) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_9) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_9) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_10) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_10) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_11) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_11) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_12) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_12) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *StructPointerMembers) XGof_ref_a() **_llcppg_anon_14 {
	return (**_llcppg_anon_14)(unsafe.Pointer(p))
}
func (p *StructPointerMembers) XGof_ref_b() **_llcppg_anon_14 {
	return (**_llcppg_anon_14)(unsafe.Pointer(p))
}
func (p *StructDoublePointerMembers) XGof_ref_a() ***_llcppg_anon_15 {
	return (***_llcppg_anon_15)(unsafe.Pointer(p))
}
func (p *StructDoublePointerMembers) XGof_ref_b() ***_llcppg_anon_15 {
	return (***_llcppg_anon_15)(unsafe.Pointer(p))
}
func (p *StructArrayMembers) XGof_ref_a() *[2]_llcppg_anon_16 {
	return (*[2]_llcppg_anon_16)(unsafe.Pointer(p))
}
func (p *StructArrayMembers) XGof_ref_b() *[2]_llcppg_anon_16 {
	return (*[2]_llcppg_anon_16)(unsafe.Pointer(p))
}
func (p *StructMatrixMembers) XGof_ref_a() *[2][3]_llcppg_anon_17 {
	return (*[2][3]_llcppg_anon_17)(unsafe.Pointer(p))
}
func (p *StructMatrixMembers) XGof_ref_b() *[2][3]_llcppg_anon_17 {
	return (*[2][3]_llcppg_anon_17)(unsafe.Pointer(p))
}
func (p *StructPointerArrayMembers) XGof_ref_a() *[2]*_llcppg_anon_18 {
	return (*[2]*_llcppg_anon_18)(unsafe.Pointer(p))
}
func (p *StructPointerArrayMembers) XGof_ref_b() *[2]*_llcppg_anon_18 {
	return (*[2]*_llcppg_anon_18)(unsafe.Pointer(p))
}
func (p *StructArrayPointerMembers) XGof_ref_a() **[2]_llcppg_anon_19 {
	return (**[2]_llcppg_anon_19)(unsafe.Pointer(p))
}
func (p *StructArrayPointerMembers) XGof_ref_b() **[2]_llcppg_anon_19 {
	return (**[2]_llcppg_anon_19)(unsafe.Pointer(p))
}
func (p *StructConstMembers) XGof_ref_a() *_llcppg_anon_20 {
	return (*_llcppg_anon_20)(unsafe.Pointer(p))
}
func (p *StructConstMembers) XGof_ref_b() *_llcppg_anon_20 {
	return (*_llcppg_anon_20)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_21) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_21) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *UnionDoublePointerMembers) XGof_ref_a() ***_llcppg_anon_21 {
	return (***_llcppg_anon_21)(unsafe.Pointer(p))
}
func (p *UnionDoublePointerMembers) XGof_ref_b() ***_llcppg_anon_21 {
	return (***_llcppg_anon_21)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_22) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_22) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *UnionArrayMembers) XGof_ref_a() *[2]_llcppg_anon_22 {
	return (*[2]_llcppg_anon_22)(unsafe.Pointer(p))
}
func (p *UnionArrayMembers) XGof_ref_b() *[2]_llcppg_anon_22 {
	return (*[2]_llcppg_anon_22)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_23) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_23) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *UnionMatrixMembers) XGof_ref_a() *[2][3]_llcppg_anon_23 {
	return (*[2][3]_llcppg_anon_23)(unsafe.Pointer(p))
}
func (p *UnionMatrixMembers) XGof_ref_b() *[2][3]_llcppg_anon_23 {
	return (*[2][3]_llcppg_anon_23)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_24) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_24) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *UnionPointerArrayMembers) XGof_ref_a() *[2]*_llcppg_anon_24 {
	return (*[2]*_llcppg_anon_24)(unsafe.Pointer(p))
}
func (p *UnionPointerArrayMembers) XGof_ref_b() *[2]*_llcppg_anon_24 {
	return (*[2]*_llcppg_anon_24)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_25) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_25) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *UnionArrayPointerMembers) XGof_ref_a() **[2]_llcppg_anon_25 {
	return (**[2]_llcppg_anon_25)(unsafe.Pointer(p))
}
func (p *UnionArrayPointerMembers) XGof_ref_b() **[2]_llcppg_anon_25 {
	return (**[2]_llcppg_anon_25)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_26) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_26) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *UnionConstMembers) XGof_ref_a() *_llcppg_anon_26 {
	return (*_llcppg_anon_26)(unsafe.Pointer(p))
}
func (p *UnionConstMembers) XGof_ref_b() *_llcppg_anon_26 {
	return (*_llcppg_anon_26)(unsafe.Pointer(p))
}
func (p *CallbackRecordMembers) XGof_ref_a() *_llcppg_anon_27 {
	return (*_llcppg_anon_27)(unsafe.Pointer(p))
}
func (p *CallbackRecordMembers) XGof_ref_b() *_llcppg_anon_27 {
	return (*_llcppg_anon_27)(unsafe.Pointer(p))
}
