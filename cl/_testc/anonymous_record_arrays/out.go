package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

const XGoPackage = true

// Like ma_device_info.nativeDataFormats[64], each element owns its fields.
type Device struct {
	NativeDataFormats [64]_llcppg_anon_0
	Sentinel          c.Int
}

// Preserve both extents: cells[1][2].y uses one shared element type.
type Grid struct {
	Cells [2][3]_llcppg_anon_1
}
type Samples struct {
	_xgo_union [6]uint64
}

// Recursive records retain the path entries[1].inner[0].code.
type Nested struct {
	Entries [2]_llcppg_anon_2
}

// Anonymous union elements use the same recursive conversion as structs.
type UnionArrays struct {
	Values [2]_llcppg_anon_5
}

// Resolve existing names before generating an anonymous element type.
type Tagged struct {
	Code c.Int
}
type Alias struct {
	Code c.Int
}
type NamedArrays struct {
	Tagged   [2]Tagged
	Aliased  [3]Alias
	Matrix   [2][3]c.Int
	Pointers [2]*Tagged
	Rows     *[2]Tagged
}

// Element pointers and array pointers retain pointers[1]->code and (*rows)[1].code.
type PointerArrays struct {
	Pointers [2]*_llcppg_anon_6
	Rows     *[2]_llcppg_anon_7
}

// Both callback paths must retain the C calling-convention marker.
//
//llgo:cdecl
type CallbackElements struct {
	Handlers  [2]_llcppg_anon_8
	Callbacks [2]func(_llcppg_param1 c.Int)
}
type _llcppg_anon_0 struct {
	Format     c.Int
	Channels   c.Uint
	SampleRate c.Uint
	Flags      c.Uint
}
type _llcppg_anon_1 struct {
	X int16
	Y int16
}
type _llcppg_anon_2 struct {
	Value _llcppg_anon_3
	Inner [2]_llcppg_anon_4
}
type _llcppg_anon_3 struct {
	_xgo_union [1]uint32
}
type _llcppg_anon_4 struct {
	Code c.Int
}
type _llcppg_anon_5 struct {
	_xgo_union [1]uint32
}
type _llcppg_anon_6 struct {
	Code c.Int
}
type _llcppg_anon_7 struct {
	Code c.Int
}

//llgo:cdecl
type _llcppg_anon_8 struct {
	Call func(_llcppg_param1 c.Int)
}
type _llcppg_anon_9 struct {
	X c.Int
	Y c.Double
}

func (p *Samples) XGof_ref_entries() *[3]_llcppg_anon_9 {
	return (*[3]_llcppg_anon_9)(unsafe.Pointer(p))
}
func (p *Samples) XGof_ref_words() *[6]c.Double {
	return (*[6]c.Double)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_3) XGof_ref_i() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_3) XGof_ref_f() *c.Float {
	return (*c.Float)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_5) XGof_ref_i() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_5) XGof_ref_f() *c.Float {
	return (*c.Float)(unsafe.Pointer(p))
}
