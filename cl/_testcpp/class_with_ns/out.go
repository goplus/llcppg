package foo

type Bar_base struct {
}
type Bar_derived struct {
	b Bar_base
}

// llgo:link (*Bar_base).XGo_Dtor C._ZN3bar4baseD1Ev
func (this *Bar_base) XGo_Dtor() {
}
