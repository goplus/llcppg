package main

import (
	"gettext"
	"unsafe"

	"github.com/goplus/lib/c"
)

func handle(severity c.Int, message gettext.MessageT, filename *int8, lineno c.Ulong, column c.Ulong, multiline_p c.Int, message_text *int8) {
	c.Printf(c.Str("Error: %s\n"), message_text)
}

func main() {

	handle := gettext.ErrorHandler{
		Error: c.Func(handle),
	}
	file := gettext.FileReadV3(c.Str("example.po"), (gettext.XerrorHandlerT)(unsafe.Pointer(&handle)))
	iter := gettext.MessageIterator(file, nil)
	message := gettext.NextMessage(iter)
	for uintptr(c.Pointer(message)) != 0 {
		msgid := gettext.MessageMsgid(message)
		msgstr := gettext.MessageMsgstr(message)
		c.Printf(c.Str("Original: %s %d\n"), msgid)
		c.Printf(c.Str("Translation: %s\n"), msgstr)
		message = gettext.NextMessage(iter)
	}
	gettext.MessageIteratorFree(iter)
	gettext.FileFree(file)
}
