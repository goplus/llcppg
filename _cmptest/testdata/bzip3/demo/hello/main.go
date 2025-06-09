package main

import (
	"fmt"

	zip "bzip3"

	"github.com/goplus/lib/c"
)

func main() {
	fmt.Println(c.GoString(zip.Version()))
	input := []byte("Hello, bzip3 compression!")
	output := make([]byte, zip.Bound(uintptr(len(input))))
	outputSize := uintptr(len(output))

	errCode := zip.Compress(1024*1024, &input[0], &output[0], uintptr(len(input)), &outputSize)
	if errCode != zip.OK {
		fmt.Println("Compression failed with error code:", errCode)
		return
	}
	fmt.Println("Compression successful. Compressed size:", outputSize)

	decompressed := make([]byte, len(input))
	decompressedSize := uintptr(len(decompressed))
	errCode = zip.Decompress(&output[0], &decompressed[0], outputSize, &decompressedSize)
	if errCode != zip.OK {
		fmt.Println("Decompression failed with error code:", errCode)
		return
	}
	fmt.Println("Decompression successful. Decompressed data:", string(decompressed))
}
