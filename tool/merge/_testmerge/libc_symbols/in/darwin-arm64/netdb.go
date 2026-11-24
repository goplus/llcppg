package libc

import (
	"github.com/goplus/lib/c"
)

//llgo:link (*Addrinfo).Free C.freeaddrinfo
func (self *Addrinfo) Free() {
}

//llgo:link (*Sockaddr).Getnameinfo C.getnameinfo
func (self *Sockaddr) Getnameinfo(_llcppg_param2 SocklenT, _llcppg_param3 *c.Char, _llcppg_param4 SocklenT, _llcppg_param5 *c.Char, _llcppg_param6 SocklenT, _llcppg_param7 c.Int) c.Int {
	return 0
}

//llgo:link (*Hostent).Free C.freehostent
func (self *Hostent) Free() {
}
