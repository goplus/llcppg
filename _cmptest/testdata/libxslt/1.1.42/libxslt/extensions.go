package libxslt

import (
	"github.com/goplus/lib/c"
	"github.com/goplus/llcppg/_cmptest/testdata/libxml2/2.13.6/libxml2"
	_ "unsafe"
)

/**
 * xsltInitGlobals:
 *
 * Initialize the global variables for extensions
 *
 */
//go:linkname InitGlobals C.xsltInitGlobals
func InitGlobals()

// llgo:type C
type StyleExtInitFunction func(StylesheetPtr, *libxml2.Char) c.Pointer

// llgo:type C
type StyleExtShutdownFunction func(StylesheetPtr, *libxml2.Char, c.Pointer)

// llgo:type C
type ExtInitFunction func(TransformContextPtr, *libxml2.Char) c.Pointer

// llgo:type C
type ExtShutdownFunction func(TransformContextPtr, *libxml2.Char, c.Pointer)

//go:linkname RegisterExtModule C.xsltRegisterExtModule
func RegisterExtModule(URI *libxml2.Char, initFunc ExtInitFunction, shutdownFunc ExtShutdownFunction) c.Int

//go:linkname RegisterExtModuleFull C.xsltRegisterExtModuleFull
func RegisterExtModuleFull(URI *libxml2.Char, initFunc ExtInitFunction, shutdownFunc ExtShutdownFunction, styleInitFunc StyleExtInitFunction, styleShutdownFunc StyleExtShutdownFunction) c.Int

//go:linkname UnregisterExtModule C.xsltUnregisterExtModule
func UnregisterExtModule(URI *libxml2.Char) c.Int

//go:linkname GetExtData C.xsltGetExtData
func GetExtData(ctxt TransformContextPtr, URI *libxml2.Char) c.Pointer

//go:linkname StyleGetExtData C.xsltStyleGetExtData
func StyleGetExtData(style StylesheetPtr, URI *libxml2.Char) c.Pointer

//go:linkname ShutdownCtxtExts C.xsltShutdownCtxtExts
func ShutdownCtxtExts(ctxt TransformContextPtr)

//go:linkname ShutdownExts C.xsltShutdownExts
func ShutdownExts(style StylesheetPtr)

//go:linkname XPathGetTransformContext C.xsltXPathGetTransformContext
func XPathGetTransformContext(ctxt libxml2.XPathParserContextPtr) TransformContextPtr

/*
 * extension functions
 */
//go:linkname RegisterExtModuleFunction C.xsltRegisterExtModuleFunction
func RegisterExtModuleFunction(name *libxml2.Char, URI *libxml2.Char, function libxml2.XPathFunction) c.Int

//go:linkname ExtModuleFunctionLookup C.xsltExtModuleFunctionLookup
func ExtModuleFunctionLookup(name *libxml2.Char, URI *libxml2.Char) libxml2.XPathFunction

//go:linkname UnregisterExtModuleFunction C.xsltUnregisterExtModuleFunction
func UnregisterExtModuleFunction(name *libxml2.Char, URI *libxml2.Char) c.Int

// llgo:type C
type PreComputeFunction func(StylesheetPtr, libxml2.NodePtr, TransformFunction) ElemPreCompPtr

//go:linkname NewElemPreComp C.xsltNewElemPreComp
func NewElemPreComp(style StylesheetPtr, inst libxml2.NodePtr, function TransformFunction) ElemPreCompPtr

//go:linkname InitElemPreComp C.xsltInitElemPreComp
func InitElemPreComp(comp ElemPreCompPtr, style StylesheetPtr, inst libxml2.NodePtr, function TransformFunction, freeFunc ElemPreCompDeallocator)

//go:linkname RegisterExtModuleElement C.xsltRegisterExtModuleElement
func RegisterExtModuleElement(name *libxml2.Char, URI *libxml2.Char, precomp PreComputeFunction, transform TransformFunction) c.Int

//go:linkname ExtElementLookup C.xsltExtElementLookup
func ExtElementLookup(ctxt TransformContextPtr, name *libxml2.Char, URI *libxml2.Char) TransformFunction

//go:linkname ExtModuleElementLookup C.xsltExtModuleElementLookup
func ExtModuleElementLookup(name *libxml2.Char, URI *libxml2.Char) TransformFunction

//go:linkname ExtModuleElementPreComputeLookup C.xsltExtModuleElementPreComputeLookup
func ExtModuleElementPreComputeLookup(name *libxml2.Char, URI *libxml2.Char) PreComputeFunction

//go:linkname UnregisterExtModuleElement C.xsltUnregisterExtModuleElement
func UnregisterExtModuleElement(name *libxml2.Char, URI *libxml2.Char) c.Int

// llgo:type C
type TopLevelFunction func(StylesheetPtr, libxml2.NodePtr)

//go:linkname RegisterExtModuleTopLevel C.xsltRegisterExtModuleTopLevel
func RegisterExtModuleTopLevel(name *libxml2.Char, URI *libxml2.Char, function TopLevelFunction) c.Int

//go:linkname ExtModuleTopLevelLookup C.xsltExtModuleTopLevelLookup
func ExtModuleTopLevelLookup(name *libxml2.Char, URI *libxml2.Char) TopLevelFunction

//go:linkname UnregisterExtModuleTopLevel C.xsltUnregisterExtModuleTopLevel
func UnregisterExtModuleTopLevel(name *libxml2.Char, URI *libxml2.Char) c.Int

/* These 2 functions are deprecated for use within modules. */
//go:linkname RegisterExtFunction C.xsltRegisterExtFunction
func RegisterExtFunction(ctxt TransformContextPtr, name *libxml2.Char, URI *libxml2.Char, function libxml2.XPathFunction) c.Int

//go:linkname RegisterExtElement C.xsltRegisterExtElement
func RegisterExtElement(ctxt TransformContextPtr, name *libxml2.Char, URI *libxml2.Char, function TransformFunction) c.Int

/*
 * Extension Prefix handling API.
 * Those are used by the XSLT (pre)processor.
 */
//go:linkname RegisterExtPrefix C.xsltRegisterExtPrefix
func RegisterExtPrefix(style StylesheetPtr, prefix *libxml2.Char, URI *libxml2.Char) c.Int

//go:linkname CheckExtPrefix C.xsltCheckExtPrefix
func CheckExtPrefix(style StylesheetPtr, URI *libxml2.Char) c.Int

//go:linkname CheckExtURI C.xsltCheckExtURI
func CheckExtURI(style StylesheetPtr, URI *libxml2.Char) c.Int

//go:linkname InitCtxtExts C.xsltInitCtxtExts
func InitCtxtExts(ctxt TransformContextPtr) c.Int

//go:linkname FreeCtxtExts C.xsltFreeCtxtExts
func FreeCtxtExts(ctxt TransformContextPtr)

//go:linkname FreeExts C.xsltFreeExts
func FreeExts(style StylesheetPtr)

//go:linkname PreComputeExtModuleElement C.xsltPreComputeExtModuleElement
func PreComputeExtModuleElement(style StylesheetPtr, inst libxml2.NodePtr) ElemPreCompPtr

/*
 * Extension Infos access.
 * Used by exslt initialisation
 */
//go:linkname GetExtInfo C.xsltGetExtInfo
func GetExtInfo(style StylesheetPtr, URI *libxml2.Char) libxml2.HashTablePtr

/**
 * Test of the extension module API
 */
//go:linkname RegisterTestModule C.xsltRegisterTestModule
func RegisterTestModule()

//go:linkname DebugDumpExtensions C.xsltDebugDumpExtensions
func DebugDumpExtensions(output *c.FILE)
