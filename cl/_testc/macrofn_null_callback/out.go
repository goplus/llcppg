package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

type Owner struct {
}

//llgo:cdecl
type CallbackT = func(_llcppg_param1 c.Int) c.Int
type HandleT uintptr
type OwnerT Owner

//go:linkname InvokeCallback C.invoke_callback
func InvokeCallback(mode c.Int, callback CallbackT) c.Int

//go:linkname SetOptions C.set_options
func SetOptions(count c.Int, pointer *c.Int, enabled bool)

//go:linkname SetFlag C.set_flag
func SetFlag(mode c.Int, enabled bool)

//go:linkname SetContext C.set_context
func SetContext(mode c.Int, context unsafe.Pointer)

//go:linkname UseHandle C.use_handle
func UseHandle(marker c.Int, value HandleT) c.Int

//go:linkname LogValue C.log_value
func LogValue(count c.Int, __llgo_va_list ...any) c.Int

//go:linkname SetFixed C.set_fixed
func SetFixed(count c.Int, values *c.Int)

//go:linkname SetOpen C.set_open
func SetOpen(count c.Int, values *c.Int)

//go:linkname SetVariable C.set_variable
func SetVariable(count c.Int, values *c.Int)

//go:linkname InvokeDirect C.invoke_direct
func InvokeDirect(mode c.Int, callback func(_llcppg_param1 c.Int) c.Int) c.Int

//go:linkname OpenCallbacks C.open_callbacks
func OpenCallbacks(cookie unsafe.Pointer, read func(_llcppg_param1 unsafe.Pointer, _llcppg_param2 *c.Char, _llcppg_param3 c.Int) c.Int, write func(_llcppg_param1 unsafe.Pointer, _llcppg_param2 *c.Char, _llcppg_param3 c.Int) c.Int, seek func(_llcppg_param1 unsafe.Pointer, _llcppg_param2 c.Long, _llcppg_param3 c.Int) c.Long, close func(_llcppg_param1 unsafe.Pointer) c.Int) c.Int

//llgo:link (*OwnerT).OwnerCallback C.owner_callback
func (owner *OwnerT) OwnerCallback(mode c.Int, callback CallbackT) c.Int {
	return 0
}

//llgo:link (*OwnerT).OwnerConfigure C.owner_configure
func (owner *OwnerT) OwnerConfigure(value *c.Int, enabled bool) {
}
func InvokeEmpty(mode c.Int) c.Int {
	return InvokeCallback(mode, nil)
}
func InvokeZero(callback CallbackT) c.Int {
	return InvokeCallback(0, callback)
}
func SetDefaultOptions() {
	SetOptions(0, nil, false)
}
func SetParenthesizedOptions() {
	SetOptions(0, nil, false)
}
func SetDefaultPointer(pointer *c.Int) {
	SetOptions(256, pointer, false)
}
func SetComputedOptions() {
	SetOptions(256, nil, false)
}
func SetFloatZero() {
	SetFlag(0, false)
}
func SetFloatNegativeZero() {
	SetFlag(0, false)
}
func SetEmptyContext() {
	SetContext(256, nil)
}
func SetComputedContext() {
	SetContext(256, nil)
}
func UseEmptyHandle() c.Int {
	return UseHandle(256, 0)
}
func UseComputedHandle() c.Int {
	return UseHandle(256, 0)
}
func LogZero() c.Int {
	return LogValue(1, 0)
}
func LogComputed() c.Int {
	return LogValue(1, 0)
}
func SetFixedEmpty() {
	SetFixed(1, nil)
}
func SetFixedComputed() {
	SetFixed(1, nil)
}
func SetOpenEmpty() {
	SetOpen(1, nil)
}
func SetVariableEmpty() {
	SetVariable(1, nil)
}
func InvokeDirectEmpty() c.Int {
	return InvokeDirect(256, nil)
}
func InvokeDirectComputed() c.Int {
	return InvokeDirect(256, nil)
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
func (owner *OwnerT) OwnerReset() {
	owner.OwnerConfigure(nil, false)
}
