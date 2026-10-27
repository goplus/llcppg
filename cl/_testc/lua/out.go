package foo

import "github.com/goplus/lib/c"

const LUA_REGISTRYINDEX = -1073742823

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

//llgo:link (*State).Getfield C.lua_getfield
func (L *State) Getfield(idx c.Int, k *c.Char) c.Int {
	return 0
}
func (L *State) LLoadfile(f *c.Char) c.Int {
	return L.LLoadfilex(f, nil)
}
func (L *State) LCheckunsigned(a c.Int) Integer {
	return L.LCheckinteger(a)
}
func (L *State) LOptunsigned(a c.Int, d Integer) Integer {
	return L.LOptinteger(a, d)
}
func (L *State) LCheckint(n c.Int) Integer {
	return L.LCheckinteger(n)
}
func (L *State) LOptint(n c.Int, d Integer) Integer {
	return L.LOptinteger(n, d)
}
func (L *State) LChecklong(n c.Int) Integer {
	return L.LCheckinteger(n)
}
func (L *State) LOptlong(n c.Int, d Integer) Integer {
	return L.LOptinteger(n, d)
}
