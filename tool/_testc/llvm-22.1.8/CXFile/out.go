package CXFile

import (
	"clang/CXString"
	"clang/cstdlib"
	"github.com/goplus/lib/c"
)

const LLGoPackage = "link: -L$(llvm-config --libdir) -lclang; -lclang"

type File uintptr
type FileUniqueID struct {
	Data [3]c.UlongLong
}

// llgo:link File.GetFileName C.clang_getFileName
func (SFile File) GetFileName() CXString.String {
	return CXString.String{}
}

// llgo:link File.GetFileTime C.clang_getFileTime
func (SFile File) GetFileTime() cstdlib.TimeT {
	return 0
}

// llgo:link File.GetFileUniqueID C.clang_getFileUniqueID
func (file File) GetFileUniqueID(outID *FileUniqueID) c.Int {
	return 0
}

// llgo:link File.FileIsEqual C.clang_File_isEqual
func (file1 File) FileIsEqual(file2 File) c.Int {
	return 0
}

// llgo:link File.FileTryGetRealPathName C.clang_File_tryGetRealPathName
func (file File) FileTryGetRealPathName() CXString.String {
	return CXString.String{}
}
