package libxslt

import (
	"github.com/goplus/llcppg/_cmptest/testdata/libxml2/2.13.6/libxml2"
	_ "unsafe"
)

//go:linkname ParseStylesheetAttributeSet C.xsltParseStylesheetAttributeSet
func ParseStylesheetAttributeSet(style StylesheetPtr, cur libxml2.NodePtr)

//go:linkname FreeAttributeSetsHashes C.xsltFreeAttributeSetsHashes
func FreeAttributeSetsHashes(style StylesheetPtr)

//go:linkname ApplyAttributeSet C.xsltApplyAttributeSet
func ApplyAttributeSet(ctxt TransformContextPtr, node libxml2.NodePtr, inst libxml2.NodePtr, attributes *libxml2.Char)

//go:linkname ResolveStylesheetAttributeSet C.xsltResolveStylesheetAttributeSet
func ResolveStylesheetAttributeSet(style StylesheetPtr)
