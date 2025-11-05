package genpkg

import (
	"fmt"

	"github.com/goplus/llcppg/cmd/internal/base"
)

var Cmd = &base.Command{
	UsageLine: "llcppg genpkg",
	Short:     "generate a go package by signature information of symbols",
}

func init() {
	Cmd.Run = runCmd
}

func runCmd(cmd *base.Command, args []string) {
	fmt.Println("todo genpkg")
}
