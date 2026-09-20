package foo

import "github.com/goplus/lib/c"

type Base struct {
	A c.Int
	B c.Uint
}
type Derived struct {
	Base
	C c.Int
}
type Mid struct {
	M c.Int
}
type Multi struct {
	Base
	Mid
	N c.Int
}
