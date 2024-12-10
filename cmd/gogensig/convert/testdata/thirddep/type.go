package thirddep

import (
	_ "unsafe"

	"github.com/goplus/llgo/c"
)

type TypeThirdDep struct {
	A c.Int
	B c.Int
}
