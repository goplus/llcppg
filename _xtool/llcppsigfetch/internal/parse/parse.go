package parse

import (
	"fmt"
	"os"
	"strings"
	"unsafe"

	"github.com/goplus/lib/c"
	"github.com/goplus/llcppg/_xtool/internal/clangtool"
	"github.com/goplus/llcppg/_xtool/internal/config"
	"github.com/goplus/llcppg/_xtool/internal/parser"
	llcppg "github.com/goplus/llcppg/config"
	"github.com/goplus/llpkg/cjson"
)

type dbgFlags = int

var debugParse bool

const (
	DbgParse   dbgFlags = 1 << iota
	DbgFlagAll          = DbgParse
)

func SetDebug(dbgFlags dbgFlags) {
	debugParse = (dbgFlags & DbgParse) != 0
}

type Config struct {
	Conf   *llcppg.Config
	Out    bool     // if gen llcppg.sigfetch.json
	Cflags []string // other cflags want to parse

	ExtractMode bool
	ExtractFile string
	IsTemp      bool
	IsCpp       bool

	CombinedFile     string
	PreprocessedFile string
	Exec             func(conf *Config, pkg *llcppg.Pkg)
}

func Do(conf *Config) error {
	if debugParse {
		parser.SetDebug(parser.DbgFlagAll)
		fmt.Fprintln(os.Stderr, "output to file:", conf.Out)
		if conf.ExtractMode {
			fmt.Fprintln(os.Stderr, "runExtract: extractFile:", conf.ExtractFile)
			fmt.Fprintln(os.Stderr, "isTemp:", conf.IsTemp)
			fmt.Fprintln(os.Stderr, "isCpp:", conf.IsCpp)
			fmt.Fprintln(os.Stderr, "out:", conf.Out)
			fmt.Fprintln(os.Stderr, "otherArgs:", conf.Cflags)
		}
	}

	var isCpp bool
	if conf.IsCpp {
		isCpp = true
	} else {
		isCpp = conf.Conf.Cplusplus
	}
	if err := createTempIfNoExist(&conf.CombinedFile, conf.Conf.Name+"*.h"); err != nil {
		return err
	}
	if err := createTempIfNoExist(&conf.PreprocessedFile, conf.Conf.Name+"*.i"); err != nil {
		return err
	}

	if debugParse {
		fmt.Fprintln(os.Stderr, "Do: combinedFile", conf.CombinedFile)
		fmt.Fprintln(os.Stderr, "Do: preprocessedFile", conf.PreprocessedFile)
	}

	// compose includes to a combined file
	err := clangtool.ComposeIncludes(conf.Conf.Include, conf.CombinedFile)
	if err != nil {
		return err
	}

	// prepare clang flags to preprocess the combined file
	clangFlags := strings.Fields(conf.Conf.CFlags)
	clangFlags = append(clangFlags, "-C")  // keep comment
	clangFlags = append(clangFlags, "-dD") // keep macro
	clangFlags = append(clangFlags, "-fparse-all-comments")

	err = clangtool.Preprocess(&clangtool.PreprocessConfig{
		File:    conf.CombinedFile,
		IsCpp:   isCpp,
		Args:    clangFlags,
		OutFile: conf.PreprocessedFile,
	})
	if err != nil {
		return err
	}

	// https://github.com/goplus/llgo/issues/603
	// we need exec.Command("clang", "-print-resource-dir").Output() in llcppsigfetch to obtain the resource directory
	// to ensure consistency between clang preprocessing and libclang-extracted header filelink cflags.
	// Currently, directly calling exec.Command in the main flow of llcppsigfetch will cause hang and fail to execute correctly.
	// As a solution, the resource directory is externally provided by llcppg.
	libclangFlags := []string{"-fparse-all-comments"}

	pkgHfiles := config.PkgHfileInfo(conf.Conf, libclangFlags)
	if debugParse {
		fmt.Fprintln(os.Stderr, "interfaces", pkgHfiles.Inters)
		fmt.Fprintln(os.Stderr, "implements", pkgHfiles.Impls)
		fmt.Fprintln(os.Stderr, "thirdhfile", pkgHfiles.Thirds)
	}
	libclangFlags = append(libclangFlags, strings.Fields(conf.Conf.CFlags)...)
	file, err := parser.Do(&parser.ConverterConfig{
		File:  conf.PreprocessedFile,
		Args:  libclangFlags,
		IsCpp: isCpp,
	})
	if err != nil {
		return err
	}

	pkg := &llcppg.Pkg{
		File:    file,
		FileMap: make(map[string]*llcppg.FileInfo),
	}

	fileTypeMappings := []struct {
		files    []string
		fileType llcppg.FileType
	}{
		{pkgHfiles.Inters, llcppg.Inter},
		{pkgHfiles.Impls, llcppg.Impl},
		{pkgHfiles.Thirds, llcppg.Third},
	}

	for _, mapping := range fileTypeMappings {
		for _, file := range mapping.files {
			pkg.FileMap[file] = &llcppg.FileInfo{
				FileType: mapping.fileType,
			}
		}
	}

	if debugParse {
		fmt.Fprintln(os.Stderr, "Have %d Macros", len(pkg.File.Macros))
		for _, macro := range pkg.File.Macros {
			fmt.Fprintf(os.Stderr, "Macro %s", macro.Name)
		}
		fmt.Fprintln(os.Stderr)
	}

	if conf.Exec != nil {
		conf.Exec(conf, pkg)
	}

	return nil
}

func OutputPkg(conf *Config, pkg *llcppg.Pkg) {
	info := MarshalPkg(pkg)
	str := info.Print()
	defer cjson.FreeCStr(unsafe.Pointer(str))
	defer info.Delete()
	outputResult(str, conf.Out)
}

func outputResult(result *c.Char, outputToFile bool) {
	if outputToFile {
		outputFile := llcppg.LLCPPG_SIGFETCH
		err := os.WriteFile(outputFile, []byte(c.GoString(result)), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to output file: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Results saved to %s\n", outputFile)
	} else {
		c.Printf(c.Str("%s"), result)
	}
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
