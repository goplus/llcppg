package gen

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

func GetOut(cmd *exec.Cmd, dir string) (string, error) {
	outBuf := new(bytes.Buffer)
	errBuff := new(bytes.Buffer)
	cmd.Stdin = os.Stdin
	cmd.Stdout = outBuf
	cmd.Stderr = errBuff
	cmd.Env = os.Environ()
	if len(dir) > 0 {
		cmd.Dir = dir
	}
	err := cmd.Run()
	if err != nil {
		return errBuff.String(), err
	}
	return outBuf.String(), nil
}

func NewExecCommand(cmdStr string, args ...string) *exec.Cmd {
	cmdStr = strings.TrimSpace(cmdStr)
	return exec.Command(cmdStr, args...)
}

func ExpandString(str string, dir string) string {
	str = strings.ReplaceAll(str, "(", "{")
	str = strings.ReplaceAll(str, ")", "}")
	expandStr := os.Expand(str, func(s string) string {
		args := strings.Fields(s)
		if len(args) == 0 {
			return ""
		}
		outString, err := GetOut(NewExecCommand(args[0], args[1:]...), dir)
		if err != nil {
			return ""
		}
		return outString
	})
	return strings.TrimSpace(expandStr)
}
