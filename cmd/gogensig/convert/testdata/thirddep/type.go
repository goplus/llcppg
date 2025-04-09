package thirddep

import (
	_ "unsafe"

	"github.com/goplus/lib/c"
)

type TypeThirdDep struct {
	A c.Int
	B c.Int
}
