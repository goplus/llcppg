package foo

import "github.com/goplus/lib/c"

type Outer struct {
	Inner Outer_Inner
}
type Outer_value_type = c.Int
type Outer_Inner struct {
	X c.Int
}
type Outer_Detail struct {
	secret c.Int
}
