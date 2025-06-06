package libxslt

import (
	"github.com/goplus/lib/c"
	"github.com/goplus/llcppg/_cmptest/testdata/libxml2/2.13.6/libxml2"
	_ "unsafe"
)

//go:linkname XPathFunctionLookup C.xsltXPathFunctionLookup
func XPathFunctionLookup(vctxt c.Pointer, name *libxml2.Char, ns_uri *libxml2.Char) libxml2.XPathFunction

/*
 * Interfaces for the functions implementations.
 */
//go:linkname DocumentFunction C.xsltDocumentFunction
func DocumentFunction(ctxt libxml2.XPathParserContextPtr, nargs c.Int)

//go:linkname KeyFunction C.xsltKeyFunction
func KeyFunction(ctxt libxml2.XPathParserContextPtr, nargs c.Int)

//go:linkname UnparsedEntityURIFunction C.xsltUnparsedEntityURIFunction
func UnparsedEntityURIFunction(ctxt libxml2.XPathParserContextPtr, nargs c.Int)

//go:linkname FormatNumberFunction C.xsltFormatNumberFunction
func FormatNumberFunction(ctxt libxml2.XPathParserContextPtr, nargs c.Int)

//go:linkname GenerateIdFunction C.xsltGenerateIdFunction
func GenerateIdFunction(ctxt libxml2.XPathParserContextPtr, nargs c.Int)

//go:linkname SystemPropertyFunction C.xsltSystemPropertyFunction
func SystemPropertyFunction(ctxt libxml2.XPathParserContextPtr, nargs c.Int)

//go:linkname ElementAvailableFunction C.xsltElementAvailableFunction
func ElementAvailableFunction(ctxt libxml2.XPathParserContextPtr, nargs c.Int)

//go:linkname FunctionAvailableFunction C.xsltFunctionAvailableFunction
func FunctionAvailableFunction(ctxt libxml2.XPathParserContextPtr, nargs c.Int)

/*
 * And the registration
 */
//go:linkname RegisterAllFunctions C.xsltRegisterAllFunctions
func RegisterAllFunctions(ctxt libxml2.XPathContextPtr)
