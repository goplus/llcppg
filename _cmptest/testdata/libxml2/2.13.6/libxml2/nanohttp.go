package libxml2

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

//go:linkname NanoHTTPInit C.xmlNanoHTTPInit
func NanoHTTPInit()

//go:linkname NanoHTTPCleanup C.xmlNanoHTTPCleanup
func NanoHTTPCleanup()

//go:linkname NanoHTTPScanProxy C.xmlNanoHTTPScanProxy
func NanoHTTPScanProxy(URL *c.Char)

//go:linkname NanoHTTPFetch C.xmlNanoHTTPFetch
func NanoHTTPFetch(URL *c.Char, filename *c.Char, contentType **c.Char) c.Int

//go:linkname NanoHTTPMethod C.xmlNanoHTTPMethod
func NanoHTTPMethod(URL *c.Char, method *c.Char, input *c.Char, contentType **c.Char, headers *c.Char, ilen c.Int) c.Pointer

//go:linkname NanoHTTPMethodRedir C.xmlNanoHTTPMethodRedir
func NanoHTTPMethodRedir(URL *c.Char, method *c.Char, input *c.Char, contentType **c.Char, redir **c.Char, headers *c.Char, ilen c.Int) c.Pointer

//go:linkname NanoHTTPOpen C.xmlNanoHTTPOpen
func NanoHTTPOpen(URL *c.Char, contentType **c.Char) c.Pointer

//go:linkname NanoHTTPOpenRedir C.xmlNanoHTTPOpenRedir
func NanoHTTPOpenRedir(URL *c.Char, contentType **c.Char, redir **c.Char) c.Pointer

//go:linkname NanoHTTPReturnCode C.xmlNanoHTTPReturnCode
func NanoHTTPReturnCode(ctx c.Pointer) c.Int

//go:linkname NanoHTTPAuthHeader C.xmlNanoHTTPAuthHeader
func NanoHTTPAuthHeader(ctx c.Pointer) *c.Char

//go:linkname NanoHTTPRedir C.xmlNanoHTTPRedir
func NanoHTTPRedir(ctx c.Pointer) *c.Char

//go:linkname NanoHTTPContentLength C.xmlNanoHTTPContentLength
func NanoHTTPContentLength(ctx c.Pointer) c.Int

//go:linkname NanoHTTPEncoding C.xmlNanoHTTPEncoding
func NanoHTTPEncoding(ctx c.Pointer) *c.Char

//go:linkname NanoHTTPMimeType C.xmlNanoHTTPMimeType
func NanoHTTPMimeType(ctx c.Pointer) *c.Char

//go:linkname NanoHTTPRead C.xmlNanoHTTPRead
func NanoHTTPRead(ctx c.Pointer, dest c.Pointer, len c.Int) c.Int

//go:linkname NanoHTTPSave C.xmlNanoHTTPSave
func NanoHTTPSave(ctxt c.Pointer, filename *c.Char) c.Int

//go:linkname NanoHTTPClose C.xmlNanoHTTPClose
func NanoHTTPClose(ctx c.Pointer)
