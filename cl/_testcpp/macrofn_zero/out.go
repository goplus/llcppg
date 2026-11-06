package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

//llgo:cdecl
type ZeroCallbackT = func(_llcppg_param1 c.Int) c.Int
type HandleT uintptr

//go:linkname SetOptions C._Z11set_optionsiPib
func SetOptions(count c.Int, pointer *c.Int, enabled bool)

//go:linkname SetCallback C._Z12set_callbackiPFiiE
func SetCallback(mode c.Int, callback ZeroCallbackT) c.Int

//go:linkname SetContext C._Z11set_contextiPv
func SetContext(mode c.Int, context unsafe.Pointer)

//go:linkname UseHandle C._Z10use_handleiPv
func UseHandle(marker c.Int, value HandleT) c.Int

//go:linkname SetArray C._Z9set_arrayiPi
func SetArray(count c.Int, values *c.Int)

//go:linkname LogValue C._Z9log_valueiz
func LogValue(count c.Int, __llgo_va_list ...any) c.Int
func DefaultOptions() {
	SetOptions(0, nil, false)
}
func FloatOptions() {
	SetOptions(256, nil, false)
}
func PointerOptions(pointer *c.Int) {
	SetOptions(256, pointer, false)
}
func DefaultCallback() c.Int {
	return SetCallback(256, nil)
}
func DefaultContext() {
	SetContext(256, nil)
}
func DefaultHandle() c.Int {
	return UseHandle(256, 0)
}
func DefaultArray() {
	SetArray(1, nil)
}
func LogZero() c.Int {
	return LogValue(1, 0)
}
