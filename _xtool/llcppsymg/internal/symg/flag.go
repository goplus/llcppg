package symg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goplus/llcppg/_xtool/internal/symbol"
)

// Note: This package is not placed under the 'config' package because 'config'
// depends on 'cjson'. The parsing of Libs and cflags is intended to be usable
// by both llgo and go, without introducing additional dependencies.

type Libs struct {
	Paths []string // Dylib Path
	Names []string
}

type CFlags struct {
	Paths []string // Include Path
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

func ParseCFlags(cflags string) *CFlags {
	parts := strings.Fields(cflags)
	cf := &CFlags{}
	for _, part := range parts {
		if strings.HasPrefix(part, "-I") {
			cf.Paths = append(cf.Paths, part[2:])
		}
	}
	return cf
}

func (cf *CFlags) GenHeaderFilePaths(files []string, defaultPaths []string) ([]string, []string, error) {
	var foundPaths []string
	var notFound []string

	searchPaths := append(cf.Paths, defaultPaths...)

	for _, file := range files {
		var found bool
		for _, path := range searchPaths {
			fullPath := filepath.Join(path, file)
			if _, err := os.Stat(fullPath); err == nil {
				foundPaths = append(foundPaths, fullPath)
				found = true
				break
			}
		}
		if !found {
			notFound = append(notFound, file)
		}
	}

	if len(foundPaths) == 0 {
		return nil, notFound, fmt.Errorf("failed to find any header files")
	}

	return foundPaths, notFound, nil
}
