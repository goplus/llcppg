package version

import (
	"fmt"

	"github.com/goplus/llcppg/cmd/internal/base"
)

var Cmd = &base.Command{
	UsageLine: "llcppg version",
	Short:     "Print llcppg version",
}

func init() {
	Cmd.Run = runCmd
}

func runCmd(cmd *base.Command, args []string) {
	fmt.Printf("todo print version")
}
