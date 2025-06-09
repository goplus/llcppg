package bzip2

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

const RUN = 0
const FLUSH = 1
const FINISH = 2
const OK = 0
const RUN_OK = 1
const FLUSH_OK = 2
const FINISH_OK = 3
const STREAM_END = 4
const MAX_UNUSED = 5000

type BzStream struct {
	NextIn       *c.Char
	AvailIn      c.Uint
	TotalInLo32  c.Uint
	TotalInHi32  c.Uint
	NextOut      *c.Char
	AvailOut     c.Uint
	TotalOutLo32 c.Uint
	TotalOutHi32 c.Uint
	State        c.Pointer
	Bzalloc      c.Pointer
	Bzfree       c.Pointer
	Opaque       c.Pointer
}

/*-- Core (low-level) library functions --*/
// llgo:link (*BzStream).CompressInit C.BZ2_bzCompressInit
func (recv_ *BzStream) CompressInit(blockSize100k c.Int, verbosity c.Int, workFactor c.Int) c.Int {
	return 0
}

// llgo:link (*BzStream).Compress C.BZ2_bzCompress
func (recv_ *BzStream) Compress(action c.Int) c.Int {
	return 0
}

// llgo:link (*BzStream).CompressEnd C.BZ2_bzCompressEnd
func (recv_ *BzStream) CompressEnd() c.Int {
	return 0
}

// llgo:link (*BzStream).DecompressInit C.BZ2_bzDecompressInit
func (recv_ *BzStream) DecompressInit(verbosity c.Int, small c.Int) c.Int {
	return 0
}

// llgo:link (*BzStream).Decompress C.BZ2_bzDecompress
func (recv_ *BzStream) Decompress() c.Int {
	return 0
}

// llgo:link (*BzStream).DecompressEnd C.BZ2_bzDecompressEnd
func (recv_ *BzStream) DecompressEnd() c.Int {
	return 0
}

type BZFILE c.Void

//go:linkname ReadOpen C.BZ2_bzReadOpen
func ReadOpen(bzerror *c.Int, f *c.FILE, verbosity c.Int, small c.Int, unused c.Pointer, nUnused c.Int) *BZFILE

//go:linkname ReadClose C.BZ2_bzReadClose
func ReadClose(bzerror *c.Int, b *BZFILE)

//go:linkname ReadGetUnused C.BZ2_bzReadGetUnused
func ReadGetUnused(bzerror *c.Int, b *BZFILE, unused *c.Pointer, nUnused *c.Int)

//go:linkname Read C.BZ2_bzRead
func Read(bzerror *c.Int, b *BZFILE, buf c.Pointer, len c.Int) c.Int

//go:linkname WriteOpen C.BZ2_bzWriteOpen
func WriteOpen(bzerror *c.Int, f *c.FILE, blockSize100k c.Int, verbosity c.Int, workFactor c.Int) *BZFILE

//go:linkname Write C.BZ2_bzWrite
func Write(bzerror *c.Int, b *BZFILE, buf c.Pointer, len c.Int)

//go:linkname WriteClose C.BZ2_bzWriteClose
func WriteClose(bzerror *c.Int, b *BZFILE, abandon c.Int, nbytes_in *c.Uint, nbytes_out *c.Uint)

//go:linkname WriteClose64 C.BZ2_bzWriteClose64
func WriteClose64(bzerror *c.Int, b *BZFILE, abandon c.Int, nbytes_in_lo32 *c.Uint, nbytes_in_hi32 *c.Uint, nbytes_out_lo32 *c.Uint, nbytes_out_hi32 *c.Uint)

/*-- Utility functions --*/
//go:linkname BuffToBuffCompress C.BZ2_bzBuffToBuffCompress
func BuffToBuffCompress(dest *c.Char, destLen *c.Uint, source *c.Char, sourceLen c.Uint, blockSize100k c.Int, verbosity c.Int, workFactor c.Int) c.Int

//go:linkname BuffToBuffDecompress C.BZ2_bzBuffToBuffDecompress
func BuffToBuffDecompress(dest *c.Char, destLen *c.Uint, source *c.Char, sourceLen c.Uint, small c.Int, verbosity c.Int) c.Int

/*--
   Code contributed by Yoshioka Tsuneo (tsuneo@rr.iij4u.or.jp)
   to support better zlib compatibility.
   This code is not _officially_ part of libbzip2 (yet);
   I haven't tested it, documented it, or considered the
   threading-safeness of it.
   If this code breaks, please contact both Yoshioka and me.
--*/
//go:linkname LibVersion C.BZ2_bzlibVersion
func LibVersion() *c.Char

//go:linkname Open C.BZ2_bzopen
func Open(path *c.Char, mode *c.Char) *BZFILE

//go:linkname Dopen C.BZ2_bzdopen
func Dopen(fd c.Int, mode *c.Char) *BZFILE

// llgo:link (*BZFILE).Read C.BZ2_bzread
func (recv_ *BZFILE) Read(buf c.Pointer, len c.Int) c.Int {
	return 0
}

// llgo:link (*BZFILE).Write C.BZ2_bzwrite
func (recv_ *BZFILE) Write(buf c.Pointer, len c.Int) c.Int {
	return 0
}

// llgo:link (*BZFILE).Flush C.BZ2_bzflush
func (recv_ *BZFILE) Flush() c.Int {
	return 0
}

// llgo:link (*BZFILE).Close C.BZ2_bzclose
func (recv_ *BZFILE) Close() {
}

// llgo:link (*BZFILE).Error C.BZ2_bzerror
func (recv_ *BZFILE) Error(errnum *c.Int) *c.Char {
	return nil
}
