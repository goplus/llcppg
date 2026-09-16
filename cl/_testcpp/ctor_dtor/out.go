package foo

import "github.com/goplus/lib/c"

const XGoPackage = true
const (
	LLGoPackage = "link: -L/path/foo -lfoo"
	LLGoFiles   = "-I/path/foo/include: _wrap/llcppg.cpp"
)

type Base struct {
}

// llgo:link (*Base).XGo_Dtor C._ZN4baseD1Ev
func (this *Base) XGo_Dtor() {
}

type Bar struct {
}

// llgo:link (*Bar).XGo_Ctor__1 C._ZN3barC1Ei
func (this *Bar) XGo_Ctor__1(a c.Int) {
}

// llgo:link (*Bar).XGo_Ctor__0 C._llcppg__ZN3barC1Ev
func (this *Bar) XGo_Ctor__0() {
}

// llgo:link (*Bar).XGo_Ctor__2 C._llcppg__ZN3barC1EPKc
func (this *Bar) XGo_Ctor__2(a *c.Char) {
}

// llgo:link (*Bar).XGo_Dtor C._llcppg__ZN3barD1Ev
func (this *Bar) XGo_Dtor() {
}
