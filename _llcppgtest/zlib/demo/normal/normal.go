package main

import (
	"unsafe"
	"zlib"

	"github.com/goplus/lib/c"
)

func main() {
	txt := []byte("zlib is a software library used for data compression. It was created by Jean-loup Gailly and Mark Adler and first released in 1995. zlib is designed to be a free, legally unencumbered—that is, not covered by any patents—alternative to the proprietary DEFLATE compression algorithm, which is often used in software applications for data compression.The library provides functions to compress and decompress data using the DEFLATE algorithm, which is a combination of the LZ77 algorithm and Huffman coding. zlib is notable for its versatility; it can be used in a wide range of applications, from web servers and web clients compressing HTTP data, to the compression of data for storage or transmission in various file formats, such as PNG, ZIP, and GZIP.")
	txtLen := zlib.ULong(len(txt))

	cmpSize := zlib.ULongf(zlib.CompressBound(txtLen))
	cmpData := make([]byte, int(cmpSize))
	data := (*zlib.Bytef)(unsafe.Pointer(unsafe.SliceData(cmpData)))
	txtData := (*zlib.Bytef)(unsafe.Pointer(unsafe.SliceData(txt)))

	res := zlib.Compress(data, &cmpSize, txtData, txtLen)
	if res != zlib.OK {
		c.Printf(c.Str("\nCompression failed: %d\n"), res)
		return
	}

	c.Printf(c.Str("Text length = %d, Compressed size = %d\n"), txtLen, cmpSize)

	ucmpSize := zlib.ULongf(txtLen)
	ucmp := make([]byte, int(ucmpSize))
	ucmpPtr := (*zlib.Bytef)(unsafe.Pointer(unsafe.SliceData(ucmp)))

	unRes := zlib.Uncompress(ucmpPtr, &ucmpSize, data, zlib.ULong(cmpSize))
	c.Printf(c.Str("Decompression result = %d, Decompressed size %d\n"), unRes, ucmpSize)

	if unRes != zlib.OK {
		c.Printf(c.Str("\nDecompression failed: %d\n"), unRes)
		return
	}

	c.Printf(c.Str("Decompressed data: \n"))
	for i := 0; i < int(ucmpSize); i++ {
		c.Printf(c.Str("%c"), ucmp[i])
	}
}
