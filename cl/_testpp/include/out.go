package foo

type BarBase struct {
}
type BarDerived struct {
	b BarBase
}

// llgo:link (*BarBase).XGo_Dtor C._ZN3bar4baseD1Ev
func (this *BarBase) XGo_Dtor() {
}
