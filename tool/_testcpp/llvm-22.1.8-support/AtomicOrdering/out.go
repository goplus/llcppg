package AtomicOrdering

import "github.com/goplus/lib/c"

const LLGoPackage = "link: -L$(llvm-config --libdir) -lLLVM; -lLLVM"

// Atomic ordering for C11 / C++11's memory models.
//
// These values cannot change because they are shared with standard library
// implementations as well as with other compilers.
type AtomicOrderingCABI c.Int

const (
	AtomicOrderingCABI_Relaxed AtomicOrderingCABI = 0
	AtomicOrderingCABI_Consume AtomicOrderingCABI = 1
	AtomicOrderingCABI_Acquire AtomicOrderingCABI = 2
	AtomicOrderingCABI_Release AtomicOrderingCABI = 3
	AtomicOrderingCABIAcqRel   AtomicOrderingCABI = 4
	AtomicOrderingCABISeqCst   AtomicOrderingCABI = 5
)

// Atomic ordering for LLVM's memory model.
//
// C++ defines ordering as a lattice. LLVM supplements this with NotAtomic and
// Unordered, which are both below the C++ orders.
//
// not_atomic-->unordered-->relaxed-->release--------------->acq_rel-->seq_cst
//                                   \-->consume-->acquire--/
type AtomicOrdering c.Int

const (
	AtomicOrdering_NotAtomic              AtomicOrdering = 0
	AtomicOrdering_Unordered              AtomicOrdering = 1
	AtomicOrdering_Monotonic              AtomicOrdering = 2
	AtomicOrdering_Acquire                AtomicOrdering = 4
	AtomicOrdering_Release                AtomicOrdering = 5
	AtomicOrdering_AcquireRelease         AtomicOrdering = 6
	AtomicOrdering_SequentiallyConsistent AtomicOrdering = 7
	AtomicOrdering_LAST                   AtomicOrdering = 7
)
