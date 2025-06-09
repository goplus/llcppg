package libxml2

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

// llgo:type C
type FreeFunc func(c.Pointer)

// llgo:type C
type MallocFunc func(c.SizeT) c.Pointer

// llgo:type C
type ReallocFunc func(c.Pointer, c.SizeT) c.Pointer

// llgo:type C
type StrdupFunc func(*c.Char) *c.Char

/*
 * The way to overload the existing functions.
 * The xmlGc function have an extra entry for atomic block
 * allocations useful for garbage collected memory allocators
 */
//go:linkname MemSetup C.xmlMemSetup
func MemSetup(freeFunc FreeFunc, mallocFunc MallocFunc, reallocFunc ReallocFunc, strdupFunc StrdupFunc) c.Int

//go:linkname MemGet C.xmlMemGet
func MemGet(freeFunc FreeFunc, mallocFunc MallocFunc, reallocFunc ReallocFunc, strdupFunc StrdupFunc) c.Int

//go:linkname GcMemSetup C.xmlGcMemSetup
func GcMemSetup(freeFunc FreeFunc, mallocFunc MallocFunc, mallocAtomicFunc MallocFunc, reallocFunc ReallocFunc, strdupFunc StrdupFunc) c.Int

//go:linkname GcMemGet C.xmlGcMemGet
func GcMemGet(freeFunc FreeFunc, mallocFunc MallocFunc, mallocAtomicFunc MallocFunc, reallocFunc ReallocFunc, strdupFunc StrdupFunc) c.Int

/*
 * Initialization of the memory layer.
 */
//go:linkname InitMemory C.xmlInitMemory
func InitMemory() c.Int

/*
 * Cleanup of the memory layer.
 */
//go:linkname CleanupMemory C.xmlCleanupMemory
func CleanupMemory()

/*
 * These are specific to the XML debug memory wrapper.
 */
//go:linkname MemSize C.xmlMemSize
func MemSize(ptr c.Pointer) c.SizeT

//go:linkname MemUsed C.xmlMemUsed
func MemUsed() c.Int

//go:linkname MemBlocks C.xmlMemBlocks
func MemBlocks() c.Int

//go:linkname MemDisplay C.xmlMemDisplay
func MemDisplay(fp *c.FILE)

//go:linkname MemDisplayLast C.xmlMemDisplayLast
func MemDisplayLast(fp *c.FILE, nbBytes c.Long)

//go:linkname MemShow C.xmlMemShow
func MemShow(fp *c.FILE, nr c.Int)

//go:linkname MemoryDump C.xmlMemoryDump
func MemoryDump()

//go:linkname MemMalloc C.xmlMemMalloc
func MemMalloc(size c.SizeT) c.Pointer

//go:linkname MemRealloc C.xmlMemRealloc
func MemRealloc(ptr c.Pointer, size c.SizeT) c.Pointer

//go:linkname MemFree C.xmlMemFree
func MemFree(ptr c.Pointer)

//go:linkname MemoryStrdup C.xmlMemoryStrdup
func MemoryStrdup(str *c.Char) *c.Char

//go:linkname MallocLoc C.xmlMallocLoc
func MallocLoc(size c.SizeT, file *c.Char, line c.Int) c.Pointer

//go:linkname ReallocLoc C.xmlReallocLoc
func ReallocLoc(ptr c.Pointer, size c.SizeT, file *c.Char, line c.Int) c.Pointer

//go:linkname MallocAtomicLoc C.xmlMallocAtomicLoc
func MallocAtomicLoc(size c.SizeT, file *c.Char, line c.Int) c.Pointer

//go:linkname MemStrdupLoc C.xmlMemStrdupLoc
func MemStrdupLoc(str *c.Char, file *c.Char, line c.Int) *c.Char
