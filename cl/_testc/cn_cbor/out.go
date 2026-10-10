package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

const XGoPackage = true

type CnCborType c.Uint

const (
	CN_CBOR_FALSE         CnCborType = 0
	CN_CBOR_TRUE          CnCborType = 1
	CN_CBOR_NULL          CnCborType = 2
	CN_CBOR_UNDEF         CnCborType = 3
	CN_CBOR_UINT          CnCborType = 4
	CN_CBOR_INT           CnCborType = 5
	CN_CBOR_BYTES         CnCborType = 6
	CN_CBOR_TEXT          CnCborType = 7
	CN_CBOR_BYTES_CHUNKED CnCborType = 8
	CN_CBOR_TEXT_CHUNKED  CnCborType = 9
	CN_CBOR_ARRAY         CnCborType = 10
	CN_CBOR_MAP           CnCborType = 11
	CN_CBOR_TAG           CnCborType = 12
	CN_CBOR_SIMPLE        CnCborType = 13
	CN_CBOR_DOUBLE        CnCborType = 14
	CN_CBOR_FLOAT         CnCborType = 15
	CN_CBOR_INVALID       CnCborType = 16
)

type CnCborFlags c.Uint

const (
	CN_CBOR_FL_COUNT CnCborFlags = 1
	CN_CBOR_FL_INDEF CnCborFlags = 2
	CN_CBOR_FL_OWNER CnCborFlags = 128
)

type Uint8T = uint8
type CnCbor struct {
	Type       CnCborType
	Flags      CnCborFlags
	V          _llcppg_anon_0
	Length     c.Int
	FirstChild *CnCbor
	LastChild  *CnCbor
	Next       *CnCbor
	Parent     *CnCbor
}
type _llcppg_anon_0 struct {
	_xgo_union [1]uint64
}

func (p *_llcppg_anon_0) XGof_ref_bytes() **Uint8T {
	return (**Uint8T)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_0) XGof_ref_str() **c.Char {
	return (**c.Char)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_0) XGof_ref_sint() *c.Long {
	return (*c.Long)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_0) XGof_ref_uint() *c.Ulong {
	return (*c.Ulong)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_0) XGof_ref_dbl() *c.Double {
	return (*c.Double)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_0) XGof_ref_f() *c.Float {
	return (*c.Float)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_0) XGof_ref_count() *c.Ulong {
	return (*c.Ulong)(unsafe.Pointer(p))
}
