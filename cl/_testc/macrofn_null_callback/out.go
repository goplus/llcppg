package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

type Owner struct {
}

//llgo:cdecl
type CallbackT = func(_llcppg_param1 c.Int) c.Int
type OwnerT Owner

//go:linkname InvokeCallback C.invoke_callback
func InvokeCallback(mode c.Int, callback CallbackT) c.Int

//go:linkname OpenCallbacks C.open_callbacks
func OpenCallbacks(cookie unsafe.Pointer, read func(_llcppg_param1 unsafe.Pointer, _llcppg_param2 *c.Char, _llcppg_param3 c.Int) c.Int, write func(_llcppg_param1 unsafe.Pointer, _llcppg_param2 *c.Char, _llcppg_param3 c.Int) c.Int, seek func(_llcppg_param1 unsafe.Pointer, _llcppg_param2 c.Long, _llcppg_param3 c.Int) c.Long, close func(_llcppg_param1 unsafe.Pointer) c.Int) c.Int

//llgo:link (*OwnerT).OwnerCallback C.owner_callback
func (owner *OwnerT) OwnerCallback(mode c.Int, callback CallbackT) c.Int {
	return 0
}
func InvokeEmpty(mode c.Int) c.Int {
	return InvokeCallback(mode, nil)
}
func InvokeZero(callback CallbackT) c.Int {
	return InvokeCallback(0, callback)
}
func OpenRead(cookie unsafe.Pointer, fn func(_llcppg_param1 unsafe.Pointer, _llcppg_param2 *c.Char, _llcppg_param3 c.Int) c.Int) c.Int {
	return OpenCallbacks(cookie, fn, nil, nil, nil)
}
func OpenWrite(cookie unsafe.Pointer, fn func(_llcppg_param1 unsafe.Pointer, _llcppg_param2 *c.Char, _llcppg_param3 c.Int) c.Int) c.Int {
	return OpenCallbacks(cookie, nil, fn, nil, nil)
}
func OpenEmpty(cookie unsafe.Pointer) c.Int {
	return OpenCallbacks(cookie, nil, nil, nil, nil)
}
func (owner *OwnerT) OwnerEmpty() c.Int {
	return owner.OwnerCallback(0, nil)
}
