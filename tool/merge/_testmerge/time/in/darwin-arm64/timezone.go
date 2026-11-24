package libc

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

//go:linkname Timezone C.timezone
var Timezone c.Long // Seconds west of UTC.
