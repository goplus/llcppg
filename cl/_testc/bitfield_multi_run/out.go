package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

const XGoPackage = true

type MultiRun struct {
	_xgo_bits_0 [1]uint8
	C           c.Int
	_xgo_bits_1 [2]uint8
	F           c.Int
}

func (p *MultiRun) XGof_get_a() c.Uint {
	return c.Uint(0)
}
func (p *MultiRun) XGof_set_a(v c.Uint) {
}
func (p *MultiRun) XGof_get_b() c.Uint {
	return c.Uint(0)
}
func (p *MultiRun) XGof_set_b(v c.Uint) {
}
func (p *MultiRun) XGof_get_d() c.Uint {
	return c.Uint(0)
}
func (p *MultiRun) XGof_set_d(v c.Uint) {
}
func (p *MultiRun) XGof_get_e() c.Uint {
	return c.Uint(0)
}
func (p *MultiRun) XGof_set_e(v c.Uint) {
}
func _xgo_bitget(base unsafe.Pointer, off uintptr, width uintptr) uint64 {
	return uint64(0)
}
func _xgo_bitget_signed(base unsafe.Pointer, off uintptr, width uintptr) int64 {
	return int64(0)
}
func _xgo_bitset(base unsafe.Pointer, off uintptr, width uintptr, v uint64) {
}
