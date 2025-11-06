package gencfg

import (
	"os"
	"strings"

	"github.com/goplus/llcppg/cmd/internal/base"
	"github.com/goplus/llcppg/cmd/llcppcfg/gen"
	"github.com/goplus/llcppg/config"
)

var Cmd = &base.Command{
	UsageLine: "llcppg init",
	Short:     "init llcppg.cfg config file",
}

func init() {
	Cmd.Run = runCmd
	addFlags(&Cmd.Flag)
}

func runCmd(cmd *base.Command, args []string) {

	if err := cmd.Flag.Parse(args); err != nil {
		return
	}

	name := ""
	if len(cmd.Flag.Args()) > 0 {
		name = cmd.Flag.Arg(0)
	}

	exts := strings.Fields(extsString)
	deps := strings.Fields(dependencies)

	excludeSubdirs := []string{}
	if len(excludes) > 0 {
		excludeSubdirs = strings.Fields(excludes)
	}
	var flagMode gen.FlagMode
	if cpp {
		flagMode |= gen.WithCpp
	}
	if tab {
		flagMode |= gen.WithTab
	}
	buf, err := gen.Do(gen.NewConfig(name, flagMode, exts, deps, excludeSubdirs))
	if err != nil {
		panic(err)
	}
	outFile := config.LLCPPG_CFG
	err = os.WriteFile(outFile, buf, 0600)
	if err != nil {
		panic(err)
	}
}
