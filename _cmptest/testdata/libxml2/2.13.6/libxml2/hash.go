package libxml2

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type X_xmlHashTable struct {
	Unused [8]uint8
}
type HashTable X_xmlHashTable
type HashTablePtr *HashTable

// llgo:type C
type HashDeallocator func(c.Pointer, *Char)

// llgo:type C
type HashCopier func(c.Pointer, *Char) c.Pointer

// llgo:type C
type HashScanner func(c.Pointer, c.Pointer, *Char)

// llgo:type C
type HashScannerFull func(c.Pointer, c.Pointer, *Char, *Char, *Char)

/*
 * Constructor and destructor.
 */
//go:linkname HashCreate C.xmlHashCreate
func HashCreate(size c.Int) HashTablePtr

//go:linkname HashCreateDict C.xmlHashCreateDict
func HashCreateDict(size c.Int, dict DictPtr) HashTablePtr

//go:linkname HashFree C.xmlHashFree
func HashFree(hash HashTablePtr, dealloc HashDeallocator)

//go:linkname HashDefaultDeallocator C.xmlHashDefaultDeallocator
func HashDefaultDeallocator(entry c.Pointer, name *Char)

/*
 * Add a new entry to the hash table.
 */
//go:linkname HashAdd C.xmlHashAdd
func HashAdd(hash HashTablePtr, name *Char, userdata c.Pointer) c.Int

//go:linkname HashAddEntry C.xmlHashAddEntry
func HashAddEntry(hash HashTablePtr, name *Char, userdata c.Pointer) c.Int

//go:linkname HashUpdateEntry C.xmlHashUpdateEntry
func HashUpdateEntry(hash HashTablePtr, name *Char, userdata c.Pointer, dealloc HashDeallocator) c.Int

//go:linkname HashAdd2 C.xmlHashAdd2
func HashAdd2(hash HashTablePtr, name *Char, name2 *Char, userdata c.Pointer) c.Int

//go:linkname HashAddEntry2 C.xmlHashAddEntry2
func HashAddEntry2(hash HashTablePtr, name *Char, name2 *Char, userdata c.Pointer) c.Int

//go:linkname HashUpdateEntry2 C.xmlHashUpdateEntry2
func HashUpdateEntry2(hash HashTablePtr, name *Char, name2 *Char, userdata c.Pointer, dealloc HashDeallocator) c.Int

//go:linkname HashAdd3 C.xmlHashAdd3
func HashAdd3(hash HashTablePtr, name *Char, name2 *Char, name3 *Char, userdata c.Pointer) c.Int

//go:linkname HashAddEntry3 C.xmlHashAddEntry3
func HashAddEntry3(hash HashTablePtr, name *Char, name2 *Char, name3 *Char, userdata c.Pointer) c.Int

//go:linkname HashUpdateEntry3 C.xmlHashUpdateEntry3
func HashUpdateEntry3(hash HashTablePtr, name *Char, name2 *Char, name3 *Char, userdata c.Pointer, dealloc HashDeallocator) c.Int

/*
 * Remove an entry from the hash table.
 */
//go:linkname HashRemoveEntry C.xmlHashRemoveEntry
func HashRemoveEntry(hash HashTablePtr, name *Char, dealloc HashDeallocator) c.Int

//go:linkname HashRemoveEntry2 C.xmlHashRemoveEntry2
func HashRemoveEntry2(hash HashTablePtr, name *Char, name2 *Char, dealloc HashDeallocator) c.Int

//go:linkname HashRemoveEntry3 C.xmlHashRemoveEntry3
func HashRemoveEntry3(hash HashTablePtr, name *Char, name2 *Char, name3 *Char, dealloc HashDeallocator) c.Int

/*
 * Retrieve the payload.
 */
//go:linkname HashLookup C.xmlHashLookup
func HashLookup(hash HashTablePtr, name *Char) c.Pointer

//go:linkname HashLookup2 C.xmlHashLookup2
func HashLookup2(hash HashTablePtr, name *Char, name2 *Char) c.Pointer

//go:linkname HashLookup3 C.xmlHashLookup3
func HashLookup3(hash HashTablePtr, name *Char, name2 *Char, name3 *Char) c.Pointer

//go:linkname HashQLookup C.xmlHashQLookup
func HashQLookup(hash HashTablePtr, prefix *Char, name *Char) c.Pointer

//go:linkname HashQLookup2 C.xmlHashQLookup2
func HashQLookup2(hash HashTablePtr, prefix *Char, name *Char, prefix2 *Char, name2 *Char) c.Pointer

//go:linkname HashQLookup3 C.xmlHashQLookup3
func HashQLookup3(hash HashTablePtr, prefix *Char, name *Char, prefix2 *Char, name2 *Char, prefix3 *Char, name3 *Char) c.Pointer

/*
 * Helpers.
 */
//go:linkname HashCopySafe C.xmlHashCopySafe
func HashCopySafe(hash HashTablePtr, copy HashCopier, dealloc HashDeallocator) HashTablePtr

//go:linkname HashCopy C.xmlHashCopy
func HashCopy(hash HashTablePtr, copy HashCopier) HashTablePtr

//go:linkname HashSize C.xmlHashSize
func HashSize(hash HashTablePtr) c.Int

//go:linkname HashScan C.xmlHashScan
func HashScan(hash HashTablePtr, scan HashScanner, data c.Pointer)

//go:linkname HashScan3 C.xmlHashScan3
func HashScan3(hash HashTablePtr, name *Char, name2 *Char, name3 *Char, scan HashScanner, data c.Pointer)

//go:linkname HashScanFull C.xmlHashScanFull
func HashScanFull(hash HashTablePtr, scan HashScannerFull, data c.Pointer)

//go:linkname HashScanFull3 C.xmlHashScanFull3
func HashScanFull3(hash HashTablePtr, name *Char, name2 *Char, name3 *Char, scan HashScannerFull, data c.Pointer)
