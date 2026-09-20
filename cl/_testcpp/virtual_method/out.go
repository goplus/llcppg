package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

type Shape struct {
	XGo_vptr unsafe.Pointer
	Id       c.Int
}
type Tag struct {
	Tag c.Int
}
type Widget struct {
	XGo_vptr unsafe.Pointer
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

// llgo:link (*Shape).Area C._ZN5Shape4areaEv
func (this *Shape) Area() c.Int {
	return 0
}

// llgo:link (*Widget).Paint C._ZN6Widget5paintEv
func (this *Widget) Paint() c.Int {
	return 0
}

// llgo:link (*Circle).Area C._ZN6Circle4areaEv
func (this *Circle) Area() c.Int {
	return 0
}

// llgo:link (*Button).Click C._ZN6Button5clickEv
func (this *Button) Click() c.Int {
	return 0
}
