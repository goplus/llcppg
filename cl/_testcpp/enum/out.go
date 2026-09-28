package foo

import "github.com/goplus/lib/c"

type Global c.Uint

const (
	GA Global = 0
	GB Global = 1
)

type BarColor c.Uint

const (
	BarRed   BarColor = 0
	BarGreen BarColor = 5
	BarBlue  BarColor = 6
)

type Shape struct {
}
type ShapeKind c.Uint

const (
	ShapeCircle   ShapeKind = 0
	ShapeSquare   ShapeKind = 1
	ShapeTriangle ShapeKind = 2
)
