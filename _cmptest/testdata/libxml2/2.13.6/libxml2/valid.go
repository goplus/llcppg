package libxml2

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type X_xmlValidState struct {
	Unused [8]uint8
}
type ValidState X_xmlValidState
type ValidStatePtr *ValidState

// llgo:type C
type ValidityErrorFunc func(__llgo_arg_0 c.Pointer, __llgo_arg_1 *c.Char, __llgo_va_list ...interface{})

// llgo:type C
type ValidityWarningFunc func(__llgo_arg_0 c.Pointer, __llgo_arg_1 *c.Char, __llgo_va_list ...interface{})

type X_xmlValidCtxt struct {
	UserData  c.Pointer
	Error     ValidityErrorFunc
	Warning   ValidityWarningFunc
	Node      NodePtr
	NodeNr    c.Int
	NodeMax   c.Int
	NodeTab   *NodePtr
	Flags     c.Uint
	Doc       DocPtr
	Valid     c.Int
	Vstate    *ValidState
	VstateNr  c.Int
	VstateMax c.Int
	VstateTab *ValidState
	Am        AutomataPtr
	State     AutomataStatePtr
}
type ValidCtxt X_xmlValidCtxt
type ValidCtxtPtr *ValidCtxt
type NotationTable X_xmlHashTable
type NotationTablePtr *NotationTable
type ElementTable X_xmlHashTable
type ElementTablePtr *ElementTable
type AttributeTable X_xmlHashTable
type AttributeTablePtr *AttributeTable
type IDTable X_xmlHashTable
type IDTablePtr *IDTable
type RefTable X_xmlHashTable
type RefTablePtr *RefTable

/* Notation */
//go:linkname AddNotationDecl C.xmlAddNotationDecl
func AddNotationDecl(ctxt ValidCtxtPtr, dtd DtdPtr, name *Char, PublicID *Char, SystemID *Char) NotationPtr

//go:linkname CopyNotationTable C.xmlCopyNotationTable
func CopyNotationTable(table NotationTablePtr) NotationTablePtr

//go:linkname FreeNotationTable C.xmlFreeNotationTable
func FreeNotationTable(table NotationTablePtr)

//go:linkname DumpNotationDecl C.xmlDumpNotationDecl
func DumpNotationDecl(buf BufferPtr, nota NotationPtr)

/* XML_DEPRECATED, still used in lxml */
//go:linkname DumpNotationTable C.xmlDumpNotationTable
func DumpNotationTable(buf BufferPtr, table NotationTablePtr)

/* Element Content */
/* the non Doc version are being deprecated */
// llgo:link (*Char).NewElementContent C.xmlNewElementContent
func (recv_ *Char) NewElementContent(type_ ElementContentType) ElementContentPtr {
	return nil
}

//go:linkname CopyElementContent C.xmlCopyElementContent
func CopyElementContent(content ElementContentPtr) ElementContentPtr

//go:linkname FreeElementContent C.xmlFreeElementContent
func FreeElementContent(cur ElementContentPtr)

/* the new versions with doc argument */
//go:linkname NewDocElementContent C.xmlNewDocElementContent
func NewDocElementContent(doc DocPtr, name *Char, type_ ElementContentType) ElementContentPtr

//go:linkname CopyDocElementContent C.xmlCopyDocElementContent
func CopyDocElementContent(doc DocPtr, content ElementContentPtr) ElementContentPtr

//go:linkname FreeDocElementContent C.xmlFreeDocElementContent
func FreeDocElementContent(doc DocPtr, cur ElementContentPtr)

//go:linkname SnprintfElementContent C.xmlSnprintfElementContent
func SnprintfElementContent(buf *c.Char, size c.Int, content ElementContentPtr, englob c.Int)

//go:linkname SprintfElementContent C.xmlSprintfElementContent
func SprintfElementContent(buf *c.Char, content ElementContentPtr, englob c.Int)

/* Element */
//go:linkname AddElementDecl C.xmlAddElementDecl
func AddElementDecl(ctxt ValidCtxtPtr, dtd DtdPtr, name *Char, type_ ElementTypeVal, content ElementContentPtr) ElementPtr

//go:linkname CopyElementTable C.xmlCopyElementTable
func CopyElementTable(table ElementTablePtr) ElementTablePtr

//go:linkname FreeElementTable C.xmlFreeElementTable
func FreeElementTable(table ElementTablePtr)

//go:linkname DumpElementTable C.xmlDumpElementTable
func DumpElementTable(buf BufferPtr, table ElementTablePtr)

//go:linkname DumpElementDecl C.xmlDumpElementDecl
func DumpElementDecl(buf BufferPtr, elem ElementPtr)

/* Enumeration */
// llgo:link (*Char).CreateEnumeration C.xmlCreateEnumeration
func (recv_ *Char) CreateEnumeration() EnumerationPtr {
	return nil
}

//go:linkname FreeEnumeration C.xmlFreeEnumeration
func FreeEnumeration(cur EnumerationPtr)

//go:linkname CopyEnumeration C.xmlCopyEnumeration
func CopyEnumeration(cur EnumerationPtr) EnumerationPtr

/* Attribute */
//go:linkname AddAttributeDecl C.xmlAddAttributeDecl
func AddAttributeDecl(ctxt ValidCtxtPtr, dtd DtdPtr, elem *Char, name *Char, ns *Char, type_ AttributeType, def AttributeDefault, defaultValue *Char, tree EnumerationPtr) AttributePtr

//go:linkname CopyAttributeTable C.xmlCopyAttributeTable
func CopyAttributeTable(table AttributeTablePtr) AttributeTablePtr

//go:linkname FreeAttributeTable C.xmlFreeAttributeTable
func FreeAttributeTable(table AttributeTablePtr)

//go:linkname DumpAttributeTable C.xmlDumpAttributeTable
func DumpAttributeTable(buf BufferPtr, table AttributeTablePtr)

//go:linkname DumpAttributeDecl C.xmlDumpAttributeDecl
func DumpAttributeDecl(buf BufferPtr, attr AttributePtr)

/* IDs */
//go:linkname AddIDSafe C.xmlAddIDSafe
func AddIDSafe(attr AttrPtr, value *Char) c.Int

//go:linkname AddID C.xmlAddID
func AddID(ctxt ValidCtxtPtr, doc DocPtr, value *Char, attr AttrPtr) IDPtr

//go:linkname FreeIDTable C.xmlFreeIDTable
func FreeIDTable(table IDTablePtr)

//go:linkname GetID C.xmlGetID
func GetID(doc DocPtr, ID *Char) AttrPtr

//go:linkname IsID C.xmlIsID
func IsID(doc DocPtr, elem NodePtr, attr AttrPtr) c.Int

//go:linkname RemoveID C.xmlRemoveID
func RemoveID(doc DocPtr, attr AttrPtr) c.Int

/* IDREFs */
//go:linkname AddRef C.xmlAddRef
func AddRef(ctxt ValidCtxtPtr, doc DocPtr, value *Char, attr AttrPtr) RefPtr

//go:linkname FreeRefTable C.xmlFreeRefTable
func FreeRefTable(table RefTablePtr)

//go:linkname IsRef C.xmlIsRef
func IsRef(doc DocPtr, elem NodePtr, attr AttrPtr) c.Int

//go:linkname RemoveRef C.xmlRemoveRef
func RemoveRef(doc DocPtr, attr AttrPtr) c.Int

//go:linkname GetRefs C.xmlGetRefs
func GetRefs(doc DocPtr, ID *Char) ListPtr

/* Allocate/Release Validation Contexts */
//go:linkname NewValidCtxt C.xmlNewValidCtxt
func NewValidCtxt() ValidCtxtPtr

//go:linkname FreeValidCtxt C.xmlFreeValidCtxt
func FreeValidCtxt(ValidCtxtPtr)

//go:linkname ValidateRoot C.xmlValidateRoot
func ValidateRoot(ctxt ValidCtxtPtr, doc DocPtr) c.Int

//go:linkname ValidateElementDecl C.xmlValidateElementDecl
func ValidateElementDecl(ctxt ValidCtxtPtr, doc DocPtr, elem ElementPtr) c.Int

//go:linkname ValidNormalizeAttributeValue C.xmlValidNormalizeAttributeValue
func ValidNormalizeAttributeValue(doc DocPtr, elem NodePtr, name *Char, value *Char) *Char

//go:linkname ValidCtxtNormalizeAttributeValue C.xmlValidCtxtNormalizeAttributeValue
func ValidCtxtNormalizeAttributeValue(ctxt ValidCtxtPtr, doc DocPtr, elem NodePtr, name *Char, value *Char) *Char

//go:linkname ValidateAttributeDecl C.xmlValidateAttributeDecl
func ValidateAttributeDecl(ctxt ValidCtxtPtr, doc DocPtr, attr AttributePtr) c.Int

// llgo:link AttributeType.ValidateAttributeValue C.xmlValidateAttributeValue
func (recv_ AttributeType) ValidateAttributeValue(value *Char) c.Int {
	return 0
}

//go:linkname ValidateNotationDecl C.xmlValidateNotationDecl
func ValidateNotationDecl(ctxt ValidCtxtPtr, doc DocPtr, nota NotationPtr) c.Int

//go:linkname ValidateDtd C.xmlValidateDtd
func ValidateDtd(ctxt ValidCtxtPtr, doc DocPtr, dtd DtdPtr) c.Int

//go:linkname ValidateDtdFinal C.xmlValidateDtdFinal
func ValidateDtdFinal(ctxt ValidCtxtPtr, doc DocPtr) c.Int

//go:linkname ValidateDocument C.xmlValidateDocument
func ValidateDocument(ctxt ValidCtxtPtr, doc DocPtr) c.Int

//go:linkname ValidateElement C.xmlValidateElement
func ValidateElement(ctxt ValidCtxtPtr, doc DocPtr, elem NodePtr) c.Int

//go:linkname ValidateOneElement C.xmlValidateOneElement
func ValidateOneElement(ctxt ValidCtxtPtr, doc DocPtr, elem NodePtr) c.Int

//go:linkname ValidateOneAttribute C.xmlValidateOneAttribute
func ValidateOneAttribute(ctxt ValidCtxtPtr, doc DocPtr, elem NodePtr, attr AttrPtr, value *Char) c.Int

//go:linkname ValidateOneNamespace C.xmlValidateOneNamespace
func ValidateOneNamespace(ctxt ValidCtxtPtr, doc DocPtr, elem NodePtr, prefix *Char, ns NsPtr, value *Char) c.Int

//go:linkname ValidateDocumentFinal C.xmlValidateDocumentFinal
func ValidateDocumentFinal(ctxt ValidCtxtPtr, doc DocPtr) c.Int

//go:linkname ValidateNotationUse C.xmlValidateNotationUse
func ValidateNotationUse(ctxt ValidCtxtPtr, doc DocPtr, notationName *Char) c.Int

//go:linkname IsMixedElement C.xmlIsMixedElement
func IsMixedElement(doc DocPtr, name *Char) c.Int

//go:linkname GetDtdAttrDesc C.xmlGetDtdAttrDesc
func GetDtdAttrDesc(dtd DtdPtr, elem *Char, name *Char) AttributePtr

//go:linkname GetDtdQAttrDesc C.xmlGetDtdQAttrDesc
func GetDtdQAttrDesc(dtd DtdPtr, elem *Char, name *Char, prefix *Char) AttributePtr

//go:linkname GetDtdNotationDesc C.xmlGetDtdNotationDesc
func GetDtdNotationDesc(dtd DtdPtr, name *Char) NotationPtr

//go:linkname GetDtdQElementDesc C.xmlGetDtdQElementDesc
func GetDtdQElementDesc(dtd DtdPtr, name *Char, prefix *Char) ElementPtr

//go:linkname GetDtdElementDesc C.xmlGetDtdElementDesc
func GetDtdElementDesc(dtd DtdPtr, name *Char) ElementPtr

// llgo:link (*ElementContent).ValidGetPotentialChildren C.xmlValidGetPotentialChildren
func (recv_ *ElementContent) ValidGetPotentialChildren(names **Char, len *c.Int, max c.Int) c.Int {
	return 0
}

// llgo:link (*Node).ValidGetValidElements C.xmlValidGetValidElements
func (recv_ *Node) ValidGetValidElements(next *Node, names **Char, max c.Int) c.Int {
	return 0
}

// llgo:link (*Char).ValidateNameValue C.xmlValidateNameValue
func (recv_ *Char) ValidateNameValue() c.Int {
	return 0
}

// llgo:link (*Char).ValidateNamesValue C.xmlValidateNamesValue
func (recv_ *Char) ValidateNamesValue() c.Int {
	return 0
}

// llgo:link (*Char).ValidateNmtokenValue C.xmlValidateNmtokenValue
func (recv_ *Char) ValidateNmtokenValue() c.Int {
	return 0
}

// llgo:link (*Char).ValidateNmtokensValue C.xmlValidateNmtokensValue
func (recv_ *Char) ValidateNmtokensValue() c.Int {
	return 0
}

/*
 * Validation based on the regexp support
 */
//go:linkname ValidBuildContentModel C.xmlValidBuildContentModel
func ValidBuildContentModel(ctxt ValidCtxtPtr, elem ElementPtr) c.Int

//go:linkname ValidatePushElement C.xmlValidatePushElement
func ValidatePushElement(ctxt ValidCtxtPtr, doc DocPtr, elem NodePtr, qname *Char) c.Int

//go:linkname ValidatePushCData C.xmlValidatePushCData
func ValidatePushCData(ctxt ValidCtxtPtr, data *Char, len c.Int) c.Int

//go:linkname ValidatePopElement C.xmlValidatePopElement
func ValidatePopElement(ctxt ValidCtxtPtr, doc DocPtr, elem NodePtr, qname *Char) c.Int
