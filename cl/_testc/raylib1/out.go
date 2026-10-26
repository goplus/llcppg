package foo

import "github.com/goplus/lib/c"

// Mouse buttons
type MouseButton c.Uint

const (
	MOUSE_BUTTON_LEFT    = 0
	MOUSE_BUTTON_RIGHT   = 1
	MOUSE_BUTTON_MIDDLE  = 2
	MOUSE_BUTTON_SIDE    = 3
	MOUSE_BUTTON_EXTRA   = 4
	MOUSE_BUTTON_FORWARD = 5
	MOUSE_BUTTON_BACK    = 6
)

// Color, 4 components, R8G8B8A8 (32bit)
type Color = uint32
