package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

const XGoPackage = true

// The holder's documentation stays on its type.
type Holder struct {
	Before   c.Int
	Callback c.Pointer // int (*callback)()
	// Keep both comments on this field.
	// int (*documented)()
	Documented      c.Pointer       // The original inline comment.
	Named           Callback        // Callback named
	Alias           CallbackAlias   // CallbackAlias alias
	Function        Function        // Function *function
	FunctionPointer FunctionPointer // FunctionPointer function_pointer
	Pointer         *c.Pointer      // int (**pointer)()
	Callbacks       [2]c.Pointer    // int (*callbacks[2])()
	Fixed           c.Pointer       // int (*const fixed)()
	First           c.Pointer       // int (*first)()
	Second          c.Pointer       // int (*second)()
	After           c.Int
}
type Nested struct {
	Entries [2]_llcppg_anon_0 // The record array's own comment.
	After   c.Int
}
type Value struct {
	_xgo_union [1]uint64
}
type Overrides struct {
	Function *CustomFunction // CustomFunction *function
	Callback CustomCallback  // CustomCallback callback
	Renamed  c.Pointer       // int (*renamed)()
}

//llgo:cdecl
type Prototypes struct {
	Empty    func() c.Int
	Typed    func(_llcppg_param1 c.Int, _llcppg_param2 c.Double) c.Int
	Variadic func(_llcppg_param1 c.Int, __llgo_va_list ...any) c.Int
}

// A callback with unspecified parameters.
// typedef int (*Callback)()
type Callback = c.Pointer

// typedef Callback CallbackAlias
type CallbackAlias = Callback

// typedef int (Function)()
type Function = c.Pointer

// typedef Function *FunctionPointer
type FunctionPointer = Function

// typedef void (*Visitor)(int (*)())
//
//llgo:cdecl
type Visitor = func(_llcppg_param1 c.Pointer)
type _llcppg_anon_0 struct {
	Callback c.Pointer // int (*callback)()
	Value    c.Int
}
type _llcppg_anon_1 struct {
	// int (*nested)()
	Nested c.Pointer // The nested callback's comment.
}

// typedef int (CustomFunction)()
type CustomFunction = int32

// typedef int (*CustomCallback)()
type CustomCallback = uintptr

// int (*callback)()
func (p *Value) XGof_ref_callback() *c.Pointer {
	return (*c.Pointer)(unsafe.Pointer(p))
}
func (p *Value) XGof_ref_number() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}

// int (*nested)()
// The nested callback's comment.
func (p *Value) XGof_ref_nested() *c.Pointer {
	return &(*_llcppg_anon_1)(unsafe.Pointer(p)).Nested
}

// The active callback.
// extern int (*global_callback)()
//
//go:linkname GlobalCallback C.global_callback
var GlobalCallback c.Pointer

// void accept_callback(int (*callback)())
//
//go:linkname AcceptCallback C.accept_callback
func AcceptCallback(callback c.Pointer)

// int (*get_callback(void))()
//
//go:linkname GetCallback C.get_callback
func GetCallback() c.Pointer

// void accept_visitor(void (*visitor)(int (*)()))
//
//go:linkname AcceptVisitor C.accept_visitor
func AcceptVisitor(visitor func(_llcppg_param1 c.Pointer))
