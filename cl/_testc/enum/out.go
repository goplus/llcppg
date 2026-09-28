package foo

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type Color c.Uint

const (
	Red   Color = 0
	Green Color = 1
	Blue  Color = 2
)

type State c.Uint

const (
	StateStopped State = 0
	StateRunning State = 10
	StatePaused  State = 11
)
const (
	FlagA = 1
	FlagB = 2
	FlagC = 4
)

type CXCompilationDatabase_Error c.Uint

const (
	CXCompilationDatabase_NoError            CXCompilationDatabase_Error = 0
	CXCompilationDatabase_CanNotLoadDatabase CXCompilationDatabase_Error = 1
)

//go:linkname F C.f
func F(_llcppg_param1 c.Int, _llcppg_param2 CXCompilationDatabase_Error)
