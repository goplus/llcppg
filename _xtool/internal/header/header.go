package header

import (
	"fmt"
	"io"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/goplus/lib/c/clang"
	clangutils "github.com/goplus/llcppg/_xtool/internal/clang"
	"github.com/goplus/llcppg/_xtool/internal/clangtool"
)

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
func PkgHfileInfo(includes []string, args []string, mix bool) *PkgHfilesInfo {
	info := &PkgHfilesInfo{
		Inters: []string{},
		Impls:  []string{},
		Thirds: []string{},
	}
	outfile, err := os.CreateTemp("", "compose_*.h")
	if err != nil {
		panic(err)
	}
	defer os.Remove(outfile.Name())

	mmOutput, err := os.CreateTemp("", "mmoutput_*")
	if err != nil {
		panic(err)
	}
	defer os.Remove(mmOutput.Name())

	includeTrie := NewTrie(WithReversePathSegmenter())

	for _, inc := range includes {
		includeTrie.Insert(inc)
		args = append(args, fmt.Sprintf("--no-system-header-prefix=%s", inc))
	}

	args = append(args, "-MM", "-MF", mmOutput.Name())

	clangtool.ComposeIncludes(includes, outfile.Name())
	index, unit, err := clangutils.CreateTranslationUnit(&clangutils.Config{
		File: outfile.Name(),
		Temp: false,
		Args: args,
	})

	defer unit.Dispose()
	defer index.Dispose()
	if err != nil {
		panic(err)
	}

	var others []string
	inters, longestPrefix := RetrieveInterfaceFromMM(outfile.Name(), mmOutput, includeTrie)

	clangutils.GetInclusions(unit, func(inced clang.File, incins []clang.SourceLocation) {
		// not in the first level include maybe impl or third hfile
		filename := filepath.Clean(clang.GoString(inced.FileName()))

		// skip the composed header
		if filename == outfile.Name() {
			return
		}

		if _, isInterface := inters[filename]; !isInterface {
			others = append(others, filename)
		}
	})

	info.Inters = slices.Collect(maps.Keys(inters))

	absLongestPrefix, err := filepath.Abs(longestPrefix)
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(os.Stderr, "tttttt", longestPrefix, CommonParentDir(info.Inters))

	for _, filename := range others {
		if mix {
			info.Thirds = append(info.Thirds, filename)
			continue
		}
		filePath, err := filepath.Abs(filename)
		if err != nil {
			panic(err)
		}
		if strings.HasPrefix(filePath, absLongestPrefix) {
			info.Impls = append(info.Impls, filename)
		} else {
			info.Thirds = append(info.Thirds, filename)
		}
	}

	sort.Strings(info.Inters)

	return info
}

// commonParentDir finds the longest common parent directory path for a given slice of paths.
// For example, given paths ["/a/b/c/d", "/a/b/e/f"], it returns "/a/b".
func CommonParentDir(paths []string) string {
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

func RetrieveInterfaceFromMM(
	composedHeaderFileName string,
	mmOutput *os.File,
	includeTrie *Trie,
) (interfaceMap map[string]struct{}, prefix string) {
	fileName := strings.TrimSuffix(filepath.Base(composedHeaderFileName), ".h")

	interfaceMap = make(map[string]struct{})

	content, _ := io.ReadAll(mmOutput)

	mmTrie := NewTrie()

	for _, line := range strings.Fields(string(content)) {
		// skip composed header file
		if strings.Contains(line, fileName) || line == `\` {
			continue
		}
		headerFile := filepath.Clean(line)

		if includeTrie.IsOnSameBranch(headerFile) {
			mmTrie.Insert(headerFile)

			interfaceMap[headerFile] = struct{}{}
		}
	}
	prefix = mmTrie.LongestPrefix()
	return
}
