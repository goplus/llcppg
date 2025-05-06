package main

import (
	"fmt"
	"os"

	"github.com/goplus/llcppg/_xtool/llcppsigfetch/parse"
	test "github.com/goplus/llcppg/_xtool/llcppsigfetch/parse/cvt_test"
	llcppg "github.com/goplus/llcppg/config"
)

func main() {
	TestPreprocess()
}

func TestPreprocess() {
	fmt.Println("=== TestPreProcess ===")
	combinedFile, err := os.Create("./combined.h")
	if err != nil {
		panic(err)
	}
	defer combinedFile.Close()
	defer os.Remove(combinedFile.Name())
	preprocessedFile, err := os.Create("./preprocessed.i")
	if err != nil {
		panic(err)
	}
	defer preprocessedFile.Close()
	defer os.Remove(preprocessedFile.Name())

	test.RunTestWithConfig(&parse.ParseConfig{
		Conf: &llcppg.Config{
			Include: []string{
				"main.h",
				"compat.h",
			},
			CFlags: "-I./hfile",
		},
		CombinedFile:     combinedFile.Name(),
		PreprocessedFile: preprocessedFile.Name(),
		OutputFile:       false,
	})
}
