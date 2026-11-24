package example

import (
	pic "image"
	"image/color"
)

type Cache[T any] struct {
	Last pic.Point
}

var Colors = map[color.Color]int{pic.Black.C: 1}

var Local = func(im struct{ X int }) int {
	return im.X
}(struct{ X int }{X: 3})
