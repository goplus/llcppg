package foo

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

const XGoPackage = true

//go:linkname Bar_detail_f__1 C._ZN3bar6detail1fEi
func Bar_detail_f__1(a c.Int) c.Uint

//go:linkname Bar_detail_f__0 C._ZN3bar6detail1fEv
func Bar_detail_f__0()

type Bar_base struct {
}

// llgo:link (*Bar_base).XGo_Dtor C._ZN3bar4baseD1Ev
func (this *Bar_base) XGo_Dtor() {
}
