package libxml2

import (
	_ "github.com/goplus/lib/c"
	_ "github.com/goplus/lib/c/os"
)

const LLGoPackage string = "link: $(pkg-config --libs libxml-2.0);"
