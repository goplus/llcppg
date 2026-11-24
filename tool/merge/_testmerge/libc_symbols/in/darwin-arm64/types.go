package libc

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

type X__Uint8T = uint8

type X__Uint32T = c.Uint

type X__Int64T = c.LongLong

type X__DarwinSizeT = c.Ulong

type X__DarwinSocklenT = X__Uint32T

type X__DarwinOffT = X__Int64T

type SizeT = X__DarwinSizeT

type SaFamilyT = X__Uint8T

type SocklenT = X__DarwinSocklenT

type Sockaddr struct {
	SaLen    X__Uint8T
	SaFamily SaFamilyT
	SaData   [14]c.Char
}

type Hostent struct {
	HName     *c.Char
	HAliases  **c.Char
	HAddrtype c.Int
	HLength   c.Int
	HAddrList **c.Char
}

type Addrinfo struct {
	AiFlags     c.Int
	AiFamily    c.Int
	AiSocktype  c.Int
	AiProtocol  c.Int
	AiAddrlen   SocklenT
	AiCanonname *c.Char
	AiAddr      *Sockaddr
	AiNext      *Addrinfo
}

type X__Sbuf struct {
	X_base *uint8
	X_size c.Int
}

type X__SFILE struct {
	X_p       *uint8
	X_r       c.Int
	X_w       c.Int
	X_flags   int16
	X_file    int16
	X_bf      X__Sbuf
	X_lbfsize c.Int
	X_cookie  unsafe.Pointer
	X_close   func(_llcppg_param1 unsafe.Pointer) c.Int
	X_read    func(_llcppg_param1 unsafe.Pointer, _llcppg_param2 *c.Char, _llcppg_param3 c.Int) c.Int
	X_seek    func(_llcppg_param1 unsafe.Pointer, _llcppg_param2 FposT, _llcppg_param3 c.Int) FposT
	X_write   func(_llcppg_param1 unsafe.Pointer, _llcppg_param2 *c.Char, _llcppg_param3 c.Int) c.Int
	X_ub      X__Sbuf
	X_extra   *X__SFILEX
	X_ur      c.Int
	X_ubuf    [3]uint8
	X_nbuf    [1]uint8
	X_lb      X__Sbuf
	X_blksize c.Int
	X_offset  FposT
}

type X__SFILEX struct {
}

type FposT = X__DarwinOffT

type FILE = X__SFILE

type X__DarwinPthreadHandlerRec struct {
	X__Routine func(_llcppg_param1 unsafe.Pointer)
	X__Arg     unsafe.Pointer
	X__Next    *X__DarwinPthreadHandlerRec
}

type X_opaquePthreadAttrT struct {
	X__Sig    c.Long
	X__Opaque [56]c.Char
}

type X_opaquePthreadT struct {
	X__Sig          c.Long
	X__CleanupStack *X__DarwinPthreadHandlerRec
	X__Opaque       [8176]c.Char
}

type X__DarwinPthreadAttrT = X_opaquePthreadAttrT

type X__DarwinPthreadT = *X_opaquePthreadT

type PthreadAttrT = X__DarwinPthreadAttrT

type PthreadT = X__DarwinPthreadT
