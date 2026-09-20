package foo

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

const XGoPackage = true

type Bar struct {
}

//go:linkname Bar_create__0 C._ZN3Bar6createEv
func Bar_create__0() c.Int

//go:linkname Bar_create__1 C._ZN3Bar6createEi
func Bar_create__1(a c.Int) c.Int

// llgo:link (*Bar).F C._ZN3Bar1fEi
func (this *Bar) F(a c.Int) c.Uint {
	return 0
}
