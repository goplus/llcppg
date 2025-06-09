package libxml2

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

//go:linkname SAX2GetPublicId C.xmlSAX2GetPublicId
func SAX2GetPublicId(ctx c.Pointer) *Char

//go:linkname SAX2GetSystemId C.xmlSAX2GetSystemId
func SAX2GetSystemId(ctx c.Pointer) *Char

//go:linkname SAX2SetDocumentLocator C.xmlSAX2SetDocumentLocator
func SAX2SetDocumentLocator(ctx c.Pointer, loc SAXLocatorPtr)

//go:linkname SAX2GetLineNumber C.xmlSAX2GetLineNumber
func SAX2GetLineNumber(ctx c.Pointer) c.Int

//go:linkname SAX2GetColumnNumber C.xmlSAX2GetColumnNumber
func SAX2GetColumnNumber(ctx c.Pointer) c.Int

//go:linkname SAX2IsStandalone C.xmlSAX2IsStandalone
func SAX2IsStandalone(ctx c.Pointer) c.Int

//go:linkname SAX2HasInternalSubset C.xmlSAX2HasInternalSubset
func SAX2HasInternalSubset(ctx c.Pointer) c.Int

//go:linkname SAX2HasExternalSubset C.xmlSAX2HasExternalSubset
func SAX2HasExternalSubset(ctx c.Pointer) c.Int

//go:linkname SAX2InternalSubset C.xmlSAX2InternalSubset
func SAX2InternalSubset(ctx c.Pointer, name *Char, ExternalID *Char, SystemID *Char)

//go:linkname SAX2ExternalSubset C.xmlSAX2ExternalSubset
func SAX2ExternalSubset(ctx c.Pointer, name *Char, ExternalID *Char, SystemID *Char)

//go:linkname SAX2GetEntity C.xmlSAX2GetEntity
func SAX2GetEntity(ctx c.Pointer, name *Char) EntityPtr

//go:linkname SAX2GetParameterEntity C.xmlSAX2GetParameterEntity
func SAX2GetParameterEntity(ctx c.Pointer, name *Char) EntityPtr

//go:linkname SAX2ResolveEntity C.xmlSAX2ResolveEntity
func SAX2ResolveEntity(ctx c.Pointer, publicId *Char, systemId *Char) ParserInputPtr

//go:linkname SAX2EntityDecl C.xmlSAX2EntityDecl
func SAX2EntityDecl(ctx c.Pointer, name *Char, type_ c.Int, publicId *Char, systemId *Char, content *Char)

//go:linkname SAX2AttributeDecl C.xmlSAX2AttributeDecl
func SAX2AttributeDecl(ctx c.Pointer, elem *Char, fullname *Char, type_ c.Int, def c.Int, defaultValue *Char, tree EnumerationPtr)

//go:linkname SAX2ElementDecl C.xmlSAX2ElementDecl
func SAX2ElementDecl(ctx c.Pointer, name *Char, type_ c.Int, content ElementContentPtr)

//go:linkname SAX2NotationDecl C.xmlSAX2NotationDecl
func SAX2NotationDecl(ctx c.Pointer, name *Char, publicId *Char, systemId *Char)

//go:linkname SAX2UnparsedEntityDecl C.xmlSAX2UnparsedEntityDecl
func SAX2UnparsedEntityDecl(ctx c.Pointer, name *Char, publicId *Char, systemId *Char, notationName *Char)

//go:linkname SAX2StartDocument C.xmlSAX2StartDocument
func SAX2StartDocument(ctx c.Pointer)

//go:linkname SAX2EndDocument C.xmlSAX2EndDocument
func SAX2EndDocument(ctx c.Pointer)

//go:linkname SAX2StartElement C.xmlSAX2StartElement
func SAX2StartElement(ctx c.Pointer, fullname *Char, atts **Char)

//go:linkname SAX2EndElement C.xmlSAX2EndElement
func SAX2EndElement(ctx c.Pointer, name *Char)

//go:linkname SAX2StartElementNs C.xmlSAX2StartElementNs
func SAX2StartElementNs(ctx c.Pointer, localname *Char, prefix *Char, URI *Char, nb_namespaces c.Int, namespaces **Char, nb_attributes c.Int, nb_defaulted c.Int, attributes **Char)

//go:linkname SAX2EndElementNs C.xmlSAX2EndElementNs
func SAX2EndElementNs(ctx c.Pointer, localname *Char, prefix *Char, URI *Char)

//go:linkname SAX2Reference C.xmlSAX2Reference
func SAX2Reference(ctx c.Pointer, name *Char)

//go:linkname SAX2Characters C.xmlSAX2Characters
func SAX2Characters(ctx c.Pointer, ch *Char, len c.Int)

//go:linkname SAX2IgnorableWhitespace C.xmlSAX2IgnorableWhitespace
func SAX2IgnorableWhitespace(ctx c.Pointer, ch *Char, len c.Int)

//go:linkname SAX2ProcessingInstruction C.xmlSAX2ProcessingInstruction
func SAX2ProcessingInstruction(ctx c.Pointer, target *Char, data *Char)

//go:linkname SAX2Comment C.xmlSAX2Comment
func SAX2Comment(ctx c.Pointer, value *Char)

//go:linkname SAX2CDataBlock C.xmlSAX2CDataBlock
func SAX2CDataBlock(ctx c.Pointer, value *Char, len c.Int)

//go:linkname SAXDefaultVersion C.xmlSAXDefaultVersion
func SAXDefaultVersion(version c.Int) c.Int

// llgo:link (*SAXHandler).SAXVersion C.xmlSAXVersion
func (recv_ *SAXHandler) SAXVersion(version c.Int) c.Int {
	return 0
}

// llgo:link (*SAXHandler).SAX2InitDefaultSAXHandler C.xmlSAX2InitDefaultSAXHandler
func (recv_ *SAXHandler) SAX2InitDefaultSAXHandler(warning c.Int) {
}

// llgo:link (*SAXHandler).SAX2InitHtmlDefaultSAXHandler C.xmlSAX2InitHtmlDefaultSAXHandler
func (recv_ *SAXHandler) SAX2InitHtmlDefaultSAXHandler() {
}

//go:linkname HtmlDefaultSAXHandlerInit C.htmlDefaultSAXHandlerInit
func HtmlDefaultSAXHandlerInit()

//go:linkname DefaultSAXHandlerInit C.xmlDefaultSAXHandlerInit
func DefaultSAXHandlerInit()
