package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

const LLGoPackage = "link: -L/path/foo -lfoo"

type Color c.Int

const (
	Red   Color = 0
	Green Color = 1
	Blue  Color = 2
)

type Matrix struct {
	Rows c.Int
	Cols c.Int
	Data [16]c.Double
	Name [32]c.Char
}

//go:linkname Primitives C.primitives
func Primitives(ch c.Char, uc uint8, sc int8, s int16, us uint16, i c.Int, ui c.Uint, l c.Long, ul c.Ulong, ll c.LongLong, ull c.UlongLong, flt c.Float) c.Double

// llgo:link Color.Pick C.pick
func (c Color) Pick() Color {
	return c
}

//go:linkname Sum C.sum
func Sum(values *c.Int, count c.Int) c.Int

//go:linkname Fill C.fill
func Fill(buf *c.Double, n c.Int)

//go:linkname Transform C.transform
func Transform(matrix *[4]c.Int, p *[10]c.Int)

//go:linkname Sort C.sort
func Sort(a unsafe.Pointer, b unsafe.Pointer, elementSize c.Int, count c.Int, cmp func(_llcppg_param1 unsafe.Pointer, _llcppg_param2 unsafe.Pointer) c.Int) c.Int

//go:linkname F C.f
func F(callback func() c.Int)

//go:linkname G C.g
func G()

//go:linkname H C.h
func H(callbackPtr *func())
