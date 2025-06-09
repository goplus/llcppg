package sqlite3

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

// llgo:type C
type LoadextEntry func(*Sqlite3, **c.Char, *ApiRoutines) c.Int
