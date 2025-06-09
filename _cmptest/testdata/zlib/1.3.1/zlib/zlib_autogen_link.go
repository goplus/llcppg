package zlib

import (
	_ "github.com/goplus/lib/c"
	_ "github.com/goplus/lib/c/os"
)

const LLGoPackage string = "link: $(pkg-config --libs zlib);"
