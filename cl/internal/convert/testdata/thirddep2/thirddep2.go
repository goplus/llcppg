package thirddep2

import (
	_ "unsafe"

	"github.com/goplus/llcppg/cl/internal/convert/testdata/basicdep"
	"github.com/goplus/llcppg/cl/internal/convert/testdata/thirddep"
)

type ThirdDep2 struct {
	A TypeThirdDep2
	B TypeThirdDep2
	C thirddep.ThirdDep
	D basicdep.BasicDep
}
