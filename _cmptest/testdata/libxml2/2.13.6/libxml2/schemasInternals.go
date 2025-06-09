package libxml2

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

const SCHEMAS_ANYATTR_SKIP = 1
const SCHEMAS_ANYATTR_LAX = 2
const SCHEMAS_ANYATTR_STRICT = 3
const SCHEMAS_ANY_SKIP = 1
const SCHEMAS_ANY_LAX = 2
const SCHEMAS_ANY_STRICT = 3
const SCHEMAS_ATTR_USE_PROHIBITED = 0
const SCHEMAS_ATTR_USE_REQUIRED = 1
const SCHEMAS_ATTR_USE_OPTIONAL = 2
const SCHEMAS_FACET_UNKNOWN = 0
const SCHEMAS_FACET_PRESERVE = 1
const SCHEMAS_FACET_REPLACE = 2
const SCHEMAS_FACET_COLLAPSE = 3

type SchemaValType c.Int

const (
	SCHEMAS_UNKNOWN       SchemaValType = 0
	SCHEMAS_STRING        SchemaValType = 1
	SCHEMAS_NORMSTRING    SchemaValType = 2
	SCHEMAS_DECIMAL       SchemaValType = 3
	SCHEMAS_TIME          SchemaValType = 4
	SCHEMAS_GDAY          SchemaValType = 5
	SCHEMAS_GMONTH        SchemaValType = 6
	SCHEMAS_GMONTHDAY     SchemaValType = 7
	SCHEMAS_GYEAR         SchemaValType = 8
	SCHEMAS_GYEARMONTH    SchemaValType = 9
	SCHEMAS_DATE          SchemaValType = 10
	SCHEMAS_DATETIME      SchemaValType = 11
	SCHEMAS_DURATION      SchemaValType = 12
	SCHEMAS_FLOAT         SchemaValType = 13
	SCHEMAS_DOUBLE        SchemaValType = 14
	SCHEMAS_BOOLEAN       SchemaValType = 15
	SCHEMAS_TOKEN         SchemaValType = 16
	SCHEMAS_LANGUAGE      SchemaValType = 17
	SCHEMAS_NMTOKEN       SchemaValType = 18
	SCHEMAS_NMTOKENS      SchemaValType = 19
	SCHEMAS_NAME          SchemaValType = 20
	SCHEMAS_QNAME         SchemaValType = 21
	SCHEMAS_NCNAME        SchemaValType = 22
	SCHEMAS_ID            SchemaValType = 23
	SCHEMAS_IDREF         SchemaValType = 24
	SCHEMAS_IDREFS        SchemaValType = 25
	SCHEMAS_ENTITY        SchemaValType = 26
	SCHEMAS_ENTITIES      SchemaValType = 27
	SCHEMAS_NOTATION      SchemaValType = 28
	SCHEMAS_ANYURI        SchemaValType = 29
	SCHEMAS_INTEGER       SchemaValType = 30
	SCHEMAS_NPINTEGER     SchemaValType = 31
	SCHEMAS_NINTEGER      SchemaValType = 32
	SCHEMAS_NNINTEGER     SchemaValType = 33
	SCHEMAS_PINTEGER      SchemaValType = 34
	SCHEMAS_INT           SchemaValType = 35
	SCHEMAS_UINT          SchemaValType = 36
	SCHEMAS_LONG          SchemaValType = 37
	SCHEMAS_ULONG         SchemaValType = 38
	SCHEMAS_SHORT         SchemaValType = 39
	SCHEMAS_USHORT        SchemaValType = 40
	SCHEMAS_BYTE          SchemaValType = 41
	SCHEMAS_UBYTE         SchemaValType = 42
	SCHEMAS_HEXBINARY     SchemaValType = 43
	SCHEMAS_BASE64BINARY  SchemaValType = 44
	SCHEMAS_ANYTYPE       SchemaValType = 45
	SCHEMAS_ANYSIMPLETYPE SchemaValType = 46
)

type SchemaTypeType c.Int

const (
	SCHEMA_TYPE_BASIC            SchemaTypeType = 1
	SCHEMA_TYPE_ANY              SchemaTypeType = 2
	SCHEMA_TYPE_FACET            SchemaTypeType = 3
	SCHEMA_TYPE_SIMPLE           SchemaTypeType = 4
	SCHEMA_TYPE_COMPLEX          SchemaTypeType = 5
	SCHEMA_TYPE_SEQUENCE         SchemaTypeType = 6
	SCHEMA_TYPE_CHOICE           SchemaTypeType = 7
	SCHEMA_TYPE_ALL              SchemaTypeType = 8
	SCHEMA_TYPE_SIMPLE_CONTENT   SchemaTypeType = 9
	SCHEMA_TYPE_COMPLEX_CONTENT  SchemaTypeType = 10
	SCHEMA_TYPE_UR               SchemaTypeType = 11
	SCHEMA_TYPE_RESTRICTION      SchemaTypeType = 12
	SCHEMA_TYPE_EXTENSION        SchemaTypeType = 13
	SCHEMA_TYPE_ELEMENT          SchemaTypeType = 14
	SCHEMA_TYPE_ATTRIBUTE        SchemaTypeType = 15
	SCHEMA_TYPE_ATTRIBUTEGROUP   SchemaTypeType = 16
	SCHEMA_TYPE_GROUP            SchemaTypeType = 17
	SCHEMA_TYPE_NOTATION         SchemaTypeType = 18
	SCHEMA_TYPE_LIST             SchemaTypeType = 19
	SCHEMA_TYPE_UNION            SchemaTypeType = 20
	SCHEMA_TYPE_ANY_ATTRIBUTE    SchemaTypeType = 21
	SCHEMA_TYPE_IDC_UNIQUE       SchemaTypeType = 22
	SCHEMA_TYPE_IDC_KEY          SchemaTypeType = 23
	SCHEMA_TYPE_IDC_KEYREF       SchemaTypeType = 24
	SCHEMA_TYPE_PARTICLE         SchemaTypeType = 25
	SCHEMA_TYPE_ATTRIBUTE_USE    SchemaTypeType = 26
	SCHEMA_FACET_MININCLUSIVE    SchemaTypeType = 1000
	SCHEMA_FACET_MINEXCLUSIVE    SchemaTypeType = 1001
	SCHEMA_FACET_MAXINCLUSIVE    SchemaTypeType = 1002
	SCHEMA_FACET_MAXEXCLUSIVE    SchemaTypeType = 1003
	SCHEMA_FACET_TOTALDIGITS     SchemaTypeType = 1004
	SCHEMA_FACET_FRACTIONDIGITS  SchemaTypeType = 1005
	SCHEMA_FACET_PATTERN         SchemaTypeType = 1006
	SCHEMA_FACET_ENUMERATION     SchemaTypeType = 1007
	SCHEMA_FACET_WHITESPACE      SchemaTypeType = 1008
	SCHEMA_FACET_LENGTH          SchemaTypeType = 1009
	SCHEMA_FACET_MAXLENGTH       SchemaTypeType = 1010
	SCHEMA_FACET_MINLENGTH       SchemaTypeType = 1011
	SCHEMA_EXTRA_QNAMEREF        SchemaTypeType = 2000
	SCHEMA_EXTRA_ATTR_USE_PROHIB SchemaTypeType = 2001
)

type SchemaContentType c.Int

const (
	SCHEMA_CONTENT_UNKNOWN           SchemaContentType = 0
	SCHEMA_CONTENT_EMPTY             SchemaContentType = 1
	SCHEMA_CONTENT_ELEMENTS          SchemaContentType = 2
	SCHEMA_CONTENT_MIXED             SchemaContentType = 3
	SCHEMA_CONTENT_SIMPLE            SchemaContentType = 4
	SCHEMA_CONTENT_MIXED_OR_ELEMENTS SchemaContentType = 5
	SCHEMA_CONTENT_BASIC             SchemaContentType = 6
	SCHEMA_CONTENT_ANY               SchemaContentType = 7
)

type X_xmlSchemaVal struct {
	Unused [8]uint8
}
type SchemaVal X_xmlSchemaVal
type SchemaValPtr *SchemaVal

type X_xmlSchemaType struct {
	Type              SchemaTypeType
	Next              *X_xmlSchemaType
	Name              *Char
	Id                *Char
	Ref               *Char
	RefNs             *Char
	Annot             SchemaAnnotPtr
	Subtypes          SchemaTypePtr
	Attributes        SchemaAttributePtr
	Node              NodePtr
	MinOccurs         c.Int
	MaxOccurs         c.Int
	Flags             c.Int
	ContentType       SchemaContentType
	Base              *Char
	BaseNs            *Char
	BaseType          SchemaTypePtr
	Facets            SchemaFacetPtr
	Redef             *X_xmlSchemaType
	Recurse           c.Int
	AttributeUses     *SchemaAttributeLinkPtr
	AttributeWildcard SchemaWildcardPtr
	BuiltInType       c.Int
	MemberTypes       SchemaTypeLinkPtr
	FacetSet          SchemaFacetLinkPtr
	RefPrefix         *Char
	ContentTypeDef    SchemaTypePtr
	ContModel         RegexpPtr
	TargetNamespace   *Char
	AttrUses          c.Pointer
}
type SchemaType X_xmlSchemaType
type SchemaTypePtr *SchemaType

type X_xmlSchemaFacet struct {
	Type       SchemaTypeType
	Next       *X_xmlSchemaFacet
	Value      *Char
	Id         *Char
	Annot      SchemaAnnotPtr
	Node       NodePtr
	Fixed      c.Int
	Whitespace c.Int
	Val        SchemaValPtr
	Regexp     RegexpPtr
}
type SchemaFacet X_xmlSchemaFacet
type SchemaFacetPtr *SchemaFacet

type X_xmlSchemaAnnot struct {
	Next    *X_xmlSchemaAnnot
	Content NodePtr
}
type SchemaAnnot X_xmlSchemaAnnot
type SchemaAnnotPtr *SchemaAnnot

type X_xmlSchemaAttribute struct {
	Type            SchemaTypeType
	Next            *X_xmlSchemaAttribute
	Name            *Char
	Id              *Char
	Ref             *Char
	RefNs           *Char
	TypeName        *Char
	TypeNs          *Char
	Annot           SchemaAnnotPtr
	Base            SchemaTypePtr
	Occurs          c.Int
	DefValue        *Char
	Subtypes        SchemaTypePtr
	Node            NodePtr
	TargetNamespace *Char
	Flags           c.Int
	RefPrefix       *Char
	DefVal          SchemaValPtr
	RefDecl         SchemaAttributePtr
}
type SchemaAttribute X_xmlSchemaAttribute
type SchemaAttributePtr *SchemaAttribute

type X_xmlSchemaAttributeLink struct {
	Next *X_xmlSchemaAttributeLink
	Attr *X_xmlSchemaAttribute
}
type SchemaAttributeLink X_xmlSchemaAttributeLink
type SchemaAttributeLinkPtr *SchemaAttributeLink

type X_xmlSchemaWildcardNs struct {
	Next  *X_xmlSchemaWildcardNs
	Value *Char
}
type SchemaWildcardNs X_xmlSchemaWildcardNs
type SchemaWildcardNsPtr *SchemaWildcardNs

type X_xmlSchemaWildcard struct {
	Type            SchemaTypeType
	Id              *Char
	Annot           SchemaAnnotPtr
	Node            NodePtr
	MinOccurs       c.Int
	MaxOccurs       c.Int
	ProcessContents c.Int
	Any             c.Int
	NsSet           SchemaWildcardNsPtr
	NegNsSet        SchemaWildcardNsPtr
	Flags           c.Int
}
type SchemaWildcard X_xmlSchemaWildcard
type SchemaWildcardPtr *SchemaWildcard

type X_xmlSchemaAttributeGroup struct {
	Type              SchemaTypeType
	Next              *X_xmlSchemaAttribute
	Name              *Char
	Id                *Char
	Ref               *Char
	RefNs             *Char
	Annot             SchemaAnnotPtr
	Attributes        SchemaAttributePtr
	Node              NodePtr
	Flags             c.Int
	AttributeWildcard SchemaWildcardPtr
	RefPrefix         *Char
	RefItem           SchemaAttributeGroupPtr
	TargetNamespace   *Char
	AttrUses          c.Pointer
}
type SchemaAttributeGroup X_xmlSchemaAttributeGroup
type SchemaAttributeGroupPtr *SchemaAttributeGroup

type X_xmlSchemaTypeLink struct {
	Next *X_xmlSchemaTypeLink
	Type SchemaTypePtr
}
type SchemaTypeLink X_xmlSchemaTypeLink
type SchemaTypeLinkPtr *SchemaTypeLink

type X_xmlSchemaFacetLink struct {
	Next  *X_xmlSchemaFacetLink
	Facet SchemaFacetPtr
}
type SchemaFacetLink X_xmlSchemaFacetLink
type SchemaFacetLinkPtr *SchemaFacetLink

type X_xmlSchemaElement struct {
	Type            SchemaTypeType
	Next            *X_xmlSchemaType
	Name            *Char
	Id              *Char
	Ref             *Char
	RefNs           *Char
	Annot           SchemaAnnotPtr
	Subtypes        SchemaTypePtr
	Attributes      SchemaAttributePtr
	Node            NodePtr
	MinOccurs       c.Int
	MaxOccurs       c.Int
	Flags           c.Int
	TargetNamespace *Char
	NamedType       *Char
	NamedTypeNs     *Char
	SubstGroup      *Char
	SubstGroupNs    *Char
	Scope           *Char
	Value           *Char
	RefDecl         *X_xmlSchemaElement
	ContModel       RegexpPtr
	ContentType     SchemaContentType
	RefPrefix       *Char
	DefVal          SchemaValPtr
	Idcs            c.Pointer
}
type SchemaElement X_xmlSchemaElement
type SchemaElementPtr *SchemaElement

type X_xmlSchemaNotation struct {
	Type            SchemaTypeType
	Name            *Char
	Annot           SchemaAnnotPtr
	Identifier      *Char
	TargetNamespace *Char
}
type SchemaNotation X_xmlSchemaNotation
type SchemaNotationPtr *SchemaNotation

/**
 * _xmlSchema:
 *
 * A Schemas definition
 */

type X_xmlSchema struct {
	Name            *Char
	TargetNamespace *Char
	Version         *Char
	Id              *Char
	Doc             DocPtr
	Annot           SchemaAnnotPtr
	Flags           c.Int
	TypeDecl        HashTablePtr
	AttrDecl        HashTablePtr
	AttrgrpDecl     HashTablePtr
	ElemDecl        HashTablePtr
	NotaDecl        HashTablePtr
	SchemasImports  HashTablePtr
	X_private       c.Pointer
	GroupDecl       HashTablePtr
	Dict            DictPtr
	Includes        c.Pointer
	Preserve        c.Int
	Counter         c.Int
	IdcDef          HashTablePtr
	Volatiles       c.Pointer
}

//go:linkname SchemaFreeType C.xmlSchemaFreeType
func SchemaFreeType(type_ SchemaTypePtr)

//go:linkname SchemaFreeWildcard C.xmlSchemaFreeWildcard
func SchemaFreeWildcard(wildcard SchemaWildcardPtr)
