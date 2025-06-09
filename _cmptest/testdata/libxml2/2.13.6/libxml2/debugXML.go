package libxml2

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

/*
 * The standard Dump routines.
 */
//go:linkname DebugDumpString C.xmlDebugDumpString
func DebugDumpString(output *c.FILE, str *Char)

//go:linkname DebugDumpAttr C.xmlDebugDumpAttr
func DebugDumpAttr(output *c.FILE, attr AttrPtr, depth c.Int)

//go:linkname DebugDumpAttrList C.xmlDebugDumpAttrList
func DebugDumpAttrList(output *c.FILE, attr AttrPtr, depth c.Int)

//go:linkname DebugDumpOneNode C.xmlDebugDumpOneNode
func DebugDumpOneNode(output *c.FILE, node NodePtr, depth c.Int)

//go:linkname DebugDumpNode C.xmlDebugDumpNode
func DebugDumpNode(output *c.FILE, node NodePtr, depth c.Int)

//go:linkname DebugDumpNodeList C.xmlDebugDumpNodeList
func DebugDumpNodeList(output *c.FILE, node NodePtr, depth c.Int)

//go:linkname DebugDumpDocumentHead C.xmlDebugDumpDocumentHead
func DebugDumpDocumentHead(output *c.FILE, doc DocPtr)

//go:linkname DebugDumpDocument C.xmlDebugDumpDocument
func DebugDumpDocument(output *c.FILE, doc DocPtr)

//go:linkname DebugDumpDTD C.xmlDebugDumpDTD
func DebugDumpDTD(output *c.FILE, dtd DtdPtr)

//go:linkname DebugDumpEntities C.xmlDebugDumpEntities
func DebugDumpEntities(output *c.FILE, doc DocPtr)

/****************************************************************
 *								*
 *			Checking routines			*
 *								*
 ****************************************************************/
//go:linkname DebugCheckDocument C.xmlDebugCheckDocument
func DebugCheckDocument(output *c.FILE, doc DocPtr) c.Int

/****************************************************************
 *								*
 *			XML shell helpers			*
 *								*
 ****************************************************************/
//go:linkname LsOneNode C.xmlLsOneNode
func LsOneNode(output *c.FILE, node NodePtr)

//go:linkname LsCountNode C.xmlLsCountNode
func LsCountNode(node NodePtr) c.Int

//go:linkname BoolToText C.xmlBoolToText
func BoolToText(boolval c.Int) *c.Char

// llgo:type C
type ShellReadlineFunc func(*c.Char) *c.Char

type X_xmlShellCtxt struct {
	Filename *c.Char
	Doc      DocPtr
	Node     NodePtr
	Pctxt    XPathContextPtr
	Loaded   c.Int
	Output   *c.FILE
	Input    c.Pointer
}
type ShellCtxt X_xmlShellCtxt
type ShellCtxtPtr *ShellCtxt

// llgo:type C
type ShellCmd func(ShellCtxtPtr, *c.Char, NodePtr, NodePtr) c.Int

//go:linkname ShellPrintXPathError C.xmlShellPrintXPathError
func ShellPrintXPathError(errorType c.Int, arg *c.Char)

//go:linkname ShellPrintXPathResult C.xmlShellPrintXPathResult
func ShellPrintXPathResult(list XPathObjectPtr)

//go:linkname ShellList C.xmlShellList
func ShellList(ctxt ShellCtxtPtr, arg *c.Char, node NodePtr, node2 NodePtr) c.Int

//go:linkname ShellBase C.xmlShellBase
func ShellBase(ctxt ShellCtxtPtr, arg *c.Char, node NodePtr, node2 NodePtr) c.Int

//go:linkname ShellDir C.xmlShellDir
func ShellDir(ctxt ShellCtxtPtr, arg *c.Char, node NodePtr, node2 NodePtr) c.Int

//go:linkname ShellLoad C.xmlShellLoad
func ShellLoad(ctxt ShellCtxtPtr, filename *c.Char, node NodePtr, node2 NodePtr) c.Int

//go:linkname ShellPrintNode C.xmlShellPrintNode
func ShellPrintNode(node NodePtr)

//go:linkname ShellCat C.xmlShellCat
func ShellCat(ctxt ShellCtxtPtr, arg *c.Char, node NodePtr, node2 NodePtr) c.Int

//go:linkname ShellWrite C.xmlShellWrite
func ShellWrite(ctxt ShellCtxtPtr, filename *c.Char, node NodePtr, node2 NodePtr) c.Int

//go:linkname ShellSave C.xmlShellSave
func ShellSave(ctxt ShellCtxtPtr, filename *c.Char, node NodePtr, node2 NodePtr) c.Int

//go:linkname ShellValidate C.xmlShellValidate
func ShellValidate(ctxt ShellCtxtPtr, dtd *c.Char, node NodePtr, node2 NodePtr) c.Int

//go:linkname ShellDu C.xmlShellDu
func ShellDu(ctxt ShellCtxtPtr, arg *c.Char, tree NodePtr, node2 NodePtr) c.Int

//go:linkname ShellPwd C.xmlShellPwd
func ShellPwd(ctxt ShellCtxtPtr, buffer *c.Char, node NodePtr, node2 NodePtr) c.Int

/*
 * The Shell interface.
 */
//go:linkname Shell C.xmlShell
func Shell(doc DocPtr, filename *c.Char, input ShellReadlineFunc, output *c.FILE)
