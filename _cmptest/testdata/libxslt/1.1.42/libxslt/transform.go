package libxslt

import (
	"github.com/goplus/lib/c"
	"github.com/goplus/llcppg/_cmptest/testdata/libxml2/2.13.6/libxml2"
	_ "unsafe"
)

/**
 * XInclude default processing.
 */
//go:linkname SetXIncludeDefault C.xsltSetXIncludeDefault
func SetXIncludeDefault(xinclude c.Int)

//go:linkname GetXIncludeDefault C.xsltGetXIncludeDefault
func GetXIncludeDefault() c.Int

/**
 * Export context to users.
 */
//go:linkname NewTransformContext C.xsltNewTransformContext
func NewTransformContext(style StylesheetPtr, doc libxml2.DocPtr) TransformContextPtr

//go:linkname FreeTransformContext C.xsltFreeTransformContext
func FreeTransformContext(ctxt TransformContextPtr)

//go:linkname ApplyStylesheetUser C.xsltApplyStylesheetUser
func ApplyStylesheetUser(style StylesheetPtr, doc libxml2.DocPtr, params **c.Char, output *c.Char, profile *c.FILE, userCtxt TransformContextPtr) libxml2.DocPtr

//go:linkname ProcessOneNode C.xsltProcessOneNode
func ProcessOneNode(ctxt TransformContextPtr, node libxml2.NodePtr, params StackElemPtr)

/**
 * Private Interfaces.
 */
//go:linkname ApplyStripSpaces C.xsltApplyStripSpaces
func ApplyStripSpaces(ctxt TransformContextPtr, node libxml2.NodePtr)

//go:linkname ApplyStylesheet C.xsltApplyStylesheet
func ApplyStylesheet(style StylesheetPtr, doc libxml2.DocPtr, params **c.Char) libxml2.DocPtr

//go:linkname ProfileStylesheet C.xsltProfileStylesheet
func ProfileStylesheet(style StylesheetPtr, doc libxml2.DocPtr, params **c.Char, output *c.FILE) libxml2.DocPtr

//go:linkname RunStylesheet C.xsltRunStylesheet
func RunStylesheet(style StylesheetPtr, doc libxml2.DocPtr, params **c.Char, output *c.Char, SAX libxml2.SAXHandlerPtr, IObuf libxml2.OutputBufferPtr) c.Int

//go:linkname RunStylesheetUser C.xsltRunStylesheetUser
func RunStylesheetUser(style StylesheetPtr, doc libxml2.DocPtr, params **c.Char, output *c.Char, SAX libxml2.SAXHandlerPtr, IObuf libxml2.OutputBufferPtr, profile *c.FILE, userCtxt TransformContextPtr) c.Int

//go:linkname ApplyOneTemplate C.xsltApplyOneTemplate
func ApplyOneTemplate(ctxt TransformContextPtr, node libxml2.NodePtr, list libxml2.NodePtr, templ TemplatePtr, params StackElemPtr)

//go:linkname DocumentElem C.xsltDocumentElem
func DocumentElem(ctxt TransformContextPtr, node libxml2.NodePtr, inst libxml2.NodePtr, comp ElemPreCompPtr)

//go:linkname Sort C.xsltSort
func Sort(ctxt TransformContextPtr, node libxml2.NodePtr, inst libxml2.NodePtr, comp ElemPreCompPtr)

//go:linkname Copy C.xsltCopy
func Copy(ctxt TransformContextPtr, node libxml2.NodePtr, inst libxml2.NodePtr, comp ElemPreCompPtr)

//go:linkname Text C.xsltText
func Text(ctxt TransformContextPtr, node libxml2.NodePtr, inst libxml2.NodePtr, comp ElemPreCompPtr)

//go:linkname Element C.xsltElement
func Element(ctxt TransformContextPtr, node libxml2.NodePtr, inst libxml2.NodePtr, comp ElemPreCompPtr)

//go:linkname Comment C.xsltComment
func Comment(ctxt TransformContextPtr, node libxml2.NodePtr, inst libxml2.NodePtr, comp ElemPreCompPtr)

//go:linkname Attribute C.xsltAttribute
func Attribute(ctxt TransformContextPtr, node libxml2.NodePtr, inst libxml2.NodePtr, comp ElemPreCompPtr)

//go:linkname ProcessingInstruction C.xsltProcessingInstruction
func ProcessingInstruction(ctxt TransformContextPtr, node libxml2.NodePtr, inst libxml2.NodePtr, comp ElemPreCompPtr)

//go:linkname CopyOf C.xsltCopyOf
func CopyOf(ctxt TransformContextPtr, node libxml2.NodePtr, inst libxml2.NodePtr, comp ElemPreCompPtr)

//go:linkname ValueOf C.xsltValueOf
func ValueOf(ctxt TransformContextPtr, node libxml2.NodePtr, inst libxml2.NodePtr, comp ElemPreCompPtr)

//go:linkname Number C.xsltNumber
func Number(ctxt TransformContextPtr, node libxml2.NodePtr, inst libxml2.NodePtr, comp ElemPreCompPtr)

//go:linkname ApplyImports C.xsltApplyImports
func ApplyImports(ctxt TransformContextPtr, node libxml2.NodePtr, inst libxml2.NodePtr, comp ElemPreCompPtr)

//go:linkname CallTemplate C.xsltCallTemplate
func CallTemplate(ctxt TransformContextPtr, node libxml2.NodePtr, inst libxml2.NodePtr, comp ElemPreCompPtr)

//go:linkname ApplyTemplates C.xsltApplyTemplates
func ApplyTemplates(ctxt TransformContextPtr, node libxml2.NodePtr, inst libxml2.NodePtr, comp ElemPreCompPtr)

//go:linkname Choose C.xsltChoose
func Choose(ctxt TransformContextPtr, node libxml2.NodePtr, inst libxml2.NodePtr, comp ElemPreCompPtr)

//go:linkname If C.xsltIf
func If(ctxt TransformContextPtr, node libxml2.NodePtr, inst libxml2.NodePtr, comp ElemPreCompPtr)

//go:linkname ForEach C.xsltForEach
func ForEach(ctxt TransformContextPtr, node libxml2.NodePtr, inst libxml2.NodePtr, comp ElemPreCompPtr)

//go:linkname RegisterAllElement C.xsltRegisterAllElement
func RegisterAllElement(ctxt TransformContextPtr)

//go:linkname CopyTextString C.xsltCopyTextString
func CopyTextString(ctxt TransformContextPtr, target libxml2.NodePtr, string *libxml2.Char, noescape c.Int) libxml2.NodePtr

/* Following 2 functions needed for libexslt/functions.c */
//go:linkname LocalVariablePop C.xsltLocalVariablePop
func LocalVariablePop(ctxt TransformContextPtr, limitNr c.Int, level c.Int)

//go:linkname LocalVariablePush C.xsltLocalVariablePush
func LocalVariablePush(ctxt TransformContextPtr, variable StackElemPtr, level c.Int) c.Int
