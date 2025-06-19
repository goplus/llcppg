package clangtool

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var _sysRootDirOnce = sync.OnceValues(sysRoot)

var _matchISysrootRegex = regexp.MustCompile(`-(resource-dir|internal-isystem|isysroot|internal-externc-isystem)\s(\S+)`)

// ComposeIncludes create Include list
// #include <file1.h>
// #include <file2.h>
func ComposeIncludes(files []string, outfile string) error {
	var str string
	for _, file := range files {
		str += ("#include <" + file + ">\n")
	}
	return os.WriteFile(outfile, []byte(str), 0644)
}

func GetIncludePaths(isCpp bool) []string {
	args := []string{"-E", "-v"}
	args = append(args, defaultArgs(isCpp)...)
	args = append(args, "/dev/null")
	cmd := exec.Command("clang", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	return ParseClangIncOutput(string(output))
}

func ParseClangIncOutput(output string) []string {
	var paths []string
	start := strings.Index(output, "#include <...> search starts here:")
	end := strings.Index(output, "End of search list.")
	if start == -1 || end == -1 {
		return paths
	}
	content := output[start:end]
	lines := strings.Split(content, "\n")
	for _, line := range lines[1:] {
		for _, item := range strings.Fields(line) {
			if path := strings.TrimSpace(item); filepath.IsAbs(path) {
				paths = append(paths, path)
			}
		}
	}
	return paths
}

func WithSysRoot(args []string) []string {
	_defaultSysRootDir, _ := _sysRootDirOnce()
	return append(args, _defaultSysRootDir...)
}

func defaultArgs(isCpp bool) []string {
	args := []string{"-x", "c"}
	if isCpp {
		args = []string{"-x", "c++"}
	}
	return args
}

// sysRoot retrieves isysroot from clang preprocessor
func sysRoot() ([]string, error) {
	var output bytes.Buffer

	// -x dones't matter, we don't care, just get the isysroot
	cmd := exec.Command("clang", "-E", "-v", "-x", "c", "/dev/null")
	cmd.Stderr = &output

	cmd.Run()

	return ParseSystemPath(output.String())
}

func ParseSystemPath(output string) ([]string, error) {
	sysRootResults := _matchISysrootRegex.FindAllStringSubmatch(output, -1)

	var result []string

	for _, sysRootResult := range sysRootResults {
		if len(sysRootResult) == 3 {
			if sysRootResult[1] == "resource-dir" {
				// the format of resource-dir must be -resource-dir=/xxx
				result = append(result, fmt.Sprintf("-%s=%s", sysRootResult[1], sysRootResult[2]))
				// append its header path also
				result = append(result, fmt.Sprintf("-I%s", filepath.Join(sysRootResult[2], "include")))
				continue
			}
			result = append(result, fmt.Sprintf("-%s%s", sysRootResult[1], sysRootResult[2]))
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("failed to find any sysRoot path")
	}

	return result, nil
}
