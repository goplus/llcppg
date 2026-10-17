package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

const XGoPackage = true

type Direct struct {
	_xgo_union [1]uint64
}
type Recursive struct {
	_xgo_union [2]uint64
}
type Mixed struct {
	_xgo_union [1]uint64
}
type NamedBoundary struct {
	_xgo_union [1]uint64
}
type AllFloats struct {
	_xgo_union [2]float32
}
type OnlyFloats struct {
	_xgo_union [2]float64
}

// The enclosing struct's union is not located at the struct's base address.
type Envelope struct {
	Kind c.Int
	_llcppg_anon_11
}
type AnonymousUnion struct {
	_xgo_union [1]uint64
}
type Alternatives struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_0 struct {
	X c.Int
	Y c.Int
}
type _llcppg_anon_1 struct {
	Tag uint8
	_llcppg_anon_2
}
type _llcppg_anon_2 struct {
	X int16
	_llcppg_anon_3
}
type _llcppg_anon_3 struct {
	_xgo_union [1]uint32
}
type _llcppg_anon_4 struct {
	Kind c.Int
	_llcppg_anon_5
}
type _llcppg_anon_5 struct {
	_xgo_union [1]uint32
}
type _llcppg_anon_6 struct {
	X int16
	Y int16
}
type _llcppg_anon_7 struct {
	Lead  c.Int
	Named _llcppg_anon_8
}
type _llcppg_anon_8 struct {
	Value c.Int
}
type _llcppg_anon_9 struct {
	X c.Float
	Y c.Float
}
type _llcppg_anon_10 struct {
	X c.Double
	Y c.Double
}
type _llcppg_anon_11 struct {
	_xgo_union [1]uint64
}
type _llcppg_anon_12 struct {
	X c.Int
	Y int16
}
type _llcppg_anon_13 struct {
	X c.Int
	Y c.Int
}
type _llcppg_anon_14 struct {
	X c.Int
	Y c.Int
}
type _llcppg_anon_15 struct {
	Low  int16
	High int16
	Tail c.Int
}

func (p *Direct) XGof_ref_raw() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *Direct) XGof_ref_x() *c.Int {
	return &(*_llcppg_anon_0)(unsafe.Pointer(p)).X
}
func (p *Direct) XGof_ref_y() *c.Int {
	return &(*_llcppg_anon_0)(unsafe.Pointer(p)).Y
}
func (p *_llcppg_anon_3) XGof_ref_number() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_3) XGof_ref_real() *c.Float {
	return (*c.Float)(unsafe.Pointer(p))
}
func (p *Recursive) XGof_ref_raw() *[2]c.Double {
	return (*[2]c.Double)(unsafe.Pointer(p))
}
func (p *Recursive) XGof_ref_tag() *uint8 {
	return &(*_llcppg_anon_1)(unsafe.Pointer(p)).Tag
}
func (p *Recursive) XGof_ref_x() *int16 {
	return &(*_llcppg_anon_1)(unsafe.Pointer(p)).X
}
func (p *Recursive) XGof_ref_number() *c.Int {
	return (*_llcppg_anon_1)(unsafe.Pointer(p)).XGof_ref_number()
}
func (p *Recursive) XGof_ref_real() *c.Float {
	return (*_llcppg_anon_1)(unsafe.Pointer(p)).XGof_ref_real()
}
func (p *_llcppg_anon_5) XGof_ref_code() *c.Uint {
	return (*c.Uint)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_5) XGof_ref_x() *int16 {
	return &(*_llcppg_anon_6)(unsafe.Pointer(p)).X
}
func (p *_llcppg_anon_5) XGof_ref_y() *int16 {
	return &(*_llcppg_anon_6)(unsafe.Pointer(p)).Y
}
func (p *Mixed) XGof_ref_raw() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *Mixed) XGof_ref_kind() *c.Int {
	return &(*_llcppg_anon_4)(unsafe.Pointer(p)).Kind
}
func (p *Mixed) XGof_ref_x() *int16 {
	return (*_llcppg_anon_4)(unsafe.Pointer(p)).XGof_ref_x()
}
func (p *Mixed) XGof_ref_y() *int16 {
	return (*_llcppg_anon_4)(unsafe.Pointer(p)).XGof_ref_y()
}
func (p *Mixed) XGof_ref_code() *c.Uint {
	return (*_llcppg_anon_4)(unsafe.Pointer(p)).XGof_ref_code()
}
func (p *NamedBoundary) XGof_ref_raw() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *NamedBoundary) XGof_ref_lead() *c.Int {
	return &(*_llcppg_anon_7)(unsafe.Pointer(p)).Lead
}
func (p *NamedBoundary) XGof_ref_named() *_llcppg_anon_8 {
	return &(*_llcppg_anon_7)(unsafe.Pointer(p)).Named
}
func (p *AllFloats) XGof_ref_raw() *[2]c.Float {
	return (*[2]c.Float)(unsafe.Pointer(p))
}
func (p *AllFloats) XGof_ref_x() *c.Float {
	return &(*_llcppg_anon_9)(unsafe.Pointer(p)).X
}
func (p *AllFloats) XGof_ref_y() *c.Float {
	return &(*_llcppg_anon_9)(unsafe.Pointer(p)).Y
}
func (p *OnlyFloats) XGof_ref_x() *c.Double {
	return &(*_llcppg_anon_10)(unsafe.Pointer(p)).X
}
func (p *OnlyFloats) XGof_ref_y() *c.Double {
	return &(*_llcppg_anon_10)(unsafe.Pointer(p)).Y
}
func (p *_llcppg_anon_11) XGof_ref_raw() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_11) XGof_ref_x() *c.Int {
	return &(*_llcppg_anon_12)(unsafe.Pointer(p)).X
}
func (p *_llcppg_anon_11) XGof_ref_y() *int16 {
	return &(*_llcppg_anon_12)(unsafe.Pointer(p)).Y
}
func (p *AnonymousUnion) XGof_ref_whole() *c.LongLong {
	return (*c.LongLong)(unsafe.Pointer(p))
}
func (p *AnonymousUnion) XGof_ref_raw() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *AnonymousUnion) XGof_ref_x() *c.Int {
	return &(*_llcppg_anon_13)(unsafe.Pointer(p)).X
}
func (p *AnonymousUnion) XGof_ref_y() *c.Int {
	return &(*_llcppg_anon_13)(unsafe.Pointer(p)).Y
}
func (p *Alternatives) XGof_ref_raw() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *Alternatives) XGof_ref_x() *c.Int {
	return &(*_llcppg_anon_14)(unsafe.Pointer(p)).X
}
func (p *Alternatives) XGof_ref_y() *c.Int {
	return &(*_llcppg_anon_14)(unsafe.Pointer(p)).Y
}
func (p *Alternatives) XGof_ref_low() *int16 {
	return &(*_llcppg_anon_15)(unsafe.Pointer(p)).Low
}
func (p *Alternatives) XGof_ref_high() *int16 {
	return &(*_llcppg_anon_15)(unsafe.Pointer(p)).High
}
func (p *Alternatives) XGof_ref_tail() *c.Int {
	return &(*_llcppg_anon_15)(unsafe.Pointer(p)).Tail
}
