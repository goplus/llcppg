package genpkg

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/goplus/llcppg/cmd/internal/base"
	"github.com/goplus/llcppg/config"
)

var Cmd = &base.Command{
	UsageLine: "llcppg genpkg",
	Short:     "generate a go package by signature information of symbols",
}

func init() {
	Cmd.Run = runCmd
	addFlags(&Cmd.Flag)
}

func runCmd(cmd *base.Command, args []string) {
	err := cmd.Flag.Parse(args)
	if err != nil {
		return
	}

	cfgFile := config.LLCPPG_CFG

	config.HandleMarshalConfigFile(cfgFile, func(b []byte, err error) {

		base.Check(err)

		r, w := io.Pipe()

		go func() {
			defer w.Close()
			llcppsigfetchCmdArgs := make([]string, 0)
			if verbose {
				llcppsigfetchCmdArgs = append(llcppsigfetchCmdArgs, "-v")
			}
			if cmd.Flag.NArg() == 0 {
				llcppsigfetchCmdArgs = append(llcppsigfetchCmdArgs, "-")
			}
			llcppsigfetchCmd := exec.Command("llcppsigfetch", llcppsigfetchCmdArgs...)
			llcppsigfetchCmd.Stdin = bytes.NewReader(b)
			llcppsigfetchCmd.Stdout = w
			llcppsigfetchCmd.Stderr = os.Stderr
			err = llcppsigfetchCmd.Run()
			base.Check(err)
		}()

		gogensigCmdArgs := make([]string, 0)
		if len(modulePath) > 0 {
			gogensigCmdArgs = append(gogensigCmdArgs, fmt.Sprintf("-mod=%s", modulePath))
		}
		if verbose {
			gogensigCmdArgs = append(gogensigCmdArgs, "-v")
		}
		gogensigCmd := exec.Command("gogensig", gogensigCmdArgs...)
		gogensigCmd.Stdin = r
		gogensigCmd.Stdout = os.Stdout
		gogensigCmd.Stderr = os.Stderr
		err = gogensigCmd.Run()
		base.Check(err)
	})
}
