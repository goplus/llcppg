package libxslt

import (
	"github.com/goplus/lib/c"
	"github.com/goplus/llcppg/_cmptest/testdata/libxml2/2.13.6/libxml2"
	_ "unsafe"
)

//go:linkname EvalXPathPredicate C.xsltEvalXPathPredicate
func EvalXPathPredicate(ctxt TransformContextPtr, comp libxml2.XPathCompExprPtr, nsList *libxml2.NsPtr, nsNr c.Int) c.Int

//go:linkname EvalTemplateString C.xsltEvalTemplateString
func EvalTemplateString(ctxt TransformContextPtr, contextNode libxml2.NodePtr, inst libxml2.NodePtr) *libxml2.Char

//go:linkname EvalAttrValueTemplate C.xsltEvalAttrValueTemplate
func EvalAttrValueTemplate(ctxt TransformContextPtr, node libxml2.NodePtr, name *libxml2.Char, ns *libxml2.Char) *libxml2.Char

//go:linkname EvalStaticAttrValueTemplate C.xsltEvalStaticAttrValueTemplate
func EvalStaticAttrValueTemplate(style StylesheetPtr, node libxml2.NodePtr, name *libxml2.Char, ns *libxml2.Char, found *c.Int) *libxml2.Char

/* TODO: this is obviously broken ... the namespaces should be passed too ! */
//go:linkname EvalXPathString C.xsltEvalXPathString
func EvalXPathString(ctxt TransformContextPtr, comp libxml2.XPathCompExprPtr) *libxml2.Char

//go:linkname EvalXPathStringNs C.xsltEvalXPathStringNs
func EvalXPathStringNs(ctxt TransformContextPtr, comp libxml2.XPathCompExprPtr, nsNr c.Int, nsList *libxml2.NsPtr) *libxml2.Char

//go:linkname TemplateProcess C.xsltTemplateProcess
func TemplateProcess(ctxt TransformContextPtr, node libxml2.NodePtr) *libxml2.NodePtr

//go:linkname AttrListTemplateProcess C.xsltAttrListTemplateProcess
func AttrListTemplateProcess(ctxt TransformContextPtr, target libxml2.NodePtr, cur libxml2.AttrPtr) libxml2.AttrPtr

//go:linkname AttrTemplateProcess C.xsltAttrTemplateProcess
func AttrTemplateProcess(ctxt TransformContextPtr, target libxml2.NodePtr, attr libxml2.AttrPtr) libxml2.AttrPtr

//go:linkname AttrTemplateValueProcess C.xsltAttrTemplateValueProcess
func AttrTemplateValueProcess(ctxt TransformContextPtr, attr *libxml2.Char) *libxml2.Char

//go:linkname AttrTemplateValueProcessNode C.xsltAttrTemplateValueProcessNode
func AttrTemplateValueProcessNode(ctxt TransformContextPtr, str *libxml2.Char, node libxml2.NodePtr) *libxml2.Char
