package CXString

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

const LLGoPackage = "link: -L$(llvm-config --libdir) -lclang; -lclang"

type String struct {
	Data         unsafe.Pointer
	PrivateFlags c.Uint
}
type StringSet struct {
	Strings *String
	Count   c.Uint
}

// llgo:link String.CStr C.clang_getCString
func (string String) CStr() *c.Char {
	return nil
}

// llgo:link String.Dispose C.clang_disposeString
func (string String) Dispose() {
}

// llgo:link (*StringSet).Dispose C.clang_disposeStringSet
func (set *StringSet) Dispose() {
}
