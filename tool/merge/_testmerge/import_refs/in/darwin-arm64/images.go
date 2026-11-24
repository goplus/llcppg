package example

import (
	im "image"
	"image/color"
)

type Cache[T any] struct {
	Last im.Point
}

var Colors = map[color.Color]int{im.Black.C: 1}

var Local = func(im struct{ X int }) int {
	return im.X
}(struct{ X int }{X: 3})

var Origin = im.Point{}
