package libxslt

import (
	"github.com/goplus/lib/c"
	"github.com/goplus/llcppg/_cmptest/testdata/libxml2/2.13.6/libxml2"
	_ "unsafe"
)

//go:linkname AddKey C.xsltAddKey
func AddKey(style StylesheetPtr, name *libxml2.Char, nameURI *libxml2.Char, match *libxml2.Char, use *libxml2.Char, inst libxml2.NodePtr) c.Int

//go:linkname GetKey C.xsltGetKey
func GetKey(ctxt TransformContextPtr, name *libxml2.Char, nameURI *libxml2.Char, value *libxml2.Char) libxml2.NodeSetPtr

//go:linkname InitCtxtKeys C.xsltInitCtxtKeys
func InitCtxtKeys(ctxt TransformContextPtr, doc DocumentPtr)

//go:linkname FreeKeys C.xsltFreeKeys
func FreeKeys(style StylesheetPtr)

//go:linkname FreeDocumentKeys C.xsltFreeDocumentKeys
func FreeDocumentKeys(doc DocumentPtr)
