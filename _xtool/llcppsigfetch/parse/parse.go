package parse

import (
	"fmt"
	"os"
	"path"
	"strings"

	clangutils "github.com/goplus/llcppg/_xtool/llcppsymg/tool/clang"
	"github.com/goplus/llcppg/_xtool/llcppsymg/tool/config"
	llcppg "github.com/goplus/llcppg/config"
)

// temp to avoid call clang in llcppsigfetch,will cause hang
var ClangSearchPath []string
var ClangResourceDir string

type ParseConfig struct {
	Conf             *llcppg.Config
	CombinedFile     string
	PreprocessedFile string
	OutputFile       bool
}

func Do(cfg *ParseConfig) (*Converter, error) {
	if err := createTempIfNoExist(&cfg.CombinedFile, cfg.Conf.Name+"*.h"); err != nil {
		return nil, err
	}
	if err := createTempIfNoExist(&cfg.PreprocessedFile, cfg.Conf.Name+"*.i"); err != nil {
		return nil, err
	}

	if debugParse {
		fmt.Fprintln(os.Stderr, "Do: combinedFile", cfg.CombinedFile)
		fmt.Fprintln(os.Stderr, "Do: preprocessedFile", cfg.PreprocessedFile)
	}

	// compose includes to a combined file
	err := clangutils.ComposeIncludes(cfg.Conf.Include, cfg.CombinedFile)
	if err != nil {
		return nil, err
	}

	// prepare clang flags to preprocess the combined file
	clangFlags := strings.Fields(cfg.Conf.CFlags)
	clangFlags = append(clangFlags, "-C")  // keep comment
	clangFlags = append(clangFlags, "-dD") // keep macro
	clangFlags = append(clangFlags, "-fparse-all-comments")

	err = clangutils.Preprocess(&clangutils.PreprocessConfig{
		File:    cfg.CombinedFile,
		IsCpp:   cfg.Conf.Cplusplus,
		Args:    clangFlags,
		OutFile: cfg.PreprocessedFile,
	})
	if err != nil {
		return nil, err
	}

	// https://github.com/goplus/llgo/issues/603
	// we need exec.Command("clang", "-print-resource-dir").Output() in llcppsigfetch to obtain the resource directory
	// to ensure consistency between clang preprocessing and libclang-extracted header filelink cflags.
	// Currently, directly calling exec.Command in the main flow of llcppsigfetch will cause hang and fail to execute correctly.
	// As a solution, the resource directory is externally provided by llcppg.
	libclangFlags := []string{"-fparse-all-comments"}
	if ClangResourceDir != "" {
		libclangFlags = append(libclangFlags, "-resource-dir="+ClangResourceDir, "-I"+path.Join(ClangResourceDir, "include"))
	}
	pkgHfiles := config.PkgHfileInfo(cfg.Conf, libclangFlags)
	if debugParse {
		fmt.Fprintln(os.Stderr, "interfaces", pkgHfiles.Inters)
		fmt.Fprintln(os.Stderr, "implements", pkgHfiles.Impls)
		fmt.Fprintln(os.Stderr, "thirdhfile", pkgHfiles.Thirds)
	}
	libclangFlags = append(libclangFlags, strings.Fields(cfg.Conf.CFlags)...)
	converter, err := NewConverter(
		&Config{
			HfileInfo: pkgHfiles,
			Cfg: &clangutils.Config{
				File:  cfg.PreprocessedFile,
				IsCpp: cfg.Conf.Cplusplus,
				Args:  libclangFlags,
			},
		})
	if err != nil {
		return nil, err
	}
	pkg, err := converter.Convert()
	if err != nil {
		return nil, err
	}
	if debugParse {
		fmt.Fprintln(os.Stderr, "Have %d Macros", len(pkg.File.Macros))
		for _, macro := range pkg.File.Macros {
			fmt.Fprintf(os.Stderr, "Macro %s", macro.Name)
		}
		fmt.Fprintln(os.Stderr)
	}
	return converter, nil
}

func createTempIfNoExist(filename *string, pattern string) error {
	if *filename != "" {
		return nil
	}
	f, err := os.CreateTemp("", pattern)
	if err != nil {
		return err
	}
	*filename = f.Name()
	return nil
}
