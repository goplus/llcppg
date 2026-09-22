package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

type Ios struct {
	State c.Int
}
type In struct {
	_xgo_vptr unsafe.Pointer
	Igcount   c.Int
	Ios
}
type Out struct {
	_xgo_vptr unsafe.Pointer
	Oputcount c.Int
	Ios
}
type IosBase struct {
	_xgo_vptr unsafe.Pointer
	Flags     c.Int
}
type IStream struct {
	_xgo_vptr unsafe.Pointer
	Gpos      c.Int
	IosBase
}
type OStream struct {
	_xgo_vptr unsafe.Pointer
	Ppos      c.Int
	IosBase
}
type IOStream struct {
	_xgo_vptr         unsafe.Pointer
	Gpos              c.Int
	_xgo_vptr_OStream unsafe.Pointer
	Ppos              c.Int
	IosBase
}

// llgo:type C
type _xgo_vtable_IosBase struct {
	XGo_dtor          func(this *IosBase)
	XGo_dtor_deleting func(this *IosBase)
}

func (p *IosBase) XGo_vptr() *_xgo_vtable_IosBase {
	return (*_xgo_vtable_IosBase)(p._xgo_vptr)
}

// llgo:link (*IosBase).XGo_Dtor C._ZN7IosBaseD1Ev
func (this *IosBase) XGo_Dtor() {
}
