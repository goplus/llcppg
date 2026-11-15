package foo

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type BuiltinTypedefFields struct {
	Args   *c.VaList
	Size   c.SizeT
	Mapped ExternalMapped
}
type ExternalMapped = uint32

//go:linkname _TakeVaList C.take_va_list
func _TakeVaList(args *c.VaList)
