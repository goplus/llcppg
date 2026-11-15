package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

const TEXT = "_"

//go:linkname SendText C.send_text
func SendText(text *c.Char)

//go:linkname SendPointer C.send_pointer
func SendPointer(pointer unsafe.Pointer)
func SendLiteral() {
	SendText(c.Str("_"))
}
func SendValue(text *c.Char) {
	SendText(text)
}
func SendNull() {
	SendPointer(nil)
}
