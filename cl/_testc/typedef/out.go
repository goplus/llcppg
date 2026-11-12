package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

// Categorizes how memory is being used by a translation unit.
type CXTUResourceUsageKind c.Uint

const (
	CXTUResourceUsage_AST         CXTUResourceUsageKind = 1
	CXTUResourceUsage_Identifiers CXTUResourceUsageKind = 2
)

type CXTUResourceUsageEntry struct {
	// The memory usage category.
	Kind CXTUResourceUsageKind
	// Amount of resources used.
	// The units will depend on the resource kind.
	Amount c.Ulong
}

// The memory usage of a CXTranslationUnit, broken into categories.
type CXTUResourceUsage struct {
	// Private data member, used for queries.
	Data unsafe.Pointer
	// The number of entries in the 'entries' array.
	NumEntries c.Uint
	// An array of key-value pairs, representing the breakdown of memory
	// usage.
	Entries *CXTUResourceUsageEntry
}
