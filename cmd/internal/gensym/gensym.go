package gensym

import (
	"github.com/goplus/llcppg/cmd/internal/base"
)

var Cmd = &base.Command{
	UsageLine: "llcppg gensym",
	Short:     "generate symbol table for a C/C++ library",
}

func init() {
	Cmd.Run = runCmd
	addFlags(&Cmd.Flag)
}

func runCmd(cmd *base.Command, args []string) {
	base.RunCmdWithName(cmd, args, "llcppsymg", nil)
}
