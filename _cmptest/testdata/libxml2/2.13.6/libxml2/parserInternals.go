package libxml2

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

const MAX_TEXT_LENGTH = 10000000
const MAX_HUGE_LENGTH = 1000000000
const MAX_NAME_LENGTH = 50000
const MAX_DICTIONARY_LIMIT = 10000000
const MAX_LOOKUP_LIMIT = 10000000
const MAX_NAMELEN = 100
const INPUT_CHUNK = 250
const SUBSTITUTE_NONE = 0
const SUBSTITUTE_REF = 1
const SUBSTITUTE_PEREF = 2
const SUBSTITUTE_BOTH = 3

/*
 * Function to finish the work of the macros where needed.
 */
//go:linkname IsLetter C.xmlIsLetter
func IsLetter(c c.Int) c.Int

/**
 * Parser context.
 */
//go:linkname CreateFileParserCtxt C.xmlCreateFileParserCtxt
func CreateFileParserCtxt(filename *c.Char) ParserCtxtPtr

//go:linkname CreateURLParserCtxt C.xmlCreateURLParserCtxt
func CreateURLParserCtxt(filename *c.Char, options c.Int) ParserCtxtPtr

//go:linkname CreateMemoryParserCtxt C.xmlCreateMemoryParserCtxt
func CreateMemoryParserCtxt(buffer *c.Char, size c.Int) ParserCtxtPtr

// llgo:link (*Char).CreateEntityParserCtxt C.xmlCreateEntityParserCtxt
func (recv_ *Char) CreateEntityParserCtxt(ID *Char, base *Char) ParserCtxtPtr {
	return nil
}

//go:linkname CtxtErrMemory C.xmlCtxtErrMemory
func CtxtErrMemory(ctxt ParserCtxtPtr)

//go:linkname SwitchEncoding C.xmlSwitchEncoding
func SwitchEncoding(ctxt ParserCtxtPtr, enc CharEncoding) c.Int

//go:linkname SwitchEncodingName C.xmlSwitchEncodingName
func SwitchEncodingName(ctxt ParserCtxtPtr, encoding *c.Char) c.Int

//go:linkname SwitchToEncoding C.xmlSwitchToEncoding
func SwitchToEncoding(ctxt ParserCtxtPtr, handler CharEncodingHandlerPtr) c.Int

//go:linkname SwitchInputEncoding C.xmlSwitchInputEncoding
func SwitchInputEncoding(ctxt ParserCtxtPtr, input ParserInputPtr, handler CharEncodingHandlerPtr) c.Int

/**
 * Input Streams.
 */
//go:linkname NewStringInputStream C.xmlNewStringInputStream
func NewStringInputStream(ctxt ParserCtxtPtr, buffer *Char) ParserInputPtr

//go:linkname NewEntityInputStream C.xmlNewEntityInputStream
func NewEntityInputStream(ctxt ParserCtxtPtr, entity EntityPtr) ParserInputPtr

//go:linkname PushInput C.xmlPushInput
func PushInput(ctxt ParserCtxtPtr, input ParserInputPtr) c.Int

//go:linkname PopInput C.xmlPopInput
func PopInput(ctxt ParserCtxtPtr) Char

//go:linkname FreeInputStream C.xmlFreeInputStream
func FreeInputStream(input ParserInputPtr)

//go:linkname NewInputFromFile C.xmlNewInputFromFile
func NewInputFromFile(ctxt ParserCtxtPtr, filename *c.Char) ParserInputPtr

//go:linkname NewInputStream C.xmlNewInputStream
func NewInputStream(ctxt ParserCtxtPtr) ParserInputPtr

/**
 * Namespaces.
 */
//go:linkname SplitQName C.xmlSplitQName
func SplitQName(ctxt ParserCtxtPtr, name *Char, prefix **Char) *Char

/**
 * Generic production rules.
 */
//go:linkname ParseName C.xmlParseName
func ParseName(ctxt ParserCtxtPtr) *Char

//go:linkname ParseNmtoken C.xmlParseNmtoken
func ParseNmtoken(ctxt ParserCtxtPtr) *Char

//go:linkname ParseEntityValue C.xmlParseEntityValue
func ParseEntityValue(ctxt ParserCtxtPtr, orig **Char) *Char

//go:linkname ParseAttValue C.xmlParseAttValue
func ParseAttValue(ctxt ParserCtxtPtr) *Char

//go:linkname ParseSystemLiteral C.xmlParseSystemLiteral
func ParseSystemLiteral(ctxt ParserCtxtPtr) *Char

//go:linkname ParsePubidLiteral C.xmlParsePubidLiteral
func ParsePubidLiteral(ctxt ParserCtxtPtr) *Char

//go:linkname ParseCharData C.xmlParseCharData
func ParseCharData(ctxt ParserCtxtPtr, cdata c.Int)

//go:linkname ParseExternalID C.xmlParseExternalID
func ParseExternalID(ctxt ParserCtxtPtr, publicID **Char, strict c.Int) *Char

//go:linkname ParseComment C.xmlParseComment
func ParseComment(ctxt ParserCtxtPtr)

//go:linkname ParsePITarget C.xmlParsePITarget
func ParsePITarget(ctxt ParserCtxtPtr) *Char

//go:linkname ParsePI C.xmlParsePI
func ParsePI(ctxt ParserCtxtPtr)

//go:linkname ParseNotationDecl C.xmlParseNotationDecl
func ParseNotationDecl(ctxt ParserCtxtPtr)

//go:linkname ParseEntityDecl C.xmlParseEntityDecl
func ParseEntityDecl(ctxt ParserCtxtPtr)

//go:linkname ParseDefaultDecl C.xmlParseDefaultDecl
func ParseDefaultDecl(ctxt ParserCtxtPtr, value **Char) c.Int

//go:linkname ParseNotationType C.xmlParseNotationType
func ParseNotationType(ctxt ParserCtxtPtr) EnumerationPtr

//go:linkname ParseEnumerationType C.xmlParseEnumerationType
func ParseEnumerationType(ctxt ParserCtxtPtr) EnumerationPtr

//go:linkname ParseEnumeratedType C.xmlParseEnumeratedType
func ParseEnumeratedType(ctxt ParserCtxtPtr, tree *EnumerationPtr) c.Int

//go:linkname ParseAttributeType C.xmlParseAttributeType
func ParseAttributeType(ctxt ParserCtxtPtr, tree *EnumerationPtr) c.Int

//go:linkname ParseAttributeListDecl C.xmlParseAttributeListDecl
func ParseAttributeListDecl(ctxt ParserCtxtPtr)

//go:linkname ParseElementMixedContentDecl C.xmlParseElementMixedContentDecl
func ParseElementMixedContentDecl(ctxt ParserCtxtPtr, inputchk c.Int) ElementContentPtr

//go:linkname ParseElementChildrenContentDecl C.xmlParseElementChildrenContentDecl
func ParseElementChildrenContentDecl(ctxt ParserCtxtPtr, inputchk c.Int) ElementContentPtr

//go:linkname ParseElementContentDecl C.xmlParseElementContentDecl
func ParseElementContentDecl(ctxt ParserCtxtPtr, name *Char, result *ElementContentPtr) c.Int

//go:linkname ParseElementDecl C.xmlParseElementDecl
func ParseElementDecl(ctxt ParserCtxtPtr) c.Int

//go:linkname ParseMarkupDecl C.xmlParseMarkupDecl
func ParseMarkupDecl(ctxt ParserCtxtPtr)

//go:linkname ParseCharRef C.xmlParseCharRef
func ParseCharRef(ctxt ParserCtxtPtr) c.Int

//go:linkname ParseEntityRef C.xmlParseEntityRef
func ParseEntityRef(ctxt ParserCtxtPtr) EntityPtr

//go:linkname ParseReference C.xmlParseReference
func ParseReference(ctxt ParserCtxtPtr)

//go:linkname ParsePEReference C.xmlParsePEReference
func ParsePEReference(ctxt ParserCtxtPtr)

//go:linkname ParseDocTypeDecl C.xmlParseDocTypeDecl
func ParseDocTypeDecl(ctxt ParserCtxtPtr)

//go:linkname ParseAttribute C.xmlParseAttribute
func ParseAttribute(ctxt ParserCtxtPtr, value **Char) *Char

//go:linkname ParseStartTag C.xmlParseStartTag
func ParseStartTag(ctxt ParserCtxtPtr) *Char

//go:linkname ParseEndTag C.xmlParseEndTag
func ParseEndTag(ctxt ParserCtxtPtr)

//go:linkname ParseCDSect C.xmlParseCDSect
func ParseCDSect(ctxt ParserCtxtPtr)

//go:linkname ParseContent C.xmlParseContent
func ParseContent(ctxt ParserCtxtPtr)

//go:linkname ParseElement C.xmlParseElement
func ParseElement(ctxt ParserCtxtPtr)

//go:linkname ParseVersionNum C.xmlParseVersionNum
func ParseVersionNum(ctxt ParserCtxtPtr) *Char

//go:linkname ParseVersionInfo C.xmlParseVersionInfo
func ParseVersionInfo(ctxt ParserCtxtPtr) *Char

//go:linkname ParseEncName C.xmlParseEncName
func ParseEncName(ctxt ParserCtxtPtr) *Char

//go:linkname ParseEncodingDecl C.xmlParseEncodingDecl
func ParseEncodingDecl(ctxt ParserCtxtPtr) *Char

//go:linkname ParseSDDecl C.xmlParseSDDecl
func ParseSDDecl(ctxt ParserCtxtPtr) c.Int

//go:linkname ParseXMLDecl C.xmlParseXMLDecl
func ParseXMLDecl(ctxt ParserCtxtPtr)

//go:linkname ParseTextDecl C.xmlParseTextDecl
func ParseTextDecl(ctxt ParserCtxtPtr)

//go:linkname ParseMisc C.xmlParseMisc
func ParseMisc(ctxt ParserCtxtPtr)

//go:linkname ParseExternalSubset C.xmlParseExternalSubset
func ParseExternalSubset(ctxt ParserCtxtPtr, ExternalID *Char, SystemID *Char)

//go:linkname StringDecodeEntities C.xmlStringDecodeEntities
func StringDecodeEntities(ctxt ParserCtxtPtr, str *Char, what c.Int, end Char, end2 Char, end3 Char) *Char

//go:linkname StringLenDecodeEntities C.xmlStringLenDecodeEntities
func StringLenDecodeEntities(ctxt ParserCtxtPtr, str *Char, len c.Int, what c.Int, end Char, end2 Char, end3 Char) *Char

/*
 * Generated by MACROS on top of parser.c c.f. PUSH_AND_POP.
 */
//go:linkname NodePush C.nodePush
func NodePush(ctxt ParserCtxtPtr, value NodePtr) c.Int

//go:linkname NodePop C.nodePop
func NodePop(ctxt ParserCtxtPtr) NodePtr

//go:linkname InputPush C.inputPush
func InputPush(ctxt ParserCtxtPtr, value ParserInputPtr) c.Int

//go:linkname InputPop C.inputPop
func InputPop(ctxt ParserCtxtPtr) ParserInputPtr

//go:linkname NamePop C.namePop
func NamePop(ctxt ParserCtxtPtr) *Char

//go:linkname NamePush C.namePush
func NamePush(ctxt ParserCtxtPtr, value *Char) c.Int

/*
 * other commodities shared between parser.c and parserInternals.
 */
//go:linkname SkipBlankChars C.xmlSkipBlankChars
func SkipBlankChars(ctxt ParserCtxtPtr) c.Int

//go:linkname StringCurrentChar C.xmlStringCurrentChar
func StringCurrentChar(ctxt ParserCtxtPtr, cur *Char, len *c.Int) c.Int

//go:linkname ParserHandlePEReference C.xmlParserHandlePEReference
func ParserHandlePEReference(ctxt ParserCtxtPtr)

// llgo:link (*Char).CheckLanguageID C.xmlCheckLanguageID
func (recv_ *Char) CheckLanguageID() c.Int {
	return 0
}

/*
 * Really core function shared with HTML parser.
 */
//go:linkname CurrentChar C.xmlCurrentChar
func CurrentChar(ctxt ParserCtxtPtr, len *c.Int) c.Int

// llgo:link (*Char).CopyCharMultiByte C.xmlCopyCharMultiByte
func (recv_ *Char) CopyCharMultiByte(val c.Int) c.Int {
	return 0
}

//go:linkname CopyChar C.xmlCopyChar
func CopyChar(len c.Int, out *Char, val c.Int) c.Int

//go:linkname NextChar C.xmlNextChar
func NextChar(ctxt ParserCtxtPtr)

//go:linkname ParserInputShrink C.xmlParserInputShrink
func ParserInputShrink(in ParserInputPtr)

// llgo:type C
type EntityReferenceFunc func(EntityPtr, NodePtr, NodePtr)

//go:linkname SetEntityReferenceFunc C.xmlSetEntityReferenceFunc
func SetEntityReferenceFunc(func_ EntityReferenceFunc)

//go:linkname ParseQuotedString C.xmlParseQuotedString
func ParseQuotedString(ctxt ParserCtxtPtr) *Char

//go:linkname ParseNamespace C.xmlParseNamespace
func ParseNamespace(ctxt ParserCtxtPtr)

//go:linkname NamespaceParseNSDef C.xmlNamespaceParseNSDef
func NamespaceParseNSDef(ctxt ParserCtxtPtr) *Char

//go:linkname ScanName C.xmlScanName
func ScanName(ctxt ParserCtxtPtr) *Char

//go:linkname NamespaceParseNCName C.xmlNamespaceParseNCName
func NamespaceParseNCName(ctxt ParserCtxtPtr) *Char

//go:linkname ParserHandleReference C.xmlParserHandleReference
func ParserHandleReference(ctxt ParserCtxtPtr)

//go:linkname NamespaceParseQName C.xmlNamespaceParseQName
func NamespaceParseQName(ctxt ParserCtxtPtr, prefix **Char) *Char

/**
 * Entities
 */
//go:linkname DecodeEntities C.xmlDecodeEntities
func DecodeEntities(ctxt ParserCtxtPtr, len c.Int, what c.Int, end Char, end2 Char, end3 Char) *Char

//go:linkname HandleEntity C.xmlHandleEntity
func HandleEntity(ctxt ParserCtxtPtr, entity EntityPtr)
