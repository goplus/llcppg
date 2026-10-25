package foo

import "github.com/goplus/lib/c"

// Get color with alpha applied, alpha goes from 0.0f to 1.0f
//
// llgo:link Color.Alpha C.ColorAlpha
func (color Color) Alpha(alpha c.Float) Color {
	return color
}
