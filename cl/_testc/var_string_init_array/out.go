package foo

import "github.com/goplus/lib/c"

// A field whose extent is inferred from the string-completed array keeps the
// resolved bound: sizeof(text) is 3, so `a` is a fixed [3]c.Char array.
type Owner struct {
	A [3]c.Char
}
