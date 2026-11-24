package example

type Header (struct {
	Ready    bool // ready to read
	Platform int  // Linux field.
})

// Version identifies the supported ABI.
const Version = 1

// See [Flags] for Linux values.
//
// [Flags]: https://example.org/flags
const Flags = 3

/* Platform limit for Linux. */
const Limit = 8

//go:generate echo "a  \"b\""
var Generate = 1

//go:generate echo "a b"
var Options = 2
