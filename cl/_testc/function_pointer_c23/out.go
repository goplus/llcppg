package foo

import "github.com/goplus/lib/c"

//llgo:cdecl
type Holder struct {
	Callback func() c.Int
	Named    Callback
	Empty    func() c.Int
	Typed    func(_llcppg_param1 c.Int, _llcppg_param2 c.Double) c.Int
	Variadic func(_llcppg_param1 c.Int, __llgo_va_list ...any) c.Int
}

//llgo:cdecl
type Callback = func() c.Int
