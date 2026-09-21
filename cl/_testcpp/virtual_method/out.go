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
type Stream struct {
	_xgo_vptr unsafe.Pointer
	Fd        c.Int
}
type Canvas struct {
	_xgo_vptr unsafe.Pointer
}
type Machine struct {
	_xgo_vptr unsafe.Pointer
}

//llgo:type C
type _xgo_vtable_Shape struct {
	Area func(this *Shape) c.Int
}

func (p *Shape) XGo_vptr() *_xgo_vtable_Shape {
	return (*_xgo_vtable_Shape)(p._xgo_vptr)
}

// llgo:link (*Shape).Area C._ZN5Shape4areaEv
func (this *Shape) Area() c.Int {
	return 0
}

//llgo:type C
type _xgo_vtable_Widget struct {
	Paint func(this *Widget) c.Int
}

func (p *Widget) XGo_vptr() *_xgo_vtable_Widget {
	return (*_xgo_vtable_Widget)(p._xgo_vptr)
}

// llgo:link (*Widget).Paint C._ZN6Widget5paintEv
func (this *Widget) Paint() c.Int {
	return 0
}

//llgo:type C
type _xgo_vtable_Circle struct {
	Area func(this *Circle) c.Int
}

func (p *Circle) XGo_vptr() *_xgo_vtable_Circle {
	return (*_xgo_vtable_Circle)(*(*unsafe.Pointer)(unsafe.Pointer(p)))
}

// llgo:link (*Circle).Area C._ZN6Circle4areaEv
func (this *Circle) Area() c.Int {
	return 0
}

//llgo:type C
type _xgo_vtable_Button struct {
	Area  func(this *Button) c.Int
	Click func(this *Button) c.Int
}

func (p *Button) XGo_vptr() *_xgo_vtable_Button {
	return (*_xgo_vtable_Button)(*(*unsafe.Pointer)(unsafe.Pointer(p)))
}

// llgo:link (*Button).Click C._ZN6Button5clickEv
func (this *Button) Click() c.Int {
	return 0
}

// llgo:link (*Stream).XGo_Dtor C._ZN6StreamD1Ev
func (this *Stream) XGo_Dtor() {
}

//llgo:type C
type _xgo_vtable_Canvas struct {
	_xgo_slot0 unsafe.Pointer
	_xgo_slot1 unsafe.Pointer
	Draw       func(this *Canvas) c.Int
}

func (p *Canvas) XGo_vptr() *_xgo_vtable_Canvas {
	return (*_xgo_vtable_Canvas)(p._xgo_vptr)
}

// llgo:link (*Canvas).XGo_Dtor C._ZN6CanvasD1Ev
func (this *Canvas) XGo_Dtor() {
}

// llgo:link (*Canvas).Draw C._ZN6Canvas4drawEv
func (this *Canvas) Draw() c.Int {
	return 0
}

//llgo:type C
type _xgo_vtable_Machine struct {
	Start      func(this *Machine) c.Int
	_xgo_slot1 unsafe.Pointer
	Stop       func(this *Machine) c.Int
}

func (p *Machine) XGo_vptr() *_xgo_vtable_Machine {
	return (*_xgo_vtable_Machine)(p._xgo_vptr)
}

// llgo:link (*Machine).Start C._ZN7Machine5startEv
func (this *Machine) Start() c.Int {
	return 0
}

// llgo:link (*Machine).Stop C._ZN7Machine4stopEv
func (this *Machine) Stop() c.Int {
	return 0
}
