package libxml2

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type SchemaWhitespaceValueType c.Int

const (
	SCHEMA_WHITESPACE_UNKNOWN  SchemaWhitespaceValueType = 0
	SCHEMA_WHITESPACE_PRESERVE SchemaWhitespaceValueType = 1
	SCHEMA_WHITESPACE_REPLACE  SchemaWhitespaceValueType = 2
	SCHEMA_WHITESPACE_COLLAPSE SchemaWhitespaceValueType = 3
)

//go:linkname SchemaInitTypes C.xmlSchemaInitTypes
func SchemaInitTypes() c.Int

//go:linkname SchemaCleanupTypes C.xmlSchemaCleanupTypes
func SchemaCleanupTypes()

// llgo:link (*Char).SchemaGetPredefinedType C.xmlSchemaGetPredefinedType
func (recv_ *Char) SchemaGetPredefinedType(ns *Char) SchemaTypePtr {
	return nil
}

//go:linkname SchemaValidatePredefinedType C.xmlSchemaValidatePredefinedType
func SchemaValidatePredefinedType(type_ SchemaTypePtr, value *Char, val *SchemaValPtr) c.Int

//go:linkname SchemaValPredefTypeNode C.xmlSchemaValPredefTypeNode
func SchemaValPredefTypeNode(type_ SchemaTypePtr, value *Char, val *SchemaValPtr, node NodePtr) c.Int

//go:linkname SchemaValidateFacet C.xmlSchemaValidateFacet
func SchemaValidateFacet(base SchemaTypePtr, facet SchemaFacetPtr, value *Char, val SchemaValPtr) c.Int

//go:linkname SchemaValidateFacetWhtsp C.xmlSchemaValidateFacetWhtsp
func SchemaValidateFacetWhtsp(facet SchemaFacetPtr, fws SchemaWhitespaceValueType, valType SchemaValType, value *Char, val SchemaValPtr, ws SchemaWhitespaceValueType) c.Int

//go:linkname SchemaFreeValue C.xmlSchemaFreeValue
func SchemaFreeValue(val SchemaValPtr)

//go:linkname SchemaNewFacet C.xmlSchemaNewFacet
func SchemaNewFacet() SchemaFacetPtr

//go:linkname SchemaCheckFacet C.xmlSchemaCheckFacet
func SchemaCheckFacet(facet SchemaFacetPtr, typeDecl SchemaTypePtr, ctxt SchemaParserCtxtPtr, name *Char) c.Int

//go:linkname SchemaFreeFacet C.xmlSchemaFreeFacet
func SchemaFreeFacet(facet SchemaFacetPtr)

//go:linkname SchemaCompareValues C.xmlSchemaCompareValues
func SchemaCompareValues(x SchemaValPtr, y SchemaValPtr) c.Int

//go:linkname SchemaGetBuiltInListSimpleTypeItemType C.xmlSchemaGetBuiltInListSimpleTypeItemType
func SchemaGetBuiltInListSimpleTypeItemType(type_ SchemaTypePtr) SchemaTypePtr

//go:linkname SchemaValidateListSimpleTypeFacet C.xmlSchemaValidateListSimpleTypeFacet
func SchemaValidateListSimpleTypeFacet(facet SchemaFacetPtr, value *Char, actualLen c.Ulong, expectedLen *c.Ulong) c.Int

// llgo:link SchemaValType.SchemaGetBuiltInType C.xmlSchemaGetBuiltInType
func (recv_ SchemaValType) SchemaGetBuiltInType() SchemaTypePtr {
	return nil
}

//go:linkname SchemaIsBuiltInTypeFacet C.xmlSchemaIsBuiltInTypeFacet
func SchemaIsBuiltInTypeFacet(type_ SchemaTypePtr, facetType c.Int) c.Int

// llgo:link (*Char).SchemaCollapseString C.xmlSchemaCollapseString
func (recv_ *Char) SchemaCollapseString() *Char {
	return nil
}

// llgo:link (*Char).SchemaWhiteSpaceReplace C.xmlSchemaWhiteSpaceReplace
func (recv_ *Char) SchemaWhiteSpaceReplace() *Char {
	return nil
}

//go:linkname SchemaGetFacetValueAsULong C.xmlSchemaGetFacetValueAsULong
func SchemaGetFacetValueAsULong(facet SchemaFacetPtr) c.Ulong

//go:linkname SchemaValidateLengthFacet C.xmlSchemaValidateLengthFacet
func SchemaValidateLengthFacet(type_ SchemaTypePtr, facet SchemaFacetPtr, value *Char, val SchemaValPtr, length *c.Ulong) c.Int

//go:linkname SchemaValidateLengthFacetWhtsp C.xmlSchemaValidateLengthFacetWhtsp
func SchemaValidateLengthFacetWhtsp(facet SchemaFacetPtr, valType SchemaValType, value *Char, val SchemaValPtr, length *c.Ulong, ws SchemaWhitespaceValueType) c.Int

//go:linkname SchemaValPredefTypeNodeNoNorm C.xmlSchemaValPredefTypeNodeNoNorm
func SchemaValPredefTypeNodeNoNorm(type_ SchemaTypePtr, value *Char, val *SchemaValPtr, node NodePtr) c.Int

//go:linkname SchemaGetCanonValue C.xmlSchemaGetCanonValue
func SchemaGetCanonValue(val SchemaValPtr, retValue **Char) c.Int

//go:linkname SchemaGetCanonValueWhtsp C.xmlSchemaGetCanonValueWhtsp
func SchemaGetCanonValueWhtsp(val SchemaValPtr, retValue **Char, ws SchemaWhitespaceValueType) c.Int

//go:linkname SchemaValueAppend C.xmlSchemaValueAppend
func SchemaValueAppend(prev SchemaValPtr, cur SchemaValPtr) c.Int

//go:linkname SchemaValueGetNext C.xmlSchemaValueGetNext
func SchemaValueGetNext(cur SchemaValPtr) SchemaValPtr

//go:linkname SchemaValueGetAsString C.xmlSchemaValueGetAsString
func SchemaValueGetAsString(val SchemaValPtr) *Char

//go:linkname SchemaValueGetAsBoolean C.xmlSchemaValueGetAsBoolean
func SchemaValueGetAsBoolean(val SchemaValPtr) c.Int

// llgo:link SchemaValType.SchemaNewStringValue C.xmlSchemaNewStringValue
func (recv_ SchemaValType) SchemaNewStringValue(value *Char) SchemaValPtr {
	return nil
}

// llgo:link (*Char).SchemaNewNOTATIONValue C.xmlSchemaNewNOTATIONValue
func (recv_ *Char) SchemaNewNOTATIONValue(ns *Char) SchemaValPtr {
	return nil
}

// llgo:link (*Char).SchemaNewQNameValue C.xmlSchemaNewQNameValue
func (recv_ *Char) SchemaNewQNameValue(localName *Char) SchemaValPtr {
	return nil
}

//go:linkname SchemaCompareValuesWhtsp C.xmlSchemaCompareValuesWhtsp
func SchemaCompareValuesWhtsp(x SchemaValPtr, xws SchemaWhitespaceValueType, y SchemaValPtr, yws SchemaWhitespaceValueType) c.Int

//go:linkname SchemaCopyValue C.xmlSchemaCopyValue
func SchemaCopyValue(val SchemaValPtr) SchemaValPtr

//go:linkname SchemaGetValType C.xmlSchemaGetValType
func SchemaGetValType(val SchemaValPtr) SchemaValType
