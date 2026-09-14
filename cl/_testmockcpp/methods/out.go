package foo

import "github.com/goplus/lib/c"

type Bar struct {
}

// llgo:link (*Bar).F C._ZN3Bar1fEi
func (this *Bar) F(a c.Int) c.Uint {
	return 0
}

// llgo:link (*Bar).Xprintf C._ZN3Bar7xprintfEPKcz
func (this *Bar) Xprintf(fmt *c.Char, __llgo_va_list ...any) c.Int {
	return 0
}
