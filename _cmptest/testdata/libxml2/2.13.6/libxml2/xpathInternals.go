package libxml2

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

//go:linkname XPathPopBoolean C.xmlXPathPopBoolean
func XPathPopBoolean(ctxt XPathParserContextPtr) c.Int

//go:linkname XPathPopNumber C.xmlXPathPopNumber
func XPathPopNumber(ctxt XPathParserContextPtr) c.Double

//go:linkname XPathPopString C.xmlXPathPopString
func XPathPopString(ctxt XPathParserContextPtr) *Char

//go:linkname XPathPopNodeSet C.xmlXPathPopNodeSet
func XPathPopNodeSet(ctxt XPathParserContextPtr) NodeSetPtr

//go:linkname XPathPopExternal C.xmlXPathPopExternal
func XPathPopExternal(ctxt XPathParserContextPtr) c.Pointer

/*
 * Variable Lookup forwarding.
 */
//go:linkname XPathRegisterVariableLookup C.xmlXPathRegisterVariableLookup
func XPathRegisterVariableLookup(ctxt XPathContextPtr, f XPathVariableLookupFunc, data c.Pointer)

/*
 * Function Lookup forwarding.
 */
//go:linkname XPathRegisterFuncLookup C.xmlXPathRegisterFuncLookup
func XPathRegisterFuncLookup(ctxt XPathContextPtr, f XPathFuncLookupFunc, funcCtxt c.Pointer)

/*
 * Error reporting.
 */
//go:linkname XPatherror C.xmlXPatherror
func XPatherror(ctxt XPathParserContextPtr, file *c.Char, line c.Int, no c.Int)

//go:linkname XPathErr C.xmlXPathErr
func XPathErr(ctxt XPathParserContextPtr, error c.Int)

//go:linkname XPathDebugDumpObject C.xmlXPathDebugDumpObject
func XPathDebugDumpObject(output *c.FILE, cur XPathObjectPtr, depth c.Int)

//go:linkname XPathDebugDumpCompExpr C.xmlXPathDebugDumpCompExpr
func XPathDebugDumpCompExpr(output *c.FILE, comp XPathCompExprPtr, depth c.Int)

/**
 * NodeSet handling.
 */
//go:linkname XPathNodeSetContains C.xmlXPathNodeSetContains
func XPathNodeSetContains(cur NodeSetPtr, val NodePtr) c.Int

//go:linkname XPathDifference C.xmlXPathDifference
func XPathDifference(nodes1 NodeSetPtr, nodes2 NodeSetPtr) NodeSetPtr

//go:linkname XPathIntersection C.xmlXPathIntersection
func XPathIntersection(nodes1 NodeSetPtr, nodes2 NodeSetPtr) NodeSetPtr

//go:linkname XPathDistinctSorted C.xmlXPathDistinctSorted
func XPathDistinctSorted(nodes NodeSetPtr) NodeSetPtr

//go:linkname XPathDistinct C.xmlXPathDistinct
func XPathDistinct(nodes NodeSetPtr) NodeSetPtr

//go:linkname XPathHasSameNodes C.xmlXPathHasSameNodes
func XPathHasSameNodes(nodes1 NodeSetPtr, nodes2 NodeSetPtr) c.Int

//go:linkname XPathNodeLeadingSorted C.xmlXPathNodeLeadingSorted
func XPathNodeLeadingSorted(nodes NodeSetPtr, node NodePtr) NodeSetPtr

//go:linkname XPathLeadingSorted C.xmlXPathLeadingSorted
func XPathLeadingSorted(nodes1 NodeSetPtr, nodes2 NodeSetPtr) NodeSetPtr

//go:linkname XPathNodeLeading C.xmlXPathNodeLeading
func XPathNodeLeading(nodes NodeSetPtr, node NodePtr) NodeSetPtr

//go:linkname XPathLeading C.xmlXPathLeading
func XPathLeading(nodes1 NodeSetPtr, nodes2 NodeSetPtr) NodeSetPtr

//go:linkname XPathNodeTrailingSorted C.xmlXPathNodeTrailingSorted
func XPathNodeTrailingSorted(nodes NodeSetPtr, node NodePtr) NodeSetPtr

//go:linkname XPathTrailingSorted C.xmlXPathTrailingSorted
func XPathTrailingSorted(nodes1 NodeSetPtr, nodes2 NodeSetPtr) NodeSetPtr

//go:linkname XPathNodeTrailing C.xmlXPathNodeTrailing
func XPathNodeTrailing(nodes NodeSetPtr, node NodePtr) NodeSetPtr

//go:linkname XPathTrailing C.xmlXPathTrailing
func XPathTrailing(nodes1 NodeSetPtr, nodes2 NodeSetPtr) NodeSetPtr

/**
 * Extending a context.
 */
//go:linkname XPathRegisterNs C.xmlXPathRegisterNs
func XPathRegisterNs(ctxt XPathContextPtr, prefix *Char, ns_uri *Char) c.Int

//go:linkname XPathNsLookup C.xmlXPathNsLookup
func XPathNsLookup(ctxt XPathContextPtr, prefix *Char) *Char

//go:linkname XPathRegisteredNsCleanup C.xmlXPathRegisteredNsCleanup
func XPathRegisteredNsCleanup(ctxt XPathContextPtr)

//go:linkname XPathRegisterFunc C.xmlXPathRegisterFunc
func XPathRegisterFunc(ctxt XPathContextPtr, name *Char, f XPathFunction) c.Int

//go:linkname XPathRegisterFuncNS C.xmlXPathRegisterFuncNS
func XPathRegisterFuncNS(ctxt XPathContextPtr, name *Char, ns_uri *Char, f XPathFunction) c.Int

//go:linkname XPathRegisterVariable C.xmlXPathRegisterVariable
func XPathRegisterVariable(ctxt XPathContextPtr, name *Char, value XPathObjectPtr) c.Int

//go:linkname XPathRegisterVariableNS C.xmlXPathRegisterVariableNS
func XPathRegisterVariableNS(ctxt XPathContextPtr, name *Char, ns_uri *Char, value XPathObjectPtr) c.Int

//go:linkname XPathFunctionLookup C.xmlXPathFunctionLookup
func XPathFunctionLookup(ctxt XPathContextPtr, name *Char) XPathFunction

//go:linkname XPathFunctionLookupNS C.xmlXPathFunctionLookupNS
func XPathFunctionLookupNS(ctxt XPathContextPtr, name *Char, ns_uri *Char) XPathFunction

//go:linkname XPathRegisteredFuncsCleanup C.xmlXPathRegisteredFuncsCleanup
func XPathRegisteredFuncsCleanup(ctxt XPathContextPtr)

//go:linkname XPathVariableLookup C.xmlXPathVariableLookup
func XPathVariableLookup(ctxt XPathContextPtr, name *Char) XPathObjectPtr

//go:linkname XPathVariableLookupNS C.xmlXPathVariableLookupNS
func XPathVariableLookupNS(ctxt XPathContextPtr, name *Char, ns_uri *Char) XPathObjectPtr

//go:linkname XPathRegisteredVariablesCleanup C.xmlXPathRegisteredVariablesCleanup
func XPathRegisteredVariablesCleanup(ctxt XPathContextPtr)

/**
 * Utilities to extend XPath.
 */
// llgo:link (*Char).XPathNewParserContext C.xmlXPathNewParserContext
func (recv_ *Char) XPathNewParserContext(ctxt XPathContextPtr) XPathParserContextPtr {
	return nil
}

//go:linkname XPathFreeParserContext C.xmlXPathFreeParserContext
func XPathFreeParserContext(ctxt XPathParserContextPtr)

/* TODO: remap to xmlXPathValuePop and Push. */
//go:linkname ValuePop C.valuePop
func ValuePop(ctxt XPathParserContextPtr) XPathObjectPtr

//go:linkname ValuePush C.valuePush
func ValuePush(ctxt XPathParserContextPtr, value XPathObjectPtr) c.Int

// llgo:link (*Char).XPathNewString C.xmlXPathNewString
func (recv_ *Char) XPathNewString() XPathObjectPtr {
	return nil
}

//go:linkname XPathNewCString C.xmlXPathNewCString
func XPathNewCString(val *c.Char) XPathObjectPtr

// llgo:link (*Char).XPathWrapString C.xmlXPathWrapString
func (recv_ *Char) XPathWrapString() XPathObjectPtr {
	return nil
}

//go:linkname XPathWrapCString C.xmlXPathWrapCString
func XPathWrapCString(val *c.Char) XPathObjectPtr

//go:linkname XPathNewFloat C.xmlXPathNewFloat
func XPathNewFloat(val c.Double) XPathObjectPtr

//go:linkname XPathNewBoolean C.xmlXPathNewBoolean
func XPathNewBoolean(val c.Int) XPathObjectPtr

//go:linkname XPathNewNodeSet C.xmlXPathNewNodeSet
func XPathNewNodeSet(val NodePtr) XPathObjectPtr

//go:linkname XPathNewValueTree C.xmlXPathNewValueTree
func XPathNewValueTree(val NodePtr) XPathObjectPtr

//go:linkname XPathNodeSetAdd C.xmlXPathNodeSetAdd
func XPathNodeSetAdd(cur NodeSetPtr, val NodePtr) c.Int

//go:linkname XPathNodeSetAddUnique C.xmlXPathNodeSetAddUnique
func XPathNodeSetAddUnique(cur NodeSetPtr, val NodePtr) c.Int

//go:linkname XPathNodeSetAddNs C.xmlXPathNodeSetAddNs
func XPathNodeSetAddNs(cur NodeSetPtr, node NodePtr, ns NsPtr) c.Int

//go:linkname XPathNodeSetSort C.xmlXPathNodeSetSort
func XPathNodeSetSort(set NodeSetPtr)

//go:linkname XPathRoot C.xmlXPathRoot
func XPathRoot(ctxt XPathParserContextPtr)

//go:linkname XPathEvalExpr C.xmlXPathEvalExpr
func XPathEvalExpr(ctxt XPathParserContextPtr)

//go:linkname XPathParseName C.xmlXPathParseName
func XPathParseName(ctxt XPathParserContextPtr) *Char

//go:linkname XPathParseNCName C.xmlXPathParseNCName
func XPathParseNCName(ctxt XPathParserContextPtr) *Char

/*
 * Existing functions.
 */
// llgo:link (*Char).XPathStringEvalNumber C.xmlXPathStringEvalNumber
func (recv_ *Char) XPathStringEvalNumber() c.Double {
	return 0
}

//go:linkname XPathEvaluatePredicateResult C.xmlXPathEvaluatePredicateResult
func XPathEvaluatePredicateResult(ctxt XPathParserContextPtr, res XPathObjectPtr) c.Int

//go:linkname XPathRegisterAllFunctions C.xmlXPathRegisterAllFunctions
func XPathRegisterAllFunctions(ctxt XPathContextPtr)

//go:linkname XPathNodeSetMerge C.xmlXPathNodeSetMerge
func XPathNodeSetMerge(val1 NodeSetPtr, val2 NodeSetPtr) NodeSetPtr

//go:linkname XPathNodeSetDel C.xmlXPathNodeSetDel
func XPathNodeSetDel(cur NodeSetPtr, val NodePtr)

//go:linkname XPathNodeSetRemove C.xmlXPathNodeSetRemove
func XPathNodeSetRemove(cur NodeSetPtr, val c.Int)

//go:linkname XPathNewNodeSetList C.xmlXPathNewNodeSetList
func XPathNewNodeSetList(val NodeSetPtr) XPathObjectPtr

//go:linkname XPathWrapNodeSet C.xmlXPathWrapNodeSet
func XPathWrapNodeSet(val NodeSetPtr) XPathObjectPtr

//go:linkname XPathWrapExternal C.xmlXPathWrapExternal
func XPathWrapExternal(val c.Pointer) XPathObjectPtr

//go:linkname XPathEqualValues C.xmlXPathEqualValues
func XPathEqualValues(ctxt XPathParserContextPtr) c.Int

//go:linkname XPathNotEqualValues C.xmlXPathNotEqualValues
func XPathNotEqualValues(ctxt XPathParserContextPtr) c.Int

//go:linkname XPathCompareValues C.xmlXPathCompareValues
func XPathCompareValues(ctxt XPathParserContextPtr, inf c.Int, strict c.Int) c.Int

//go:linkname XPathValueFlipSign C.xmlXPathValueFlipSign
func XPathValueFlipSign(ctxt XPathParserContextPtr)

//go:linkname XPathAddValues C.xmlXPathAddValues
func XPathAddValues(ctxt XPathParserContextPtr)

//go:linkname XPathSubValues C.xmlXPathSubValues
func XPathSubValues(ctxt XPathParserContextPtr)

//go:linkname XPathMultValues C.xmlXPathMultValues
func XPathMultValues(ctxt XPathParserContextPtr)

//go:linkname XPathDivValues C.xmlXPathDivValues
func XPathDivValues(ctxt XPathParserContextPtr)

//go:linkname XPathModValues C.xmlXPathModValues
func XPathModValues(ctxt XPathParserContextPtr)

// llgo:link (*Char).XPathIsNodeType C.xmlXPathIsNodeType
func (recv_ *Char) XPathIsNodeType() c.Int {
	return 0
}

/*
 * Some of the axis navigation routines.
 */
//go:linkname XPathNextSelf C.xmlXPathNextSelf
func XPathNextSelf(ctxt XPathParserContextPtr, cur NodePtr) NodePtr

//go:linkname XPathNextChild C.xmlXPathNextChild
func XPathNextChild(ctxt XPathParserContextPtr, cur NodePtr) NodePtr

//go:linkname XPathNextDescendant C.xmlXPathNextDescendant
func XPathNextDescendant(ctxt XPathParserContextPtr, cur NodePtr) NodePtr

//go:linkname XPathNextDescendantOrSelf C.xmlXPathNextDescendantOrSelf
func XPathNextDescendantOrSelf(ctxt XPathParserContextPtr, cur NodePtr) NodePtr

//go:linkname XPathNextParent C.xmlXPathNextParent
func XPathNextParent(ctxt XPathParserContextPtr, cur NodePtr) NodePtr

//go:linkname XPathNextAncestorOrSelf C.xmlXPathNextAncestorOrSelf
func XPathNextAncestorOrSelf(ctxt XPathParserContextPtr, cur NodePtr) NodePtr

//go:linkname XPathNextFollowingSibling C.xmlXPathNextFollowingSibling
func XPathNextFollowingSibling(ctxt XPathParserContextPtr, cur NodePtr) NodePtr

//go:linkname XPathNextFollowing C.xmlXPathNextFollowing
func XPathNextFollowing(ctxt XPathParserContextPtr, cur NodePtr) NodePtr

//go:linkname XPathNextNamespace C.xmlXPathNextNamespace
func XPathNextNamespace(ctxt XPathParserContextPtr, cur NodePtr) NodePtr

//go:linkname XPathNextAttribute C.xmlXPathNextAttribute
func XPathNextAttribute(ctxt XPathParserContextPtr, cur NodePtr) NodePtr

//go:linkname XPathNextPreceding C.xmlXPathNextPreceding
func XPathNextPreceding(ctxt XPathParserContextPtr, cur NodePtr) NodePtr

//go:linkname XPathNextAncestor C.xmlXPathNextAncestor
func XPathNextAncestor(ctxt XPathParserContextPtr, cur NodePtr) NodePtr

//go:linkname XPathNextPrecedingSibling C.xmlXPathNextPrecedingSibling
func XPathNextPrecedingSibling(ctxt XPathParserContextPtr, cur NodePtr) NodePtr

/*
 * The official core of XPath functions.
 */
//go:linkname XPathLastFunction C.xmlXPathLastFunction
func XPathLastFunction(ctxt XPathParserContextPtr, nargs c.Int)

//go:linkname XPathPositionFunction C.xmlXPathPositionFunction
func XPathPositionFunction(ctxt XPathParserContextPtr, nargs c.Int)

//go:linkname XPathCountFunction C.xmlXPathCountFunction
func XPathCountFunction(ctxt XPathParserContextPtr, nargs c.Int)

//go:linkname XPathIdFunction C.xmlXPathIdFunction
func XPathIdFunction(ctxt XPathParserContextPtr, nargs c.Int)

//go:linkname XPathLocalNameFunction C.xmlXPathLocalNameFunction
func XPathLocalNameFunction(ctxt XPathParserContextPtr, nargs c.Int)

//go:linkname XPathNamespaceURIFunction C.xmlXPathNamespaceURIFunction
func XPathNamespaceURIFunction(ctxt XPathParserContextPtr, nargs c.Int)

//go:linkname XPathStringFunction C.xmlXPathStringFunction
func XPathStringFunction(ctxt XPathParserContextPtr, nargs c.Int)

//go:linkname XPathStringLengthFunction C.xmlXPathStringLengthFunction
func XPathStringLengthFunction(ctxt XPathParserContextPtr, nargs c.Int)

//go:linkname XPathConcatFunction C.xmlXPathConcatFunction
func XPathConcatFunction(ctxt XPathParserContextPtr, nargs c.Int)

//go:linkname XPathContainsFunction C.xmlXPathContainsFunction
func XPathContainsFunction(ctxt XPathParserContextPtr, nargs c.Int)

//go:linkname XPathStartsWithFunction C.xmlXPathStartsWithFunction
func XPathStartsWithFunction(ctxt XPathParserContextPtr, nargs c.Int)

//go:linkname XPathSubstringFunction C.xmlXPathSubstringFunction
func XPathSubstringFunction(ctxt XPathParserContextPtr, nargs c.Int)

//go:linkname XPathSubstringBeforeFunction C.xmlXPathSubstringBeforeFunction
func XPathSubstringBeforeFunction(ctxt XPathParserContextPtr, nargs c.Int)

//go:linkname XPathSubstringAfterFunction C.xmlXPathSubstringAfterFunction
func XPathSubstringAfterFunction(ctxt XPathParserContextPtr, nargs c.Int)

//go:linkname XPathNormalizeFunction C.xmlXPathNormalizeFunction
func XPathNormalizeFunction(ctxt XPathParserContextPtr, nargs c.Int)

//go:linkname XPathTranslateFunction C.xmlXPathTranslateFunction
func XPathTranslateFunction(ctxt XPathParserContextPtr, nargs c.Int)

//go:linkname XPathNotFunction C.xmlXPathNotFunction
func XPathNotFunction(ctxt XPathParserContextPtr, nargs c.Int)

//go:linkname XPathTrueFunction C.xmlXPathTrueFunction
func XPathTrueFunction(ctxt XPathParserContextPtr, nargs c.Int)

//go:linkname XPathFalseFunction C.xmlXPathFalseFunction
func XPathFalseFunction(ctxt XPathParserContextPtr, nargs c.Int)

//go:linkname XPathLangFunction C.xmlXPathLangFunction
func XPathLangFunction(ctxt XPathParserContextPtr, nargs c.Int)

//go:linkname XPathNumberFunction C.xmlXPathNumberFunction
func XPathNumberFunction(ctxt XPathParserContextPtr, nargs c.Int)

//go:linkname XPathSumFunction C.xmlXPathSumFunction
func XPathSumFunction(ctxt XPathParserContextPtr, nargs c.Int)

//go:linkname XPathFloorFunction C.xmlXPathFloorFunction
func XPathFloorFunction(ctxt XPathParserContextPtr, nargs c.Int)

//go:linkname XPathCeilingFunction C.xmlXPathCeilingFunction
func XPathCeilingFunction(ctxt XPathParserContextPtr, nargs c.Int)

//go:linkname XPathRoundFunction C.xmlXPathRoundFunction
func XPathRoundFunction(ctxt XPathParserContextPtr, nargs c.Int)

//go:linkname XPathBooleanFunction C.xmlXPathBooleanFunction
func XPathBooleanFunction(ctxt XPathParserContextPtr, nargs c.Int)

/**
 * Really internal functions
 */
//go:linkname XPathNodeSetFreeNs C.xmlXPathNodeSetFreeNs
func XPathNodeSetFreeNs(ns NsPtr)
