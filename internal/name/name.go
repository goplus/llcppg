package name

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode"
)

func GoName(name string, trimPrefixes []string, inCurPkg bool) string {
	if inCurPkg {
		name = RemovePrefixedName(name, trimPrefixes)
	}
	return PubName(name)
}

func RemovePrefixedName(name string, trimPrefixes []string) string {
	for _, prefix := range trimPrefixes {
		if strings.HasPrefix(name, prefix) {
			return strings.TrimPrefix(name, prefix)
		}
	}
	return name
}

func PubName(name string) string {
	if len(name) == 0 {
		return name
	}
	baseName := strings.Trim(name, "_")
	if len(baseName) == 0 {
		return "X" + name
	}
	prefix := preUScore(name)
	suffix := sufUScore(name)

	if len(prefix) != 0 || unicode.IsDigit(rune(baseName[0])) {
		return "X" + prefix + ToCamelCase(baseName, false) + suffix
	}
	return ToCamelCase(baseName, true) + suffix
}

func sufUScore(name string) string {
	return strings.Repeat("_", len(name)-len(strings.TrimRight(name, "_")))
}

func preUScore(name string) string {
	return strings.Repeat("_", len(name)-len(strings.TrimLeft(name, "_")))
}

func ToCamelCase(s string, firstPartUpper bool) string {
	parts := strings.Split(s, "_")
	result := []string{}
	for i, part := range parts {
		if i == 0 && !firstPartUpper {
			result = append(result, part)
			continue
		}
		if len(part) > 0 {
			result = append(result, strings.ToUpper(part[:1])+part[1:])
		}
	}
	return strings.Join(result, "")
}

// Only Make it Public, no turn to other camel method
func ExportName(name string) string {
	fChar := name[0]
	if fChar == '_' || unicode.IsDigit(rune(fChar)) {
		return "X" + name
	}
	return UpperFirst(name)
}

func UpperFirst(name string) string {
	return strings.ToUpper(name[:1]) + name[1:]
}

// /path/to/foo.h -> foo.go
// /path/to/_intptr.h -> X_intptr.go
func HeaderFileToGo(incPath string) string {
	_, fileName := filepath.Split(incPath)
	ext := filepath.Ext(fileName)
	if len(ext) > 0 {
		fileName = strings.TrimSuffix(fileName, ext)
	}
	if strings.HasPrefix(fileName, "_") {
		fileName = "X" + fileName
	}
	return fileName + ".go"
}

func SuffixCount(name string, count int) string {
	if count > 1 {
		return fmt.Sprintf("%s__%d", name, count-1)
	}
	return name
}
