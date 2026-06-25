//go:build !darwin

package main

import "lua"

func openlibs(L *lua.State) {
	L.Openlibs()
}
