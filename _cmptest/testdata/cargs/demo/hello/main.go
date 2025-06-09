package main

import (
	"fmt"
	"os"

	"cargs"

	C "github.com/goplus/lib/c"
)

func main() {
	// Define command-line options
	options := []cargs.Option{
		{
			Identifier:    'h',
			AccessLetters: C.Str("h"),
			AccessName:    C.Str("help"),
			ValueName:     nil,
			Description:   C.Str("Show help information"),
		},
		{
			Identifier:    'v',
			AccessLetters: C.Str("v"),
			AccessName:    C.Str("version"),
			ValueName:     nil,
			Description:   C.Str("Show version information"),
		},
	}

	args := os.Args

	// Convert Go string array to C-style argv
	argv := make([]*int8, len(args))
	for i, arg := range args {
		argv[i] = C.AllocaCStr(arg)
	}

	// Initialize option context
	var context cargs.OptionContext
	context.OptionInit(&options[0], uintptr(len(options)), C.Int(len(args)), &argv[0])

	// Process all options
	identifierFound := false
	for context.OptionFetch() {
		identifierFound = true
		identifier := context.OptionGetIdentifier()
		switch identifier {
		case 'h':
			fmt.Println("Help: This is a simple command-line parser demo")
		case 'v':
			fmt.Println("Version: 1.0.0")
		}
	}

	// Default output if no identifier is found
	if !identifierFound {
		fmt.Println("Demo Command-line Tool\nIdentifier:\n\t-h: Help\n\t-v: Version")
	}
}
