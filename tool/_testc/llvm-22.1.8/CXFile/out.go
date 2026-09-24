package CXFile

import (
	"clang/CXString"
	"clang/cstdlib"
	"github.com/goplus/lib/c"
	"unsafe"
)

const LLGoPackage = "link: -L$(llvm-config --libdir) -lclang; -lclang"

type File = unsafe.Pointer
type FileUniqueID struct {
	Data [3]c.UlongLong
}

//go:linkname GetFileName C.clang_getFileName
func GetFileName(SFile File) CXString.String

//go:linkname GetFileTime C.clang_getFileTime
func GetFileTime(SFile File) cstdlib.TimeT

//go:linkname GetFileUniqueID C.clang_getFileUniqueID
func GetFileUniqueID(file File, outID *FileUniqueID) c.Int

//go:linkname FileIsEqual C.clang_File_isEqual
func FileIsEqual(file1 File, file2 File) c.Int

//go:linkname FileTryGetRealPathName C.clang_File_tryGetRealPathName
func FileTryGetRealPathName(file File) CXString.String
