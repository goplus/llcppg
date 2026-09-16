package foo

import "github.com/goplus/lib/c"

const XGoPackage = true

type Bar struct {
}

// llgo:link (*Bar).XGo_Ctor__1 C._ZN3barC1Ei
func (this *Bar) XGo_Ctor__1(a c.Int) {
}

// llgo:link (*Bar).XGo_Ctor__0 C._ZN3barC1Ev
func (this *Bar) XGo_Ctor__0() {
}

// llgo:link (*Bar).XGo_Dtor C._ZN3barD1Ev
func (this *Bar) XGo_Dtor() {
}
