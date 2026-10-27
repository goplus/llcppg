package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

const (
	LLGoPackage = "link: -L/path/foo -lfoo"
	LLGoFiles   = "-I/path/foo/include: _wrap/llcppg.c"
)
const ZLIB_VERSION = "1.3.1"
const Z_DEFAULT_LEVEL = -1

type GzFileS struct {
	Have c.Uint
	Next *uint8
	Pos  c.Long
}
type StreamS struct {
	NextIn  *uint8
	AvailIn c.Uint
	TotalIn c.Ulong
}
type OpaqueS struct {
}
type GzFile = *GzFileS
type Stream = StreamS
type Streamp = *Stream

// An incomplete (opaque) struct: its layout is unknown, so sizeof over it must
// not be detected (rule 3).
type OpaqueT OpaqueS

//llgo:link (*OpaqueT).OpenOpaque_ C.openOpaque_
func (o *OpaqueT) OpenOpaque_(size c.Int) c.Int {
	return 0
}

// Read one byte, falling back to the real function.
//
//llgo:link GzFile.Gzgetc_ C.gzgetc_
func (file GzFile) Gzgetc_() c.Int {
	return 0
}

//llgo:link Streamp.DeflateInit_ C.deflateInit_
func (strm Streamp) DeflateInit_(level c.Int, version *c.Char, stream_size c.Int) c.Int {
	return 0
}

//llgo:link Streamp.InflateInit_ C.inflateInit_
func (strm Streamp) InflateInit_(version *c.Char, stream_size c.Int) c.Int {
	return 0
}

// Initialize for deflation with extra options.
//
//llgo:link Streamp.InflateBackInit_ C.inflateBackInit_
func (strm Streamp) InflateBackInit_(windowBits c.Int, window *uint8, version *c.Char, stream_size c.Int) c.Int {
	return 0
}

//llgo:link Streamp.Reset_ C.reset_
func (strm Streamp) Reset_(flags c.Int) {
}

// Return the library version.
//
//go:linkname Version_ C.version_
func Version_() c.Int

// A package-level callee: its first parameter is a scalar, so it is not
// converted to a method. Used to cover the package-function call path and
// numeric/constant-macro arguments.
//
//go:linkname Configure_ C.configure_
func Configure_(mode c.Int, level c.Int) c.Int

// gzgetc has real logic, so it stays an explicit MacroFunc entry (C wrapper).
//
//llgo:link GzFile.Gzgetc C._llcppg_z_gzgetc
func (g GzFile) Gzgetc() c.Int {
	return 0
}

// Initialize for deflation, filling in the version and stream size.
func (strm Streamp) DeflateInit(level c.Int) c.Int {
	return strm.DeflateInit_(level, c.Str("1.3.1"), c.Int(unsafe.Sizeof(Stream{})))
}
func (strm Streamp) InflateInit() c.Int {
	return strm.InflateInit_(c.Str("1.3.1"), c.Int(unsafe.Sizeof(Stream{})))
}
func (strm Streamp) InflateBackInit(windowBits c.Int, window *uint8) c.Int {
	return strm.InflateBackInit_(windowBits, window, c.Str("1.3.1"), c.Int(unsafe.Sizeof(Stream{})))
}

// void-returning forwarder: the body is a bare call with no result.
func (strm Streamp) Reset() {
	strm.Reset_(0)
}

// Package-level callee with a numeric literal and an object-like constant macro.
func Configure(mode c.Int) c.Int {
	return Configure_(mode, -1)
}

// Skipped: sizeof over an incomplete type whose layout Go cannot know (rule 3).
func (o *OpaqueT) OpenOpaque() c.Int {
	return o.OpenOpaque_(c.Int(unsafe.Sizeof(OpaqueT{})))
}
