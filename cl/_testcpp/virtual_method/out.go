package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

// Case 1: no base class; the class itself has virtual methods.
//
// It also declares a virtual destructor and a private virtual method to
// exercise the two extra vtable slot rules:
//   - a virtual destructor takes two consecutive Itanium slots, XGo_dtor
//     (complete object) then XGo_dtor_deleting, both func(this *Shape);
//   - a private virtual method still occupies a slot, emitted as an unexported
//     placeholder so later slots keep their index.
// Declaration order fixes the vtable order: area, ~Shape (two slots), reset.
type Shape struct {
	_xgo_vptr unsafe.Pointer
	Id        c.Int
}

// A plain base class without any virtual methods.
type Tag struct {
	Tag c.Int
}

// Case 2: has a base class, but the base has no virtual methods; the class
// itself has virtual methods. It also declares a virtual destructor, so its own
// vtable holds paint plus the two destructor slots.
type Widget struct {
	_xgo_vptr unsafe.Pointer
	Tag
	W c.Int
}

// Case 3: single base class, and the base class has virtual methods. The
// derived class reuses the base's vptr, so it introduces no new vptr. It does
// not declare its own destructor, so Shape's destructor and private-method slots
// are inherited unchanged (with this re-typed to *Circle).
type Circle struct {
	Shape
	R c.Int
}

// Case 4: multiple base classes; an earlier base has no virtual methods, while
// a later one does. The polymorphic base (Shape) is the primary base, so the
// class reuses its vptr and introduces no new one; the primary base is laid out
// first (before the non-polymorphic Tag) to keep the shared vptr at offset 0.
// Shape's inherited slots (area, the two destructor slots, the private-method
// placeholder) precede the new click slot.
type Button struct {
	Shape
	Tag
	B c.Int
}

//llgo:cdecl
type X_vtable_Shape struct {
	Area              func(this *Shape) c.Int
	XGo_dtor          func(this *Shape)
	XGo_dtor_deleting func(this *Shape)
	_xgo_slot3        unsafe.Pointer
}

func (p *Shape) XGo_vptr() *X_vtable_Shape {
	return (*X_vtable_Shape)(p._xgo_vptr)
}

//llgo:link (*Shape).Area C._ZN5Shape4areaEv
func (this *Shape) Area() c.Int {
	return 0
}

//llgo:link (*Shape).XGo_Dtor C._ZN5ShapeD1Ev
func (this *Shape) XGo_Dtor() {
}

//llgo:cdecl
type X_vtable_Widget struct {
	Paint             func(this *Widget) c.Int
	XGo_dtor          func(this *Widget)
	XGo_dtor_deleting func(this *Widget)
}

func (p *Widget) XGo_vptr() *X_vtable_Widget {
	return (*X_vtable_Widget)(p._xgo_vptr)
}

//llgo:link (*Widget).Paint C._ZN6Widget5paintEv
func (this *Widget) Paint() c.Int {
	return 0
}

//llgo:link (*Widget).XGo_Dtor C._ZN6WidgetD1Ev
func (this *Widget) XGo_Dtor() {
}

//llgo:cdecl
type X_vtable_Circle struct {
	Area              func(this *Circle) c.Int
	XGo_dtor          func(this *Circle)
	XGo_dtor_deleting func(this *Circle)
	_xgo_slot3        unsafe.Pointer
}

func (p *Circle) XGo_vptr() *X_vtable_Circle {
	return (*X_vtable_Circle)(*(*unsafe.Pointer)(unsafe.Pointer(p)))
}

//llgo:link (*Circle).Area C._ZN6Circle4areaEv
func (this *Circle) Area() c.Int {
	return 0
}

//llgo:cdecl
type X_vtable_Button struct {
	Area              func(this *Button) c.Int
	XGo_dtor          func(this *Button)
	XGo_dtor_deleting func(this *Button)
	_xgo_slot3        unsafe.Pointer
	Click             func(this *Button) c.Int
}

func (p *Button) XGo_vptr() *X_vtable_Button {
	return (*X_vtable_Button)(*(*unsafe.Pointer)(unsafe.Pointer(p)))
}

//llgo:link (*Button).Click C._ZN6Button5clickEv
func (this *Button) Click() c.Int {
	return 0
}
