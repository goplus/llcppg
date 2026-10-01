package foo

import "github.com/goplus/lib/c"

type Global c.Int

const (
	GA Global = 0
	GB Global = 1
)

type Local c.Int

const (
	Local_LA Local = 0
	Local_LB Local = 1
)

type BarColor c.Uint

const (
	Bar_Red   BarColor = 0
	Bar_Green BarColor = 5
	Bar_Blue  BarColor = 6
)

type Shape struct {
}
type ShapeKind c.Uint

const (
	Shape_Circle   ShapeKind = 0
	Shape_Square   ShapeKind = 1
	Shape_Triangle ShapeKind = 2
)
