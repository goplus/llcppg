package gencfg

import "flag"

var dependencies string
var extsString string
var excludes string
var cpp, help, tab bool

func addFlags(fs *flag.FlagSet) {
	fs.BoolVar(&cpp, "cpp", false, "if it is C++ lib")
	fs.BoolVar(&help, "help", false, "print help message")
	fs.BoolVar(&tab, "tab", true, "generate .cfg config file with tab indent")
	fs.StringVar(&excludes, "excludes", "", "exclude all header files in subdir of include example -excludes=\"internal impl\"")
	fs.StringVar(&extsString, "exts", ".h", "extra include file extensions for example -exts=\".h .hpp .hh\"")
	fs.StringVar(&dependencies, "deps", "", "deps for autofilling dependencies")
}
