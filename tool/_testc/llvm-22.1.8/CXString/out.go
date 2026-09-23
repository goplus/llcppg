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

//go:linkname Clang_getCString C.clang_getCString
func Clang_getCString(string String) *c.Char

//go:linkname Clang_disposeString C.clang_disposeString
func Clang_disposeString(string String)

//go:linkname Clang_disposeStringSet C.clang_disposeStringSet
func Clang_disposeStringSet(set *StringSet)
