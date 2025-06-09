package libxml2

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type SchemaValidError c.Int

const (
	SCHEMAS_ERR_OK             SchemaValidError = 0
	SCHEMAS_ERR_NOROOT         SchemaValidError = 1
	SCHEMAS_ERR_UNDECLAREDELEM SchemaValidError = 2
	SCHEMAS_ERR_NOTTOPLEVEL    SchemaValidError = 3
	SCHEMAS_ERR_MISSING        SchemaValidError = 4
	SCHEMAS_ERR_WRONGELEM      SchemaValidError = 5
	SCHEMAS_ERR_NOTYPE         SchemaValidError = 6
	SCHEMAS_ERR_NOROLLBACK     SchemaValidError = 7
	SCHEMAS_ERR_ISABSTRACT     SchemaValidError = 8
	SCHEMAS_ERR_NOTEMPTY       SchemaValidError = 9
	SCHEMAS_ERR_ELEMCONT       SchemaValidError = 10
	SCHEMAS_ERR_HAVEDEFAULT    SchemaValidError = 11
	SCHEMAS_ERR_NOTNILLABLE    SchemaValidError = 12
	SCHEMAS_ERR_EXTRACONTENT   SchemaValidError = 13
	SCHEMAS_ERR_INVALIDATTR    SchemaValidError = 14
	SCHEMAS_ERR_INVALIDELEM    SchemaValidError = 15
	SCHEMAS_ERR_NOTDETERMINIST SchemaValidError = 16
	SCHEMAS_ERR_CONSTRUCT      SchemaValidError = 17
	SCHEMAS_ERR_INTERNAL       SchemaValidError = 18
	SCHEMAS_ERR_NOTSIMPLE      SchemaValidError = 19
	SCHEMAS_ERR_ATTRUNKNOWN    SchemaValidError = 20
	SCHEMAS_ERR_ATTRINVALID    SchemaValidError = 21
	SCHEMAS_ERR_VALUE          SchemaValidError = 22
	SCHEMAS_ERR_FACET          SchemaValidError = 23
	SCHEMAS_ERR_               SchemaValidError = 24
	SCHEMAS_ERR_XXX            SchemaValidError = 25
)

type SchemaValidOption c.Int

const SCHEMA_VAL_VC_I_CREATE SchemaValidOption = 1

type Schema X_xmlSchema
type SchemaPtr *Schema

// llgo:type C
type SchemaValidityErrorFunc func(__llgo_arg_0 c.Pointer, __llgo_arg_1 *c.Char, __llgo_va_list ...interface{})

// llgo:type C
type SchemaValidityWarningFunc func(__llgo_arg_0 c.Pointer, __llgo_arg_1 *c.Char, __llgo_va_list ...interface{})

type X_xmlSchemaParserCtxt struct {
	Unused [8]uint8
}
type SchemaParserCtxt X_xmlSchemaParserCtxt
type SchemaParserCtxtPtr *SchemaParserCtxt

type X_xmlSchemaValidCtxt struct {
	Unused [8]uint8
}
type SchemaValidCtxt X_xmlSchemaValidCtxt
type SchemaValidCtxtPtr *SchemaValidCtxt

// llgo:type C
type SchemaValidityLocatorFunc func(c.Pointer, **c.Char, *c.Ulong) c.Int

/*
 * Interfaces for parsing.
 */
//go:linkname SchemaNewParserCtxt C.xmlSchemaNewParserCtxt
func SchemaNewParserCtxt(URL *c.Char) SchemaParserCtxtPtr

//go:linkname SchemaNewMemParserCtxt C.xmlSchemaNewMemParserCtxt
func SchemaNewMemParserCtxt(buffer *c.Char, size c.Int) SchemaParserCtxtPtr

//go:linkname SchemaNewDocParserCtxt C.xmlSchemaNewDocParserCtxt
func SchemaNewDocParserCtxt(doc DocPtr) SchemaParserCtxtPtr

//go:linkname SchemaFreeParserCtxt C.xmlSchemaFreeParserCtxt
func SchemaFreeParserCtxt(ctxt SchemaParserCtxtPtr)

//go:linkname SchemaSetParserErrors C.xmlSchemaSetParserErrors
func SchemaSetParserErrors(ctxt SchemaParserCtxtPtr, err SchemaValidityErrorFunc, warn SchemaValidityWarningFunc, ctx c.Pointer)

//go:linkname SchemaSetParserStructuredErrors C.xmlSchemaSetParserStructuredErrors
func SchemaSetParserStructuredErrors(ctxt SchemaParserCtxtPtr, serror StructuredErrorFunc, ctx c.Pointer)

//go:linkname SchemaGetParserErrors C.xmlSchemaGetParserErrors
func SchemaGetParserErrors(ctxt SchemaParserCtxtPtr, err SchemaValidityErrorFunc, warn SchemaValidityWarningFunc, ctx *c.Pointer) c.Int

//go:linkname SchemaIsValid C.xmlSchemaIsValid
func SchemaIsValid(ctxt SchemaValidCtxtPtr) c.Int

//go:linkname SchemaParse C.xmlSchemaParse
func SchemaParse(ctxt SchemaParserCtxtPtr) SchemaPtr

//go:linkname SchemaFree C.xmlSchemaFree
func SchemaFree(schema SchemaPtr)

//go:linkname SchemaDump C.xmlSchemaDump
func SchemaDump(output *c.FILE, schema SchemaPtr)

/*
 * Interfaces for validating
 */
//go:linkname SchemaSetValidErrors C.xmlSchemaSetValidErrors
func SchemaSetValidErrors(ctxt SchemaValidCtxtPtr, err SchemaValidityErrorFunc, warn SchemaValidityWarningFunc, ctx c.Pointer)

//go:linkname SchemaSetValidStructuredErrors C.xmlSchemaSetValidStructuredErrors
func SchemaSetValidStructuredErrors(ctxt SchemaValidCtxtPtr, serror StructuredErrorFunc, ctx c.Pointer)

//go:linkname SchemaGetValidErrors C.xmlSchemaGetValidErrors
func SchemaGetValidErrors(ctxt SchemaValidCtxtPtr, err SchemaValidityErrorFunc, warn SchemaValidityWarningFunc, ctx *c.Pointer) c.Int

//go:linkname SchemaSetValidOptions C.xmlSchemaSetValidOptions
func SchemaSetValidOptions(ctxt SchemaValidCtxtPtr, options c.Int) c.Int

//go:linkname SchemaValidateSetFilename C.xmlSchemaValidateSetFilename
func SchemaValidateSetFilename(vctxt SchemaValidCtxtPtr, filename *c.Char)

//go:linkname SchemaValidCtxtGetOptions C.xmlSchemaValidCtxtGetOptions
func SchemaValidCtxtGetOptions(ctxt SchemaValidCtxtPtr) c.Int

//go:linkname SchemaNewValidCtxt C.xmlSchemaNewValidCtxt
func SchemaNewValidCtxt(schema SchemaPtr) SchemaValidCtxtPtr

//go:linkname SchemaFreeValidCtxt C.xmlSchemaFreeValidCtxt
func SchemaFreeValidCtxt(ctxt SchemaValidCtxtPtr)

//go:linkname SchemaValidateDoc C.xmlSchemaValidateDoc
func SchemaValidateDoc(ctxt SchemaValidCtxtPtr, instance DocPtr) c.Int

//go:linkname SchemaValidateOneElement C.xmlSchemaValidateOneElement
func SchemaValidateOneElement(ctxt SchemaValidCtxtPtr, elem NodePtr) c.Int

//go:linkname SchemaValidateStream C.xmlSchemaValidateStream
func SchemaValidateStream(ctxt SchemaValidCtxtPtr, input ParserInputBufferPtr, enc CharEncoding, sax SAXHandlerPtr, user_data c.Pointer) c.Int

//go:linkname SchemaValidateFile C.xmlSchemaValidateFile
func SchemaValidateFile(ctxt SchemaValidCtxtPtr, filename *c.Char, options c.Int) c.Int

//go:linkname SchemaValidCtxtGetParserCtxt C.xmlSchemaValidCtxtGetParserCtxt
func SchemaValidCtxtGetParserCtxt(ctxt SchemaValidCtxtPtr) ParserCtxtPtr

type X_xmlSchemaSAXPlug struct {
	Unused [8]uint8
}
type SchemaSAXPlugStruct X_xmlSchemaSAXPlug
type SchemaSAXPlugPtr *SchemaSAXPlugStruct

//go:linkname SchemaSAXPlug C.xmlSchemaSAXPlug
func SchemaSAXPlug(ctxt SchemaValidCtxtPtr, sax *SAXHandlerPtr, user_data *c.Pointer) SchemaSAXPlugPtr

//go:linkname SchemaSAXUnplug C.xmlSchemaSAXUnplug
func SchemaSAXUnplug(plug SchemaSAXPlugPtr) c.Int

//go:linkname SchemaValidateSetLocator C.xmlSchemaValidateSetLocator
func SchemaValidateSetLocator(vctxt SchemaValidCtxtPtr, f SchemaValidityLocatorFunc, ctxt c.Pointer)
