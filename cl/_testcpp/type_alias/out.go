package foo

import "github.com/goplus/lib/c"

const LLGoPackage = "link: -L/path/foo -lfoo"

type StdBasicString[CharT any, Traits any, Allocator any] struct {
}
type StdPtrdiffT = c.PtrdiffT
type StdString = StdBasicString[byte, c.Void, c.Void]
