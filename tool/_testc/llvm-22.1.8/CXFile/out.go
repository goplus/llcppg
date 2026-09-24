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

// llgo:link File.Name C.clang_getFileName
func (SFile File) Name() CXString.String {
	return CXString.String{}
}

// llgo:link File.Time C.clang_getFileTime
func (SFile File) Time() cstdlib.TimeT {
	return 0
}

// llgo:link File.UniqueID C.clang_getFileUniqueID
func (file File) UniqueID(outID *FileUniqueID) c.Int {
	return 0
}

// llgo:link File.IsEqual C.clang_File_isEqual
func (file1 File) IsEqual(file2 File) c.Int {
	return 0
}

// llgo:link File.TryGetRealPathName C.clang_File_tryGetRealPathName
func (file File) TryGetRealPathName() CXString.String {
	return CXString.String{}
}
