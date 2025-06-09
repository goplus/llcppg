package main

import (
	"fmt"

	"libtool"

	"github.com/goplus/lib/c"
)

func main() {
	fmt.Println("Simple libtool demonstration")

	// Initialize libtool
	ret := libtool.Dlinit()
	if ret != 0 {
		fmt.Println("Failed to initialize libtool:", c.GoString(libtool.Dlerror()))
		return
	}
	fmt.Println("Successfully initialized libtool")

	// Try to load a common library (libc)
	libName := "libc.so.6" // Linux style
	handle := libtool.Dlopen(c.Str(libName))
	if handle == nil {
		libName = "libc.dylib" // macOS style
		handle = libtool.Dlopen(c.Str(libName))
	}
	if handle == nil {
		libName = "c" // Generic style
		handle = libtool.Dlopen(c.Str(libName))
	}

	if handle != nil {
		fmt.Printf("Successfully opened %s\n", libName)

		// Try to find a common function (printf)
		symPtr := libtool.Dlsym(handle, c.Str("printf"))
		if symPtr != nil {
			fmt.Println("Found 'printf' function")
		} else {
			fmt.Println("Symbol 'printf' not found:", c.GoString(libtool.Dlerror()))
		}

		// Close the library
		libtool.Dlclose(handle)
		fmt.Println("Closed library")
	} else {
		fmt.Println("Could not open any standard library:", c.GoString(libtool.Dlerror()))
	}

	// Clean up libtool
	libtool.Dlexit()
	fmt.Println("Successfully cleaned up libtool")
}
