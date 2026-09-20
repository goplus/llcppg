package foo

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

const XGoPackage = true
const (
	LLGoPackage = "link: -L/path/foo -lfoo"
	LLGoFiles   = "-I/path/foo/include: _wrap/llcppg.cpp"
)

type Bar_base struct {
}
type Bar_boolean = c.Char

//go:linkname Bar_detail_f__1 C._ZN3bar6detail1fEi
func Bar_detail_f__1(a c.Int) c.Uint

//go:linkname Bar_detail_f__0 C._llcppg__ZN3bar6detail1fEv
func Bar_detail_f__0()

// llgo:link (*Bar_base).XGo_Dtor C._ZN3bar4baseD1Ev
func (this *Bar_base) XGo_Dtor() {
}

//go:linkname Bar_print C._llcppg__ZN3bar5printENS_4baseE
func Bar_print(b Bar_base)
