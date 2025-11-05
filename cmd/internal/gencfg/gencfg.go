package gencfg

import (
	"fmt"

	"github.com/goplus/llcppg/cmd/internal/base"
)

var Cmd = &base.Command{
	UsageLine: "llcppg gencfg",
	Short:     "generate llcpp.cfg",
}

func init() {
	Cmd.Run = runCmd
}

func runCmd(cmd *base.Command, args []string) {
	fmt.Printf("todo generate llcpp.cfg")
}
