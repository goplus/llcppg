package main

import (
	"fmt"
	"unsafe"
	"zlibstatic"
)

func main() {
	ul := zlibstatic.ULong(0)
	data := "Hello world"
	res := ul.Crc32Z(
		(*zlibstatic.Bytef)(unsafe.Pointer(unsafe.StringData(data))),
		zlibstatic.ZSizeT(uintptr(len(data))),
	)
	fmt.Printf("%08x\n", res)
}
