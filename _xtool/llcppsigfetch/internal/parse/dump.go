package parse

import (
	"github.com/goplus/llcppg/_xtool/internal/parser"
	llcppg "github.com/goplus/llcppg/config"
)

func MarshalPkg(pkg *llcppg.Pkg) map[string]any {
	return map[string]any{
		"File":    parser.XMarshalASTFile(pkg.File),
		"FileMap": MarshalFileMap(pkg.FileMap),
	}
}

func MarshalFileMap(fmap map[string]*llcppg.FileInfo) map[string]any {
	root := make(map[string]any)
	for path, info := range fmap {
		root[path] = MarshalFileInfo(info)
	}
	return root
}

func MarshalFileInfo(info *llcppg.FileInfo) map[string]any {
	return map[string]any{
		"FileType": float64(info.FileType),
	}
}
