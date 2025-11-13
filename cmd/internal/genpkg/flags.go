package genpkg

import (
	"flag"
)

var modulePath string
var verbose bool

func addFlags(fs *flag.FlagSet) {
	fs.BoolVar(&verbose, "v", false, "enable verbose output")
	fs.StringVar(&modulePath, "mod", "", "the module path of the generated code,if not set,will not init a new module")
}
