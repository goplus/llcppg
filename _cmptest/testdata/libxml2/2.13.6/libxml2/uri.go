package libxml2

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type X_xmlURI struct {
	Scheme    *c.Char
	Opaque    *c.Char
	Authority *c.Char
	Server    *c.Char
	User      *c.Char
	Port      c.Int
	Path      *c.Char
	Query     *c.Char
	Fragment  *c.Char
	Cleanup   c.Int
	QueryRaw  *c.Char
}
type URI X_xmlURI
type URIPtr *URI

/*
 * This function is in tree.h:
 * xmlChar *	xmlNodeGetBase	(xmlDocPtr doc,
 *                               xmlNodePtr cur);
 */
//go:linkname CreateURI C.xmlCreateURI
func CreateURI() URIPtr

// llgo:link (*Char).BuildURISafe C.xmlBuildURISafe
func (recv_ *Char) BuildURISafe(base *Char, out **Char) c.Int {
	return 0
}

// llgo:link (*Char).BuildURI C.xmlBuildURI
func (recv_ *Char) BuildURI(base *Char) *Char {
	return nil
}

// llgo:link (*Char).BuildRelativeURISafe C.xmlBuildRelativeURISafe
func (recv_ *Char) BuildRelativeURISafe(base *Char, out **Char) c.Int {
	return 0
}

// llgo:link (*Char).BuildRelativeURI C.xmlBuildRelativeURI
func (recv_ *Char) BuildRelativeURI(base *Char) *Char {
	return nil
}

//go:linkname ParseURI C.xmlParseURI
func ParseURI(str *c.Char) URIPtr

//go:linkname ParseURISafe C.xmlParseURISafe
func ParseURISafe(str *c.Char, uri *URIPtr) c.Int

//go:linkname ParseURIRaw C.xmlParseURIRaw
func ParseURIRaw(str *c.Char, raw c.Int) URIPtr

//go:linkname ParseURIReference C.xmlParseURIReference
func ParseURIReference(uri URIPtr, str *c.Char) c.Int

//go:linkname SaveUri C.xmlSaveUri
func SaveUri(uri URIPtr) *Char

//go:linkname PrintURI C.xmlPrintURI
func PrintURI(stream *c.FILE, uri URIPtr)

// llgo:link (*Char).URIEscapeStr C.xmlURIEscapeStr
func (recv_ *Char) URIEscapeStr(list *Char) *Char {
	return nil
}

//go:linkname URIUnescapeString C.xmlURIUnescapeString
func URIUnescapeString(str *c.Char, len c.Int, target *c.Char) *c.Char

//go:linkname NormalizeURIPath C.xmlNormalizeURIPath
func NormalizeURIPath(path *c.Char) c.Int

// llgo:link (*Char).URIEscape C.xmlURIEscape
func (recv_ *Char) URIEscape() *Char {
	return nil
}

//go:linkname FreeURI C.xmlFreeURI
func FreeURI(uri URIPtr)

// llgo:link (*Char).CanonicPath C.xmlCanonicPath
func (recv_ *Char) CanonicPath() *Char {
	return nil
}

// llgo:link (*Char).PathToURI C.xmlPathToURI
func (recv_ *Char) PathToURI() *Char {
	return nil
}
