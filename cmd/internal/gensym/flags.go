package gensym

import "flag"

var verbose bool

func addFlags(fs *flag.FlagSet) {
	fs.BoolVar(&verbose, "v", false, "enable verbose output")
}
