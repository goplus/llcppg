package clang

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

// A character string.
//
// The \c CXString type is used to return strings from the interface when
// the ownership of that string might differ from one call to the next.
// Use \c clang_getCString() to retrieve the string data and, once finished
// with the string data, call \c clang_disposeString() to free the string.
type String struct {
	Data         unsafe.Pointer
	PrivateFlags c.Uint
}
type StringSet struct {
	Strings *String
	Count   c.Uint
}

// Retrieve the character data associated with the given string.
//
// The returned data is a reference and not owned by the user. This data
// is only valid while the `CXString` is valid. This function is similar
// to `std::string::c_str()`.
//
// llgo:link String.CStr C.clang_getCString
func (string String) CStr() *c.Char {
	return nil
}

// Free the given string.
//
// llgo:link String.Dispose C.clang_disposeString
func (string String) Dispose() {
}

// Free the given string set.
//
// llgo:link (*StringSet).Dispose C.clang_disposeStringSet
func (set *StringSet) Dispose() {
}
