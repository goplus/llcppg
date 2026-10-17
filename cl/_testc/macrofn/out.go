package foo

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

const (
	LLGoPackage = "link: -L/path/foo -lfoo"
	LLGoFiles   = "-I/path/foo/include: _wrap/llcppg.c"
)

type GzFileS struct {
	Have c.Uint
	Next *uint8
	Pos  c.Long
}
type StreamS struct {
}
type GzFile = *GzFileS
type Stream StreamS
type Streamp = *Stream

// llgo:link GzFile.Gzgetc C._llcppg_z_gzgetc
func (g GzFile) Gzgetc() c.Int {
	return 0
}

// Initialize for deflation, filling in the version and stream size.
//
// llgo:link Streamp.DeflateInit C._llcppg_z_deflateInit
func (strm Streamp) DeflateInit(level c.Int) c.Int {
	return 0
}

// llgo:link Streamp.Reset C._llcppg_z_reset
func (strm Streamp) Reset(flags c.Int) {
}

// Flush any pending output.
//
// llgo:link Streamp.Flush C._llcppg_z_flush
func (strm Streamp) Flush() {
}

// Read one byte, falling back to the real function.
//
// llgo:link GzFile.Gzgetc_ C.gzgetc_
func (file GzFile) Gzgetc_() c.Int {
	return 0
}

// llgo:link Streamp.DeflateInit_ C.deflateInit_
func (strm Streamp) DeflateInit_(level c.Int, version *c.Char, stream_size c.Int) c.Int {
	return 0
}

// llgo:link Streamp.Reset_ C.reset_
func (strm Streamp) Reset_(flags c.Int, version *c.Char) {
}

// llgo:link Streamp.Flush_ C.flush_
func (strm Streamp) Flush_() {
}

// Return the library version.
//
//go:linkname Version_ C.version_
func Version_() c.Int

// Release internal buffers.
//
// llgo:link Streamp.Cleanup_ C.cleanup_
func (strm Streamp) Cleanup_() {
}
