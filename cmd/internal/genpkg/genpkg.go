package genpkg

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/goplus/llcppg/cmd/internal/base"
	"github.com/goplus/llcppg/config"
	"golang.org/x/mod/module"
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

	err = module.CheckPath(modulePath)
	base.Check(err)

	cfgFile := config.LLCPPG_CFG

	config.HandleMarshalConfigFile(cfgFile, func(b []byte, err error) {

		base.Check(err)

		r, w := io.Pipe()

		errCh := make(chan error, 1)
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
			errCh <- llcppsigfetchCmd.Run()
		}()

		gogensigCmdArgs := make([]string, 0)
		gogensigCmdArgs = append(gogensigCmdArgs, fmt.Sprintf("-mod=%s", modulePath))
		if verbose {
			gogensigCmdArgs = append(gogensigCmdArgs, "-v")
		}
		gogensigCmd := exec.Command("gogensig", gogensigCmdArgs...)
		gogensigCmd.Stdin = r
		gogensigCmd.Stdout = os.Stdout
		gogensigCmd.Stderr = os.Stderr
		err = gogensigCmd.Run()
		base.Check(err)

		if fetchErr := <-errCh; fetchErr != nil {
			base.Check(fetchErr)
		}
	})
}
