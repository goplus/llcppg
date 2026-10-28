package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

type SSL struct {
}

//go:linkname SSL_ctrl C._Z8SSL_ctrlP3SSLilPv
func SSL_ctrl(ssl *SSL, cmd c.Int, larg c.Long, parg unsafe.Pointer) c.Long
func SSL_set_tlsext_status_type(ssl *SSL, type_ c.Long) c.Long {
	return SSL_ctrl(ssl, 0, type_, nil)
}
