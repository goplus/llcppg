package foo

import "github.com/goplus/lib/c"

const (
	LLGoPackage = "link: -L/path/foo -lfoo"
	LLGoFiles   = "-I/path/foo/include: _wrap/foo.cpp"
)

type Bar struct {
}

// llgo:link (*Bar).F C._llcppg__ZN3bar1fEi
func (this *Bar) F(a c.Int) c.Uint {
	return 0
}

// llgo:link (*Bar).X_g C._llcppg__ZN3bar2_gEv
func (this *Bar) X_g() {
}

// llgo:link (*Bar).Xprintf C._ZN3bar7xprintfEPKcz
func (this *Bar) Xprintf(fmt *c.Char, __llgo_va_list ...any) c.Int {
	return 0
}
