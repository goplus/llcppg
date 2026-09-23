package foo

const LLGoPackage = "link: -L/path/foo -lfoo"

type Bar_base struct {
}
type Bar_derived struct {
	b Bar_base
}
type Bar_detail_Bar struct {
	Bar_base
}

// llgo:link (*Bar_base).XGo_Dtor C._ZN3bar4baseD1Ev
func (this *Bar_base) XGo_Dtor() {
}
