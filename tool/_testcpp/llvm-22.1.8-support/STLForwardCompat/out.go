package STLForwardCompat

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

const LLGoPackage = "link: -L$(llvm-config --libdir) -lLLVM; -lLLVM"

type Identity struct {
}
type IdentityIsTransparent = c.Void
type FromRangeT struct {
}

//go:linkname NumbersE C._ZN4llvm7numbers1eE
var NumbersE c.Double

//go:linkname NumbersEgamma C._ZN4llvm7numbers6egammaE
var NumbersEgamma c.Double

//go:linkname NumbersLn2 C._ZN4llvm7numbers3ln2E
var NumbersLn2 c.Double

//go:linkname NumbersLn10 C._ZN4llvm7numbers4ln10E
var NumbersLn10 c.Double

//go:linkname NumbersLog2e C._ZN4llvm7numbers5log2eE
var NumbersLog2e c.Double

//go:linkname NumbersLog10e C._ZN4llvm7numbers6log10eE
var NumbersLog10e c.Double

//go:linkname NumbersPi C._ZN4llvm7numbers2piE
var NumbersPi c.Double

//go:linkname NumbersInvPi C._ZN4llvm7numbers6inv_piE
var NumbersInvPi c.Double

//go:linkname NumbersInvSqrtpi C._ZN4llvm7numbers10inv_sqrtpiE
var NumbersInvSqrtpi c.Double

//go:linkname NumbersSqrt2 C._ZN4llvm7numbers5sqrt2E
var NumbersSqrt2 c.Double

//go:linkname NumbersInvSqrt2 C._ZN4llvm7numbers9inv_sqrt2E
var NumbersInvSqrt2 c.Double

//go:linkname NumbersSqrt3 C._ZN4llvm7numbers5sqrt3E
var NumbersSqrt3 c.Double

//go:linkname NumbersInvSqrt3 C._ZN4llvm7numbers9inv_sqrt3E
var NumbersInvSqrt3 c.Double

//go:linkname NumbersPhi C._ZN4llvm7numbers3phiE
var NumbersPhi c.Double

type RemoveCvref[T any] struct {
}
type TypeIdentity[T any] struct {
}
type TypeIdentityType[T any] = T

//go:linkname FromRange C._ZN4llvm10from_rangeE
var FromRange FromRangeT
