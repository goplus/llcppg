package libxml2

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type X_xmlTextWriter struct {
	Unused [8]uint8
}
type TextWriter X_xmlTextWriter
type TextWriterPtr *TextWriter

/*
 * Constructors & Destructor
 */
//go:linkname NewTextWriter C.xmlNewTextWriter
func NewTextWriter(out OutputBufferPtr) TextWriterPtr

//go:linkname NewTextWriterFilename C.xmlNewTextWriterFilename
func NewTextWriterFilename(uri *c.Char, compression c.Int) TextWriterPtr

//go:linkname NewTextWriterMemory C.xmlNewTextWriterMemory
func NewTextWriterMemory(buf BufferPtr, compression c.Int) TextWriterPtr

//go:linkname NewTextWriterPushParser C.xmlNewTextWriterPushParser
func NewTextWriterPushParser(ctxt ParserCtxtPtr, compression c.Int) TextWriterPtr

//go:linkname NewTextWriterDoc C.xmlNewTextWriterDoc
func NewTextWriterDoc(doc *DocPtr, compression c.Int) TextWriterPtr

//go:linkname NewTextWriterTree C.xmlNewTextWriterTree
func NewTextWriterTree(doc DocPtr, node NodePtr, compression c.Int) TextWriterPtr

//go:linkname FreeTextWriter C.xmlFreeTextWriter
func FreeTextWriter(writer TextWriterPtr)

/*
 * Document
 */
//go:linkname TextWriterStartDocument C.xmlTextWriterStartDocument
func TextWriterStartDocument(writer TextWriterPtr, version *c.Char, encoding *c.Char, standalone *c.Char) c.Int

//go:linkname TextWriterEndDocument C.xmlTextWriterEndDocument
func TextWriterEndDocument(writer TextWriterPtr) c.Int

/*
 * Comments
 */
//go:linkname TextWriterStartComment C.xmlTextWriterStartComment
func TextWriterStartComment(writer TextWriterPtr) c.Int

//go:linkname TextWriterEndComment C.xmlTextWriterEndComment
func TextWriterEndComment(writer TextWriterPtr) c.Int

//go:linkname TextWriterWriteFormatComment C.xmlTextWriterWriteFormatComment
func TextWriterWriteFormatComment(writer TextWriterPtr, format *c.Char, __llgo_va_list ...interface{}) c.Int

//go:linkname TextWriterWriteVFormatComment C.xmlTextWriterWriteVFormatComment
func TextWriterWriteVFormatComment(writer TextWriterPtr, format *c.Char, argptr c.VaList) c.Int

//go:linkname TextWriterWriteComment C.xmlTextWriterWriteComment
func TextWriterWriteComment(writer TextWriterPtr, content *Char) c.Int

/*
 * Elements
 */
//go:linkname TextWriterStartElement C.xmlTextWriterStartElement
func TextWriterStartElement(writer TextWriterPtr, name *Char) c.Int

//go:linkname TextWriterStartElementNS C.xmlTextWriterStartElementNS
func TextWriterStartElementNS(writer TextWriterPtr, prefix *Char, name *Char, namespaceURI *Char) c.Int

//go:linkname TextWriterEndElement C.xmlTextWriterEndElement
func TextWriterEndElement(writer TextWriterPtr) c.Int

//go:linkname TextWriterFullEndElement C.xmlTextWriterFullEndElement
func TextWriterFullEndElement(writer TextWriterPtr) c.Int

/*
 * Elements conveniency functions
 */
//go:linkname TextWriterWriteFormatElement C.xmlTextWriterWriteFormatElement
func TextWriterWriteFormatElement(writer TextWriterPtr, name *Char, format *c.Char, __llgo_va_list ...interface{}) c.Int

//go:linkname TextWriterWriteVFormatElement C.xmlTextWriterWriteVFormatElement
func TextWriterWriteVFormatElement(writer TextWriterPtr, name *Char, format *c.Char, argptr c.VaList) c.Int

//go:linkname TextWriterWriteElement C.xmlTextWriterWriteElement
func TextWriterWriteElement(writer TextWriterPtr, name *Char, content *Char) c.Int

//go:linkname TextWriterWriteFormatElementNS C.xmlTextWriterWriteFormatElementNS
func TextWriterWriteFormatElementNS(writer TextWriterPtr, prefix *Char, name *Char, namespaceURI *Char, format *c.Char, __llgo_va_list ...interface{}) c.Int

//go:linkname TextWriterWriteVFormatElementNS C.xmlTextWriterWriteVFormatElementNS
func TextWriterWriteVFormatElementNS(writer TextWriterPtr, prefix *Char, name *Char, namespaceURI *Char, format *c.Char, argptr c.VaList) c.Int

//go:linkname TextWriterWriteElementNS C.xmlTextWriterWriteElementNS
func TextWriterWriteElementNS(writer TextWriterPtr, prefix *Char, name *Char, namespaceURI *Char, content *Char) c.Int

/*
 * Text
 */
//go:linkname TextWriterWriteFormatRaw C.xmlTextWriterWriteFormatRaw
func TextWriterWriteFormatRaw(writer TextWriterPtr, format *c.Char, __llgo_va_list ...interface{}) c.Int

//go:linkname TextWriterWriteVFormatRaw C.xmlTextWriterWriteVFormatRaw
func TextWriterWriteVFormatRaw(writer TextWriterPtr, format *c.Char, argptr c.VaList) c.Int

//go:linkname TextWriterWriteRawLen C.xmlTextWriterWriteRawLen
func TextWriterWriteRawLen(writer TextWriterPtr, content *Char, len c.Int) c.Int

//go:linkname TextWriterWriteRaw C.xmlTextWriterWriteRaw
func TextWriterWriteRaw(writer TextWriterPtr, content *Char) c.Int

//go:linkname TextWriterWriteFormatString C.xmlTextWriterWriteFormatString
func TextWriterWriteFormatString(writer TextWriterPtr, format *c.Char, __llgo_va_list ...interface{}) c.Int

//go:linkname TextWriterWriteVFormatString C.xmlTextWriterWriteVFormatString
func TextWriterWriteVFormatString(writer TextWriterPtr, format *c.Char, argptr c.VaList) c.Int

//go:linkname TextWriterWriteString C.xmlTextWriterWriteString
func TextWriterWriteString(writer TextWriterPtr, content *Char) c.Int

//go:linkname TextWriterWriteBase64 C.xmlTextWriterWriteBase64
func TextWriterWriteBase64(writer TextWriterPtr, data *c.Char, start c.Int, len c.Int) c.Int

//go:linkname TextWriterWriteBinHex C.xmlTextWriterWriteBinHex
func TextWriterWriteBinHex(writer TextWriterPtr, data *c.Char, start c.Int, len c.Int) c.Int

/*
 * Attributes
 */
//go:linkname TextWriterStartAttribute C.xmlTextWriterStartAttribute
func TextWriterStartAttribute(writer TextWriterPtr, name *Char) c.Int

//go:linkname TextWriterStartAttributeNS C.xmlTextWriterStartAttributeNS
func TextWriterStartAttributeNS(writer TextWriterPtr, prefix *Char, name *Char, namespaceURI *Char) c.Int

//go:linkname TextWriterEndAttribute C.xmlTextWriterEndAttribute
func TextWriterEndAttribute(writer TextWriterPtr) c.Int

/*
 * Attributes conveniency functions
 */
//go:linkname TextWriterWriteFormatAttribute C.xmlTextWriterWriteFormatAttribute
func TextWriterWriteFormatAttribute(writer TextWriterPtr, name *Char, format *c.Char, __llgo_va_list ...interface{}) c.Int

//go:linkname TextWriterWriteVFormatAttribute C.xmlTextWriterWriteVFormatAttribute
func TextWriterWriteVFormatAttribute(writer TextWriterPtr, name *Char, format *c.Char, argptr c.VaList) c.Int

//go:linkname TextWriterWriteAttribute C.xmlTextWriterWriteAttribute
func TextWriterWriteAttribute(writer TextWriterPtr, name *Char, content *Char) c.Int

//go:linkname TextWriterWriteFormatAttributeNS C.xmlTextWriterWriteFormatAttributeNS
func TextWriterWriteFormatAttributeNS(writer TextWriterPtr, prefix *Char, name *Char, namespaceURI *Char, format *c.Char, __llgo_va_list ...interface{}) c.Int

//go:linkname TextWriterWriteVFormatAttributeNS C.xmlTextWriterWriteVFormatAttributeNS
func TextWriterWriteVFormatAttributeNS(writer TextWriterPtr, prefix *Char, name *Char, namespaceURI *Char, format *c.Char, argptr c.VaList) c.Int

//go:linkname TextWriterWriteAttributeNS C.xmlTextWriterWriteAttributeNS
func TextWriterWriteAttributeNS(writer TextWriterPtr, prefix *Char, name *Char, namespaceURI *Char, content *Char) c.Int

/*
 * PI's
 */
//go:linkname TextWriterStartPI C.xmlTextWriterStartPI
func TextWriterStartPI(writer TextWriterPtr, target *Char) c.Int

//go:linkname TextWriterEndPI C.xmlTextWriterEndPI
func TextWriterEndPI(writer TextWriterPtr) c.Int

/*
 * PI conveniency functions
 */
//go:linkname TextWriterWriteFormatPI C.xmlTextWriterWriteFormatPI
func TextWriterWriteFormatPI(writer TextWriterPtr, target *Char, format *c.Char, __llgo_va_list ...interface{}) c.Int

//go:linkname TextWriterWriteVFormatPI C.xmlTextWriterWriteVFormatPI
func TextWriterWriteVFormatPI(writer TextWriterPtr, target *Char, format *c.Char, argptr c.VaList) c.Int

//go:linkname TextWriterWritePI C.xmlTextWriterWritePI
func TextWriterWritePI(writer TextWriterPtr, target *Char, content *Char) c.Int

/*
 * CDATA
 */
//go:linkname TextWriterStartCDATA C.xmlTextWriterStartCDATA
func TextWriterStartCDATA(writer TextWriterPtr) c.Int

//go:linkname TextWriterEndCDATA C.xmlTextWriterEndCDATA
func TextWriterEndCDATA(writer TextWriterPtr) c.Int

/*
 * CDATA conveniency functions
 */
//go:linkname TextWriterWriteFormatCDATA C.xmlTextWriterWriteFormatCDATA
func TextWriterWriteFormatCDATA(writer TextWriterPtr, format *c.Char, __llgo_va_list ...interface{}) c.Int

//go:linkname TextWriterWriteVFormatCDATA C.xmlTextWriterWriteVFormatCDATA
func TextWriterWriteVFormatCDATA(writer TextWriterPtr, format *c.Char, argptr c.VaList) c.Int

//go:linkname TextWriterWriteCDATA C.xmlTextWriterWriteCDATA
func TextWriterWriteCDATA(writer TextWriterPtr, content *Char) c.Int

/*
 * DTD
 */
//go:linkname TextWriterStartDTD C.xmlTextWriterStartDTD
func TextWriterStartDTD(writer TextWriterPtr, name *Char, pubid *Char, sysid *Char) c.Int

//go:linkname TextWriterEndDTD C.xmlTextWriterEndDTD
func TextWriterEndDTD(writer TextWriterPtr) c.Int

/*
 * DTD conveniency functions
 */
//go:linkname TextWriterWriteFormatDTD C.xmlTextWriterWriteFormatDTD
func TextWriterWriteFormatDTD(writer TextWriterPtr, name *Char, pubid *Char, sysid *Char, format *c.Char, __llgo_va_list ...interface{}) c.Int

//go:linkname TextWriterWriteVFormatDTD C.xmlTextWriterWriteVFormatDTD
func TextWriterWriteVFormatDTD(writer TextWriterPtr, name *Char, pubid *Char, sysid *Char, format *c.Char, argptr c.VaList) c.Int

//go:linkname TextWriterWriteDTD C.xmlTextWriterWriteDTD
func TextWriterWriteDTD(writer TextWriterPtr, name *Char, pubid *Char, sysid *Char, subset *Char) c.Int

/*
 * DTD element definition
 */
//go:linkname TextWriterStartDTDElement C.xmlTextWriterStartDTDElement
func TextWriterStartDTDElement(writer TextWriterPtr, name *Char) c.Int

//go:linkname TextWriterEndDTDElement C.xmlTextWriterEndDTDElement
func TextWriterEndDTDElement(writer TextWriterPtr) c.Int

/*
 * DTD element definition conveniency functions
 */
//go:linkname TextWriterWriteFormatDTDElement C.xmlTextWriterWriteFormatDTDElement
func TextWriterWriteFormatDTDElement(writer TextWriterPtr, name *Char, format *c.Char, __llgo_va_list ...interface{}) c.Int

//go:linkname TextWriterWriteVFormatDTDElement C.xmlTextWriterWriteVFormatDTDElement
func TextWriterWriteVFormatDTDElement(writer TextWriterPtr, name *Char, format *c.Char, argptr c.VaList) c.Int

//go:linkname TextWriterWriteDTDElement C.xmlTextWriterWriteDTDElement
func TextWriterWriteDTDElement(writer TextWriterPtr, name *Char, content *Char) c.Int

/*
 * DTD attribute list definition
 */
//go:linkname TextWriterStartDTDAttlist C.xmlTextWriterStartDTDAttlist
func TextWriterStartDTDAttlist(writer TextWriterPtr, name *Char) c.Int

//go:linkname TextWriterEndDTDAttlist C.xmlTextWriterEndDTDAttlist
func TextWriterEndDTDAttlist(writer TextWriterPtr) c.Int

/*
 * DTD attribute list definition conveniency functions
 */
//go:linkname TextWriterWriteFormatDTDAttlist C.xmlTextWriterWriteFormatDTDAttlist
func TextWriterWriteFormatDTDAttlist(writer TextWriterPtr, name *Char, format *c.Char, __llgo_va_list ...interface{}) c.Int

//go:linkname TextWriterWriteVFormatDTDAttlist C.xmlTextWriterWriteVFormatDTDAttlist
func TextWriterWriteVFormatDTDAttlist(writer TextWriterPtr, name *Char, format *c.Char, argptr c.VaList) c.Int

//go:linkname TextWriterWriteDTDAttlist C.xmlTextWriterWriteDTDAttlist
func TextWriterWriteDTDAttlist(writer TextWriterPtr, name *Char, content *Char) c.Int

/*
 * DTD entity definition
 */
//go:linkname TextWriterStartDTDEntity C.xmlTextWriterStartDTDEntity
func TextWriterStartDTDEntity(writer TextWriterPtr, pe c.Int, name *Char) c.Int

//go:linkname TextWriterEndDTDEntity C.xmlTextWriterEndDTDEntity
func TextWriterEndDTDEntity(writer TextWriterPtr) c.Int

/*
 * DTD entity definition conveniency functions
 */
//go:linkname TextWriterWriteFormatDTDInternalEntity C.xmlTextWriterWriteFormatDTDInternalEntity
func TextWriterWriteFormatDTDInternalEntity(writer TextWriterPtr, pe c.Int, name *Char, format *c.Char, __llgo_va_list ...interface{}) c.Int

//go:linkname TextWriterWriteVFormatDTDInternalEntity C.xmlTextWriterWriteVFormatDTDInternalEntity
func TextWriterWriteVFormatDTDInternalEntity(writer TextWriterPtr, pe c.Int, name *Char, format *c.Char, argptr c.VaList) c.Int

//go:linkname TextWriterWriteDTDInternalEntity C.xmlTextWriterWriteDTDInternalEntity
func TextWriterWriteDTDInternalEntity(writer TextWriterPtr, pe c.Int, name *Char, content *Char) c.Int

//go:linkname TextWriterWriteDTDExternalEntity C.xmlTextWriterWriteDTDExternalEntity
func TextWriterWriteDTDExternalEntity(writer TextWriterPtr, pe c.Int, name *Char, pubid *Char, sysid *Char, ndataid *Char) c.Int

//go:linkname TextWriterWriteDTDExternalEntityContents C.xmlTextWriterWriteDTDExternalEntityContents
func TextWriterWriteDTDExternalEntityContents(writer TextWriterPtr, pubid *Char, sysid *Char, ndataid *Char) c.Int

//go:linkname TextWriterWriteDTDEntity C.xmlTextWriterWriteDTDEntity
func TextWriterWriteDTDEntity(writer TextWriterPtr, pe c.Int, name *Char, pubid *Char, sysid *Char, ndataid *Char, content *Char) c.Int

/*
 * DTD notation definition
 */
//go:linkname TextWriterWriteDTDNotation C.xmlTextWriterWriteDTDNotation
func TextWriterWriteDTDNotation(writer TextWriterPtr, name *Char, pubid *Char, sysid *Char) c.Int

/*
 * Indentation
 */
//go:linkname TextWriterSetIndent C.xmlTextWriterSetIndent
func TextWriterSetIndent(writer TextWriterPtr, indent c.Int) c.Int

//go:linkname TextWriterSetIndentString C.xmlTextWriterSetIndentString
func TextWriterSetIndentString(writer TextWriterPtr, str *Char) c.Int

//go:linkname TextWriterSetQuoteChar C.xmlTextWriterSetQuoteChar
func TextWriterSetQuoteChar(writer TextWriterPtr, quotechar Char) c.Int

/*
 * misc
 */
//go:linkname TextWriterFlush C.xmlTextWriterFlush
func TextWriterFlush(writer TextWriterPtr) c.Int

//go:linkname TextWriterClose C.xmlTextWriterClose
func TextWriterClose(writer TextWriterPtr) c.Int
