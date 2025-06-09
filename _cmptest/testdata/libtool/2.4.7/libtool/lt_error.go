package libtool

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

const ERROR_H = 1
const (
	ERROR_UNKNOWN               c.Int = 0
	ERROR_DLOPEN_NOT_SUPPORTED  c.Int = 1
	ERROR_INVALID_LOADER        c.Int = 2
	ERROR_INIT_LOADER           c.Int = 3
	ERROR_REMOVE_LOADER         c.Int = 4
	ERROR_FILE_NOT_FOUND        c.Int = 5
	ERROR_DEPLIB_NOT_FOUND      c.Int = 6
	ERROR_NO_SYMBOLS            c.Int = 7
	ERROR_CANNOT_OPEN           c.Int = 8
	ERROR_CANNOT_CLOSE          c.Int = 9
	ERROR_SYMBOL_NOT_FOUND      c.Int = 10
	ERROR_NO_MEMORY             c.Int = 11
	ERROR_INVALID_HANDLE        c.Int = 12
	ERROR_BUFFER_OVERFLOW       c.Int = 13
	ERROR_INVALID_ERRORCODE     c.Int = 14
	ERROR_SHUTDOWN              c.Int = 15
	ERROR_CLOSE_RESIDENT_MODULE c.Int = 16
	ERROR_INVALID_MUTEX_ARGS    c.Int = 17
	ERROR_INVALID_POSITION      c.Int = 18
	ERROR_CONFLICTING_FLAGS     c.Int = 19
	ERROR_MAX                   c.Int = 20
)

/* These functions are only useful from inside custom module loaders. */
//go:linkname Dladderror C.lt_dladderror
func Dladderror(diagnostic *c.Char) c.Int

//go:linkname Dlseterror C.lt_dlseterror
func Dlseterror(errorcode c.Int) c.Int
