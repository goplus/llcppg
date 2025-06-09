package libxml2

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type HtmlParserCtxt ParserCtxt
type HtmlParserCtxtPtr ParserCtxtPtr
type HtmlParserNodeInfo ParserNodeInfo
type HtmlSAXHandler SAXHandler
type HtmlSAXHandlerPtr SAXHandlerPtr
type HtmlParserInput ParserInput
type HtmlParserInputPtr ParserInputPtr
type HtmlDocPtr DocPtr
type HtmlNodePtr NodePtr

type X_htmlElemDesc struct {
	Name          *c.Char
	StartTag      c.Char
	EndTag        c.Char
	SaveEndTag    c.Char
	Empty         c.Char
	Depr          c.Char
	Dtd           c.Char
	Isinline      c.Char
	Desc          *c.Char
	Subelts       **c.Char
	Defaultsubelt *c.Char
	AttrsOpt      **c.Char
	AttrsDepr     **c.Char
	AttrsReq      **c.Char
}
type HtmlElemDesc X_htmlElemDesc
type HtmlElemDescPtr *HtmlElemDesc

type X_htmlEntityDesc struct {
	Value c.Uint
	Name  *c.Char
	Desc  *c.Char
}
type HtmlEntityDesc X_htmlEntityDesc
type HtmlEntityDescPtr *HtmlEntityDesc

//go:linkname X__htmlDefaultSAXHandler C.__htmlDefaultSAXHandler
func X__htmlDefaultSAXHandler() *SAXHandlerV1

/*
 * There is only few public functions.
 */
//go:linkname HtmlInitAutoClose C.htmlInitAutoClose
func HtmlInitAutoClose()

// llgo:link (*Char).HtmlTagLookup C.htmlTagLookup
func (recv_ *Char) HtmlTagLookup() *HtmlElemDesc {
	return nil
}

// llgo:link (*Char).HtmlEntityLookup C.htmlEntityLookup
func (recv_ *Char) HtmlEntityLookup() *HtmlEntityDesc {
	return nil
}

//go:linkname HtmlEntityValueLookup C.htmlEntityValueLookup
func HtmlEntityValueLookup(value c.Uint) *HtmlEntityDesc

//go:linkname HtmlIsAutoClosed C.htmlIsAutoClosed
func HtmlIsAutoClosed(doc HtmlDocPtr, elem HtmlNodePtr) c.Int

//go:linkname HtmlAutoCloseTag C.htmlAutoCloseTag
func HtmlAutoCloseTag(doc HtmlDocPtr, name *Char, elem HtmlNodePtr) c.Int

//go:linkname HtmlParseEntityRef C.htmlParseEntityRef
func HtmlParseEntityRef(ctxt HtmlParserCtxtPtr, str **Char) *HtmlEntityDesc

//go:linkname HtmlParseCharRef C.htmlParseCharRef
func HtmlParseCharRef(ctxt HtmlParserCtxtPtr) c.Int

//go:linkname HtmlParseElement C.htmlParseElement
func HtmlParseElement(ctxt HtmlParserCtxtPtr)

//go:linkname HtmlNewParserCtxt C.htmlNewParserCtxt
func HtmlNewParserCtxt() HtmlParserCtxtPtr

// llgo:link (*HtmlSAXHandler).HtmlNewSAXParserCtxt C.htmlNewSAXParserCtxt
func (recv_ *HtmlSAXHandler) HtmlNewSAXParserCtxt(userData c.Pointer) HtmlParserCtxtPtr {
	return nil
}

//go:linkname HtmlCreateMemoryParserCtxt C.htmlCreateMemoryParserCtxt
func HtmlCreateMemoryParserCtxt(buffer *c.Char, size c.Int) HtmlParserCtxtPtr

//go:linkname HtmlParseDocument C.htmlParseDocument
func HtmlParseDocument(ctxt HtmlParserCtxtPtr) c.Int

// llgo:link (*Char).HtmlSAXParseDoc C.htmlSAXParseDoc
func (recv_ *Char) HtmlSAXParseDoc(encoding *c.Char, sax HtmlSAXHandlerPtr, userData c.Pointer) HtmlDocPtr {
	return nil
}

// llgo:link (*Char).HtmlParseDoc C.htmlParseDoc
func (recv_ *Char) HtmlParseDoc(encoding *c.Char) HtmlDocPtr {
	return nil
}

//go:linkname HtmlCreateFileParserCtxt C.htmlCreateFileParserCtxt
func HtmlCreateFileParserCtxt(filename *c.Char, encoding *c.Char) HtmlParserCtxtPtr

//go:linkname HtmlSAXParseFile C.htmlSAXParseFile
func HtmlSAXParseFile(filename *c.Char, encoding *c.Char, sax HtmlSAXHandlerPtr, userData c.Pointer) HtmlDocPtr

//go:linkname HtmlParseFile C.htmlParseFile
func HtmlParseFile(filename *c.Char, encoding *c.Char) HtmlDocPtr

//go:linkname UTF8ToHtml C.UTF8ToHtml
func UTF8ToHtml(out *c.Char, outlen *c.Int, in *c.Char, inlen *c.Int) c.Int

//go:linkname HtmlEncodeEntities C.htmlEncodeEntities
func HtmlEncodeEntities(out *c.Char, outlen *c.Int, in *c.Char, inlen *c.Int, quoteChar c.Int) c.Int

// llgo:link (*Char).HtmlIsScriptAttribute C.htmlIsScriptAttribute
func (recv_ *Char) HtmlIsScriptAttribute() c.Int {
	return 0
}

//go:linkname HtmlHandleOmittedElem C.htmlHandleOmittedElem
func HtmlHandleOmittedElem(val c.Int) c.Int

/**
 * Interfaces for the Push mode.
 */
//go:linkname HtmlCreatePushParserCtxt C.htmlCreatePushParserCtxt
func HtmlCreatePushParserCtxt(sax HtmlSAXHandlerPtr, user_data c.Pointer, chunk *c.Char, size c.Int, filename *c.Char, enc CharEncoding) HtmlParserCtxtPtr

//go:linkname HtmlParseChunk C.htmlParseChunk
func HtmlParseChunk(ctxt HtmlParserCtxtPtr, chunk *c.Char, size c.Int, terminate c.Int) c.Int

//go:linkname HtmlFreeParserCtxt C.htmlFreeParserCtxt
func HtmlFreeParserCtxt(ctxt HtmlParserCtxtPtr)

type HtmlParserOption c.Int

const (
	HTML_PARSE_RECOVER    HtmlParserOption = 1
	HTML_PARSE_NODEFDTD   HtmlParserOption = 4
	HTML_PARSE_NOERROR    HtmlParserOption = 32
	HTML_PARSE_NOWARNING  HtmlParserOption = 64
	HTML_PARSE_PEDANTIC   HtmlParserOption = 128
	HTML_PARSE_NOBLANKS   HtmlParserOption = 256
	HTML_PARSE_NONET      HtmlParserOption = 2048
	HTML_PARSE_NOIMPLIED  HtmlParserOption = 8192
	HTML_PARSE_COMPACT    HtmlParserOption = 65536
	HTML_PARSE_IGNORE_ENC HtmlParserOption = 2097152
)

//go:linkname HtmlCtxtReset C.htmlCtxtReset
func HtmlCtxtReset(ctxt HtmlParserCtxtPtr)

//go:linkname HtmlCtxtUseOptions C.htmlCtxtUseOptions
func HtmlCtxtUseOptions(ctxt HtmlParserCtxtPtr, options c.Int) c.Int

// llgo:link (*Char).HtmlReadDoc C.htmlReadDoc
func (recv_ *Char) HtmlReadDoc(URL *c.Char, encoding *c.Char, options c.Int) HtmlDocPtr {
	return nil
}

//go:linkname HtmlReadFile C.htmlReadFile
func HtmlReadFile(URL *c.Char, encoding *c.Char, options c.Int) HtmlDocPtr

//go:linkname HtmlReadMemory C.htmlReadMemory
func HtmlReadMemory(buffer *c.Char, size c.Int, URL *c.Char, encoding *c.Char, options c.Int) HtmlDocPtr

//go:linkname HtmlReadFd C.htmlReadFd
func HtmlReadFd(fd c.Int, URL *c.Char, encoding *c.Char, options c.Int) HtmlDocPtr

//go:linkname HtmlReadIO C.htmlReadIO
func HtmlReadIO(ioread InputReadCallback, ioclose InputCloseCallback, ioctx c.Pointer, URL *c.Char, encoding *c.Char, options c.Int) HtmlDocPtr

//go:linkname HtmlCtxtParseDocument C.htmlCtxtParseDocument
func HtmlCtxtParseDocument(ctxt HtmlParserCtxtPtr, input ParserInputPtr) HtmlDocPtr

//go:linkname HtmlCtxtReadDoc C.htmlCtxtReadDoc
func HtmlCtxtReadDoc(ctxt ParserCtxtPtr, cur *Char, URL *c.Char, encoding *c.Char, options c.Int) HtmlDocPtr

//go:linkname HtmlCtxtReadFile C.htmlCtxtReadFile
func HtmlCtxtReadFile(ctxt ParserCtxtPtr, filename *c.Char, encoding *c.Char, options c.Int) HtmlDocPtr

//go:linkname HtmlCtxtReadMemory C.htmlCtxtReadMemory
func HtmlCtxtReadMemory(ctxt ParserCtxtPtr, buffer *c.Char, size c.Int, URL *c.Char, encoding *c.Char, options c.Int) HtmlDocPtr

//go:linkname HtmlCtxtReadFd C.htmlCtxtReadFd
func HtmlCtxtReadFd(ctxt ParserCtxtPtr, fd c.Int, URL *c.Char, encoding *c.Char, options c.Int) HtmlDocPtr

//go:linkname HtmlCtxtReadIO C.htmlCtxtReadIO
func HtmlCtxtReadIO(ctxt ParserCtxtPtr, ioread InputReadCallback, ioclose InputCloseCallback, ioctx c.Pointer, URL *c.Char, encoding *c.Char, options c.Int) HtmlDocPtr

type HtmlStatus c.Int

const (
	HTML_NA         HtmlStatus = 0
	HTML_INVALID    HtmlStatus = 1
	HTML_DEPRECATED HtmlStatus = 2
	HTML_VALID      HtmlStatus = 4
	HTML_REQUIRED   HtmlStatus = 12
)

/* Using htmlElemDesc rather than name here, to emphasise the fact
   that otherwise there's a lookup overhead
*/
// llgo:link (*HtmlElemDesc).HtmlAttrAllowed C.htmlAttrAllowed
func (recv_ *HtmlElemDesc) HtmlAttrAllowed(*Char, c.Int) HtmlStatus {
	return 0
}

// llgo:link (*HtmlElemDesc).HtmlElementAllowedHere C.htmlElementAllowedHere
func (recv_ *HtmlElemDesc) HtmlElementAllowedHere(*Char) c.Int {
	return 0
}

// llgo:link (*HtmlElemDesc).HtmlElementStatusHere C.htmlElementStatusHere
func (recv_ *HtmlElemDesc) HtmlElementStatusHere(*HtmlElemDesc) HtmlStatus {
	return 0
}

//go:linkname HtmlNodeStatus C.htmlNodeStatus
func HtmlNodeStatus(HtmlNodePtr, c.Int) HtmlStatus
