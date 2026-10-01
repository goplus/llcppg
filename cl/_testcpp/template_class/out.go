package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

type Foo[T any] struct {
	Value T
}
type Bar[T1 any, T2 any] struct {
	_xgo_vptr unsafe.Pointer
	v1        T1
	v2        T2
	v3        c.Int
}

// llgo:type C
type X_vtable_Bar[T1 any, T2 any] struct {
	XGo_dtor          func(this *Bar)
	XGo_dtor_deleting func(this *Bar)
	G                 func(this *Bar, val1 T1, val2 T2, val3 c.Int) T2
}

func (p *Bar[T1, T2]) XGo_vptr() *X_vtable_Bar[T1, T2] {
	return (*X_vtable_Bar[T1, T2])(p._xgo_vptr)
}
