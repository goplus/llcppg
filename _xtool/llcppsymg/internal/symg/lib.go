package symg

import (
	"fmt"
	"strings"

	"github.com/goplus/llcppg/_xtool/internal/symbol"
)

type Libs struct {
	Paths []string
	Names []string
}

func ParseLibs(libs string) *Libs {
	parts := strings.Fields(libs)
	lbs := &Libs{}
	for _, part := range parts {
		if strings.HasPrefix(part, "-L") {
			lbs.Paths = append(lbs.Paths, part[2:])
		} else if strings.HasPrefix(part, "-l") {
			lbs.Names = append(lbs.Names, part[2:])
		}
	}
	return lbs
}

type LibMode = symbol.Mode

// searches for each library name in the provided paths and default paths,
// appending the appropriate file extension (.dylib for macOS, .so for Linux).
//
// Example: For "-L/opt/homebrew/lib -llua -lm":
// - It will search for liblua.dylib (on macOS) or liblua.so (on Linux)
// - System libs like -lm are ignored and included in notFound
//
// So error is returned if no libraries found at all.
func (l *Libs) Files(findPaths []string, mode LibMode) ([]string, []string, error) {
	var foundPaths []string
	var notFound []string
	searchPaths := append(l.Paths, findPaths...)
	for _, name := range l.Names {
		var foundPath string
		for _, path := range searchPaths {
			libPath, err := symbol.FindLibFile(path, name, mode)
			if err != nil {
				continue
			}
			if libPath != "" {
				foundPath = libPath
				break
			}
		}
		if foundPath != "" {
			foundPaths = append(foundPaths, foundPath)
		} else {
			notFound = append(notFound, name)
		}
	}
	if len(foundPaths) == 0 {
		return nil, notFound, fmt.Errorf("failed to find any libraries")
	}
	return foundPaths, notFound, nil
}
