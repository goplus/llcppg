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

type PyStatus struct {
	X_type   _llcppg_anon_0
	Func     *c.Char
	ErrMsg   *c.Char
	Exitcode c.Int
}
type _llcppg_anon_0 c.Uint

const (
	X_PyStatus_TYPE_OK    _llcppg_anon_0 = 0
	X_PyStatus_TYPE_ERROR _llcppg_anon_0 = 1
	X_PyStatus_TYPE_EXIT  _llcppg_anon_0 = 2
)

//go:linkname F C.f
func F(_llcppg_param1 c.Int, _llcppg_param2 CXCompilationDatabase_Error)
