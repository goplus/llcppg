package foo

import "github.com/goplus/lib/c"

type SizeT = c.Ulong
type WeakPtr[_Tp any] struct {
	__Ptr_   *_Tp
	__Cntrl_ *SizeT
}
type SharedPtr[_Tp any] struct {
	__Ptr_   *_Tp
	__Cntrl_ *SizeT
}
type SharedPtrWeakType[_Tp any] = WeakPtr[_Tp]
