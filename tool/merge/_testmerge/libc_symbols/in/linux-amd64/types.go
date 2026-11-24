package libc

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

type X__Uint64T = c.Ulong

type X__OffT = c.Long

type X__Off64T = c.Long

type X__SocklenT = c.Uint

type SizeT = c.Ulong

type SaFamilyT = uint16

type Sockaddr struct {
	SaFamily SaFamilyT
	SaData   [14]c.Char
}

type SocklenT = X__SocklenT

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
	AiAddr      *Sockaddr
	AiCanonname *c.Char
	AiNext      *Addrinfo
}

type X_IO_FILE struct {
	X_flags          c.Int
	X_IO_read_ptr    *c.Char
	X_IO_read_end    *c.Char
	X_IO_read_base   *c.Char
	X_IO_write_base  *c.Char
	X_IO_write_ptr   *c.Char
	X_IO_write_end   *c.Char
	X_IO_buf_base    *c.Char
	X_IO_buf_end     *c.Char
	X_IO_save_base   *c.Char
	X_IO_backup_base *c.Char
	X_IO_save_end    *c.Char
	X_markers        *X_IO_marker
	X_chain          *X_IO_FILE
	X_fileno         c.Int
	_xgo_bits_0      [3]uint8
	X_shortBackupbuf [1]c.Char
	X_oldOffset      X__OffT
	X_curColumn      uint16
	X_vtableOffset   int8
	X_shortbuf       [1]c.Char
	X_lock           *X_IO_lock_t
	X_offset         X__Off64T
	X_codecvt        *X_IO_codecvt
	X_wideData       *X_IO_wide_data
	X_freeresList    *X_IO_FILE
	X_freeresBuf     unsafe.Pointer
	X_prevchain      **X_IO_FILE
	X_mode           c.Int
	X_unused3        c.Int
	X_totalWritten   X__Uint64T
	X_unused2        [8]c.Char
}

type X_IO_codecvt struct {
}

type X_IO_marker struct {
}

type X_IO_wide_data struct {
}

type X_IO_lock_t = c.Void

type FILE = X_IO_FILE

type PthreadAttrT struct {
	_xgo_union [7]uint64
}

type PthreadT = c.Ulong
