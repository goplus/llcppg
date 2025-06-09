package libxml2

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type Dict X_xmlDict
type DictPtr *Dict

/*
 * Initializer
 */
//go:linkname InitializeDict C.xmlInitializeDict
func InitializeDict() c.Int

/*
 * Constructor and destructor.
 */
//go:linkname DictCreate C.xmlDictCreate
func DictCreate() DictPtr

//go:linkname DictSetLimit C.xmlDictSetLimit
func DictSetLimit(dict DictPtr, limit c.SizeT) c.SizeT

//go:linkname DictGetUsage C.xmlDictGetUsage
func DictGetUsage(dict DictPtr) c.SizeT

//go:linkname DictCreateSub C.xmlDictCreateSub
func DictCreateSub(sub DictPtr) DictPtr

//go:linkname DictReference C.xmlDictReference
func DictReference(dict DictPtr) c.Int

//go:linkname DictFree C.xmlDictFree
func DictFree(dict DictPtr)

/*
 * Lookup of entry in the dictionary.
 */
//go:linkname DictLookup C.xmlDictLookup
func DictLookup(dict DictPtr, name *Char, len c.Int) *Char

//go:linkname DictExists C.xmlDictExists
func DictExists(dict DictPtr, name *Char, len c.Int) *Char

//go:linkname DictQLookup C.xmlDictQLookup
func DictQLookup(dict DictPtr, prefix *Char, name *Char) *Char

//go:linkname DictOwns C.xmlDictOwns
func DictOwns(dict DictPtr, str *Char) c.Int

//go:linkname DictSize C.xmlDictSize
func DictSize(dict DictPtr) c.Int

/*
 * Cleanup function
 */
//go:linkname DictCleanup C.xmlDictCleanup
func DictCleanup()
