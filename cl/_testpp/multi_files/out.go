package foo

const LLGoPackage = "link: -L/path/foo -lfoo"

type BarBase struct {
}
type BarDerived struct {
	b BarBase
}
type BarDetailBar struct {
	BarBase
}

// llgo:link (*BarBase).XGo_Dtor C._ZN3bar4baseD1Ev
func (this *BarBase) XGo_Dtor() {
}
