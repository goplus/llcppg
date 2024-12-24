package main

import (
	test "github.com/goplus/llcppg/_xtool/llcppsigfetch/parse/cvt_test"
	"github.com/goplus/llcppg/_xtool/llcppsymg/clangutils"
)

func main() {
	TestInclude()
}

func TestInclude() {
	test.RunTestWithConfig(&clangutils.Config{
		File:  "./hfile/temp.h",
		Temp:  false,
		IsCpp: false,
	})
}
