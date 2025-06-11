package gen

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	llcppg "github.com/goplus/llcppg/config"
)

type Config struct {
	Name           string
	IsCpp          bool
	Exts           []string
	Deps           []string
	ExcludeSubdirs []string
}

type llcppCfgKey string

const (
	cfgLibsKey   llcppCfgKey = "libs"
	cfgCflagsKey llcppCfgKey = "cflags"
)

type emptyStringError struct {
	name string
}

func (p *emptyStringError) Error() string {
	return p.name + " can't be empty"
}

func newEmptyStringError(name string) *emptyStringError {
	return &emptyStringError{name: name}
}

func getDir(relPath string) string {
	index := strings.IndexRune(relPath, filepath.Separator)
	if index < 0 {
		return relPath
	}
	return relPath[:index]
}

func isExcludeDir(relPath string, excludeSubdirs []string) bool {
	if len(excludeSubdirs) == 0 {
		return false
	}
	dir := getDir(relPath)
	for _, subdir := range excludeSubdirs {
		if subdir == dir {
			return true
		}
	}
	return false
}

func ExpandName(name string, dir string, cfgKey llcppCfgKey) string {
	originString := fmt.Sprintf("$(pkg-config --%s %s)", cfgKey, name)
	return ExpandString(originString, dir)
}

func findDepSlice(lines []string) ([]string, string) {
	objFileString := ""
	iStart := 0
	numLines := len(lines)
	complete := false
	for i := 0; i < numLines && !complete; i++ {
		line := lines[i]
		if strings.ContainsRune(line, rune(':')) && !strings.HasSuffix(line, ":") {
			objFileString = line
			iStart = i + 1
			break
		}
		complete = true
		for j := i + 1; j < numLines; j++ {
			line2 := lines[j]
			if len(line2) > 0 {
				iStart = j + 1
				objFileString = line + line2
				break
			}
		}
	}
	if iStart < numLines {
		return lines[iStart:], objFileString
	}
	return []string{}, objFileString
}

func getClangArgs(cflags string, relpath string) []string {
	args := make([]string, 0)
	cflagsField := strings.Fields(cflags)
	args = append(args, cflagsField...)
	args = append(args, "-MM")
	args = append(args, relpath)
	return args
}

func parseFileEntry(cflags, trimCflag, path string, d fs.DirEntry, exts []string, excludeSubdirs []string) (*ObjFile, error) {
	if d.IsDir() || strings.HasPrefix(d.Name(), ".") {
		return nil, errors.New("invalid file entry")
	}
	idx := len(exts)
	for i, ext := range exts {
		if strings.HasSuffix(d.Name(), ext) {
			idx = i
			break
		}
	}
	if idx == len(exts) {
		return nil, errors.New("invalid file ext")
	}
	relPath, err := filepath.Rel(trimCflag, path)
	if err != nil {
		relPath = path
	}
	if isExcludeDir(relPath, excludeSubdirs) {
		return nil, errors.New("file in excluded directory")
	}
	args := getClangArgs(cflags, relPath)
	clangCmd := NewExecCommand("clang", args...)
	outString, err := GetOut(clangCmd, trimCflag)
	if err != nil {
		log.Println(outString)
		return NewObjFile(relPath, relPath), errors.New(outString)
	}
	outString = strings.ReplaceAll(outString, "\\\n", "\n")
	fields := strings.Fields(outString)
	lines, objFileStr := findDepSlice(fields)
	objFile := NewObjFileString(objFileStr)
	objFile.Deps = append(objFile.Deps, lines...)
	return objFile, nil
}

func parseCFlagsEntry(cflags, cflag string, exts []string, excludeSubdirs []string) *CflagEntry {
	if !strings.HasPrefix(cflag, "-I") {
		return nil
	}
	trimCflag := strings.TrimPrefix(cflag, "-I")
	if !strings.HasSuffix(trimCflag, string(filepath.Separator)) {
		trimCflag += string(filepath.Separator)
	}
	var cflagEntry CflagEntry
	cflagEntry.Include = trimCflag
	cflagEntry.ObjFiles = make([]*ObjFile, 0)
	cflagEntry.InvalidObjFiles = make([]*ObjFile, 0)
	err := filepath.WalkDir(trimCflag, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		pObjFile, err := parseFileEntry(cflags, trimCflag, path, d, exts, excludeSubdirs)
		if err != nil {
			if pObjFile != nil {
				cflagEntry.InvalidObjFiles = append(cflagEntry.InvalidObjFiles, pObjFile)
			}
			return nil
		}
		if pObjFile != nil {
			cflagEntry.ObjFiles = append(cflagEntry.ObjFiles, pObjFile)
		}
		return nil
	})
	sort.Slice(cflagEntry.ObjFiles, func(i, j int) bool {
		return len(cflagEntry.ObjFiles[i].Deps) > len(cflagEntry.ObjFiles[j].Deps)
	})
	if err != nil {
		return nil
	}
	return &cflagEntry
}

func NormalizePackageName(name string) string {
	fields := strings.FieldsFunc(name, func(r rune) bool {
		return !unicode.IsLetter(r) && r != '_' && !unicode.IsDigit(r)
	})
	if len(fields) > 0 {
		if len(fields[0]) > 0 && unicode.IsDigit(rune(fields[0][0])) {
			fields[0] = "_" + fields[0]
		}
	}
	return strings.Join(fields, "_")
}

func (c *Config) toLLCppg() *llcppg.Config {
	cfg := llcppg.NewDefault()
	cfg.Name = NormalizePackageName(c.Name)
	cfg.CFlags = fmt.Sprintf("$(pkg-config --cflags %s)", c.Name)
	cfg.Libs = fmt.Sprintf("$(pkg-config --libs %s)", c.Name)
	cfg.Cplusplus = c.IsCpp
	cfg.Deps = c.Deps

	expandCFlags := ExpandName(c.Name, "", cfgCflagsKey)
	cfg.Include = c.sortIncludes(expandCFlags, c.Exts, c.ExcludeSubdirs)

	return cfg
}

func (c *Config) sortIncludes(expandCflags string, exts []string, excludeSubdirs []string) []string {
	list := strings.Fields(expandCflags)
	includeList := NewIncludeList()
	for i, cflag := range list {
		pCflagEntry := parseCFlagsEntry(expandCflags, cflag, exts, excludeSubdirs)
		includeList.AddCflagEntry(i, pCflagEntry)
	}
	return includeList.include
}

func Do(genCfg *Config) ([]byte, error) {
	if len(genCfg.Name) == 0 {
		return nil, newEmptyStringError("name")
	}
	cfg := genCfg.toLLCppg()

	return json.MarshalIndent(cfg, "", "  ")
}
