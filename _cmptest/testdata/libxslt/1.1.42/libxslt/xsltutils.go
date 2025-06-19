package libxslt

import (
	"github.com/goplus/lib/c"
	"github.com/goplus/llcppg/_cmptest/testdata/libxml2/2.13.6/libxml2"
	_ "unsafe"
)

/*
 * Our own version of namespaced attributes lookup.
 */
//go:linkname GetNsProp C.xsltGetNsProp
func GetNsProp(node libxml2.NodePtr, name *libxml2.Char, nameSpace *libxml2.Char) *libxml2.Char

//go:linkname GetCNsProp C.xsltGetCNsProp
func GetCNsProp(style StylesheetPtr, node libxml2.NodePtr, name *libxml2.Char, nameSpace *libxml2.Char) *libxml2.Char

//go:linkname GetUTF8Char C.xsltGetUTF8Char
func GetUTF8Char(utf *c.Char, len *c.Int) c.Int

type DebugTraceCodes c.Int

const (
	TRACE_ALL             DebugTraceCodes = -1
	TRACE_NONE            DebugTraceCodes = 0
	TRACE_COPY_TEXT       DebugTraceCodes = 1
	TRACE_PROCESS_NODE    DebugTraceCodes = 2
	TRACE_APPLY_TEMPLATE  DebugTraceCodes = 4
	TRACE_COPY            DebugTraceCodes = 8
	TRACE_COMMENT         DebugTraceCodes = 16
	TRACE_PI              DebugTraceCodes = 32
	TRACE_COPY_OF         DebugTraceCodes = 64
	TRACE_VALUE_OF        DebugTraceCodes = 128
	TRACE_CALL_TEMPLATE   DebugTraceCodes = 256
	TRACE_APPLY_TEMPLATES DebugTraceCodes = 512
	TRACE_CHOOSE          DebugTraceCodes = 1024
	TRACE_IF              DebugTraceCodes = 2048
	TRACE_FOR_EACH        DebugTraceCodes = 4096
	TRACE_STRIP_SPACES    DebugTraceCodes = 8192
	TRACE_TEMPLATES       DebugTraceCodes = 16384
	TRACE_KEYS            DebugTraceCodes = 32768
	TRACE_VARIABLES       DebugTraceCodes = 65536
)

// llgo:link DebugTraceCodes.DebugSetDefaultTrace C.xsltDebugSetDefaultTrace
func (recv_ DebugTraceCodes) DebugSetDefaultTrace() {
}

//go:linkname DebugGetDefaultTrace C.xsltDebugGetDefaultTrace
func DebugGetDefaultTrace() DebugTraceCodes

//go:linkname PrintErrorContext C.xsltPrintErrorContext
func PrintErrorContext(ctxt TransformContextPtr, style StylesheetPtr, node libxml2.NodePtr)

//go:linkname Message C.xsltMessage
func Message(ctxt TransformContextPtr, node libxml2.NodePtr, inst libxml2.NodePtr)

//go:linkname SetGenericErrorFunc C.xsltSetGenericErrorFunc
func SetGenericErrorFunc(ctx c.Pointer, handler libxml2.GenericErrorFunc)

//go:linkname SetGenericDebugFunc C.xsltSetGenericDebugFunc
func SetGenericDebugFunc(ctx c.Pointer, handler libxml2.GenericErrorFunc)

//go:linkname SetTransformErrorFunc C.xsltSetTransformErrorFunc
func SetTransformErrorFunc(ctxt TransformContextPtr, ctx c.Pointer, handler libxml2.GenericErrorFunc)

//go:linkname TransformError C.xsltTransformError
func TransformError(ctxt TransformContextPtr, style StylesheetPtr, node libxml2.NodePtr, msg *c.Char, __llgo_va_list ...interface{})

//go:linkname SetCtxtParseOptions C.xsltSetCtxtParseOptions
func SetCtxtParseOptions(ctxt TransformContextPtr, options c.Int) c.Int

/*
 * Sorting.
 */
//go:linkname DocumentSortFunction C.xsltDocumentSortFunction
func DocumentSortFunction(list libxml2.NodeSetPtr)

//go:linkname SetSortFunc C.xsltSetSortFunc
func SetSortFunc(handler SortFunc)

//go:linkname SetCtxtSortFunc C.xsltSetCtxtSortFunc
func SetCtxtSortFunc(ctxt TransformContextPtr, handler SortFunc)

//go:linkname DefaultSortFunction C.xsltDefaultSortFunction
func DefaultSortFunction(ctxt TransformContextPtr, sorts *libxml2.NodePtr, nbsorts c.Int)

//go:linkname DoSortFunction C.xsltDoSortFunction
func DoSortFunction(ctxt TransformContextPtr, sorts *libxml2.NodePtr, nbsorts c.Int)

//go:linkname ComputeSortResult C.xsltComputeSortResult
func ComputeSortResult(ctxt TransformContextPtr, sort libxml2.NodePtr) *libxml2.XPathObjectPtr

/*
 * QNames handling.
 */
//go:linkname SplitQName C.xsltSplitQName
func SplitQName(dict libxml2.DictPtr, name *libxml2.Char, prefix **libxml2.Char) *libxml2.Char

//go:linkname GetQNameURI C.xsltGetQNameURI
func GetQNameURI(node libxml2.NodePtr, name **libxml2.Char) *libxml2.Char

//go:linkname GetQNameURI2 C.xsltGetQNameURI2
func GetQNameURI2(style StylesheetPtr, node libxml2.NodePtr, name **libxml2.Char) *libxml2.Char

/*
 * Output, reuse libxml I/O buffers.
 */
//go:linkname SaveResultTo C.xsltSaveResultTo
func SaveResultTo(buf libxml2.OutputBufferPtr, result libxml2.DocPtr, style StylesheetPtr) c.Int

//go:linkname SaveResultToFilename C.xsltSaveResultToFilename
func SaveResultToFilename(URI *c.Char, result libxml2.DocPtr, style StylesheetPtr, compression c.Int) c.Int

//go:linkname SaveResultToFile C.xsltSaveResultToFile
func SaveResultToFile(file *c.FILE, result libxml2.DocPtr, style StylesheetPtr) c.Int

//go:linkname SaveResultToFd C.xsltSaveResultToFd
func SaveResultToFd(fd c.Int, result libxml2.DocPtr, style StylesheetPtr) c.Int

//go:linkname SaveResultToString C.xsltSaveResultToString
func SaveResultToString(doc_txt_ptr **libxml2.Char, doc_txt_len *c.Int, result libxml2.DocPtr, style StylesheetPtr) c.Int

/*
 * XPath interface
 */
//go:linkname XPathCompile C.xsltXPathCompile
func XPathCompile(style StylesheetPtr, str *libxml2.Char) libxml2.XPathCompExprPtr

//go:linkname XPathCompileFlags C.xsltXPathCompileFlags
func XPathCompileFlags(style StylesheetPtr, str *libxml2.Char, flags c.Int) libxml2.XPathCompExprPtr

type DebugStatusCodes c.Int

const (
	DEBUG_NONE        DebugStatusCodes = 0
	DEBUG_INIT        DebugStatusCodes = 1
	DEBUG_STEP        DebugStatusCodes = 2
	DEBUG_STEPOUT     DebugStatusCodes = 3
	DEBUG_NEXT        DebugStatusCodes = 4
	DEBUG_STOP        DebugStatusCodes = 5
	DEBUG_CONT        DebugStatusCodes = 6
	DEBUG_RUN         DebugStatusCodes = 7
	DEBUG_RUN_RESTART DebugStatusCodes = 8
	DEBUG_QUIT        DebugStatusCodes = 9
)

// llgo:type C
type HandleDebuggerCallback func(libxml2.NodePtr, libxml2.NodePtr, TemplatePtr, TransformContextPtr)

// llgo:type C
type AddCallCallback func(TemplatePtr, libxml2.NodePtr) c.Int

// llgo:type C
type DropCallCallback func()

//go:linkname GetDebuggerStatus C.xsltGetDebuggerStatus
func GetDebuggerStatus() c.Int
