package main

import (
	"fmt"
	"log"
	"os"
	"unsafe"

	"bzip2"

	"github.com/goplus/lib/c"
)

func DecompressFile(inPath, outPath string) error {
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
	bzfile := bzip2.ReadOpen(&bzerr, inFile, 0, 0, nil, 0)
	if bzfile == nil || bzerr != bzip2.OK {
		return fmt.Errorf("BzReadOpen error, code=%d", bzerr)
	}

	buf := make([]byte, 4096)
	for {
		n := bzip2.Read(&bzerr, bzfile, unsafe.Pointer(&buf[0]), c.Int(len(buf)))
		if bzerr == bzip2.STREAM_END {
			if n > 0 {
				c.Fwrite(unsafe.Pointer(&buf[0]), 1, uintptr(n), outFile)
			}
			break
		}

		if bzerr != bzip2.OK && bzerr != bzip2.STREAM_END {
			return fmt.Errorf("BzRead error, code=%d", bzerr)
		}
		if n > 0 {
			c.Fwrite(unsafe.Pointer(&buf[0]), 1, uintptr(n), outFile)
		} else {
			break
		}
	}

	bzip2.ReadClose(&bzerr, bzfile)
	if bzerr != bzip2.OK {
		return fmt.Errorf("BzReadClose error, code=%d", bzerr)
	}

	return nil
}

func main() {
	err := DecompressFile("test.bz2", "ttt.test")
	if err != nil {
		panic(err)
	}

	b, err := os.ReadFile("ttt.test")
	if err != nil {
		panic(err)
	}

	log.Println(string(b))
}
