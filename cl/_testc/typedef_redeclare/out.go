package foo

import "github.com/goplus/lib/c"

type X_FrameHeader struct {
	Size c.Int
}
type FrameHeader = X_FrameHeader
type MyInt = c.Int
