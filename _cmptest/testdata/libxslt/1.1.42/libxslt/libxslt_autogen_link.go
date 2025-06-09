package libxslt

import (
	_ "github.com/goplus/lib/c"
	_ "github.com/goplus/lib/c/os"
	_ "github.com/goplus/llcppg/_cmptest/testdata/libxml2/2.13.6/libxml2"
)

const LLGoPackage string = "link: $(pkg-config --libs libxslt);"
