package libxml2

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

//go:linkname GetPublicId C.getPublicId
func GetPublicId(ctx c.Pointer) *Char

//go:linkname GetSystemId C.getSystemId
func GetSystemId(ctx c.Pointer) *Char

//go:linkname SetDocumentLocator C.setDocumentLocator
func SetDocumentLocator(ctx c.Pointer, loc SAXLocatorPtr)

//go:linkname GetLineNumber C.getLineNumber
func GetLineNumber(ctx c.Pointer) c.Int

//go:linkname GetColumnNumber C.getColumnNumber
func GetColumnNumber(ctx c.Pointer) c.Int

//go:linkname IsStandalone C.isStandalone
func IsStandalone(ctx c.Pointer) c.Int

//go:linkname HasInternalSubset C.hasInternalSubset
func HasInternalSubset(ctx c.Pointer) c.Int

//go:linkname HasExternalSubset C.hasExternalSubset
func HasExternalSubset(ctx c.Pointer) c.Int

//go:linkname InternalSubset C.internalSubset
func InternalSubset(ctx c.Pointer, name *Char, ExternalID *Char, SystemID *Char)

//go:linkname ExternalSubset C.externalSubset
func ExternalSubset(ctx c.Pointer, name *Char, ExternalID *Char, SystemID *Char)

//go:linkname GetEntity C.getEntity
func GetEntity(ctx c.Pointer, name *Char) EntityPtr

//go:linkname GetParameterEntity__1 C.getParameterEntity
func GetParameterEntity__1(ctx c.Pointer, name *Char) EntityPtr

//go:linkname ResolveEntity C.resolveEntity
func ResolveEntity(ctx c.Pointer, publicId *Char, systemId *Char) ParserInputPtr

//go:linkname EntityDecl C.entityDecl
func EntityDecl(ctx c.Pointer, name *Char, type_ c.Int, publicId *Char, systemId *Char, content *Char)

//go:linkname AttributeDecl C.attributeDecl
func AttributeDecl(ctx c.Pointer, elem *Char, fullname *Char, type_ c.Int, def c.Int, defaultValue *Char, tree EnumerationPtr)

//go:linkname ElementDecl C.elementDecl
func ElementDecl(ctx c.Pointer, name *Char, type_ c.Int, content ElementContentPtr)

//go:linkname NotationDecl C.notationDecl
func NotationDecl(ctx c.Pointer, name *Char, publicId *Char, systemId *Char)

//go:linkname UnparsedEntityDecl C.unparsedEntityDecl
func UnparsedEntityDecl(ctx c.Pointer, name *Char, publicId *Char, systemId *Char, notationName *Char)

//go:linkname StartDocument C.startDocument
func StartDocument(ctx c.Pointer)

//go:linkname EndDocument C.endDocument
func EndDocument(ctx c.Pointer)

//go:linkname GetAttribute C.attribute
func GetAttribute(ctx c.Pointer, fullname *Char, value *Char)

//go:linkname StartElement C.startElement
func StartElement(ctx c.Pointer, fullname *Char, atts **Char)

//go:linkname EndElement C.endElement
func EndElement(ctx c.Pointer, name *Char)

//go:linkname Reference C.reference
func Reference(ctx c.Pointer, name *Char)

//go:linkname Characters C.characters
func Characters(ctx c.Pointer, ch *Char, len c.Int)

//go:linkname IgnorableWhitespace C.ignorableWhitespace
func IgnorableWhitespace(ctx c.Pointer, ch *Char, len c.Int)

//go:linkname ProcessingInstruction C.processingInstruction
func ProcessingInstruction(ctx c.Pointer, target *Char, data *Char)

//go:linkname GlobalNamespace C.globalNamespace
func GlobalNamespace(ctx c.Pointer, href *Char, prefix *Char)

//go:linkname SetNamespace C.setNamespace
func SetNamespace(ctx c.Pointer, name *Char)

//go:linkname GetNamespace C.getNamespace
func GetNamespace(ctx c.Pointer) NsPtr

//go:linkname CheckNamespace C.checkNamespace
func CheckNamespace(ctx c.Pointer, nameSpace *Char) c.Int

//go:linkname NamespaceDecl C.namespaceDecl
func NamespaceDecl(ctx c.Pointer, href *Char, prefix *Char)

//go:linkname Comment C.comment
func Comment(ctx c.Pointer, value *Char)

//go:linkname CdataBlock C.cdataBlock
func CdataBlock(ctx c.Pointer, value *Char, len c.Int)

// llgo:link (*SAXHandlerV1).InitxmlDefaultSAXHandler C.initxmlDefaultSAXHandler
func (recv_ *SAXHandlerV1) InitxmlDefaultSAXHandler(warning c.Int) {
}

// llgo:link (*SAXHandlerV1).InithtmlDefaultSAXHandler C.inithtmlDefaultSAXHandler
func (recv_ *SAXHandlerV1) InithtmlDefaultSAXHandler() {
}
