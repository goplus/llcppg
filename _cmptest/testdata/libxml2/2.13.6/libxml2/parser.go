package libxml2

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

const DEFAULT_VERSION = "1.0"
const DETECT_IDS = 2
const COMPLETE_ATTRS = 4
const SKIP_IDS = 8
const SAX2_MAGIC = 0xDEEDBEAF

// llgo:type C
type ParserInputDeallocate func(*Char)

type X_xmlParserNodeInfo struct {
	Node      *X_xmlNode
	BeginPos  c.Ulong
	BeginLine c.Ulong
	EndPos    c.Ulong
	EndLine   c.Ulong
}
type ParserNodeInfo X_xmlParserNodeInfo
type ParserNodeInfoPtr *ParserNodeInfo

type X_xmlParserNodeInfoSeq struct {
	Maximum c.Ulong
	Length  c.Ulong
	Buffer  *ParserNodeInfo
}
type ParserNodeInfoSeq X_xmlParserNodeInfoSeq
type ParserNodeInfoSeqPtr *ParserNodeInfoSeq
type ParserInputState c.Int

const (
	PARSER_EOF             ParserInputState = -1
	PARSER_START           ParserInputState = 0
	PARSER_MISC            ParserInputState = 1
	PARSER_PI              ParserInputState = 2
	PARSER_DTD             ParserInputState = 3
	PARSER_PROLOG          ParserInputState = 4
	PARSER_COMMENT         ParserInputState = 5
	PARSER_START_TAG       ParserInputState = 6
	PARSER_CONTENT         ParserInputState = 7
	PARSER_CDATA_SECTION   ParserInputState = 8
	PARSER_END_TAG         ParserInputState = 9
	PARSER_ENTITY_DECL     ParserInputState = 10
	PARSER_ENTITY_VALUE    ParserInputState = 11
	PARSER_ATTRIBUTE_VALUE ParserInputState = 12
	PARSER_SYSTEM_LITERAL  ParserInputState = 13
	PARSER_EPILOG          ParserInputState = 14
	PARSER_IGNORE          ParserInputState = 15
	PARSER_PUBLIC_LITERAL  ParserInputState = 16
	PARSER_XML_DECL        ParserInputState = 17
)

type ParserMode c.Int

const (
	PARSE_UNKNOWN  ParserMode = 0
	PARSE_DOM      ParserMode = 1
	PARSE_SAX      ParserMode = 2
	PARSE_PUSH_DOM ParserMode = 3
	PARSE_PUSH_SAX ParserMode = 4
	PARSE_READER   ParserMode = 5
)

type X_xmlStartTag struct {
	Unused [8]uint8
}
type StartTag X_xmlStartTag

type X_xmlParserNsData struct {
	Unused [8]uint8
}
type ParserNsData X_xmlParserNsData

type X_xmlAttrHashBucket struct {
	Unused [8]uint8
}
type AttrHashBucket X_xmlAttrHashBucket

// llgo:type C
type ResolveEntitySAXFunc func(c.Pointer, *Char, *Char) ParserInputPtr

// llgo:type C
type InternalSubsetSAXFunc func(c.Pointer, *Char, *Char, *Char)

// llgo:type C
type ExternalSubsetSAXFunc func(c.Pointer, *Char, *Char, *Char)

// llgo:type C
type GetEntitySAXFunc func(c.Pointer, *Char) EntityPtr

// llgo:type C
type GetParameterEntitySAXFunc func(c.Pointer, *Char) EntityPtr

// llgo:type C
type EntityDeclSAXFunc func(c.Pointer, *Char, c.Int, *Char, *Char, *Char)

// llgo:type C
type NotationDeclSAXFunc func(c.Pointer, *Char, *Char, *Char)

// llgo:type C
type AttributeDeclSAXFunc func(c.Pointer, *Char, *Char, c.Int, c.Int, *Char, EnumerationPtr)

// llgo:type C
type ElementDeclSAXFunc func(c.Pointer, *Char, c.Int, ElementContentPtr)

// llgo:type C
type UnparsedEntityDeclSAXFunc func(c.Pointer, *Char, *Char, *Char, *Char)

// llgo:type C
type SetDocumentLocatorSAXFunc func(c.Pointer, SAXLocatorPtr)

// llgo:type C
type StartDocumentSAXFunc func(c.Pointer)

// llgo:type C
type EndDocumentSAXFunc func(c.Pointer)

// llgo:type C
type StartElementSAXFunc func(c.Pointer, *Char, **Char)

// llgo:type C
type EndElementSAXFunc func(c.Pointer, *Char)

// llgo:type C
type AttributeSAXFunc func(c.Pointer, *Char, *Char)

// llgo:type C
type ReferenceSAXFunc func(c.Pointer, *Char)

// llgo:type C
type CharactersSAXFunc func(c.Pointer, *Char, c.Int)

// llgo:type C
type IgnorableWhitespaceSAXFunc func(c.Pointer, *Char, c.Int)

// llgo:type C
type ProcessingInstructionSAXFunc func(c.Pointer, *Char, *Char)

// llgo:type C
type CommentSAXFunc func(c.Pointer, *Char)

// llgo:type C
type CdataBlockSAXFunc func(c.Pointer, *Char, c.Int)

// llgo:type C
type WarningSAXFunc func(__llgo_arg_0 c.Pointer, __llgo_arg_1 *c.Char, __llgo_va_list ...interface{})

// llgo:type C
type ErrorSAXFunc func(__llgo_arg_0 c.Pointer, __llgo_arg_1 *c.Char, __llgo_va_list ...interface{})

// llgo:type C
type FatalErrorSAXFunc func(__llgo_arg_0 c.Pointer, __llgo_arg_1 *c.Char, __llgo_va_list ...interface{})

// llgo:type C
type IsStandaloneSAXFunc func(c.Pointer) c.Int

// llgo:type C
type HasInternalSubsetSAXFunc func(c.Pointer) c.Int

// llgo:type C
type HasExternalSubsetSAXFunc func(c.Pointer) c.Int

// llgo:type C
type StartElementNsSAX2Func func(c.Pointer, *Char, *Char, *Char, c.Int, **Char, c.Int, c.Int, **Char)

// llgo:type C
type EndElementNsSAX2Func func(c.Pointer, *Char, *Char, *Char)

type X_xmlSAXHandlerV1 struct {
	InternalSubset        InternalSubsetSAXFunc
	IsStandalone          IsStandaloneSAXFunc
	HasInternalSubset     HasInternalSubsetSAXFunc
	HasExternalSubset     HasExternalSubsetSAXFunc
	ResolveEntity         ResolveEntitySAXFunc
	GetEntity             GetEntitySAXFunc
	EntityDecl            EntityDeclSAXFunc
	NotationDecl          NotationDeclSAXFunc
	AttributeDecl         AttributeDeclSAXFunc
	ElementDecl           ElementDeclSAXFunc
	UnparsedEntityDecl    UnparsedEntityDeclSAXFunc
	SetDocumentLocator    SetDocumentLocatorSAXFunc
	StartDocument         StartDocumentSAXFunc
	EndDocument           EndDocumentSAXFunc
	StartElement          StartElementSAXFunc
	EndElement            EndElementSAXFunc
	Reference             ReferenceSAXFunc
	Characters            CharactersSAXFunc
	IgnorableWhitespace   IgnorableWhitespaceSAXFunc
	ProcessingInstruction ProcessingInstructionSAXFunc
	Comment               CommentSAXFunc
	Warning               WarningSAXFunc
	Error                 ErrorSAXFunc
	FatalError            FatalErrorSAXFunc
	GetParameterEntity    GetParameterEntitySAXFunc
	CdataBlock            CdataBlockSAXFunc
	ExternalSubset        ExternalSubsetSAXFunc
	Initialized           c.Uint
}
type SAXHandlerV1 X_xmlSAXHandlerV1
type SAXHandlerV1Ptr *SAXHandlerV1

// llgo:type C
type ExternalEntityLoader func(*c.Char, *c.Char, ParserCtxtPtr) ParserInputPtr

/* backward compatibility */
//go:linkname X__xmlParserVersion C.__xmlParserVersion
func X__xmlParserVersion() **c.Char

//go:linkname X__oldXMLWDcompatibility C.__oldXMLWDcompatibility
func X__oldXMLWDcompatibility() *c.Int

//go:linkname X__xmlParserDebugEntities C.__xmlParserDebugEntities
func X__xmlParserDebugEntities() *c.Int

//go:linkname X__xmlDefaultSAXLocator C.__xmlDefaultSAXLocator
func X__xmlDefaultSAXLocator() *SAXLocator

//go:linkname X__xmlDefaultSAXHandler C.__xmlDefaultSAXHandler
func X__xmlDefaultSAXHandler() *SAXHandlerV1

//go:linkname X__xmlDoValidityCheckingDefaultValue C.__xmlDoValidityCheckingDefaultValue
func X__xmlDoValidityCheckingDefaultValue() *c.Int

//go:linkname X__xmlGetWarningsDefaultValue C.__xmlGetWarningsDefaultValue
func X__xmlGetWarningsDefaultValue() *c.Int

//go:linkname X__xmlKeepBlanksDefaultValue C.__xmlKeepBlanksDefaultValue
func X__xmlKeepBlanksDefaultValue() *c.Int

//go:linkname X__xmlLineNumbersDefaultValue C.__xmlLineNumbersDefaultValue
func X__xmlLineNumbersDefaultValue() *c.Int

//go:linkname X__xmlLoadExtDtdDefaultValue C.__xmlLoadExtDtdDefaultValue
func X__xmlLoadExtDtdDefaultValue() *c.Int

//go:linkname X__xmlPedanticParserDefaultValue C.__xmlPedanticParserDefaultValue
func X__xmlPedanticParserDefaultValue() *c.Int

//go:linkname X__xmlSubstituteEntitiesDefaultValue C.__xmlSubstituteEntitiesDefaultValue
func X__xmlSubstituteEntitiesDefaultValue() *c.Int

//go:linkname X__xmlIndentTreeOutput C.__xmlIndentTreeOutput
func X__xmlIndentTreeOutput() *c.Int

//go:linkname X__xmlTreeIndentString C.__xmlTreeIndentString
func X__xmlTreeIndentString() **c.Char

//go:linkname X__xmlSaveNoEmptyTags C.__xmlSaveNoEmptyTags
func X__xmlSaveNoEmptyTags() *c.Int

/*
 * Init/Cleanup
 */
//go:linkname InitParser C.xmlInitParser
func InitParser()

//go:linkname CleanupParser C.xmlCleanupParser
func CleanupParser()

//go:linkname InitGlobals C.xmlInitGlobals
func InitGlobals()

//go:linkname CleanupGlobals C.xmlCleanupGlobals
func CleanupGlobals()

/*
 * Input functions
 */
//go:linkname ParserInputRead C.xmlParserInputRead
func ParserInputRead(in ParserInputPtr, len c.Int) c.Int

//go:linkname ParserInputGrow C.xmlParserInputGrow
func ParserInputGrow(in ParserInputPtr, len c.Int) c.Int

/*
 * Basic parsing Interfaces
 */
// llgo:link (*Char).ParseDoc C.xmlParseDoc
func (recv_ *Char) ParseDoc() DocPtr {
	return nil
}

//go:linkname ParseFile C.xmlParseFile
func ParseFile(filename *c.Char) DocPtr

//go:linkname ParseMemory C.xmlParseMemory
func ParseMemory(buffer *c.Char, size c.Int) DocPtr

//go:linkname SubstituteEntitiesDefault C.xmlSubstituteEntitiesDefault
func SubstituteEntitiesDefault(val c.Int) c.Int

//go:linkname ThrDefSubstituteEntitiesDefaultValue C.xmlThrDefSubstituteEntitiesDefaultValue
func ThrDefSubstituteEntitiesDefaultValue(v c.Int) c.Int

//go:linkname KeepBlanksDefault C.xmlKeepBlanksDefault
func KeepBlanksDefault(val c.Int) c.Int

//go:linkname ThrDefKeepBlanksDefaultValue C.xmlThrDefKeepBlanksDefaultValue
func ThrDefKeepBlanksDefaultValue(v c.Int) c.Int

//go:linkname StopParser C.xmlStopParser
func StopParser(ctxt ParserCtxtPtr)

//go:linkname PedanticParserDefault C.xmlPedanticParserDefault
func PedanticParserDefault(val c.Int) c.Int

//go:linkname ThrDefPedanticParserDefaultValue C.xmlThrDefPedanticParserDefaultValue
func ThrDefPedanticParserDefaultValue(v c.Int) c.Int

//go:linkname LineNumbersDefault C.xmlLineNumbersDefault
func LineNumbersDefault(val c.Int) c.Int

//go:linkname ThrDefLineNumbersDefaultValue C.xmlThrDefLineNumbersDefaultValue
func ThrDefLineNumbersDefaultValue(v c.Int) c.Int

//go:linkname ThrDefDoValidityCheckingDefaultValue C.xmlThrDefDoValidityCheckingDefaultValue
func ThrDefDoValidityCheckingDefaultValue(v c.Int) c.Int

//go:linkname ThrDefGetWarningsDefaultValue C.xmlThrDefGetWarningsDefaultValue
func ThrDefGetWarningsDefaultValue(v c.Int) c.Int

//go:linkname ThrDefLoadExtDtdDefaultValue C.xmlThrDefLoadExtDtdDefaultValue
func ThrDefLoadExtDtdDefaultValue(v c.Int) c.Int

//go:linkname ThrDefParserDebugEntities C.xmlThrDefParserDebugEntities
func ThrDefParserDebugEntities(v c.Int) c.Int

/*
 * Recovery mode
 */
// llgo:link (*Char).RecoverDoc C.xmlRecoverDoc
func (recv_ *Char) RecoverDoc() DocPtr {
	return nil
}

//go:linkname RecoverMemory C.xmlRecoverMemory
func RecoverMemory(buffer *c.Char, size c.Int) DocPtr

//go:linkname RecoverFile C.xmlRecoverFile
func RecoverFile(filename *c.Char) DocPtr

/*
 * Less common routines and SAX interfaces
 */
//go:linkname ParseDocument C.xmlParseDocument
func ParseDocument(ctxt ParserCtxtPtr) c.Int

//go:linkname ParseExtParsedEnt C.xmlParseExtParsedEnt
func ParseExtParsedEnt(ctxt ParserCtxtPtr) c.Int

//go:linkname SAXUserParseFile C.xmlSAXUserParseFile
func SAXUserParseFile(sax SAXHandlerPtr, user_data c.Pointer, filename *c.Char) c.Int

//go:linkname SAXUserParseMemory C.xmlSAXUserParseMemory
func SAXUserParseMemory(sax SAXHandlerPtr, user_data c.Pointer, buffer *c.Char, size c.Int) c.Int

//go:linkname SAXParseDoc C.xmlSAXParseDoc
func SAXParseDoc(sax SAXHandlerPtr, cur *Char, recovery c.Int) DocPtr

//go:linkname SAXParseMemory C.xmlSAXParseMemory
func SAXParseMemory(sax SAXHandlerPtr, buffer *c.Char, size c.Int, recovery c.Int) DocPtr

//go:linkname SAXParseMemoryWithData C.xmlSAXParseMemoryWithData
func SAXParseMemoryWithData(sax SAXHandlerPtr, buffer *c.Char, size c.Int, recovery c.Int, data c.Pointer) DocPtr

//go:linkname SAXParseFile C.xmlSAXParseFile
func SAXParseFile(sax SAXHandlerPtr, filename *c.Char, recovery c.Int) DocPtr

//go:linkname SAXParseFileWithData C.xmlSAXParseFileWithData
func SAXParseFileWithData(sax SAXHandlerPtr, filename *c.Char, recovery c.Int, data c.Pointer) DocPtr

//go:linkname SAXParseEntity C.xmlSAXParseEntity
func SAXParseEntity(sax SAXHandlerPtr, filename *c.Char) DocPtr

//go:linkname ParseEntity C.xmlParseEntity
func ParseEntity(filename *c.Char) DocPtr

//go:linkname SAXParseDTD C.xmlSAXParseDTD
func SAXParseDTD(sax SAXHandlerPtr, ExternalID *Char, SystemID *Char) DtdPtr

// llgo:link (*Char).ParseDTD C.xmlParseDTD
func (recv_ *Char) ParseDTD(SystemID *Char) DtdPtr {
	return nil
}

//go:linkname IOParseDTD C.xmlIOParseDTD
func IOParseDTD(sax SAXHandlerPtr, input ParserInputBufferPtr, enc CharEncoding) DtdPtr

//go:linkname ParseBalancedChunkMemory C.xmlParseBalancedChunkMemory
func ParseBalancedChunkMemory(doc DocPtr, sax SAXHandlerPtr, user_data c.Pointer, depth c.Int, string *Char, lst *NodePtr) c.Int

//go:linkname ParseInNodeContext C.xmlParseInNodeContext
func ParseInNodeContext(node NodePtr, data *c.Char, datalen c.Int, options c.Int, lst *NodePtr) ParserErrors

//go:linkname ParseBalancedChunkMemoryRecover C.xmlParseBalancedChunkMemoryRecover
func ParseBalancedChunkMemoryRecover(doc DocPtr, sax SAXHandlerPtr, user_data c.Pointer, depth c.Int, string *Char, lst *NodePtr, recover c.Int) c.Int

//go:linkname ParseExternalEntity C.xmlParseExternalEntity
func ParseExternalEntity(doc DocPtr, sax SAXHandlerPtr, user_data c.Pointer, depth c.Int, URL *Char, ID *Char, lst *NodePtr) c.Int

//go:linkname ParseCtxtExternalEntity C.xmlParseCtxtExternalEntity
func ParseCtxtExternalEntity(ctx ParserCtxtPtr, URL *Char, ID *Char, lst *NodePtr) c.Int

/*
 * Parser contexts handling.
 */
//go:linkname NewParserCtxt C.xmlNewParserCtxt
func NewParserCtxt() ParserCtxtPtr

// llgo:link (*SAXHandler).NewSAXParserCtxt C.xmlNewSAXParserCtxt
func (recv_ *SAXHandler) NewSAXParserCtxt(userData c.Pointer) ParserCtxtPtr {
	return nil
}

//go:linkname InitParserCtxt C.xmlInitParserCtxt
func InitParserCtxt(ctxt ParserCtxtPtr) c.Int

//go:linkname ClearParserCtxt C.xmlClearParserCtxt
func ClearParserCtxt(ctxt ParserCtxtPtr)

//go:linkname FreeParserCtxt C.xmlFreeParserCtxt
func FreeParserCtxt(ctxt ParserCtxtPtr)

//go:linkname SetupParserForBuffer C.xmlSetupParserForBuffer
func SetupParserForBuffer(ctxt ParserCtxtPtr, buffer *Char, filename *c.Char)

// llgo:link (*Char).CreateDocParserCtxt C.xmlCreateDocParserCtxt
func (recv_ *Char) CreateDocParserCtxt() ParserCtxtPtr {
	return nil
}

/*
 * Reading/setting optional parsing features.
 */
//go:linkname GetFeaturesList C.xmlGetFeaturesList
func GetFeaturesList(len *c.Int, result **c.Char) c.Int

//go:linkname GetFeature C.xmlGetFeature
func GetFeature(ctxt ParserCtxtPtr, name *c.Char, result c.Pointer) c.Int

//go:linkname SetFeature C.xmlSetFeature
func SetFeature(ctxt ParserCtxtPtr, name *c.Char, value c.Pointer) c.Int

/*
 * Interfaces for the Push mode.
 */
//go:linkname CreatePushParserCtxt C.xmlCreatePushParserCtxt
func CreatePushParserCtxt(sax SAXHandlerPtr, user_data c.Pointer, chunk *c.Char, size c.Int, filename *c.Char) ParserCtxtPtr

//go:linkname ParseChunk C.xmlParseChunk
func ParseChunk(ctxt ParserCtxtPtr, chunk *c.Char, size c.Int, terminate c.Int) c.Int

/*
 * Special I/O mode.
 */
//go:linkname CreateIOParserCtxt C.xmlCreateIOParserCtxt
func CreateIOParserCtxt(sax SAXHandlerPtr, user_data c.Pointer, ioread InputReadCallback, ioclose InputCloseCallback, ioctx c.Pointer, enc CharEncoding) ParserCtxtPtr

//go:linkname NewIOInputStream C.xmlNewIOInputStream
func NewIOInputStream(ctxt ParserCtxtPtr, input ParserInputBufferPtr, enc CharEncoding) ParserInputPtr

/*
 * Node infos.
 */
//go:linkname ParserFindNodeInfo C.xmlParserFindNodeInfo
func ParserFindNodeInfo(ctxt ParserCtxtPtr, node NodePtr) *ParserNodeInfo

//go:linkname InitNodeInfoSeq C.xmlInitNodeInfoSeq
func InitNodeInfoSeq(seq ParserNodeInfoSeqPtr)

//go:linkname ClearNodeInfoSeq C.xmlClearNodeInfoSeq
func ClearNodeInfoSeq(seq ParserNodeInfoSeqPtr)

//go:linkname ParserFindNodeInfoIndex C.xmlParserFindNodeInfoIndex
func ParserFindNodeInfoIndex(seq ParserNodeInfoSeqPtr, node NodePtr) c.Ulong

//go:linkname ParserAddNodeInfo C.xmlParserAddNodeInfo
func ParserAddNodeInfo(ctxt ParserCtxtPtr, info ParserNodeInfoPtr)

/*
 * External entities handling actually implemented in xmlIO.
 */
//go:linkname SetExternalEntityLoader C.xmlSetExternalEntityLoader
func SetExternalEntityLoader(f ExternalEntityLoader)

//go:linkname GetExternalEntityLoader C.xmlGetExternalEntityLoader
func GetExternalEntityLoader() ExternalEntityLoader

//go:linkname LoadExternalEntity C.xmlLoadExternalEntity
func LoadExternalEntity(URL *c.Char, ID *c.Char, ctxt ParserCtxtPtr) ParserInputPtr

/*
 * Index lookup, actually implemented in the encoding module
 */
//go:linkname ByteConsumed C.xmlByteConsumed
func ByteConsumed(ctxt ParserCtxtPtr) c.Long

type ParserOption c.Int

const (
	PARSE_RECOVER    ParserOption = 1
	PARSE_NOENT      ParserOption = 2
	PARSE_DTDLOAD    ParserOption = 4
	PARSE_DTDATTR    ParserOption = 8
	PARSE_DTDVALID   ParserOption = 16
	PARSE_NOERROR    ParserOption = 32
	PARSE_NOWARNING  ParserOption = 64
	PARSE_PEDANTIC   ParserOption = 128
	PARSE_NOBLANKS   ParserOption = 256
	PARSE_SAX1       ParserOption = 512
	PARSE_XINCLUDE   ParserOption = 1024
	PARSE_NONET      ParserOption = 2048
	PARSE_NODICT     ParserOption = 4096
	PARSE_NSCLEAN    ParserOption = 8192
	PARSE_NOCDATA    ParserOption = 16384
	PARSE_NOXINCNODE ParserOption = 32768
	PARSE_COMPACT    ParserOption = 65536
	PARSE_OLD10      ParserOption = 131072
	PARSE_NOBASEFIX  ParserOption = 262144
	PARSE_HUGE       ParserOption = 524288
	PARSE_OLDSAX     ParserOption = 1048576
	PARSE_IGNORE_ENC ParserOption = 2097152
	PARSE_BIG_LINES  ParserOption = 4194304
	PARSE_NO_XXE     ParserOption = 8388608
)

//go:linkname CtxtReset C.xmlCtxtReset
func CtxtReset(ctxt ParserCtxtPtr)

//go:linkname CtxtResetPush C.xmlCtxtResetPush
func CtxtResetPush(ctxt ParserCtxtPtr, chunk *c.Char, size c.Int, filename *c.Char, encoding *c.Char) c.Int

//go:linkname CtxtSetOptions C.xmlCtxtSetOptions
func CtxtSetOptions(ctxt ParserCtxtPtr, options c.Int) c.Int

//go:linkname CtxtUseOptions C.xmlCtxtUseOptions
func CtxtUseOptions(ctxt ParserCtxtPtr, options c.Int) c.Int

//go:linkname CtxtSetErrorHandler C.xmlCtxtSetErrorHandler
func CtxtSetErrorHandler(ctxt ParserCtxtPtr, handler StructuredErrorFunc, data c.Pointer)

//go:linkname CtxtSetMaxAmplification C.xmlCtxtSetMaxAmplification
func CtxtSetMaxAmplification(ctxt ParserCtxtPtr, maxAmpl c.Uint)

// llgo:link (*Char).ReadDoc C.xmlReadDoc
func (recv_ *Char) ReadDoc(URL *c.Char, encoding *c.Char, options c.Int) DocPtr {
	return nil
}

//go:linkname ReadFile C.xmlReadFile
func ReadFile(URL *c.Char, encoding *c.Char, options c.Int) DocPtr

//go:linkname ReadMemory C.xmlReadMemory
func ReadMemory(buffer *c.Char, size c.Int, URL *c.Char, encoding *c.Char, options c.Int) DocPtr

//go:linkname ReadFd C.xmlReadFd
func ReadFd(fd c.Int, URL *c.Char, encoding *c.Char, options c.Int) DocPtr

//go:linkname ReadIO C.xmlReadIO
func ReadIO(ioread InputReadCallback, ioclose InputCloseCallback, ioctx c.Pointer, URL *c.Char, encoding *c.Char, options c.Int) DocPtr

//go:linkname CtxtParseDocument C.xmlCtxtParseDocument
func CtxtParseDocument(ctxt ParserCtxtPtr, input ParserInputPtr) DocPtr

//go:linkname CtxtReadDoc C.xmlCtxtReadDoc
func CtxtReadDoc(ctxt ParserCtxtPtr, cur *Char, URL *c.Char, encoding *c.Char, options c.Int) DocPtr

//go:linkname CtxtReadFile C.xmlCtxtReadFile
func CtxtReadFile(ctxt ParserCtxtPtr, filename *c.Char, encoding *c.Char, options c.Int) DocPtr

//go:linkname CtxtReadMemory C.xmlCtxtReadMemory
func CtxtReadMemory(ctxt ParserCtxtPtr, buffer *c.Char, size c.Int, URL *c.Char, encoding *c.Char, options c.Int) DocPtr

//go:linkname CtxtReadFd C.xmlCtxtReadFd
func CtxtReadFd(ctxt ParserCtxtPtr, fd c.Int, URL *c.Char, encoding *c.Char, options c.Int) DocPtr

//go:linkname CtxtReadIO C.xmlCtxtReadIO
func CtxtReadIO(ctxt ParserCtxtPtr, ioread InputReadCallback, ioclose InputCloseCallback, ioctx c.Pointer, URL *c.Char, encoding *c.Char, options c.Int) DocPtr

type Feature c.Int

const (
	WITH_THREAD     Feature = 1
	WITH_TREE       Feature = 2
	WITH_OUTPUT     Feature = 3
	WITH_PUSH       Feature = 4
	WITH_READER     Feature = 5
	WITH_PATTERN    Feature = 6
	WITH_WRITER     Feature = 7
	WITH_SAX1       Feature = 8
	WITH_FTP        Feature = 9
	WITH_HTTP       Feature = 10
	WITH_VALID      Feature = 11
	WITH_HTML       Feature = 12
	WITH_LEGACY     Feature = 13
	WITH_C14N       Feature = 14
	WITH_CATALOG    Feature = 15
	WITH_XPATH      Feature = 16
	WITH_XPTR       Feature = 17
	WITH_XINCLUDE   Feature = 18
	WITH_ICONV      Feature = 19
	WITH_ISO8859X   Feature = 20
	WITH_UNICODE    Feature = 21
	WITH_REGEXP     Feature = 22
	WITH_AUTOMATA   Feature = 23
	WITH_EXPR       Feature = 24
	WITH_SCHEMAS    Feature = 25
	WITH_SCHEMATRON Feature = 26
	WITH_MODULES    Feature = 27
	WITH_DEBUG      Feature = 28
	WITH_DEBUG_MEM  Feature = 29
	WITH_DEBUG_RUN  Feature = 30
	WITH_ZLIB       Feature = 31
	WITH_ICU        Feature = 32
	WITH_LZMA       Feature = 33
	WITH_NONE       Feature = 99999
)

// llgo:link Feature.HasFeature C.xmlHasFeature
func (recv_ Feature) HasFeature() c.Int {
	return 0
}
