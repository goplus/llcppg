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
type StructPointerFields struct {
	A *_llcppg_anon_10
	B *_llcppg_anon_10
}
type StructDoublePointerFields struct {
	A **_llcppg_anon_11
	B **_llcppg_anon_11
}
type StructArrayFields struct {
	A [2]_llcppg_anon_12
	B [2]_llcppg_anon_12
}
type StructMatrixFields struct {
	A [2][3]_llcppg_anon_13
	B [2][3]_llcppg_anon_13
}
type StructPointerArrayFields struct {
	A [2]*_llcppg_anon_14
	B [2]*_llcppg_anon_14
}
type StructArrayPointerFields struct {
	A *[2]_llcppg_anon_15
	B *[2]_llcppg_anon_15
}
type StructConstFields struct {
	A _llcppg_anon_16
	B _llcppg_anon_16
}
type UnionDoublePointerFields struct {
	A **_llcppg_anon_17
	B **_llcppg_anon_17
}
type UnionArrayFields struct {
	A [2]_llcppg_anon_18
	B [2]_llcppg_anon_18
}
type UnionMatrixFields struct {
	A [2][3]_llcppg_anon_19
	B [2][3]_llcppg_anon_19
}
type UnionPointerArrayFields struct {
	A [2]*_llcppg_anon_20
	B [2]*_llcppg_anon_20
}
type UnionArrayPointerFields struct {
	A *[2]_llcppg_anon_21
	B *[2]_llcppg_anon_21
}
type UnionConstFields struct {
	A _llcppg_anon_22
	B _llcppg_anon_22
}
type CallbackRecordFields struct {
	A _llcppg_anon_23
	B _llcppg_anon_23
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
type StructValuePointerFields struct {
	A _llcppg_anon_38
	B *_llcppg_anon_38
}
type StructPointerValueFields struct {
	A *_llcppg_anon_39
	B _llcppg_anon_39
}
type StructPointerDepthFields struct {
	A *_llcppg_anon_40
	B **_llcppg_anon_40
}
type StructValueArrayFields struct {
	A _llcppg_anon_41
	B [2]_llcppg_anon_41
}
type StructArrayExtentFields struct {
	A [2]_llcppg_anon_42
	B [3]_llcppg_anon_42
}
type StructMixedPointerArrayFields struct {
	A [2]*_llcppg_anon_43
	B *[2]_llcppg_anon_43
}
type StructPointerConstFields struct {
	A *_llcppg_anon_44
	B *_llcppg_anon_44
}
type StructAlternatingFields struct {
	A _llcppg_anon_45
	B *_llcppg_anon_45
	C _llcppg_anon_45
}
type UnionValuePointerFields struct {
	A _llcppg_anon_46
	B *_llcppg_anon_46
}
type UnionPointerValueFields struct {
	A *_llcppg_anon_47
	B _llcppg_anon_47
}
type UnionPointerDepthFields struct {
	A *_llcppg_anon_48
	B **_llcppg_anon_48
}
type UnionValueArrayFields struct {
	A _llcppg_anon_49
	B [2]_llcppg_anon_49
}
type UnionArrayExtentFields struct {
	A [2]_llcppg_anon_50
	B [3]_llcppg_anon_50
}
type UnionMixedPointerArrayFields struct {
	A [2]*_llcppg_anon_51
	B *[2]_llcppg_anon_51
}
type UnionPointerConstFields struct {
	A *_llcppg_anon_52
	B *_llcppg_anon_52
}
type UnionAlternatingFields struct {
	A _llcppg_anon_53
	B *_llcppg_anon_53
	C _llcppg_anon_53
}
type StructValuePointerMembers struct {
	_xgo_union [1]uint64
}
type StructPointerValueMembers struct {
	_xgo_union [1]uint64
}
type StructPointerDepthMembers struct {
	_xgo_union [1]uint64
}
type StructValueArrayMembers struct {
	_xgo_union [2]uint32
}
type StructArrayExtentMembers struct {
	_xgo_union [3]uint32
}
type StructMixedPointerArrayMembers struct {
	_xgo_union [2]uint64
}
type StructPointerConstMembers struct {
	_xgo_union [1]uint64
}
type StructAlternatingMembers struct {
	_xgo_union [1]uint64
}
type UnionValuePointerMembers struct {
	_xgo_union [1]uint64
}
type UnionPointerValueMembers struct {
	_xgo_union [1]uint64
}
type UnionPointerDepthMembers struct {
	_xgo_union [1]uint64
}
type UnionValueArrayMembers struct {
	_xgo_union [2]uint64
}
type UnionArrayExtentMembers struct {
	_xgo_union [3]uint64
}
type UnionMixedPointerArrayMembers struct {
	_xgo_union [2]uint64
}
type UnionPointerConstMembers struct {
	_xgo_union [1]uint64
}
type UnionAlternatingMembers struct {
	_xgo_union [1]uint64
}
type ReferenceFields struct {
	A *_llcppg_anon_70
	B *_llcppg_anon_70
}
type RvalueReferenceFields struct {
	A *_llcppg_anon_71
	B *_llcppg_anon_71
}
type MixedReferenceFields struct {
	A *_llcppg_anon_72
	B *_llcppg_anon_72
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
	Value c.Int
}
type _llcppg_anon_11 struct {
	Value c.Int
}
type _llcppg_anon_12 struct {
	Value c.Int
}
type _llcppg_anon_13 struct {
	Value c.Int
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
	_xgo_union [1]uint64
}
type _llcppg_anon_18 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_19 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_20 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_21 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_22 struct {
	_xgo_union [1]uint64
}

//llgo:cdecl
type _llcppg_anon_23 struct {
	Call func(_llcppg_param1 c.Int)
}
type _llcppg_anon_24 struct {
	Value c.Int
}
type _llcppg_anon_25 struct {
	Value c.Int
}
type _llcppg_anon_26 struct {
	Value c.Int
}
type _llcppg_anon_27 struct {
	Value c.Int
}
type _llcppg_anon_28 struct {
	Value c.Int
}
type _llcppg_anon_29 struct {
	Value c.Int
}
type _llcppg_anon_30 struct {
	Value c.Int
}
type _llcppg_anon_31 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_32 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_33 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_34 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_35 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_36 struct {
	_xgo_union [1]uint64
}

//llgo:cdecl
type _llcppg_anon_37 struct {
	Call func(_llcppg_param1 c.Int)
}
type _llcppg_anon_38 struct {
	Value c.Int
}
type _llcppg_anon_39 struct {
	Value c.Int
}
type _llcppg_anon_40 struct {
	Value c.Int
}
type _llcppg_anon_41 struct {
	Value c.Int
}
type _llcppg_anon_42 struct {
	Value c.Int
}
type _llcppg_anon_43 struct {
	Value c.Int
}
type _llcppg_anon_44 struct {
	Value c.Int
}
type _llcppg_anon_45 struct {
	Value c.Int
}
type _llcppg_anon_46 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_47 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_48 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_49 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_50 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_51 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_52 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_53 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_54 struct {
	Value c.Int
}
type _llcppg_anon_55 struct {
	Value c.Int
}
type _llcppg_anon_56 struct {
	Value c.Int
}
type _llcppg_anon_57 struct {
	Value c.Int
}
type _llcppg_anon_58 struct {
	Value c.Int
}
type _llcppg_anon_59 struct {
	Value c.Int
}
type _llcppg_anon_60 struct {
	Value c.Int
}
type _llcppg_anon_61 struct {
	Value c.Int
}
type _llcppg_anon_62 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_63 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_64 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_65 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_66 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_67 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_68 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_69 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_70 struct {
	Value c.Int
}
type _llcppg_anon_71 struct {
	Value c.Int
}
type _llcppg_anon_72 struct {
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
func (p *_llcppg_anon_17) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_17) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_18) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_18) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_19) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_19) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_20) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_20) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_21) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_21) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_22) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_22) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *StructPointerMembers) XGof_ref_a() **_llcppg_anon_24 {
	return (**_llcppg_anon_24)(unsafe.Pointer(p))
}
func (p *StructPointerMembers) XGof_ref_b() **_llcppg_anon_24 {
	return (**_llcppg_anon_24)(unsafe.Pointer(p))
}
func (p *StructDoublePointerMembers) XGof_ref_a() ***_llcppg_anon_25 {
	return (***_llcppg_anon_25)(unsafe.Pointer(p))
}
func (p *StructDoublePointerMembers) XGof_ref_b() ***_llcppg_anon_25 {
	return (***_llcppg_anon_25)(unsafe.Pointer(p))
}
func (p *StructArrayMembers) XGof_ref_a() *[2]_llcppg_anon_26 {
	return (*[2]_llcppg_anon_26)(unsafe.Pointer(p))
}
func (p *StructArrayMembers) XGof_ref_b() *[2]_llcppg_anon_26 {
	return (*[2]_llcppg_anon_26)(unsafe.Pointer(p))
}
func (p *StructMatrixMembers) XGof_ref_a() *[2][3]_llcppg_anon_27 {
	return (*[2][3]_llcppg_anon_27)(unsafe.Pointer(p))
}
func (p *StructMatrixMembers) XGof_ref_b() *[2][3]_llcppg_anon_27 {
	return (*[2][3]_llcppg_anon_27)(unsafe.Pointer(p))
}
func (p *StructPointerArrayMembers) XGof_ref_a() *[2]*_llcppg_anon_28 {
	return (*[2]*_llcppg_anon_28)(unsafe.Pointer(p))
}
func (p *StructPointerArrayMembers) XGof_ref_b() *[2]*_llcppg_anon_28 {
	return (*[2]*_llcppg_anon_28)(unsafe.Pointer(p))
}
func (p *StructArrayPointerMembers) XGof_ref_a() **[2]_llcppg_anon_29 {
	return (**[2]_llcppg_anon_29)(unsafe.Pointer(p))
}
func (p *StructArrayPointerMembers) XGof_ref_b() **[2]_llcppg_anon_29 {
	return (**[2]_llcppg_anon_29)(unsafe.Pointer(p))
}
func (p *StructConstMembers) XGof_ref_a() *_llcppg_anon_30 {
	return (*_llcppg_anon_30)(unsafe.Pointer(p))
}
func (p *StructConstMembers) XGof_ref_b() *_llcppg_anon_30 {
	return (*_llcppg_anon_30)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_31) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_31) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *UnionDoublePointerMembers) XGof_ref_a() ***_llcppg_anon_31 {
	return (***_llcppg_anon_31)(unsafe.Pointer(p))
}
func (p *UnionDoublePointerMembers) XGof_ref_b() ***_llcppg_anon_31 {
	return (***_llcppg_anon_31)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_32) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_32) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *UnionArrayMembers) XGof_ref_a() *[2]_llcppg_anon_32 {
	return (*[2]_llcppg_anon_32)(unsafe.Pointer(p))
}
func (p *UnionArrayMembers) XGof_ref_b() *[2]_llcppg_anon_32 {
	return (*[2]_llcppg_anon_32)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_33) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_33) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *UnionMatrixMembers) XGof_ref_a() *[2][3]_llcppg_anon_33 {
	return (*[2][3]_llcppg_anon_33)(unsafe.Pointer(p))
}
func (p *UnionMatrixMembers) XGof_ref_b() *[2][3]_llcppg_anon_33 {
	return (*[2][3]_llcppg_anon_33)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_34) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_34) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *UnionPointerArrayMembers) XGof_ref_a() *[2]*_llcppg_anon_34 {
	return (*[2]*_llcppg_anon_34)(unsafe.Pointer(p))
}
func (p *UnionPointerArrayMembers) XGof_ref_b() *[2]*_llcppg_anon_34 {
	return (*[2]*_llcppg_anon_34)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_35) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_35) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *UnionArrayPointerMembers) XGof_ref_a() **[2]_llcppg_anon_35 {
	return (**[2]_llcppg_anon_35)(unsafe.Pointer(p))
}
func (p *UnionArrayPointerMembers) XGof_ref_b() **[2]_llcppg_anon_35 {
	return (**[2]_llcppg_anon_35)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_36) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_36) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *UnionConstMembers) XGof_ref_a() *_llcppg_anon_36 {
	return (*_llcppg_anon_36)(unsafe.Pointer(p))
}
func (p *UnionConstMembers) XGof_ref_b() *_llcppg_anon_36 {
	return (*_llcppg_anon_36)(unsafe.Pointer(p))
}
func (p *CallbackRecordMembers) XGof_ref_a() *_llcppg_anon_37 {
	return (*_llcppg_anon_37)(unsafe.Pointer(p))
}
func (p *CallbackRecordMembers) XGof_ref_b() *_llcppg_anon_37 {
	return (*_llcppg_anon_37)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_46) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_46) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_47) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_47) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_48) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_48) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_49) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_49) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_50) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_50) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_51) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_51) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_52) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_52) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_53) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_53) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *StructValuePointerMembers) XGof_ref_a() *_llcppg_anon_54 {
	return (*_llcppg_anon_54)(unsafe.Pointer(p))
}
func (p *StructValuePointerMembers) XGof_ref_b() **_llcppg_anon_54 {
	return (**_llcppg_anon_54)(unsafe.Pointer(p))
}
func (p *StructPointerValueMembers) XGof_ref_a() **_llcppg_anon_55 {
	return (**_llcppg_anon_55)(unsafe.Pointer(p))
}
func (p *StructPointerValueMembers) XGof_ref_b() *_llcppg_anon_55 {
	return (*_llcppg_anon_55)(unsafe.Pointer(p))
}
func (p *StructPointerDepthMembers) XGof_ref_a() **_llcppg_anon_56 {
	return (**_llcppg_anon_56)(unsafe.Pointer(p))
}
func (p *StructPointerDepthMembers) XGof_ref_b() ***_llcppg_anon_56 {
	return (***_llcppg_anon_56)(unsafe.Pointer(p))
}
func (p *StructValueArrayMembers) XGof_ref_a() *_llcppg_anon_57 {
	return (*_llcppg_anon_57)(unsafe.Pointer(p))
}
func (p *StructValueArrayMembers) XGof_ref_b() *[2]_llcppg_anon_57 {
	return (*[2]_llcppg_anon_57)(unsafe.Pointer(p))
}
func (p *StructArrayExtentMembers) XGof_ref_a() *[2]_llcppg_anon_58 {
	return (*[2]_llcppg_anon_58)(unsafe.Pointer(p))
}
func (p *StructArrayExtentMembers) XGof_ref_b() *[3]_llcppg_anon_58 {
	return (*[3]_llcppg_anon_58)(unsafe.Pointer(p))
}
func (p *StructMixedPointerArrayMembers) XGof_ref_a() *[2]*_llcppg_anon_59 {
	return (*[2]*_llcppg_anon_59)(unsafe.Pointer(p))
}
func (p *StructMixedPointerArrayMembers) XGof_ref_b() **[2]_llcppg_anon_59 {
	return (**[2]_llcppg_anon_59)(unsafe.Pointer(p))
}
func (p *StructPointerConstMembers) XGof_ref_a() **_llcppg_anon_60 {
	return (**_llcppg_anon_60)(unsafe.Pointer(p))
}
func (p *StructPointerConstMembers) XGof_ref_b() **_llcppg_anon_60 {
	return (**_llcppg_anon_60)(unsafe.Pointer(p))
}
func (p *StructAlternatingMembers) XGof_ref_a() *_llcppg_anon_61 {
	return (*_llcppg_anon_61)(unsafe.Pointer(p))
}
func (p *StructAlternatingMembers) XGof_ref_b() **_llcppg_anon_61 {
	return (**_llcppg_anon_61)(unsafe.Pointer(p))
}
func (p *StructAlternatingMembers) XGof_ref_c() *_llcppg_anon_61 {
	return (*_llcppg_anon_61)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_62) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_62) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *UnionValuePointerMembers) XGof_ref_a() *_llcppg_anon_62 {
	return (*_llcppg_anon_62)(unsafe.Pointer(p))
}
func (p *UnionValuePointerMembers) XGof_ref_b() **_llcppg_anon_62 {
	return (**_llcppg_anon_62)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_63) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_63) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *UnionPointerValueMembers) XGof_ref_a() **_llcppg_anon_63 {
	return (**_llcppg_anon_63)(unsafe.Pointer(p))
}
func (p *UnionPointerValueMembers) XGof_ref_b() *_llcppg_anon_63 {
	return (*_llcppg_anon_63)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_64) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_64) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *UnionPointerDepthMembers) XGof_ref_a() **_llcppg_anon_64 {
	return (**_llcppg_anon_64)(unsafe.Pointer(p))
}
func (p *UnionPointerDepthMembers) XGof_ref_b() ***_llcppg_anon_64 {
	return (***_llcppg_anon_64)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_65) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_65) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *UnionValueArrayMembers) XGof_ref_a() *_llcppg_anon_65 {
	return (*_llcppg_anon_65)(unsafe.Pointer(p))
}
func (p *UnionValueArrayMembers) XGof_ref_b() *[2]_llcppg_anon_65 {
	return (*[2]_llcppg_anon_65)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_66) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_66) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *UnionArrayExtentMembers) XGof_ref_a() *[2]_llcppg_anon_66 {
	return (*[2]_llcppg_anon_66)(unsafe.Pointer(p))
}
func (p *UnionArrayExtentMembers) XGof_ref_b() *[3]_llcppg_anon_66 {
	return (*[3]_llcppg_anon_66)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_67) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_67) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *UnionMixedPointerArrayMembers) XGof_ref_a() *[2]*_llcppg_anon_67 {
	return (*[2]*_llcppg_anon_67)(unsafe.Pointer(p))
}
func (p *UnionMixedPointerArrayMembers) XGof_ref_b() **[2]_llcppg_anon_67 {
	return (**[2]_llcppg_anon_67)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_68) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_68) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *UnionPointerConstMembers) XGof_ref_a() **_llcppg_anon_68 {
	return (**_llcppg_anon_68)(unsafe.Pointer(p))
}
func (p *UnionPointerConstMembers) XGof_ref_b() **_llcppg_anon_68 {
	return (**_llcppg_anon_68)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_69) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_69) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *UnionAlternatingMembers) XGof_ref_a() *_llcppg_anon_69 {
	return (*_llcppg_anon_69)(unsafe.Pointer(p))
}
func (p *UnionAlternatingMembers) XGof_ref_b() **_llcppg_anon_69 {
	return (**_llcppg_anon_69)(unsafe.Pointer(p))
}
func (p *UnionAlternatingMembers) XGof_ref_c() *_llcppg_anon_69 {
	return (*_llcppg_anon_69)(unsafe.Pointer(p))
}
