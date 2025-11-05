package gensig

import (
	"fmt"

	"github.com/goplus/llcppg/cmd/internal/base"
)

var Cmd = &base.Command{
	UsageLine: "llcppg gensig",
	Short:     "generate signature information of C/C++ symbols",
}

func init() {
	Cmd.Run = runCmd
}

func runCmd(cmd *base.Command, args []string) {
	fmt.Printf("todo gensig")
}
