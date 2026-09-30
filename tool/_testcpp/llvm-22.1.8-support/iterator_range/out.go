package iterator_range

const LLGoPackage = "link: -L$(llvm-config --libdir) -lLLVM; -lLLVM"

// A range adaptor for a pair of iterators.
//
// This just wraps two iterators into a range-compatible interface. Nothing
// fancy at all.
type IteratorRange[IteratorT any] struct {
	beginIterator IteratorT
	endIterator   IteratorT
}
