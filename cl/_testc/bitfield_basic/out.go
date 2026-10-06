package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

const XGoPackage = true

type Packet struct {
	_xgo_bits_0 [5]uint8
	Id          c.Int
}

func (p *Packet) XGof_get_ready() c.Uint {
	return c.Uint(0)
}
func (p *Packet) XGof_set_ready(v c.Uint) {
}
func (p *Packet) XGof_get_mode() c.Uint {
	return c.Uint(0)
}
func (p *Packet) XGof_set_mode(v c.Uint) {
}
func (p *Packet) XGof_get_delta() c.Int {
	return c.Int(0)
}
func (p *Packet) XGof_set_delta(v c.Int) {
}
func (p *Packet) XGof_get_level() c.Uint {
	return c.Uint(0)
}
func (p *Packet) XGof_set_level(v c.Uint) {
}

func _xgo_bitget(base unsafe.Pointer, off, width uintptr) uint64 {
	return uint64(0)
}
func _xgo_bitget_signed(base unsafe.Pointer, off, width uintptr) int64 {
	return int64(0)
}
func _xgo_bitset(base unsafe.Pointer, off, width uintptr, v uint64) {
}
