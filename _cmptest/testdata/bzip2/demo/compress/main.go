package main

import (
	"fmt"
	"unsafe"

	"bzip2"

	"github.com/goplus/lib/c"
)

func CompressFile(inPath, outPath string) error {
	inPathC := c.AllocaCStr(inPath)

	inFile := c.Fopen(inPathC, c.Str("rb"))
	if inFile == nil {
		return fmt.Errorf("failed to open input file: %s", inPath)
	}
	defer c.Fclose(inFile)

	outPathC := c.AllocaCStr(outPath)
	outFile := c.Fopen(outPathC, c.Str("wb"))
	if outFile == nil {
		return fmt.Errorf("failed to open output file: %s", outPath)
	}
	defer c.Fclose(outFile)

	var bzerr c.Int
	bzfile := bzip2.WriteOpen(&bzerr, outFile, 9, 0, 30)
	if bzfile == nil || bzerr != bzip2.OK {
		return fmt.Errorf("BzWriteOpen error, code=%d", bzerr)
	}

	buf := make([]byte, 4096)
	for {
		n := c.Fread(unsafe.Pointer(&buf[0]), 1, uintptr(len(buf)), inFile)
		if n == 0 {
			break
		}

		bzip2.Write(&bzerr, bzfile, unsafe.Pointer(&buf[0]), c.Int(n))
		if bzerr != bzip2.OK {
			return fmt.Errorf("BzWrite error, code=%d", bzerr)
		}

		if n < uintptr(len(buf)) {
			break
		}
	}

	bzip2.WriteClose(&bzerr, bzfile, 0, nil, nil)
	if bzerr != bzip2.OK {
		return fmt.Errorf("BzWriteClose error, code=%d", bzerr)
	}

	return nil
}

func main() {
	CompressFile("test.txt", "../decompress/test.bz2")
}
