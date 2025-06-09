package zlib

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

const MAX_MEM_LEVEL = 9
const MAX_WBITS = 15

type ZSizeT c.SizeT
type Byte c.Char
type UInt c.Uint
type ULong c.Ulong
type Bytef Byte
type Charf c.Char
type Intf c.Int
type UIntf UInt
type ULongf ULong
type Voidpc c.Pointer
type Voidpf c.Pointer
type Voidp c.Pointer
type ZCrcT c.Uint
