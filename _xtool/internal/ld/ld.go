package ld

import (
	"os/exec"
	"regexp"
	"runtime"
)

// GetLibSearchPaths returns the library paths from the ld command.
// With linux, it will use ld --verbose to get the library paths.
func GetLibSearchPaths() []string {
	var paths []string
	if runtime.GOOS == "linux" {
		//resolution from https://github.com/goplus/llcppg/commit/02307485db9269481297a4dc5e8449fffaa4f562
		cmd := exec.Command("ld", "--verbose")
		output, err := cmd.Output()
		if err != nil {
			panic(err)
		}
		return ParseOutput(string(output))
	}
	return paths
}

// ParseOutput parses the output of the ld command.
// It returns the search library paths from the ld command.
func ParseOutput(output string) []string {
	var paths []string
	matches := regexp.MustCompile(`SEARCH_DIR\("=([^"]+)"\)`).FindAllStringSubmatch(output, -1)
	for _, match := range matches {
		paths = append(paths, match[1])
	}
	return paths
}
