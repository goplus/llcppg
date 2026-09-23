package CXString

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

const LLGoPackage = "link: -L$(llvm-config --libdir) -lclang; -lclang"

type String struct {
	Data          unsafe.Pointer
	Private_flags c.Uint
}
type StringSet struct {
	Strings *String
	Count   c.Uint
}

//go:linkname GetCString C.clang_getCString
func GetCString(string String) *c.Char

//go:linkname DisposeString C.clang_disposeString
func DisposeString(string String)

//go:linkname DisposeStringSet C.clang_disposeStringSet
func DisposeStringSet(set *StringSet)
