package thirddep

import (
	_ "unsafe"

	"github.com/goplus/llgo/c"
)

type ThirdDep struct {
	A TypeThirdDep
	B TypeThirdDep
}

type Stream c.Long
