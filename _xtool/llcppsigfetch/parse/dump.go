package parse

import (
	"github.com/goplus/lib/c"
	"github.com/goplus/llcppg/_xtool/internal/parser"
	llcppg "github.com/goplus/llcppg/config"
	"github.com/goplus/llpkg/cjson"
)

func MarshalPkg(pkg *llcppg.Pkg) *cjson.JSON {
	root := cjson.Object()
	root.SetItem(c.Str("File"), parser.MarshalASTFile(pkg.File))
	root.SetItem(c.Str("FileMap"), MarshalFileMap(pkg.FileMap))
	return root
}

func MarshalFileMap(fmap map[string]*llcppg.FileInfo) *cjson.JSON {
	root := cjson.Object()
	for path, info := range fmap {
		root.SetItem(c.AllocaCStr(path), MarshalFileInfo(info))
	}
	return root
}

func MarshalFileInfo(info *llcppg.FileInfo) *cjson.JSON {
	root := cjson.Object()
	root.SetItem(c.Str("FileType"), cjson.Number(float64(info.FileType)))
	return root
}
