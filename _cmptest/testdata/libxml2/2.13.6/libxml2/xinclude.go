package libxml2

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type X_xmlXIncludeCtxt struct {
	Unused [8]uint8
}
type XIncludeCtxt X_xmlXIncludeCtxt
type XIncludeCtxtPtr *XIncludeCtxt

/*
 * standalone processing
 */
//go:linkname XIncludeProcess C.xmlXIncludeProcess
func XIncludeProcess(doc DocPtr) c.Int

//go:linkname XIncludeProcessFlags C.xmlXIncludeProcessFlags
func XIncludeProcessFlags(doc DocPtr, flags c.Int) c.Int

//go:linkname XIncludeProcessFlagsData C.xmlXIncludeProcessFlagsData
func XIncludeProcessFlagsData(doc DocPtr, flags c.Int, data c.Pointer) c.Int

//go:linkname XIncludeProcessTreeFlagsData C.xmlXIncludeProcessTreeFlagsData
func XIncludeProcessTreeFlagsData(tree NodePtr, flags c.Int, data c.Pointer) c.Int

//go:linkname XIncludeProcessTree C.xmlXIncludeProcessTree
func XIncludeProcessTree(tree NodePtr) c.Int

//go:linkname XIncludeProcessTreeFlags C.xmlXIncludeProcessTreeFlags
func XIncludeProcessTreeFlags(tree NodePtr, flags c.Int) c.Int

/*
 * contextual processing
 */
//go:linkname XIncludeNewContext C.xmlXIncludeNewContext
func XIncludeNewContext(doc DocPtr) XIncludeCtxtPtr

//go:linkname XIncludeSetFlags C.xmlXIncludeSetFlags
func XIncludeSetFlags(ctxt XIncludeCtxtPtr, flags c.Int) c.Int

//go:linkname XIncludeSetErrorHandler C.xmlXIncludeSetErrorHandler
func XIncludeSetErrorHandler(ctxt XIncludeCtxtPtr, handler StructuredErrorFunc, data c.Pointer)

//go:linkname XIncludeGetLastError C.xmlXIncludeGetLastError
func XIncludeGetLastError(ctxt XIncludeCtxtPtr) c.Int

//go:linkname XIncludeFreeContext C.xmlXIncludeFreeContext
func XIncludeFreeContext(ctxt XIncludeCtxtPtr)

//go:linkname XIncludeProcessNode C.xmlXIncludeProcessNode
func XIncludeProcessNode(ctxt XIncludeCtxtPtr, tree NodePtr) c.Int
