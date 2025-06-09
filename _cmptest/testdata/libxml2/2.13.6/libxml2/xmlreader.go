package libxml2

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type ParserSeverities c.Int

const (
	PARSER_SEVERITY_VALIDITY_WARNING ParserSeverities = 1
	PARSER_SEVERITY_VALIDITY_ERROR   ParserSeverities = 2
	PARSER_SEVERITY_WARNING          ParserSeverities = 3
	PARSER_SEVERITY_ERROR            ParserSeverities = 4
)

type TextReaderMode c.Int

const (
	TEXTREADER_MODE_INITIAL     TextReaderMode = 0
	TEXTREADER_MODE_INTERACTIVE TextReaderMode = 1
	TEXTREADER_MODE_ERROR       TextReaderMode = 2
	TEXTREADER_MODE_EOF         TextReaderMode = 3
	TEXTREADER_MODE_CLOSED      TextReaderMode = 4
	TEXTREADER_MODE_READING     TextReaderMode = 5
)

type ParserProperties c.Int

const (
	PARSER_LOADDTD        ParserProperties = 1
	PARSER_DEFAULTATTRS   ParserProperties = 2
	PARSER_VALIDATE       ParserProperties = 3
	PARSER_SUBST_ENTITIES ParserProperties = 4
)

type ReaderTypes c.Int

const (
	READER_TYPE_NONE                   ReaderTypes = 0
	READER_TYPE_ELEMENT                ReaderTypes = 1
	READER_TYPE_ATTRIBUTE              ReaderTypes = 2
	READER_TYPE_TEXT                   ReaderTypes = 3
	READER_TYPE_CDATA                  ReaderTypes = 4
	READER_TYPE_ENTITY_REFERENCE       ReaderTypes = 5
	READER_TYPE_ENTITY                 ReaderTypes = 6
	READER_TYPE_PROCESSING_INSTRUCTION ReaderTypes = 7
	READER_TYPE_COMMENT                ReaderTypes = 8
	READER_TYPE_DOCUMENT               ReaderTypes = 9
	READER_TYPE_DOCUMENT_TYPE          ReaderTypes = 10
	READER_TYPE_DOCUMENT_FRAGMENT      ReaderTypes = 11
	READER_TYPE_NOTATION               ReaderTypes = 12
	READER_TYPE_WHITESPACE             ReaderTypes = 13
	READER_TYPE_SIGNIFICANT_WHITESPACE ReaderTypes = 14
	READER_TYPE_END_ELEMENT            ReaderTypes = 15
	READER_TYPE_END_ENTITY             ReaderTypes = 16
	READER_TYPE_XML_DECLARATION        ReaderTypes = 17
)

type X_xmlTextReader struct {
	Unused [8]uint8
}
type TextReader X_xmlTextReader
type TextReaderPtr *TextReader

/*
 * Constructors & Destructor
 */
//go:linkname NewTextReader C.xmlNewTextReader
func NewTextReader(input ParserInputBufferPtr, URI *c.Char) TextReaderPtr

//go:linkname NewTextReaderFilename C.xmlNewTextReaderFilename
func NewTextReaderFilename(URI *c.Char) TextReaderPtr

//go:linkname FreeTextReader C.xmlFreeTextReader
func FreeTextReader(reader TextReaderPtr)

//go:linkname TextReaderSetup C.xmlTextReaderSetup
func TextReaderSetup(reader TextReaderPtr, input ParserInputBufferPtr, URL *c.Char, encoding *c.Char, options c.Int) c.Int

//go:linkname TextReaderSetMaxAmplification C.xmlTextReaderSetMaxAmplification
func TextReaderSetMaxAmplification(reader TextReaderPtr, maxAmpl c.Uint)

//go:linkname TextReaderGetLastError C.xmlTextReaderGetLastError
func TextReaderGetLastError(reader TextReaderPtr) *Error

/*
 * Iterators
 */
//go:linkname TextReaderRead C.xmlTextReaderRead
func TextReaderRead(reader TextReaderPtr) c.Int

//go:linkname TextReaderReadInnerXml C.xmlTextReaderReadInnerXml
func TextReaderReadInnerXml(reader TextReaderPtr) *Char

//go:linkname TextReaderReadOuterXml C.xmlTextReaderReadOuterXml
func TextReaderReadOuterXml(reader TextReaderPtr) *Char

//go:linkname TextReaderReadString C.xmlTextReaderReadString
func TextReaderReadString(reader TextReaderPtr) *Char

//go:linkname TextReaderReadAttributeValue C.xmlTextReaderReadAttributeValue
func TextReaderReadAttributeValue(reader TextReaderPtr) c.Int

/*
 * Attributes of the node
 */
//go:linkname TextReaderAttributeCount C.xmlTextReaderAttributeCount
func TextReaderAttributeCount(reader TextReaderPtr) c.Int

//go:linkname TextReaderDepth C.xmlTextReaderDepth
func TextReaderDepth(reader TextReaderPtr) c.Int

//go:linkname TextReaderHasAttributes C.xmlTextReaderHasAttributes
func TextReaderHasAttributes(reader TextReaderPtr) c.Int

//go:linkname TextReaderHasValue C.xmlTextReaderHasValue
func TextReaderHasValue(reader TextReaderPtr) c.Int

//go:linkname TextReaderIsDefault C.xmlTextReaderIsDefault
func TextReaderIsDefault(reader TextReaderPtr) c.Int

//go:linkname TextReaderIsEmptyElement C.xmlTextReaderIsEmptyElement
func TextReaderIsEmptyElement(reader TextReaderPtr) c.Int

//go:linkname TextReaderNodeType C.xmlTextReaderNodeType
func TextReaderNodeType(reader TextReaderPtr) c.Int

//go:linkname TextReaderQuoteChar C.xmlTextReaderQuoteChar
func TextReaderQuoteChar(reader TextReaderPtr) c.Int

//go:linkname TextReaderReadState C.xmlTextReaderReadState
func TextReaderReadState(reader TextReaderPtr) c.Int

//go:linkname TextReaderIsNamespaceDecl C.xmlTextReaderIsNamespaceDecl
func TextReaderIsNamespaceDecl(reader TextReaderPtr) c.Int

//go:linkname TextReaderConstBaseUri C.xmlTextReaderConstBaseUri
func TextReaderConstBaseUri(reader TextReaderPtr) *Char

//go:linkname TextReaderConstLocalName C.xmlTextReaderConstLocalName
func TextReaderConstLocalName(reader TextReaderPtr) *Char

//go:linkname TextReaderConstName C.xmlTextReaderConstName
func TextReaderConstName(reader TextReaderPtr) *Char

//go:linkname TextReaderConstNamespaceUri C.xmlTextReaderConstNamespaceUri
func TextReaderConstNamespaceUri(reader TextReaderPtr) *Char

//go:linkname TextReaderConstPrefix C.xmlTextReaderConstPrefix
func TextReaderConstPrefix(reader TextReaderPtr) *Char

//go:linkname TextReaderConstXmlLang C.xmlTextReaderConstXmlLang
func TextReaderConstXmlLang(reader TextReaderPtr) *Char

//go:linkname TextReaderConstString C.xmlTextReaderConstString
func TextReaderConstString(reader TextReaderPtr, str *Char) *Char

//go:linkname TextReaderConstValue C.xmlTextReaderConstValue
func TextReaderConstValue(reader TextReaderPtr) *Char

/*
 * use the Const version of the routine for
 * better performance and simpler code
 */
//go:linkname TextReaderBaseUri C.xmlTextReaderBaseUri
func TextReaderBaseUri(reader TextReaderPtr) *Char

//go:linkname TextReaderLocalName C.xmlTextReaderLocalName
func TextReaderLocalName(reader TextReaderPtr) *Char

//go:linkname TextReaderName C.xmlTextReaderName
func TextReaderName(reader TextReaderPtr) *Char

//go:linkname TextReaderNamespaceUri C.xmlTextReaderNamespaceUri
func TextReaderNamespaceUri(reader TextReaderPtr) *Char

//go:linkname TextReaderPrefix C.xmlTextReaderPrefix
func TextReaderPrefix(reader TextReaderPtr) *Char

//go:linkname TextReaderXmlLang C.xmlTextReaderXmlLang
func TextReaderXmlLang(reader TextReaderPtr) *Char

//go:linkname TextReaderValue C.xmlTextReaderValue
func TextReaderValue(reader TextReaderPtr) *Char

/*
 * Methods of the XmlTextReader
 */
//go:linkname TextReaderClose C.xmlTextReaderClose
func TextReaderClose(reader TextReaderPtr) c.Int

//go:linkname TextReaderGetAttributeNo C.xmlTextReaderGetAttributeNo
func TextReaderGetAttributeNo(reader TextReaderPtr, no c.Int) *Char

//go:linkname TextReaderGetAttribute C.xmlTextReaderGetAttribute
func TextReaderGetAttribute(reader TextReaderPtr, name *Char) *Char

//go:linkname TextReaderGetAttributeNs C.xmlTextReaderGetAttributeNs
func TextReaderGetAttributeNs(reader TextReaderPtr, localName *Char, namespaceURI *Char) *Char

//go:linkname TextReaderGetRemainder C.xmlTextReaderGetRemainder
func TextReaderGetRemainder(reader TextReaderPtr) ParserInputBufferPtr

//go:linkname TextReaderLookupNamespace C.xmlTextReaderLookupNamespace
func TextReaderLookupNamespace(reader TextReaderPtr, prefix *Char) *Char

//go:linkname TextReaderMoveToAttributeNo C.xmlTextReaderMoveToAttributeNo
func TextReaderMoveToAttributeNo(reader TextReaderPtr, no c.Int) c.Int

//go:linkname TextReaderMoveToAttribute C.xmlTextReaderMoveToAttribute
func TextReaderMoveToAttribute(reader TextReaderPtr, name *Char) c.Int

//go:linkname TextReaderMoveToAttributeNs C.xmlTextReaderMoveToAttributeNs
func TextReaderMoveToAttributeNs(reader TextReaderPtr, localName *Char, namespaceURI *Char) c.Int

//go:linkname TextReaderMoveToFirstAttribute C.xmlTextReaderMoveToFirstAttribute
func TextReaderMoveToFirstAttribute(reader TextReaderPtr) c.Int

//go:linkname TextReaderMoveToNextAttribute C.xmlTextReaderMoveToNextAttribute
func TextReaderMoveToNextAttribute(reader TextReaderPtr) c.Int

//go:linkname TextReaderMoveToElement C.xmlTextReaderMoveToElement
func TextReaderMoveToElement(reader TextReaderPtr) c.Int

//go:linkname TextReaderNormalization C.xmlTextReaderNormalization
func TextReaderNormalization(reader TextReaderPtr) c.Int

//go:linkname TextReaderConstEncoding C.xmlTextReaderConstEncoding
func TextReaderConstEncoding(reader TextReaderPtr) *Char

/*
 * Extensions
 */
//go:linkname TextReaderSetParserProp C.xmlTextReaderSetParserProp
func TextReaderSetParserProp(reader TextReaderPtr, prop c.Int, value c.Int) c.Int

//go:linkname TextReaderGetParserProp C.xmlTextReaderGetParserProp
func TextReaderGetParserProp(reader TextReaderPtr, prop c.Int) c.Int

//go:linkname TextReaderCurrentNode C.xmlTextReaderCurrentNode
func TextReaderCurrentNode(reader TextReaderPtr) NodePtr

//go:linkname TextReaderGetParserLineNumber C.xmlTextReaderGetParserLineNumber
func TextReaderGetParserLineNumber(reader TextReaderPtr) c.Int

//go:linkname TextReaderGetParserColumnNumber C.xmlTextReaderGetParserColumnNumber
func TextReaderGetParserColumnNumber(reader TextReaderPtr) c.Int

//go:linkname TextReaderPreserve C.xmlTextReaderPreserve
func TextReaderPreserve(reader TextReaderPtr) NodePtr

//go:linkname TextReaderPreservePattern C.xmlTextReaderPreservePattern
func TextReaderPreservePattern(reader TextReaderPtr, pattern *Char, namespaces **Char) c.Int

//go:linkname TextReaderCurrentDoc C.xmlTextReaderCurrentDoc
func TextReaderCurrentDoc(reader TextReaderPtr) DocPtr

//go:linkname TextReaderExpand C.xmlTextReaderExpand
func TextReaderExpand(reader TextReaderPtr) NodePtr

//go:linkname TextReaderNext C.xmlTextReaderNext
func TextReaderNext(reader TextReaderPtr) c.Int

//go:linkname TextReaderNextSibling C.xmlTextReaderNextSibling
func TextReaderNextSibling(reader TextReaderPtr) c.Int

//go:linkname TextReaderIsValid C.xmlTextReaderIsValid
func TextReaderIsValid(reader TextReaderPtr) c.Int

//go:linkname TextReaderRelaxNGValidate C.xmlTextReaderRelaxNGValidate
func TextReaderRelaxNGValidate(reader TextReaderPtr, rng *c.Char) c.Int

//go:linkname TextReaderRelaxNGValidateCtxt C.xmlTextReaderRelaxNGValidateCtxt
func TextReaderRelaxNGValidateCtxt(reader TextReaderPtr, ctxt RelaxNGValidCtxtPtr, options c.Int) c.Int

//go:linkname TextReaderRelaxNGSetSchema C.xmlTextReaderRelaxNGSetSchema
func TextReaderRelaxNGSetSchema(reader TextReaderPtr, schema RelaxNGPtr) c.Int

//go:linkname TextReaderSchemaValidate C.xmlTextReaderSchemaValidate
func TextReaderSchemaValidate(reader TextReaderPtr, xsd *c.Char) c.Int

//go:linkname TextReaderSchemaValidateCtxt C.xmlTextReaderSchemaValidateCtxt
func TextReaderSchemaValidateCtxt(reader TextReaderPtr, ctxt SchemaValidCtxtPtr, options c.Int) c.Int

//go:linkname TextReaderSetSchema C.xmlTextReaderSetSchema
func TextReaderSetSchema(reader TextReaderPtr, schema SchemaPtr) c.Int

//go:linkname TextReaderConstXmlVersion C.xmlTextReaderConstXmlVersion
func TextReaderConstXmlVersion(reader TextReaderPtr) *Char

//go:linkname TextReaderStandalone C.xmlTextReaderStandalone
func TextReaderStandalone(reader TextReaderPtr) c.Int

/*
 * Index lookup
 */
//go:linkname TextReaderByteConsumed C.xmlTextReaderByteConsumed
func TextReaderByteConsumed(reader TextReaderPtr) c.Long

/*
 * New more complete APIs for simpler creation and reuse of readers
 */
//go:linkname ReaderWalker C.xmlReaderWalker
func ReaderWalker(doc DocPtr) TextReaderPtr

// llgo:link (*Char).ReaderForDoc C.xmlReaderForDoc
func (recv_ *Char) ReaderForDoc(URL *c.Char, encoding *c.Char, options c.Int) TextReaderPtr {
	return nil
}

//go:linkname ReaderForFile C.xmlReaderForFile
func ReaderForFile(filename *c.Char, encoding *c.Char, options c.Int) TextReaderPtr

//go:linkname ReaderForMemory C.xmlReaderForMemory
func ReaderForMemory(buffer *c.Char, size c.Int, URL *c.Char, encoding *c.Char, options c.Int) TextReaderPtr

//go:linkname ReaderForFd C.xmlReaderForFd
func ReaderForFd(fd c.Int, URL *c.Char, encoding *c.Char, options c.Int) TextReaderPtr

//go:linkname ReaderForIO C.xmlReaderForIO
func ReaderForIO(ioread InputReadCallback, ioclose InputCloseCallback, ioctx c.Pointer, URL *c.Char, encoding *c.Char, options c.Int) TextReaderPtr

//go:linkname ReaderNewWalker C.xmlReaderNewWalker
func ReaderNewWalker(reader TextReaderPtr, doc DocPtr) c.Int

//go:linkname ReaderNewDoc C.xmlReaderNewDoc
func ReaderNewDoc(reader TextReaderPtr, cur *Char, URL *c.Char, encoding *c.Char, options c.Int) c.Int

//go:linkname ReaderNewFile C.xmlReaderNewFile
func ReaderNewFile(reader TextReaderPtr, filename *c.Char, encoding *c.Char, options c.Int) c.Int

//go:linkname ReaderNewMemory C.xmlReaderNewMemory
func ReaderNewMemory(reader TextReaderPtr, buffer *c.Char, size c.Int, URL *c.Char, encoding *c.Char, options c.Int) c.Int

//go:linkname ReaderNewFd C.xmlReaderNewFd
func ReaderNewFd(reader TextReaderPtr, fd c.Int, URL *c.Char, encoding *c.Char, options c.Int) c.Int

//go:linkname ReaderNewIO C.xmlReaderNewIO
func ReaderNewIO(reader TextReaderPtr, ioread InputReadCallback, ioclose InputCloseCallback, ioctx c.Pointer, URL *c.Char, encoding *c.Char, options c.Int) c.Int

type TextReaderLocatorPtr c.Pointer

// llgo:type C
type TextReaderErrorFunc func(c.Pointer, *c.Char, ParserSeverities, TextReaderLocatorPtr)

//go:linkname TextReaderLocatorLineNumber C.xmlTextReaderLocatorLineNumber
func TextReaderLocatorLineNumber(locator TextReaderLocatorPtr) c.Int

//go:linkname TextReaderLocatorBaseURI C.xmlTextReaderLocatorBaseURI
func TextReaderLocatorBaseURI(locator TextReaderLocatorPtr) *Char

//go:linkname TextReaderSetErrorHandler C.xmlTextReaderSetErrorHandler
func TextReaderSetErrorHandler(reader TextReaderPtr, f TextReaderErrorFunc, arg c.Pointer)

//go:linkname TextReaderSetStructuredErrorHandler C.xmlTextReaderSetStructuredErrorHandler
func TextReaderSetStructuredErrorHandler(reader TextReaderPtr, f StructuredErrorFunc, arg c.Pointer)

//go:linkname TextReaderGetErrorHandler C.xmlTextReaderGetErrorHandler
func TextReaderGetErrorHandler(reader TextReaderPtr, f TextReaderErrorFunc, arg *c.Pointer)
