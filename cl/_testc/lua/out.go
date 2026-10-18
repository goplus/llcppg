package foo

import "github.com/goplus/lib/c"

type LBuffer struct {
	B    *c.Char
	Size SizeT
	N    SizeT
}
type SizeT = c.Ulong

// llgo:link (*LBuffer).Prepbuffsize C.luaL_prepbuffsize
func (B *LBuffer) Prepbuffsize(sz SizeT) *c.Char {
	return nil
}

// llgo:link (*LBuffer).Addlstring C.luaL_addlstring
func (B *LBuffer) Addlstring(s *c.Char, l SizeT) {
}

// llgo:link (*LBuffer).MAddlstring C.luaM_addlstring
func (B *LBuffer) MAddlstring(s *c.Char, l SizeT) {
}
