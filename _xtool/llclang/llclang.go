package main

import (
	"fmt"
	"os"

	"github.com/goplus/llcppg/_xtool/llclang/internal/parser"
)

/*
1) input: llclang parseIntermediateFile <filename> <mode>
* output: AST in JSON format to stdout
* error: non-zero exit code, and error message to stderr

2) input: XXX
* output: XXX
*/

type Command struct {
	Name     string
	Usage    string
	Short    string
	Run      func(args []string)
	Commands []*Command
}

var llclang *Command

func init() {
	llclang = &Command{
		Name:  "llclang",
		Usage: "llclang <command> [arguments]",
		Short: "llclang is a tool for parsing C/C++ source files",
		Commands: []*Command{
			{
				Name:  "parseIntermediateFile",
				Usage: "llclang parseIntermediateFile <filename> <mode>",
				Short: "Parse an intermediate file and output AST in JSON format",
				Run:   parser.RunParseIntermediateFile,
			},
			{
				Name:  "help",
				Usage: "llclang help <command>",
				Short: "Show help for commands",
				Run:   runHelp,
			},
		},
	}
}

func main() {
	if len(os.Args) < 2 {
		mainUsage()
		os.Exit(1)
	}

	cmdName := os.Args[1]
	args := os.Args[2:]

	for _, cmd := range llclang.Commands {
		if cmd.Name == cmdName {
			cmd.Run(args)
			return
		}
	}
}

func runHelp(args []string) {
	if len(args) == 0 {
		mainUsage()
		return
	}

	for _, cmd := range llclang.Commands {
		if cmd.Name == args[0] {
			fmt.Println(cmd.Usage)
			return
		}
	}
}

func mainUsage() {
	fmt.Println(llclang.Short)
	fmt.Println("Usage:", llclang.Usage)
	fmt.Println("Commands:")
	for _, cmd := range llclang.Commands {
		fmt.Printf("  %-30s %s\n", cmd.Name, cmd.Short)
	}
}
