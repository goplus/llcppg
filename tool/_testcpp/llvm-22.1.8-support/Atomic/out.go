package Atomic

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

const LLGoPackage = "link: -L$(llvm-config --libdir) -lLLVM; -lLLVM"

type SysCasFlag = c.Uint32T

//go:linkname SysMemoryFence C._ZN4llvm3sys11MemoryFenceEv
func SysMemoryFence()

//go:linkname SysCompareAndSwap C._ZN4llvm3sys14CompareAndSwapEPVjjj
func SysCompareAndSwap(ptr *SysCasFlag, new_value SysCasFlag, old_value SysCasFlag) SysCasFlag
