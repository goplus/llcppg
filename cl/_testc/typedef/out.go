package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

type CXTUResourceUsageKind c.Int

const (
	CXTUResourceUsage_AST         CXTUResourceUsageKind = 1
	CXTUResourceUsage_Identifiers CXTUResourceUsageKind = 2
)

type CXTUResourceUsageEntry struct {
	Kind   CXTUResourceUsageKind
	Amount c.Ulong
}
type CXTUResourceUsage struct {
	Data       unsafe.Pointer
	NumEntries c.Uint
	Entries    *CXTUResourceUsageEntry
}
