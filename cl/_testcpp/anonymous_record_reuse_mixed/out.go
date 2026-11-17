package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

const XGoPackage = true

type StructValuePointerFields struct {
	A _llcppg_anon_0
	B *_llcppg_anon_0
}
type StructPointerValueFields struct {
	A *_llcppg_anon_1
	B _llcppg_anon_1
}
type StructPointerDepthFields struct {
	A *_llcppg_anon_2
	B **_llcppg_anon_2
}
type StructValueArrayFields struct {
	A _llcppg_anon_3
	B [2]_llcppg_anon_3
}
type StructArrayExtentFields struct {
	A [2]_llcppg_anon_4
	B [3]_llcppg_anon_4
}
type StructPointerArrayFields struct {
	A [2]*_llcppg_anon_5
	B *[2]_llcppg_anon_5
}
type StructPointerConstFields struct {
	A *_llcppg_anon_6
	B *_llcppg_anon_6
}
type StructAlternatingFields struct {
	A _llcppg_anon_7
	B *_llcppg_anon_7
	C _llcppg_anon_7
}
type UnionValuePointerFields struct {
	A _llcppg_anon_8
	B *_llcppg_anon_8
}
type UnionPointerValueFields struct {
	A *_llcppg_anon_9
	B _llcppg_anon_9
}
type UnionPointerDepthFields struct {
	A *_llcppg_anon_10
	B **_llcppg_anon_10
}
type UnionValueArrayFields struct {
	A _llcppg_anon_11
	B [2]_llcppg_anon_11
}
type UnionArrayExtentFields struct {
	A [2]_llcppg_anon_12
	B [3]_llcppg_anon_12
}
type UnionPointerArrayFields struct {
	A [2]*_llcppg_anon_13
	B *[2]_llcppg_anon_13
}
type UnionPointerConstFields struct {
	A *_llcppg_anon_14
	B *_llcppg_anon_14
}
type UnionAlternatingFields struct {
	A _llcppg_anon_15
	B *_llcppg_anon_15
	C _llcppg_anon_15
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
type StructPointerArrayMembers struct {
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
type UnionPointerArrayMembers struct {
	_xgo_union [2]uint64
}
type UnionPointerConstMembers struct {
	_xgo_union [1]uint64
}
type UnionAlternatingMembers struct {
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
	Value c.Int
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
type _llcppg_anon_13 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_14 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_15 struct {
	_xgo_union [1]uint64
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
	Value c.Int
}
type _llcppg_anon_22 struct {
	Value c.Int
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
type _llcppg_anon_27 struct {
	_xgo_union [1]uint64
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
	_xgo_union [1]uint64
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
func (p *_llcppg_anon_13) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_13) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_14) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_14) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_15) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_15) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *StructValuePointerMembers) XGof_ref_a() *_llcppg_anon_16 {
	return (*_llcppg_anon_16)(unsafe.Pointer(p))
}
func (p *StructValuePointerMembers) XGof_ref_b() **_llcppg_anon_16 {
	return (**_llcppg_anon_16)(unsafe.Pointer(p))
}
func (p *StructPointerValueMembers) XGof_ref_a() **_llcppg_anon_17 {
	return (**_llcppg_anon_17)(unsafe.Pointer(p))
}
func (p *StructPointerValueMembers) XGof_ref_b() *_llcppg_anon_17 {
	return (*_llcppg_anon_17)(unsafe.Pointer(p))
}
func (p *StructPointerDepthMembers) XGof_ref_a() **_llcppg_anon_18 {
	return (**_llcppg_anon_18)(unsafe.Pointer(p))
}
func (p *StructPointerDepthMembers) XGof_ref_b() ***_llcppg_anon_18 {
	return (***_llcppg_anon_18)(unsafe.Pointer(p))
}
func (p *StructValueArrayMembers) XGof_ref_a() *_llcppg_anon_19 {
	return (*_llcppg_anon_19)(unsafe.Pointer(p))
}
func (p *StructValueArrayMembers) XGof_ref_b() *[2]_llcppg_anon_19 {
	return (*[2]_llcppg_anon_19)(unsafe.Pointer(p))
}
func (p *StructArrayExtentMembers) XGof_ref_a() *[2]_llcppg_anon_20 {
	return (*[2]_llcppg_anon_20)(unsafe.Pointer(p))
}
func (p *StructArrayExtentMembers) XGof_ref_b() *[3]_llcppg_anon_20 {
	return (*[3]_llcppg_anon_20)(unsafe.Pointer(p))
}
func (p *StructPointerArrayMembers) XGof_ref_a() *[2]*_llcppg_anon_21 {
	return (*[2]*_llcppg_anon_21)(unsafe.Pointer(p))
}
func (p *StructPointerArrayMembers) XGof_ref_b() **[2]_llcppg_anon_21 {
	return (**[2]_llcppg_anon_21)(unsafe.Pointer(p))
}
func (p *StructPointerConstMembers) XGof_ref_a() **_llcppg_anon_22 {
	return (**_llcppg_anon_22)(unsafe.Pointer(p))
}
func (p *StructPointerConstMembers) XGof_ref_b() **_llcppg_anon_22 {
	return (**_llcppg_anon_22)(unsafe.Pointer(p))
}
func (p *StructAlternatingMembers) XGof_ref_a() *_llcppg_anon_23 {
	return (*_llcppg_anon_23)(unsafe.Pointer(p))
}
func (p *StructAlternatingMembers) XGof_ref_b() **_llcppg_anon_23 {
	return (**_llcppg_anon_23)(unsafe.Pointer(p))
}
func (p *StructAlternatingMembers) XGof_ref_c() *_llcppg_anon_23 {
	return (*_llcppg_anon_23)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_24) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_24) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *UnionValuePointerMembers) XGof_ref_a() *_llcppg_anon_24 {
	return (*_llcppg_anon_24)(unsafe.Pointer(p))
}
func (p *UnionValuePointerMembers) XGof_ref_b() **_llcppg_anon_24 {
	return (**_llcppg_anon_24)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_25) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_25) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *UnionPointerValueMembers) XGof_ref_a() **_llcppg_anon_25 {
	return (**_llcppg_anon_25)(unsafe.Pointer(p))
}
func (p *UnionPointerValueMembers) XGof_ref_b() *_llcppg_anon_25 {
	return (*_llcppg_anon_25)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_26) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_26) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *UnionPointerDepthMembers) XGof_ref_a() **_llcppg_anon_26 {
	return (**_llcppg_anon_26)(unsafe.Pointer(p))
}
func (p *UnionPointerDepthMembers) XGof_ref_b() ***_llcppg_anon_26 {
	return (***_llcppg_anon_26)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_27) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_27) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *UnionValueArrayMembers) XGof_ref_a() *_llcppg_anon_27 {
	return (*_llcppg_anon_27)(unsafe.Pointer(p))
}
func (p *UnionValueArrayMembers) XGof_ref_b() *[2]_llcppg_anon_27 {
	return (*[2]_llcppg_anon_27)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_28) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_28) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *UnionArrayExtentMembers) XGof_ref_a() *[2]_llcppg_anon_28 {
	return (*[2]_llcppg_anon_28)(unsafe.Pointer(p))
}
func (p *UnionArrayExtentMembers) XGof_ref_b() *[3]_llcppg_anon_28 {
	return (*[3]_llcppg_anon_28)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_29) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_29) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *UnionPointerArrayMembers) XGof_ref_a() *[2]*_llcppg_anon_29 {
	return (*[2]*_llcppg_anon_29)(unsafe.Pointer(p))
}
func (p *UnionPointerArrayMembers) XGof_ref_b() **[2]_llcppg_anon_29 {
	return (**[2]_llcppg_anon_29)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_30) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_30) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *UnionPointerConstMembers) XGof_ref_a() **_llcppg_anon_30 {
	return (**_llcppg_anon_30)(unsafe.Pointer(p))
}
func (p *UnionPointerConstMembers) XGof_ref_b() **_llcppg_anon_30 {
	return (**_llcppg_anon_30)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_31) XGof_ref_integer() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_31) XGof_ref_decimal() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *UnionAlternatingMembers) XGof_ref_a() *_llcppg_anon_31 {
	return (*_llcppg_anon_31)(unsafe.Pointer(p))
}
func (p *UnionAlternatingMembers) XGof_ref_b() **_llcppg_anon_31 {
	return (**_llcppg_anon_31)(unsafe.Pointer(p))
}
func (p *UnionAlternatingMembers) XGof_ref_c() *_llcppg_anon_31 {
	return (*_llcppg_anon_31)(unsafe.Pointer(p))
}
