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

type Bar struct {
}

//go:linkname BarCreate__0 C._llcppg__ZN3Bar6createEv
func BarCreate__0() c.Int

//go:linkname BarCreate__1 C._ZN3Bar6createEi
func BarCreate__1(a c.Int) c.Int

// llgo:link (*Bar).F C._ZN3Bar1fEi
func (this *Bar) F(a c.Int) c.Uint {
	return 0
}
