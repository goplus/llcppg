package libxml2

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type X_xmlRelaxNG struct {
	Unused [8]uint8
}
type RelaxNG X_xmlRelaxNG
type RelaxNGPtr *RelaxNG

// llgo:type C
type RelaxNGValidityErrorFunc func(__llgo_arg_0 c.Pointer, __llgo_arg_1 *c.Char, __llgo_va_list ...interface{})

// llgo:type C
type RelaxNGValidityWarningFunc func(__llgo_arg_0 c.Pointer, __llgo_arg_1 *c.Char, __llgo_va_list ...interface{})

type X_xmlRelaxNGParserCtxt struct {
	Unused [8]uint8
}
type RelaxNGParserCtxt X_xmlRelaxNGParserCtxt
type RelaxNGParserCtxtPtr *RelaxNGParserCtxt

type X_xmlRelaxNGValidCtxt struct {
	Unused [8]uint8
}
type RelaxNGValidCtxt X_xmlRelaxNGValidCtxt
type RelaxNGValidCtxtPtr *RelaxNGValidCtxt
type RelaxNGValidErr c.Int

const (
	RELAXNG_OK               RelaxNGValidErr = 0
	RELAXNG_ERR_MEMORY       RelaxNGValidErr = 1
	RELAXNG_ERR_TYPE         RelaxNGValidErr = 2
	RELAXNG_ERR_TYPEVAL      RelaxNGValidErr = 3
	RELAXNG_ERR_DUPID        RelaxNGValidErr = 4
	RELAXNG_ERR_TYPECMP      RelaxNGValidErr = 5
	RELAXNG_ERR_NOSTATE      RelaxNGValidErr = 6
	RELAXNG_ERR_NODEFINE     RelaxNGValidErr = 7
	RELAXNG_ERR_LISTEXTRA    RelaxNGValidErr = 8
	RELAXNG_ERR_LISTEMPTY    RelaxNGValidErr = 9
	RELAXNG_ERR_INTERNODATA  RelaxNGValidErr = 10
	RELAXNG_ERR_INTERSEQ     RelaxNGValidErr = 11
	RELAXNG_ERR_INTEREXTRA   RelaxNGValidErr = 12
	RELAXNG_ERR_ELEMNAME     RelaxNGValidErr = 13
	RELAXNG_ERR_ATTRNAME     RelaxNGValidErr = 14
	RELAXNG_ERR_ELEMNONS     RelaxNGValidErr = 15
	RELAXNG_ERR_ATTRNONS     RelaxNGValidErr = 16
	RELAXNG_ERR_ELEMWRONGNS  RelaxNGValidErr = 17
	RELAXNG_ERR_ATTRWRONGNS  RelaxNGValidErr = 18
	RELAXNG_ERR_ELEMEXTRANS  RelaxNGValidErr = 19
	RELAXNG_ERR_ATTREXTRANS  RelaxNGValidErr = 20
	RELAXNG_ERR_ELEMNOTEMPTY RelaxNGValidErr = 21
	RELAXNG_ERR_NOELEM       RelaxNGValidErr = 22
	RELAXNG_ERR_NOTELEM      RelaxNGValidErr = 23
	RELAXNG_ERR_ATTRVALID    RelaxNGValidErr = 24
	RELAXNG_ERR_CONTENTVALID RelaxNGValidErr = 25
	RELAXNG_ERR_EXTRACONTENT RelaxNGValidErr = 26
	RELAXNG_ERR_INVALIDATTR  RelaxNGValidErr = 27
	RELAXNG_ERR_DATAELEM     RelaxNGValidErr = 28
	RELAXNG_ERR_VALELEM      RelaxNGValidErr = 29
	RELAXNG_ERR_LISTELEM     RelaxNGValidErr = 30
	RELAXNG_ERR_DATATYPE     RelaxNGValidErr = 31
	RELAXNG_ERR_VALUE        RelaxNGValidErr = 32
	RELAXNG_ERR_LIST         RelaxNGValidErr = 33
	RELAXNG_ERR_NOGRAMMAR    RelaxNGValidErr = 34
	RELAXNG_ERR_EXTRADATA    RelaxNGValidErr = 35
	RELAXNG_ERR_LACKDATA     RelaxNGValidErr = 36
	RELAXNG_ERR_INTERNAL     RelaxNGValidErr = 37
	RELAXNG_ERR_ELEMWRONG    RelaxNGValidErr = 38
	RELAXNG_ERR_TEXTWRONG    RelaxNGValidErr = 39
)

type RelaxNGParserFlag c.Int

const (
	RELAXNGP_NONE     RelaxNGParserFlag = 0
	RELAXNGP_FREE_DOC RelaxNGParserFlag = 1
	RELAXNGP_CRNG     RelaxNGParserFlag = 2
)

//go:linkname RelaxNGInitTypes C.xmlRelaxNGInitTypes
func RelaxNGInitTypes() c.Int

//go:linkname RelaxNGCleanupTypes C.xmlRelaxNGCleanupTypes
func RelaxNGCleanupTypes()

/*
 * Interfaces for parsing.
 */
//go:linkname RelaxNGNewParserCtxt C.xmlRelaxNGNewParserCtxt
func RelaxNGNewParserCtxt(URL *c.Char) RelaxNGParserCtxtPtr

//go:linkname RelaxNGNewMemParserCtxt C.xmlRelaxNGNewMemParserCtxt
func RelaxNGNewMemParserCtxt(buffer *c.Char, size c.Int) RelaxNGParserCtxtPtr

//go:linkname RelaxNGNewDocParserCtxt C.xmlRelaxNGNewDocParserCtxt
func RelaxNGNewDocParserCtxt(doc DocPtr) RelaxNGParserCtxtPtr

//go:linkname RelaxParserSetFlag C.xmlRelaxParserSetFlag
func RelaxParserSetFlag(ctxt RelaxNGParserCtxtPtr, flag c.Int) c.Int

//go:linkname RelaxNGFreeParserCtxt C.xmlRelaxNGFreeParserCtxt
func RelaxNGFreeParserCtxt(ctxt RelaxNGParserCtxtPtr)

//go:linkname RelaxNGSetParserErrors C.xmlRelaxNGSetParserErrors
func RelaxNGSetParserErrors(ctxt RelaxNGParserCtxtPtr, err RelaxNGValidityErrorFunc, warn RelaxNGValidityWarningFunc, ctx c.Pointer)

//go:linkname RelaxNGGetParserErrors C.xmlRelaxNGGetParserErrors
func RelaxNGGetParserErrors(ctxt RelaxNGParserCtxtPtr, err RelaxNGValidityErrorFunc, warn RelaxNGValidityWarningFunc, ctx *c.Pointer) c.Int

//go:linkname RelaxNGSetParserStructuredErrors C.xmlRelaxNGSetParserStructuredErrors
func RelaxNGSetParserStructuredErrors(ctxt RelaxNGParserCtxtPtr, serror StructuredErrorFunc, ctx c.Pointer)

//go:linkname RelaxNGParse C.xmlRelaxNGParse
func RelaxNGParse(ctxt RelaxNGParserCtxtPtr) RelaxNGPtr

//go:linkname RelaxNGFree C.xmlRelaxNGFree
func RelaxNGFree(schema RelaxNGPtr)

//go:linkname RelaxNGDump C.xmlRelaxNGDump
func RelaxNGDump(output *c.FILE, schema RelaxNGPtr)

//go:linkname RelaxNGDumpTree C.xmlRelaxNGDumpTree
func RelaxNGDumpTree(output *c.FILE, schema RelaxNGPtr)

/*
 * Interfaces for validating
 */
//go:linkname RelaxNGSetValidErrors C.xmlRelaxNGSetValidErrors
func RelaxNGSetValidErrors(ctxt RelaxNGValidCtxtPtr, err RelaxNGValidityErrorFunc, warn RelaxNGValidityWarningFunc, ctx c.Pointer)

//go:linkname RelaxNGGetValidErrors C.xmlRelaxNGGetValidErrors
func RelaxNGGetValidErrors(ctxt RelaxNGValidCtxtPtr, err RelaxNGValidityErrorFunc, warn RelaxNGValidityWarningFunc, ctx *c.Pointer) c.Int

//go:linkname RelaxNGSetValidStructuredErrors C.xmlRelaxNGSetValidStructuredErrors
func RelaxNGSetValidStructuredErrors(ctxt RelaxNGValidCtxtPtr, serror StructuredErrorFunc, ctx c.Pointer)

//go:linkname RelaxNGNewValidCtxt C.xmlRelaxNGNewValidCtxt
func RelaxNGNewValidCtxt(schema RelaxNGPtr) RelaxNGValidCtxtPtr

//go:linkname RelaxNGFreeValidCtxt C.xmlRelaxNGFreeValidCtxt
func RelaxNGFreeValidCtxt(ctxt RelaxNGValidCtxtPtr)

//go:linkname RelaxNGValidateDoc C.xmlRelaxNGValidateDoc
func RelaxNGValidateDoc(ctxt RelaxNGValidCtxtPtr, doc DocPtr) c.Int

/*
 * Interfaces for progressive validation when possible
 */
//go:linkname RelaxNGValidatePushElement C.xmlRelaxNGValidatePushElement
func RelaxNGValidatePushElement(ctxt RelaxNGValidCtxtPtr, doc DocPtr, elem NodePtr) c.Int

//go:linkname RelaxNGValidatePushCData C.xmlRelaxNGValidatePushCData
func RelaxNGValidatePushCData(ctxt RelaxNGValidCtxtPtr, data *Char, len c.Int) c.Int

//go:linkname RelaxNGValidatePopElement C.xmlRelaxNGValidatePopElement
func RelaxNGValidatePopElement(ctxt RelaxNGValidCtxtPtr, doc DocPtr, elem NodePtr) c.Int

//go:linkname RelaxNGValidateFullElement C.xmlRelaxNGValidateFullElement
func RelaxNGValidateFullElement(ctxt RelaxNGValidCtxtPtr, doc DocPtr, elem NodePtr) c.Int
