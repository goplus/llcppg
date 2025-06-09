package libxml2

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type X_xmlMutex struct {
	Unused [8]uint8
}
type Mutex X_xmlMutex
type MutexPtr *Mutex

type X_xmlRMutex struct {
	Unused [8]uint8
}
type RMutex X_xmlRMutex
type RMutexPtr *RMutex

//go:linkname CheckThreadLocalStorage C.xmlCheckThreadLocalStorage
func CheckThreadLocalStorage() c.Int

//go:linkname NewMutex C.xmlNewMutex
func NewMutex() MutexPtr

//go:linkname MutexLock C.xmlMutexLock
func MutexLock(tok MutexPtr)

//go:linkname MutexUnlock C.xmlMutexUnlock
func MutexUnlock(tok MutexPtr)

//go:linkname FreeMutex C.xmlFreeMutex
func FreeMutex(tok MutexPtr)

//go:linkname NewRMutex C.xmlNewRMutex
func NewRMutex() RMutexPtr

//go:linkname RMutexLock C.xmlRMutexLock
func RMutexLock(tok RMutexPtr)

//go:linkname RMutexUnlock C.xmlRMutexUnlock
func RMutexUnlock(tok RMutexPtr)

//go:linkname FreeRMutex C.xmlFreeRMutex
func FreeRMutex(tok RMutexPtr)

/*
 * Library wide APIs.
 */
//go:linkname InitThreads C.xmlInitThreads
func InitThreads()

//go:linkname LockLibrary C.xmlLockLibrary
func LockLibrary()

//go:linkname UnlockLibrary C.xmlUnlockLibrary
func UnlockLibrary()

//go:linkname GetThreadId C.xmlGetThreadId
func GetThreadId() c.Int

//go:linkname IsMainThread C.xmlIsMainThread
func IsMainThread() c.Int

//go:linkname CleanupThreads C.xmlCleanupThreads
func CleanupThreads()
