package STLForwardCompat

import "github.com/goplus/lib/c"

const LLGoPackage = "link: -L$(llvm-config --libdir) -lLLVM; -lLLVM"

type Identity struct {
}
type IdentityIsTransparent = c.Void
type FromRangeT struct {
}
type RemoveCvref[T any] struct {
}
type TypeIdentity[T any] struct {
}
