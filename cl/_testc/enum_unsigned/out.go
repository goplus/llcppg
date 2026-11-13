package foo

import "github.com/goplus/lib/c"

type UnsignedFlags c.Uint

const (
	FlagHigh UnsignedFlags = 2147483648
	FlagAll  UnsignedFlags = 4294967295
)

type UnsignedWide c.Ulong

const (
	WideHigh UnsignedWide = 9223372036854775808
	WideAll  UnsignedWide = 18446744073709551615
)

type SignedValues c.Int

const (
	Negative SignedValues = -1
	Positive SignedValues = 1
)

type MixedValues c.Long

const (
	MixedNegative MixedValues = -1
	MixedHigh     MixedValues = 2147483648
)
const (
	AnonymousHigh = 2147483648
	AnonymousAll  = 4294967295
)
const AnonymousWide = 18446744073709551615
