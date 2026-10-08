package foo

import (
	"github.com/goplus/lib/c"
	"github.com/goplus/lib/c/bitfield"
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

// unsigned int a : 1
func (p *Flags) XGof_get_a() c.Uint {
	return c.Uint(bitfield.Unsigned(unsafe.Pointer(p), 0, 1))
}

// unsigned int a : 1
func (p *Flags) XGof_set_a(v c.Uint) {
	bitfield.Set(unsafe.Pointer(p), 0, 1, uint64(v))
}

// unsigned int b : 3
func (p *Flags) XGof_get_b() c.Uint {
	return c.Uint(bitfield.Unsigned(unsafe.Pointer(p), 1, 3))
}

// unsigned int b : 3
func (p *Flags) XGof_set_b(v c.Uint) {
	bitfield.Set(unsafe.Pointer(p), 1, 3, uint64(v))
}

// unsigned int c : 12
func (p *Flags) XGof_get_c() c.Uint {
	return c.Uint(bitfield.Unsigned(unsafe.Pointer(p), 4, 12))
}

// unsigned int c : 12
func (p *Flags) XGof_set_c(v c.Uint) {
	bitfield.Set(unsafe.Pointer(p), 4, 12, uint64(v))
}

// int sx : 5
func (p *Mixed) XGof_get_sx() c.Int {
	return c.Int(bitfield.Signed(unsafe.Pointer(p), 0, 5))
}

// int sx : 5
func (p *Mixed) XGof_set_sx(v c.Int) {
	bitfield.Set(unsafe.Pointer(p), 0, 5, uint64(v))
}

// _Bool flag : 1
func (p *Mixed) XGof_get_flag() bool {
	return bitfield.Unsigned(unsafe.Pointer(p), 32, 1) != 0
}

// _Bool flag : 1
func (p *Mixed) XGof_set_flag(v bool) {
	bitfield.Set(unsafe.Pointer(p), 32, 1, uint64(func() uint64 {
		if v {
			return 1
		}
		return 0
	}()))
}

// unsigned int rest : 20
func (p *Mixed) XGof_get_rest() c.Uint {
	return c.Uint(bitfield.Unsigned(unsafe.Pointer(p), 33, 20))
}

// unsigned int rest : 20
func (p *Mixed) XGof_set_rest(v c.Uint) {
	bitfield.Set(unsafe.Pointer(p), 33, 20, uint64(v))
}
