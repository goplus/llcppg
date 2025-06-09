package libxml2

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type SaveOption c.Int

const (
	SAVE_FORMAT   SaveOption = 1
	SAVE_NO_DECL  SaveOption = 2
	SAVE_NO_EMPTY SaveOption = 4
	SAVE_NO_XHTML SaveOption = 8
	SAVE_XHTML    SaveOption = 16
	SAVE_AS_XML   SaveOption = 32
	SAVE_AS_HTML  SaveOption = 64
	SAVE_WSNONSIG SaveOption = 128
)

type X_xmlSaveCtxt struct {
	Unused [8]uint8
}
type SaveCtxt X_xmlSaveCtxt
type SaveCtxtPtr *SaveCtxt

//go:linkname SaveToFd C.xmlSaveToFd
func SaveToFd(fd c.Int, encoding *c.Char, options c.Int) SaveCtxtPtr

//go:linkname SaveToFilename C.xmlSaveToFilename
func SaveToFilename(filename *c.Char, encoding *c.Char, options c.Int) SaveCtxtPtr

//go:linkname SaveToBuffer C.xmlSaveToBuffer
func SaveToBuffer(buffer BufferPtr, encoding *c.Char, options c.Int) SaveCtxtPtr

//go:linkname SaveToIO C.xmlSaveToIO
func SaveToIO(iowrite OutputWriteCallback, ioclose OutputCloseCallback, ioctx c.Pointer, encoding *c.Char, options c.Int) SaveCtxtPtr

//go:linkname SaveDoc C.xmlSaveDoc
func SaveDoc(ctxt SaveCtxtPtr, doc DocPtr) c.Long

//go:linkname SaveTree C.xmlSaveTree
func SaveTree(ctxt SaveCtxtPtr, node NodePtr) c.Long

//go:linkname SaveFlush C.xmlSaveFlush
func SaveFlush(ctxt SaveCtxtPtr) c.Int

//go:linkname SaveClose C.xmlSaveClose
func SaveClose(ctxt SaveCtxtPtr) c.Int

//go:linkname SaveFinish C.xmlSaveFinish
func SaveFinish(ctxt SaveCtxtPtr) c.Int

//go:linkname SaveSetEscape C.xmlSaveSetEscape
func SaveSetEscape(ctxt SaveCtxtPtr, escape CharEncodingOutputFunc) c.Int

//go:linkname SaveSetAttrEscape C.xmlSaveSetAttrEscape
func SaveSetAttrEscape(ctxt SaveCtxtPtr, escape CharEncodingOutputFunc) c.Int

//go:linkname ThrDefIndentTreeOutput C.xmlThrDefIndentTreeOutput
func ThrDefIndentTreeOutput(v c.Int) c.Int

//go:linkname ThrDefTreeIndentString C.xmlThrDefTreeIndentString
func ThrDefTreeIndentString(v *c.Char) *c.Char

//go:linkname ThrDefSaveNoEmptyTags C.xmlThrDefSaveNoEmptyTags
func ThrDefSaveNoEmptyTags(v c.Int) c.Int
