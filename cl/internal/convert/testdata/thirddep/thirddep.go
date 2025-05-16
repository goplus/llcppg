package thirddep

import (
	_ "unsafe"

	"github.com/goplus/lib/c"
)

type ThirdDep struct {
	A TypeThirdDep
	B TypeThirdDep
}

type Stream c.Long
