package foo

import "github.com/goplus/lib/c"

type Color c.Int

const (
	Red   Color = 0
	Green Color = 1
	Blue  Color = 2
)

type State c.Int

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
