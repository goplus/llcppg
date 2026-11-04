package foo

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type ClntStat c.Uint

const (
	RPC_SUCCESS ClntStat = 0
	RPC_FAILED  ClntStat = 1
)

//go:linkname RpcStatus C.rpc_status
func RpcStatus() *ClntStat
