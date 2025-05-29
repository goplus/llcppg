package symbol

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// FindLibs finds the library file in the given path & the given name.
type Mode int

const (
	ModeDynamic Mode = iota
	ModeStatic
)

func FindLibFile(path string, name string, mode Mode) (string, error) {
	affix := libAffix(mode)
	libPath := filepath.Join(path, fmt.Sprintf("lib%s%s", name, affix))
	_, err := os.Stat(libPath)
	if err != nil {
		return "", err
	}
	return libPath, nil
}

func libAffix(mode Mode) (affix string) {
	if mode == ModeStatic {
		return ".a"
	}
	if runtime.GOOS == "linux" {
		return ".so"
	}
	return ".dylib"
}
