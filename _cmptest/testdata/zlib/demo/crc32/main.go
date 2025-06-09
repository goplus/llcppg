package main

import (
	"fmt"
	"unsafe"

	"zlib"
)

func main() {
	ul := zlib.ULong(0)
	data := "Hello world"
	res := ul.Crc32Z(
		(*zlib.Bytef)(unsafe.Pointer(unsafe.StringData(data))),
		zlib.ZSizeT(uintptr(len(data))),
	)
	fmt.Printf("%08x\n", res)
}
