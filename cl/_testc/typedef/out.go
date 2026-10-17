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
	Kind   CXTUResourceUsageKind
	Amount c.Ulong
}

// The memory usage of a CXTranslationUnit, broken into categories.
type CXTUResourceUsage struct {
	Data       unsafe.Pointer
	NumEntries c.Uint
	Entries    *CXTUResourceUsageEntry
}
