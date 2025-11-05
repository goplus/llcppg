package gensym

import (
	"fmt"

	"github.com/goplus/llcppg/cmd/internal/base"
)

var Cmd = &base.Command{
	UsageLine: "llcppg gensym",
	Short:     "generate symbol table for a C/C++ library",
}

func init() {
	Cmd.Run = runCmd
}

func runCmd(cmd *base.Command, args []string) {
	fmt.Println("todo gensym")
}
