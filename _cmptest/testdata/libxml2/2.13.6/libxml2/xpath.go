package libxml2

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

const XPATH_POINT = 5
const XPATH_RANGE = 6
const XPATH_LOCATIONSET = 7

type X_xmlXPathContext struct {
	Doc                DocPtr
	Node               NodePtr
	NbVariablesUnused  c.Int
	MaxVariablesUnused c.Int
	VarHash            HashTablePtr
	NbTypes            c.Int
	MaxTypes           c.Int
	Types              XPathTypePtr
	NbFuncsUnused      c.Int
	MaxFuncsUnused     c.Int
	FuncHash           HashTablePtr
	NbAxis             c.Int
	MaxAxis            c.Int
	Axis               XPathAxisPtr
	Namespaces         *NsPtr
	NsNr               c.Int
	User               c.Pointer
	ContextSize        c.Int
	ProximityPosition  c.Int
	Xptr               c.Int
	Here               NodePtr
	Origin             NodePtr
	NsHash             HashTablePtr
	VarLookupFunc      c.Pointer
	VarLookupData      c.Pointer
	Extra              c.Pointer
	Function           *Char
	FunctionURI        *Char
	FuncLookupFunc     c.Pointer
	FuncLookupData     c.Pointer
	TmpNsList          *NsPtr
	TmpNsNr            c.Int
	UserData           c.Pointer
	Error              c.Pointer
	LastError          Error
	DebugNode          NodePtr
	Dict               DictPtr
	Flags              c.Int
	Cache              c.Pointer
	OpLimit            c.Ulong
	OpCount            c.Ulong
	Depth              c.Int
}
type XPathContext X_xmlXPathContext
type XPathContextPtr *XPathContext

type X_xmlXPathParserContext struct {
	Cur        *Char
	Base       *Char
	Error      c.Int
	Context    XPathContextPtr
	Value      XPathObjectPtr
	ValueNr    c.Int
	ValueMax   c.Int
	ValueTab   *XPathObjectPtr
	Comp       XPathCompExprPtr
	Xptr       c.Int
	Ancestor   NodePtr
	ValueFrame c.Int
}
type XPathParserContext X_xmlXPathParserContext
type XPathParserContextPtr *XPathParserContext
type XPathError c.Int

const (
	XPATH_EXPRESSION_OK__1            XPathError = 0
	XPATH_NUMBER_ERROR__1             XPathError = 1
	XPATH_UNFINISHED_LITERAL_ERROR__1 XPathError = 2
	XPATH_START_LITERAL_ERROR__1      XPathError = 3
	XPATH_VARIABLE_REF_ERROR__1       XPathError = 4
	XPATH_UNDEF_VARIABLE_ERROR__1     XPathError = 5
	XPATH_INVALID_PREDICATE_ERROR__1  XPathError = 6
	XPATH_EXPR_ERROR__1               XPathError = 7
	XPATH_UNCLOSED_ERROR__1           XPathError = 8
	XPATH_UNKNOWN_FUNC_ERROR__1       XPathError = 9
	XPATH_INVALID_OPERAND__1          XPathError = 10
	XPATH_INVALID_TYPE__1             XPathError = 11
	XPATH_INVALID_ARITY__1            XPathError = 12
	XPATH_INVALID_CTXT_SIZE__1        XPathError = 13
	XPATH_INVALID_CTXT_POSITION__1    XPathError = 14
	XPATH_MEMORY_ERROR__1             XPathError = 15
	XPTR_SYNTAX_ERROR__1              XPathError = 16
	XPTR_RESOURCE_ERROR__1            XPathError = 17
	XPTR_SUB_RESOURCE_ERROR__1        XPathError = 18
	XPATH_UNDEF_PREFIX_ERROR__1       XPathError = 19
	XPATH_ENCODING_ERROR__1           XPathError = 20
	XPATH_INVALID_CHAR_ERROR__1       XPathError = 21
	XPATH_INVALID_CTXT                XPathError = 22
	XPATH_STACK_ERROR                 XPathError = 23
	XPATH_FORBID_VARIABLE_ERROR       XPathError = 24
	XPATH_OP_LIMIT_EXCEEDED           XPathError = 25
	XPATH_RECURSION_LIMIT_EXCEEDED    XPathError = 26
)

type X_xmlNodeSet struct {
	NodeNr  c.Int
	NodeMax c.Int
	NodeTab *NodePtr
}
type NodeSet X_xmlNodeSet
type NodeSetPtr *NodeSet
type XPathObjectType c.Int

const (
	XPATH_UNDEFINED XPathObjectType = 0
	XPATH_NODESET   XPathObjectType = 1
	XPATH_BOOLEAN   XPathObjectType = 2
	XPATH_NUMBER    XPathObjectType = 3
	XPATH_STRING    XPathObjectType = 4
	XPATH_USERS     XPathObjectType = 8
	XPATH_XSLT_TREE XPathObjectType = 9
)

type X_xmlXPathObject struct {
	Type       XPathObjectType
	Nodesetval NodeSetPtr
	Boolval    c.Int
	Floatval   c.Double
	Stringval  *Char
	User       c.Pointer
	Index      c.Int
	User2      c.Pointer
	Index2     c.Int
}
type XPathObject X_xmlXPathObject
type XPathObjectPtr *XPathObject

// llgo:type C
type XPathConvertFunc func(XPathObjectPtr, c.Int) c.Int

type X_xmlXPathType struct {
	Name *Char
	Func c.Pointer
}
type XPathType X_xmlXPathType
type XPathTypePtr *XPathType

type X_xmlXPathVariable struct {
	Name  *Char
	Value XPathObjectPtr
}
type XPathVariable X_xmlXPathVariable
type XPathVariablePtr *XPathVariable

// llgo:type C
type XPathEvalFunc func(XPathParserContextPtr, c.Int)

type X_xmlXPathFunct struct {
	Name *Char
	Func c.Pointer
}
type XPathFunct X_xmlXPathFunct
type XPathFuncPtr *XPathFunct

// llgo:type C
type XPathAxisFunc func(XPathParserContextPtr, XPathObjectPtr) XPathObjectPtr

type X_xmlXPathAxis struct {
	Name *Char
	Func c.Pointer
}
type XPathAxis X_xmlXPathAxis
type XPathAxisPtr *XPathAxis

// llgo:type C
type XPathFunction func(XPathParserContextPtr, c.Int)

// llgo:type C
type XPathVariableLookupFunc func(c.Pointer, *Char, *Char) XPathObjectPtr

// llgo:type C
type XPathFuncLookupFunc func(c.Pointer, *Char, *Char) XPathFunction

type X_xmlXPathCompExpr struct {
	Unused [8]uint8
}
type XPathCompExpr X_xmlXPathCompExpr
type XPathCompExprPtr *XPathCompExpr

//go:linkname XPathFreeObject C.xmlXPathFreeObject
func XPathFreeObject(obj XPathObjectPtr)

//go:linkname XPathNodeSetCreate C.xmlXPathNodeSetCreate
func XPathNodeSetCreate(val NodePtr) NodeSetPtr

//go:linkname XPathFreeNodeSetList C.xmlXPathFreeNodeSetList
func XPathFreeNodeSetList(obj XPathObjectPtr)

//go:linkname XPathFreeNodeSet C.xmlXPathFreeNodeSet
func XPathFreeNodeSet(obj NodeSetPtr)

//go:linkname XPathObjectCopy C.xmlXPathObjectCopy
func XPathObjectCopy(val XPathObjectPtr) XPathObjectPtr

//go:linkname XPathCmpNodes C.xmlXPathCmpNodes
func XPathCmpNodes(node1 NodePtr, node2 NodePtr) c.Int

/**
 * Conversion functions to basic types.
 */
//go:linkname XPathCastNumberToBoolean C.xmlXPathCastNumberToBoolean
func XPathCastNumberToBoolean(val c.Double) c.Int

// llgo:link (*Char).XPathCastStringToBoolean C.xmlXPathCastStringToBoolean
func (recv_ *Char) XPathCastStringToBoolean() c.Int {
	return 0
}

//go:linkname XPathCastNodeSetToBoolean C.xmlXPathCastNodeSetToBoolean
func XPathCastNodeSetToBoolean(ns NodeSetPtr) c.Int

//go:linkname XPathCastToBoolean C.xmlXPathCastToBoolean
func XPathCastToBoolean(val XPathObjectPtr) c.Int

//go:linkname XPathCastBooleanToNumber C.xmlXPathCastBooleanToNumber
func XPathCastBooleanToNumber(val c.Int) c.Double

// llgo:link (*Char).XPathCastStringToNumber C.xmlXPathCastStringToNumber
func (recv_ *Char) XPathCastStringToNumber() c.Double {
	return 0
}

//go:linkname XPathCastNodeToNumber C.xmlXPathCastNodeToNumber
func XPathCastNodeToNumber(node NodePtr) c.Double

//go:linkname XPathCastNodeSetToNumber C.xmlXPathCastNodeSetToNumber
func XPathCastNodeSetToNumber(ns NodeSetPtr) c.Double

//go:linkname XPathCastToNumber C.xmlXPathCastToNumber
func XPathCastToNumber(val XPathObjectPtr) c.Double

//go:linkname XPathCastBooleanToString C.xmlXPathCastBooleanToString
func XPathCastBooleanToString(val c.Int) *Char

//go:linkname XPathCastNumberToString C.xmlXPathCastNumberToString
func XPathCastNumberToString(val c.Double) *Char

//go:linkname XPathCastNodeToString C.xmlXPathCastNodeToString
func XPathCastNodeToString(node NodePtr) *Char

//go:linkname XPathCastNodeSetToString C.xmlXPathCastNodeSetToString
func XPathCastNodeSetToString(ns NodeSetPtr) *Char

//go:linkname XPathCastToString C.xmlXPathCastToString
func XPathCastToString(val XPathObjectPtr) *Char

//go:linkname XPathConvertBoolean C.xmlXPathConvertBoolean
func XPathConvertBoolean(val XPathObjectPtr) XPathObjectPtr

//go:linkname XPathConvertNumber C.xmlXPathConvertNumber
func XPathConvertNumber(val XPathObjectPtr) XPathObjectPtr

//go:linkname XPathConvertString C.xmlXPathConvertString
func XPathConvertString(val XPathObjectPtr) XPathObjectPtr

/**
 * Context handling.
 */
//go:linkname XPathNewContext C.xmlXPathNewContext
func XPathNewContext(doc DocPtr) XPathContextPtr

//go:linkname XPathFreeContext C.xmlXPathFreeContext
func XPathFreeContext(ctxt XPathContextPtr)

//go:linkname XPathSetErrorHandler C.xmlXPathSetErrorHandler
func XPathSetErrorHandler(ctxt XPathContextPtr, handler StructuredErrorFunc, context c.Pointer)

//go:linkname XPathContextSetCache C.xmlXPathContextSetCache
func XPathContextSetCache(ctxt XPathContextPtr, active c.Int, value c.Int, options c.Int) c.Int

/**
 * Evaluation functions.
 */
//go:linkname XPathOrderDocElems C.xmlXPathOrderDocElems
func XPathOrderDocElems(doc DocPtr) c.Long

//go:linkname XPathSetContextNode C.xmlXPathSetContextNode
func XPathSetContextNode(node NodePtr, ctx XPathContextPtr) c.Int

//go:linkname XPathNodeEval C.xmlXPathNodeEval
func XPathNodeEval(node NodePtr, str *Char, ctx XPathContextPtr) XPathObjectPtr

// llgo:link (*Char).XPathEval C.xmlXPathEval
func (recv_ *Char) XPathEval(ctx XPathContextPtr) XPathObjectPtr {
	return nil
}

// llgo:link (*Char).XPathEvalExpression C.xmlXPathEvalExpression
func (recv_ *Char) XPathEvalExpression(ctxt XPathContextPtr) XPathObjectPtr {
	return nil
}

//go:linkname XPathEvalPredicate C.xmlXPathEvalPredicate
func XPathEvalPredicate(ctxt XPathContextPtr, res XPathObjectPtr) c.Int

/**
 * Separate compilation/evaluation entry points.
 */
// llgo:link (*Char).XPathCompile C.xmlXPathCompile
func (recv_ *Char) XPathCompile() XPathCompExprPtr {
	return nil
}

//go:linkname XPathCtxtCompile C.xmlXPathCtxtCompile
func XPathCtxtCompile(ctxt XPathContextPtr, str *Char) XPathCompExprPtr

//go:linkname XPathCompiledEval C.xmlXPathCompiledEval
func XPathCompiledEval(comp XPathCompExprPtr, ctx XPathContextPtr) XPathObjectPtr

//go:linkname XPathCompiledEvalToBoolean C.xmlXPathCompiledEvalToBoolean
func XPathCompiledEvalToBoolean(comp XPathCompExprPtr, ctxt XPathContextPtr) c.Int

//go:linkname XPathFreeCompExpr C.xmlXPathFreeCompExpr
func XPathFreeCompExpr(comp XPathCompExprPtr)

//go:linkname XPathInit C.xmlXPathInit
func XPathInit()

//go:linkname XPathIsNaN C.xmlXPathIsNaN
func XPathIsNaN(val c.Double) c.Int

//go:linkname XPathIsInf C.xmlXPathIsInf
func XPathIsInf(val c.Double) c.Int
