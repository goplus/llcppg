package gensym

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/goplus/llcppg/cmd/internal/base"
	"github.com/goplus/llcppg/config"
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
	err := cmd.Flag.Parse(args)
	check(err)

	cfgFile := config.LLCPPG_CFG
	bytesOfConf, err := config.MarshalConfigFile(cfgFile)
	check(err)

	if cmd.Flag.NArg() == 0 {
		args = append(args, "-")
	}

	cmdForLlcppsymg := exec.Command("llcppsymg", args...)
	cmdForLlcppsymg.Stdin = bytes.NewReader(bytesOfConf)
	cmdForLlcppsymg.Stdout = os.Stdout
	cmdForLlcppsymg.Stderr = os.Stderr
	cmdForLlcppsymg.Run()
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
