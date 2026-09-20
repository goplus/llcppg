package foo

import "github.com/goplus/lib/c"

type Global c.Int

const (
	GA Global = 0
	GB Global = 1
)

type Bar_Color c.Int

const (
	Bar_Red   Bar_Color = 0
	Bar_Green Bar_Color = 5
	Bar_Blue  Bar_Color = 6
)

type Shape struct {
}
type Shape_Kind c.Int

const (
	Shape_Circle   Shape_Kind = 0
	Shape_Square   Shape_Kind = 1
	Shape_Triangle Shape_Kind = 2
)
