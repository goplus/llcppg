package libxslt

import (
	"github.com/goplus/lib/c"
	"github.com/goplus/llcppg/_cmptest/testdata/libxml2/2.13.6/libxml2"
	_ "unsafe"
)

type CompMatch X_xsltCompMatch
type CompMatchPtr *CompMatch

/*
 * Pattern related interfaces.
 */
//go:linkname CompilePattern C.xsltCompilePattern
func CompilePattern(pattern *libxml2.Char, doc libxml2.DocPtr, node libxml2.NodePtr, style StylesheetPtr, runtime TransformContextPtr) CompMatchPtr

//go:linkname FreeCompMatchList C.xsltFreeCompMatchList
func FreeCompMatchList(comp CompMatchPtr)

//go:linkname TestCompMatchList C.xsltTestCompMatchList
func TestCompMatchList(ctxt TransformContextPtr, node libxml2.NodePtr, comp CompMatchPtr) c.Int

//go:linkname CompMatchClearCache C.xsltCompMatchClearCache
func CompMatchClearCache(ctxt TransformContextPtr, comp CompMatchPtr)

//go:linkname NormalizeCompSteps C.xsltNormalizeCompSteps
func NormalizeCompSteps(payload c.Pointer, data c.Pointer, name *libxml2.Char)

/*
 * Template related interfaces.
 */
//go:linkname AddTemplate C.xsltAddTemplate
func AddTemplate(style StylesheetPtr, cur TemplatePtr, mode *libxml2.Char, modeURI *libxml2.Char) c.Int

//go:linkname GetTemplate C.xsltGetTemplate
func GetTemplate(ctxt TransformContextPtr, node libxml2.NodePtr, style StylesheetPtr) TemplatePtr

//go:linkname FreeTemplateHashes C.xsltFreeTemplateHashes
func FreeTemplateHashes(style StylesheetPtr)

//go:linkname CleanupTemplates C.xsltCleanupTemplates
func CleanupTemplates(style StylesheetPtr)
