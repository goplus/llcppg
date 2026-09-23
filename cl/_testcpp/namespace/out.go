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

type BarBase struct {
}
type BarBoolean = c.Char

//go:linkname BarDetailF__1 C._ZN3bar6detail1fEi
func BarDetailF__1(a c.Int) c.Uint

//go:linkname BarDetailF__0 C._llcppg__ZN3bar6detail1fEv
func BarDetailF__0()

// llgo:link (*BarBase).XGo_Dtor C._llcppg__ZN3bar4baseD1Ev
func (this *BarBase) XGo_Dtor() {
}

//go:linkname BarPrint C._llcppg__ZN3bar5printENS_4baseE
func BarPrint(b BarBase)
