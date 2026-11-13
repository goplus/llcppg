package foo

import "github.com/goplus/lib/c"

type ByteFlags uint8

const (
	ByteFlags_High ByteFlags = 128
	ByteFlags_All  ByteFlags = 255
)

type ShortFlags uint16

const (
	ShortFlags_High ShortFlags = 32768
	ShortFlags_All  ShortFlags = 65535
)

type WideFlags c.UlongLong

const (
	WideFlags_High WideFlags = 9223372036854775808
	WideFlags_All  WideFlags = 18446744073709551615
)

type SignedValues c.LongLong

const (
	SignedValues_Negative SignedValues = -1
	SignedValues_Minimum  SignedValues = -9223372036854775808
)
