package convert

import (
	"github.com/goplus/llcppg/_xtool/llcppsymg/tool/name"
	llcppg "github.com/goplus/llcppg/config"
)

type HeaderFile struct {
	File     string
	FileType llcppg.FileType
}

// Note:third hfile should not set to gogen.Package
func (p *HeaderFile) ToGoFileName(pkgName string) string {
	switch p.FileType {
	case llcppg.Inter:
		return name.HeaderFileToGo(p.File)
	case llcppg.Impl, llcppg.Third:
		return pkgName + "_autogen.go"
	default:
		panic("unkown FileType")
	}
}

func (p *HeaderFile) InCurPkg() bool {
	return p.FileType == llcppg.Inter || p.FileType == llcppg.Impl
}

func NewHeaderFile(file string, fileType llcppg.FileType) *HeaderFile {
	return &HeaderFile{
		File:     file,
		FileType: fileType,
	}
}
