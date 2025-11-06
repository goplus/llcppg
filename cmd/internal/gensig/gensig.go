package gensig

import (
	"github.com/goplus/llcppg/cmd/internal/base"
)

var Cmd = &base.Command{
	UsageLine: "llcppg gensig",
	Short:     "generate signature information of C/C++ symbols",
}

func init() {
	Cmd.Run = runCmd
	addFlags(&Cmd.Flag)
}

func runCmd(cmd *base.Command, args []string) {
	base.RunCmdWithName(cmd, args, "llcppsigfetch", nil)
}
