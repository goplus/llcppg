package libxslt

import (
	"github.com/goplus/lib/c"
	"github.com/goplus/llcppg/_cmptest/testdata/libxml2/2.13.6/libxml2"
	_ "unsafe"
)

/*
 * Module interfaces
 */
//go:linkname ParseStylesheetImport C.xsltParseStylesheetImport
func ParseStylesheetImport(style StylesheetPtr, cur libxml2.NodePtr) c.Int

//go:linkname ParseStylesheetInclude C.xsltParseStylesheetInclude
func ParseStylesheetInclude(style StylesheetPtr, cur libxml2.NodePtr) c.Int

//go:linkname NextImport C.xsltNextImport
func NextImport(style StylesheetPtr) StylesheetPtr

//go:linkname NeedElemSpaceHandling C.xsltNeedElemSpaceHandling
func NeedElemSpaceHandling(ctxt TransformContextPtr) c.Int

//go:linkname FindElemSpaceHandling C.xsltFindElemSpaceHandling
func FindElemSpaceHandling(ctxt TransformContextPtr, node libxml2.NodePtr) c.Int

//go:linkname FindTemplate C.xsltFindTemplate
func FindTemplate(ctxt TransformContextPtr, name *libxml2.Char, nameURI *libxml2.Char) TemplatePtr
