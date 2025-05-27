package config

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/goplus/lib/c/clang"
	clangutils "github.com/goplus/llcppg/_xtool/internal/clang"
	"github.com/goplus/llcppg/_xtool/internal/clangtool"
	llcppg "github.com/goplus/llcppg/config"
)

type Conf struct {
	*llcppg.Config
}

func GetConf(useStdin bool, cfgFile string) (conf Conf, err error) {
	var data []byte
	if useStdin {
		data, err = io.ReadAll(os.Stdin)
	} else {
		data, err = os.ReadFile(cfgFile)
	}
	if err != nil {
		return
	}
	conf, err = GetConfByByte(data)
	if err != nil {
		return
	}
	return conf, nil
}

func GetConfByByte(data []byte) (Conf, error) {
	var config llcppg.Config

	if err := json.Unmarshal(data, &config); err != nil {
		return Conf{}, err
	}

	return Conf{
		Config: &config,
	}, nil
}

type PkgHfilesInfo struct {
	Inters []string // From types.Config.Include
	Impls  []string // From same root of types.Config.Include
	Thirds []string // Not Current Pkg's Files
}

func (p *PkgHfilesInfo) CurPkgFiles() []string {
	return append(p.Inters, p.Impls...)
}

// PkgHfileInfo analyzes header files dependencies and categorizes them into three groups:
// 1. Inters: Direct includes from types.Config.Include
// 2. Impls: Header files from the same root directory as Inters
// 3. Thirds: Header files from external sources
//
// The function works by:
// 1. Creating a temporary header file that includes all headers from conf.Include
// 2. Using clang to parse the translation unit and analyze includes
// 3. Categorizing includes based on their inclusion level and path relationship
func PkgHfileInfo(conf *llcppg.Config, args []string) *PkgHfilesInfo {
	info := &PkgHfilesInfo{
		Inters: []string{},
		Impls:  []string{},
		Thirds: []string{},
	}
	outfile, err := os.CreateTemp("", "compose_*.h")
	if err != nil {
		panic(err)
	}

	cflags := append(args, strings.Fields(conf.CFlags)...)
	inters := make(map[string]struct{})
	others := []string{} // impl & third
	for _, f := range conf.Include {
		content := "#include <" + f + ">"
		index, unit, err := clangutils.CreateTranslationUnit(&clangutils.Config{
			File: content,
			Temp: true,
			Args: cflags,
		})
		if err != nil {
			panic(err)
		}
		clangutils.GetInclusions(unit, func(inced clang.File, incins []clang.SourceLocation) {
			if len(incins) == 1 {
				filename := filepath.Clean(clang.GoString(inced.FileName()))
				info.Inters = append(info.Inters, filename)
				inters[filename] = struct{}{}
			}
		})
		unit.Dispose()
		index.Dispose()
	}

	clangtool.ComposeIncludes(conf.Include, outfile.Name())
	index, unit, err := clangutils.CreateTranslationUnit(&clangutils.Config{
		File: outfile.Name(),
		Temp: false,
		Args: cflags,
	})
	defer unit.Dispose()
	defer index.Dispose()
	if err != nil {
		panic(err)
	}
	clangutils.GetInclusions(unit, func(inced clang.File, incins []clang.SourceLocation) {
		// not in the first level include maybe impl or third hfile
		filename := filepath.Clean(clang.GoString(inced.FileName()))
		_, inter := inters[filename]
		if len(incins) > 1 && !inter {
			others = append(others, filename)
		}
	})

	if conf.Mix {
		info.Thirds = others
		return info
	}

	root, err := filepath.Abs(commonParentDir(info.Inters))
	if err != nil {
		panic(err)
	}
	for _, f := range others {
		file, err := filepath.Abs(f)
		if err != nil {
			panic(err)
		}
		if strings.HasPrefix(file, root) {
			info.Impls = append(info.Impls, f)
		} else {
			info.Thirds = append(info.Thirds, f)
		}
	}
	return info
}

// commonParentDir finds the longest common parent directory path for a given slice of paths.
// For example, given paths ["/a/b/c/d", "/a/b/e/f"], it returns "/a/b".
func commonParentDir(paths []string) string {
	if len(paths) == 0 {
		return ""
	}

	parts := make([][]string, len(paths))
	for i, path := range paths {
		parts[i] = strings.Split(filepath.Dir(path), string(filepath.Separator))
	}

	for i := 0; i < len(parts[0]); i++ {
		for j := 1; j < len(parts); j++ {
			if i == len(parts[j]) || parts[j][i] != parts[0][i] {
				return filepath.Join(parts[0][:i]...)
			}
		}
	}
	return filepath.Dir(paths[0])
}
