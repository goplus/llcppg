package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

const XGoPackage = true

type Supported struct {
	_xgo_union [1]uint32
}

func (p *Supported) XGof_ref_value() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
