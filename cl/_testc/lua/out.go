package foo

import "github.com/goplus/lib/c"

type LBuffer struct {
	B    *c.Char
	Size SizeT
	N    SizeT
}
type State struct {
}
type SizeT = c.Ulong

// type for integer functions
type Integer = c.Int

// unsigned integer type
type Unsigned = c.Uint

//llgo:link (*LBuffer).Prepbuffsize C.luaL_prepbuffsize
func (B *LBuffer) Prepbuffsize(sz SizeT) *c.Char {
	return nil
}

//llgo:link (*LBuffer).Addlstring C.luaL_addlstring
func (B *LBuffer) Addlstring(s *c.Char, l SizeT) {
}

//llgo:link (*LBuffer).MAddlstring C.luaM_addlstring
func (B *LBuffer) MAddlstring(s *c.Char, l SizeT) {
}

//llgo:link (*State).LCheckinteger C.luaL_checkinteger
func (L *State) LCheckinteger(arg c.Int) Integer {
	return 0
}

//llgo:link (*State).LOptinteger C.luaL_optinteger
func (L *State) LOptinteger(arg c.Int, def Integer) Integer {
	return 0
}

//llgo:link (*State).LLoadfilex C.luaL_loadfilex
func (L *State) LLoadfilex(filename *c.Char, mode *c.Char) c.Int {
	return 0
}
func (L *State) LLoadfile(f *c.Char) c.Int {
	return L.LLoadfilex(f, nil)
}
