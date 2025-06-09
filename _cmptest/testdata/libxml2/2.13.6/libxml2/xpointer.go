package libxml2

import _ "unsafe"

/*
 * Functions.
 */
//go:linkname XPtrNewContext C.xmlXPtrNewContext
func XPtrNewContext(doc DocPtr, here NodePtr, origin NodePtr) XPathContextPtr

// llgo:link (*Char).XPtrEval C.xmlXPtrEval
func (recv_ *Char) XPtrEval(ctx XPathContextPtr) XPathObjectPtr {
	return nil
}
