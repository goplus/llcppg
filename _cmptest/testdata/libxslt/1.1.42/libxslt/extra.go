package libxslt

import (
	"github.com/goplus/lib/c"
	"github.com/goplus/llcppg/_cmptest/testdata/libxml2/2.13.6/libxml2"
	_ "unsafe"
)

//go:linkname FunctionNodeSet C.xsltFunctionNodeSet
func FunctionNodeSet(ctxt libxml2.XPathParserContextPtr, nargs c.Int)

//go:linkname Debug C.xsltDebug
func Debug(ctxt TransformContextPtr, node libxml2.NodePtr, inst libxml2.NodePtr, comp ElemPreCompPtr)

//go:linkname RegisterExtras C.xsltRegisterExtras
func RegisterExtras(ctxt TransformContextPtr)

//go:linkname RegisterAllExtras C.xsltRegisterAllExtras
func RegisterAllExtras()
