package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

const XGoPackage = true

type Members struct {
	_xgo_union [4]uint64
}
type StructUnion struct {
	_xgo_union [2]uint64
}

// Like yaml_token_t.data: a struct owns a union of named, tagless structs.
type Token struct {
	Type c.Int
	Data _llcppg_anon_0
}
type Nested struct {
	_xgo_union [4]uint32
}

// Named records and typedefs must continue to use their declared Go types.
type Tagged struct {
	Code  c.Int
	Value c.Double
}
type References struct {
	_xgo_union [2]uint64
}
type _llcppg_anon_0 struct {
	_xgo_union [3]uint64
}
type Alias = Tagged
type _llcppg_anon_1 struct {
	Code c.Int
}

func (p *Members) XGof_ref_small() *_llcppg_anon_1 {
	return (*_llcppg_anon_1)(unsafe.Pointer(p))
}

type _llcppg_anon_2 struct {
	Values [2]c.Double
	Flags  c.Int
}

func (p *Members) XGof_ref_large() *_llcppg_anon_2 {
	return (*_llcppg_anon_2)(unsafe.Pointer(p))
}
func (p *Members) XGof_ref_bytes() *[32]uint8 {
	return (*[32]uint8)(unsafe.Pointer(p))
}

type _llcppg_anon_3 struct {
	Kind uint8
	Data _llcppg_anon_4
}
type _llcppg_anon_4 struct {
	_xgo_union [1]uint64
}

func (p *StructUnion) XGof_ref_record() *_llcppg_anon_3 {
	return (*_llcppg_anon_3)(unsafe.Pointer(p))
}
func (p *StructUnion) XGof_ref_words() *[2]c.Double {
	return (*[2]c.Double)(unsafe.Pointer(p))
}

type _llcppg_anon_5 struct {
	Encoding c.Int
}

func (p *_llcppg_anon_0) XGof_ref_stream_start() *_llcppg_anon_5 {
	return (*_llcppg_anon_5)(unsafe.Pointer(p))
}

type _llcppg_anon_6 struct {
	Value  *uint8
	Length c.Ulong
	Style  c.Int
}

func (p *_llcppg_anon_0) XGof_ref_scalar() *_llcppg_anon_6 {
	return (*_llcppg_anon_6)(unsafe.Pointer(p))
}

type _llcppg_anon_7 struct {
	Kind  c.Int
	Point _llcppg_anon_8
	Value _llcppg_anon_9
}
type _llcppg_anon_8 struct {
	X int16
	Y int16
}
type _llcppg_anon_9 struct {
	_xgo_union [2]uint32
}

func (p *Nested) XGof_ref_record() *_llcppg_anon_7 {
	return (*_llcppg_anon_7)(unsafe.Pointer(p))
}
func (p *Nested) XGof_ref_bytes() *[16]uint8 {
	return (*[16]uint8)(unsafe.Pointer(p))
}
func (p *References) XGof_ref_tagged() *Tagged {
	return (*Tagged)(unsafe.Pointer(p))
}
func (p *References) XGof_ref_alias() *Alias {
	return (*Alias)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_4) XGof_ref_code() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_4) XGof_ref_value() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_9) XGof_ref_code() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}

type _llcppg_anon_10 struct {
	X c.Float
	Y c.Float
}

func (p *_llcppg_anon_9) XGof_ref_pair() *_llcppg_anon_10 {
	return (*_llcppg_anon_10)(unsafe.Pointer(p))
}
