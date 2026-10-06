package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

const XGoPackage = true

type Flags struct {
	_xgo_bits_0 [2]uint8
	Value       c.Int
}
type Mixed struct {
	_xgo_align  [0]uint32
	_xgo_bits_0 [7]uint8
	_           [1]uint8
}

func _xgo_bitget(base unsafe.Pointer, off uintptr, width uintptr) uint64 {
	s := unsafe.Slice((*uint8)(base), (off+width+7)/8)
	var r uint64
	for i := uintptr(0); i < width; i++ {
		if s[(off+i)/8]&(uint8(1)<<((off+i)%8)) != 0 {
			r |= uint64(1) << i
		}
	}
	return r
}
func _xgo_bitget_signed(base unsafe.Pointer, off uintptr, width uintptr) int64 {
	r := _xgo_bitget(base, off, width)
	shift := 64 - width
	return int64(r<<shift) >> shift
}
func _xgo_bitset(base unsafe.Pointer, off uintptr, width uintptr, v uint64) {
	s := unsafe.Slice((*uint8)(base), (off+width+7)/8)
	for i := uintptr(0); i < width; i++ {
		idx := (off + i) / 8
		mask := uint8(1) << ((off + i) % 8)
		if v>>i&1 != 0 {
			s[idx] |= mask
		} else {
			s[idx] &^= mask
		}
	}
}

// unsigned int a : 1
func (p *Flags) XGof_get_a() c.Uint {
	return c.Uint(_xgo_bitget(unsafe.Pointer(p), 0, 1))
}

// unsigned int a : 1
func (p *Flags) XGof_set_a(v c.Uint) {
	_xgo_bitset(unsafe.Pointer(p), 0, 1, uint64(v))
}

// unsigned int b : 3
func (p *Flags) XGof_get_b() c.Uint {
	return c.Uint(_xgo_bitget(unsafe.Pointer(p), 1, 3))
}

// unsigned int b : 3
func (p *Flags) XGof_set_b(v c.Uint) {
	_xgo_bitset(unsafe.Pointer(p), 1, 3, uint64(v))
}

// unsigned int c : 12
func (p *Flags) XGof_get_c() c.Uint {
	return c.Uint(_xgo_bitget(unsafe.Pointer(p), 4, 12))
}

// unsigned int c : 12
func (p *Flags) XGof_set_c(v c.Uint) {
	_xgo_bitset(unsafe.Pointer(p), 4, 12, uint64(v))
}

// int sx : 5
func (p *Mixed) XGof_get_sx() c.Int {
	return c.Int(_xgo_bitget_signed(unsafe.Pointer(p), 0, 5))
}

// int sx : 5
func (p *Mixed) XGof_set_sx(v c.Int) {
	_xgo_bitset(unsafe.Pointer(p), 0, 5, uint64(v))
}

// _Bool flag : 1
func (p *Mixed) XGof_get_flag() bool {
	return _xgo_bitget(unsafe.Pointer(p), 32, 1) != 0
}

// _Bool flag : 1
func (p *Mixed) XGof_set_flag(v bool) {
	_xgo_bitset(unsafe.Pointer(p), 32, 1, uint64(func() uint64 {
		if v {
			return 1
		}
		return 0
	}()))
}

// unsigned int rest : 20
func (p *Mixed) XGof_get_rest() c.Uint {
	return c.Uint(_xgo_bitget(unsafe.Pointer(p), 33, 20))
}

// unsigned int rest : 20
func (p *Mixed) XGof_set_rest(v c.Uint) {
	_xgo_bitset(unsafe.Pointer(p), 33, 20, uint64(v))
}
