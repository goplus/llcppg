package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

type Shape struct {
	_xgo_vptr unsafe.Pointer
	Id        c.Int
}
type Tag struct {
	Tag c.Int
}
type Widget struct {
	_xgo_vptr unsafe.Pointer
	Tag
	W c.Int
}
type Circle struct {
	Shape
	R c.Int
}
type Button struct {
	Shape
	Tag
	B c.Int
}

// llgo:type C
type _xgo_vtable_Shape struct {
	Area              func(this *Shape) c.Int
	XGo_dtor          func(this *Shape)
	XGo_dtor_deleting func(this *Shape)
	_xgo_slot3        unsafe.Pointer
}

func (p *Shape) XGo_vptr() *_xgo_vtable_Shape {
	return (*_xgo_vtable_Shape)(p._xgo_vptr)
}

// llgo:link (*Shape).Area C._ZN5Shape4areaEv
func (this *Shape) Area() c.Int {
	return 0
}

// llgo:link (*Shape).XGo_Dtor C._ZN5ShapeD1Ev
func (this *Shape) XGo_Dtor() {
}

// llgo:type C
type _xgo_vtable_Widget struct {
	Paint             func(this *Widget) c.Int
	XGo_dtor          func(this *Widget)
	XGo_dtor_deleting func(this *Widget)
}

func (p *Widget) XGo_vptr() *_xgo_vtable_Widget {
	return (*_xgo_vtable_Widget)(p._xgo_vptr)
}

// llgo:link (*Widget).Paint C._ZN6Widget5paintEv
func (this *Widget) Paint() c.Int {
	return 0
}

// llgo:link (*Widget).XGo_Dtor C._ZN6WidgetD1Ev
func (this *Widget) XGo_Dtor() {
}

// llgo:type C
type _xgo_vtable_Circle struct {
	Area              func(this *Circle) c.Int
	XGo_dtor          func(this *Circle)
	XGo_dtor_deleting func(this *Circle)
	_xgo_slot3        unsafe.Pointer
}

func (p *Circle) XGo_vptr() *_xgo_vtable_Circle {
	return (*_xgo_vtable_Circle)(*(*unsafe.Pointer)(unsafe.Pointer(p)))
}

// llgo:link (*Circle).Area C._ZN6Circle4areaEv
func (this *Circle) Area() c.Int {
	return 0
}

// llgo:type C
type _xgo_vtable_Button struct {
	Area              func(this *Button) c.Int
	XGo_dtor          func(this *Button)
	XGo_dtor_deleting func(this *Button)
	_xgo_slot3        unsafe.Pointer
	Click             func(this *Button) c.Int
}

func (p *Button) XGo_vptr() *_xgo_vtable_Button {
	return (*_xgo_vtable_Button)(*(*unsafe.Pointer)(unsafe.Pointer(p)))
}

// llgo:link (*Button).Click C._ZN6Button5clickEv
func (this *Button) Click() c.Int {
	return 0
}
