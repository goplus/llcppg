package libxslt

import (
	"github.com/goplus/lib/c"
	"github.com/goplus/llcppg/_cmptest/testdata/libxml2/2.13.6/libxml2"
	_ "unsafe"
)

const MAX_SORT = 15

type X_xsltRuntimeExtra struct {
	Info       c.Pointer
	Deallocate libxml2.FreeFunc
	Val        struct {
		Ptr c.Pointer
	}
}
type RuntimeExtra X_xsltRuntimeExtra
type RuntimeExtraPtr *RuntimeExtra

type X_xsltTemplate struct {
	Next           *X_xsltTemplate
	Style          *X_xsltStylesheet
	Match          *libxml2.Char
	Priority       c.Float
	Name           *libxml2.Char
	NameURI        *libxml2.Char
	Mode           *libxml2.Char
	ModeURI        *libxml2.Char
	Content        libxml2.NodePtr
	Elem           libxml2.NodePtr
	InheritedNsNr  c.Int
	InheritedNs    *libxml2.NsPtr
	NbCalls        c.Int
	Time           c.Ulong
	Params         c.Pointer
	TemplNr        c.Int
	TemplMax       c.Int
	TemplCalledTab *TemplatePtr
	TemplCountTab  *c.Int
	Position       c.Int
}
type Template X_xsltTemplate
type TemplatePtr *Template

type X_xsltStylesheet struct {
	Parent             *X_xsltStylesheet
	Next               *X_xsltStylesheet
	Imports            *X_xsltStylesheet
	DocList            DocumentPtr
	Doc                libxml2.DocPtr
	StripSpaces        libxml2.HashTablePtr
	StripAll           c.Int
	CdataSection       libxml2.HashTablePtr
	Variables          StackElemPtr
	Templates          TemplatePtr
	TemplatesHash      libxml2.HashTablePtr
	RootMatch          *X_xsltCompMatch
	KeyMatch           *X_xsltCompMatch
	ElemMatch          *X_xsltCompMatch
	AttrMatch          *X_xsltCompMatch
	ParentMatch        *X_xsltCompMatch
	TextMatch          *X_xsltCompMatch
	PiMatch            *X_xsltCompMatch
	CommentMatch       *X_xsltCompMatch
	NsAliases          libxml2.HashTablePtr
	AttributeSets      libxml2.HashTablePtr
	NsHash             libxml2.HashTablePtr
	NsDefs             c.Pointer
	Keys               c.Pointer
	Method             *libxml2.Char
	MethodURI          *libxml2.Char
	Version            *libxml2.Char
	Encoding           *libxml2.Char
	OmitXmlDeclaration c.Int
	DecimalFormat      DecimalFormatPtr
	Standalone         c.Int
	DoctypePublic      *libxml2.Char
	DoctypeSystem      *libxml2.Char
	Indent             c.Int
	MediaType          *libxml2.Char
	PreComps           ElemPreCompPtr
	Warnings           c.Int
	Errors             c.Int
	ExclPrefix         *libxml2.Char
	ExclPrefixTab      **libxml2.Char
	ExclPrefixNr       c.Int
	ExclPrefixMax      c.Int
	X_private          c.Pointer
	ExtInfos           libxml2.HashTablePtr
	ExtrasNr           c.Int
	Includes           DocumentPtr
	Dict               libxml2.DictPtr
	AttVTs             c.Pointer
	DefaultAlias       *libxml2.Char
	Nopreproc          c.Int
	Internalized       c.Int
	LiteralResult      c.Int
	Principal          StylesheetPtr
	ForwardsCompatible c.Int
	NamedTemplates     libxml2.HashTablePtr
	XpathCtxt          libxml2.XPathContextPtr
	OpLimit            c.Ulong
	OpCount            c.Ulong
}

type X_xsltDecimalFormat struct {
	Next             *X_xsltDecimalFormat
	Name             *libxml2.Char
	Digit            *libxml2.Char
	PatternSeparator *libxml2.Char
	MinusSign        *libxml2.Char
	Infinity         *libxml2.Char
	NoNumber         *libxml2.Char
	DecimalPoint     *libxml2.Char
	Grouping         *libxml2.Char
	Percent          *libxml2.Char
	Permille         *libxml2.Char
	ZeroDigit        *libxml2.Char
	NsUri            *libxml2.Char
}
type DecimalFormat X_xsltDecimalFormat
type DecimalFormatPtr *DecimalFormat

type X_xsltDocument struct {
	Next           *X_xsltDocument
	Main           c.Int
	Doc            libxml2.DocPtr
	Keys           c.Pointer
	Includes       *X_xsltDocument
	Preproc        c.Int
	NbKeysComputed c.Int
}
type Document X_xsltDocument
type DocumentPtr *Document

type X_xsltKeyDef struct {
	Next    *X_xsltKeyDef
	Inst    libxml2.NodePtr
	Name    *libxml2.Char
	NameURI *libxml2.Char
	Match   *libxml2.Char
	Use     *libxml2.Char
	Comp    libxml2.XPathCompExprPtr
	Usecomp libxml2.XPathCompExprPtr
	NsList  *libxml2.NsPtr
	NsNr    c.Int
}
type KeyDef X_xsltKeyDef
type KeyDefPtr *KeyDef

type X_xsltKeyTable struct {
	Next    *X_xsltKeyTable
	Name    *libxml2.Char
	NameURI *libxml2.Char
	Keys    libxml2.HashTablePtr
}
type KeyTable X_xsltKeyTable
type KeyTablePtr *KeyTable
type Stylesheet X_xsltStylesheet
type StylesheetPtr *Stylesheet

type X_xsltTransformContext struct {
	Style               StylesheetPtr
	Type                OutputType
	Templ               TemplatePtr
	TemplNr             c.Int
	TemplMax            c.Int
	TemplTab            *TemplatePtr
	Vars                StackElemPtr
	VarsNr              c.Int
	VarsMax             c.Int
	VarsTab             *StackElemPtr
	VarsBase            c.Int
	ExtFunctions        libxml2.HashTablePtr
	ExtElements         libxml2.HashTablePtr
	ExtInfos            libxml2.HashTablePtr
	Mode                *libxml2.Char
	ModeURI             *libxml2.Char
	DocList             DocumentPtr
	Document            DocumentPtr
	Node                libxml2.NodePtr
	NodeList            libxml2.NodeSetPtr
	Output              libxml2.DocPtr
	Insert              libxml2.NodePtr
	XpathCtxt           libxml2.XPathContextPtr
	State               TransformState
	GlobalVars          libxml2.HashTablePtr
	Inst                libxml2.NodePtr
	Xinclude            c.Int
	OutputFile          *c.Char
	Profile             c.Int
	Prof                c.Long
	ProfNr              c.Int
	ProfMax             c.Int
	ProfTab             *c.Long
	X_private           c.Pointer
	ExtrasNr            c.Int
	ExtrasMax           c.Int
	Extras              RuntimeExtraPtr
	StyleList           DocumentPtr
	Sec                 c.Pointer
	Error               libxml2.GenericErrorFunc
	Errctx              c.Pointer
	Sortfunc            SortFunc
	TmpRVT              libxml2.DocPtr
	PersistRVT          libxml2.DocPtr
	Ctxtflags           c.Int
	Lasttext            *libxml2.Char
	Lasttsize           c.Int
	Lasttuse            c.Int
	DebugStatus         c.Int
	TraceCode           *c.Ulong
	ParserOptions       c.Int
	Dict                libxml2.DictPtr
	TmpDoc              libxml2.DocPtr
	Internalized        c.Int
	NbKeys              c.Int
	HasTemplKeyPatterns c.Int
	CurrentTemplateRule TemplatePtr
	InitialContextNode  libxml2.NodePtr
	InitialContextDoc   libxml2.DocPtr
	Cache               TransformCachePtr
	ContextVariable     c.Pointer
	LocalRVT            libxml2.DocPtr
	LocalRVTBase        libxml2.DocPtr
	KeyInitLevel        c.Int
	Depth               c.Int
	MaxTemplateDepth    c.Int
	MaxTemplateVars     c.Int
	OpLimit             c.Ulong
	OpCount             c.Ulong
	SourceDocDirty      c.Int
	CurrentId           c.Ulong
	NewLocale           NewLocaleFunc
	FreeLocale          FreeLocaleFunc
	GenSortKey          GenSortKeyFunc
}
type TransformContext X_xsltTransformContext
type TransformContextPtr *TransformContext

type X_xsltElemPreComp struct {
	Next ElemPreCompPtr
	Type StyleType
	Func TransformFunction
	Inst libxml2.NodePtr
	Free ElemPreCompDeallocator
}
type ElemPreComp X_xsltElemPreComp
type ElemPreCompPtr *ElemPreComp

// llgo:type C
type TransformFunction func(TransformContextPtr, libxml2.NodePtr, libxml2.NodePtr, ElemPreCompPtr)

// llgo:type C
type SortFunc func(TransformContextPtr, *libxml2.NodePtr, c.Int)
type StyleType c.Int

const (
	FUNC_COPY           StyleType = 1
	FUNC_SORT           StyleType = 2
	FUNC_TEXT           StyleType = 3
	FUNC_ELEMENT        StyleType = 4
	FUNC_ATTRIBUTE      StyleType = 5
	FUNC_COMMENT        StyleType = 6
	FUNC_PI             StyleType = 7
	FUNC_COPYOF         StyleType = 8
	FUNC_VALUEOF        StyleType = 9
	FUNC_NUMBER         StyleType = 10
	FUNC_APPLYIMPORTS   StyleType = 11
	FUNC_CALLTEMPLATE   StyleType = 12
	FUNC_APPLYTEMPLATES StyleType = 13
	FUNC_CHOOSE         StyleType = 14
	FUNC_IF             StyleType = 15
	FUNC_FOREACH        StyleType = 16
	FUNC_DOCUMENT       StyleType = 17
	FUNC_WITHPARAM      StyleType = 18
	FUNC_PARAM          StyleType = 19
	FUNC_VARIABLE       StyleType = 20
	FUNC_WHEN           StyleType = 21
	FUNC_EXTENSION      StyleType = 22
)

// llgo:type C
type ElemPreCompDeallocator func(ElemPreCompPtr)

type X_xsltStylePreComp struct {
	Next        ElemPreCompPtr
	Type        StyleType
	Func        TransformFunction
	Inst        libxml2.NodePtr
	Stype       *libxml2.Char
	HasStype    c.Int
	Number      c.Int
	Order       *libxml2.Char
	HasOrder    c.Int
	Descending  c.Int
	Lang        *libxml2.Char
	HasLang     c.Int
	CaseOrder   *libxml2.Char
	LowerFirst  c.Int
	Use         *libxml2.Char
	HasUse      c.Int
	Noescape    c.Int
	Name        *libxml2.Char
	HasName     c.Int
	Ns          *libxml2.Char
	HasNs       c.Int
	Mode        *libxml2.Char
	ModeURI     *libxml2.Char
	Test        *libxml2.Char
	Templ       TemplatePtr
	Select      *libxml2.Char
	Ver11       c.Int
	Filename    *libxml2.Char
	HasFilename c.Int
	Numdata     NumberData
	Comp        libxml2.XPathCompExprPtr
	NsList      *libxml2.NsPtr
	NsNr        c.Int
}
type StylePreComp X_xsltStylePreComp
type StylePreCompPtr *StylePreComp

type X_xsltStackElem struct {
	Next     *X_xsltStackElem
	Comp     StylePreCompPtr
	Computed c.Int
	Name     *libxml2.Char
	NameURI  *libxml2.Char
	Select   *libxml2.Char
	Tree     libxml2.NodePtr
	Value    libxml2.XPathObjectPtr
	Fragment libxml2.DocPtr
	Level    c.Int
	Context  TransformContextPtr
	Flags    c.Int
}
type StackElem X_xsltStackElem
type StackElemPtr *StackElem

type X_xsltTransformCache struct {
	RVT          libxml2.DocPtr
	NbRVT        c.Int
	StackItems   StackElemPtr
	NbStackItems c.Int
}
type TransformCache X_xsltTransformCache
type TransformCachePtr *TransformCache
type OutputType c.Int

const (
	OUTPUT_XML  OutputType = 0
	OUTPUT_HTML OutputType = 1
	OUTPUT_TEXT OutputType = 2
)

// llgo:type C
type NewLocaleFunc func(*libxml2.Char, c.Int) c.Pointer

// llgo:type C
type FreeLocaleFunc func(c.Pointer)

// llgo:type C
type GenSortKeyFunc func(c.Pointer, *libxml2.Char) *libxml2.Char
type TransformState c.Int

const (
	STATE_OK      TransformState = 0
	STATE_ERROR   TransformState = 1
	STATE_STOPPED TransformState = 2
)

/*
 * Functions associated to the internal types
xsltDecimalFormatPtr	xsltDecimalFormatGetByName(xsltStylesheetPtr sheet,
						   xmlChar *name);
*/
//go:linkname NewStylesheet C.xsltNewStylesheet
func NewStylesheet() StylesheetPtr

//go:linkname ParseStylesheetFile C.xsltParseStylesheetFile
func ParseStylesheetFile(filename *libxml2.Char) StylesheetPtr

//go:linkname FreeStylesheet C.xsltFreeStylesheet
func FreeStylesheet(style StylesheetPtr)

//go:linkname IsBlank C.xsltIsBlank
func IsBlank(str *libxml2.Char) c.Int

//go:linkname FreeStackElemList C.xsltFreeStackElemList
func FreeStackElemList(elem StackElemPtr)

//go:linkname DecimalFormatGetByName C.xsltDecimalFormatGetByName
func DecimalFormatGetByName(style StylesheetPtr, name *libxml2.Char) DecimalFormatPtr

//go:linkname DecimalFormatGetByQName C.xsltDecimalFormatGetByQName
func DecimalFormatGetByQName(style StylesheetPtr, nsUri *libxml2.Char, name *libxml2.Char) DecimalFormatPtr

//go:linkname ParseStylesheetProcess C.xsltParseStylesheetProcess
func ParseStylesheetProcess(ret StylesheetPtr, doc libxml2.DocPtr) StylesheetPtr

//go:linkname ParseStylesheetOutput C.xsltParseStylesheetOutput
func ParseStylesheetOutput(style StylesheetPtr, cur libxml2.NodePtr)

//go:linkname ParseStylesheetDoc C.xsltParseStylesheetDoc
func ParseStylesheetDoc(doc libxml2.DocPtr) StylesheetPtr

//go:linkname ParseStylesheetImportedDoc C.xsltParseStylesheetImportedDoc
func ParseStylesheetImportedDoc(doc libxml2.DocPtr, style StylesheetPtr) StylesheetPtr

//go:linkname ParseStylesheetUser C.xsltParseStylesheetUser
func ParseStylesheetUser(style StylesheetPtr, doc libxml2.DocPtr) c.Int

//go:linkname LoadStylesheetPI C.xsltLoadStylesheetPI
func LoadStylesheetPI(doc libxml2.DocPtr) StylesheetPtr

//go:linkname NumberFormat C.xsltNumberFormat
func NumberFormat(ctxt TransformContextPtr, data NumberDataPtr, node libxml2.NodePtr)

//go:linkname FormatNumberConversion C.xsltFormatNumberConversion
func FormatNumberConversion(self DecimalFormatPtr, format *libxml2.Char, number c.Double, result **libxml2.Char) libxml2.XPathError

//go:linkname ParseTemplateContent C.xsltParseTemplateContent
func ParseTemplateContent(style StylesheetPtr, templ libxml2.NodePtr)

//go:linkname AllocateExtra C.xsltAllocateExtra
func AllocateExtra(style StylesheetPtr) c.Int

//go:linkname AllocateExtraCtxt C.xsltAllocateExtraCtxt
func AllocateExtraCtxt(ctxt TransformContextPtr) c.Int

/*
 * Extra functions for Result Value Trees
 */
//go:linkname CreateRVT C.xsltCreateRVT
func CreateRVT(ctxt TransformContextPtr) libxml2.DocPtr

//go:linkname RegisterTmpRVT C.xsltRegisterTmpRVT
func RegisterTmpRVT(ctxt TransformContextPtr, RVT libxml2.DocPtr) c.Int

//go:linkname RegisterLocalRVT C.xsltRegisterLocalRVT
func RegisterLocalRVT(ctxt TransformContextPtr, RVT libxml2.DocPtr) c.Int

//go:linkname RegisterPersistRVT C.xsltRegisterPersistRVT
func RegisterPersistRVT(ctxt TransformContextPtr, RVT libxml2.DocPtr) c.Int

//go:linkname ExtensionInstructionResultRegister C.xsltExtensionInstructionResultRegister
func ExtensionInstructionResultRegister(ctxt TransformContextPtr, obj libxml2.XPathObjectPtr) c.Int

//go:linkname ExtensionInstructionResultFinalize C.xsltExtensionInstructionResultFinalize
func ExtensionInstructionResultFinalize(ctxt TransformContextPtr) c.Int

//go:linkname FlagRVTs C.xsltFlagRVTs
func FlagRVTs(ctxt TransformContextPtr, obj libxml2.XPathObjectPtr, val c.Int) c.Int

//go:linkname FreeRVTs C.xsltFreeRVTs
func FreeRVTs(ctxt TransformContextPtr)

//go:linkname ReleaseRVT C.xsltReleaseRVT
func ReleaseRVT(ctxt TransformContextPtr, RVT libxml2.DocPtr)

/*
 * Extra functions for Attribute Value Templates
 */
//go:linkname CompileAttr C.xsltCompileAttr
func CompileAttr(style StylesheetPtr, attr libxml2.AttrPtr)

//go:linkname EvalAVT C.xsltEvalAVT
func EvalAVT(ctxt TransformContextPtr, avt c.Pointer, node libxml2.NodePtr) *libxml2.Char

//go:linkname FreeAVTList C.xsltFreeAVTList
func FreeAVTList(avt c.Pointer)

/*
 * Extra function for successful xsltCleanupGlobals / xsltInit sequence.
 */
//go:linkname Uninit C.xsltUninit
func Uninit()

/************************************************************************
 *									*
 *  Transformation-time functions for *internal* use only               *
 *									*
 ************************************************************************/
//go:linkname InitCtxtKey C.xsltInitCtxtKey
func InitCtxtKey(ctxt TransformContextPtr, doc DocumentPtr, keyd KeyDefPtr) c.Int

//go:linkname InitAllDocKeys C.xsltInitAllDocKeys
func InitAllDocKeys(ctxt TransformContextPtr) c.Int
