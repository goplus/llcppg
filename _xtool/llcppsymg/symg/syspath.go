package symg

import (
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

func GetLibPaths() []string {
	var paths []string
	if runtime.GOOS == "linux" {
		//resolution from https://github.com/goplus/llcppg/commit/02307485db9269481297a4dc5e8449fffaa4f562
		cmd := exec.Command("ld", "--verbose")
		output, err := cmd.Output()
		if err != nil {
			panic(err)
		}
		return ParseLdOutput(string(output))
	}
	return paths
}

// Note:this function is only use in this package
// The public function name is for llgo test
func ParseLdOutput(output string) []string {
	var paths []string
	matches := regexp.MustCompile(`SEARCH_DIR\("=([^"]+)"\)`).FindAllStringSubmatch(output, -1)
	for _, match := range matches {
		paths = append(paths, match[1])
	}
	return paths
}

func GetIncludePaths() []string {
	var paths []string
	if runtime.GOOS == "linux" {
		cmd := exec.Command("clang", "-E", "-v", "-x", "c", "/dev/null")
		output, err := cmd.CombinedOutput()
		if err != nil {
			panic(err)
		}
		return ParseClangIncOutput(string(output))
	}
	return paths
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
		if path := strings.TrimSpace(line); path != "" {
			paths = append(paths, path)
		}
	}
	return paths
}
