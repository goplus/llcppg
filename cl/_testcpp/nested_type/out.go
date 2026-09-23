package foo

import "github.com/goplus/lib/c"

type Outer struct {
	Inner OuterInner
}
type OuterValueType = c.Int
type OuterInner struct {
	X c.Int
}
type OuterDetail struct {
	secret c.Int
}
