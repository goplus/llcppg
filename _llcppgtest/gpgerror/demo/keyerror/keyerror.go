package main

import (
	"gpgerror"

	"github.com/goplus/lib/c"
)

var SOURCEMASK = gpgerror.SOURCEDIM - 1
var CODEMASK = gpgerror.CODEDIM - 1

func ErrMake(source gpgerror.SourceT, code gpgerror.CodeT) gpgerror.ErrorT {
	if code == gpgerror.NOERROR {
		return gpgerror.ErrorT(gpgerror.NOERROR)
	}

	return gpgerror.ErrorT(((c.Int(source) & c.Int(SOURCEMASK)) << gpgerror.SOURCE_SHIFT) | (c.Int(code) & c.Int(CODEMASK)))
}

func main() {
	gpgerror.Init()
	err := ErrMake(gpgerror.SOURCEUSER1, gpgerror.BADKEY)
	errStr := err.Strerror()
	source := err.Strsource()
	c.Printf(c.Str("Error: %s\n"), errStr)
	c.Printf(c.Str("Source: %s\n"), source)
}
