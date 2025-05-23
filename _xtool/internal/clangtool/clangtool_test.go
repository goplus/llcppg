package clangtool_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/goplus/llcppg/_xtool/internal/clangtool"
)

func TestComposeIncludes(t *testing.T) {
	fmt.Println("=== Test ComposeIncludes ===")
	testCases := []struct {
		name   string
		files  []string
		expect string
	}{
		{
			name:  "One file",
			files: []string{"file1.h"},
			expect: `#include <file1.h>
`,
		},
		{
			name:  "Two files",
			files: []string{"file1.h", "file2.h"},
			expect: `#include <file1.h>
#include <file2.h>
`,
		},
		{
			name:   "Empty files",
			files:  []string{},
			expect: "",
		},
	}
	for _, tc := range testCases {
		outfile, err := os.CreateTemp("", "compose_*.h")
		if err != nil {
			t.Fatal(err)
		}

		err = clangtool.ComposeIncludes(tc.files, outfile.Name())
		if err != nil {
			t.Fatal(err)
		}
		content, err := os.ReadFile(outfile.Name())
		if err != nil {
			t.Fatal(err)
		}
		if string(content) != tc.expect {
			t.Fatalf("expect %s, but got %s", tc.expect, string(content))
		}
		outfile.Close()
		os.Remove(outfile.Name())
	}
}
