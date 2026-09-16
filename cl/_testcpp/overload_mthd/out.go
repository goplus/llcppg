package foo

import "github.com/goplus/lib/c"

const XGoPackage = true

type Bar struct {
}

// llgo:link (*Bar).F__1 C._ZN3bar1fEi
func (this *Bar) F__1(a c.Int) c.Uint {
	return 0
}

// llgo:link (*Bar).F__0 C._ZN3bar1fEv
func (this *Bar) F__0() {
}
