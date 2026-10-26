package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

type Asn1ObjectSt struct {
}
type Asn1StringSt struct {
}
type Asn1TypeSt struct {
}
type ASN1_STRING Asn1StringSt
type ASN1_BIT_STRING Asn1StringSt
type ASN1_TYPE Asn1TypeSt
type ASN1_OBJECT Asn1ObjectSt

//go:linkname ASN1_STRING_new C.ASN1_STRING_new
func ASN1_STRING_new() *ASN1_STRING

//llgo:link (*ASN1_STRING).Free C.ASN1_STRING_free
func (a *ASN1_STRING) Free() {
}

//llgo:link (*ASN1_STRING).ClearFree C.ASN1_STRING_clear_free
func (a *ASN1_STRING) ClearFree() {
}

//llgo:link (*ASN1_STRING).Copy C.ASN1_STRING_copy
func (dst *ASN1_STRING) Copy(str *ASN1_STRING) c.Int {
	return 0
}

//llgo:link (*ASN1_STRING).I2d_DISPLAYTEXT C.i2d_DISPLAYTEXT
func (a *ASN1_STRING) I2d_DISPLAYTEXT(pp **c.Char) c.Int {
	return 0
}

//go:linkname ASN1_BIT_STRING_new C.ASN1_BIT_STRING_new
func ASN1_BIT_STRING_new() *ASN1_BIT_STRING

//llgo:link (*ASN1_BIT_STRING).Free C.ASN1_BIT_STRING_free
func (a *ASN1_BIT_STRING) Free() {
}

//llgo:link (*ASN1_BIT_STRING).Set C.ASN1_BIT_STRING_set
func (a *ASN1_BIT_STRING) Set(d *uint8, length c.Int) c.Int {
	return 0
}

//llgo:link (*ASN1_BIT_STRING).SetBit C.ASN1_BIT_STRING_set_bit
func (a *ASN1_BIT_STRING) SetBit(n c.Int, value c.Int) c.Int {
	return 0
}

//llgo:link (*ASN1_BIT_STRING).Bit C.ASN1_BIT_STRING_get_bit
func (a *ASN1_BIT_STRING) Bit(n c.Int) c.Int {
	return 0
}

//llgo:link (*ASN1_BIT_STRING).Check C.ASN1_BIT_STRING_check
func (a *ASN1_BIT_STRING) Check(flags *uint8, flags_len c.Int) c.Int {
	return 0
}

//go:linkname ASN1_generate_nconf C.ASN1_generate_nconf
func ASN1_generate_nconf(str *c.Char, nconf unsafe.Pointer) *ASN1_TYPE

//llgo:link (*ASN1_OBJECT).I2d C.i2d_ASN1_OBJECT
func (obj *ASN1_OBJECT) I2d(pp **c.Char) c.Int {
	return 0
}

//llgo:link (*ASN1_OBJECT).AddObject C.OBJ_add_object
func (obj *ASN1_OBJECT) AddObject() c.Int {
	return 0
}
