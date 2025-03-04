/*
 * @Author: Zhang Zhi Yang
 * @Date: 2025-03-04 11:16:39
 * @LastEditors: Zhang Zhi Yang
 * @LastEditTime: 2025-03-04 11:20:09
 * @FilePath: /llcppg/_xtool/llcppsigfetch/parse/cvt_test/complex_test/forwarddecl_test/forwarddecl.go
 * @Description:
 */
package main

import (
	"github.com/goplus/llcppg/_xtool/llcppsigfetch/parse"
	test "github.com/goplus/llcppg/_xtool/llcppsigfetch/parse/cvt_test"
	"github.com/goplus/llcppg/llcppg"
)

func main() {
	TestForwardDecl()
	TestForwardDeclCrossFile()
}

func TestForwardDecl() {
	test.RunTestWithConfig(&parse.ParseConfig{
		Conf: &llcppg.Config{
			Include: []string{"forwarddecl.h"},
			CFlags:  "-I./hfile/",
		},
	})
}

func TestForwardDeclCrossFile() {
	test.RunTestWithConfig(&parse.ParseConfig{
		Conf: &llcppg.Config{
			Include: []string{"def.h"},
			CFlags:  "-I./hfile/",
		},
	})
}
