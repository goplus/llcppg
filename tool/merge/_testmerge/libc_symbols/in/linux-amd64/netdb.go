package libc

import (
	libc "github.com/goplus/lib/c"
)

//llgo:link (*Addrinfo).Free C.freeaddrinfo
func (self *Addrinfo) Free() {
}

//llgo:link (*Sockaddr).Getnameinfo C.getnameinfo
func (self *Sockaddr) Getnameinfo(__salen SocklenT, __host *libc.Char, __hostlen SocklenT, __serv *libc.Char, __servlen SocklenT, __flags libc.Int) libc.Int {
	return 0
}
