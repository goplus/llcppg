package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

const XGoPackage = true

type Tag struct {
	Tag c.Int // The non-polymorphic base's field.
}
type Shape struct {
	_xgo_vptr unsafe.Pointer
	// The shape identifier.
	Id c.Int // After the implicit vptr.
}
type Derived struct {
	Shape
	Tag
	Value c.Int // After the reordered base classes.
}
type Access struct {
	hidden  c.Int // A private field.
	Visible c.Int // A public field.
}
type Value struct {
	_xgo_union [2]uint32
}
type Box[T any] struct {
	// The template's stored value.
	Value T // An inline template field comment.
}
type _llcppg_anon_0 struct {
	First  c.Int // A promoted C++ struct member.
	Second c.Int // A promoted C++ struct member at a nonzero offset.
}

//llgo:cdecl
type X_vtable_Shape struct {
	Size func(this *Shape) c.Int
}

func (p *Shape) XGo_vptr() *X_vtable_Shape {
	return (*X_vtable_Shape)(p._xgo_vptr)
}

//llgo:link (*Shape).Size C._ZN5Shape4sizeEv
func (this *Shape) Size() c.Int {
	return 0
}

//llgo:cdecl
type X_vtable_Derived struct {
	Size func(this *Derived) c.Int
}

func (p *Derived) XGo_vptr() *X_vtable_Derived {
	return (*X_vtable_Derived)(*(*unsafe.Pointer)(unsafe.Pointer(p)))
}

// A C++ union member.
func (p *Value) XGof_ref_number() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}

// A promoted C++ struct member.
func (p *Value) XGof_ref_first() *c.Int {
	return &(*_llcppg_anon_0)(unsafe.Pointer(p)).First
}

// A promoted C++ struct member at a nonzero offset.
func (p *Value) XGof_ref_second() *c.Int {
	return &(*_llcppg_anon_0)(unsafe.Pointer(p)).Second
}
