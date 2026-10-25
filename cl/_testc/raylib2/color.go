package foo

import (
	_ "unsafe"

	"github.com/goplus/lib/c"
)

type FloatPtr = *c.Float

// Color, 4 components, R8G8B8A8 (32bit)
// R, G, B, A uint8
type Color c.Uint

const (
	LIGHTGRAY = Color(200 | 200<<8 | 200<<16 | 255<<24) // Light Gray
	GRAY      = Color(130 | 130<<8 | 130<<16 | 255<<24) // Gray
	DARKGRAY  = Color(80 | 80<<8 | 80<<16 | 255<<24)    // Dark Gray
	YELLOW    = Color(253 | 249<<8 | 0<<16 | 255<<24)   // Yellow
	GOLD      = Color(255 | 203<<8 | 0<<16 | 255<<24)   // Gold
	ORANGE    = Color(255 | 161<<8 | 0<<16 | 255<<24)   // Orange
	PINK      = Color(255 | 109<<8 | 194<<16 | 255<<24) // Pink
	RED       = Color(230 | 41<<8 | 55<<16 | 255<<24)   // Red
	MAROON    = Color(190 | 33<<8 | 55<<16 | 255<<24)   // Maroon
	GREEN     = Color(0 | 228<<8 | 48<<16 | 255<<24)    // Green
	LIME      = Color(0 | 158<<8 | 47<<16 | 255<<24)    // Lime
	DARKGREEN = Color(0 | 117<<8 | 44<<16 | 255<<24)    // Dark Green
	SKYBLUE   = Color(102 | 191<<8 | 255<<16 | 255<<24) // Sky Blue
	BLUE      = Color(0 | 121<<8 | 241<<16 | 255<<24)   // Blue

	DARKBLUE   = Color(0 | 82<<8 | 172<<16 | 255<<24)    // Dark Blue
	PURPLE     = Color(200 | 122<<8 | 255<<16 | 255<<24) // Purple
	VIOLET     = Color(135 | 60<<8 | 190<<16 | 255<<24)  // Violet
	DARKPURPLE = Color(112 | 31<<8 | 126<<16 | 255<<24)  // Dark Purple
	BEIGE      = Color(211 | 176<<8 | 131<<16 | 255<<24) // Beige
	BROWN      = Color(127 | 106<<8 | 79<<16 | 255<<24)  // Brown
	DARKBROWN  = Color(76 | 63<<8 | 47<<16 | 255<<24)    // Dark Brown

	WHITE    = Color(255 | 255<<8 | 255<<16 | 255<<24) // White
	BLACK    = Color(0 | 0<<8 | 0<<16 | 255<<24)       // Black
	BLANK    = Color(0 | 0<<8 | 0<<16 | 0<<24)         // Blank (Transparent)
	MAGENTA  = Color(255 | 0<<8 | 255<<16 | 255<<24)   // Magenta
	RAYWHITE = Color(245 | 245<<8 | 245<<16 | 255<<24) // My own White (raylib logo)
)

type Foo struct {
	color Color
	FloatPtr
}

type Point3D struct {
	X [3](float64)
}

type FooT[T any] struct {
	Foo
	Value T
}

type BarT[T any] = FooT[T]
