package libc

import (
	"github.com/goplus/lib/c"
)

//llgo:link (*FILE).Fprintf C.fprintf
func (self *FILE) Fprintf(_llcppg_param2 *c.Char, __llgo_va_list ...any) c.Int {
	return 0
}
