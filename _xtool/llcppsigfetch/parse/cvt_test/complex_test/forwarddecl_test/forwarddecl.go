package main

import (
	test "github.com/goplus/llcppg/_xtool/llcppsigfetch/parse/cvt_test"
	"github.com/goplus/llcppg/llcppg"
)

func main() {
	TestForwardDecl()
	TestForwardDeclCrossFile()
}

func TestForwardDecl() {
	test.RunTestWithConfig(&llcppg.Config{
		Include: []string{"forwarddecl.h"},
		CFlags:  "-I./hfile",
	})
}

func TestForwardDeclCrossFile() {
	test.RunTestWithConfig(&llcppg.Config{
		Include: []string{"def.h"},
		CFlags:  "-I./hfile",
	})
}
